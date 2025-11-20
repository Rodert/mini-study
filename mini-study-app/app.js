const api = require("./services/api");

App({
  async onLaunch() {
    console.log("Mini Study app launched");
    await this.restoreSession();
  },
  async restoreSession() {
    const token = wx.getStorageSync("token");
    if (!token) {
      return;
    }
    try {
      const res = await api.user.getCurrentUser();
      if (res.code === 200) {
        this.globalData.user = res.data;
        wx.setStorageSync("user", res.data);
      } else {
        this.clearSession();
      }
    } catch (error) {
      console.error("restore session failed", error);
      this.clearSession();
    }
  },
  clearSession() {
    wx.removeStorageSync("token");
    wx.removeStorageSync("refresh_token");
    wx.removeStorageSync("user");
    this.globalData.user = null;
  },
  globalData: {
    theme: {
      brand: "#2563eb",
      brandDark: "#1e3a8a",
      accent: "#f97316",
      text: "#1f2937",
      muted: "#6b7280",
      card: "#ffffff",
      background: "#f4f6fb"
    },
    user: null
  }
});

