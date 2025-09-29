package domain

type Order struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Customer_id int     `json:"customer_id" gorm:"type:integer;not null;default:0"`
	Total       float64 `json:"total" gorm:"type:decimal(10,2);not null;default:0"`
	Status      string  `json:"status" gorm:"type:varchar(255)"`
}

func (o *Order) Validate() error {
	if o.Customer_id < 0 {
		return ErrInvalidOrderCustomer_id
	}
	if o.Total < 0 {
		return ErrInvalidOrderTotal
	}
	if o.Status == "" {
		return ErrInvalidOrderStatus
	}
	return nil
}
