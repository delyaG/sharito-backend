package viewmodels

import (
	"backend/internal/domain"
	"time"
)

type Order struct {
	OrderStart time.Time `json:"order_start"`
	OrderEnd   time.Time `json:"order_end"`
	User       *User     `json:"user"`
	Product    *Product  `json:"product"`
	Price      float64   `json:"price"`
}

func (o *Order) ViewModel(d *domain.Order) {
	o.OrderStart = d.OrderStart
	o.OrderEnd = d.OrderEnd
	o.User = &User{}
	o.User.ViewModel(d.User)
	o.Product = &Product{}
	o.Product.ViewModel(d.Product)
	o.Price = d.Price * d.Product.PerHour
}

type Orders []*Order

func (oo *Orders) ViewModel(dd []*domain.Order) {
	*oo = make([]*Order, 0)
	for _, d := range dd {
		var o Order
		o.ViewModel(d)
		*oo = append(*oo, &o)
	}
}
