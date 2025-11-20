const api = require("../../../services/api");
const app = getApp();

Page({
  data: {
    employees: [],
    filteredEmployees: [],
    managers: [],
    searchText: "",
    showMenu: false
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

  onShow() {
    // 从新增页面返回时刷新数据
    this.loadData();
  },

  async loadData() {
    try {
      const res = await api.admin.listUsers();
      if (res.code !== 200) {
        wx.showToast({ title: res.message || "加载数据失败", icon: "none" });
        return;
      }

      const employees = (res.data || []).map((user) => {
        const managers = (user.managers || []).map((manager) => ({
          id: manager.id,
          name: manager.name || manager.work_no,
          workNo: manager.work_no,
          phone: manager.phone
        }));
        return {
          id: user.id,
          name: user.name,
          username: user.work_no,
          mobile: user.phone,
          store: user.store || "—",
          role: user.role,
          managerIds: user.manager_ids || [],
          managers
        };
      });

      this.setData({
        employees,
        filteredEmployees: employees
      });
    } catch (err) {
      console.error("load data error", err);
      wx.showToast({ title: "加载数据失败", icon: "none" });
    }
  },

  handleSearch(e) {
    const searchText = e.detail.value.toLowerCase();
    this.setData({ searchText });

    const filtered = this.data.employees.filter(emp => {
      const name = (emp.name || "").toLowerCase();
      const mobile = emp.mobile || "";
      const workNo = (emp.username || "").toLowerCase();
      return (
        name.includes(searchText) ||
        mobile.includes(searchText) ||
        workNo.includes(searchText)
      );
    });

    this.setData({ filteredEmployees: filtered });
  },

  goEmployeeDetail(e) {
    const { userId } = e.currentTarget.dataset;
    wx.navigateTo({
      url: `/pages/admin/employees/detail/index?userId=${userId}`
    });
  },

  showAddMenu() {
    this.setData({ showMenu: true });
  },

  hideAddMenu() {
    this.setData({ showMenu: false });
  },

  stopPropagation() {
    // 阻止事件冒泡
  },

  goAddEmployee() {
    this.setData({ showMenu: false });
    wx.navigateTo({
      url: "/pages/admin/employees/add/index"
    });
  },

  goAddManager() {
    this.setData({ showMenu: false });
    wx.navigateTo({
      url: "/pages/admin/managers/add/index"
    });
  },

  goBack() {
    wx.navigateBack();
  }
});
