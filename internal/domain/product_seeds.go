package domain

import (
	"time"
)

// GetProductSeeds retorna datos de ejemplo para product
func GetProductSeeds() []Product {
	return []Product{
		{
			Name:        "Juan Pérez",
			Price:       99.99,
			Description: "Descripción detallada del primer elemento",
		},
		{
			Name:        "María García",
			Price:       149.50,
			Description: "Información completa del segundo item",
		},
		{
			Name:        "Carlos López",
			Price:       199.99,
			Description: "Detalles específicos del tercer registro",
		},
	}
}

// GetSQLProductSeeds retorna sentencias SQL INSERT para product
func GetSQLProductSeeds() string {
	return `-- Datos de ejemplo para tabla product
INSERT INTO products (name, price, description) VALUES ('Juan Pérez', 99.99, 'Descripción detallada del primer elemento');\nINSERT INTO products (name, price, description) VALUES ('María García', 149.50, 'Información completa del segundo item');\nINSERT INTO products (name, price, description) VALUES ('Carlos López', 199.99, 'Detalles específicos del tercer registro');\n`
}
