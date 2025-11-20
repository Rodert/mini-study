package service

import (
	"encoding/json"
	"errors"
	"math"
	"sort"
	"time"

	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
)

// ExamService handles exam workflows.
type ExamService struct {
	exams     *repository.ExamRepository
	attempts  *repository.ExamAttemptRepository
	users     *repository.UserRepository
	relations *repository.ManagerEmployeeRepository
	learning  *repository.LearningRecordRepository
	contents  *repository.ContentRepository
}

// NewExamService builds ExamService.
func NewExamService(
	examRepo *repository.ExamRepository,
	attemptRepo *repository.ExamAttemptRepository,
	userRepo *repository.UserRepository,
	relationRepo *repository.ManagerEmployeeRepository,
	learningRepo *repository.LearningRecordRepository,
	contentRepo *repository.ContentRepository,
) *ExamService {
	return &ExamService{
		exams:     examRepo,
		attempts:  attemptRepo,
		users:     userRepo,
		relations: relationRepo,
		learning:  learningRepo,
		contents:  contentRepo,
	}
}

// AdminCreateExam allows admin to create an exam paper.
func (s *ExamService) AdminCreateExam(adminID uint, req dto.AdminExamUpsertRequest) (*dto.ExamDetailResponse, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	questions, totalScore, err := s.buildQuestionModels(req.Questions)
	if err != nil {
		return nil, err
	}

	exam := &model.ExamPaper{
		Title:            req.Title,
		Description:      req.Description,
		Status:           s.normalizeExamStatus(req.Status),
		TargetRole:       model.Role(s.normalizeTargetRole(req.TargetRole)),
		TimeLimitMinutes: req.TimeLimitMinutes,
		PassScore:        req.PassScore,
		TotalScore:       totalScore,
		CreatorID:        adminID,
		Questions:        questions,
	}

	if exam.PassScore > exam.TotalScore {
		return nil, errors.New("及格分不能高于总分")
	}

	if err := s.exams.Create(exam); err != nil {
		return nil, err
	}

	return s.buildExamDetailDTO(exam), nil
}

// AdminListExams returns all exams for admin.
func (s *ExamService) AdminListExams(adminID uint) ([]dto.ExamDetailResponse, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	exams, err := s.exams.ListAll()
	if err != nil {
		return nil, err
	}

	resp := make([]dto.ExamDetailResponse, 0, len(exams))
	for _, exam := range exams {
		examWithQuestions, err := s.exams.FindWithQuestions(exam.ID)
		if err != nil {
			continue
		}
		// Include correct answers for admin list view
		resp = append(resp, *s.buildExamDetailDTOWithAnswers(examWithQuestions, true))
	}

	return resp, nil
}

// AdminGetExam returns exam detail for admin editing.
func (s *ExamService) AdminGetExam(adminID, examID uint) (*dto.ExamDetailResponse, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	exam, err := s.exams.FindWithQuestions(examID)
	if err != nil {
		return nil, err
	}

	// Include correct answers for admin editing
	return s.buildExamDetailDTOWithAnswers(exam, true), nil
}

// AdminUpdateExam allows admin to update exam and questions.
func (s *ExamService) AdminUpdateExam(adminID, examID uint, req dto.AdminExamUpsertRequest) (*dto.ExamDetailResponse, error) {
	if err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	exam, err := s.exams.FindWithQuestions(examID)
	if err != nil {
		return nil, err
	}

	questions, totalScore, err := s.buildQuestionModels(req.Questions)
	if err != nil {
		return nil, err
	}

	exam.Title = req.Title
	exam.Description = req.Description
	if req.Status != "" {
		exam.Status = s.normalizeExamStatus(req.Status)
	}
	if req.TargetRole != "" {
		exam.TargetRole = model.Role(s.normalizeTargetRole(req.TargetRole))
	}
	if req.TimeLimitMinutes >= 0 {
		exam.TimeLimitMinutes = req.TimeLimitMinutes
	}
	exam.PassScore = req.PassScore
	exam.TotalScore = totalScore

	if exam.PassScore > exam.TotalScore {
		return nil, errors.New("及格分不能高于总分")
	}

	if err := s.exams.UpdateExam(exam); err != nil {
		return nil, err
	}
	if err := s.exams.ReplaceQuestions(exam.ID, questions); err != nil {
		return nil, err
	}

	updated, err := s.exams.FindWithQuestions(exam.ID)
	if err != nil {
		return nil, err
	}

	return s.buildExamDetailDTO(updated), nil
}

