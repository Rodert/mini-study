const api = require("../../../services/api");

Page({
  data: {
    examId: null,
    exam: null,
    loading: false,
    answers: {},
    optionChecked: {}, // 存储每个选项的选中状态: { "questionId_optionId": true/false }
    submitting: false,
    readyToSubmit: false
  },

  onLoad(options) {
    const { id } = options;
    if (!id) {
      wx.showToast({ title: "缺少考试ID", icon: "none" });
      setTimeout(() => {
        this.handleBack();
      }, 1500);
      return;
    }
    this.setData({ examId: Number(id) });
    this.fetchExam();
  },

  async fetchExam() {
    if (!this.data.examId) return;
    this.setData({ loading: true });
    try {
      const res = await api.exam.getDetail(this.data.examId);
      if (res.code === 200) {
        // 初始化选项选中状态
        const optionChecked = {};
        if (res.data && res.data.questions) {
          res.data.questions.forEach(q => {
            q.options.forEach(opt => {
              optionChecked[`${q.id}_${opt.id}`] = false;
            });
          });
        }
        this.setData({ 
          exam: res.data,
          optionChecked
        }, () => {
          this.updateReadyState();
        });
      } else {
        wx.showToast({ title: res.message || "获取考试信息失败", icon: "none" });
      }
    } catch (err) {
      console.error("fetch exam detail error", err);
      wx.showToast({ title: "获取考试信息失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  handleSingleChange(event) {
    const { qid } = event.currentTarget.dataset;
    const value = Number(event.detail.value);
    const questionId = Number(qid);
    
    // 更新答案
    const newAnswers = { ...this.data.answers };
    newAnswers[questionId] = value ? [value] : [];
    
    // 更新选项选中状态
    const newOptionChecked = { ...this.data.optionChecked };
    // 先清除该题所有选项的选中状态
    const question = this.data.exam?.questions?.find(q => q.id === questionId);
    if (question) {
      question.options.forEach(opt => {
        newOptionChecked[`${questionId}_${opt.id}`] = false;
      });
    }
    // 设置当前选中项
    if (value) {
      newOptionChecked[`${questionId}_${value}`] = true;
    }
    
    this.setData(
      {
        answers: newAnswers,
        optionChecked: newOptionChecked
      },
      () => {
        this.updateReadyState();
      }
    );
  },

  handleMultiChange(event) {
    const { qid } = event.currentTarget.dataset;
    const values = (event.detail.value || []).map((val) => Number(val));
    const questionId = Number(qid);
    
    // 更新答案
    const newAnswers = { ...this.data.answers };
    newAnswers[questionId] = values;
    
    // 更新选项选中状态
    const newOptionChecked = { ...this.data.optionChecked };
    const question = this.data.exam?.questions?.find(q => q.id === questionId);
    if (question) {
      question.options.forEach(opt => {
        newOptionChecked[`${questionId}_${opt.id}`] = values.includes(opt.id);
      });
    }
    
    this.setData(
      {
        answers: newAnswers,
        optionChecked: newOptionChecked
      },
      () => {
        this.updateReadyState();
      }
    );
  },

  canSubmit() {
    const exam = this.data.exam;
    if (!exam || !exam.questions || !exam.questions.length) {
      return false;
    }
    return exam.questions.every((q) => {
      const selected = this.data.answers[q.id];
      return Array.isArray(selected) && selected.length > 0;
    });
  },

  async handleSubmit() {
    if (!this.canSubmit()) {
      wx.showToast({ title: "请完成所有题目", icon: "none" });
      return;
    }
    if (this.data.submitting) return;

    const payload = {
      answers: this.data.exam.questions.map((question) => ({
        question_id: question.id,
        option_ids: this.data.answers[question.id]
      }))
    };

    this.setData({ submitting: true });
    wx.showLoading({ title: "提交中..." });
    try {
      const res = await api.exam.submit(this.data.examId, payload);
      if (res.code === 200) {
        this.jumpToResult(res.data);
      } else {
        wx.showToast({ title: res.message || "提交失败", icon: "none" });
      }
    } catch (err) {
      console.error("submit exam error", err);
      wx.showToast({ title: "提交失败", icon: "none" });
    } finally {
      wx.hideLoading();
      this.setData({ submitting: false });
    }
  },

  jumpToResult(result) {
    const exam = this.data.exam;
    wx.navigateTo({
      url: "/pages/exams/result/index",
      success: (navRes) => {
        navRes.eventChannel.emit("examResult", {
          exam,
          result
        });
      }
    });
  },

  updateReadyState() {
    const ready = this.canSubmit();
    if (ready !== this.data.readyToSubmit) {
      this.setData({ readyToSubmit: ready });
    }
  },

  handleBack() {
    // 返回到考试中心（考试列表页面）
    wx.redirectTo({
      url: '/pages/exams/list/index'
    });
  }
});


