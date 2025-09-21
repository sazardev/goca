package domain

import (
	"time"
)

// GetTestPerfSeeds retorna datos de ejemplo para testperf
func GetTestPerfSeeds() []TestPerf {
	return []TestPerf{
		{
			Name:   "Juan Pérez",
			Email:  "juan@ejemplo.com",
			Age:    25,
			Score:  10.50,
			Active: true,
		},
		{
			Name:   "María García",
			Email:  "maria@ejemplo.com",
			Age:    30,
			Score:  21.00,
			Active: false,
		},
		{
			Name:   "Carlos López",
			Email:  "carlos@ejemplo.com",
			Age:    35,
			Score:  31.50,
			Active: true,
		},
	}
}

// GetSQLTestPerfSeeds retorna sentencias SQL INSERT para testperf
func GetSQLTestPerfSeeds() string {
	return `-- Datos de ejemplo para tabla testperf
INSERT INTO testperfs (name, email, age, score, active) VALUES ('Juan Pérez', 'juan@ejemplo.com', 25, 10.50, true);\nINSERT INTO testperfs (name, email, age, score, active) VALUES ('María García', 'maria@ejemplo.com', 30, 21.00, false);\nINSERT INTO testperfs (name, email, age, score, active) VALUES ('Carlos López', 'carlos@ejemplo.com', 35, 31.50, true);\n`
}
