const api = require("../../../services/api");

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

  formatDuration(seconds) {
    if (!seconds || seconds <= 0) {
      return "--";
    }
    const mins = Math.round(seconds / 60);
    if (mins < 60) {
      return `${mins} 分钟`;
    }
    const hours = Math.floor(mins / 60);
    const remains = mins % 60;
    return `${hours} 小时 ${remains} 分`;
  },

  async loadCourses() {
    if (!this.data.categoryId) return;
    this.setData({ loading: true });
    try {
      const res = await api.content.listPublished({
        category_id: this.data.categoryId
      });
      if (res.code === 200) {
        const courses = (res.data || []).map((item) => ({
          id: item.id,
          title: item.title,
          summary: item.summary,
          type: item.type,
          cover: item.cover_url ? api.buildFileUrl(item.cover_url) : "",
          duration: this.formatDuration(item.duration_seconds)
        }));
        this.setData({ courses });
      } else {
        wx.showToast({ title: res.message || "课程加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load courses error", err);
      wx.showToast({ title: "课程加载失败", icon: "none" });
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

