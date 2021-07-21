package domain

import "time"

type ContextKey string

const ContextUserID ContextKey = "ctx_user_id"

type User struct {
	ID           int
	FirstName    string
	LastName     string
	Login        string
	Email        string
	Password     string
	PasswordHash []byte
	Salt         []byte
}

type Product struct {
	ID          int
	OwnerID     int
	Name        string
	PerHour     float64
	Description *string
	Photos      []string
}

type Order struct {
	ID         int
	OrderStart time.Time
	OrderEnd   time.Time
	UserID     int
	ProductID  int
	Price      float64
	User       *User
	Product    *Product
}
