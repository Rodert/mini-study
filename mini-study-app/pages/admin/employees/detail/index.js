const mockService = require("../../../../services/mockService");
const app = getApp();

Page({
  data: {
    userId: null,
    employee: {},
    managers: [],
    selectedManagerIds: [],
    originalRole: "",
    originalManagerIds: []
  },

  onLoad(options) {
    const admin = app.globalData.user || wx.getStorageSync("user");
    if (!admin || admin.role !== "admin") {
      wx.showToast({ title: "无权限访问", icon: "none" });
      setTimeout(() => {
        wx.navigateBack();
      }, 1000);
      return;
    }

    const { userId } = options;
    if (!userId) {
      wx.showToast({ title: "参数错误", icon: "none" });
      wx.navigateBack();
      return;
    }

    this.setData({ userId: parseInt(userId) });
    this.loadData();
  },

  async loadData() {
    try {
      const [usersRes, managersRes] = await Promise.all([
        mockService.fetchAllUsers(),
        mockService.fetchManagers()
      ]);

      const employee = usersRes.data.find(u => u.id === this.data.userId);
      if (!employee) {
        wx.showToast({ title: "员工不存在", icon: "none" });
        wx.navigateBack();
        return;
      }

      const selectedManagerIds = employee.managerIds || [];
      this.setData({
        employee,
        managers: managersRes.data || [],
        selectedManagerIds,
        originalRole: employee.role,
        originalManagerIds: [...selectedManagerIds]
      });
    } catch (err) {
      console.error("load data error", err);
      wx.showToast({ title: "加载数据失败", icon: "none" });
    }
  },

  selectRole(e) {
    const { role } = e.currentTarget.dataset;
    const employee = { ...this.data.employee, role };
    
    // 如果切换为非员工角色，清空店长绑定
    if (role !== "employee") {
      this.setData({
        employee,
        selectedManagerIds: []
      });
    } else {
      this.setData({ employee });
    }
  },

  handleManagerSelect(e) {
    const { managerId } = e.currentTarget.dataset;
    const { selectedManagerIds } = this.data;
    const index = selectedManagerIds.indexOf(managerId);

    if (index > -1) {
      selectedManagerIds.splice(index, 1);
    } else {
      selectedManagerIds.push(managerId);
    }

    this.setData({ selectedManagerIds });
  },

  async handleSave() {
    const { employee, selectedManagerIds, originalRole, originalManagerIds } = this.data;

    // 检查是否有修改
    const roleChanged = employee.role !== originalRole;
    const managerChanged = JSON.stringify(selectedManagerIds) !== JSON.stringify(originalManagerIds);

    if (!roleChanged && !managerChanged) {
      wx.showToast({ title: "没有修改任何信息", icon: "none" });
      return;
    }

    // 如果是员工，验证是否选择了店长
    if (employee.role === "employee" && selectedManagerIds.length === 0) {
      wx.showToast({ title: "员工必须绑定至少一个店长", icon: "none" });
      return;
    }

    wx.showLoading({ title: "保存中..." });

    try {
      let response;

      if (roleChanged && managerChanged) {
        // 同时修改角色和店长绑定
        await mockService.updateUserRole(this.data.userId, employee.role);
        response = await mockService.updateUserManagers(this.data.userId, selectedManagerIds);
      } else if (roleChanged) {
        // 只修改角色
        response = await mockService.updateUserRole(this.data.userId, employee.role);
      } else {
        // 只修改店长绑定
        response = await mockService.updateUserManagers(this.data.userId, selectedManagerIds);
      }

      if (response.success) {
        wx.hideLoading();
        wx.showToast({ title: "员工已更新", icon: "success" });
        setTimeout(() => {
          wx.navigateBack();
        }, 500);
      } else {
        wx.hideLoading();
        wx.showToast({ title: response.message || "更新失败", icon: "none" });
      }
    } catch (err) {
      wx.hideLoading();
      wx.showToast({ title: "更新异常", icon: "none" });
      console.error(err);
    }
  },

  goBack() {
    wx.navigateBack();
  }
});
