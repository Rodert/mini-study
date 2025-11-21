package service

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
)

// UserService handles user business logic.
type UserService struct {
	repo         *repository.UserRepository
	relationRepo *repository.ManagerEmployeeRepository
	audit        *AuditService
	points       *PointService
}

// NewUserService builds a user service.
func NewUserService(userRepo *repository.UserRepository, relationRepo *repository.ManagerEmployeeRepository, audit *AuditService, pointSvc *PointService) *UserService {
	return &UserService{repo: userRepo, relationRepo: relationRepo, audit: audit, points: pointSvc}
}

// Register creates a new user.
func (s *UserService) Register(req dto.RegisterRequest) (*model.User, error) {
	if _, err := s.repo.FindByWorkNo(req.WorkNo); err == nil {
		return nil, errors.New("工号已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		WorkNo:       req.WorkNo,
		Phone:        req.Phone,
		Name:         req.Name,
		Role:         model.RoleEmployee,
		PasswordHash: hash,
		Status:       true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// 将传入的店长工号转换为店长 ID，并建立店长-员工关系
	managerIDs, err := s.resolveManagerWorkNos(req.ManagerWorkNos)
	if err != nil {
		return nil, err
	}
	if err := s.relationRepo.CreateRelations(user.ID, managerIDs); err != nil {
		return nil, err
	}

	_ = s.audit.Record(user.ID, "register", "users", utils.ToJSONString(req), http.StatusText(http.StatusCreated))
	return user, nil
}

// Login validates user credentials and returns the user.
func (s *UserService) Login(req dto.LoginRequest) (*model.User, error) {
	user, err := s.repo.FindByWorkNo(req.WorkNo)
	if err != nil {
		return nil, errors.New("工号或密码错误")
	}

	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("工号或密码错误")
	}

	_ = s.audit.Record(user.ID, "login", "users", "{}", http.StatusText(http.StatusOK))
	return user, nil
}

// ListManagers returns all active managers.
func (s *UserService) ListManagers() ([]model.User, error) {
	return s.repo.ListManagers()
}

// AdminListUsers returns users for admin panel.
func (s *UserService) AdminListUsers(adminID uint, filter dto.AdminListUsersQuery) ([]dto.AdminUserResponse, error) {
	if _, err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	users, err := s.repo.ListUsers(filter.Role, filter.Keyword)
	if err != nil {
		return nil, err
	}

	return s.buildAdminUserResponses(users)
}

// AdminGetUser returns a single user for admin panel.
func (s *UserService) AdminGetUser(adminID, targetID uint) (*dto.AdminUserResponse, error) {
	if _, err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, err
	}

	resp, err := s.buildAdminUserResponses([]model.User{*user})
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("用户不存在")
	}

	return &resp[0], nil
}

// AdminUpdateUserRole updates role for a specific user.
func (s *UserService) AdminUpdateUserRole(adminID, targetID uint, role model.Role) (*model.User, error) {
	if _, err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	if role != model.RoleEmployee && role != model.RoleManager && role != model.RoleAdmin {
		return nil, errors.New("无效的角色")
	}

	user, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, err
	}

	user.Role = role
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	if role != model.RoleEmployee {
		if err := s.relationRepo.ReplaceRelations(user.ID, nil); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *UserService) buildAdminUserResponses(users []model.User) ([]dto.AdminUserResponse, error) {
	if len(users) == 0 {
		return []dto.AdminUserResponse{}, nil
	}

	resp := make([]dto.AdminUserResponse, 0, len(users))
	employeeIDs := make([]uint, 0)
	allUserIDs := make([]uint, 0, len(users))
	for _, user := range users {
		allUserIDs = append(allUserIDs, user.ID)
		if user.Role == model.RoleEmployee {
			employeeIDs = append(employeeIDs, user.ID)
		}
	}

	relations, err := s.relationRepo.ListByEmployeeIDs(employeeIDs)
	if err != nil {
		return nil, err
	}

	bindingMap := make(map[uint][]uint)
	managerIDSet := make(map[uint]struct{})
	for _, rel := range relations {
		bindingMap[rel.EmployeeID] = append(bindingMap[rel.EmployeeID], rel.ManagerID)
		managerIDSet[rel.ManagerID] = struct{}{}
	}

	managerIDs := make([]uint, 0, len(managerIDSet))
	for id := range managerIDSet {
		managerIDs = append(managerIDs, id)
	}

	managerUsers, err := s.repo.FindByIDs(managerIDs)
	if err != nil {
		return nil, err
	}

	managerMap := make(map[uint]dto.ManagerBrief, len(managerUsers))
	for _, manager := range managerUsers {
		managerMap[manager.ID] = dto.ManagerBrief{
			ID:     manager.ID,
			WorkNo: manager.WorkNo,
			Name:   manager.Name,
			Phone:  manager.Phone,
		}
	}

	var pointTotals map[uint]int64
	if s.points != nil && len(allUserIDs) > 0 {
		pt, err := s.points.GetTotalsMap(allUserIDs)
		if err != nil {
			return nil, err
		}
		pointTotals = pt
	}

	for _, user := range users {
		boundManagerIDs := bindingMap[user.ID]
		var copiedIDs []uint
		var managerBriefs []dto.ManagerBrief

		if len(boundManagerIDs) > 0 {
			copiedIDs = make([]uint, 0, len(boundManagerIDs))
			managerBriefs = make([]dto.ManagerBrief, 0, len(boundManagerIDs))
			for _, mid := range boundManagerIDs {
				copiedIDs = append(copiedIDs, mid)
				if brief, ok := managerMap[mid]; ok {
					managerBriefs = append(managerBriefs, brief)
				}
			}
		} else {
			copiedIDs = []uint{}
			managerBriefs = []dto.ManagerBrief{}
		}

		points := int64(0)
		if pointTotals != nil {
			if val, ok := pointTotals[user.ID]; ok {
				points = val
			}
		}

		resp = append(resp, dto.AdminUserResponse{
			UserResponse: dto.UserResponse{
				ID:     user.ID,
				WorkNo: user.WorkNo,
				Phone:  user.Phone,
				Name:   user.Name,
				Role:   user.Role,
				Status: user.Status,
			},
			ManagerIDs: copiedIDs,
			Managers:   managerBriefs,
			Points:     points,
		})
	}

	return resp, nil
}

