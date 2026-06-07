package domain

type Order struct {
	ID         uint    `gorm:"primaryKey"                            json:"id"`
	CustomerID int     `gorm:"type:integer;not null;default:0"       json:"customer_id"`
	Total      float64 `gorm:"type:decimal(10,2);not null;default:0" json:"total"`
	Status     string  `gorm:"type:varchar(255)"                     json:"status"`
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
