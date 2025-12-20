package dto

import "time"

// AdminExamQuestionOption defines option payload for admin upsert.
type AdminExamQuestionOption struct {
	Label     string `json:"label" binding:"required"`
	Content   string `json:"content" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
	SortOrder int    `json:"sort_order"`
}

// AdminExamQuestionUpsert defines question payload for admin upsert.
type AdminExamQuestionUpsert struct {
	Type     string                    `json:"type" binding:"required,oneof=single multiple"`
	Stem     string                    `json:"stem" binding:"required"`
	Score    int                       `json:"score" binding:"required,min=1"`
	Analysis string                    `json:"analysis"`
	Options  []AdminExamQuestionOption `json:"options" binding:"required,min=2,dive"`
}

// AdminExamUpsertRequest represents admin create/update exam request.
type AdminExamUpsertRequest struct {
	Title            string                    `json:"title" binding:"required"`
	Description      string                    `json:"description"`
	Status           string                    `json:"status" binding:"omitempty,oneof=draft published archived"`
	TargetRole       string                    `json:"target_role" binding:"omitempty,oneof=employee manager all"`
	TimeLimitMinutes int                       `json:"time_limit_minutes" binding:"omitempty,min=0"`
	PassScore        int                       `json:"pass_score" binding:"required,min=0"`
	Questions        []AdminExamQuestionUpsert `json:"questions" binding:"required,min=1,dive"`
}

// ExamListItem represents summary for available exams.
type ExamListItem struct {
	ID               uint       `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	TotalScore       int        `json:"total_score"`
	PassScore        int        `json:"pass_score"`
	TimeLimitMinutes int        `json:"time_limit_minutes"`
	QuestionCount    int        `json:"question_count"`
	Status           string     `json:"status"`
	AttemptStatus    string     `json:"attempt_status"`
	LastScore        int        `json:"last_score"`
	LastPassed       bool       `json:"last_passed"`
	LastSubmittedAt  *time.Time `json:"last_submitted_at"`
}

// ExamDetailQuestionOption is returned to exam detail API.
type ExamDetailQuestionOption struct {
	ID        uint   `json:"id"`
	Label     string `json:"label"`
	Content   string `json:"content"`
	IsCorrect bool   `json:"is_correct,omitempty"` // Only for admin
}

// ExamDetailQuestion describes question for display.
type ExamDetailQuestion struct {
	ID        uint                       `json:"id"`
	Type      string                     `json:"type"`
	Stem      string                     `json:"stem"`
	Score     int                        `json:"score"`
	Analysis  string                     `json:"analysis,omitempty"` // Only for admin
	Options   []ExamDetailQuestionOption `json:"options"`
}

// ExamDetailResponse describes exam for answering.
type ExamDetailResponse struct {
	ID               uint                 `json:"id"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	TimeLimitMinutes int                  `json:"time_limit_minutes"`
	PassScore        int                  `json:"pass_score"`
	TotalScore       int                  `json:"total_score"`
	QuestionCount    int                  `json:"question_count"`
	Questions        []ExamDetailQuestion `json:"questions"`
}

// ExamSubmitAnswer describes a user's answer payload.
type ExamSubmitAnswer struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	OptionIDs  []uint `json:"option_ids" binding:"required,min=1,dive,required"`
}

// ExamSubmitRequest holds submission data.
type ExamSubmitRequest struct {
	Answers         []ExamSubmitAnswer `json:"answers" binding:"required,min=1,dive"`
	DurationSeconds int64              `json:"duration_seconds"`
}

// ExamAnswerReview is returned after submission.
type ExamAnswerReview struct {
	QuestionID        uint   `json:"question_id"`
	Stem              string `json:"stem"`
	Type              string `json:"type"`
	Score             int    `json:"score"`
	ObtainedScore     int    `json:"obtained_score"`
	IsCorrect         bool   `json:"is_correct"`
	SelectedOptionIDs []uint `json:"selected_option_ids"`
	CorrectOptionIDs  []uint `json:"correct_option_ids"`
}

// ExamSubmitResponse is returned when exam submission succeeds.
type ExamSubmitResponse struct {
	AttemptID       uint               `json:"attempt_id"`
	ExamID          uint               `json:"exam_id"`
	Score           int                `json:"score"`
	TotalScore      int                `json:"total_score"`
	Pass            bool               `json:"pass"`
	CorrectCount    int                `json:"correct_count"`
	TotalCount      int                `json:"total_count"`
	DurationSeconds int64              `json:"duration_seconds"`
	Answers         []ExamAnswerReview `json:"answers"`
}

// ExamResultSummary represents a simplified attempt record.
type ExamResultSummary struct {
	AttemptID   uint      `json:"attempt_id"`
	ExamID      uint      `json:"exam_id"`
	ExamTitle   string    `json:"exam_title"`
	Score       int       `json:"score"`
	TotalScore  int       `json:"total_score"`
	PassScore   int       `json:"pass_score"`
	Pass        bool      `json:"pass"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// ManagerExamProgressItem summarises exam stats for manager dashboard.
type ManagerExamProgressItem struct {
	ExamID       uint    `json:"exam_id"`
	Title        string  `json:"title"`
	AttemptCount int64   `json:"attempt_count"`
	PassRate     float64 `json:"pass_rate"`
	AvgScore     float64 `json:"avg_score"`
}

// EmployeeLatestExamResult summarises latest exam for employee.
type EmployeeLatestExamResult struct {
	ExamID      uint      `json:"exam_id"`
	ExamTitle   string    `json:"exam_title"`
	Score       int       `json:"score"`
	Pass        bool      `json:"pass"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// EmployeeLearningProgress is returned for manager to view progress.
type EmployeeLearningProgress struct {
	Completed int `json:"completed"`
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Percent   int `json:"percent"`
}

// ManagerEmployeeExamRecord summarises each employee row.
type ManagerEmployeeExamRecord struct {
	EmployeeID       uint                      `json:"employee_id"`
	Name             string                    `json:"name"`
	WorkNo           string                    `json:"work_no"`
	LatestExam       *EmployeeLatestExamResult `json:"latest_exam"`
	LearningProgress EmployeeLearningProgress  `json:"learning_progress"`
}

// ManagerExamOverviewResponse returns exam & learning overview for manager.
type ManagerExamOverviewResponse struct {
	ExamProgress []ManagerExamProgressItem   `json:"exam_progress"`
	Employees    []ManagerEmployeeExamRecord `json:"employees"`
}

// AdminExamOverviewQuery controls filtering and pagination for admin exam overview.
type AdminExamOverviewQuery struct {
	Role      string `form:"role" binding:"omitempty,oneof=employee manager admin all" example:"all"`
	ManagerID uint   `form:"manager_id" binding:"omitempty,min=1" example:"2"`
	Keyword   string `form:"keyword" binding:"omitempty,max=100" example:"张三"`
	Page      int    `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
}

// AdminUserExamRecord summarises each user row for admin overview.
type AdminUserExamRecord struct {
	UserID           uint                      `json:"user_id"`
	Name             string                    `json:"name"`
	WorkNo           string                    `json:"work_no"`
	Role             string                    `json:"role"`
	LatestExam       *EmployeeLatestExamResult `json:"latest_exam"`
	LearningProgress EmployeeLearningProgress  `json:"learning_progress"`
}

// AdminExamOverviewResponse returns exam & learning overview for admin.
type AdminExamOverviewResponse struct {
	ExamProgress []ManagerExamProgressItem `json:"exam_progress"`
	Users        []AdminUserExamRecord     `json:"users"`
	Pagination   Pagination                `json:"pagination"`
}
