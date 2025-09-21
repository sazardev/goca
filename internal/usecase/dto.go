package usecase

import (
	"github.com/sazardev/goca/internal/domain"
)

// CreateTestPerfRequest DTO para crear un nuevo testperf
type CreateTestPerfRequest struct {
	Name   string  `json:"name" validate:"required,min=1"`
	Email  string  `json:"email" validate:"required,min=1"`
	Age    int     `json:"age" validate:"required,min=1"`
	Score  float64 `json:"score" validate:"required,min=0"`
	Active bool    `json:"active" validate:""`
}

// Validate valida los datos del DTO CreateTestPerfRequest
func (r *CreateTestPerfRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("el nombre es requerido")
	}
	if r.Email == "" {
		return errors.New("el email es requerido")
	}
	if !strings.Contains(r.Email, "@") {
		return errors.New("formato de el email inválido")
	}
	if r.Age < 0 {
		return errors.New("la edad debe ser un número positivo")
	}
	if r.Score < 0 {
		return errors.New("el campo score debe ser un número positivo")
	}
	return nil
}

// CreateTestPerfResponse DTO para la respuesta de creación
type CreateTestPerfResponse struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Age     int     `json:"age"`
	Score   float64 `json:"score"`
	Active  bool    `json:"active"`
	Message string  `json:"message"`
}

type UpdateTestPerfInput struct {
	Name   *string  `json:"name,omitempty" validate:"omitempty,required,min=1"`
	Email  *string  `json:"email,omitempty" validate:"omitempty,required,min=1"`
	Age    *int     `json:"age,omitempty" validate:"omitempty,required,min=1"`
	Score  *float64 `json:"score,omitempty" validate:"omitempty,required,min=0"`
	Active *bool    `json:"active,omitempty" validate:"omitempty,"`
}

type ListTestPerfOutput struct {
	TestPerfs []domain.TestPerf `json:"testperfs"`
	Total     int               `json:"total"`
	Message   string            `json:"message"`
}
