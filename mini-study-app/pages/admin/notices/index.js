const api = require("../../../services/api");
const app = getApp();

function createDefaultForm() {
  return {
    title: "",
    content: "",
    image_url: "",
    status: true
  };
}

Page({
  data: {
    user: {},
    notices: [],
    loading: false,
    showForm: false,
    editingId: null,
    uploadingImage: false,
    imagePreviewUrl: "",
    form: createDefaultForm(),
    // 公告确认列表
    showConfirmList: false,
    confirmList: [],
    confirmNoticeTitle: ""
  },

  // 时间格式化：YYYY-MM-DD HH:mm:ss
  formatDateTime(str) {
    if (!str) return "";
    const d = new Date(str);
    if (isNaN(d.getTime())) return str;
    const pad = (n) => (n < 10 ? "0" + n : "" + n);
    const Y = d.getFullYear();
    const M = pad(d.getMonth() + 1);
    const D = pad(d.getDate());
    const h = pad(d.getHours());
    const m = pad(d.getMinutes());
    const s = pad(d.getSeconds());
    return `${Y}-${M}-${D} ${h}:${m}:${s}`;
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || !user.id) {
      wx.reLaunch({ url: "/pages/login/index" });
      return;
    }
    this.setData({ user });
    this.loadNotices();
  },

  async loadNotices() {
    this.setData({ loading: true });
    try {
      // 所有角色共用公告列表接口
      const res = await api.notice.list();
      if (res.code === 200) {
        const notices = (res.data || [])
          .map((item) => ({
            ...item,
            image_preview_url: item.image_url ? api.buildFileUrl(item.image_url) : ""
          }))
          .sort((a, b) => {
            const ta = new Date(a.start_at || a.created_at || 0).getTime();
            const tb = new Date(b.start_at || b.created_at || 0).getTime();
            return tb - ta; // 倒序：时间新的在前
          });
        this.setData({ notices });
      } else {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load notices error", err);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  openCreateForm() {
    const form = createDefaultForm();
    this.setData({
      showForm: true,
      editingId: null,
      form,
      imagePreviewUrl: ""
    });
  },

  openEditForm(e) {
    const index = Number(e.currentTarget.dataset.index);
    const notice = this.data.notices[index];
    if (!notice) return;

    const form = {
      title: notice.title || "",
      content: notice.content || "",
      image_url: notice.image_url || "",
      status: notice.status
    };

    this.setData({
      showForm: true,
      editingId: notice.id,
      form
    });
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

  validateForm() {
    const { form } = this.data;
    if (!form.title || !form.title.trim()) {
      wx.showToast({ title: "请输入标题", icon: "none" });
      return false;
    }
    if ((!form.content || !form.content.trim()) && (!form.image_url || !form.image_url.trim())) {
      wx.showToast({ title: "请输入内容或上传图片", icon: "none" });
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
      content: form.content ? form.content.trim() : "",
      image_url: form.image_url ? form.image_url.trim() : "",
      status: form.status
    };

    wx.showLoading({ title: editingId ? "保存中..." : "创建中..." });
    try {
      let res;
      if (editingId) {
        res = await api.admin.updateNotice(editingId, payload);
      } else {
        res = await api.admin.createNotice(payload);
      }
      if (res.code === 200) {
        wx.showToast({ title: "操作成功", icon: "success" });
        this.setData({ showForm: false });
        this.loadNotices();
      } else {
        wx.showToast({ title: res.message || "保存失败", icon: "none" });
      }
    } catch (err) {
      console.error("save notice error", err);
      wx.showToast({ title: "保存失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  },

  // 员工/店长：确认公告已读
  async handleConfirmNotice(e) {
    const id = Number(e.currentTarget.dataset.id);
    const index = Number(e.currentTarget.dataset.index);
    if (!id && id !== 0) return;

    wx.showLoading({ title: "提交中..." });
    try {
      const res = await api.notice.confirm(id);
      if (res.code === 200) {
        wx.showToast({ title: "已确认", icon: "success" });
        const { notices } = this.data;
        const item = notices[index];
        if (item) {
          const confirmedAt = (res.data && res.data.confirmed_at) || new Date().toISOString();
          item.confirmed = true;
          item.confirmed_at = confirmedAt;
          this.setData({ [`notices[${index}]`]: item });
        }
      } else {
        wx.showToast({ title: res.message || "操作失败", icon: "none" });
      }
    } catch (err) {
      console.error("confirm notice error", err);
      wx.showToast({ title: "操作失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  },

  // 管理员：查看某条公告的确认情况
  async viewConfirmations(e) {
    const id = Number(e.currentTarget.dataset.id);
    if (!id && id !== 0) return;

    wx.showLoading({ title: "加载中..." });
    try {
      const res = await api.admin.getNoticeConfirmations(id);
      if (res.code !== 200) {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
        return;
      }
      const rawList = res.data || [];

      const list = rawList.map((item) => ({
        ...item,
        display_time: this.formatDateTime(item.confirmed_at)
      }));

      const notice = (this.data.notices || []).find((n) => n.id === id);
      this.setData({
        confirmList: list,
        confirmNoticeTitle: notice ? notice.title : "",
        showConfirmList: true
      });
    } catch (err) {
      console.error("get confirmations error", err);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  },

  closeConfirmList() {
    this.setData({ showConfirmList: false });
  },

  async toggleNoticeStatus(e) {
    const index = Number(e.currentTarget.dataset.index);
    const notice = this.data.notices[index];
    if (!notice) return;

    wx.showLoading({ title: "处理中..." });
    try {
      const res = await api.admin.updateNotice(notice.id, {
        status: !notice.status
      });
      if (res.code === 200) {
        wx.showToast({ title: "已更新状态", icon: "success" });
        this.loadNotices();
      } else {
        wx.showToast({ title: res.message || "操作失败", icon: "none" });
      }
    } catch (err) {
      console.error("toggle notice error", err);
      wx.showToast({ title: "操作失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  },

  goBack() {
    wx.navigateBack();
  }
});