func (s *UserService) resolveManagerWorkNos(workNos []string) ([]uint, error) {
	if len(workNos) == 0 {
		return nil, nil
	}

	var managerIDs []uint
	for _, workNo := range workNos {
		manager, err := s.repo.FindByWorkNo(workNo)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("店长工号不存在: " + workNo)
			}
			return nil, err
		}
		if manager.Role != model.RoleManager {
			return nil, errors.New("用户不是店长角色: " + workNo)
		}
		managerIDs = append(managerIDs, manager.ID)
	}
	return managerIDs, nil
}

// ensureAdmin checks whether a given user is admin.
func (s *UserService) ensureAdmin(userID uint) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user.Role != model.RoleAdmin {
		return nil, errors.New("无权限，仅管理员可操作")
	}
	return user, nil
}

// GetCurrentUser returns the current user by ID with manager information if the user is an employee.
func (s *UserService) GetCurrentUser(userID uint) (*dto.AdminUserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	resp, err := s.buildAdminUserResponses([]model.User{*user})
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("用户不存在")
	}

	return &resp[0], nil
}

// UpdateProfile allows user to update own name and phone.
func (s *UserService) UpdateProfile(userID uint, req dto.UpdateProfileRequest) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// CreateManager creates a new manager user; only admin can call this.
func (s *UserService) CreateManager(adminID uint, req dto.AdminCreateManagerRequest) (*model.User, error) {
	if _, err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	if _, err := s.repo.FindByWorkNo(req.WorkNo); err == nil {
		return nil, errors.New("工号已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		WorkNo:       req.WorkNo,
		Phone:        req.Phone,
		Name:         req.Name,
		Role:         model.RoleManager,
		PasswordHash: hash,
		Status:       true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// CreateEmployee creates a new employee user; only admin can call this.
func (s *UserService) CreateEmployee(adminID uint, req dto.AdminCreateEmployeeRequest) (*model.User, error) {
	if _, err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	if _, err := s.repo.FindByWorkNo(req.WorkNo); err == nil {
		return nil, errors.New("工号已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		WorkNo:       req.WorkNo,
		Phone:        req.Phone,
		Name:         req.Name,
		Role:         model.RoleEmployee,
		PasswordHash: hash,
		Status:       true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// 将传入的店长工号转换为店长 ID，并建立店长-员工关系
	managerIDs, err := s.resolveManagerWorkNos(req.ManagerWorkNos)
	if err != nil {
		return nil, err
	}
	if len(managerIDs) > 0 {
		if err := s.relationRepo.CreateRelations(user.ID, managerIDs); err != nil {
			return nil, err
		}
	}

	_ = s.audit.Record(adminID, "create_employee", "users", utils.ToJSONString(req), http.StatusText(http.StatusCreated))
	return user, nil
}

// PromoteToManager changes an existing employee to manager; only admin can call.
func (s *UserService) PromoteToManager(adminID, targetUserID uint) (*model.User, error) {
	if _, err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByID(targetUserID)
	if err != nil {
		return nil, err
	}

	if user.Role == model.RoleAdmin {
		return nil, errors.New("不能修改管理员角色")
	}

	user.Role = model.RoleManager

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateEmployeeManagers updates the manager bindings for an employee or manager (non-admin).
func (s *UserService) UpdateEmployeeManagers(adminID, targetUserID uint, managerWorkNos []string) (*model.User, error) {
	if _, err := s.ensureAdmin(adminID); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByID(targetUserID)
	if err != nil {
		return nil, err
	}

	if user.Role == model.RoleAdmin {
		return nil, errors.New("不能修改管理员的店长绑定")
	}

	managerIDs, err := s.resolveManagerWorkNos(managerWorkNos)
	if err != nil {
		return nil, err
	}

	if err := s.relationRepo.ReplaceRelations(user.ID, managerIDs); err != nil {
		return nil, err
	}

	return user, nil
}
