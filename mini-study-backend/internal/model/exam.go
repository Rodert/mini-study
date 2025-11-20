package model

import "time"

// ExamPaper represents an exam paper that contains multiple questions.
type ExamPaper struct {
	Base
	Title            string         `gorm:"size:200;not null" json:"title"`
	Description      string         `gorm:"type:text" json:"description"`
	Status           string         `gorm:"size:20;default:'draft'" json:"status"`
	TargetRole       Role           `gorm:"size:16;default:'employee'" json:"target_role"`
	TimeLimitMinutes int            `gorm:"default:0" json:"time_limit_minutes"`
	PassScore        int            `gorm:"default:0" json:"pass_score"`
	TotalScore       int            `gorm:"default:0" json:"total_score"`
	CreatorID        uint           `json:"creator_id"`
	Questions        []ExamQuestion `json:"questions" gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	QuestionCount int `gorm:"-" json:"question_count"`
}

// ExamQuestion represents a single exam question.
type ExamQuestion struct {
	Base
	ExamID   uint         `gorm:"not null;index" json:"exam_id"`
	Type     string       `gorm:"size:16;not null" json:"type"` // single / multiple
	Stem     string       `gorm:"type:text" json:"stem"`
	Score    int          `gorm:"default:1" json:"score"`
	Analysis string       `gorm:"type:text" json:"analysis"`
	Options  []ExamOption `json:"options" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// ExamOption represents an option belonging to a question.
type ExamOption struct {
	Base
	QuestionID uint   `gorm:"not null;index" json:"question_id"`
	Label      string `gorm:"size:8" json:"label"`
	Content    string `gorm:"type:text" json:"content"`
	IsCorrect  bool   `gorm:"default:false" json:"is_correct"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
}

// ExamAttempt records a user's submission for an exam.
type ExamAttempt struct {
	Base
	ExamID          uint           `json:"exam_id" gorm:"uniqueIndex:idx_user_exam"`
	UserID          uint           `json:"user_id" gorm:"uniqueIndex:idx_user_exam"`
	Status          string         `gorm:"size:16;default:'submitted'" json:"status"`
	Score           int            `json:"score"`
	CorrectCount    int            `json:"correct_count"`
	TotalCount      int            `json:"total_count"`
	Pass            bool           `json:"pass"`
	DurationSeconds int64          `json:"duration_seconds"`
	AnswerSnapshot  []byte      `gorm:"type:json" json:"answer_snapshot"`
	SubmittedAt     *time.Time     `json:"submitted_at"`

	Exam ExamPaper `json:"exam" gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
