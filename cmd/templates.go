package cmd

// Entity template for dynamic generation
const entityTemplate = `package domain

{{range .Imports}}import "{{.}}"
{{end}}

type {{.Entity.Name}} struct {
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`" + `{{.JSONTag}} {{.GormTag}}` + "`" + `
{{end}}{{if .Features.Timestamps}}	CreatedAt time.Time ` + "`" + `json:"created_at" gorm:"autoCreateTime"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at" gorm:"autoUpdateTime"` + "`" + `
{{end}}{{if .Features.SoftDelete}}	DeletedAt gorm.DeletedAt ` + "`" + `json:"deleted_at,omitempty" gorm:"index"` + "`" + `
{{end}}}

{{if .Features.Validation}}
func ({{.Entity.NameLower}} *{{.Entity.Name}}) Validate() error {
{{range .Validations}}	if {{.Field}} == "" {
		return errors.New("{{.Message}}")
	}
{{end}}	return nil
}
{{end}}

{{range .Methods}}
func ({{$.Entity.NameLower}} *{{$.Entity.Name}}) {{.Name}}({{range .Params}}{{.Name}} {{.Type}}, {{end}}) {{.ReturnType}} {
	{{.Body}}
}
{{end}}
`

// UseCase template for dynamic generation
const useCaseTemplate = `package usecase

import (
	"{{.Module}}/internal/domain"
	"{{.Module}}/internal/repository"
)

type {{.Entity.Name}}UseCase interface {
	Create{{.Entity.Name}}(input Create{{.Entity.Name}}Input) (*Create{{.Entity.Name}}Output, error)
	Get{{.Entity.Name}}ByID(id int) (*{{.Entity.Name}}Output, error)
	Update{{.Entity.Name}}(id int, input Update{{.Entity.Name}}Input) error
	Delete{{.Entity.Name}}(id int) error
	List{{.Entity.Name}}s() (*List{{.Entity.Name}}sOutput, error)
}

type {{.Entity.NameLower}}Service struct {
	repo repository.{{.Entity.Name}}Repository
}

func New{{.Entity.Name}}Service(repo repository.{{.Entity.Name}}Repository) {{.Entity.Name}}UseCase {
	return &{{.Entity.NameLower}}Service{repo: repo}
}

type Create{{.Entity.Name}}Input struct {
{{range .Fields}}{{if ne .Name "ID"}}	{{.Name}} {{.Type}} ` + "`" + `json:"{{.JSONTag}}"{{if .ValidateTag}} validate:"{{.ValidateTag}}"{{end}}` + "`" + `
{{end}}{{end}}}

type Create{{.Entity.Name}}Output struct {
	{{.Entity.Name}} *domain.{{.Entity.Name}} ` + "`" + `json:"{{.Entity.NameLower}}"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
}

type Update{{.Entity.Name}}Input struct {
{{range .Fields}}{{if ne .Name "ID"}}	{{.Name}} *{{.Type}} ` + "`" + `json:"{{.JSONTag}},omitempty"` + "`" + `
{{end}}{{end}}}

type {{.Entity.Name}}Output struct {
	{{.Entity.Name}} *domain.{{.Entity.Name}} ` + "`" + `json:"{{.Entity.NameLower}}"` + "`" + `
}

type List{{.Entity.Name}}sOutput struct {
	{{.Entity.NamePlural}} []domain.{{.Entity.Name}} ` + "`" + `json:"{{.Entity.NameLower}}s"` + "`" + `
	Total int ` + "`" + `json:"total"` + "`" + `
}
`

// Repository template for dynamic generation
const repositoryTemplate = `package repository

import "{{.Module}}/internal/domain"

type {{.Entity.Name}}Repository interface {
	Save({{.Entity.NameLower}} *domain.{{.Entity.Name}}) error
	FindByID(id int) (*domain.{{.Entity.Name}}, error)
{{range .Fields}}{{if .IsSearchable}}	FindBy{{.Name}}({{.Name}} {{.Type}}) (*domain.{{.Entity.Name}}, error)
{{end}}{{end}}	Update({{.Entity.NameLower}} *domain.{{.Entity.Name}}) error
	Delete(id int) error
	FindAll() ([]domain.{{.Entity.Name}}, error)
{{if .Features.Transactions}}	SaveWithTx(tx interface{}, {{.Entity.NameLower}} *domain.{{.Entity.Name}}) error
	UpdateWithTx(tx interface{}, {{.Entity.NameLower}} *domain.{{.Entity.Name}}) error
	DeleteWithTx(tx interface{}, id int) error
{{end}}}
`

// Handler template for dynamic generation
const handlerTemplate = `package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	
	"github.com/gorilla/mux"
	"{{.Module}}/internal/usecase"
)

type {{.Entity.Name}}Handler struct {
	usecase usecase.{{.Entity.Name}}UseCase
}

func New{{.Entity.Name}}Handler(uc usecase.{{.Entity.Name}}UseCase) *{{.Entity.Name}}Handler {
	return &{{.Entity.Name}}Handler{usecase: uc}
}

func (h *{{.Entity.Name}}Handler) Create{{.Entity.Name}}(w http.ResponseWriter, r *http.Request) {
	var input usecase.Create{{.Entity.Name}}Input
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	output, err := h.usecase.Create{{.Entity.Name}}(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *{{.Entity.Name}}Handler) Get{{.Entity.Name}}(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid {{.Entity.NameLower}} ID", http.StatusBadRequest)
		return
	}
	
	output, err := h.usecase.Get{{.Entity.Name}}ByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
`
