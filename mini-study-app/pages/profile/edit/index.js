const api = require("../../../services/api");
const app = getApp();

Page({
  data: {
    mobile: "",
    originalMobile: "",
    userManagers: [],
    user: {}
  },

  onLoad() {
    this.loadUserInfo();
  },

  async loadUserInfo() {
    try {
      // 从 API 获取最新的用户信息（包含店长信息）
      const res = await api.user.getCurrentUser();
      if (res.code === 200 && res.data) {
        const user = res.data;
        // 更新全局用户信息
        app.globalData.user = user;
        wx.setStorageSync("user", user);
        
        this.setData({
          user: user,
          mobile: user.phone || user.mobile || "",
          originalMobile: user.phone || user.mobile || "",
          // 直接使用 API 返回的 managers 字段
          userManagers: user.managers || []
        });
      } else {
        // 如果 API 失败，使用缓存的用户信息
        const cachedUser = app.globalData.user || wx.getStorageSync("user");
        if (!cachedUser || !cachedUser.id) {
          wx.showToast({ title: "请先登录", icon: "none" });
          setTimeout(() => {
            wx.navigateBack();
          }, 1000);
          return;
        }
        
        this.setData({
          user: cachedUser,
          mobile: cachedUser.phone || cachedUser.mobile || "",
          originalMobile: cachedUser.phone || cachedUser.mobile || "",
          userManagers: cachedUser.managers || []
        });
      }
    } catch (err) {
      console.error("load user info error", err);
      // 如果 API 失败，使用缓存的用户信息
      const cachedUser = app.globalData.user || wx.getStorageSync("user");
      if (!cachedUser || !cachedUser.id) {
        wx.showToast({ title: "请先登录", icon: "none" });
        setTimeout(() => {
          wx.navigateBack();
        }, 1000);
        return;
      }
      
      this.setData({
        user: cachedUser,
        mobile: cachedUser.phone || cachedUser.mobile || "",
        originalMobile: cachedUser.phone || cachedUser.mobile || "",
        userManagers: cachedUser.managers || []
      });
    }
  },

  handleMobileInput(e) {
    this.setData({ mobile: e.detail.value });
  },

  async handleSave() {
    // 验证手机号
    if (!this.data.mobile || this.data.mobile.length !== 11) {
      wx.showToast({ title: "请输入11位手机号", icon: "none" });
      return;
    }

    // 检查是否有修改
    if (this.data.mobile === this.data.originalMobile) {
      wx.showToast({ title: "手机号未修改", icon: "none" });
      return;
    }

    wx.showLoading({ title: "保存中..." });

    try {
      const response = await api.user.updateProfile({
        phone: this.data.mobile
      });

      if (response.code === 200) {
        // 重新获取用户信息以获取最新数据（包括店长信息）
        const userRes = await api.user.getCurrentUser();
        if (userRes.code === 200 && userRes.data) {
          const user = userRes.data;
          app.globalData.user = user;
          wx.setStorageSync("user", user);
        }
        wx.hideLoading();
        wx.showToast({ title: "手机号已更新", icon: "success" });
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
