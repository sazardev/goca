package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DataGenerator generates realistic sample data for testing
type DataGenerator struct {
	rand *rand.Rand
}

// NewDataGenerator creates a new data generator
func NewDataGenerator() *DataGenerator {
	return &DataGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateSampleData generates sample data based on field types
func (g *DataGenerator) GenerateSampleData(fields []Field, entity string) map[string]interface{} {
	data := make(map[string]interface{})
	entityLower := strings.ToLower(entity)

	for _, field := range fields {
		if field.Name == "ID" {
			continue // Skip ID as it's auto-generated
		}

		data[field.Name] = g.generateValueForField(field, entityLower)
	}

	return data
}

// generateValueForField generates a realistic value for a specific field
func (g *DataGenerator) generateValueForField(field Field, entity string) interface{} {
	fieldLower := strings.ToLower(field.Name)

	// Context-aware generation based on field name and entity
	switch {
	case fieldLower == "email":
		return g.generateEmail(entity)
	case fieldLower == "name" || fieldLower == "title":
		return g.generateName(entity, fieldLower)
	case fieldLower == "age":
		return g.rand.Intn(60) + 18 // 18-77 years
	case fieldLower == "price":
		return float64(g.rand.Intn(10000)) + g.rand.Float64()*100
	case fieldLower == "description":
		return g.generateDescription(entity)
	case fieldLower == "status":
		return g.generateStatus(entity)
	case fieldLower == "phone":
		return g.generatePhone()
	case fieldLower == "address":
		return g.generateAddress()
	case strings.Contains(fieldLower, "url") || strings.Contains(fieldLower, "website"):
		return g.generateURL(entity)
	case strings.Contains(fieldLower, "code") || strings.Contains(fieldLower, "sku"):
		return g.generateCode(entity)
	default:
		return g.generateByType(field.Type)
	}
}

// generateEmail generates realistic email addresses
func (g *DataGenerator) generateEmail(entity string) string {
	names := []string{"john", "maria", "carlos", "ana", "miguel", "lucia", "david", "sofia"}
	domains := []string{"gmail.com", "hotmail.com", "company.com", "test.org"}

	name := names[g.rand.Intn(len(names))]
	domain := domains[g.rand.Intn(len(domains))]

	return fmt.Sprintf("%s.%s@%s", name, entity, domain)
}

// generateName generates contextual names
func (g *DataGenerator) generateName(entity, field string) string {
	switch entity {
	case "user", "employee", "customer":
		names := []string{"Juan Pérez", "María García", "Carlos López", "Ana Martín", "Luis Rodríguez"}
		return names[g.rand.Intn(len(names))]
	case "product":
		if field == "name" {
			products := []string{"Laptop Pro", "Smartphone X", "Tablet Ultra", "Auriculares Premium", "Monitor 4K"}
			return products[g.rand.Intn(len(products))]
		}
	case "project":
		projects := []string{"Sistema de Gestión", "App Móvil", "Portal Web", "API Rest", "Dashboard Analytics"}
		return projects[g.rand.Intn(len(projects))]
	case "order":
		return fmt.Sprintf("Pedido #%d", g.rand.Intn(10000))
	}

	// Use cases.Title for proper word capitalization
	caser := cases.Title(language.English)
	return fmt.Sprintf("Sample %s %d", caser.String(entity), g.rand.Intn(1000))
}

// generateDescription generates context-aware descriptions
func (g *DataGenerator) generateDescription(entity string) string {
	switch entity {
	case "product":
		descriptions := []string{
			"Producto de alta calidad con excelentes características",
			"Diseño innovador y funcionalidad superior",
			"Perfecto para uso profesional y personal",
			"Tecnología avanzada y fácil de usar",
		}
		return descriptions[g.rand.Intn(len(descriptions))]
	case "project":
		descriptions := []string{
			"Proyecto estratégico para mejorar la eficiencia operacional",
			"Iniciativa de transformación digital",
			"Desarrollo de nueva funcionalidad core",
			"Optimización de procesos existentes",
		}
		return descriptions[g.rand.Intn(len(descriptions))]
	default:
		return fmt.Sprintf("Descripción detallada del %s con información relevante", entity)
	}
}

// generateStatus generates appropriate status values
func (g *DataGenerator) generateStatus(entity string) string {
	switch entity {
	case "user", "employee":
		statuses := []string{"active", "inactive", "pending"}
		return statuses[g.rand.Intn(len(statuses))]
	case "order":
		statuses := []string{"pending", "processing", "shipped", "delivered", "cancelled"}
		return statuses[g.rand.Intn(len(statuses))]
	case "project":
		statuses := []string{"planning", "in_progress", "testing", "completed", "on_hold"}
		return statuses[g.rand.Intn(len(statuses))]
	default:
		statuses := []string{"active", "inactive", "draft", "published"}
		return statuses[g.rand.Intn(len(statuses))]
	}
}

// generatePhone generates realistic phone numbers
func (g *DataGenerator) generatePhone() string {
	return fmt.Sprintf("+34 %d%d%d %d%d%d %d%d%d",
		g.rand.Intn(10), g.rand.Intn(10), g.rand.Intn(10),
		g.rand.Intn(10), g.rand.Intn(10), g.rand.Intn(10),
		g.rand.Intn(10), g.rand.Intn(10), g.rand.Intn(10))
}

// generateAddress generates realistic addresses
func (g *DataGenerator) generateAddress() string {
	streets := []string{"Calle Mayor", "Avenida Principal", "Plaza Central", "Paseo del Parque"}
	street := streets[g.rand.Intn(len(streets))]
	number := g.rand.Intn(200) + 1

	return fmt.Sprintf("%s %d, Madrid, España", street, number)
}

// generateURL generates realistic URLs
func (g *DataGenerator) generateURL(entity string) string {
	domains := []string{"example.com", "company.org", "business.net"}
	domain := domains[g.rand.Intn(len(domains))]

	return fmt.Sprintf("https://www.%s/%s", domain, strings.ToLower(entity))
}

// generateCode generates realistic codes/SKUs
func (g *DataGenerator) generateCode(entity string) string {
	prefix := strings.ToUpper(entity[:min(3, len(entity))])
	number := g.rand.Intn(10000)

	return fmt.Sprintf("%s-%04d", prefix, number)
}

// generateByType generates values based on Go types
func (g *DataGenerator) generateByType(fieldType string) interface{} {
	switch fieldType {
	case "string":
		return "Sample text"
	case "int", "int32", "int64":
		return g.rand.Intn(1000)
	case "uint", "uint32", "uint64":
		return uint(g.rand.Intn(1000))
	case "float32", "float64":
		return g.rand.Float64() * 1000
	case "bool":
		return g.rand.Intn(2) == 1
	case "time.Time":
		return time.Now().Add(-time.Duration(g.rand.Intn(365*24)) * time.Hour)
	case "[]byte":
		return []byte("sample data")
	default:
		return "default value"
	}
}

// GenerateTestData generates a complete test dataset
func (g *DataGenerator) GenerateTestData(entity string, fields []Field, count int) []map[string]interface{} {
	var dataset []map[string]interface{}

	for i := 0; i < count; i++ {
		data := g.generateSampleData(fields, entity)
		dataset = append(dataset, data)
	}

	return dataset
}

// generateSampleData is a wrapper for GenerateSampleData to maintain consistency
func (g *DataGenerator) generateSampleData(fields []Field, entity string) map[string]interface{} {
	return g.GenerateSampleData(fields, entity)
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GenerateInsertSQL generates SQL INSERT statements with sample data
func (g *DataGenerator) GenerateInsertSQL(tableName string, fields []Field, count int) string {
	var result strings.Builder
	entity := strings.TrimSuffix(tableName, "s") // Simple pluralization removal

	// Generate field names (excluding ID)
	var fieldNames []string
	for _, field := range fields {
		if field.Name != "ID" {
			fieldNames = append(fieldNames, strings.ToLower(field.Name))
		}
	}

	result.WriteString(fmt.Sprintf("-- Sample data for %s\n", tableName))

	for i := 0; i < count; i++ {
		data := g.GenerateSampleData(fields, entity)

		result.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES (",
			tableName, strings.Join(fieldNames, ", ")))

		var values []string
		caser := cases.Title(language.English)
		for _, fieldName := range fieldNames {
			value := data[caser.String(fieldName)]
			values = append(values, g.formatSQLValue(value))
		}

		result.WriteString(strings.Join(values, ", "))
		result.WriteString(");\n")
	}

	return result.String()
}

// formatSQLValue formats a value for SQL insertion
func (g *DataGenerator) formatSQLValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case int, int32, int64, uint, uint32, uint64:
		return fmt.Sprintf("%v", v)
	case float32, float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	default:
		return fmt.Sprintf("'%v'", v)
	}
}
