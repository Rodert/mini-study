const api = require("../../../services/api");

Page({
  data: {
    categoryId: null,
    categoryName: "",
    courses: [],
    loading: false
  },

  onLoad(options) {
    const { categoryId, name } = options;
    this.setData({
      categoryId: Number(categoryId),
      categoryName: name || ""
    });
    this.loadCourses();
  },

  formatDuration(seconds) {
    if (!seconds || seconds <= 0) {
      return "--";
    }
    const mins = Math.round(seconds / 60);
    if (mins < 60) {
      return `${mins} 分钟`;
    }
    const hours = Math.floor(mins / 60);
    const remains = mins % 60;
    return `${hours} 小时 ${remains} 分`;
  },

  async loadCourses() {
    if (!this.data.categoryId) return;
    this.setData({ loading: true });
    try {
      const res = await api.content.listPublished({
        category_id: this.data.categoryId
      });
      if (res.code === 200) {
        let courses = (res.data || []).map((item) => ({
          id: item.id,
          title: item.title,
          summary: item.summary,
          type: item.type,
          cover: item.cover_url ? api.buildFileUrl(item.cover_url) : "",
          duration: this.formatDuration(item.duration_seconds)
        }));

        // 加载当前用户的学习进度列表，并合并到课程数据中
        try {
          const progressRes = await api.learning.listProgress();
          if (progressRes.code === 200 && Array.isArray(progressRes.data)) {
            const progressMap = {};
            (progressRes.data || []).forEach((p) => {
              const contentId = p.content_id;
              if (!contentId) {
                return;
              }
              progressMap[contentId] = p;
            });

            courses = courses.map((course) => {
              const progress = progressMap[course.id];
              let progressStatus = "not_started";
              let progressPercent = 0;
              let progressText = "";

              if (course.type === "video") {
                if (progress) {
                  if (typeof progress.progress === "number") {
                    progressPercent = progress.progress;
                  }
                  if (progress.status) {
                    progressStatus = progress.status;
                  }
                }

                if (progressStatus === "completed") {
                  progressText = "已完成";
                } else if (progressPercent > 0) {
                  progressText = `已学习 ${progressPercent}%`;
                } else {
                  progressText = "未开始";
                }
              }

              course.progressStatus = progressStatus;
              course.progressPercent = progressPercent;
              course.progressText = progressText;
              return course;
            });
          }
        } catch (e) {
          console.error("load learning progress list error", e);
        }
        this.setData({ courses });
      } else {
        wx.showToast({ title: res.message || "课程加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load courses error", err);
      wx.showToast({ title: "课程加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  handleSelectCourse(e) {
    const { id } = e.currentTarget.dataset;
    if (!id) return;
    wx.navigateTo({
      url: `/pages/learning/detail/index?id=${id}`
    });
  }
});

