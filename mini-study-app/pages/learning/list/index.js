const mockService = require("../../../services/mockService");

Page({
  data: {
    categoryId: null,
    categoryName: "",
    courses: [],
    loading: false
  },

  onLoad(options) {
    const { categoryId, name } = options;
    this.setData({
      categoryId: Number(categoryId),
      categoryName: name || ""
    });
    this.loadCourses();
  },

  async loadCourses() {
    if (!this.data.categoryId) return;
    this.setData({ loading: true });
    try {
      const res = await mockService.fetchCoursesByCategory(this.data.categoryId);
      this.setData({ courses: res.data || [] });
    } catch (err) {
      console.error("load courses error", err);
    } finally {
      this.setData({ loading: false });
    }
  },

  handleSelectCourse(e) {
    const { id } = e.currentTarget.dataset;
    if (!id) return;
    wx.navigateTo({
      url: `/pages/learning/detail/index?id=${id}`
    });
  }
});

