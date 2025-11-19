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
}

// NewUserService builds a user service.
func NewUserService(userRepo *repository.UserRepository, relationRepo *repository.ManagerEmployeeRepository, audit *AuditService) *UserService {
	return &UserService{repo: userRepo, relationRepo: relationRepo, audit: audit}
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
