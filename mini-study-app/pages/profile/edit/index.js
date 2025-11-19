const mockService = require("../../../services/mockService");
const app = getApp();

Page({
  data: {
    mobile: "",
    originalMobile: "",
    managers: [],
    userManagers: [],
    user: {}
  },

  onLoad() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || !user.id) {
      wx.showToast({ title: "请先登录", icon: "none" });
      setTimeout(() => {
        wx.navigateBack();
      }, 1000);
      return;
    }

    this.setData({
      user,
      mobile: user.mobile,
      originalMobile: user.mobile
    });

    this.loadManagers(user);
  },

  async loadManagers(user) {
    try {
      const res = await mockService.fetchManagers();
      const managers = res.data || [];
      const userManagers = managers.filter(m => user.managerIds && user.managerIds.includes(m.id));
      this.setData({ 
        managers,
        userManagers
      });
    } catch (err) {
      console.error("fetch managers error", err);
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
      const response = await mockService.updateUserProfile(this.data.user.id, {
        mobile: this.data.mobile
      });

      if (response.success) {
        const user = response.data;
        app.globalData.user = user;
        wx.setStorageSync("user", user);
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
