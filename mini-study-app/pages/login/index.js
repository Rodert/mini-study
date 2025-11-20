const api = require("../../services/api");
const app = getApp();

const PRESETS = {
  employee: {
    workNo: "employee001",
    password: "123456"
  },
  manager: {
    workNo: "manager001",
    password: "123456"
  },
  admin: {
    workNo: "admin",
    password: "admin123456"
  }
};

Page({
  data: {
    workNo: "",
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
      [field]: e.detail.value.trim(),
      error: ""
    });
  },

  async handleSubmit() {
    if (!this.validate()) {
      return;
    }
    this.setData({ loading: true, error: "" });

    try {
      const response = await api.user.login({
        work_no: this.data.workNo,
        password: this.data.password
      });

      if (response.code !== 200) {
        this.setData({ error: response.message || "登录失败" });
        return;
      }

      const profileRes = await api.user.getCurrentUser();
      if (profileRes.code !== 200) {
        this.setData({ error: profileRes.message || "获取用户信息失败" });
        return;
      }

      const user = profileRes.data;
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
    if (!this.data.workNo) {
      this.setData({ error: "请输入工号" });
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

