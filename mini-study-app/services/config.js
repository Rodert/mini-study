// 应用配置
// 是否使用 Mock 数据（开发环境可以设置为 true，生产环境应设置为 false）
const USE_MOCK = true; // 设置为 true 使用 mock 数据， false 使用真实 API

// API 基础地址配置
const API_BASE_URL = 'http://localhost:8080/api/v1'; // 根据实际环境修改

module.exports = {
  USE_MOCK,
  API_BASE_URL
};

