package viewmodels

import "backend/internal/domain"

type Product struct {
	ID          int      `json:"id"`
	OwnerID     int      `json:"owner_id"`
	Name        string   `json:"name"`
	PerHour     float64  `json:"per_hour"`
	Description *string  `json:"description,omitempty"`
	Photos      []string `json:"photos"`
}

func (p *Product) Domain() *domain.Product {
	return &domain.Product{
		ID:          p.ID,
		OwnerID:     p.OwnerID,
		Name:        p.Name,
		PerHour:     p.PerHour,
		Description: p.Description,
		Photos:      p.Photos,
	}
}

func (p *Product) ViewModel(d *domain.Product) {
	p.ID = d.ID
	p.OwnerID = d.OwnerID
	p.Name = d.Name
	p.PerHour = d.PerHour
	p.Description = d.Description
	p.Photos = d.Photos
}

type Products []*Product

func (oo *Products) ViewModel(dd []*domain.Product) {
	*oo = make([]*Product, 0)
	for _, d := range dd {
		var p Product
		p.ViewModel(d)
		*oo = append(*oo, &p)
	}
}

type ProductsWithCount struct {
	Products Products `json:"products"`
	Count    int      `json:"count"`
}

func (p *ProductsWithCount) ViewModel(pp []*domain.Product, count int) {
	p.Products.ViewModel(pp)
	p.Count = count
}
