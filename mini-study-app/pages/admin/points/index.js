const api = require("../../../services/api");
const app = getApp();

Page({
  data: {
    points: [],
    filteredPoints: [],
    searchText: "",
    page: 1,
    pageSize: 20,
    total: 0,
    loading: false,
    hasMore: true
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

    this.loadData();
  },

  onShow() {
    // 从详情页返回时不自动刷新，用户可手动下拉刷新
  },

  async loadData(refresh = false) {
    if (this.data.loading) return;

    const page = refresh ? 1 : this.data.page;
    if (!refresh && !this.data.hasMore) return;

    this.setData({ loading: true });

    try {
      // 构建请求参数，过滤掉空值
      const params = {
        page,
        page_size: this.data.pageSize
      };
      
      // 只有当搜索文本不为空时才添加 keyword 参数
      if (this.data.searchText && this.data.searchText.trim()) {
        params.keyword = this.data.searchText.trim();
      }

      const res = await api.admin.listPoints(params);

      console.log("积分列表API响应:", JSON.stringify(res, null, 2));

      if (res.code !== 200) {
        wx.showToast({ title: res.message || "加载数据失败", icon: "none" });
        this.setData({ loading: false });
        return;
      }

      // 确保正确访问嵌套的数据结构
      const data = res.data || {};
      const items = data.items || data.Items || [];
      const pagination = data.pagination || data.Pagination || {};

      console.log("解析后的items数量:", items.length);
      console.log("解析后的items:", items);
      console.log("解析后的pagination:", pagination);

      const newPoints = (Array.isArray(items) ? items : []).map((item) => ({
        id: item.id,
        name: item.name || "",
        workNo: item.work_no || item.workNo || "",
        phone: item.phone || "",
        role: item.role || "",
        status: item.status !== undefined ? item.status : true,
        points: item.points || 0
      }));

      console.log("映射后的newPoints数量:", newPoints.length);

      const total = pagination.total || pagination.Total || 0;
      const hasMore = page * this.data.pageSize < total;

      this.setData({
        points: refresh ? newPoints : [...this.data.points, ...newPoints],
        filteredPoints: refresh ? newPoints : [...this.data.filteredPoints, ...newPoints],
        page: page,
        total: total,
        hasMore: hasMore,
        loading: false
      });

      console.log("设置后的filteredPoints数量:", this.data.filteredPoints.length);
    } catch (err) {
      console.error("load points error", err);
      console.error("错误详情:", JSON.stringify(err, null, 2));
      wx.showToast({ title: err.message || "加载数据失败", icon: "none" });
      this.setData({ loading: false });
    }
  },

  handleSearch(e) {
    const searchText = e.detail.value;
    this.setData({ 
      searchText,
      page: 1,
      hasMore: true
    });

    // 重新加载数据（后端会进行搜索）
    this.loadData(true);
  },

  goPointDetail(e) {
    const { userId } = e.currentTarget.dataset;
    if (!userId) {
      wx.showToast({ title: "用户ID不存在", icon: "none" });
      return;
    }
    wx.navigateTo({
      url: `/pages/admin/points/detail/index?userId=${userId}`
    });
  },

  goBack() {
    wx.navigateBack();
  },

  onReachBottom() {
    // 加载更多
    if (this.data.hasMore && !this.data.loading && !this.data.searchText) {
      this.setData({ page: this.data.page + 1 });
      this.loadData();
    }
  },

  onPullDownRefresh() {
    // 下拉刷新
    this.setData({
      page: 1,
      hasMore: true,
      searchText: ""
    });
    this.loadData(true).then(() => {
      wx.stopPullDownRefresh();
    }).catch(() => {
      wx.stopPullDownRefresh();
    });
  }
});
