const api = require("../../../services/api");

Page({
  data: {
    courseId: null,
    course: {},
    progress: null,
    statusText: "未开始",
    docLoading: false,
    docTempPath: "",
    docSourceUrl: "",
    videoContext: null,
    lastAllowedPosition: 0, // 记录最后允许的播放位置（秒）
    lastUpdateTime: 0 // 上次上报进度的时间戳
  },

  onLoad(options) {
    const { id } = options;
    this.setData({ courseId: Number(id) });
    this.loadCourse();
  },

  onReady() {
    // 获取视频上下文
    const videoContext = wx.createVideoContext('courseVideo', this);
    this.setData({
      videoContext: videoContext
    });
    // 首次加载时，如果有进度，设置播放位置（仅首次，不限制拖动）
    if (this.data.lastAllowedPosition > 0) {
      setTimeout(() => {
        if (videoContext && this.data.lastAllowedPosition > 0) {
          videoContext.seek(this.data.lastAllowedPosition);
        }
      }, 500);
    }
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

  getStatusText(status) {
    switch (status) {
      case "completed":
        return "已完成";
      case "in_progress":
        return "学习中";
      default:
        return "未开始";
    }
  },

  async loadCourse() {
    if (!this.data.courseId) return;
    try {
      const res = await api.content.getDetail(this.data.courseId);
      if (res.code === 200) {
        const detail = res.data || {};
        const mediaUrl = api.buildFileUrl(detail.file_path);
        const normalizedPath = detail.file_path ? detail.file_path.replace(/\\/g, "/") : "";
        const course = {
          id: detail.id,
          title: detail.title,
          type: detail.type,
          summary: detail.summary,
          cover: detail.cover_url,
          duration: this.formatDuration(detail.duration_seconds),
          durationSeconds: detail.duration_seconds || 0, // 保存原始时长（秒）
          mediaUrl,
          rawFilePath: detail.file_path,
          fileName: normalizedPath ? normalizedPath.split("/").pop() : "",
          content: detail.summary || ""
        };
        this.setData({
          course,
          docTempPath:
            this.data.docSourceUrl === mediaUrl ? this.data.docTempPath : "",
          docSourceUrl: this.data.docSourceUrl === mediaUrl ? this.data.docSourceUrl : ""
        });
        this.loadProgress();
      } else {
        wx.showToast({ title: res.message || "获取课程失败", icon: "none" });
      }
    } catch (err) {
      console.error("load course detail error", err);
      wx.showToast({ title: "获取课程失败", icon: "none" });
    }
  },

  async loadProgress() {
    try {
      const res = await api.learning.getProgress(this.data.courseId);
      if (res.code === 200 && res.data) {
        const videoPosition = res.data?.video_position || 0;
        const status = res.data?.status || "not_started";
        this.setData({
          progress: res.data,
          statusText: this.getStatusText(status),
          lastAllowedPosition: videoPosition
        });
        // 首次加载时，如果有进度且视频上下文已初始化，设置视频播放位置（仅首次）
        // 不限制用户后续拖动
        if (videoPosition > 0 && this.data.videoContext) {
          // 延迟执行，确保视频已加载
          setTimeout(() => {
            if (this.data.videoContext) {
              this.data.videoContext.seek(videoPosition);
            }
          }, 500);
        }
      } else {
        // 没有学习记录，设置为未开始状态
        this.setData({ 
          progress: { status: "not_started" }, 
          statusText: "未开始", 
          lastAllowedPosition: 0 
        });
      }
    } catch (err) {
      console.error("load progress error", err);
      // 加载失败时也设置为未开始状态
      this.setData({ 
        progress: { status: "not_started" }, 
        statusText: "未开始", 
        lastAllowedPosition: 0 
      });
    }
  },

  handleOpenDocument() {
    if (this.data.docLoading) return;
    const url = this.data.course.mediaUrl;
    if (!url) {
      wx.showToast({ title: "文件地址缺失", icon: "none" });
      return;
    }
    
    // 触发开始学习（如果是首次打开，会创建学习记录）
    // 文档类型会在后端自动标记为完成
    const isNotStarted = !this.data.progress || 
                        this.data.progress.status === "not_started" || 
                        !this.data.progress.status;
    if (isNotStarted) {
      this.updateProgress(0, true); // 强制立即上报，不等待节流
    }
    
    if (this.data.docTempPath && this.data.docSourceUrl === url) {
      this.openDocument(this.data.docTempPath);
      return;
    }

    this.setData({ docLoading: true });
    wx.showLoading({ title: "获取文档..." });
    wx.downloadFile({
      url,
      success: (res) => {
        if (res.statusCode === 200 && res.tempFilePath) {
          this.setData({ docTempPath: res.tempFilePath, docSourceUrl: url });
          this.openDocument(res.tempFilePath);
        } else {
          wx.showToast({ title: "下载文档失败", icon: "none" });
        }
      },
      fail: (err) => {
        console.error("download document error", err);
        wx.showToast({ title: "下载文档失败", icon: "none" });
      },
      complete: () => {
        wx.hideLoading();
        this.setData({ docLoading: false });
      }
    });
  },

  openDocument(filePath) {
    if (!filePath) return;
    
    // 确保 courseId 已加载
    if (!this.data.courseId) {
      console.warn('courseId not ready, skip updateProgress');
      wx.showToast({ title: "课程信息加载中", icon: "none" });
      return;
    }
    
    // 文档打开时触发开始学习（后端会自动标记为完成）
    const isNotStarted = !this.data.progress || 
                        this.data.progress.status === "not_started" || 
                        !this.data.progress.status;
    if (isNotStarted) {
      this.updateProgress(0, true); // 强制立即上报
    }
    
    wx.openDocument({
      filePath,
      showMenu: true,
      success: () => {
        // 文档打开成功后，再次上报进度确保状态更新为已完成
        setTimeout(() => {
          this.updateProgress(0, true); // 强制立即上报
          // 刷新进度显示
          setTimeout(() => {
            this.loadProgress();
          }, 500);
        }, 500);
      },
      fail: (err) => {
        console.error("open document error", err);
        wx.showToast({ title: "打开失败", icon: "none" });
      }
    });
  },

  // 处理视频拖动事件
  handleVideoSeek(e) {
    // 完全允许用户自由拖动，不做任何限制
    // 观看时长记录由 handleTimeUpdate 和 updateProgress 处理
  },

  // 处理视频播放时间更新
  handleTimeUpdate(e) {
    const currentTime = e.detail.currentTime || 0;
    const lastPos = this.data.lastAllowedPosition;
    const duration = this.data.course?.durationSeconds || 0;
    
    // 如果当前播放位置大于已记录的最大位置，更新并上报进度
    // 允许用户自由拖动，但只记录最大观看进度（观看时长）
    if (currentTime > lastPos) {
      this.setData({ lastAllowedPosition: currentTime });
      // 定期上报进度（记录观看时长）
      this.updateProgress(currentTime);
      
      // 如果进度接近完成（>= 95%），提前上报确保标记完成
      if (duration > 0 && currentTime >= duration * 0.95 && 
          this.data.progress?.status !== "completed") {
        // 接近完成时，强制立即上报
        this.updateProgress(currentTime, true);
      }
    }
  },

  // 处理视频播放结束
  handleVideoEnded(e) {
    const duration = this.data.course?.durationSeconds || 0;
    console.log('Video ended, duration:', duration);
    
    // 视频播放结束时，上报最终位置（使用视频总时长）
    if (duration > 0) {
      // 确保上报的进度是100%
      this.updateProgress(duration, true);
      
      // 延迟刷新进度，确保后端处理完成
      setTimeout(() => {
        this.loadProgress();
        // 如果已标记为完成，显示提示
        if (this.data.progress?.status === "completed") {
          wx.showToast({ 
            title: "恭喜完成学习！", 
            icon: "success",
            duration: 2000
          });
        }
      }, 500);
    }
  },

  // 处理视频开始播放
  handleVideoPlay() {
    // 首次播放时触发开始学习（创建学习记录）
    // 确保 courseId 已加载
    if (!this.data.courseId) {
      console.warn('courseId not ready, skip updateProgress');
      return;
    }
    
    if (!this.data.progress || this.data.progress.status === "not_started" || !this.data.progress.status) {
      this.updateProgress(0, true); // 强制立即上报，不等待节流
    }
    // 不限制播放位置，允许用户自由拖动和播放
    // 观看时长记录由 handleTimeUpdate 处理
  },

  // 上报学习进度
  async updateProgress(videoPosition, force = false) {
    // 确保 videoPosition 是有效的数字
    const position = Math.max(0, Math.floor(Number(videoPosition) || 0));
    
    // 节流：每5秒上报一次（除非强制更新）
    const now = Date.now();
    if (!force && this.data.lastUpdateTime && now - this.data.lastUpdateTime < 5000) {
      return;
    }
    this.setData({ lastUpdateTime: now });

    if (!this.data.courseId) {
      console.error('courseId is required');
      return;
    }

    try {
      // 确保所有字段都是有效值
      const requestData = {
        content_id: Number(this.data.courseId),
        video_position: Number(position)
      };
      
      // 验证数据有效性
      if (isNaN(requestData.content_id) || requestData.content_id <= 0) {
        console.error('Invalid courseId:', this.data.courseId);
        return;
      }
      
      if (isNaN(requestData.video_position) || requestData.video_position < 0) {
        console.error('Invalid video_position:', position);
        return;
      }
      
      console.log('Updating progress:', requestData);
      
      const res = await api.learning.updateProgress(requestData);
      
      // 更新本地进度和状态
      if (res.code === 200 && res.data) {
        this.setData({
          progress: res.data,
          statusText: this.getStatusText(res.data.status)
        });
      }
    } catch (err) {
      console.error('update progress error', err);
      // 如果是验证错误，记录详细信息
      if (err.code === 400) {
        console.error('Validation error details:', {
          courseId: this.data.courseId,
          position: position,
          error: err.message
        });
      }
    }
  }
});

