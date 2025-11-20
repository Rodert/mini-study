package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/model"
)

// ExamRepository handles CRUD for exam papers and questions.
type ExamRepository struct {
	db *gorm.DB
}

// NewExamRepository builds a new ExamRepository.
func NewExamRepository(db *gorm.DB) *ExamRepository {
	return &ExamRepository{db: db}
}

// Create creates an exam with associations.
func (r *ExamRepository) Create(exam *model.ExamPaper) error {
	if err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(exam).Error; err != nil {
		return errors.Wrap(err, "create exam")
	}
	return nil
}

// UpdateExam updates exam metadata (without questions).
func (r *ExamRepository) UpdateExam(exam *model.ExamPaper) error {
	if err := r.db.Model(&model.ExamPaper{}).Where("id = ?", exam.ID).
		Updates(map[string]interface{}{
			"title":              exam.Title,
			"description":        exam.Description,
			"status":             exam.Status,
			"target_role":        exam.TargetRole,
			"time_limit_minutes": exam.TimeLimitMinutes,
			"pass_score":         exam.PassScore,
			"total_score":        exam.TotalScore,
		}).Error; err != nil {
		return errors.Wrap(err, "update exam")
	}
	return nil
}

// ReplaceQuestions removes existing questions/options and inserts new ones.
func (r *ExamRepository) ReplaceQuestions(examID uint, questions []model.ExamQuestion) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先获取该考试的所有题目ID
		var questionIDs []uint
		if err := tx.Model(&model.ExamQuestion{}).
			Where("exam_id = ?", examID).
			Pluck("id", &questionIDs).Error; err != nil {
			return errors.Wrap(err, "find exam question ids")
		}

		// 删除这些题目对应的所有选项
		if len(questionIDs) > 0 {
			if err := tx.Where("question_id IN ?", questionIDs).Delete(&model.ExamOption{}).Error; err != nil {
				return errors.Wrap(err, "delete exam options")
			}
		}

		// 删除所有题目
		if err := tx.Where("exam_id = ?", examID).Delete(&model.ExamQuestion{}).Error; err != nil {
			return errors.Wrap(err, "delete exam questions")
		}

		// 创建新题目和选项
		if len(questions) == 0 {
			return nil
		}
		for idx := range questions {
			questions[idx].ExamID = examID
		}
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&questions).Error; err != nil {
			return errors.Wrap(err, "create exam questions")
		}
		return nil
	})
}

// FindByID returns exam without associations.
func (r *ExamRepository) FindByID(id uint) (*model.ExamPaper, error) {
	var exam model.ExamPaper
	if err := r.db.First(&exam, id).Error; err != nil {
		return nil, errors.Wrap(err, "find exam by id")
	}
	return &exam, nil
}

// FindWithQuestions loads exam along with questions and options.
func (r *ExamRepository) FindWithQuestions(id uint) (*model.ExamPaper, error) {
	var exam model.ExamPaper
	if err := r.db.
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("exam_questions.id ASC")
		}).
		Preload("Questions.Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("exam_options.sort_order ASC, exam_options.id ASC")
		}).
		First(&exam, id).Error; err != nil {
		return nil, errors.Wrap(err, "find exam with questions")
	}
	return &exam, nil
}

// FindPublishedWithQuestions ensures exam is published and loads associations.
func (r *ExamRepository) FindPublishedWithQuestions(id uint) (*model.ExamPaper, error) {
	var exam model.ExamPaper
	if err := r.db.
		Where("status = ?", "published").
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("exam_questions.id ASC")
		}).
		Preload("Questions.Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("exam_options.sort_order ASC, exam_options.id ASC")
		}).
		First(&exam, id).Error; err != nil {
		return nil, errors.Wrap(err, "find published exam with questions")
	}
	return &exam, nil
}

// ListPublishedByRole lists published exams that match role scope.
func (r *ExamRepository) ListPublishedByRole(role model.Role) ([]model.ExamPaper, error) {
	var exams []model.ExamPaper
	query := r.db.
		Where("status = ?", "published").
		Order("id DESC")

	query = query.Where("target_role = ? OR target_role = ?", role, "all")

	if err := query.Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, exam_id").Order("exam_questions.id ASC")
	}).Find(&exams).Error; err != nil {
		return nil, errors.Wrap(err, "list published exams")
	}

	for idx := range exams {
		exams[idx].QuestionCount = len(exams[idx].Questions)
		exams[idx].Questions = nil
	}

	return exams, nil
}

// CountQuestions returns how many questions belong to an exam.
func (r *ExamRepository) CountQuestions(examID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.ExamQuestion{}).Where("exam_id = ?", examID).Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, "count exam questions")
	}
	return count, nil
}

// ListAll loads all exams for admin with question count.
func (r *ExamRepository) ListAll() ([]model.ExamPaper, error) {
	var exams []model.ExamPaper
	if err := r.db.Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, exam_id").Order("exam_questions.id ASC")
	}).Find(&exams).Error; err != nil {
		return nil, errors.Wrap(err, "list all exams")
	}
	for idx := range exams {
		exams[idx].QuestionCount = len(exams[idx].Questions)
		exams[idx].Questions = nil
	}
	return exams, nil
}

// FindByIDs returns exams by ID list.
func (r *ExamRepository) FindByIDs(ids []uint) ([]model.ExamPaper, error) {
	if len(ids) == 0 {
		return []model.ExamPaper{}, nil
	}
	var exams []model.ExamPaper
	if err := r.db.Where("id IN ?", ids).Find(&exams).Error; err != nil {
		return nil, errors.Wrap(err, "find exams by ids")
	}
	return exams, nil
}
