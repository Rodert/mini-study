const api = require("../../../../services/api");
const app = getApp();

Page({
  data: {
    userId: null,
    user: {},
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
      const [userRes, managersRes] = await Promise.all([
        api.admin.getUser(this.data.userId),
        api.user.getManagers()
      ]);

      if (userRes.code !== 200 || !userRes.data) {
        wx.showToast({ title: "用户不存在", icon: "none" });
        wx.navigateBack();
        return;
      }

      const user = userRes.data;
      const selectedManagerIds = user.manager_ids || [];
      this.setData({
        user,
        managers: managersRes.code === 200 ? (managersRes.data || []) : [],
        selectedManagerIds,
        originalRole: user.role,
        originalManagerIds: [...selectedManagerIds]
      });
    } catch (err) {
      console.error("load data error", err);
      wx.showToast({ title: "加载数据失败", icon: "none" });
    }
  },

  selectRole(e) {
    const { role } = e.currentTarget.dataset;
    const user = { ...this.data.user, role };
    
    // 如果切换为非员工角色，清空店长绑定
    if (role !== "employee") {
      this.setData({
        user,
        selectedManagerIds: []
      });
    } else {
      this.setData({ user });
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
    const { user, selectedManagerIds, originalRole, originalManagerIds } = this.data;

    // 检查是否有修改
    const roleChanged = user.role !== originalRole;
    const managerChanged = JSON.stringify(selectedManagerIds) !== JSON.stringify(originalManagerIds);

    if (!roleChanged && !managerChanged) {
      wx.showToast({ title: "没有修改任何信息", icon: "none" });
      return;
    }

    // 如果是员工，验证是否选择了店长
    if (user.role === "employee" && selectedManagerIds.length === 0) {
      wx.showToast({ title: "员工必须绑定至少一个店长", icon: "none" });
      return;
    }

    wx.showLoading({ title: "保存中..." });

    try {
      let response;

      if (roleChanged && managerChanged) {
        // 同时修改角色和店长绑定
        await api.admin.updateUserRole(this.data.userId, user.role);
        // 获取店长工号列表（从 managers 中查找）
        const managers = this.data.managers;
        const managerWorkNos = selectedManagerIds
          .map(id => {
            const m = managers.find(m => m.id === id);
            return m ? m.work_no : null;
          })
          .filter(Boolean);
        response = await api.admin.updateEmployeeManagers(this.data.userId, managerWorkNos);
      } else if (roleChanged) {
        // 只修改角色
        response = await api.admin.updateUserRole(this.data.userId, user.role);
      } else {
        // 只修改店长绑定
        const managers = this.data.managers;
        const managerWorkNos = selectedManagerIds
          .map(id => {
            const m = managers.find(m => m.id === id);
            return m ? m.work_no : null;
          })
          .filter(Boolean);
        response = await api.admin.updateEmployeeManagers(this.data.userId, managerWorkNos);
      }

      if (response.code === 200) {
        wx.hideLoading();
        wx.showToast({ title: "用户已更新", icon: "success" });
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
