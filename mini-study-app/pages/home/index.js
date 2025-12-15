const api = require("../../services/api");
const app = getApp();

Page({
  data: {
    user: {},
    banners: [],
    categories: []
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || !user.id) {
      wx.reLaunch({ url: "/pages/login/index" });
      return;
    }
    this.setData({ user });
    this.loadInitialData();
  },

  async loadInitialData() {
    if (this.data.user.role !== "admin") {
      await Promise.all([this.loadBanners(), this.loadCategories()]);
    }
  },

  async loadBanners() {
    try {
      console.log("[Banner] 开始加载轮播图数据");
      const res = await api.banner.listVisible();
      console.log("[Banner] API响应:", res);
      if (res.code === 200) {
        const banners = (res.data || []).map((item) => ({
          id: item.id,
          title: item.title,
          cover: item.image_url ? api.buildFileUrl(item.image_url) : "",
          url: item.link_url,
          type: item.visible_roles
        }));
        console.log("[Banner] 处理后的轮播图数据:", banners);
        console.log("[Banner] 每个轮播图的URL:", banners.map(b => ({ title: b.title, url: b.url })));
        this.setData({ banners });
        console.log("[Banner] 轮播图数据已设置到页面，数量:", banners.length);
      } else {
        console.error("[Banner] API返回错误:", res);
        wx.showToast({ title: res.message || "轮播加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("[Banner] 加载轮播图异常:", err);
      wx.showToast({ title: "轮播加载失败", icon: "none" });
    }
  },

  async loadCategories() {
    try {
      const res = await api.content.listCategories();
      if (res.code === 200) {
        const categories = (res.data || []).map((item) => ({
          id: item.id,
          name: item.name,
          icon: "📖",
          count: item.count || 0
        }));
        this.setData({ categories });
      } else {
        wx.showToast({ title: res.message || "分类加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("fetch categories error", err);
      wx.showToast({ title: "分类加载失败", icon: "none" });
    }
  },

  reloadBanners() {
    this.loadBanners();
  },

  handleBannerTap(e) {
    console.log("[Banner] 点击事件触发");
    console.log("[Banner] 事件对象:", e);
    console.log("[Banner] currentTarget:", e.currentTarget);
    console.log("[Banner] dataset:", e.currentTarget.dataset);
    
    const { item } = e.currentTarget.dataset;
    console.log("[Banner] 提取的item:", item);
    console.log("[Banner] item的类型:", typeof item);
    console.log("[Banner] item的URL:", item ? item.url : "item为空");
    
    if (!item) {
      console.warn("[Banner] item为空，无法处理点击");
      wx.showToast({ title: "数据异常，稍后重试", icon: "none" });
      return;
    }
    
    if (!item.url) {
      console.warn("[Banner] 轮播图没有链接URL:", item);
      wx.showToast({ title: "该轮播图暂无链接", icon: "none" });
      return;
    }
    
    console.log("[Banner] 准备跳转到webview，URL:", item.url);
    const targetUrl = `/pages/webview/index?url=${encodeURIComponent(item.url)}`;
    console.log("[Banner] 目标页面路径:", targetUrl);
    
    try {
      wx.navigateTo({
        url: targetUrl,
        success: (res) => {
          console.log("[Banner] 跳转成功:", res);
        },
        fail: (err) => {
          console.error("[Banner] 跳转失败:", err);
          wx.showToast({ title: `打开失败: ${err.errMsg || "未知错误"}`, icon: "none", duration: 2000 });
        }
      });
    } catch (err) {
      console.error("[Banner] navigateTo异常:", err);
      wx.showToast({ title: "打开失败", icon: "none" });
    }
  },

  goProgress() {
    wx.navigateTo({ url: "/pages/manager/progress/index" });
  },

  handleSelectCategory(e) {
    const { item } = e.currentTarget.dataset;
    if (!item) return;
    wx.navigateTo({
      url: `/pages/learning/list/index?categoryId=${item.id}&name=${item.name}`
    });
  },

  goProfile() {
    wx.navigateTo({ url: "/pages/profile/index" });
  },

  goUserManagement() {
    wx.navigateTo({ url: "/pages/manager/users/index" });
  },

  goEmployeesList() {
    wx.navigateTo({ url: "/pages/admin/employees/index" });
  },

  goBannerManagement() {
    wx.navigateTo({ url: "/pages/admin/banners/index" });
  },

  goContentCreate() {
    wx.navigateTo({ url: "/pages/admin/contents/index/index" });
  },

  goExamManagement() {
    wx.navigateTo({ url: "/pages/admin/exams/index" });
  },

  goExamList() {
    wx.navigateTo({ url: "/pages/exams/list/index" });
  },

  goPointsManagement() {
    wx.navigateTo({ url: "/pages/admin/points/index" });
  },

  goGrowthManagement() {
    wx.navigateTo({ url: "/pages/admin/growth/index/index" });
  }
});

