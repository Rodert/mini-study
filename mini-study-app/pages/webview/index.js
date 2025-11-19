Page({
  data: {
    url: ""
  },

  onLoad(options) {
    const { url } = options || {};
    if (url) {
      this.setData({ url: decodeURIComponent(url) });
    } else {
      console.warn("webview url is empty");
      wx.showToast({ title: "无效链接", icon: "none" });
    }
  }
});

