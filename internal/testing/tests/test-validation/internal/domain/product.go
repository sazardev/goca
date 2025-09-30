package domain

type Product struct {
	ID    uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string  `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
	Price float64 `json:"price" gorm:"type:decimal(10,2);not null;default:0" validate:"required,gte=0"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrInvalidProductName
	}
	if p.Price < 0 {
		return ErrInvalidProductPrice
	}
	return nil
}