// ListAvailableExams returns published exams for current user.
func (s *ExamService) ListAvailableExams(userID uint) ([]dto.ExamListItem, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}

	exams, err := s.exams.ListPublishedByRole(user.Role)
	if err != nil {
		return nil, err
	}

	attempts, err := s.attempts.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	latestByExam := make(map[uint]model.ExamAttempt)
	for _, attempt := range attempts {
		if _, exists := latestByExam[attempt.ExamID]; exists {
			continue
		}
		latestByExam[attempt.ExamID] = attempt
	}

	resp := make([]dto.ExamListItem, 0, len(exams))
	for _, exam := range exams {
		item := dto.ExamListItem{
			ID:               exam.ID,
			Title:            exam.Title,
			Description:      exam.Description,
			TotalScore:       exam.TotalScore,
			PassScore:        exam.PassScore,
			TimeLimitMinutes: exam.TimeLimitMinutes,
			QuestionCount:    exam.QuestionCount,
			Status:           exam.Status,
			AttemptStatus:    "not_started",
		}

		if attempt, ok := latestByExam[exam.ID]; ok {
			item.LastScore = attempt.Score
			item.LastPassed = attempt.Pass
			if attempt.SubmittedAt != nil {
				item.LastSubmittedAt = attempt.SubmittedAt
			} else {
				submitted := attempt.CreatedAt
				item.LastSubmittedAt = &submitted
			}

			switch {
			case attempt.Pass:
				item.AttemptStatus = "passed"
			default:
				item.AttemptStatus = "attempted"
			}
		}

		resp = append(resp, item)
	}

	return resp, nil
}

// GetExamDetail returns exam detail for answering.
func (s *ExamService) GetExamDetail(userID, examID uint) (*dto.ExamDetailResponse, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}

	exam, err := s.exams.FindPublishedWithQuestions(examID)
	if err != nil {
		return nil, err
	}

	if err := s.ensureExamAccessible(user.Role, exam); err != nil {
		return nil, err
	}

	return s.buildExamDetailDTO(exam), nil
}

// SubmitExam evaluates answers and stores attempt.
func (s *ExamService) SubmitExam(userID, examID uint, req dto.ExamSubmitRequest) (*dto.ExamSubmitResponse, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}

	exam, err := s.exams.FindPublishedWithQuestions(examID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureExamAccessible(user.Role, exam); err != nil {
		return nil, err
	}

	// 检查用户是否已经参加过这个考试
	existingAttempt, err := s.attempts.FindLatestByUserAndExam(userID, examID)
	if err == nil && existingAttempt != nil {
		return nil, errors.New("您已经参加过该考试，不能重复参加")
	}
	// 如果错误是记录不存在，可以继续；其他错误需要返回
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if len(req.Answers) != len(exam.Questions) {
		return nil, errors.New("请完成所有题目后再提交")
	}

	answerMap := make(map[uint][]uint, len(req.Answers))
	for _, ans := range req.Answers {
		if len(ans.OptionIDs) == 0 {
			return nil, errors.New("请选择该题的答案")
		}
		answerMap[ans.QuestionID] = s.normalizeOptionIDs(ans.OptionIDs)
	}

	totalScore := 0
	correctCount := 0
	reviews := make([]dto.ExamAnswerReview, 0, len(exam.Questions))

	for idx := range exam.Questions {
		question := exam.Questions[idx]
		selected, ok := answerMap[question.ID]
		if !ok {
			return nil, errors.New("存在未作答的题目")
		}

		correctOptionIDs := s.extractCorrectOptionIDs(question.Options)
		if len(correctOptionIDs) == 0 {
			return nil, errors.New("试题未配置正确答案")
		}
		if question.Type == "single" && len(correctOptionIDs) != 1 {
			return nil, errors.New("单选题必须设置唯一正确答案")
		}

		if err := s.ensureSelectedOptionsValid(question, selected); err != nil {
			return nil, err
		}

		isCorrect := s.compareOptionSets(selected, correctOptionIDs)
		earned := 0
		if isCorrect {
			earned = question.Score
			correctCount++
		}
		totalScore += earned

		reviews = append(reviews, dto.ExamAnswerReview{
			QuestionID:        question.ID,
			Stem:              question.Stem,
			Type:              question.Type,
			Score:             question.Score,
			ObtainedScore:     earned,
			IsCorrect:         isCorrect,
			SelectedOptionIDs: selected,
			CorrectOptionIDs:  correctOptionIDs,
		})
	}

	pass := totalScore >= exam.PassScore
	now := time.Now()
	payload, err := json.Marshal(reviews)
	if err != nil {
		return nil, err
	}

	attempt := &model.ExamAttempt{
		ExamID:          exam.ID,
		UserID:          userID,
		Status:          "submitted",
		Score:           totalScore,
		CorrectCount:    correctCount,
		TotalCount:      len(exam.Questions),
		Pass:            pass,
		DurationSeconds: req.DurationSeconds,
		AnswerSnapshot:  payload,
		SubmittedAt:     &now,
	}

	if err := s.attempts.Create(attempt); err != nil {
		return nil, err
	}

	return &dto.ExamSubmitResponse{
		AttemptID:       attempt.ID,
		ExamID:          exam.ID,
		Score:           totalScore,
		TotalScore:      exam.TotalScore,
		Pass:            pass,
		CorrectCount:    correctCount,
		TotalCount:      len(exam.Questions),
		DurationSeconds: req.DurationSeconds,
		Answers:         reviews,
	}, nil
}

