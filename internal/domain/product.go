package domain

type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string  `json:"name" gorm:"type:varchar(255);not null"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2);not null;default:0"`
	Description string  `json:"description" gorm:"type:text"`
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
