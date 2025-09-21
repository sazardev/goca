package domain

type TestPerf struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	ID     uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string  `json:"name" gorm:"type:varchar(255);not null"`
	Email  string  `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Age    int     `json:"age" gorm:"type:integer;not null;default:0"`
	Score  float64 `json:"score" gorm:"type:decimal(10,2);not null;default:0"`
	Active bool    `json:"active" gorm:"type:boolean;not null;default:false"`
}

func (t *TestPerf) Validate() error {
	if t.Name == "" {
		return ErrInvalidTestPerfName
	}
	if t.Email == "" {
		return ErrInvalidTestPerfEmail
	}
	if t.Age < 0 {
		return ErrInvalidTestPerfAge
	}
	if t.Score < 0 {
		return ErrInvalidTestPerfScore
	}
	return nil
}
