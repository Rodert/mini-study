// pages/admin/growth/index/index.js
const api = require("../../../../services/api");
const app = getApp();

Page({

  /**
   * 页面的初始数据
   */
  data: {
    user: {},
    posts: [],
    keyword: "",
    statusFilterOptions: [
      { label: "全部", value: "" },
      { label: "待审核", value: "pending" },
      { label: "已通过", value: "approved" },
      { label: "已拒绝", value: "rejected" }
    ],
    statusFilterIndex: 0,
    loading: false
  },

  onLoad() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || user.role !== "admin") {
      wx.showToast({ title: "无权限访问", icon: "none" });
      setTimeout(() => wx.navigateBack(), 800);
      return;
    }
    this.setData({ user });
    this.loadPosts();
  },

  onPullDownRefresh() {
    this.loadPosts();
  },

  async loadPosts() {
    if (this.data.loading) return;
    this.setData({ loading: true });
    try {
      const statusOpt = this.data.statusFilterOptions[this.data.statusFilterIndex];
      const params = { keyword: this.data.keyword };
      if (statusOpt && statusOpt.value) {
        params.status = statusOpt.value;
      }
      const res = await api.admin.listGrowth(params);
      if (res.code === 200) {
        const posts = (res.data || []).map((item) => this.transformPost(item));
        this.setData({ posts });
      } else {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("admin load growth posts error", err);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
      wx.stopPullDownRefresh();
    }
  },

  transformPost(item) {
    const publisherName = item.publisher_name || "";
    const avatarText = publisherName ? publisherName.charAt(0) : "?";
    let statusText = "";
    if (item.status === "pending") statusText = "待审核";
    else if (item.status === "approved") statusText = "已通过";
    else if (item.status === "rejected") statusText = "已拒绝";

    const createdAtText = this.formatDateTime(item.created_at);

    return {
      id: item.id,
      content: item.content || "",
      status: item.status,
      statusText,
      publisherName,
      publisherRole: item.publisher_role || "",
      avatarText,
      createdAtText
    };
  },

  formatDateTime(isoString) {
    if (!isoString) return "";
    const d = new Date(isoString);
    if (Number.isNaN(d.getTime())) return isoString;
    const pad = (n) => (n < 10 ? `0${n}` : `${n}`);
    const y = d.getFullYear();
    const m = pad(d.getMonth() + 1);
    const day = pad(d.getDate());
    const h = pad(d.getHours());
    const mi = pad(d.getMinutes());
    return `${y}-${m}-${day} ${h}:${mi}`;
  },

  handleStatusChange(e) {
    const index = Number(e.detail.value);
    this.setData({ statusFilterIndex: index });
    this.loadPosts();
  },

  handleSearchInput(e) {
    this.setData({ keyword: e.detail.value || "" });
  },

  handleSearchConfirm() {
    this.loadPosts();
  },

  async handleApprove(e) {
    const id = Number(e.currentTarget.dataset.id);
    if (!id) return;
    try {
      const res = await api.admin.approveGrowth(id);
      if (res.code === 200) {
        wx.showToast({ title: "已通过", icon: "success" });
        this.loadPosts();
      } else {
        wx.showToast({ title: res.message || "操作失败", icon: "none" });
      }
    } catch (err) {
      console.error("approve growth error", err);
      wx.showToast({ title: "操作失败", icon: "none" });
    }
  },

  async handleReject(e) {
    const id = Number(e.currentTarget.dataset.id);
    if (!id) return;
    try {
      const res = await api.admin.rejectGrowth(id);
      if (res.code === 200) {
        wx.showToast({ title: "已拒绝", icon: "success" });
        this.loadPosts();
      } else {
        wx.showToast({ title: res.message || "操作失败", icon: "none" });
      }
    } catch (err) {
      console.error("reject growth error", err);
      wx.showToast({ title: "操作失败", icon: "none" });
    }
  },

  async handleDelete(e) {
    const id = Number(e.currentTarget.dataset.id);
    if (!id) return;
    wx.showModal({
      title: "删除动态",
      content: "确定要删除这条成长圈动态吗？",
      success: async (res) => {
        if (!res.confirm) return;
        try {
          const resp = await api.growth.delete(id);
          if (resp.code === 200) {
            wx.showToast({ title: "已删除", icon: "success" });
            this.loadPosts();
          } else {
            wx.showToast({ title: resp.message || "删除失败", icon: "none" });
          }
        } catch (err) {
          console.error("admin delete growth error", err);
          wx.showToast({ title: "删除失败", icon: "none" });
        }
      }
    });
  },

  goBack() {
    wx.navigateBack();
  }
});