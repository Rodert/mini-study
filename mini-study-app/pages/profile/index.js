const mockService = require("../../services/mockService");
const app = getApp();

Page({
  data: {
    user: {},
    stats: null,
    managers: [],
    userManagers: []
  },

  onShow() {
    const cachedUser = app.globalData.user || wx.getStorageSync("user");
    if (cachedUser && cachedUser.id) {
      this.setData({ user: cachedUser });
      this.loadStats(cachedUser.id);
      this.loadManagers(cachedUser);
    } else {
      this.setData({ user: {}, stats: null });
    }
  },

  async loadStats(userId) {
    try {
      const res = await mockService.fetchLearningStats(userId);
      this.setData({ stats: res.data });
    } catch (err) {
      console.error("load stats error", err);
      this.setData({ stats: null });
    }
  },

  async loadManagers(user) {
    try {
      const res = await mockService.fetchManagers();
      const managers = res.data || [];
      this.setData({ managers });
      this.updateUserManagers(user, managers);
    } catch (err) {
      console.error("fetch managers error", err);
      this.updateUserManagers(user, []);
    }
  },

  updateUserManagers(user, managers = this.data.managers) {
    if (!user || !user.id) {
      this.setData({ userManagers: [] });
      return;
    }
    if (user.managerIds && user.managerIds.length > 0) {
      const userManagers = managers.filter((m) => user.managerIds.includes(m.id));
      this.setData({ userManagers });
    } else {
      this.setData({ userManagers: [] });
    }
  },

  goEditPage() {
    wx.navigateTo({ url: "/pages/profile/edit/index" });
  },

  goLearning() {
    wx.navigateTo({ url: "/pages/learning/index/index" });
  },

  handleLogout() {
    wx.showModal({
      title: "退出登录",
      content: "确定要退出当前账号吗？",
      success: (res) => {
        if (res.confirm) {
          app.globalData.user = null;
          wx.removeStorageSync("user");
          wx.showToast({ title: "已退出", icon: "none" });
          setTimeout(() => {
            wx.reLaunch({ url: "/pages/login/index" });
          }, 500);
        }
      }
    });
  },

  goBack() {
    wx.navigateBack();
  },

  goLogin() {
    wx.navigateTo({ url: "/pages/login/index" });
  }
});