// ListMyResults returns attempt summaries for user.
func (s *ExamService) ListMyResults(userID uint) ([]dto.ExamResultSummary, error) {
	attempts, err := s.attempts.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	results := make([]dto.ExamResultSummary, 0, len(attempts))
	for _, attempt := range attempts {
		if attempt.Exam.ID == 0 {
			continue
		}
		submittedAt := attempt.CreatedAt
		if attempt.SubmittedAt != nil {
			submittedAt = *attempt.SubmittedAt
		}
		results = append(results, dto.ExamResultSummary{
			AttemptID:   attempt.ID,
			ExamID:      attempt.ExamID,
			ExamTitle:   attempt.Exam.Title,
			Score:       attempt.Score,
			TotalScore:  attempt.Exam.TotalScore,
			PassScore:   attempt.Exam.PassScore,
			Pass:        attempt.Pass,
			SubmittedAt: submittedAt,
		})
	}

	return results, nil
}

// GetManagerOverview returns learning/exam summary for manager employees.
func (s *ExamService) GetManagerOverview(managerID uint) (*dto.ManagerExamOverviewResponse, error) {
	manager, err := s.users.FindByID(managerID)
	if err != nil {
		return nil, err
	}
	if manager.Role != model.RoleManager && manager.Role != model.RoleAdmin {
		return nil, errors.New("仅店长可查看")
	}

	employeeIDs, err := s.relations.ListEmployeeIDsByManager(managerID)
	if err != nil {
		return nil, err
	}
	if len(employeeIDs) == 0 {
		return &dto.ManagerExamOverviewResponse{
			ExamProgress: []dto.ManagerExamProgressItem{},
			Employees:    []dto.ManagerEmployeeExamRecord{},
		}, nil
	}

	employees, err := s.users.FindByIDs(employeeIDs)
	if err != nil {
		return nil, err
	}

	learningAgg, err := s.learning.AggregateByUsers(employeeIDs)
	if err != nil {
		return nil, err
	}

	latestAttempts, err := s.attempts.ListLatestByUsers(employeeIDs)
	if err != nil {
		return nil, err
	}

	examStats, err := s.attempts.AggregateByExamForUsers(employeeIDs)
	if err != nil {
		return nil, err
	}

	totalContents, err := s.contents.CountPublishedForRole(string(model.RoleEmployee))
	if err != nil {
		return nil, err
	}

	examIDSet := make(map[uint]struct{}, len(examStats))
	examIDs := make([]uint, 0, len(examStats))
	for _, row := range examStats {
		if _, exists := examIDSet[row.ExamID]; exists {
			continue
		}
		examIDSet[row.ExamID] = struct{}{}
		examIDs = append(examIDs, row.ExamID)
	}

	examInfos := make(map[uint]model.ExamPaper)
	if len(examIDs) > 0 {
		exams, err := s.exams.FindByIDs(examIDs)
		if err != nil {
			return nil, err
		}
		for _, exam := range exams {
			examInfos[exam.ID] = exam
		}
	}

	progressList := make([]dto.ManagerExamProgressItem, 0, len(examStats))
	for _, row := range examStats {
		info, ok := examInfos[row.ExamID]
		if !ok {
			continue
		}
		passRate := 0.0
		if row.AttemptCount > 0 {
			passRate = float64(row.PassCount) / float64(row.AttemptCount)
		}
		progressList = append(progressList, dto.ManagerExamProgressItem{
			ExamID:       row.ExamID,
			Title:        info.Title,
			AttemptCount: row.AttemptCount,
			PassRate:     math.Round(passRate*1000) / 1000,
			AvgScore:     math.Round(row.AvgScore*10) / 10,
		})
	}
	sort.Slice(progressList, func(i, j int) bool {
		return progressList[i].AttemptCount > progressList[j].AttemptCount
	})

	sort.Slice(employees, func(i, j int) bool {
		return employees[i].Name < employees[j].Name
	})

	employeeRecords := make([]dto.ManagerEmployeeExamRecord, 0, len(employees))
	total := int(totalContents)

	for _, emp := range employees {
		agg := learningAgg[emp.ID]
		completed := int(agg.Completed)
		pending := total - completed
		if pending < 0 {
			pending = 0
		}
		percent := 0
		if total > 0 {
			percent = int(math.Round(float64(completed) * 100 / float64(total)))
		}

		progress := dto.EmployeeLearningProgress{
			Completed: completed,
			Total:     total,
			Pending:   pending,
			Percent:   percent,
		}

		var latest *dto.EmployeeLatestExamResult
		if att, ok := latestAttempts[emp.ID]; ok {
			if att.Exam.ID > 0 {
				submitted := att.CreatedAt
				if att.SubmittedAt != nil {
					submitted = *att.SubmittedAt
				}
				latest = &dto.EmployeeLatestExamResult{
					ExamID:      att.ExamID,
					ExamTitle:   att.Exam.Title,
					Score:       att.Score,
					Pass:        att.Pass,
					SubmittedAt: submitted,
				}
			}
		}

		employeeRecords = append(employeeRecords, dto.ManagerEmployeeExamRecord{
			EmployeeID:       emp.ID,
			Name:             emp.Name,
			WorkNo:           emp.WorkNo,
			LatestExam:       latest,
			LearningProgress: progress,
		})
	}

	return &dto.ManagerExamOverviewResponse{
		ExamProgress: progressList,
		Employees:    employeeRecords,
	}, nil
}

