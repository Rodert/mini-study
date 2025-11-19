const mockService = require("../../services/mockService");
const app = getApp();

Page({
  data: {
    mobile: "",
    employeeId: "",
    name: "",
    store: "",
    password: "",
    managers: [],
    selectedManagerIds: [],
    loading: false,
    error: ""
  },

  onLoad() {
    this.loadManagers();
  },

  async loadManagers() {
    try {
      const res = await mockService.fetchManagers();
      this.setData({ managers: res.data || [] });
    } catch (err) {
      console.error("fetch managers error", err);
    }
  },

  handleInput(e) {
    const { field } = e.currentTarget.dataset;
    this.setData({
      [field]: e.detail.value,
      error: ""
    });
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
    
    this.setData({ selectedManagerIds, error: "" });
  },

  async handleSubmit() {
    if (!this.validate()) {
      return;
    }
    this.setData({ loading: true, error: "" });

    try {
      const payload = {
        mobile: this.data.mobile,
        employeeId: this.data.employeeId,
        name: this.data.name,
        store: this.data.store,
        password: this.data.password,
        role: "employee",
        managerIds: this.data.selectedManagerIds
      };
      const response = await mockService.register(payload);

      if (!response.success) {
        this.setData({ error: response.message || "注册失败" });
        return;
      }

      const user = response.data;
      app.globalData.user = user;
      wx.setStorageSync("user", user);
      wx.showToast({ title: "注册成功", icon: "success" });
      setTimeout(() => {
        wx.reLaunch({ url: "/pages/home/index" });
      }, 500);
    } catch (error) {
      this.setData({ error: "注册异常，请稍后再试" });
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
    if (!this.data.employeeId) {
      this.setData({ error: "请输入工号" });
      return false;
    }
    if (!this.data.name) {
      this.setData({ error: "请输入姓名" });
      return false;
    }
    if (!this.data.password) {
      this.setData({ error: "请输入密码" });
      return false;
    }
    if (this.data.selectedManagerIds.length === 0) {
      this.setData({ error: "请至少选择一个店长" });
      return false;
    }
    return true;
  },

  goLogin() {
    wx.navigateBack({ delta: 1 });
  }
});

