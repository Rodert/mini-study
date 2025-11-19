const api = require("../../../../services/api");
const app = getApp();

Page({
  data: {
    form: {
      work_no: "",
      name: "",
      phone: "",
      password: ""
    },
    loading: false,
    errors: {}
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
  },

  handleInput(e) {
    const { field } = e.currentTarget.dataset;
    const value = e.detail.value;
    this.setData({
      [`form.${field}`]: value,
      [`errors.${field}`]: ""
    });
  },

  validate() {
    const { form } = this.data;
    const errors = {};

    if (!form.work_no || form.work_no.trim().length < 2) {
      errors.work_no = "工号至少2位";
    }

    if (!form.name || form.name.trim().length === 0) {
      errors.name = "请输入姓名";
    }

    if (!form.password || form.password.length < 6) {
      errors.password = "密码至少6位";
    }

    if (form.phone && form.phone.length !== 11) {
      errors.phone = "手机号格式不正确";
    }

    this.setData({ errors });
    return Object.keys(errors).length === 0;
  },

  async handleSubmit() {
    if (!this.validate()) {
      wx.showToast({ title: "请检查表单", icon: "none" });
      return;
    }

    this.setData({ loading: true });

    try {
      const res = await api.admin.createManager({
        work_no: this.data.form.work_no.trim(),
        name: this.data.form.name.trim(),
        phone: this.data.form.phone.trim() || undefined,
        password: this.data.form.password
      });

      if (res.code === 200) {
        wx.showToast({ title: "创建成功", icon: "success" });
        setTimeout(() => {
          wx.navigateBack();
        }, 1500);
      } else {
        wx.showToast({ title: res.message || "创建失败", icon: "none" });
      }
    } catch (err) {
      console.error("create manager error", err);
      wx.showToast({ 
        title: err.message || "创建失败，请稍后重试", 
        icon: "none" 
      });
    } finally {
      this.setData({ loading: false });
    }
  },

  goBack() {
    wx.navigateBack();
  }
});

