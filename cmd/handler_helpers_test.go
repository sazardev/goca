package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCreateHandlerMethod(t *testing.T) {
	t.Parallel()

	t.Run("without validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateCreateHandlerMethod(&b, "Product", "ProductHandler", false, false)
		output := b.String()
		assert.Contains(t, output, "func (p *ProductHandler) CreateProduct(")
		assert.Contains(t, output, "CreateProductInput")
		assert.Contains(t, output, "json.NewDecoder")
		assert.Contains(t, output, "StatusCreated")
		assert.NotContains(t, output, "validator.New()")
	})

	t.Run("with validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateCreateHandlerMethod(&b, "Product", "ProductHandler", true, false)
		output := b.String()
		assert.Contains(t, output, "validator.New().Struct(input)")
		assert.Contains(t, output, "StatusUnprocessableEntity")
	})
}

// Regression test: an entity name starting with "W" or "R" (Widget, Wallet,
// Report, ...) must not collide the receiver variable with the fixed
// "w http.ResponseWriter"/"r *http.Request" parameter names.
func TestHandlerReceiverVar_AvoidsParamCollisions(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "h", handlerReceiverVar("WidgetHandler"))
	assert.Equal(t, "h", handlerReceiverVar("ReportHandler"))
	assert.Equal(t, "p", handlerReceiverVar("ProductHandler"))

	var b strings.Builder
	generateCreateHandlerMethod(&b, "Widget", "WidgetHandler", false, false)
	output := b.String()
	assert.Contains(t, output, "func (h *WidgetHandler) CreateWidget(w http.ResponseWriter, r *http.Request)")
	assert.NotContains(t, output, "func (w *WidgetHandler)")
}

func TestGenerateGetHandlerMethod(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	generateGetHandlerMethod(&b, "Product", "ProductHandler", false)
	output := b.String()
	assert.Contains(t, output, "func (p *ProductHandler) GetProduct(")
	assert.Contains(t, output, "mux.Vars(r)")
	assert.Contains(t, output, "strconv.Atoi")
	assert.Contains(t, output, "Invalid product ID")
	assert.Contains(t, output, "StatusNotFound")
}

func TestGenerateUpdateHandlerMethod(t *testing.T) {
	t.Parallel()

	t.Run("without validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateUpdateHandlerMethod(&b, "Product", "ProductHandler", false, false)
		output := b.String()
		assert.Contains(t, output, "func (p *ProductHandler) UpdateProduct(")
		assert.Contains(t, output, "UpdateProductInput")
		assert.Contains(t, output, "StatusNoContent")
		assert.NotContains(t, output, "validator.New()")
	})

	t.Run("with validation", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateUpdateHandlerMethod(&b, "Product", "ProductHandler", true, false)
		output := b.String()
		assert.Contains(t, output, "validator.New().Struct(input)")
	})
}

func TestGenerateDeleteHandlerMethod(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	generateDeleteHandlerMethod(&b, "Product", "ProductHandler", false)
	output := b.String()
	assert.Contains(t, output, "func (p *ProductHandler) DeleteProduct(")
	assert.Contains(t, output, "mux.Vars(r)")
	assert.Contains(t, output, "StatusNoContent")
}

func TestGenerateListHandlerMethod(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	generateListHandlerMethod(&b, "Product", "ProductHandler", false)
	output := b.String()
	assert.Contains(t, output, "func (p *ProductHandler) ListProducts(")
	assert.Contains(t, output, "ListProducts()")
	assert.Contains(t, output, "json.NewEncoder(w).Encode(output)")
}

// HANDLER-5: --swagger must add @Summary/@Router/@Success godoc annotations.
func TestHandlerMethods_SwaggerAnnotations(t *testing.T) {
	t.Parallel()

	t.Run("create has annotations", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateCreateHandlerMethod(&b, "Product", "ProductHandler", false, true)
		out := b.String()
		assert.Contains(t, out, "@Summary Create product")
		assert.Contains(t, out, "@Router /products [post]")
		assert.Contains(t, out, "@Success 201 {object} usecase.CreateProductOutput")
		assert.Contains(t, out, "@Param body body usecase.CreateProductInput")
	})

	t.Run("get has path param and router", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateGetHandlerMethod(&b, "Product", "ProductHandler", true)
		out := b.String()
		assert.Contains(t, out, "@Router /products/{id} [get]")
		assert.Contains(t, out, "@Param id path int true")
	})

	t.Run("no annotations when swagger disabled", func(t *testing.T) {
		t.Parallel()
		var b strings.Builder
		generateCreateHandlerMethod(&b, "Product", "ProductHandler", false, false)
		assert.NotContains(t, b.String(), "@Summary")
	})
}

// HANDLER-3: protoType maps Go scalar types to proto3 types.
func TestProtoType(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"string": "string", "int": "int32", "int64": "int64",
		"bool": "bool", "float64": "double", "time.Time": "",
	}
	for in, want := range cases {
		assert.Equal(t, want, protoType(in), "protoType(%q)", in)
	}
}
