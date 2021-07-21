package viewmodels

import "backend/internal/domain"

type User struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Login     string  `json:"login"`
	Email     string  `json:"email"`
	Password  *string `json:"password,omitempty"`
}

func (u *User) Domain() *domain.User {
	return &domain.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Login:     u.Login,
		Email:     u.Email,
	}
}

func (u *User) ViewModel(d *domain.User) {
	u.FirstName = d.FirstName
	u.LastName = d.LastName
	u.Login = d.Login
	u.Email = d.Email
}
