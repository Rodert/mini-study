const banners = [
  {
    id: 1,
    type: "web",
    cover: "https://inews.gtimg.com/newsapp_bt/0/15956949227/1000",
    url: "https://news.qq.com/rain/a/20251118A014A400",
    title: "企业新闻 · 最新资讯"
  },
  {
    id: 2,
    type: "web",
    cover: "https://example.com/banner2.jpg",
    url: "https://www.qq.com",
    title: "新人训 · 快速入口"
  },
  {
    id: 3,
    type: "web",
    cover: "https://example.com/banner3.jpg",
    url: "https://www.tencent.com",
    title: "视频教学 · 销售技巧"
  }
];

const courseCategories = [
  { id: 1, name: "专题特训", role: "employee" },
  { id: 2, name: "金牌课程", role: "employee" },
  { id: 3, name: "新人养成", role: "employee" },
  { id: 4, name: "短视频教学", role: "employee" },
  { id: 5, name: "营销推广", role: "employee" },
  { id: 6, name: "学习执行", role: "employee" },
  { id: 7, name: "销售技巧", role: "employee" },
  { id: 101, name: "店长特训", role: "manager" },
  { id: 102, name: "新员工培训", role: "manager" },
  { id: 103, name: "老员工进阶", role: "manager" },
  { id: 104, name: "管理知识与自我提升", role: "manager" }
];

const courses = [
  {
    id: 10,
    category_id: 1,
    title: "销售基础特训 01",
    cover: "https://example.com/course1.jpg",
    type: "video",
    duration: "12min",
    url: "https://example.com/video1.mp4",
    summary: "覆盖销售模型、客户洞察等核心内容"
  },
  {
    id: 20,
    category_id: 101,
    title: "店长管理基础 01",
    cover: "https://example.com/course2.jpg",
    type: "article",
    duration: "15min",
    content: "<p>文章内容……</p>",
    summary: "聚焦门店管理、团队带教"
  }
];

const managers = [
  {
    id: 2,
    name: "李婷",
    store: "上海一店",
    mobile: "13911110000"
  },
  {
    id: 3,
    name: "张三",
    store: "上海二店",
    mobile: "13922220000"
  },
  {
    id: 4,
    name: "王五",
    store: "上海三店",
    mobile: "13933330000"
  }
];

const users = [
  {
    id: 1,
    username: "1001",
    name: "王明",
    mobile: "13900000000",
    role: "employee",
    store: "上海一店",
    managerIds: [2],
    password: "123456"
  },
  {
    id: 2,
    username: "2001",
    name: "李婷",
    mobile: "13911110000",
    role: "manager",
    store: "上海一店",
    password: "123456"
  },
  {
    id: 3,
    username: "admin",
    name: "管理员",
    mobile: "13900000001",
    role: "admin",
    store: "总部",
    password: "123456"
  }
];

const managerProgress = [
  {
    courseId: 1010,
    courseName: "店长管理基础 01",
    completion: 0.78
  },
  {
    courseId: 1011,
    courseName: "店长管理基础 02",
    completion: 0.42
  }
];

const managerProgressEmployees = [
  {
    id: 11,
    name: "王晓",
    store: "门店 01",
    completed: 4,
    total: 5,
    pending: "老员工进阶"
  },
  {
    id: 12,
    name: "李静",
    store: "门店 03",
    completed: 2,
    total: 5,
    pending: "管理知识与自我提升"
  },
  {
    id: 13,
    name: "赵晨",
    store: "门店 02",
    completed: 3,
    total: 5,
    pending: "店长特训"
  }
];

const learningStats = {
  1: {
    totalHours: 12,
    completedCourses: 6,
    totalCourses: 8,
    streakDays: 5
  },
  2: {
    totalHours: 18,
    completedCourses: 9,
    totalCourses: 10,
    streakDays: 8
  }
};

module.exports = {
  banners,
  courseCategories,
  courses,
  users,
  managers,
  managerProgress,
  managerProgressEmployees,
  learningStats
};

