const mockService = require("../../services/mockService");
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
      const res = await mockService.fetchBanners();
      this.setData({ banners: res.data || [] });
    } catch (err) {
      console.error("fetch banners error", err);
    }
  },

  async loadCategories() {
    try {
      const res = await mockService.fetchCourseCategories(this.data.user.role);
      this.setData({ categories: res.data || [] });
    } catch (err) {
      console.error("fetch categories error", err);
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
  }
});

