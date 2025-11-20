const api = require("../../services/api");
const mockService = require("../../services/mockService");
const app = getApp();

Page({
  data: {
    user: {},
    stats: null,
    userManagers: []
  },

  onShow() {
    this.loadUserInfo();
  },

  async loadUserInfo() {
    try {
      // 从 API 获取最新的用户信息（包含店长信息）
      const res = await api.user.getCurrentUser();
      console.log("getCurrentUser response:", res);
      if (res.code === 200 && res.data) {
        const user = res.data;
        console.log("user data:", user);
        console.log("user managers:", user.managers);
        // 更新全局用户信息
        app.globalData.user = user;
        wx.setStorageSync("user", user);
        
        const isAdmin = user.role === "admin";
        const managers = user.managers || [];
        console.log("setting userManagers:", managers);
        this.setData({ 
          user: user,
          // 直接使用 API 返回的 managers 字段
          userManagers: managers
        });
        
        if (!isAdmin) {
          this.loadStats(user.id);
        } else {
          this.setData({ stats: null });
        }
      } else {
        // 如果 API 失败，使用缓存的用户信息
        const cachedUser = app.globalData.user || wx.getStorageSync("user");
        if (cachedUser && cachedUser.id) {
          this.setData({ 
            user: cachedUser,
            userManagers: cachedUser.managers || []
          });
          if (cachedUser.role !== "admin") {
            this.loadStats(cachedUser.id);
          }
        } else {
          this.setData({ user: {}, stats: null, userManagers: [] });
        }
      }
    } catch (err) {
      console.error("load user info error", err);
      // 如果 API 失败，使用缓存的用户信息
      const cachedUser = app.globalData.user || wx.getStorageSync("user");
      if (cachedUser && cachedUser.id) {
        this.setData({ 
          user: cachedUser,
          userManagers: cachedUser.managers || []
        });
        if (cachedUser.role !== "admin") {
          this.loadStats(cachedUser.id);
        }
      } else {
        this.setData({ user: {}, stats: null, userManagers: [] });
      }
    }
  },

  async loadStats(userId) {
    try {
      const res = await api.learning.getUserStats();
      if (res.code === 200 && res.data) {
        // 转换API返回的数据格式以匹配页面显示
        const statsData = res.data;
        const completionRate = statsData.completion_rate || 0;
        this.setData({ 
          stats: {
            completedCourses: statsData.completed_count || 0,
            totalCount: statsData.total_count || 0, // 已开始学习的数量
            totalCourses: statsData.total_contents || 0, // 总课程数
            completionRate: completionRate.toFixed(1) // 完成率，保留一位小数
          }
        });
      } else {
        this.setData({ stats: null });
      }
    } catch (err) {
      console.error("load stats error", err);
      this.setData({ stats: null });
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

