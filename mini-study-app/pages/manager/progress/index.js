const api = require("../../../services/api");
const app = getApp();

Page({
  data: {
    progress: [],
    employees: [],
    loading: false
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || user.role !== "manager") {
      wx.showToast({ title: "仅店长可访问", icon: "none" });
      wx.navigateBack({ delta: 1 });
      return;
    }
    this.loadData();
  },

  async loadData() {
    this.setData({ loading: true });
    try {
      const res = await api.exam.managerOverview();
      if (res.code === 200) {
        const overview = res.data || {};
        const progress = (overview.exam_progress || []).map((item) => ({
          id: item.exam_id,
          title: item.title,
          avgScore: (item.avg_score || 0).toFixed(1),
          attemptCount: item.attempt_count || 0,
          passRate: Math.round((item.pass_rate || 0) * 100)
      }));
        const employees = (overview.employees || []).map((item) => {
          const progressInfo = item.learning_progress || {};
          const latestExam = item.latest_exam ? {
            examId: item.latest_exam.exam_id,
            examTitle: item.latest_exam.exam_title,
            score: item.latest_exam.score,
            pass: item.latest_exam.pass,
            submittedAt: item.latest_exam.submitted_at
          } : null;
          return {
            id: item.employee_id,
            name: item.name,
            avatar: item.name ? item.name.charAt(0) : 'U',
            workNo: item.work_no,
            completed: progressInfo.completed || 0,
            total: progressInfo.total || 0,
            pending: progressInfo.pending || 0,
            percent: progressInfo.percent || 0,
            latestExam
          };
        });
      this.setData({ progress, employees });
      } else {
        wx.showToast({
          title: res.message || "加载进度失败",
          icon: "none"
        });
      }
    } catch (err) {
      console.error("load manager overview error", err);
      wx.showToast({ title: "加载进度失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  }
});

