package viewmodels

import "backend/internal/domain"

type ProductWithUser struct {
	Product Product `json:"product"`
	User    User    `json:"user"`
}

func (p *ProductWithUser) ViewModel(dp *domain.Product, du *domain.User) {
	p.Product.ViewModel(dp)
	p.User.ViewModel(du)
}
