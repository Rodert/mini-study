Page({
  data: {
    exam: null,
    result: null
  },

  onLoad(options) {
    // 尝试从事件通道获取数据
    try {
      const channel = this.getOpenerEventChannel && this.getOpenerEventChannel();
      if (channel && typeof channel.once === 'function') {
        channel.once("examResult", (payload) => {
          this.mergeResult(payload?.exam, payload?.result);
        });
      }
    } catch (err) {
      console.warn("getOpenerEventChannel failed", err);
      // 如果事件通道不可用，显示提示
      wx.showToast({
        title: "数据加载失败，请重试",
        icon: "none"
      });
      setTimeout(() => {
        wx.redirectTo({
          url: '/pages/exams/list/index'
        });
      }, 1500);
    }
  },

  mergeResult(exam, result) {
    if (!exam || !result) {
      this.setData({ exam, result });
      return;
    }
    const questionMap = new Map();
    (exam.questions || []).forEach((question) => {
      const optionMap = new Map();
      (question.options || []).forEach((option) => {
        optionMap.set(option.id, option);
      });
      questionMap.set(question.id, { ...question, optionMap });
    });

    const answers = (result.answers || []).map((answer) => {
      const question = questionMap.get(answer.question_id) || {};
      const optionMap = question.optionMap || new Map();
      const buildLabel = (ids = []) =>
        ids
          .map((id) => optionMap.get(id)?.label || "")
          .filter(Boolean)
          .join("、");
      return {
        ...answer,
        stem: question.stem || answer.stem,
        type: question.type || answer.type,
        selectedLabels: buildLabel(answer.selected_option_ids),
        correctLabels: buildLabel(answer.correct_option_ids)
      };
    });

    this.setData({
      exam,
      result: {
        ...result,
        answers
      }
    });
  },

  handleBack() {
    // 返回到考试中心（考试列表页面）
    wx.redirectTo({
      url: '/pages/exams/list/index'
    });
  }
});


