const mock = require("../mock/mockData");

const delay = (ms = 300) =>
  new Promise((resolve) => setTimeout(resolve, ms));

function respond(data) {
  return delay().then(() => ({ code: 200, message: "success", data }));
}

function errorResponse(message, code = 400) {
  return delay().then(() => ({ code, message, data: null }));
}

module.exports = {
  login({ work_no, password }) {
    const user = mock.users.find(
      (item) =>
        item.work_no === work_no &&
        item.password === password
    );
    if (!user) {
      return errorResponse("工号或密码不正确", 401);
    }
    // 返回 token 对象，而不是用户对象
    return respond({
      access_token: `mock_access_token_${user.id}_${Date.now()}`,
      refresh_token: `mock_refresh_token_${user.id}_${Date.now()}`
    });
  },
  register(payload) {
    const exists = mock.users.find(
      (item) =>
        item.phone === payload.phone || item.work_no === payload.work_no
    );

    if (exists) {
      return errorResponse("手机号或工号已存在", 400);
    }

    const newUser = {
      id: Date.now(),
      work_no: payload.work_no,
      name: payload.name,
      phone: payload.phone,
      role: "employee",
      status: true,
      manager_ids: [],
      managers: [],
      points: 0,
      password: payload.password
    };

    mock.users.push(newUser);

    const { password: _ignored, ...safeUser } = newUser;
    return respond({
      id: safeUser.id,
      work_no: safeUser.work_no,
      phone: safeUser.phone,
      name: safeUser.name,
      role: safeUser.role,
      status: safeUser.status
    });
  },
  fetchBanners() {
    return respond(mock.banners);
  },
  fetchCourseCategories(role) {
    // 根据角色过滤分类
    const roleScope = role === "manager" ? "manager" : role === "admin" ? "both" : "employee";
    const list = mock.courseCategories.filter((item) => 
      item.role_scope === roleScope || item.role_scope === "both"
    );
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
      return errorResponse("用户不存在", 404);
    }
    if (updates.name) user.name = updates.name;
    if (updates.phone) user.phone = updates.phone;
    return respond({
      id: user.id,
      work_no: user.work_no,
      phone: user.phone,
      name: user.name,
      role: user.role,
      status: user.status
    });
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
      return errorResponse("用户不存在", 404);
    }
    user.role = role;
    const { password: _ignored, ...safeUser } = user;
    return respond(safeUser);
  },
  updateUserManagers(userId, managerIds) {
    const user = mock.users.find(item => item.id === userId);
    if (!user) {
      return errorResponse("用户不存在", 404);
    }
    user.manager_ids = managerIds;
    // 更新 managers 数组
    user.managers = mock.managers.filter(m => managerIds.includes(m.id));
    const { password: _ignored, ...safeUser } = user;
    return respond(safeUser);
  }
};

