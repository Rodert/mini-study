// API 服务配置
const API_BASE_URL = 'http://localhost:8080/api/v1'; // 根据实际环境修改
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

    wx.request({
      url: `${API_BASE_URL}${url}`,
      method,
      data: cleanData,
      header,
      success: (res) => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data);
        } else if (res.statusCode === 401) {
          // 未授权，清除 token 并跳转登录
          handleUnauthorized();
          reject({
            code: 401,
            message: '登录已过期，请重新登录',
            data: res.data
          });
        } else {
          reject({
            code: res.statusCode,
            message: res.data?.message || '请求失败',
            data: res.data
          });
        }
      },
      fail: (err) => {
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
      const res = await request({
        url: '/users/login',
        method: 'POST',
        data: {
          work_no: data.work_no,
          password: data.password
        }
      });
      
      // 保存 token
      if (res.code === 200 && res.data) {
        if (res.data.access_token) {
          wx.setStorageSync('token', res.data.access_token);
        }
        if (res.data.refresh_token) {
          wx.setStorageSync('refresh_token', res.data.refresh_token);
        }
      }
      
      return res;
    },
    // 获取当前用户信息
    getCurrentUser() {
      return request({
        url: '/users/me',
        method: 'GET'
      });
    },
    // 获取店长列表
    getManagers() {
      return request({
        url: '/users/managers',
        method: 'GET'
      });
    },
    // 员工注册
    register(data) {
      return request({
        url: '/users/register',
        method: 'POST',
        data
      });
    },
    // 更新个人资料
    updateProfile(data) {
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
      return request({
        url: '/contents/categories',
        method: 'GET'
      });
    },
    listPublished(params = {}) {
      return request({
        url: '/contents',
        method: 'GET',
        data: params
      });
    },
    getDetail(id) {
      return request({
        url: `/contents/${id}`,
        method: 'GET'
      });
    }
  },

  // 学习相关
  learning: {
    listProgress() {
      return request({
        url: '/learning',
        method: 'GET'
      });
    },
    getProgress(contentId) {
      return request({
        url: `/learning/${contentId}`,
        method: 'GET'
      });
    },
    updateProgress(data) {
      return request({
        url: '/learning',
        method: 'POST',
        data
      });
    },
    // 获取用户学习统计
    getUserStats() {
      return request({
        url: '/learning/stats',
        method: 'GET'
      });
    },
    // 获取内容完成统计（管理员用）
    getContentStats(contentId) {
      return request({
        url: `/learning/content/${contentId}/stats`,
        method: 'GET'
      });
    }
  },

  // 考试相关
  exam: {
    listAvailable() {
      return request({
        url: '/exams',
        method: 'GET'
      });
    },
    getDetail(id) {
      return request({
        url: `/exams/${id}`,
        method: 'GET'
      });
    },
    submit(id, data) {
      return request({
        url: `/exams/${id}/submit`,
        method: 'POST',
        data
      });
    },
    listMyResults() {
      return request({
        url: '/exams/my/results',
        method: 'GET'
      });
    },
    managerOverview() {
      return request({
        url: '/manager/exams/overview',
        method: 'GET'
      });
    }
  },

  // 轮播图
  banner: {
    listVisible() {
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
      return request({
        url: '/admin/users',
        method: 'GET',
        data: params
      });
    },
    // 查询单个用户
    getUser(id) {
      return request({
        url: `/admin/users/${id}`,
        method: 'GET'
      });
    },
    // 更新用户角色
    updateUserRole(id, role) {
      return request({
        url: `/admin/users/${id}/role`,
        method: 'PUT',
        data: { role }
      });
    },
    // 更新员工的店长绑定
    updateEmployeeManagers(id, managerWorkNos) {
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
      return request({
        url: '/admin/employees',
        method: 'POST',
        data
      });
    },
    // 创建店长
    createManager(data) {
      return request({
        url: '/admin/managers',
        method: 'POST',
        data
      });
    },
    // 轮播管理
    listBanners(params = {}) {
      return request({
        url: '/admin/banners',
        method: 'GET',
        data: params
      });
    },
    createBanner(data) {
      return request({
        url: '/admin/banners',
        method: 'POST',
        data
      });
    },
    updateBanner(id, data) {
      return request({
        url: `/admin/banners/${id}`,
        method: 'PUT',
        data
      });
    },
    // 学习内容管理
    listContents(params = {}) {
      return request({
        url: '/admin/contents',
        method: 'GET',
        data: params
      });
    },
    createContent(data) {
      return request({
        url: '/admin/contents',
        method: 'POST',
        data
      });
    },
    updateContent(id, data) {
      return request({
        url: `/admin/contents/${id}`,
        method: 'PUT',
        data
      });
    },
    // 考试管理
    listExams(params = {}) {
      return request({
        url: '/admin/exams',
        method: 'GET',
        data: params
      });
    },
    createExam(data) {
      return request({
        url: '/admin/exams',
        method: 'POST',
        data
      });
    },
    updateExam(id, data) {
      return request({
        url: `/admin/exams/${id}`,
        method: 'PUT',
        data
      });
    },
    getExamDetail(id) {
      return request({
        url: `/admin/exams/${id}`,
        method: 'GET'
      });
    },
    updateExam(id, data) {
      return request({
        url: `/admin/exams/${id}`,
        method: 'PUT',
        data
      });
    },
    // 积分管理
    listPoints(params = {}) {
      return request({
        url: '/admin/points',
        method: 'GET',
        data: params
      });
    },
    getUserPoints(userId, params = {}) {
      return request({
        url: `/admin/users/${userId}/points`,
        method: 'GET',
        data: params
      });
    }
  }
};

