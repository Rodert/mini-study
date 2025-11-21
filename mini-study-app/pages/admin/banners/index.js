const api = require("../../../services/api");
const app = getApp();

function createDefaultForm() {
  return {
    title: "",
    image_url: "",
    link_url: "",
    visible_roles: "both",
    sort_order: 1,
    status: true
  };
}

Page({
  data: {
    user: {},
    banners: [],
    loading: false,
    showForm: false,
    editingId: null,
    uploadingImage: false,
    imagePreviewUrl: "",
    form: createDefaultForm(),
    visibleRoleOptions: [
      { label: "全部可见", value: "both" },
      { label: "仅员工可见", value: "employee" },
      { label: "仅店长可见", value: "manager" }
    ],
    visibleRoleIndex: 0,
    roleTextMap: {
      both: "全部可见",
      employee: "仅员工",
      manager: "仅店长"
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
    this.loadBanners();
  },

  async loadBanners() {
    this.setData({ loading: true });
    try {
      const res = await api.admin.listBanners();
      if (res.code === 200) {
        const banners = (res.data || []).map((item) => ({
          ...item,
          image_preview_url: item.image_url ? api.buildFileUrl(item.image_url) : ""
        }));
        this.setData({ banners });
      } else {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load banners error", err);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  getRoleIndex(value) {
    const idx = this.data.visibleRoleOptions.findIndex(
      (item) => item.value === value
    );
    return idx >= 0 ? idx : 0;
  },

  openCreateForm() {
    const form = createDefaultForm();
    this.setData({
      showForm: true,
      editingId: null,
      form,
      visibleRoleIndex: this.getRoleIndex(form.visible_roles),
      imagePreviewUrl: ""
    });
  },

  openEditForm(e) {
    const index = Number(e.currentTarget.dataset.index);
    const banner = this.data.banners[index];
    if (!banner) return;

    const form = {
      title: banner.title || "",
      image_url: banner.image_url || "",
      link_url: banner.link_url || "",
      visible_roles: banner.visible_roles || "both",
      sort_order: banner.sort_order || 1,
      status: banner.status
    };

    this.setData({
      showForm: true,
      editingId: banner.id,
      form,
      visibleRoleIndex: this.getRoleIndex(form.visible_roles)
    });
    // 更新预览URL
    this.updateImagePreview(form.image_url);
  },

  closeForm() {
    this.setData({ showForm: false });
  },

  handleInput(e) {
    const { field } = e.currentTarget.dataset;
    const value = e.detail.value;
    if (!field) return;
    this.setData({
      [`form.${field}`]: value
    });
    // 如果修改的是图片URL，更新预览URL
    if (field === "image_url") {
      this.updateImagePreview(value);
    }
  },

  updateImagePreview(imageUrl) {
    const previewUrl = imageUrl ? api.buildFileUrl(imageUrl) : "";
    this.setData({ imagePreviewUrl: previewUrl });
  },

  handleUploadImage() {
    if (this.data.uploadingImage) return;
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
        this.uploadImageToServer(filePath);
      },
      fail: (err) => {
        if (err && err.errMsg !== "chooseImage:fail cancel") {
          wx.showToast({ title: "选择图片失败", icon: "none" });
        }
      }
    });
  },

  async uploadImageToServer(filePath) {
    if (!filePath) return;
    let usablePath = filePath;
    try {
      usablePath = await this.ensureWxFilePath(filePath);
    } catch (err) {
      console.error("prepare image path error", err);
      wx.showToast({ title: err.message || "图片不可用，请重新选择", icon: "none" });
      return;
    }

    this.setData({ uploadingImage: true });
    wx.showLoading({ title: "上传图片中..." });
    try {
      const res = await api.file.upload(usablePath);
      if (res.code === 200 && res.data && res.data.path) {
        const imagePath = res.data.path;
        this.setData({
          "form.image_url": imagePath
        });
        this.updateImagePreview(imagePath);
        wx.showToast({ title: "图片上传成功", icon: "success" });
      } else {
        wx.showToast({ title: res.message || "上传失败", icon: "none" });
      }
    } catch (err) {
      console.error("upload image error", err);
      wx.showToast({ title: err.message || "上传失败", icon: "none" });
    } finally {
      wx.hideLoading();
      this.setData({ uploadingImage: false });
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

  handleSwitchChange(e) {
    this.setData({
      "form.status": e.detail.value
    });
  },

  handleRolePicker(e) {
    const index = Number(e.detail.value);
    const option = this.data.visibleRoleOptions[index];
    if (!option) return;
    this.setData({
      visibleRoleIndex: index,
      "form.visible_roles": option.value
    });
  },

  validateForm() {
    const { form } = this.data;
    if (!form.title || !form.title.trim()) {
      wx.showToast({ title: "请输入标题", icon: "none" });
      return false;
    }
    if (!form.image_url || !form.image_url.trim()) {
      wx.showToast({ title: "请上传图片", icon: "none" });
      return false;
    }
    if (!form.link_url || !form.link_url.trim()) {
      wx.showToast({ title: "请输入跳转链接", icon: "none" });
      return false;
    }
    return true;
  },

  async handleSubmit() {
    if (!this.validateForm()) {
      return;
    }
    const { form, editingId } = this.data;
    const payload = {
      title: form.title.trim(),
      image_url: form.image_url.trim(),
      link_url: form.link_url.trim(),
      visible_roles: form.visible_roles,
      sort_order: Number(form.sort_order) || 0,
      status: form.status
    };

    wx.showLoading({ title: editingId ? "保存中..." : "创建中..." });
    try {
      let res;
      if (editingId) {
        res = await api.admin.updateBanner(editingId, payload);
      } else {
        res = await api.admin.createBanner(payload);
      }
      if (res.code === 200) {
        wx.showToast({ title: "操作成功", icon: "success" });
        this.setData({ showForm: false });
        this.loadBanners();
      } else {
        wx.showToast({ title: res.message || "保存失败", icon: "none" });
      }
    } catch (err) {
      console.error("save banner error", err);
      wx.showToast({ title: "保存失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  },

  async toggleBannerStatus(e) {
    const index = Number(e.currentTarget.dataset.index);
    const banner = this.data.banners[index];
    if (!banner) return;

    wx.showLoading({ title: "处理中..." });
    try {
      const res = await api.admin.updateBanner(banner.id, {
        status: !banner.status
      });
      if (res.code === 200) {
        wx.showToast({ title: "已更新状态", icon: "success" });
        this.loadBanners();
      } else {
        wx.showToast({ title: res.message || "操作失败", icon: "none" });
      }
    } catch (err) {
      console.error("toggle banner error", err);
      wx.showToast({ title: "操作失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  },

  goBack() {
    wx.navigateBack();
  }
});


