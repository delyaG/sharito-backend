package models

import "backend/internal/domain"

type User struct {
	ID           int    `db:"id"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	Login        string `db:"login"`
	Email        string `db:"email"`
	PasswordHash []byte `db:"password_hash"`
	Salt         []byte `db:"salt"`
}

func (u *User) Domain() *domain.User {
	return &domain.User{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Login:        u.Login,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Salt:         u.Salt,
	}
}
