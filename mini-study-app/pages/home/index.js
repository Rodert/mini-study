const api = require("../../services/api");
const app = getApp();

Page({
  data: {
    user: {},
    banners: [],
    categories: []
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || !user.id) {
      wx.reLaunch({ url: "/pages/login/index" });
      return;
    }
    this.setData({ user });
    this.loadInitialData();
  },

  async loadInitialData() {
    if (this.data.user.role !== "admin") {
      await Promise.all([this.loadBanners(), this.loadCategories()]);
    }
  },

  async loadBanners() {
    try {
      const res = await api.banner.listVisible();
      if (res.code === 200) {
        const banners = (res.data || []).map((item) => ({
          id: item.id,
          title: item.title,
          cover: item.image_url,
          url: item.link_url,
          type: item.visible_roles
        }));
        this.setData({ banners });
      } else {
        wx.showToast({ title: res.message || "轮播加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("fetch banners error", err);
      wx.showToast({ title: "轮播加载失败", icon: "none" });
    }
  },

  async loadCategories() {
    try {
      const res = await api.content.listCategories();
      if (res.code === 200) {
        const categories = (res.data || []).map((item) => ({
          id: item.id,
          name: item.name,
          icon: "📖",
          count: item.count || 0
        }));
        this.setData({ categories });
      } else {
        wx.showToast({ title: res.message || "分类加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("fetch categories error", err);
      wx.showToast({ title: "分类加载失败", icon: "none" });
    }
  },

  reloadBanners() {
    this.loadBanners();
  },

  handleBannerTap(e) {
    const { item } = e.currentTarget.dataset;
    if (!item || !item.url) {
      console.warn("banner url missing", item);
      wx.showToast({ title: "链接异常，稍后重试", icon: "none" });
      return;
    }
    try {
      wx.navigateTo({
        url: `/pages/webview/index?url=${encodeURIComponent(item.url)}`
      });
    } catch (err) {
      console.error("navigateTo webview error", err);
      wx.showToast({ title: "打开失败", icon: "none" });
    }
  },

  goProgress() {
    wx.navigateTo({ url: "/pages/manager/progress/index" });
  },

  handleSelectCategory(e) {
    const { item } = e.currentTarget.dataset;
    if (!item) return;
    wx.navigateTo({
      url: `/pages/learning/list/index?categoryId=${item.id}&name=${item.name}`
    });
  },

  goProfile() {
    wx.navigateTo({ url: "/pages/profile/index" });
  },

  goUserManagement() {
    wx.navigateTo({ url: "/pages/manager/users/index" });
  },

  goEmployeesList() {
    wx.navigateTo({ url: "/pages/admin/employees/index" });
  },

  goBannerManagement() {
    wx.navigateTo({ url: "/pages/admin/banners/index" });
  },

  goContentCreate() {
    wx.navigateTo({ url: "/pages/admin/contents/add/index" });
  },

  goExamManagement() {
    wx.navigateTo({ url: "/pages/admin/exams/index" });
  },

  goExamList() {
    wx.navigateTo({ url: "/pages/exams/list/index" });
  }
});

