const api = require("../../../../services/api");
const app = getApp();

Page({
  data: {
    user: {},
    categories: [],
    categoryNames: [],
    categoryIndex: -1,
    typeOptions: [
      { label: "文档", value: "doc" },
      { label: "视频", value: "video" },
      { label: "图文", value: "article" }
    ],
    typeIndex: 0,
    visibleRoleOptions: [
      { label: "全部可见", value: "both" },
      { label: "仅员工", value: "employee" },
      { label: "仅店长", value: "manager" }
    ],
    visibleRoleIndex: 0,
    statusOptions: [
      { label: "草稿", value: "draft" },
      { label: "发布", value: "published" }
    ],
    statusIndex: 0,
    uploadingFile: false,
    uploadingCover: false,
    coverPreviewUrl: "",
    form: {
      title: "",
      type: "doc",
      category_id: null,
      file_path: "",
      cover_url: "",
      summary: "",
      visible_roles: "both",
      status: "draft",
      duration_seconds: "",
      article_blocks: []
    },
    articleEditorText: "",
    isEdit: false,
    contentId: null
  },

  onLoad(options) {
    const id = options && options.id ? Number(options.id) : 0;
    if (id) {
      this.setData({
        isEdit: true,
        contentId: id
      });
      wx.setNavigationBarTitle({ title: "编辑内容" });
      this.loadContentDetail(id);
    } else {
      wx.setNavigationBarTitle({ title: "新增内容" });
    }
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || user.role !== "admin") {
      wx.showToast({ title: "仅管理员可访问", icon: "none" });
      setTimeout(() => wx.navigateBack(), 800);
      return;
    }
    this.setData({ user });
    this.loadCategories();
  },

  async loadCategories() {
    try {
      const res = await api.content.listCategories();
      if (res.code === 200) {
        const categories = res.data || [];
        let categoryIndex = this.data.categoryIndex;
        const form = this.data.form || {};
        if (form.category_id && categoryIndex < 0) {
          const idx = categories.findIndex((c) => c.id === form.category_id);
          if (idx !== -1) {
            categoryIndex = idx;
          }
        }
        this.setData({
          categories,
          categoryNames: categories.map((c) => c.name),
          categoryIndex
        });
      } else {
        wx.showToast({ title: res.message || "加载分类失败", icon: "none" });
      }
    } catch (err) {
      console.error("load categories error", err);
      wx.showToast({ title: "加载分类失败", icon: "none" });
    }
  },

  async loadContentDetail(id) {
    if (!id) return;
    try {
      const res = await api.content.getDetail(id);
      if (res.code !== 200 || !res.data) {
        wx.showToast({ title: res.message || "加载内容失败", icon: "none" });
        return;
      }
      const detail = res.data;
      const rawBlocks = detail.article_blocks || [];
      const articleBlocks = rawBlocks.map((block) => ({
        type: block.type,
        text: block.text || "",
        image_path: block.image_path || ""
      }));

      const form = {
        ...this.data.form,
        title: detail.title || "",
        type: detail.type || "doc",
        category_id: detail.category_id || null,
        file_path: detail.file_path || "",
        cover_url: detail.cover_url || "",
        summary: detail.summary || "",
        visible_roles: detail.visible_roles || "both",
        status: detail.status || "draft",
        duration_seconds: detail.duration_seconds || "",
        article_blocks: articleBlocks
      };

      const typeIndex = this.data.typeOptions.findIndex((o) => o.value === form.type);
      const visibleRoleIndex = this.data.visibleRoleOptions.findIndex(
        (o) => o.value === form.visible_roles
      );
      const statusIndex = this.data.statusOptions.findIndex((o) => o.value === form.status);

      this.setData({
        form,
        typeIndex: typeIndex >= 0 ? typeIndex : 0,
        visibleRoleIndex: visibleRoleIndex >= 0 ? visibleRoleIndex : 0,
        statusIndex: statusIndex >= 0 ? statusIndex : 0,
        coverPreviewUrl: form.cover_url ? api.buildFileUrl(form.cover_url) : ""
      });

      // 如果分类已经加载过，尝试设置分类下标
      const categories = this.data.categories || [];
      if (categories.length && form.category_id && this.data.categoryIndex < 0) {
        const idx = categories.findIndex((c) => c.id === form.category_id);
        if (idx !== -1) {
          this.setData({ categoryIndex: idx });
        }
      }
    } catch (err) {
      console.error("load content detail error", err);
      wx.showToast({ title: "加载内容失败", icon: "none" });
    }
  },

  handleInput(e) {
    const { field } = e.currentTarget.dataset;
    if (!field) return;
    const value = e.detail.value;
    this.setData({
      [`form.${field}`]: value
    });
    // 如果修改的是封面URL，更新预览URL
    if (field === "cover_url") {
      this.updateCoverPreview(value);
    }
  },

  handleArticleTextInput(e) {
    const value = e.detail.value || "";
    this.setData({ articleEditorText: value });
  },

  addArticleTextBlock() {
    const text = (this.data.articleEditorText || "").trim();
    if (!text) {
      wx.showToast({ title: "请输入文本内容", icon: "none" });
      return;
    }
    const blocks = this.data.form.article_blocks || [];
    blocks.push({ type: "text", text });
    this.setData({
      "form.article_blocks": blocks,
      articleEditorText: ""
    });
  },

  handleAddArticleImageBlock() {
    if (this.data.uploadingFile) return;
    wx.chooseImage({
      count: 1,
      sizeType: ["compressed", "original"],
      sourceType: ["album", "camera"],
      success: (res) => {
        const filePath = res.tempFilePaths && res.tempFilePaths[0];
        if (!filePath) {
          wx.showToast({ title: "图片路径无效", icon: "none" });
          return;
        }
        this.uploadArticleImageToServer(filePath);
      },
      fail: (err) => {
        if (err && err.errMsg !== "chooseImage:fail cancel") {
          wx.showToast({ title: "选择图片失败", icon: "none" });
        }
      }
    });
  },

  async uploadArticleImageToServer(filePath) {
    if (!filePath) return;
    let usablePath = filePath;
    try {
      usablePath = await this.ensureWxFilePath(filePath);
    } catch (err) {
      console.error("prepare article image path error", err);
      wx.showToast({ title: err.message || "图片不可用，请重新选择", icon: "none" });
      return;
    }

    this.setData({ uploadingFile: true });
    wx.showLoading({ title: "上传图片中..." });
    try {
      const res = await api.file.upload(usablePath);
      if (res.code === 200 && res.data && res.data.path) {
        const path = res.data.path;
        const blocks = this.data.form.article_blocks || [];
        blocks.push({ type: "image", image_path: path });
        this.setData({ "form.article_blocks": blocks });
        wx.showToast({ title: "图片块已添加", icon: "success" });
      } else {
        wx.showToast({ title: res.message || "上传失败", icon: "none" });
      }
    } catch (err) {
      console.error("upload article image error", err);
      wx.showToast({ title: err.message || "上传失败", icon: "none" });
    } finally {
      wx.hideLoading();
      this.setData({ uploadingFile: false });
    }
  },

  removeArticleBlock(e) {
    const index = Number(e.currentTarget.dataset.index);
    const blocks = this.data.form.article_blocks || [];
    if (index < 0 || index >= blocks.length) return;
    blocks.splice(index, 1);
    this.setData({ "form.article_blocks": blocks });
  },

  updateCoverPreview(coverUrl) {
    const previewUrl = coverUrl ? api.buildFileUrl(coverUrl) : "";
    this.setData({ coverPreviewUrl: previewUrl });
  },

  handleCategoryChange(e) {
    const index = Number(e.detail.value);
    const category = this.data.categories[index];
    if (!category) return;
    this.setData({
      categoryIndex: index,
      "form.category_id": category.id
    });
  },

  handleTypeChange(e) {
    const index = Number(e.detail.value);
    const option = this.data.typeOptions[index];
    if (!option) return;
    this.setData({
      typeIndex: index,
      "form.type": option.value
    });
  },

  handleRoleChange(e) {
    const index = Number(e.detail.value);
    const option = this.data.visibleRoleOptions[index];
    if (!option) return;
    this.setData({
      visibleRoleIndex: index,
      "form.visible_roles": option.value
    });
  },

  handleStatusChange(e) {
    const index = Number(e.detail.value);
    const option = this.data.statusOptions[index];
    if (!option) return;
    this.setData({
      statusIndex: index,
      "form.status": option.value
    });
  },

  validateForm() {
    const { form } = this.data;
    if (!form.title || !form.title.trim()) {
      wx.showToast({ title: "请输入标题", icon: "none" });
      return false;
    }
    if (!form.category_id) {
      wx.showToast({ title: "请选择分类", icon: "none" });
      return false;
    }
    if (form.type === "doc" || form.type === "video") {
      if (!form.file_path || !form.file_path.trim()) {
        wx.showToast({ title: "请输入文件地址", icon: "none" });
        return false;
      }
    }
    // 封面图改为可选，可以通过上传获取
    if (!form.summary || !form.summary.trim()) {
      wx.showToast({ title: "请输入内容简介", icon: "none" });
      return false;
    }
    if (form.type === "video") {
      const duration = Number(form.duration_seconds);
      if (!duration || duration <= 0) {
        wx.showToast({ title: "请输入视频时长(秒)", icon: "none" });
        return false;
      }
    }
    if (form.type === "article") {
      const blocks = form.article_blocks || [];
      if (!blocks.length) {
        wx.showToast({ title: "请至少添加一个图文内容块", icon: "none" });
        return false;
      }
    }
    return true;
  },

  handleUploadFile() {
    if (this.data.uploadingFile) return;
    if (this.data.form.type === "video") {
      this.chooseVideo();
    } else if (this.data.form.type === "doc") {
      this.chooseDocument();
    } else {
      // 图文类型不使用文件上传作为主体内容
      return;
    }
  },

  chooseDocument() {
    wx.chooseMessageFile({
      count: 1,
      type: "file",
      success: (res) => {
        const file = res.tempFiles && res.tempFiles[0];
        if (!file) return;
        const filePath = file.path || file.tempFilePath;
        if (!filePath) {
          wx.showToast({ title: "文件路径无效", icon: "none" });
          return;
        }
        this.uploadFileToServer(filePath);
      },
      fail: (err) => {
        if (err && err.errMsg !== "chooseMessageFile:fail cancel") {
          wx.showToast({ title: "选择文件失败", icon: "none" });
        }
      }
    });
  },

  chooseVideo() {
    wx.chooseMedia({
      count: 1,
      mediaType: ["video"],
      sourceType: ["album", "camera"],
      success: (res) => {
        const file = res.tempFiles && res.tempFiles[0];
        if (!file) return;
        const filePath = file.tempFilePath || file.path;
        if (!filePath) {
          wx.showToast({ title: "视频路径无效", icon: "none" });
          return;
        }
        const duration = Math.round(file.duration || 0);
        if (duration && !this.data.form.duration_seconds) {
          this.setData({
            "form.duration_seconds": duration
          });
        }
        this.uploadFileToServer(filePath);
      },
      fail: (err) => {
        if (err && err.errMsg !== "chooseMedia:fail cancel") {
          wx.showToast({ title: "选择视频失败", icon: "none" });
        }
      }
    });
  },

  async uploadFileToServer(filePath) {
    if (!filePath) return;
    let usablePath = filePath;
    try {
      usablePath = await this.ensureWxFilePath(filePath);
    } catch (err) {
      console.error("prepare file path error", err);
      wx.showToast({ title: err.message || "文件不可用，请重新选择", icon: "none" });
      return;
    }

    this.setData({ uploadingFile: true });
    wx.showLoading({ title: "上传中..." });
    try {
      const res = await api.file.upload(usablePath);
      if (res.code === 200 && res.data && res.data.path) {
        this.setData({
          "form.file_path": res.data.path
        });
        wx.showToast({ title: "上传成功", icon: "success" });
      } else {
        wx.showToast({ title: res.message || "上传失败", icon: "none" });
      }
    } catch (err) {
      console.error("upload file error", err);
      wx.showToast({ title: err.message || "上传失败", icon: "none" });
    } finally {
      wx.hideLoading();
      this.setData({ uploadingFile: false });
    }
  },

  ensureWxFilePath(filePath) {
    return new Promise((resolve, reject) => {
      if (!filePath) {
        reject(new Error("文件路径无效"));
        return;
      }

      const normalizedPath = filePath.trim();
      const isTempScheme =
        normalizedPath.startsWith("wxfile://") || normalizedPath.startsWith("http://tmp/");
      if (isTempScheme) {
        resolve(normalizedPath);
        return;
      }

      const fs = wx.getFileSystemManager();
      const dotIndex = normalizedPath.lastIndexOf(".");
      const ext = dotIndex !== -1 ? normalizedPath.slice(dotIndex) : "";
      const targetPath = `${wx.env.USER_DATA_PATH}/upload_${Date.now()}${ext}`;
      fs.readFile({
        filePath: normalizedPath,
        success: (readRes) => {
          fs.writeFile({
            filePath: targetPath,
            data: readRes.data,
            success: () => resolve(targetPath),
            fail: (err) => {
              reject(err || new Error("写入临时文件失败"));
            }
          });
        },
        fail: (err) => {
          reject(err || new Error("无法读取文件"));
        }
      });
    });
  },

  handleUploadCover() {
    if (this.data.uploadingCover) return;
    wx.chooseImage({
      count: 1,
      sizeType: ["compressed", "original"],
      sourceType: ["album", "camera"],
      success: (res) => {
        const filePath = res.tempFilePaths && res.tempFilePaths[0];
        if (!filePath) {
          wx.showToast({ title: "图片路径无效", icon: "none" });
          return;
        }
        this.uploadCoverToServer(filePath);
      },
      fail: (err) => {
        if (err && err.errMsg !== "chooseImage:fail cancel") {
          wx.showToast({ title: "选择图片失败", icon: "none" });
        }
      }
    });
  },

  async uploadCoverToServer(filePath) {
    if (!filePath) return;
    let usablePath = filePath;
    try {
      usablePath = await this.ensureWxFilePath(filePath);
    } catch (err) {
      console.error("prepare cover path error", err);
      wx.showToast({ title: err.message || "图片不可用，请重新选择", icon: "none" });
      return;
    }

    this.setData({ uploadingCover: true });
    wx.showLoading({ title: "上传封面中..." });
    try {
      const res = await api.file.upload(usablePath);
      if (res.code === 200 && res.data && res.data.path) {
        const coverPath = res.data.path;
        this.setData({
          "form.cover_url": coverPath
        });
        this.updateCoverPreview(coverPath);
        wx.showToast({ title: "封面上传成功", icon: "success" });
      } else {
        wx.showToast({ title: res.message || "上传失败", icon: "none" });
      }
    } catch (err) {
      console.error("upload cover error", err);
      wx.showToast({ title: err.message || "上传失败", icon: "none" });
    } finally {
      wx.hideLoading();
      this.setData({ uploadingCover: false });
    }
  },

  async handleSubmit() {
    if (!this.validateForm()) return;

    const { form } = this.data;
    const payload = {
      title: form.title.trim(),
      type: form.type,
      category_id: form.category_id,
      file_path: form.type === "article" ? "" : form.file_path.trim(),
      cover_url: form.cover_url ? form.cover_url.trim() : "",
      summary: form.summary.trim(),
      visible_roles: form.visible_roles,
      status: form.status,
      duration_seconds:
        form.type === "video" ? Number(form.duration_seconds) || 0 : undefined,
      article_blocks:
        form.type === "article" ? (form.article_blocks || []) : undefined
    };

    wx.showLoading({ title: "提交中..." });
    try {
      let res;
      if (this.data.isEdit && this.data.contentId) {
        res = await api.admin.updateContent(this.data.contentId, payload);
      } else {
        res = await api.admin.createContent(payload);
      }
      if (res.code === 200) {
        wx.showToast({ title: this.data.isEdit ? "保存成功" : "创建成功", icon: "success" });
        setTimeout(() => {
          wx.navigateBack();
        }, 600);
      } else {
        wx.showToast({ title: res.message || "创建失败", icon: "none" });
      }
    } catch (err) {
      console.error("create content error", err);
      wx.showToast({ title: "创建失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  },

  goBack() {
    wx.navigateBack();
  }
});


