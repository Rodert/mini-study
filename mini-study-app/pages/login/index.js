const mockService = require("../../services/mockService");
const app = getApp();

const PRESETS = {
  employee: {
    mobile: "13900000000",
    password: "123456"
  },
  manager: {
    mobile: "13911110000",
    password: "123456"
  },
  admin: {
    mobile: "13900000001",
    password: "123456"
  }
};

Page({
  data: {
    mobile: "",
    password: "",
    loading: false,
    error: "",
    lastPreset: "manager"
  },

  onShow() {
    const cachedUser = app.globalData.user || wx.getStorageSync("user");
    if (cachedUser && cachedUser.id) {
      wx.reLaunch({
        url: "/pages/home/index"
      });
    }
  },

  handleInput(e) {
    const { field } = e.currentTarget.dataset;
    this.setData({
      [field]: e.detail.value,
      error: ""
    });
  },

  async handleSubmit() {
    if (!this.validate()) {
      return;
    }
    this.setData({ loading: true, error: "" });

    try {
      const response = await mockService.login({
        mobile: this.data.mobile,
        password: this.data.password
      });

      if (!response.success) {
        this.setData({ error: response.message || "登录失败" });
        return;
      }

      const user = response.data;
      app.globalData.user = user;
      wx.setStorageSync("user", user);
      wx.showToast({ title: "登录成功", icon: "success" });
      setTimeout(() => {
        wx.reLaunch({ url: "/pages/home/index" });
      }, 500);
    } catch (error) {
      this.setData({ error: "登录异常，请稍后重试" });
      console.error(error);
    } finally {
      this.setData({ loading: false });
    }
  },

  validate() {
    if (!this.data.mobile || this.data.mobile.length !== 11) {
      this.setData({ error: "请输入 11 位手机号" });
      return false;
    }
    if (!this.data.password) {
      this.setData({ error: "请输入密码" });
      return false;
    }
    return true;
  },

  goRegister() {
    wx.navigateTo({ url: "/pages/register/index" });
  },

  quickSwitch() {
    const roles = ["employee", "manager", "admin"];
    const currentIndex = roles.indexOf(this.data.lastPreset);
    const nextIndex = (currentIndex + 1) % roles.length;
    const nextRole = roles[nextIndex];
    
    this.setData({
      ...PRESETS[nextRole],
      lastPreset: nextRole,
      error: ""
    });
  }
});

