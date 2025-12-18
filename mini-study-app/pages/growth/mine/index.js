// pages/growth/mine/index.js
const api = require("../../../services/api");
const app = getApp();

Page({

  /**
   * 页面的初始数据
   */
  data: {
    user: {},
    posts: [],
    keyword: "",
    statusFilterOptions: [
      { label: "全部", value: "" },
      { label: "待审核", value: "pending" },
      { label: "已通过", value: "approved" },
      { label: "已拒绝", value: "rejected" }
    ],
    statusFilterIndex: 0,
    loading: false,
    form: {
      content: "",
      image_paths: []
    },
    imagePreviewList: [],
    uploading: false
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || !user.id) {
      wx.reLaunch({ url: "/pages/login/index" });
      return;
    }
    this.setData({ user });
    this.loadMyPosts();
  },

  onPullDownRefresh() {
    this.loadMyPosts();
  },

  async loadMyPosts() {
    this.setData({ loading: true });
    try {
      const statusOpt = this.data.statusFilterOptions[this.data.statusFilterIndex];
      const params = { keyword: this.data.keyword };
      if (statusOpt && statusOpt.value) {
        params.status = statusOpt.value;
      }
      const res = await api.growth.listMine(params);
      if (res.code === 200) {
        const posts = (res.data || []).map((item) => this.transformPost(item));
        this.setData({ posts });
      } else {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load my growth posts error", err);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
      wx.stopPullDownRefresh();
    }
  },

  transformPost(item) {
    const publisherName = item.publisher_name || "";
    const avatarText = publisherName ? publisherName.charAt(0) : "?";
    const images = (item.image_paths || []).map((p) => api.buildFileUrl(p));
    let statusText = "";
    if (item.status === "pending") statusText = "待审核";
    else if (item.status === "approved") statusText = "已通过";
    else if (item.status === "rejected") statusText = "已拒绝";

    return {
      id: item.id,
      content: item.content || "",
      images,
      rawImagePaths: item.image_paths || [],
      status: item.status,
      statusText,
      publisherName,
      publisherRole: item.publisher_role || "",
      avatarText,
      createdAt: item.created_at,
      createdAtText: this.formatDateTime(item.created_at)
    };
  },

  formatDateTime(isoString) {
    if (!isoString) return "";
    const d = new Date(isoString);
    if (Number.isNaN(d.getTime())) return isoString;
    const pad = (n) => (n < 10 ? `0${n}` : `${n}`);
    const y = d.getFullYear();
    const m = pad(d.getMonth() + 1);
    const day = pad(d.getDate());
    const h = pad(d.getHours());
    const mi = pad(d.getMinutes());
    return `${y}-${m}-${day} ${h}:${mi}`;
  },

  handleStatusChange(e) {
    const index = Number(e.detail.value);
    this.setData({ statusFilterIndex: index });
    this.loadMyPosts();
  },

  handleSearchInput(e) {
    this.setData({ keyword: e.detail.value || "" });
  },

  handleSearchConfirm() {
    this.loadMyPosts();
  },

  handleContentInput(e) {
    this.setData({ "form.content": e.detail.value });
  },

  handleChooseImages() {
    if (this.data.uploading) return;
    const current = this.data.form.image_paths || [];
    const remain = 9 - current.length;
    if (remain <= 0) {
      wx.showToast({ title: "最多9张图片", icon: "none" });
      return;
    }
    wx.chooseImage({
      count: remain,
      sizeType: ["compressed", "original"],
      sourceType: ["album", "camera"],
      success: (res) => {
        const paths = res.tempFilePaths || [];
        if (!paths.length) return;
        this.uploadSelectedImages(paths);
      }
    });
  },

  async uploadSelectedImages(paths) {
    if (!paths || !paths.length) return;
    this.setData({ uploading: true });
    wx.showLoading({ title: "上传图片中..." });
    const imagePaths = [...(this.data.form.image_paths || [])];
    const previews = [...(this.data.imagePreviewList || [])];

    for (const filePath of paths) {
      try {
        const res = await api.file.upload(filePath);
        if (res.code === 200 && res.data && res.data.path) {
          const p = res.data.path;
          imagePaths.push(p);
          previews.push(api.buildFileUrl(p));
        } else {
          wx.showToast({ title: res.message || "上传失败", icon: "none" });
        }
      } catch (err) {
        console.error("upload image error", err);
        wx.showToast({ title: "上传失败", icon: "none" });
      }
    }

    wx.hideLoading();
    this.setData({
      uploading: false,
      "form.image_paths": imagePaths,
      imagePreviewList: previews
    });
  },

  async handleSubmit() {
    const { user, form } = this.data;
    if (user.role !== "manager") {
      wx.showToast({ title: "仅店长可发布", icon: "none" });
      return;
    }
    if (!form.content || !form.content.trim()) {
      wx.showToast({ title: "请输入内容", icon: "none" });
      return;
    }

    wx.showLoading({ title: "发布中..." });
    try {
      const payload = {
        content: form.content.trim(),
        image_paths: form.image_paths || []
      };
      const res = await api.growth.create(payload);
      if (res.code === 200) {
        wx.showToast({ title: "发布成功", icon: "success" });
        this.resetForm();
        this.loadMyPosts();
      } else {
        wx.showToast({ title: res.message || "发布失败", icon: "none" });
      }
    } catch (err) {
      console.error("create growth post error", err);
      wx.showToast({ title: "发布失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  },

  resetForm() {
    this.setData({
      form: { content: "", image_paths: [] },
      imagePreviewList: [],
      uploading: false
    });
  },

  async handleDelete(e) {
    const id = Number(e.currentTarget.dataset.id);
    if (!id) return;
    wx.showModal({
      title: "删除动态",
      content: "确定要删除这条成长圈动态吗？",
      success: async (res) => {
        if (!res.confirm) return;
        try {
          const resp = await api.growth.delete(id);
          if (resp.code === 200) {
            wx.showToast({ title: "已删除", icon: "success" });
            this.loadMyPosts();
          } else {
            wx.showToast({ title: resp.message || "删除失败", icon: "none" });
          }
        } catch (err) {
          console.error("delete growth post error", err);
          wx.showToast({ title: "删除失败", icon: "none" });
        }
      }
    });
  },

  async handleCopyAndDownload(e) {
    const index = Number(e.currentTarget.dataset.index);
    const post = this.data.posts[index];
    if (!post) return;

    const text = post.content || "";
    wx.setClipboardData({
      data: text,
      success: () => {
        wx.showToast({ title: "文案已复制", icon: "none" });
        if (post.images && post.images.length) {
          this.saveImages(post.images);
        } else {
          wx.showToast({ title: "此动态无图片", icon: "none" });
        }
      }
    });
  },

  saveImages(imageUrls) {
    if (!imageUrls || !imageUrls.length) return;
    wx.showLoading({ title: "保存图片中..." });
    let successCount = 0;
    let failCount = 0;
    const total = imageUrls.length;

    const saveNext = (index) => {
      if (index >= total) {
        wx.hideLoading();
        let msg = "";
        if (successCount > 0) {
          msg = `已保存${successCount}张图片`;
        } else {
          msg = "未能保存图片";
        }
        if (failCount > 0) {
          msg += `，${failCount}张保存失败`;
        }
        wx.showToast({ title: msg, icon: "none" });
        return;
      }

      const url = imageUrls[index];
      wx.getImageInfo({
        src: url,
        success: (res) => {
          wx.saveImageToPhotosAlbum({
            filePath: res.path,
            success: () => {
              successCount += 1;
            },
            fail: (err) => {
              console.error("saveImageToPhotosAlbum fail", err);
              failCount += 1;
            },
            complete: () => {
              saveNext(index + 1);
            }
          });
        },
        fail: (err) => {
          console.error("getImageInfo fail", err);
          failCount += 1;
          saveNext(index + 1);
        }
      });
    };

    saveNext(0);
  },

  handlePreviewImage(e) {
    const postIndex = Number(e.currentTarget.dataset.postIndex);
    const imageIndex = Number(e.currentTarget.dataset.imageIndex);
    const post = this.data.posts[postIndex];
    if (!post || !post.images || !post.images.length) return;

    wx.previewImage({
      current: post.images[imageIndex] || post.images[0],
      urls: post.images
    });
  },

  goBack() {
    wx.navigateBack();
  }
});