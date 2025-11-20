const api = require("../../../services/api");
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
    this.loadCategories();
  },

  async loadCategories() {
    this.setData({ loading: true });
    try {
      const res = await api.content.listCategories();
      if (res.code === 200) {
        const categories = (res.data || []).map((item) => ({
          id: item.id,
          name: item.name,
          icon: "📘"
        }));
        this.setData({ categories });
      } else {
        wx.showToast({ title: res.message || "分类加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load categories error", err);
      wx.showToast({ title: "分类加载失败", icon: "none" });
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

