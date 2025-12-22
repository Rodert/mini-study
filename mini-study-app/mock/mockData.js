const banners = [
  {
    id: 1,
    title: "企业新闻 · 最新资讯",
    image_url: "https://inews.gtimg.com/newsapp_bt/0/15956949227/1000",
    link_url: "https://news.qq.com/rain/a/20251118A014A400",
    visible_roles: "both",
    sort_order: 1,
    status: true,
    start_at: null,
    end_at: null
  },
  {
    id: 2,
    title: "新人训 · 快速入口",
    image_url: "https://example.com/banner2.jpg",
    link_url: "https://www.qq.com",
    visible_roles: "employee",
    sort_order: 2,
    status: true,
    start_at: null,
    end_at: null
  },
  {
    id: 3,
    title: "视频教学 · 销售技巧",
    image_url: "https://example.com/banner3.jpg",
    link_url: "https://www.tencent.com",
    visible_roles: "both",
    sort_order: 3,
    status: true,
    start_at: null,
    end_at: null
  }
];

const courseCategories = [
  { id: 1, name: "专题特训", role_scope: "employee", sort_order: 1, count: 5 },
  { id: 2, name: "金牌课程", role_scope: "employee", sort_order: 2, count: 8 },
  { id: 3, name: "新人养成", role_scope: "employee", sort_order: 3, count: 3 },
  { id: 4, name: "短视频教学", role_scope: "employee", sort_order: 4, count: 12 },
  { id: 5, name: "营销推广", role_scope: "employee", sort_order: 5, count: 6 },
  { id: 6, name: "学习执行", role_scope: "employee", sort_order: 6, count: 4 },
  { id: 7, name: "销售技巧", role_scope: "employee", sort_order: 7, count: 10 },
  { id: 101, name: "店长特训", role_scope: "manager", sort_order: 1, count: 3 },
  { id: 102, name: "新员工培训", role_scope: "manager", sort_order: 2, count: 5 },
  { id: 103, name: "老员工进阶", role_scope: "manager", sort_order: 3, count: 4 },
  { id: 104, name: "管理知识与自我提升", role_scope: "manager", sort_order: 4, count: 6 }
];

const courses = [
  {
    id: 10,
    category_id: 1,
    category_name: "专题特训",
    title: "销售基础特训 01",
    cover_url: "https://example.com/course1.jpg",
    type: "video",
    file_path: "/uploads/video1.mp4",
    duration_seconds: 720,
    summary: "覆盖销售模型、客户洞察等核心内容",
    status: "published",
    visible_roles: "employee",
    publish_at: "2024-01-01T00:00:00Z"
  },
  {
    id: 20,
    category_id: 101,
    category_name: "店长特训",
    title: "店长管理基础 01",
    cover_url: "https://example.com/course2.jpg",
    type: "doc",
    file_path: "/uploads/doc1.pdf",
    duration_seconds: 0,
    summary: "聚焦运营管理、团队带教",
    status: "published",
    visible_roles: "manager",
    publish_at: "2024-01-01T00:00:00Z"
  }
];

const managers = [
  {
    id: 2,
    work_no: "manager001",
    name: "李婷",
    phone: "13911110000",
    role: "manager",
    status: true
  },
  {
    id: 3,
    work_no: "manager002",
    name: "张三",
    phone: "13922220000",
    role: "manager",
    status: true
  },
  {
    id: 4,
    work_no: "manager003",
    name: "王五",
    phone: "13933330000",
    role: "manager",
    status: true
  }
];

const users = [
  {
    id: 1,
    work_no: "employee001",
    name: "王明",
    phone: "13900000000",
    role: "employee",
    status: true,
    manager_ids: [2],
    managers: [
      {
        id: 2,
        work_no: "manager001",
        name: "李婷",
        phone: "13911110000"
      }
    ],
    points: 100,
    password: "123456"
  },
  {
    id: 2,
    work_no: "manager001",
    name: "李婷",
    phone: "13911110000",
    role: "manager",
    status: true,
    manager_ids: [],
    managers: [],
    points: 200,
    password: "123456"
  },
  {
    id: 3,
    work_no: "admin",
    name: "管理员",
    phone: "13900000001",
    role: "admin",
    status: true,
    manager_ids: [],
    managers: [],
    points: 0,
    password: "admin123456"
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
    store: "部门 01",
    completed: 4,
    total: 5,
    pending: "老员工进阶"
  },
  {
    id: 12,
    name: "李静",
    store: "部门 03",
    completed: 2,
    total: 5,
    pending: "管理知识与自我提升"
  },
  {
    id: 13,
    name: "赵晨",
    store: "部门 02",
    completed: 3,
    total: 5,
    pending: "店长特训"
  }
];

const learningStats = {
  1: {
    user_id: 1,
    completed_count: 6,
    total_count: 8,
    total_contents: 20,
    completion_rate: 30.0
  },
  2: {
    user_id: 2,
    completed_count: 9,
    total_count: 10,
    total_contents: 15,
    completion_rate: 60.0
  }
};

const notices = [
  {
    id: 1,
    title: "系统维护公告",
    content: "本周日 22:00-23:00 将进行系统维护，期间部分功能可能短暂不可用，请合理安排学习与考试时间。",
    image_url: "",
    status: true,
    start_at: null,
    end_at: null,
    created_at: "2024-01-01T10:00:00Z"
  },
  {
    id: 2,
    title: "新品培训上线",
    content: "新品培训课程已上线，请各位员工在本周内完成相关学习内容。",
    image_url: "",
    status: true,
    start_at: null,
    end_at: null,
    created_at: "2024-02-01T09:00:00Z"
  }
];

module.exports = {
  banners,
  courseCategories,
  courses,
  users,
  managers,
  managerProgress,
  managerProgressEmployees,
  learningStats,
  notices
};

