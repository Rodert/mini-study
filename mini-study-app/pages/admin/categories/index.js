const api = require("../../../services/api");
const app = getApp();
const DEFAULT_CATEGORY_COVER = "https://img2024.cnblogs.com/blog/1326459/202512/1326459-20251221224420399-2001683945.png";

Page({
  data: {
    user: {},
    categories: [],
    loading: false
  },

  onShow() {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || user.role !== "admin") {
      wx.showToast({ title: "仅管理员可访问", icon: "none" });
      setTimeout(() => wx.navigateBack(), 800);
      return;
    }
    this.setData({ user });
    this.loadCategories();
  },

  async loadCategories() {
    this.setData({ loading: true });
    try {
      const res = await api.admin.listCategories();
      if (res.code === 200) {
        const categories = (res.data || []).map((item) => {
          const hasCustomCover = !!item.cover_url;
          const coverUrl = item.cover_url || "";
          const coverPreview = hasCustomCover ? api.buildFileUrl(coverUrl) : "";
          const displayCover = hasCustomCover ? coverPreview : DEFAULT_CATEGORY_COVER;

          return {
            id: item.id,
            name: item.name,
            role_scope: item.role_scope,
            sort_order: item.sort_order,
            count: item.count || 0,
            cover_url: coverUrl,
            cover_preview: coverPreview,
            display_cover: displayCover,
            hasCustomCover,
            roleText: this.mapRoleText(item.role_scope)
          };
        });
        this.setData({ categories });
      } else {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load categories error", err);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  mapRoleText(roleScope) {
    if (roleScope === "employee") return "员工可见";
    if (roleScope === "manager") return "店长可见";
    if (roleScope === "both") return "全部可见";
    return roleScope || "";
  },

  handleChangeCover(e) {
    const id = Number(e.currentTarget.dataset.id);
    if (!id) return;
    wx.chooseImage({
      count: 1,
      sizeType: ["compressed", "original"],
      sourceType: ["album", "camera"],
      success: (res) => {
        const filePath = res.tempFilePaths && res.tempFilePaths[0];
        if (!filePath) {
          wx.showToast({ title: "图片路径无效", icon: "none" });
          return;
        }
        this.uploadCover(id, filePath);
      },
      fail: (err) => {
        if (err && err.errMsg !== "chooseImage:fail cancel") {
          wx.showToast({ title: "选择图片失败", icon: "none" });
        }
      }
    });
  },

  async uploadCover(id, filePath) {
    wx.showLoading({ title: "上传中..." });
    try {
      const res = await api.file.upload(filePath);
      if (res.code === 200 && res.data && res.data.path) {
        const coverPath = res.data.path;
        const saveRes = await api.admin.updateCategory(id, { cover_url: coverPath });
        if (saveRes.code === 200) {
          wx.showToast({ title: "已更新封面", icon: "success" });
          this.loadCategories();
        } else {
          wx.showToast({ title: saveRes.message || "保存失败", icon: "none" });
        }
      } else {
        wx.showToast({ title: res.message || "上传失败", icon: "none" });
      }
    } catch (err) {
      console.error("upload cover error", err);
      wx.showToast({ title: "操作失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  },

  async handleClearCover(e) {
    const id = Number(e.currentTarget.dataset.id);
    if (!id) return;
    wx.showModal({
      title: "确认操作",
      content: "确定要清空该分类封面，使用默认图吗？",
      success: async (modalRes) => {
        if (!modalRes.confirm) return;
        wx.showLoading({ title: "处理中..." });
        try {
          const res = await api.admin.updateCategory(id, { cover_url: "" });
          if (res.code === 200) {
            wx.showToast({ title: "已恢复默认", icon: "success" });
            this.loadCategories();
          } else {
            wx.showToast({ title: res.message || "操作失败", icon: "none" });
          }
        } catch (err) {
          console.error("clear cover error", err);
          wx.showToast({ title: "操作失败", icon: "none" });
        } finally {
          wx.hideLoading();
        }
      }
    });
  }
});
