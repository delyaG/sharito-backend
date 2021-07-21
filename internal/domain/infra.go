package domain

import (
	"context"
	"time"
)

type Delivery interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type Database interface {
	UserRepository
	ProductRepository
	OrderRepository
}

type UserRepository interface {
	SaveUser(user *User) (int, error)
	GetUserByLogin(login string) (*User, error)
	GetUserByID(id int) (*User, error)
}

type ProductRepository interface {
	SaveProduct(product *Product) (int, error)
	GetProductByID(id int) (*Product, error)
	GetProductsWithPagination(limit, offset int, search string) ([]*Product, int, error)
	RentProduct(productID, userID int, from, to time.Time) error
}

type OrderRepository interface {
	GetOrders(userID int, isMine bool) ([]*Order, error)
}

type Security interface {
	HashPassword(password string) ([]byte, []byte, error)
	VerifyPassword(salt []byte, passwordHash []byte, password string) bool
	GenerateNewJWT(id int, duration time.Duration) (string, error)
}
