const api = require("../../../services/api");
const app = getApp();

Page({
  data: {
    exams: [],
    loading: false
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || user.role !== "admin") {
      wx.showToast({ title: "仅管理员可访问", icon: "none" });
      setTimeout(() => wx.navigateBack(), 800);
      return;
    }
    this.loadExams();
  },

  async loadExams() {
    this.setData({ loading: true });
    try {
      const res = await api.admin.listExams();
      if (res.code === 200) {
        const exams = (res.data || []).map((item) => ({
          id: item.id,
          title: item.title,
          description: item.description,
          status: item.status || "draft",
          statusText: item.status === "published" ? "已发布" : item.status === "archived" ? "已归档" : "草稿",
          targetRole: item.target_role || "employee",
          targetRoleText: item.target_role === "employee" ? "员工" : item.target_role === "manager" ? "店长" : "全部",
          questionCount: item.question_count || 0,
          totalScore: item.total_score || 0,
          passScore: item.pass_score || 0,
          timeLimit: item.time_limit_minutes || 0
        }));
        this.setData({ exams });
      } else {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load exams error", err);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  handleCreate() {
    wx.navigateTo({ url: "/pages/admin/exams/add/index" });
  },

  handleEdit(e) {
    const { id } = e.currentTarget.dataset;
    if (!id) return;
    wx.navigateTo({ url: `/pages/admin/exams/add/index?id=${id}` });
  }
});

