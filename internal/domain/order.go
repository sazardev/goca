package domain

type Order struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	CustomerID int     `json:"customer_id" gorm:"type:integer;not null;default:0"`
	Total      float64 `json:"total" gorm:"type:decimal(10,2);not null;default:0"`
	Status     string  `json:"status" gorm:"type:varchar(255)"`
}

func (o *Order) Validate() error {
	if o.CustomerID < 0 {
		return ErrInvalidOrderCustomerID
	}
	if o.Total < 0 {
		return ErrInvalidOrderTotal
	}
	if o.Status == "" {
		return ErrInvalidOrderStatus
	}
	return nil
}
