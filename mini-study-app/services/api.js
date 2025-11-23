// API 服务配置
const config = require('./config');
const mockService = require('./mockService');
const API_BASE_URL = config.API_BASE_URL;
const USE_MOCK = config.USE_MOCK;
const API_ORIGIN = API_BASE_URL.replace(/\/api\/v1$/, '');

function buildFileUrl(path = '') {
  if (!path) return '';
  if (/^https?:\/\//i.test(path)) return path;

  let normalized = path.trim().replace(/\\/g, '/');
  if (!normalized) return '';

  const lower = normalized.toLowerCase();
  const uploadsIndex = lower.indexOf('/uploads/');
  if (uploadsIndex !== -1) {
    normalized = normalized.slice(uploadsIndex);
  } else {
    const storageIndex = lower.indexOf('storage/uploads');
    if (storageIndex !== -1) {
      normalized = normalized.slice(storageIndex + 'storage'.length);
    } else if (!normalized.startsWith('/')) {
      normalized = `/${normalized}`;
    }
  }

  if (!normalized.startsWith('/')) {
    normalized = `/${normalized}`;
  }

  return `${API_ORIGIN}${normalized}`;
}
// const API_BASE_URL = 'http://192.168.65.1:8080/api/v1';


// 获取 token
function getToken() {
  return wx.getStorageSync('token') || '';
}

// 通用请求方法
function request(options) {
  return new Promise((resolve, reject) => {
    const { url, method = 'GET', data = {}, header = {} } = options;
    
    // 过滤掉 undefined、null 和空字符串的值（GET 请求时）
    const cleanData = {};
    if (method === 'GET') {
      Object.keys(data).forEach(key => {
        const value = data[key];
        if (value !== undefined && value !== null && value !== '') {
          cleanData[key] = value;
        }
      });
    } else {
      Object.assign(cleanData, data);
    }
    
    // 添加 token
    const token = getToken();
    if (token) {
      header['Authorization'] = `Bearer ${token}`;
    }
    header['Content-Type'] = 'application/json';

    // 判断是否是登录相关接口（登录失败时不应该清除 token 和跳转）
    const isLoginEndpoint = url.includes('/users/login') || url.includes('/users/register');

    wx.request({
      url: `${API_BASE_URL}${url}`,
      method,
      data: cleanData,
      header,
      success: (res) => {
        console.log('请求响应:', url, res.statusCode, res.data);
        if (res.statusCode >= 200 && res.statusCode < 300) {
          // 成功响应，直接返回数据
          resolve(res.data);
        } else if (res.statusCode === 401) {
          // 401 未授权
          // 如果是登录/注册接口，不清除 token（因为还没有 token 或登录失败）
          if (!isLoginEndpoint) {
            handleUnauthorized();
          }
          // 返回错误响应，但格式统一
          const errorData = res.data || {};
          reject({
            code: errorData.code || 401,
            message: errorData.message || (isLoginEndpoint ? '工号或密码不正确' : '登录已过期，请重新登录'),
            data: errorData
          });
        } else {
          // 其他错误状态码
          const errorData = res.data || {};
          reject({
            code: errorData.code || res.statusCode,
            message: errorData.message || '请求失败',
            data: errorData
          });
        }
      },
      fail: (err) => {
        console.error('请求失败:', url, err);
        reject({
          code: -1,
          message: err.errMsg || '网络请求失败',
          err
        });
      }
    });
  });
}

// 处理 401 错误，跳转到登录页
function handleUnauthorized() {
  wx.removeStorageSync('token');
  wx.removeStorageSync('refresh_token');
  wx.removeStorageSync('user');
  wx.reLaunch({
    url: '/pages/login/index'
  });
}

module.exports = {
  $request: request,
  buildFileUrl,
  file: {
    upload(filePath) {
      const token = getToken();
      return new Promise((resolve, reject) => {
        wx.uploadFile({
          url: `${API_BASE_URL}/files/upload`,
          filePath,
          name: "file",
          header: token
            ? {
                Authorization: `Bearer ${token}`
              }
            : {},
          success: (res) => {
            let data;
            try {
              data = typeof res.data === "string" ? JSON.parse(res.data) : res.data;
            } catch (err) {
              reject({
                code: -1,
                message: "上传返回结果解析失败",
                raw: res.data,
                err
              });
              return;
            }

            if (res.statusCode >= 200 && res.statusCode < 300 && data.code === 200) {
              resolve(data);
            } else {
              reject({
                code: data.code || res.statusCode,
                message: data.message || "上传失败",
                data
              });
            }
          },
          fail: (err) => {
            reject({
              code: -1,
              message: err.errMsg || "上传失败",
              err
            });
          }
        });
      });
    }
  },
  // 用户相关
  user: {
    // 登录
    async login(data) {
      if (USE_MOCK) {
        const res = await mockService.login(data);
        // 保存 token
        if (res.code === 200 && res.data) {
          if (res.data.access_token) {
            wx.setStorageSync('token', res.data.access_token);
          }
          if (res.data.refresh_token) {
            wx.setStorageSync('refresh_token', res.data.refresh_token);
          }
          // Mock 模式下，登录成功后需要保存用户信息
          // 根据 work_no 从 mock 数据中找到用户
          const mock = require('../mock/mockData');
          const user = mock.users.find(u => u.work_no === data.work_no);
          if (user) {
            // 移除密码字段，保存用户信息
            const { password: _ignored, ...safeUser } = user;
            wx.setStorageSync('user', safeUser);
          }
        }
        return res;
      }
      
      try {
        const res = await request({
          url: '/users/login',
          method: 'POST',
          data: {
            work_no: data.work_no,
            password: data.password
          }
        });
        
        console.log('登录 API 响应:', res);
        
        // 保存 token
        if (res && res.code === 200 && res.data) {
          if (res.data.access_token) {
            wx.setStorageSync('token', res.data.access_token);
            console.log('Token 已保存');
          }
          if (res.data.refresh_token) {
            wx.setStorageSync('refresh_token', res.data.refresh_token);
          }
        }
        
        return res;
      } catch (error) {
        console.error('登录请求失败:', error);
        // 将错误转换为统一的响应格式
        if (error.code) {
          return error;
        }
        return {
          code: error.code || -1,
          message: error.message || '网络请求失败',
          data: null
        };
      }
    },
    // 获取当前用户信息
    getCurrentUser() {
      if (USE_MOCK) {
        // Mock 模式下需要从存储中获取用户信息
        const user = wx.getStorageSync('user');
        if (user) {
          return Promise.resolve({ code: 200, message: 'success', data: user });
        }
        return Promise.resolve({ code: 401, message: '未登录', data: null });
      }
      return request({
        url: '/users/me',
        method: 'GET'
      }).catch(error => {
        console.error('获取用户信息失败:', error);
        // 将错误转换为统一的响应格式
        if (error.code) {
          return error;
        }
        return {
          code: error.code || -1,
          message: error.message || '获取用户信息失败',
          data: null
        };
      });
    },
    // 获取店长列表
    getManagers() {
      if (USE_MOCK) {
        return mockService.fetchManagers();
      }
      return request({
        url: '/users/managers',
        method: 'GET'
      });
    },
    // 员工注册
    register(data) {
      if (USE_MOCK) {
        return mockService.register(data);
      }
      return request({
        url: '/users/register',
        method: 'POST',
        data
      });
    },
    // 更新个人资料
    updateProfile(data) {
      if (USE_MOCK) {
        const user = wx.getStorageSync('user');
        if (user && user.id) {
          return mockService.updateUserProfile(user.id, data);
        }
        return Promise.resolve({ code: 401, message: '未登录', data: null });
      }
      return request({
        url: '/users/me/profile',
        method: 'PATCH',
        data
      });
    }
  },
  
  // 内容相关
  content: {
    listCategories() {
      if (USE_MOCK) {
        const user = wx.getStorageSync('user');
        const role = user ? user.role : 'employee';
        return mockService.fetchCourseCategories(role);
      }
      return request({
        url: '/contents/categories',
        method: 'GET'
      });
    },
    listPublished(params = {}) {
      if (USE_MOCK) {
        if (params.category_id) {
          return mockService.fetchCoursesByCategory(params.category_id);
        }
        // 返回所有课程
        const mock = require('../mock/mockData');
        return Promise.resolve({ code: 200, message: 'success', data: mock.courses });
      }
      return request({
        url: '/contents',
        method: 'GET',
        data: params
      });
    },
    getDetail(id) {
      if (USE_MOCK) {
        return mockService.fetchCourseDetail(id);
      }
      return request({
        url: `/contents/${id}`,
        method: 'GET'
      });
    }
  },

  // 学习相关
  learning: {
    listProgress() {
      if (USE_MOCK) {
        return mockService.fetchProgress();
      }
      return request({
        url: '/learning',
        method: 'GET'
      });
    },
    getProgress(contentId) {
      if (USE_MOCK) {
        // Mock 模式下返回空进度
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: { 
            content_id: contentId, 
            video_position: 0, 
            duration_seconds: 0, 
            progress: 0, 
            status: 'not_started' 
          } 
        });
      }
      return request({
        url: `/learning/${contentId}`,
        method: 'GET'
      });
    },
    updateProgress(data) {
      if (USE_MOCK) {
        // Mock 模式下简单返回成功
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: { 
            content_id: data.content_id, 
            video_position: data.video_position || 0, 
            duration_seconds: 3600, 
            progress: Math.floor((data.video_position || 0) / 3600 * 100), 
            status: 'in_progress' 
          } 
        });
      }
      return request({
        url: '/learning',
        method: 'POST',
        data
      });
    },
    // 获取用户学习统计
    getUserStats() {
      if (USE_MOCK) {
        const user = wx.getStorageSync('user');
        if (user && user.id) {
          return mockService.fetchLearningStats(user.id);
        }
        return Promise.resolve({ code: 401, message: '未登录', data: null });
      }
      return request({
        url: '/learning/stats',
        method: 'GET'
      });
    },
    // 获取内容完成统计（管理员用）
    getContentStats(contentId) {
      if (USE_MOCK) {
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: { 
            content_id: contentId, 
            content_title: '示例内容', 
            completed_count: 10, 
            total_count: 20, 
            completion_rate: 50.0 
          } 
        });
      }
      return request({
        url: `/learning/content/${contentId}/stats`,
        method: 'GET'
      });
    }
  },

  // 考试相关
  exam: {
    listAvailable() {
      if (USE_MOCK) {
        // Mock 模式下返回示例考试列表
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: [
            {
              id: 1,
              title: '销售技巧考核',
              description: '测试你对销售技巧的掌握程度',
              total_score: 100,
              pass_score: 60,
              question_count: 10,
              time_limit_minutes: 60,
              status: 'published',
              attempt_status: 'not_started',
              last_score: null,
              last_passed: false,
              last_submitted_at: null
            },
            {
              id: 2,
              title: '产品知识测试',
              description: '检验产品相关知识的掌握情况',
              total_score: 100,
              pass_score: 70,
              question_count: 15,
              time_limit_minutes: 45,
              status: 'published',
              attempt_status: 'attempted',
              last_score: 65,
              last_passed: false,
              last_submitted_at: '2024-01-15T10:30:00Z'
            }
          ]
        });
      }
      return request({
        url: '/exams',
        method: 'GET'
      });
    },
    getDetail(id) {
      if (USE_MOCK) {
        // Mock 模式下返回示例考试详情
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: {
            id: id,
            title: '销售技巧考核',
            description: '测试你对销售技巧的掌握程度',
            time_limit_minutes: 60,
            pass_score: 60,
            total_score: 100,
            question_count: 3,
            questions: [
              {
                id: 1,
                type: 'single',
                stem: '以下哪个是销售过程中的关键步骤？',
                score: 10,
                options: [
                  { id: 1, label: 'A', content: '了解客户需求' },
                  { id: 2, label: 'B', content: '直接报价' },
                  { id: 3, label: 'C', content: '催促下单' },
                  { id: 4, label: 'D', content: '忽略客户反馈' }
                ]
              },
              {
                id: 2,
                type: 'single',
                stem: '客户提出异议时，应该如何处理？',
                score: 10,
                options: [
                  { id: 5, label: 'A', content: '立即反驳' },
                  { id: 6, label: 'B', content: '认真倾听并解答' },
                  { id: 7, label: 'C', content: '转移话题' },
                  { id: 8, label: 'D', content: '不予理睬' }
                ]
              },
              {
                id: 3,
                type: 'multiple',
                stem: '以下哪些是有效的销售技巧？（多选）',
                score: 20,
                options: [
                  { id: 9, label: 'A', content: '建立信任关系' },
                  { id: 10, label: 'B', content: '提供专业建议' },
                  { id: 11, label: 'C', content: '夸大产品功能' },
                  { id: 12, label: 'D', content: '跟进客户需求' }
                ]
              }
            ]
          }
        });
      }
      return request({
        url: `/exams/${id}`,
        method: 'GET'
      });
    },
    submit(id, data) {
      if (USE_MOCK) {
        // Mock 模式下模拟提交考试
        const totalScore = 100;
        const correctCount = Math.floor(Math.random() * 3) + 2; // 2-4 题正确
        const score = Math.floor((correctCount / 3) * totalScore);
        const pass = score >= 60;
        
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: {
            attempt_id: Date.now(),
            exam_id: id,
            score: score,
            total_score: totalScore,
            pass: pass,
            correct_count: correctCount,
            total_count: 3,
            duration_seconds: data.duration_seconds || 0,
            answers: [
              {
                question_id: 1,
                stem: '以下哪个是销售过程中的关键步骤？',
                type: 'single',
                score: 10,
                obtained_score: 10,
                is_correct: true,
                selected_option_ids: [1],
                correct_option_ids: [1]
              },
              {
                question_id: 2,
                stem: '客户提出异议时，应该如何处理？',
                type: 'single',
                score: 10,
                obtained_score: 10,
                is_correct: true,
                selected_option_ids: [6],
                correct_option_ids: [6]
              },
              {
                question_id: 3,
                stem: '以下哪些是有效的销售技巧？（多选）',
                type: 'multiple',
                score: 20,
                obtained_score: pass ? 20 : 0,
                is_correct: pass,
                selected_option_ids: [9, 10, 12],
                correct_option_ids: [9, 10, 12]
              }
            ]
          }
        });
      }
      return request({
        url: `/exams/${id}/submit`,
        method: 'POST',
        data
      });
    },
    listMyResults() {
      if (USE_MOCK) {
        // Mock 模式下返回考试结果列表
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: [
            {
              attempt_id: 1,
              exam_id: 2,
              exam_title: '产品知识测试',
              score: 65,
              total_score: 100,
              pass_score: 70,
              pass: false,
              submitted_at: '2024-01-15T10:30:00Z'
            }
          ]
        });
      }
      return request({
        url: '/exams/my/results',
        method: 'GET'
      });
    },
    managerOverview() {
      if (USE_MOCK) {
        // Mock 模式下返回店长考核概览
        const mock = require('../mock/mockData');
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: {
            exam_progress: [
              {
                exam_id: 1,
                title: '销售技巧考核',
                attempt_count: 5,
                pass_rate: 60.0,
                avg_score: 72.5
              },
              {
                exam_id: 2,
                title: '产品知识测试',
                attempt_count: 3,
                pass_rate: 33.3,
                avg_score: 65.0
              }
            ],
            employees: mock.managerProgressEmployees.map(emp => ({
              employee_id: emp.id,
              name: emp.name,
              work_no: `E${String(emp.id).padStart(3, '0')}`,
              latest_exam: {
                exam_id: 1,
                exam_title: '销售技巧考核',
                score: 75,
                pass: true,
                submitted_at: '2024-01-20T14:30:00Z'
              },
              learning_progress: {
                completed: emp.completed,
                total: emp.total,
                pending: emp.total - emp.completed,
                percent: Math.round((emp.completed / emp.total) * 100)
              }
            }))
          }
        });
      }
      return request({
        url: '/manager/exams/overview',
        method: 'GET'
      });
    }
  },

  // 轮播图
  banner: {
    listVisible() {
      if (USE_MOCK) {
        return mockService.fetchBanners();
      }
      return request({
        url: '/banners',
        method: 'GET'
      });
    }
  },
  
  // 管理员相关
  admin: {
    // 查询用户列表
    listUsers(params = {}) {
      if (USE_MOCK) {
        return mockService.fetchAllUsers();
      }
      return request({
        url: '/admin/users',
        method: 'GET',
        data: params
      });
    },
    // 查询单个用户
    getUser(id) {
      if (USE_MOCK) {
        const mock = require('../mock/mockData');
        const user = mock.users.find(u => u.id === id);
        if (!user) {
          return Promise.resolve({ code: 404, message: '用户不存在', data: null });
        }
        const { password: _ignored, ...safeUser } = user;
        return Promise.resolve({ code: 200, message: 'success', data: safeUser });
      }
      return request({
        url: `/admin/users/${id}`,
        method: 'GET'
      });
    },
    // 更新用户角色
    updateUserRole(id, role) {
      if (USE_MOCK) {
        return mockService.updateUserRole(id, role);
      }
      return request({
        url: `/admin/users/${id}/role`,
        method: 'PUT',
        data: { role }
      });
    },
    // 更新员工的店长绑定
    updateEmployeeManagers(id, managerWorkNos) {
      if (USE_MOCK) {
        // Mock 模式下，managerWorkNos 应该是 ID 数组
        return mockService.updateUserManagers(id, managerWorkNos);
      }
      return request({
        url: `/admin/users/${id}/managers`,
        method: 'PUT',
        data: {
          manager_ids: managerWorkNos
        }
      });
    },
    // 创建员工
    createEmployee(data) {
      if (USE_MOCK) {
        const mock = require('../mock/mockData');
        // 检查工号是否已存在
        const exists = mock.users.find(u => u.work_no === data.work_no);
        if (exists) {
          return Promise.resolve({ code: 400, message: '工号已存在', data: null });
        }
        const newEmployee = {
          id: Date.now(),
          work_no: data.work_no || '',
          name: data.name || '',
          phone: data.phone || '',
          role: 'employee',
          status: true,
          manager_ids: data.manager_ids || [],
          managers: [],
          points: 0
        };
        mock.users.push(newEmployee);
        const { password: _ignored, ...safeUser } = newEmployee;
        return Promise.resolve({ code: 200, message: 'success', data: safeUser });
      }
      return request({
        url: '/admin/employees',
        method: 'POST',
        data
      });
    },
    // 创建店长
    createManager(data) {
      if (USE_MOCK) {
        const mock = require('../mock/mockData');
        // 检查工号是否已存在
        const exists = mock.users.find(u => u.work_no === data.work_no);
        if (exists) {
          return Promise.resolve({ code: 400, message: '工号已存在', data: null });
        }
        const newManager = {
          id: Date.now(),
          work_no: data.work_no || '',
          name: data.name || '',
          phone: data.phone || '',
          role: 'manager',
          status: true,
          manager_ids: [],
          managers: [],
          points: 0
        };
        mock.users.push(newManager);
        mock.managers.push({
          id: newManager.id,
          work_no: newManager.work_no,
          name: newManager.name,
          phone: newManager.phone,
          role: 'manager',
          status: true
        });
        const { password: _ignored, ...safeUser } = newManager;
        return Promise.resolve({ code: 200, message: 'success', data: safeUser });
      }
      return request({
        url: '/admin/managers',
        method: 'POST',
        data
      });
    },
    // 轮播管理
    listBanners(params = {}) {
      if (USE_MOCK) {
        // Mock 模式下返回所有轮播图
        const mock = require('../mock/mockData');
        return Promise.resolve({ code: 200, message: 'success', data: mock.banners });
      }
      return request({
        url: '/admin/banners',
        method: 'GET',
        data: params
      });
    },
    createBanner(data) {
      if (USE_MOCK) {
        // Mock 模式下创建新轮播图
        const mock = require('../mock/mockData');
        const newBanner = {
          id: Date.now(),
          title: data.title || '',
          image_url: data.image_url || '',
          link_url: data.link_url || '',
          visible_roles: data.visible_roles || 'both',
          sort_order: data.sort_order || 1,
          status: data.status !== undefined ? data.status : true,
          start_at: data.start_at || null,
          end_at: data.end_at || null
        };
        mock.banners.push(newBanner);
        return Promise.resolve({ code: 200, message: 'success', data: newBanner });
      }
      return request({
        url: '/admin/banners',
        method: 'POST',
        data
      });
    },
    updateBanner(id, data) {
      if (USE_MOCK) {
        // Mock 模式下更新轮播图
        const mock = require('../mock/mockData');
        const banner = mock.banners.find(b => b.id === id);
        if (!banner) {
          return Promise.resolve({ code: 404, message: '轮播图不存在', data: null });
        }
        // 更新字段
        if (data.title !== undefined) banner.title = data.title;
        if (data.image_url !== undefined) banner.image_url = data.image_url;
        if (data.link_url !== undefined) banner.link_url = data.link_url;
        if (data.visible_roles !== undefined) banner.visible_roles = data.visible_roles;
        if (data.sort_order !== undefined) banner.sort_order = data.sort_order;
        if (data.status !== undefined) banner.status = data.status;
        if (data.start_at !== undefined) banner.start_at = data.start_at;
        if (data.end_at !== undefined) banner.end_at = data.end_at;
        return Promise.resolve({ code: 200, message: 'success', data: banner });
      }
      return request({
        url: `/admin/banners/${id}`,
        method: 'PUT',
        data
      });
    },
    // 学习内容管理
    listContents(params = {}) {
      if (USE_MOCK) {
        const mock = require('../mock/mockData');
        let contents = [...mock.courses];
        // 根据分类过滤
        if (params.category_id) {
          contents = contents.filter(c => c.category_id === params.category_id);
        }
        // 根据类型过滤
        if (params.type) {
          contents = contents.filter(c => c.type === params.type);
        }
        return Promise.resolve({ code: 200, message: 'success', data: contents });
      }
      return request({
        url: '/admin/contents',
        method: 'GET',
        data: params
      });
    },
    createContent(data) {
      if (USE_MOCK) {
        const mock = require('../mock/mockData');
        const newContent = {
          id: Date.now(),
          category_id: data.category_id || 1,
          category_name: mock.courseCategories.find(c => c.id === data.category_id)?.name || '',
          title: data.title || '',
          type: data.type || 'doc',
          file_path: data.file_path || '',
          cover_url: data.cover_url || '',
          duration_seconds: data.duration_seconds || 0,
          summary: data.summary || '',
          status: data.status || 'published',
          visible_roles: data.visible_roles || 'both',
          publish_at: new Date().toISOString()
        };
        mock.courses.push(newContent);
        return Promise.resolve({ code: 200, message: 'success', data: newContent });
      }
      return request({
        url: '/admin/contents',
        method: 'POST',
        data
      });
    },
    updateContent(id, data) {
      if (USE_MOCK) {
        const mock = require('../mock/mockData');
        const content = mock.courses.find(c => c.id === id);
        if (!content) {
          return Promise.resolve({ code: 404, message: '内容不存在', data: null });
        }
        // 更新字段
        if (data.title !== undefined) content.title = data.title;
        if (data.type !== undefined) content.type = data.type;
        if (data.category_id !== undefined) {
          content.category_id = data.category_id;
          content.category_name = mock.courseCategories.find(c => c.id === data.category_id)?.name || '';
        }
        if (data.file_path !== undefined) content.file_path = data.file_path;
        if (data.cover_url !== undefined) content.cover_url = data.cover_url;
        if (data.duration_seconds !== undefined) content.duration_seconds = data.duration_seconds;
        if (data.summary !== undefined) content.summary = data.summary;
        if (data.status !== undefined) content.status = data.status;
        if (data.visible_roles !== undefined) content.visible_roles = data.visible_roles;
        return Promise.resolve({ code: 200, message: 'success', data: content });
      }
      return request({
        url: `/admin/contents/${id}`,
        method: 'PUT',
        data
      });
    },
    // 考试管理
    listExams(params = {}) {
      if (USE_MOCK) {
        // Mock 模式下返回空列表或示例数据
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: [] 
        });
      }
      return request({
        url: '/admin/exams',
        method: 'GET',
        data: params
      });
    },
    createExam(data) {
      if (USE_MOCK) {
        // Mock 模式下创建考试
        const newExam = {
          id: Date.now(),
          title: data.title || '',
          description: data.description || '',
          status: data.status || 'published',
          target_role: data.target_role || 'all',
          time_limit_minutes: data.time_limit_minutes || 0,
          pass_score: data.pass_score || 60,
          questions: data.questions || []
        };
        return Promise.resolve({ code: 200, message: 'success', data: newExam });
      }
      return request({
        url: '/admin/exams',
        method: 'POST',
        data
      });
    },
    updateExam(id, data) {
      if (USE_MOCK) {
        // Mock 模式下简单返回成功
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: { id, ...data } 
        });
      }
      return request({
        url: `/admin/exams/${id}`,
        method: 'PUT',
        data
      });
    },
    getExamDetail(id) {
      if (USE_MOCK) {
        // Mock 模式下返回示例考试
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: {
            id,
            title: '示例考试',
            description: '这是一个示例考试',
            time_limit_minutes: 60,
            pass_score: 60,
            total_score: 100,
            question_count: 10,
            questions: []
          }
        });
      }
      return request({
        url: `/admin/exams/${id}`,
        method: 'GET'
      });
    },
    // 积分管理
    listPoints(params = {}) {
      if (USE_MOCK) {
        // Mock 模式下返回用户积分列表
        const mock = require('../mock/mockData');
        const users = mock.users.map(user => ({
          id: user.id,
          work_no: user.work_no,
          name: user.name,
          phone: user.phone,
          role: user.role,
          status: user.status,
          points: user.points || 0
        }));
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: {
            items: users,
            pagination: {
              page: 1,
              page_size: 20,
              total: users.length
            }
          }
        });
      }
      return request({
        url: '/admin/points',
        method: 'GET',
        data: params
      });
    },
    getUserPoints(userId, params = {}) {
      if (USE_MOCK) {
        const mock = require('../mock/mockData');
        const user = mock.users.find(u => u.id === userId);
        if (!user) {
          return Promise.resolve({ code: 404, message: '用户不存在', data: null });
        }
        // 返回用户积分详情
        return Promise.resolve({ 
          code: 200, 
          message: 'success', 
          data: {
            user: {
              id: user.id,
              work_no: user.work_no,
              name: user.name,
              phone: user.phone,
              role: user.role,
              status: user.status
            },
            total_points: user.points || 0,
            transactions: [],
            pagination: {
              page: 1,
              page_size: 20,
              total: 0
            }
          }
        });
      }
      return request({
        url: `/admin/users/${userId}/points`,
        method: 'GET',
        data: params
      });
    }
  }
};

