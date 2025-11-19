const mockService = require("../../../services/mockService");
const app = getApp();

Page({
  data: {
    employees: [],
    filteredEmployees: [],
    managers: [],
    searchText: ""
  },

  onLoad() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || user.role !== "admin") {
      wx.showToast({ title: "无权限访问", icon: "none" });
      setTimeout(() => {
        wx.navigateBack();
      }, 1000);
      return;
    }

    this.loadData();
  },

  async loadData() {
    try {
      const [usersRes, managersRes] = await Promise.all([
        mockService.fetchAllUsers(),
        mockService.fetchManagers()
      ]);

      const managers = managersRes.data || [];
      const employees = (usersRes.data || []).map(user => {
        const userManagers = managers.filter(m => user.managerIds && user.managerIds.includes(m.id));
        return {
          ...user,
          managers: userManagers
        };
      });

      this.setData({
        employees,
        filteredEmployees: employees,
        managers
      });
    } catch (err) {
      console.error("load data error", err);
      wx.showToast({ title: "加载数据失败", icon: "none" });
    }
  },

  handleSearch(e) {
    const searchText = e.detail.value.toLowerCase();
    this.setData({ searchText });

    const filtered = this.data.employees.filter(emp => 
      emp.name.toLowerCase().includes(searchText) ||
      emp.mobile.includes(searchText)
    );

    this.setData({ filteredEmployees: filtered });
  },

  goEmployeeDetail(e) {
    const { userId } = e.currentTarget.dataset;
    wx.navigateTo({
      url: `/pages/admin/employees/detail/index?userId=${userId}`
    });
  },

  goBack() {
    wx.navigateBack();
  }
});
