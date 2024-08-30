package hashutil

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordManager chứa các hàm để làm việc với mật khẩu
type passwordManager struct {
	cost int
}

// NewPasswordManager tạo một PasswordManager mới với cost chỉ định
func NewPasswordManager(cost int) *passwordManager {
	return &passwordManager{cost: cost}
}

// HashPassword tạo hash cho mật khẩu
func (pm *passwordManager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword so sánh mật khẩu với hash đã lưu
func (pm *passwordManager) VerifyPassword(hash, passwordLogin string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordLogin))
	return err == nil
}
