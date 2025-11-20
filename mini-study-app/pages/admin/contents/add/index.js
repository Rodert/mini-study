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
      { label: "视频", value: "video" }
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
    form: {
      title: "",
      type: "doc",
      category_id: null,
      file_path: "",
      cover_url: "",
      summary: "",
      visible_roles: "both",
      status: "draft",
      duration_seconds: ""
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
        this.setData({
          categories,
          categoryNames: categories.map((c) => c.name)
        });
      } else {
        wx.showToast({ title: res.message || "加载分类失败", icon: "none" });
      }
    } catch (err) {
      console.error("load categories error", err);
      wx.showToast({ title: "加载分类失败", icon: "none" });
    }
  },

  handleInput(e) {
    const { field } = e.currentTarget.dataset;
    if (!field) return;
    this.setData({
      [`form.${field}`]: e.detail.value
    });
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
    if (!form.file_path || !form.file_path.trim()) {
      wx.showToast({ title: "请输入文件地址", icon: "none" });
      return false;
    }
    if (!form.cover_url || !form.cover_url.trim()) {
      wx.showToast({ title: "请输入封面链接", icon: "none" });
      return false;
    }
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
    return true;
  },

  handleUploadFile() {
    if (this.data.uploadingFile) return;
    if (this.data.form.type === "video") {
      this.chooseVideo();
    } else {
      this.chooseDocument();
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

  async handleSubmit() {
    if (!this.validateForm()) return;

    const { form } = this.data;
    const payload = {
      title: form.title.trim(),
      type: form.type,
      category_id: form.category_id,
      file_path: form.file_path.trim(),
      cover_url: form.cover_url.trim(),
      summary: form.summary.trim(),
      visible_roles: form.visible_roles,
      status: form.status,
      duration_seconds:
        form.type === "video" ? Number(form.duration_seconds) || 0 : undefined
    };

    wx.showLoading({ title: "提交中..." });
    try {
      const res = await api.admin.createContent(payload);
      if (res.code === 200) {
        wx.showToast({ title: "创建成功", icon: "success" });
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


