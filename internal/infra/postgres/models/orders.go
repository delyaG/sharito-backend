package models

import (
	"backend/internal/domain"
	"time"
)

type Order struct {
	ID         int       `db:"id"`
	OrderStart time.Time `db:"order_start"`
	OrderEnd   time.Time `db:"order_end"`
	UserID     int       `db:"user_id"`
	ProductID  int       `db:"product_id"`
	Price      float64   `db:"price"`
}

func (o *Order) Domain() *domain.Order {
	return &domain.Order{
		ID:         o.ID,
		OrderStart: o.OrderStart,
		OrderEnd:   o.OrderEnd,
		UserID:     o.UserID,
		ProductID:  o.ProductID,
		Price:      o.Price,
	}
}

type Orders []*Order

func (oo Orders) Domain() []*domain.Order {
	dd := make([]*domain.Order, 0)
	for _, v := range oo {
		dd = append(dd, v.Domain())
	}

	return dd
}
