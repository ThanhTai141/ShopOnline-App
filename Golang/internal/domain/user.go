package domain

import "fmt"

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}
type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	GetAll() ([]*User, error)
}

func (u *User) ValidateLogin() error {
	if len(u.Email) == 0 {
		return fmt.Errorf("email is required")
	}
	if len(u.Password) == 0 {
		return fmt.Errorf("password is required")
	}
	return nil
}