func (s *ExamService) ensureAdmin(userID uint) error {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return err
	}
	if user.Role != model.RoleAdmin {
		return errors.New("仅管理员可操作")
	}
	return nil
}

func (s *ExamService) ensureExamAccessible(role model.Role, exam *model.ExamPaper) error {
	if exam.TargetRole == "all" {
		return nil
	}
	if role == model.RoleAdmin {
		return nil
	}
	if string(role) != exam.TargetRole {
		return errors.New("当前考试不在您的学习范围内")
	}
	return nil
}

func (s *ExamService) normalizeExamStatus(status string) string {
	switch status {
	case "published", "archived":
		return status
	default:
		return "draft"
	}
}

func (s *ExamService) normalizeTargetRole(role string) string {
	switch role {
	case "manager", "all":
		return role
	default:
		return "employee"
	}
}

func (s *ExamService) buildExamDetailDTO(exam *model.ExamPaper) *dto.ExamDetailResponse {
	return s.buildExamDetailDTOWithAnswers(exam, false)
}

// buildExamDetailDTOWithAnswers builds exam DTO, optionally including correct answers.
func (s *ExamService) buildExamDetailDTOWithAnswers(exam *model.ExamPaper, includeAnswers bool) *dto.ExamDetailResponse {
	resp := &dto.ExamDetailResponse{
		ID:               exam.ID,
		Title:            exam.Title,
		Description:      exam.Description,
		TimeLimitMinutes: exam.TimeLimitMinutes,
		PassScore:        exam.PassScore,
		TotalScore:       exam.TotalScore,
		QuestionCount:    len(exam.Questions),
		Questions:        make([]dto.ExamDetailQuestion, 0, len(exam.Questions)),
	}

	for _, question := range exam.Questions {
		opts := make([]dto.ExamDetailQuestionOption, 0, len(question.Options))
		for _, opt := range question.Options {
			optDTO := dto.ExamDetailQuestionOption{
				ID:      opt.ID,
				Label:   opt.Label,
				Content: opt.Content,
			}
			if includeAnswers {
				optDTO.IsCorrect = opt.IsCorrect
			}
			opts = append(opts, optDTO)
		}
		questionDTO := dto.ExamDetailQuestion{
			ID:      question.ID,
			Type:    question.Type,
			Stem:    question.Stem,
			Score:   question.Score,
			Options: opts,
		}
		if includeAnswers {
			questionDTO.Analysis = question.Analysis
		}
		resp.Questions = append(resp.Questions, questionDTO)
	}
	return resp
}

