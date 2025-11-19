const mockService = require("../../../services/mockService");
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
      const [progressRes, employeesRes] = await Promise.all([
        mockService.fetchProgress(),
        mockService.fetchProgressEmployees()
      ]);
      const progress = (progressRes.data || []).map((item) => ({
        ...item,
        completionPercent: Math.round(item.completion * 100)
      }));
      const employees = (employeesRes.data || []).map((item) => ({
        ...item,
        percent: Math.round((item.completed / item.total) * 100)
      }));
      this.setData({ progress, employees });
    } catch (err) {
      console.error("load progress error", err);
    } finally {
      this.setData({ loading: false });
    }
  }
});

