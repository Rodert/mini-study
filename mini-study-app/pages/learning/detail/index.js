const mockService = require("../../../services/mockService");

Page({
  data: {
    courseId: null,
    course: {},
    finished: false
  },

  onLoad(options) {
    const { id } = options;
    this.setData({ courseId: Number(id) });
    this.loadCourse();
  },

  async loadCourse() {
    if (!this.data.courseId) return;
    try {
      const res = await mockService.fetchCourseDetail(this.data.courseId);
      this.setData({
        course: res.data || {},
        finished: Math.random() > 0.5 // mock completion
      });
    } catch (err) {
      console.error("load course detail error", err);
    }
  }
});

