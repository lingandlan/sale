package rbac

import (
	"fmt"
	"strings"
	"sync"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"marketplace/backend/internal/config"
)

// CasbinService casbin 权限服务
type CasbinService struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
	mu       sync.RWMutex
}

// NewCasbinService 创建 casbin 服务
func NewCasbinService(sqlDB *sqlx.DB, cfg *config.DatabaseConfig) (*CasbinService, error) {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB.DB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("create gorm db for casbin failed: %w", err)
	}

	adapter, err := gormadapter.NewAdapterByDBUseTableName(gormDB, "casbin", "")
	if err != nil {
		return nil, fmt.Errorf("create casbin adapter failed: %w", err)
	}

	e, err := casbin.NewEnforcer("configs/rbac_model.conf", adapter)
	if err != nil {
		return nil, fmt.Errorf("create casbin enforcer failed: %w", err)
	}

	if err := e.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("load casbin policy failed: %w", err)
	}

	svc := &CasbinService{
		enforcer: e,
		db:       gormDB,
	}

	// 如果没有策略，初始化基础权限
	policy, _ := e.GetPolicy()
	if len(policy) == 0 {
		if err := svc.initDefaultPolicies(); err != nil {
			return nil, fmt.Errorf("init default policies failed: %w", err)
		}
	}

	return svc, nil
}

// Enforce 检查权限
func (s *CasbinService) Enforce(sub interface{}, obj string, act string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.enforcer.Enforce(sub, obj, act)
}

// AddPolicy 添加策略
func (s *CasbinService) AddPolicy(sub, obj, act string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		return fmt.Errorf("add policy failed: %w", err)
	}
	return s.enforcer.SavePolicy()
}

// RemovePolicy 移除策略
func (s *CasbinService) RemovePolicy(sub, obj, act string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		return fmt.Errorf("remove policy failed: %w", err)
	}
	return s.enforcer.SavePolicy()
}

// AddRoleForUser 为用户添加角色
func (s *CasbinService) AddRoleForUser(userID, role string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.enforcer.AddGroupingPolicy(userID, role)
	if err != nil {
		return fmt.Errorf("add role for user failed: %w", err)
	}
	return s.enforcer.SavePolicy()
}

// RemoveRoleForUser 移除用户的角色
func (s *CasbinService) RemoveRoleForUser(userID, role string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.enforcer.RemoveGroupingPolicy(userID, role)
	if err != nil {
		return fmt.Errorf("remove role for user failed: %w", err)
	}
	return s.enforcer.SavePolicy()
}

// GetRolesForUser 获取用户的所有角色
func (s *CasbinService) GetRolesForUser(userID string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.enforcer.GetRolesForUser(userID)
}

// GetPermissionsForUser 获取用户的权限
func (s *CasbinService) GetPermissionsForUser(userID string) ([][]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.enforcer.GetPermissionsForUser(userID)
}

// ClearCasbinForUser 清除用户的所有角色和权限
func (s *CasbinService) ClearCasbinForUser(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.enforcer.RemoveFilteredGroupingPolicy(0, userID)
	if err != nil {
		return fmt.Errorf("clear user roles failed: %w", err)
	}
	return s.enforcer.SavePolicy()
}

// UpdateUserRole 更新用户角色
func (s *CasbinService) UpdateUserRole(userID, newRole string) error {
	// 清除旧角色
	if err := s.ClearCasbinForUser(userID); err != nil {
		return err
	}
	// 添加新角色
	return s.AddRoleForUser(userID, newRole)
}

// EnforcePath 简化版检查（支持路径通配符）
func (s *CasbinService) EnforcePath(role, method, path string) (bool, error) {
	// 将路径中的数字 ID 替换为通配符，兼容 /api/v1/users/123 -> /api/v1/users/*
	pattern := convertPathToPattern(path)

	s.mu.RLock()
	defer s.mu.RUnlock()

	// 先尝试精确匹配
	ok, err := s.enforcer.Enforce(role, pattern, method)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	// 尝试通配符匹配
	ok, err = s.enforcer.Enforce(role, convertPathToWildcard(pattern), method)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// convertPathToPattern 将路径中的 ID 替换为通配符
func convertPathToPattern(path string) string {
	// 匹配 /api/v1/users/123 或 /api/v1/users/abc-123 这种模式
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if isDynamicSegment(part) {
			parts[i] = "*"
		}
	}
	return strings.Join(parts, "/")
}

// convertPathToWildcard 将路径转为通配符模式
func convertPathToWildcard(path string) string {
	parts := strings.Split(path, "/")
	for i := range parts {
		if parts[i] == "*" {
			parts[i] = "*"
		}
	}
	// 添加通配符路径
	return strings.Join(parts, "/")
}

// isDynamicSegment 判断是否为动态路径段（纯数字或 UUID 等 ID 格式）
func isDynamicSegment(segment string) bool {
	if segment == "" {
		return false
	}
	// 纯数字（如 123, 456）
	allDigits := true
	for _, c := range segment {
		if c < '0' || c > '9' {
			allDigits = false
			break
		}
	}
	if allDigits && len(segment) > 0 {
		return true
	}
	// UUID 格式（包含连字符，如 abc-123-def）
	hasLetter := false
	hasDash := false
	for _, c := range segment {
		if c == '-' {
			hasDash = true
		} else if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		} else if c < '0' || c > '9' {
			return false
		}
	}
	// 只有包含连字符且有字母的才视为动态 ID（如 center-bj-cy）
	return hasDash && hasLetter
}

// ReloadPolicy 重新加载策略
func (s *CasbinService) ReloadPolicy() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.enforcer.LoadPolicy()
}

// initDefaultPolicies 初始化默认角色权限策略
func (s *CasbinService) initDefaultPolicies() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// center_admin 和 operator 可访问的基础路径
	policies := [][]string{
		// 用户自身信息
		{"center_admin", "/api/v1/user/*", "*"},
		{"operator", "/api/v1/user/*", "*"},
		// Dashboard
		{"center_admin", "/api/v1/dashboard/*", "GET"},
		{"operator", "/api/v1/dashboard/*", "GET"},
		// 充值相关
		{"center_admin", "/api/v1/recharge/*", "*"},
		{"operator", "/api/v1/recharge/*", "*"},
		// 门店卡
		{"center_admin", "/api/v1/card/*", "*"},
		{"operator", "/api/v1/card/*", "*"},
		// 充值中心（只读）
		{"center_admin", "/api/v1/center", "GET"},
		{"center_admin", "/api/v1/center/*", "GET"},
		{"operator", "/api/v1/center", "GET"},
		{"operator", "/api/v1/center/*", "GET"},
		// 操作员（只读）
		{"center_admin", "/api/v1/operator", "GET"},
		{"operator", "/api/v1/operator", "GET"},
		// 管理员接口（center_admin 部分权限）
		{"center_admin", "/api/v1/admin/users", "*"},
		{"center_admin", "/api/v1/admin/users/*", "*"},
	}

	for _, p := range policies {
		if _, err := s.enforcer.AddPolicy(p); err != nil {
			return err
		}
	}

	return s.enforcer.SavePolicy()
}
