package models

import "backend/internal/domain"

type Product struct {
	ID          int      `json:"id"`
	OwnerID     int      `db:"owner_id"`
	Name        string   `db:"name"`
	PerHour     float64  `db:"per_hour"`
	Description *string  `db:"description"`
	Photos      []string `db:"photos"`
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

type Products []*Product

func (pp Products) Domain() []*domain.Product {
	dd := make([]*domain.Product, 0)
	for _, v := range pp {
		dd = append(dd, v.Domain())
	}

	return dd
}
