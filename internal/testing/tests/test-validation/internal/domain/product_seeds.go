package domain

import (
	"time"
)

// GetProductSeeds retorna datos de ejemplo para product
func GetProductSeeds() []Product {
	return []Product{
		{
			Name:  "Juan Pérez",
			Price: 99.99,
		},
		{
			Name:  "María García",
			Price: 149.50,
		},
		{
			Name:  "Carlos López",
			Price: 199.99,
		},
	}
}

// GetSQLProductSeeds retorna sentencias SQL INSERT para product
func GetSQLProductSeeds() string {
	return `-- Datos de ejemplo para tabla product
INSERT INTO products (name, price) VALUES ('Juan Pérez', 99.99);\nINSERT INTO products (name, price) VALUES ('María García', 149.50);\nINSERT INTO products (name, price) VALUES ('Carlos López', 199.99);\n`
}
