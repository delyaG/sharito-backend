package domain

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Service interface {
	UserService
	AuthService
	ProductService
}

type AuthService interface {
	Register(user *User) (string, error)
	Login(user *User) (string, error)
}

type UserService interface {
	GetUser(ctx context.Context) (*User, error)
}

type ProductService interface {
	AddProduct(ctx context.Context, product *Product) (int, error)
	GetProductAndOwnerUserByProductID(productID int) (*Product, *User, error)
	GetProductsWithPagination(page, count int, search string) ([]*Product, int, error)
	RentProduct(ctx context.Context, productID int, from, to time.Time) error
	GetOrders(ctx context.Context, isMine bool) ([]*Order, error)
}

type service struct {
	logger   logrus.FieldLogger
	db       Database
	security Security
}

func NewService(logger logrus.FieldLogger, db Database, security Security) Service {
	s := &service{
		logger:   logger,
		db:       db,
		security: security,
	}

	return s
}

func (s *service) Register(user *User) (string, error) {
	// get hash password
	passwordHash, salt, err := s.security.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.PasswordHash = passwordHash
	user.Salt = salt

	userID, err := s.db.SaveUser(user)
	if err != nil {
		return "", err
	}

	token, err := s.security.GenerateNewJWT(userID, 24*time.Hour)
	if err != nil {
		return "", nil
	}

	return token, nil
}

func (s *service) Login(u *User) (string, error) {
	user, err := s.db.GetUserByLogin(u.Login)
	if err != nil {
		return "", err
	}

	if !s.security.VerifyPassword(user.Salt, user.PasswordHash, u.Password) {
		return "", ErrInvalidInputData
	}

	token, err := s.security.GenerateNewJWT(user.ID, 24*time.Hour)
	if err != nil {
		return "", nil
	}

	return token, nil
}

func (s *service) GetUser(ctx context.Context) (*User, error) {
	userID := ctx.Value(ContextUserID).(int)
	return s.db.GetUserByID(userID)
}

func (s *service) AddProduct(ctx context.Context, product *Product) (int, error) {
	// get user_id like owner
	userID := ctx.Value(ContextUserID).(int)
	product.OwnerID = userID

	return s.db.SaveProduct(product)
}

func (s *service) GetProductAndOwnerUserByProductID(productID int) (*Product, *User, error) {
	product, err := s.db.GetProductByID(productID)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.db.GetUserByID(product.OwnerID)
	if err != nil {
		return nil, nil, err
	}

	return product, user, nil
}

func (s *service) GetProductsWithPagination(page, count int, search string) ([]*Product, int, error) {
	return s.db.GetProductsWithPagination(count, page*count, search)
}

func (s *service) GetOrders(ctx context.Context, isMine bool) ([]*Order, error) {
	userID := ctx.Value(ContextUserID).(int)

	orders, err := s.db.GetOrders(userID, isMine)
	if err != nil {
		return nil, err
	}

	for i, v := range orders {
		product, err := s.db.GetProductByID(v.ProductID)
		if err != nil {
			return nil, err
		}
		orders[i].Product = product

		user, err := s.db.GetUserByID(v.UserID)
		if err != nil {
			return nil, err
		}
		orders[i].User = user
	}

	return orders, nil
}

func (s *service) RentProduct(ctx context.Context, productID int, from, to time.Time) error {
	userID := ctx.Value(ContextUserID).(int)
	return s.db.RentProduct(productID, userID, from, to)
}
