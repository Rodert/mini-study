const api = require("../../../services/api");

Page({
  data: {
    loading: false,
    exams: [],
    results: []
  },

  onShow() {
    this.refreshData();
  },

  async refreshData() {
    this.setData({ loading: true });
    await Promise.all([this.loadExams(), this.loadResults()]);
    this.setData({ loading: false });
  },

  async loadExams() {
    try {
      const res = await api.exam.listAvailable();
      if (res.code === 200) {
        const exams = (res.data || []).map((item) => ({
          id: item.id,
          title: item.title,
          description: item.description,
          totalScore: item.total_score,
          passScore: item.pass_score,
          questionCount: item.question_count,
          attemptStatus: item.attempt_status,
          timeLimit: item.time_limit_minutes,
          lastScore: item.last_score,
          lastPassed: item.last_passed,
          lastSubmittedText: item.last_submitted_at
            ? this.formatDateTime(item.last_submitted_at)
            : ""
        }));
        this.setData({ exams });
      } else {
        wx.showToast({
          title: res.message || "获取考试列表失败",
          icon: "none"
        });
      }
    } catch (err) {
      console.error("load exams error", err);
      wx.showToast({ title: "获取考试列表失败", icon: "none" });
    }
  },

  async loadResults() {
    try {
      const res = await api.exam.listMyResults();
      if (res.code === 200) {
        const results = (res.data || []).map((item) => ({
          id: item.attempt_id,
          title: item.exam_title,
          scoreText: `${item.score}/${item.total_score}`,
          pass: item.pass,
          submittedAt: this.formatDateTime(item.submitted_at)
        }));
        this.setData({ results });
      }
    } catch (err) {
      console.error("load exam results error", err);
    }
  },

  handleExamTap(event) {
    const { id } = event.currentTarget.dataset;
    if (!id) return;
    wx.navigateTo({
      url: `/pages/exams/detail/index?id=${id}`
    });
  },

  handleRetry(event) {
    const { id } = event.currentTarget.dataset;
    if (!id) return;
    wx.navigateTo({
      url: `/pages/exams/detail/index?id=${id}`
    });
  },

  formatDateTime(input) {
    if (!input) return "";
    const date = new Date(input);
    if (Number.isNaN(date.getTime())) {
      return "";
    }
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const d = String(date.getDate()).padStart(2, "0");
    const hh = String(date.getHours()).padStart(2, "0");
    const mm = String(date.getMinutes()).padStart(2, "0");
    return `${y}-${m}-${d} ${hh}:${mm}`;
  }
});



