const api = require("../../services/api");
const app = getApp();

Page({
  data: {
    phone: "",
    workNo: "",
    name: "",
    password: "",
    managers: [],
    selectedManagerWorkNos: [],
    loading: false,
    error: ""
  },

  onLoad() {
    this.loadManagers();
  },

  onShow() {
    // 确保每次进入页面时都从未选择状态开始
    this.setData({ selectedManagerWorkNos: [] });
  },

  async loadManagers() {
    try {
      const res = await api.user.getManagers();
      if (res.code === 200) {
        const managers = res.data || [];
        this.setData({
          managers,
          // 确保每次加载列表时都清空已选
          selectedManagerWorkNos: []
        });
      } else {
        this.setData({ error: res.message || "获取店长列表失败" });
      }
    } catch (err) {
      console.error("fetch managers error", err);
      this.setData({ error: "获取店长列表失败" });
    }
  },

  handleInput(e) {
    const { field } = e.currentTarget.dataset;
    this.setData({
      [field]: e.detail.value.trim(),
      error: ""
    });
  },

  handleManagerSelect(e) {
    console.log("handleManagerSelect tap event:", e);
    const workNo = String(e.currentTarget.dataset.workNo || "");
    console.log("handleManagerSelect workNo from dataset:", workNo, "current selected:", this.data.selectedManagerWorkNos);

    // 简单可视化反馈，便于确认点击是否触发
    wx.showToast({
      title: workNo ? `选择店长: ${workNo}` : "点击了店长项",
      icon: "none",
      duration: 800
    });
    const selectedManagerWorkNos = [...this.data.selectedManagerWorkNos];
    const index = selectedManagerWorkNos.indexOf(workNo);

    if (index > -1) {
      selectedManagerWorkNos.splice(index, 1);
    } else {
      selectedManagerWorkNos.push(workNo);
    }

    this.setData({ selectedManagerWorkNos, error: "" }, () => {
      console.log("handleManagerSelect updated selectedManagerWorkNos:", this.data.selectedManagerWorkNos);
    });
  },

  async handleSubmit() {
    if (!this.validate()) {
      return;
    }
    this.setData({ loading: true, error: "" });

    try {
      const payload = {
        work_no: this.data.workNo,
        phone: this.data.phone,
        name: this.data.name,
        password: this.data.password,
        manager_ids: this.data.selectedManagerWorkNos
      };
      const response = await api.user.register(payload);

      if (response.code !== 200) {
        this.setData({ error: response.message || "注册失败" });
        return;
      }

      // 自动登录
      const loginRes = await api.user.login({
        work_no: payload.work_no,
        password: payload.password
      });
      if (loginRes.code !== 200) {
        this.setData({ error: loginRes.message || "登录失败，请手动登录" });
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
      wx.showToast({ title: "注册并登录成功", icon: "success" });
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
    if (!this.data.phone || this.data.phone.length !== 11) {
      this.setData({ error: "请输入 11 位手机号" });
      return false;
    }
    if (!this.data.workNo) {
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
    if (this.data.selectedManagerWorkNos.length === 0) {
      this.setData({ error: "请至少选择一个店长" });
      return false;
    }
    return true;
  },

  goLogin() {
    wx.navigateBack({ delta: 1 });
  }
});

