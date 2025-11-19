// API 服务配置
const API_BASE_URL = 'http://localhost:8080/api/v1'; // 根据实际环境修改

// 获取 token
function getToken() {
  return wx.getStorageSync('token') || '';
}

// 通用请求方法
function request(options) {
  return new Promise((resolve, reject) => {
    const { url, method = 'GET', data = {}, header = {} } = options;
    
    // 添加 token
    const token = getToken();
    if (token) {
      header['Authorization'] = `Bearer ${token}`;
    }
    header['Content-Type'] = 'application/json';

    wx.request({
      url: `${API_BASE_URL}${url}`,
      method,
      data,
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
  // 用户相关
  user: {
    // 登录
    async login(data) {
      const res = await request({
        url: '/users/login',
        method: 'POST',
        data: {
          work_no: data.work_no || data.mobile, // 兼容 mobile 字段
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
    }
  },
  
  // 管理员相关
  admin: {
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
    }
  }
};

