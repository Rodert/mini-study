const api = require("../../../../services/api");
const app = getApp();

Page({
  data: {
    userId: null,
    user: {},
    totalPoints: 0,
    transactions: [],
    page: 1,
    pageSize: 20,
    total: 0,
    hasMore: true,
    loading: false
  },

  onLoad(options) {
    const admin = app.globalData.user || wx.getStorageSync("user");
    if (!admin || admin.role !== "admin") {
      wx.showToast({ title: "无权限访问", icon: "none" });
      setTimeout(() => {
        wx.navigateBack();
      }, 1000);
      return;
    }

    const { userId } = options;
    if (!userId) {
      wx.showToast({ title: "参数错误", icon: "none" });
      wx.navigateBack();
      return;
    }

    this.setData({ userId: parseInt(userId) });
    this.loadData();
  },

  async loadData(refresh = false) {
    if (this.data.loading) return;

    const page = refresh ? 1 : this.data.page;
    if (!refresh && !this.data.hasMore) return;

    this.setData({ loading: true });

    try {
      const params = {
        page,
        page_size: this.data.pageSize
      };

      const res = await api.admin.getUserPoints(this.data.userId, params);

      console.log("积分明细API响应:", JSON.stringify(res, null, 2));

      if (res.code !== 200) {
        wx.showToast({ title: res.message || "加载数据失败", icon: "none" });
        this.setData({ loading: false });
        return;
      }

      const data = res.data || {};
      const user = data.user || {};
      const transactions = data.transactions || [];
      const pagination = data.pagination || {};

      const newTransactions = (Array.isArray(transactions) ? transactions : []).map((item) => {
        const createdAt = item.created_at || item.createdAt || "";
        const source = item.source || "";
        
        // 格式化日期
        let formattedDate = "";
        if (createdAt) {
          try {
            const date = new Date(createdAt);
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1);
            const day = String(date.getDate());
            const hours = String(date.getHours());
            const minutes = String(date.getMinutes());
            // 兼容性处理：padStart 可能不支持，手动补零
            const pad = (str) => str.length < 2 ? "0" + str : str;
            formattedDate = `${year}-${pad(month)}-${pad(day)} ${pad(hours)}:${pad(minutes)}`;
          } catch (e) {
            formattedDate = createdAt;
          }
        }
        
        // 格式化来源
        const sourceMap = {
          content_completion: "完成学习",
          exam_pass: "通过考试",
          admin_adjust: "管理员调整",
          system_reward: "系统奖励"
        };
        const formattedSource = sourceMap[source] || source || "未知";
        
        return {
          id: item.id,
          change: item.change || 0,
          source: source,
          sourceText: formattedSource,
          description: item.description || "",
          memo: item.memo || "",
          createdAt: createdAt,
          createdAtText: formattedDate
        };
      });

      const total = pagination.total || pagination.Total || 0;
      const hasMore = page * this.data.pageSize < total;

      this.setData({
        user: {
          id: user.id,
          name: user.name || "",
          workNo: user.work_no || user.workNo || "",
          phone: user.phone || "",
          role: user.role || ""
        },
        totalPoints: data.total_points || data.totalPoints || 0,
        transactions: refresh ? newTransactions : [...this.data.transactions, ...newTransactions],
        page: page,
        total: total,
        hasMore: hasMore,
        loading: false
      });
    } catch (err) {
      console.error("load points detail error", err);
      console.error("错误详情:", JSON.stringify(err, null, 2));
      wx.showToast({ title: err.message || "加载数据失败", icon: "none" });
      this.setData({ loading: false });
    }
  },

  goBack() {
    wx.navigateBack();
  },

  onReachBottom() {
    // 加载更多
    if (this.data.hasMore && !this.data.loading) {
      this.setData({ page: this.data.page + 1 });
      this.loadData();
    }
  },

  onPullDownRefresh() {
    // 下拉刷新
    this.setData({
      page: 1,
      hasMore: true
    });
    this.loadData(true).then(() => {
      wx.stopPullDownRefresh();
    }).catch(() => {
      wx.stopPullDownRefresh();
    });
  },

});

