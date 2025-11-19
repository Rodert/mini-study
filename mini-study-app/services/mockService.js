const mock = require("../mock/mockData");

const delay = (ms = 300) =>
  new Promise((resolve) => setTimeout(resolve, ms));

function respond(data) {
  return delay().then(() => ({ data, success: true }));
}

module.exports = {
  login({ mobile, password }) {
    const user = mock.users.find(
      (item) =>
        item.mobile === mobile &&
        item.password === password
    );
    if (!user) {
      return delay().then(() => ({
        success: false,
        message: "手机号或密码不正确"
      }));
    }
    const { password: _ignored, ...safeUser } = user;
    return respond(safeUser);
  },
  register(payload) {
    const exists = mock.users.find(
      (item) =>
        item.mobile === payload.mobile || item.username === payload.employeeId
    );

    if (exists) {
      return delay().then(() => ({
        success: false,
        message: "手机号或工号已存在"
      }));
    }

    const newUser = {
      id: Date.now(),
      username: payload.employeeId,
      name: payload.name,
      mobile: payload.mobile,
      role: payload.role,
      store: payload.store || "未分配门店",
      password: payload.password
    };

    mock.users.push(newUser);

    const { password: _ignored, ...safeUser } = newUser;
    return respond(safeUser);
  },
  fetchBanners() {
    return respond(mock.banners);
  },
  fetchCourseCategories(role) {
    const list = mock.courseCategories.filter((item) => item.role === role);
    return respond(list);
  },
  fetchCoursesByCategory(categoryId) {
    const list = mock.courses.filter(
      (course) => course.category_id === categoryId
    );
    return respond(list);
  },
  fetchCourseDetail(courseId) {
    const course = mock.courses.find((item) => item.id === courseId);
    return respond(course);
  },
  fetchProgress() {
    return respond(mock.managerProgress);
  },
  fetchProgressEmployees() {
    return respond(mock.managerProgressEmployees);
  },
  fetchLearningStats(userId) {
    const stats = mock.learningStats[userId];
    return respond(stats || null);
  },
  fetchManagers() {
    return respond(mock.managers);
  },
  updateUserProfile(userId, updates) {
    const user = mock.users.find(item => item.id === userId);
    if (!user) {
      return delay().then(() => ({
        success: false,
        message: "用户不存在"
      }));
    }
    Object.assign(user, updates);
    const { password: _ignored, ...safeUser } = user;
    return respond(safeUser);
  },
  fetchAllUsers() {
    const users = mock.users.map(user => {
      const { password: _ignored, ...safeUser } = user;
      return safeUser;
    });
    return respond(users);
  },
  updateUserRole(userId, role) {
    const user = mock.users.find(item => item.id === userId);
    if (!user) {
      return delay().then(() => ({
        success: false,
        message: "用户不存在"
      }));
    }
    user.role = role;
    const { password: _ignored, ...safeUser } = user;
    return respond(safeUser);
  },
  updateUserManagers(userId, managerIds) {
    const user = mock.users.find(item => item.id === userId);
    if (!user) {
      return delay().then(() => ({
        success: false,
        message: "用户不存在"
      }));
    }
    user.managerIds = managerIds;
    const { password: _ignored, ...safeUser } = user;
    return respond(safeUser);
  }
};

