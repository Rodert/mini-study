package model

import "time"

// TableName 指定表名
func (ExamPaper) TableName() string {
	return "exam_papers"
}

// ExamPaper represents an exam paper that contains multiple questions.
type ExamPaper struct {
	Base
	Title            string         `gorm:"size:200;not null;comment:试卷标题" json:"title"`
	Description      string         `gorm:"type:text;comment:试卷描述" json:"description"`
	Status           string         `gorm:"size:20;default:'draft';comment:状态(draft草稿/published已发布)" json:"status"`
	TargetRole       Role           `gorm:"size:16;default:'employee';comment:目标角色(employee员工/manager店长)" json:"target_role"`
	TimeLimitMinutes int            `gorm:"default:0;comment:时间限制(分钟)" json:"time_limit_minutes"`
	PassScore        int            `gorm:"default:0;comment:及格分数" json:"pass_score"`
	TotalScore       int            `gorm:"default:0;comment:总分" json:"total_score"`
	CreatorID        uint           `gorm:"comment:创建者ID" json:"creator_id"`
	Questions        []ExamQuestion `json:"questions" gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	QuestionCount int `gorm:"-" json:"question_count"`
}

// TableName 指定表名
func (ExamQuestion) TableName() string {
	return "exam_questions"
}

// ExamQuestion represents a single exam question.
type ExamQuestion struct {
	Base
	ExamID   uint         `gorm:"not null;index;comment:试卷ID" json:"exam_id"`
	Type     string       `gorm:"size:16;not null;comment:题型(single单选/multiple多选)" json:"type"` // single / multiple
	Stem     string       `gorm:"type:text;comment:题干" json:"stem"`
	Score    int          `gorm:"default:1;comment:分值" json:"score"`
	Analysis string       `gorm:"type:text;comment:解析" json:"analysis"`
	Options  []ExamOption `json:"options" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 指定表名
func (ExamOption) TableName() string {
	return "exam_options"
}

// ExamOption represents an option belonging to a question.
type ExamOption struct {
	Base
	QuestionID uint   `gorm:"not null;index;comment:题目ID" json:"question_id"`
	Label      string `gorm:"size:8;comment:选项标签(A/B/C/D)" json:"label"`
	Content    string `gorm:"type:text;comment:选项内容" json:"content"`
	IsCorrect  bool   `gorm:"default:false;comment:是否正确答案" json:"is_correct"`
	SortOrder  int    `gorm:"default:0;comment:排序顺序" json:"sort_order"`
}

// TableName 指定表名
func (ExamAttempt) TableName() string {
	return "exam_attempts"
}

// ExamAttempt records a user's submission for an exam.
type ExamAttempt struct {
	Base
	ExamID          uint       `json:"exam_id" gorm:"uniqueIndex:idx_user_exam;comment:试卷ID"`
	UserID          uint       `json:"user_id" gorm:"uniqueIndex:idx_user_exam;comment:用户ID"`
	Status          string     `gorm:"size:16;default:'submitted';comment:状态(submitted已提交/grading评分中/completed已完成)" json:"status"`
	Score           int        `gorm:"comment:得分" json:"score"`
	CorrectCount    int        `gorm:"comment:正确题数" json:"correct_count"`
	TotalCount      int        `gorm:"comment:总题数" json:"total_count"`
	Pass            bool       `gorm:"comment:是否通过" json:"pass"`
	DurationSeconds int64      `gorm:"comment:答题时长(秒)" json:"duration_seconds"`
	AnswerSnapshot  []byte     `gorm:"type:json;comment:答案快照(JSON格式)" json:"answer_snapshot"`
	SubmittedAt     *time.Time `gorm:"comment:提交时间" json:"submitted_at"`

	Exam ExamPaper `json:"exam" gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
