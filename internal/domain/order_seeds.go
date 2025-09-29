package domain

import (
	"time"
)

// GetOrderSeeds retorna datos de ejemplo para order
func GetOrderSeeds() []Order {
	return []Order{
		{
			Customer_id: 10,
			Total:       10.50,
			Status:      "activo",
		},
		{
			Customer_id: 20,
			Total:       21.00,
			Status:      "pendiente",
		},
		{
			Customer_id: 30,
			Total:       31.50,
			Status:      "completado",
		},
	}
}

// GetSQLOrderSeeds retorna sentencias SQL INSERT para order
func GetSQLOrderSeeds() string {
	return `-- Datos de ejemplo para tabla order
INSERT INTO orders (customer_id, total, status) VALUES (10, 10.50, 'activo');\nINSERT INTO orders (customer_id, total, status) VALUES (20, 21.00, 'pendiente');\nINSERT INTO orders (customer_id, total, status) VALUES (30, 31.50, 'completado');\n`
}
