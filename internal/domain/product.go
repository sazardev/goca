package domain

type Product struct {
	ID          uint    `gorm:"primaryKey"                            json:"id"`
	Name        string  `gorm:"type:varchar(255);not null"            json:"name"`
	Price       float64 `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	Description string  `gorm:"type:text"                             json:"description"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrInvalidProductName
	}
	if p.Price < 0 {
		return ErrInvalidProductPrice
	}
	if p.Description == "" {
		return ErrInvalidProductDescription
	}
	return nil
}
