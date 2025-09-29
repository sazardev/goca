package domain

// GetTestFeatureSeeds retorna datos de ejemplo para testfeature
func GetTestFeatureSeeds() []TestFeature {
	return []TestFeature{
		{
			Name:  "Juan Pérez",
			Email: "Ejemplo Email 1",
			Age:   25,
		},
		{
			Name:  "María García",
			Email: "Ejemplo Email 2",
			Age:   30,
		},
		{
			Name:  "Carlos López",
			Email: "Ejemplo Email 3",
			Age:   35,
		},
	}
}

// GetSQLTestFeatureSeeds retorna sentencias SQL INSERT para testfeature
func GetSQLTestFeatureSeeds() string {
	return `-- Datos de ejemplo para tabla testfeature
INSERT INTO testfeatures (name, email, age) VALUES ('Juan Pérez', 'juan@ejemplo.com', 25);\nINSERT INTO testfeatures (name, email, age) VALUES ('María García', 'maria@ejemplo.com', 30);\nINSERT INTO testfeatures (name, email, age) VALUES ('Carlos López', 'carlos@ejemplo.com', 35);\n`
}
