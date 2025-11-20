const api = require("../../../../services/api");
const app = getApp();

Page({
  data: {
    userId: null,
    employee: {},
    managers: [],
    selectedManagerIds: [],
    originalRole: "",
    originalManagerIds: [],
    loading: false
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
      this.setData({ loading: true });
      const [userRes, managersRes] = await Promise.all([
        api.admin.getUser(this.data.userId),
        api.user.getManagers()
      ]);

      if (userRes.code !== 200) {
        wx.showToast({ title: userRes.message || "获取用户失败", icon: "none" });
        wx.navigateBack();
        return;
      }

      const employeeData = userRes.data || {};
      const employee = {
        id: employeeData.id,
        name: employeeData.name,
        username: employeeData.work_no,
        mobile: employeeData.phone || "--",
        role: employeeData.role,
        store: employeeData.store || "—",
        managers: employeeData.managers || []
      };

      const managerOptions = (managersRes?.data || [])
        .filter(manager => manager.id !== employeeData.id)
        .map(manager => ({
          id: manager.id,
          name: manager.name || manager.work_no,
          workNo: manager.work_no,
          phone: manager.phone || "",
          store: manager.phone ? `电话：${manager.phone}` : ""
        }));

      const selectedManagerIds = (employeeData.manager_ids || []).map(id => Number(id));

      this.setData({
        employee,
        managers: managerOptions,
        selectedManagerIds,
        originalRole: employee.role,
        originalManagerIds: [...selectedManagerIds]
      });
    } catch (err) {
      console.error("load data error", err);
      wx.showToast({ title: "加载数据失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
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
    const managerId = Number(e.currentTarget.dataset.managerId);
    const selectedManagerIds = [...this.data.selectedManagerIds];
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
    const managerChanged = !this.arraysEqual(selectedManagerIds, originalManagerIds);

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
      if (roleChanged) {
        const res = await api.admin.updateUserRole(this.data.userId, employee.role);
        if (res.code !== 200) {
          throw new Error(res.message || "角色更新失败");
        }
      }

      if (employee.role === "employee" && managerChanged) {
        const managerWorkNos = this.getSelectedManagerWorkNos();
        const res = await api.admin.updateEmployeeManagers(this.data.userId, managerWorkNos);
        if (res.code !== 200) {
          throw new Error(res.message || "店长绑定失败");
        }
      }

      wx.hideLoading();
      wx.showToast({ title: "员工已更新", icon: "success" });
      setTimeout(() => {
        wx.navigateBack();
      }, 600);
    } catch (err) {
      wx.hideLoading();
      wx.showToast({ title: err.message || "更新异常", icon: "none" });
      console.error(err);
    }
  },

  arraysEqual(a = [], b = []) {
    if (a.length !== b.length) {
      return false;
    }
    const sortedA = [...a].sort((x, y) => x - y);
    const sortedB = [...b].sort((x, y) => x - y);
    return sortedA.every((val, idx) => val === sortedB[idx]);
  },

  getSelectedManagerWorkNos() {
    const { managers, selectedManagerIds } = this.data;
    const idToWorkNo = new Map();
    managers.forEach(manager => {
      idToWorkNo.set(manager.id, manager.workNo);
    });
    return selectedManagerIds
      .map(id => idToWorkNo.get(id))
      .filter(Boolean);
  },

  goBack() {
    wx.navigateBack();
  }
});
