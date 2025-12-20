const api = require("../../../services/api");
const app = getApp();

Page({
  data: {
    loading: false,
    filters: {
      role: "all",
      keyword: "",
      page: 1,
      pageSize: 20
    },
    roles: [
      { label: "全部角色", value: "all" },
      { label: "员工", value: "employee" },
      { label: "店长", value: "manager" },
      { label: "管理员", value: "admin" }
    ],
    examProgress: [],
    users: [],
    pagination: {
      page: 1,
      pageSize: 20,
      total: 0,
      maxPage: 1
    },
    roleIndex: 0
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || user.role !== "admin") {
      wx.showToast({ title: "仅管理员可访问", icon: "none" });
      wx.navigateBack({ delta: 1 });
      return;
    }
    this.initPage();
  },

  async initPage() {
    this.loadData();
  },

  async loadData() {
    const { filters } = this.data;
    this.setData({ loading: true });
    try {
      const params = {
        role: filters.role === "all" ? "" : filters.role,
        keyword: filters.keyword || "",
        page: filters.page,
        page_size: filters.pageSize
      };
      const res = await api.admin.examOverview(params);
      if (res.code === 200) {
        const data = res.data || {};
        const examProgress = (data.exam_progress || []).map((item) => ({
          id: item.exam_id,
          title: item.title,
          avgScore: (item.avg_score || 0).toFixed(1),
          attemptCount: item.attempt_count || 0,
          passRate: Math.round((item.pass_rate || 0) * 100)
        }));
        const users = (data.users || []).map((item) => {
          const progress = item.learning_progress || {};
          const latestExam = item.latest_exam
            ? {
                examId: item.latest_exam.exam_id,
                examTitle: item.latest_exam.exam_title,
                score: item.latest_exam.score,
                pass: item.latest_exam.pass,
                submittedAt: item.latest_exam.submitted_at
              }
            : null;
          return {
            id: item.user_id,
            name: item.name,
            avatar: item.name ? item.name.charAt(0) : "U",
            workNo: item.work_no,
            role: item.role,
            completed: progress.completed || 0,
            total: progress.total || 0,
            pending: progress.pending || 0,
            percent: progress.percent || 0,
            latestExam
          };
        });
        const pagination = data.pagination || {};
        const page = pagination.page || filters.page;
        const pageSize = pagination.page_size || filters.pageSize;
        const total = pagination.total || 0;
        const maxPage = total > 0 ? Math.ceil(total / pageSize) : 1;
        this.setData({
          examProgress,
          users,
          pagination: {
            page,
            pageSize,
            total,
            maxPage
          }
        });
      } else {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
      }
    } catch (e) {
      console.error("load admin exam overview error", e);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  handleRoleChange(e) {
    const idx = Number(e.detail.value || 0);
    const roleItem = this.data.roles[idx] || this.data.roles[0];
    this.setData({
      roleIndex: idx,
      "filters.role": roleItem.value,
      "filters.page": 1
    });
    this.loadData();
  },

  handleKeywordInput(e) {
    this.setData({ "filters.keyword": e.detail.value || "" });
  },

  handleSearch() {
    this.setData({ "filters.page": 1 });
    this.loadData();
  },

  handlePrevPage() {
    const { page } = this.data.pagination;
    if (page <= 1) return;
    this.setData({ "filters.page": page - 1 });
    this.loadData();
  },

  handleNextPage() {
    const { page, maxPage } = this.data.pagination;
    if (page >= maxPage) return;
    this.setData({ "filters.page": page + 1 });
    this.loadData();
  }
});
