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
        let courses = (res.data || []).map((item) => {
          const normalizedPath = item.file_path ? item.file_path.replace(/\\/g, "/") : "";
          return {
            id: item.id,
            title: item.title,
            summary: item.summary,
            type: item.type,
            cover: item.cover_url ? api.buildFileUrl(item.cover_url) : "",
            duration: this.formatDuration(item.duration_seconds),
            mediaUrl: item.file_path ? api.buildFileUrl(item.file_path) : "",
            rawFilePath: item.file_path || "",
            fileName: normalizedPath ? normalizedPath.split("/").pop() : ""
          };
        });

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
    const courseId = Number(id);
    const courses = this.data.courses || [];
    const course = courses.find((c) => c.id === courseId);

    if (course && course.type === "doc") {
      this.openDocumentFromList(course);
      return;
    }

    wx.navigateTo({
      url: `/pages/learning/detail/index?id=${courseId}`
    });
  },

  async openDocumentFromList(course) {
    if (!course || !course.id) return;
    const mediaUrl = course.mediaUrl;
    if (!mediaUrl) {
      wx.showToast({ title: "文件地址缺失", icon: "none" });
      return;
    }

    wx.showLoading({ title: "获取文档..." });
    try {
      const downloadRes = await new Promise((resolve, reject) => {
        wx.downloadFile({
          url: mediaUrl,
          success: (res) => {
            if (res.statusCode === 200 && res.tempFilePath) {
              resolve(res);
            } else {
              reject(new Error("download document failed"));
            }
          },
          fail: reject
        });
      });

      await new Promise((resolve, reject) => {
        wx.openDocument({
          filePath: downloadRes.tempFilePath,
          showMenu: true,
          success: () => {
            resolve();
          },
          fail: (err) => {
            reject(err);
          }
        });
      });

      try {
        await api.learning.updateProgress({
          content_id: Number(course.id),
          video_position: 0
        });
      } catch (err) {
        console.error("update document progress error", err);
      }
    } catch (err) {
      console.error("open document from list error", err);
      wx.showToast({ title: "打开失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  }
});

