const api = require("../../../../services/api");
const app = getApp();

Page({
  data: {
    form: {
      work_no: "",
      name: "",
      phone: "",
      password: "",
      manager_ids: []
    },
    managers: [],
    selectedManagers: [],
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

    this.loadManagers();
  },

  async loadManagers() {
    try {
      wx.showLoading({ title: "加载中..." });
      const res = await api.user.getManagers();
      if (res.code === 200) {
        // 为每个店长添加 selected 属性
        const managers = (res.data || []).map(manager => ({
          ...manager,
          selected: false
        }));
        this.setData({
          managers
        });
      } else {
        wx.showToast({ title: res.message || "加载店长列表失败", icon: "none" });
      }
    } catch (err) {
      console.error("load managers error", err);
      wx.showToast({ title: err.message || "加载店长列表失败", icon: "none" });
    } finally {
      wx.hideLoading();
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

  toggleManager(e) {
    const { managerId } = e.currentTarget.dataset;
    const managerIdNum = parseInt(managerId);
    
    // 更新 managers 数组中对应项的 selected 状态
    const managers = this.data.managers.map(m => {
      if (m.id === managerIdNum) {
        return { ...m, selected: !m.selected };
      }
      return m;
    });

    // 计算选中的店长列表
    const selectedManagers = managers.filter(m => m.selected);
    const manager_ids = selectedManagers.map(m => m.work_no);

    this.setData({
      managers,
      selectedManagers,
      "form.manager_ids": manager_ids
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
      const res = await api.admin.createEmployee({
        work_no: this.data.form.work_no.trim(),
        name: this.data.form.name.trim(),
        phone: this.data.form.phone.trim() || undefined,
        password: this.data.form.password,
        manager_ids: this.data.form.manager_ids
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
      console.error("create employee error", err);
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

