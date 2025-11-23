const api = require("../../../services/api");
const app = getApp();

Page({
  data: {
    users: []
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

    this.loadUsers();
  },

  async loadUsers() {
    try {
      const res = await api.admin.listUsers();
      if (res.code === 200) {
        this.setData({ users: res.data || [] });
      } else {
        wx.showToast({ title: res.message || "加载用户失败", icon: "none" });
      }
    } catch (err) {
      console.error("fetch users error", err);
      wx.showToast({ title: "加载用户失败", icon: "none" });
    }
  },

  goEditUser(e) {
    const { userId } = e.currentTarget.dataset;
    wx.navigateTo({ 
      url: `/pages/manager/users/edit/index?userId=${userId}`
    });
  },

  goBack() {
    wx.navigateBack();
  }
});
