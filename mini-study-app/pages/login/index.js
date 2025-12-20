const api = require("../../services/api");
const config = require("../../services/config");
const app = getApp();

const PRESETS = {
  employee: {
    workNo: "employee001",
    password: "123456"  // 后端数据库默认密码
  },
  manager: {
    workNo: "manager001",
    password: "123456"  // 后端数据库默认密码
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
    lastPreset: "manager",
    useMock: config.USE_MOCK,  // 显示当前使用的模式
    apiUrl: config.API_BASE_URL  // 显示 API 地址
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
      console.log("开始登录，工号:", this.data.workNo);
      const response = await api.user.login({
        work_no: this.data.workNo,
        password: this.data.password
      });

      console.log("登录响应:", response);

      if (response.code !== 200) {
        const errorMsg = response.message || "登录失败";
        console.error("登录失败:", errorMsg, response);
        this.setData({ error: errorMsg });
        this.setData({ loading: false });
        return;
      }

      // 登录成功，获取用户信息
      console.log("登录成功，获取用户信息...");
      const profileRes = await api.user.getCurrentUser();
      console.log("用户信息响应:", profileRes);

      if (profileRes.code !== 200) {
        const errorMsg = profileRes.message || "获取用户信息失败";
        console.error("获取用户信息失败:", errorMsg, profileRes);
        this.setData({ error: errorMsg });
        this.setData({ loading: false });
        return;
      }

      const user = profileRes.data;
      console.log("用户信息:", user);
      app.globalData.user = user;
      wx.setStorageSync("user", user);
      app.globalData.hasShownNotice = false;
      wx.showToast({ title: "登录成功", icon: "success" });
      setTimeout(() => {
        wx.reLaunch({ url: "/pages/home/index" });
      }, 500);
    } catch (error) {
      console.error("登录异常:", error);
      const errorMsg = error.message || error.errMsg || "登录异常，请稍后重试";
      this.setData({ error: errorMsg });
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

