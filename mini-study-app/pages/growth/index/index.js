const api = require("../../../services/api");
const app = getApp();

Page({
  data: {
    user: {},
    posts: [],
    keyword: "",
    loading: false
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || !user.id) {
      wx.reLaunch({ url: "/pages/login/index" });
      return;
    }
    this.setData({ user });
    this.loadPosts();
  },

  async loadPosts() {
    this.setData({ loading: true });
    try {
      const res = await api.growth.list({ keyword: this.data.keyword });
      if (res.code === 200) {
        const posts = (res.data || []).map((item) => this.transformPost(item));
        this.setData({ posts });
      } else {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load growth posts error", err);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
      wx.stopPullDownRefresh();
    }
  },

  transformPost(item) {
    const publisherName = item.publisher_name || "未知";
    const avatarText = publisherName ? publisherName.charAt(0) : "?";
    const images = (item.image_paths || []).map((p) => api.buildFileUrl(p));
    return {
      id: item.id,
      content: item.content || "",
      images,
      rawImagePaths: item.image_paths || [],
      status: item.status,
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

  onPullDownRefresh() {
    this.loadPosts();
  },

  handleSearchInput(e) {
    this.setData({ keyword: e.detail.value || "" });
  },

  handleSearchConfirm() {
    this.loadPosts();
  },

  goMyPosts() {
    wx.navigateTo({ url: "/pages/growth/mine/index" });
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
  }
});
