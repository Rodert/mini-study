const api = require("../../../../services/api");
const app = getApp();

Page({
  data: {
    user: {},
    contents: [],
    loading: false,
    typeFilterOptions: [
      { label: "全部类型", value: "" },
      { label: "文档", value: "doc" },
      { label: "视频", value: "video" },
      { label: "图文", value: "article" }
    ],
    typeFilterIndex: 0,
    statusFilterOptions: [
      { label: "全部状态", value: "" },
      { label: "草稿", value: "draft" },
      { label: "已发布", value: "published" },
      { label: "已下线", value: "offline" }
    ],
    statusFilterIndex: 0
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || user.role !== "admin") {
      wx.showToast({ title: "仅管理员可访问", icon: "none" });
      setTimeout(() => wx.navigateBack(), 800);
      return;
    }
    this.setData({ user });
    this.loadContents();
  },

  async loadContents() {
    this.setData({ loading: true });
    try {
      const typeOption = this.data.typeFilterOptions[this.data.typeFilterIndex];
      const statusOption = this.data.statusFilterOptions[this.data.statusFilterIndex];
      const params = {};
      if (typeOption && typeOption.value) {
        params.type = typeOption.value;
      }
      if (statusOption && statusOption.value) {
        params.status = statusOption.value;
      }

      const res = await api.admin.listContents(params);
      if (res.code === 200) {
        const list = (res.data || []).map((item) => ({
          id: item.id,
          title: item.title || "",
          categoryName: item.category_name || "",
          type: item.type || "",
          typeText:
            item.type === "video"
              ? "视频"
              : item.type === "article"
              ? "图文"
              : "文档",
          status: item.status || "",
          statusText:
            item.status === "published"
              ? "已发布"
              : item.status === "offline"
              ? "已下线"
              : "草稿",
          visibleRoles: item.visible_roles || "both",
          visibleRolesText:
            item.visible_roles === "employee"
              ? "仅员工"
              : item.visible_roles === "manager"
              ? "仅店长"
              : "全部",
          publishAt: item.publish_at || null
        }));
        this.setData({ contents: list });
      } else {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load contents error", err);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  handleTypeFilterChange(e) {
    const index = Number(e.detail.value);
    this.setData({ typeFilterIndex: index });
    this.loadContents();
  },

  handleStatusFilterChange(e) {
    const index = Number(e.detail.value);
    this.setData({ statusFilterIndex: index });
    this.loadContents();
  },

  goCreate() {
    wx.navigateTo({ url: "/pages/admin/contents/add/index" });
  },

  goEdit(e) {
    const id = Number(e.currentTarget.dataset.id);
    if (!id) return;
    wx.navigateTo({ url: `/pages/admin/contents/add/index?id=${id}` });
  },

  async handleToggleStatus(e) {
    const id = Number(e.currentTarget.dataset.id);
    const currentStatus = e.currentTarget.dataset.status;
    if (!id || !currentStatus) return;

    let targetStatus = "";
    let actionText = "";
    if (currentStatus === "published") {
      targetStatus = "offline";
      actionText = "下线";
    } else if (currentStatus === "offline") {
      targetStatus = "published";
      actionText = "重新上线";
    } else {
      return;
    }

    wx.showModal({
      title: "确认操作",
      content: `确定要${actionText}该内容吗？`,
      success: async (res) => {
        if (!res.confirm) return;
        try {
          wx.showLoading({ title: "处理中..." });
          const resp = await api.admin.updateContent(id, { status: targetStatus });
          if (resp.code === 200) {
            wx.showToast({ title: `${actionText}成功`, icon: "success" });
            this.loadContents();
          } else {
            wx.showToast({ title: resp.message || `${actionText}失败`, icon: "none" });
          }
        } catch (err) {
          console.error("toggle content status error", err);
          wx.showToast({ title: `${actionText}失败`, icon: "none" });
        } finally {
          wx.hideLoading();
        }
      }
    });
  },

  goBack() {
    wx.navigateBack();
  }
});