func (s *ExamService) buildQuestionModels(payload []dto.AdminExamQuestionUpsert) ([]model.ExamQuestion, int, error) {
	questions := make([]model.ExamQuestion, 0, len(payload))
	totalScore := 0

	for _, item := range payload {
		if len(item.Options) < 2 {
			return nil, 0, errors.New("每道题至少需要两个选项")
		}
		hasCorrect := false
		options := make([]model.ExamOption, 0, len(item.Options))
		for idx, opt := range item.Options {
			label := opt.Label
			if label == "" {
				label = string('A' + rune(idx))
			}
			sortOrder := opt.SortOrder
			if sortOrder == 0 {
				sortOrder = idx
			}
			options = append(options, model.ExamOption{
				Label:     label,
				Content:   opt.Content,
				IsCorrect: opt.IsCorrect,
				SortOrder: sortOrder,
			})
			if opt.IsCorrect {
				hasCorrect = true
			}
		}
		if !hasCorrect {
			return nil, 0, errors.New("每道题至少需要一个正确答案")
		}

		question := model.ExamQuestion{
			Type:     item.Type,
			Stem:     item.Stem,
			Score:    item.Score,
			Analysis: item.Analysis,
			Options:  options,
		}
		totalScore += item.Score
		questions = append(questions, question)
	}

	return questions, totalScore, nil
}

func (s *ExamService) extractCorrectOptionIDs(options []model.ExamOption) []uint {
	ids := make([]uint, 0, len(options))
	for _, opt := range options {
		if opt.IsCorrect {
			ids = append(ids, opt.ID)
		}
	}
	return ids
}

func (s *ExamService) normalizeOptionIDs(optionIDs []uint) []uint {
	set := make(map[uint]struct{}, len(optionIDs))
	for _, id := range optionIDs {
		if id == 0 {
			continue
		}
		set[id] = struct{}{}
	}

	normalized := make([]uint, 0, len(set))
	for id := range set {
		normalized = append(normalized, id)
	}
	sort.Slice(normalized, func(i, j int) bool { return normalized[i] < normalized[j] })
	return normalized
}

func (s *ExamService) compareOptionSets(selected, correct []uint) bool {
	if len(selected) != len(correct) {
		return false
	}
	correctSet := make(map[uint]struct{}, len(correct))
	for _, id := range correct {
		correctSet[id] = struct{}{}
	}
	for _, id := range selected {
		if _, ok := correctSet[id]; !ok {
			return false
		}
	}
	return true
}

func (s *ExamService) ensureSelectedOptionsValid(question model.ExamQuestion, selected []uint) error {
	if question.Type == "single" && len(selected) != 1 {
		return errors.New("单选题只能选择一个选项")
	}

	allowed := make(map[uint]struct{}, len(question.Options))
	for _, opt := range question.Options {
		allowed[opt.ID] = struct{}{}
	}
	for _, id := range selected {
		if _, ok := allowed[id]; !ok {
			return errors.New("存在非法的选项")
		}
	}
	return nil
}
