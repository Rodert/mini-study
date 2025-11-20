const api = require("../../../../services/api");
const app = getApp();

Page({
  data: {
    examId: null,
    isEdit: false,
    submitting: false,
    form: {
      title: "",
      description: "",
      status: "draft",
      target_role: "employee",
      time_limit_minutes: 0,
      pass_score: 60,
      questions: []
    },
    statusOptions: [
      { label: "草稿", value: "draft" },
      { label: "已发布", value: "published" },
      { label: "已归档", value: "archived" }
    ],
    statusIndex: 0,
    roleOptions: [
      { label: "员工", value: "employee" },
      { label: "店长", value: "manager" },
      { label: "全部", value: "all" }
    ],
    roleIndex: 0
  },

  onLoad(options) {
    const user = app.globalData.user || wx.getStorageSync("user");
    if (!user || user.role !== "admin") {
      wx.showToast({ title: "仅管理员可访问", icon: "none" });
      setTimeout(() => wx.navigateBack(), 800);
      return;
    }

    if (options.id) {
      this.setData({ examId: Number(options.id), isEdit: true });
      this.loadExam();
    }
  },

  async loadExam() {
    if (!this.data.examId) return;
    wx.showLoading({ title: "加载中..." });
    try {
      const res = await api.admin.getExamDetail(this.data.examId);
      if (res.code === 200) {
        const exam = res.data;
        const statusIndex = this.data.statusOptions.findIndex(
          (opt) => opt.value === exam.status
        );
        const roleIndex = this.data.roleOptions.findIndex(
          (opt) => opt.value === exam.target_role
        );
        this.setData({
          form: {
            title: exam.title || "",
            description: exam.description || "",
            status: exam.status || "draft",
            target_role: exam.target_role || "employee",
            time_limit_minutes: exam.time_limit_minutes || 0,
            pass_score: exam.pass_score || 60,
            questions: (exam.questions || []).map((q, qIdx) => ({
              id: q.id || Date.now() + qIdx,
              type: q.type,
              stem: q.stem,
              score: q.score,
              analysis: q.analysis || "",
              options: (q.options || []).map((opt, oIdx) => ({
                id: opt.id || Date.now() + qIdx * 1000 + oIdx,
                label: opt.label,
                content: opt.content,
                is_correct: opt.is_correct || false,
                sort_order: opt.sort_order || oIdx
              }))
            }))
          },
          statusIndex: statusIndex >= 0 ? statusIndex : 0,
          roleIndex: roleIndex >= 0 ? roleIndex : 0
        });
      } else {
        wx.showToast({ title: res.message || "加载失败", icon: "none" });
      }
    } catch (err) {
      console.error("load exam error", err);
      wx.showToast({ title: "加载失败", icon: "none" });
    } finally {
      wx.hideLoading();
    }
  },

  handleInput(e) {
    const { field } = e.currentTarget.dataset;
    if (!field) return;
    this.setData({
      [`form.${field}`]: e.detail.value
    });
  },

  handleNumberInput(e) {
    const { field } = e.currentTarget.dataset;
    if (!field) return;
    const value = Number(e.detail.value) || 0;
    this.setData({
      [`form.${field}`]: value
    });
  },

  handleStatusChange(e) {
    const index = Number(e.detail.value);
    const option = this.data.statusOptions[index];
    if (!option) return;
    this.setData({
      statusIndex: index,
      "form.status": option.value
    });
  },

  handleRoleChange(e) {
    const index = Number(e.detail.value);
    const option = this.data.roleOptions[index];
    if (!option) return;
    this.setData({
      roleIndex: index,
      "form.target_role": option.value
    });
  },

  addQuestion() {
    const baseId = Date.now();
    const newQuestion = {
      id: baseId + Math.random(), // 添加唯一ID
      type: "single",
      stem: "",
      score: 5,
      analysis: "",
      options: [
        { id: baseId + Math.random() + 1, label: "A", content: "", is_correct: false, sort_order: 0 },
        { id: baseId + Math.random() + 2, label: "B", content: "", is_correct: false, sort_order: 1 }
      ]
    };
    this.setData({ 
      "form.questions": [...this.data.form.questions, newQuestion]
    });
  },

  removeQuestion(e) {
    const { index } = e.currentTarget.dataset;
    const questions = this.data.form.questions.filter((q, idx) => idx !== Number(index));
    this.setData({ "form.questions": questions });
  },

  handleQuestionInput(e) {
    const { index, field } = e.currentTarget.dataset;
    const questions = this.data.form.questions.map((q, idx) => {
      if (idx === Number(index)) {
        return { ...q, [field]: e.detail.value };
      }
      return { ...q };
    });
    this.setData({ "form.questions": questions });
  },

  handleQuestionNumberInput(e) {
    const { index, field } = e.currentTarget.dataset;
    const questions = this.data.form.questions.map((q, idx) => {
      if (idx === Number(index)) {
        return { ...q, [field]: Number(e.detail.value) || 0 };
      }
      return { ...q };
    });
    this.setData({ "form.questions": questions });
  },

  handleQuestionTypeChange(e) {
    const { index } = e.currentTarget.dataset;
    const type = e.detail.value;
    const questions = this.data.form.questions.map((q, idx) => {
      if (idx === Number(index)) {
        const newQ = { ...q, type };
        // 如果是单选题，确保只有一个正确答案
        if (type === "single") {
          const hasCorrect = newQ.options.some((opt) => opt.is_correct);
          if (hasCorrect) {
            // 保留第一个正确答案，其他设为false
            let foundFirst = false;
            newQ.options = newQ.options.map((opt) => {
              if (opt.is_correct && !foundFirst) {
                foundFirst = true;
                return { ...opt };
              }
              return { ...opt, is_correct: false };
            });
          }
        }
        return newQ;
      }
      return { ...q };
    });
    this.setData({ "form.questions": questions });
  },

  addOption(e) {
    const { index } = e.currentTarget.dataset;
    const questions = this.data.form.questions.map((q, idx) => {
      if (idx === Number(index)) {
        const newQ = { ...q };
        const labels = ["A", "B", "C", "D", "E", "F", "G", "H"];
        const usedLabels = newQ.options.map((opt) => opt.label);
        const nextLabel = labels.find((l) => !usedLabels.includes(l)) || String.fromCharCode(65 + usedLabels.length);
        newQ.options = [...newQ.options, {
          id: Date.now() + Math.random(), // 添加唯一ID用于key
          label: nextLabel,
          content: "",
          is_correct: false,
          sort_order: newQ.options.length
        }];
        return newQ;
      }
      return { ...q };
    });
    this.setData({ "form.questions": questions });
  },

  removeOption(e) {
    const { qIndex, oIndex } = e.currentTarget.dataset;
    const questions = this.data.form.questions.map((q, qIdx) => {
      if (qIdx === Number(qIndex)) {
        const newQ = { ...q };
        newQ.options = q.options.filter((opt, optIdx) => optIdx !== Number(oIndex));
        return newQ;
      }
      return { ...q };
    });
    this.setData({ "form.questions": questions });
  },

  handleOptionInput(e) {
    const { qIndex, oIndex, field } = e.currentTarget.dataset;
    const questions = this.data.form.questions.map((q, qIdx) => {
      if (qIdx === Number(qIndex)) {
        const newQ = { ...q };
        newQ.options = q.options.map((opt, optIdx) => {
          if (optIdx === Number(oIndex)) {
            return { ...opt, [field]: e.detail.value };
          }
          return { ...opt };
        });
        return newQ;
      }
      return { ...q };
    });
    this.setData({ "form.questions": questions });
  },

  toggleCorrect(e) {
    const { qIndex, oIndex } = e.currentTarget.dataset;
    const qIdx = Number(qIndex);
    const oIdx = Number(oIndex);
    
    if (isNaN(qIdx) || isNaN(oIdx)) {
      console.error("Invalid index", qIndex, oIndex);
      return;
    }

    // 使用深拷贝确保数据变化能被检测到
    const questions = this.data.form.questions.map((q, idx) => {
      if (idx === qIdx) {
        const newQ = { ...q };
        newQ.options = q.options.map((opt, optIdx) => {
          if (optIdx === oIdx) {
            // 当前选项
            if (q.type === "single") {
              // 单选题：先全部设为false，再设置当前为true
              return { ...opt, is_correct: true };
            } else {
              // 多选题：切换状态
              return { ...opt, is_correct: !opt.is_correct };
            }
          } else if (q.type === "single") {
            // 单选题：其他选项设为false
            return { ...opt, is_correct: false };
          } else {
            // 多选题：其他选项不变
            return { ...opt };
          }
        });
        return newQ;
      }
      return { ...q };
    });
    
    // 使用完整路径更新，确保视图刷新
    this.setData({ 
      "form.questions": questions 
    });
  },

  validateForm() {
    const { form } = this.data;
    if (!form.title || !form.title.trim()) {
      wx.showToast({ title: "请输入考试标题", icon: "none" });
      return false;
    }
    if (form.pass_score <= 0) {
      wx.showToast({ title: "及格分必须大于0", icon: "none" });
      return false;
    }
    if (!form.questions || form.questions.length === 0) {
      wx.showToast({ title: "请至少添加一道题目", icon: "none" });
      return false;
    }

    for (let i = 0; i < form.questions.length; i++) {
      const q = form.questions[i];
      if (!q.stem || !q.stem.trim()) {
        wx.showToast({ title: `第${i + 1}题：请输入题目内容`, icon: "none" });
        return false;
      }
      if (q.score <= 0) {
        wx.showToast({ title: `第${i + 1}题：分数必须大于0`, icon: "none" });
        return false;
      }
      if (!q.options || q.options.length < 2) {
        wx.showToast({ title: `第${i + 1}题：至少需要2个选项`, icon: "none" });
        return false;
      }
      for (let j = 0; j < q.options.length; j++) {
        const opt = q.options[j];
        if (!opt.content || !opt.content.trim()) {
          wx.showToast({ title: `第${i + 1}题选项${opt.label}：请输入选项内容`, icon: "none" });
          return false;
        }
      }
      const correctCount = q.options.filter((opt) => opt.is_correct).length;
      if (correctCount === 0) {
        wx.showToast({ title: `第${i + 1}题：请设置正确答案`, icon: "none" });
        return false;
      }
      if (q.type === "single" && correctCount > 1) {
        wx.showToast({ title: `第${i + 1}题：单选题只能有一个正确答案`, icon: "none" });
        return false;
      }
    }

    // 计算总分
    const totalScore = form.questions.reduce((sum, q) => sum + q.score, 0);
    if (form.pass_score > totalScore) {
      wx.showToast({ title: `及格分(${form.pass_score})不能超过总分(${totalScore})`, icon: "none" });
      return false;
    }

    return true;
  },

  async handleSubmit() {
    if (!this.validateForm()) return;
    if (this.data.submitting) return;

    this.setData({ submitting: true });
    wx.showLoading({ title: "提交中..." });

    try {
      const payload = {
        title: this.data.form.title.trim(),
        description: this.data.form.description.trim(),
        status: this.data.form.status,
        target_role: this.data.form.target_role,
        time_limit_minutes: this.data.form.time_limit_minutes,
        pass_score: this.data.form.pass_score,
        questions: this.data.form.questions.map((q) => ({
          type: q.type,
          stem: q.stem.trim(),
          score: q.score,
          analysis: q.analysis ? q.analysis.trim() : "",
          options: q.options.map((opt) => ({
            label: opt.label,
            content: opt.content.trim(),
            is_correct: opt.is_correct,
            sort_order: opt.sort_order
          }))
        }))
      };

      let res;
      if (this.data.isEdit) {
        res = await api.admin.updateExam(this.data.examId, payload);
      } else {
        res = await api.admin.createExam(payload);
      }

      if (res.code === 200) {
        wx.showToast({
          title: this.data.isEdit ? "更新成功" : "创建成功",
          icon: "success"
        });
        setTimeout(() => {
          wx.navigateBack();
        }, 1500);
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
  }
});

