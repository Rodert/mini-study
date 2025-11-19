const mockService = require("../../../services/mockService");
const app = getApp();

Page({
  data: {
    categories: [],
    role: "",
    loading: false
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || !user.id) {
      wx.reLaunch({ url: "/pages/login/index" });
      return;
    }
    this.setData({ role: user.role });
    this.loadCategories(user.role);
  },

  async loadCategories(role) {
    this.setData({ loading: true });
    try {
      const res = await mockService.fetchCourseCategories(role);
      this.setData({ categories: res.data || [] });
    } catch (err) {
      console.error("load categories error", err);
    } finally {
      this.setData({ loading: false });
    }
  },

  handleSelectCategory(e) {
    const { item } = e.currentTarget.dataset;
    if (!item) return;
    wx.navigateTo({
      url: `/pages/learning/list/index?categoryId=${item.id}&name=${item.name}`
    });
  }
});

