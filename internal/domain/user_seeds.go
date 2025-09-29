package domain

// GetUserSeeds retorna datos de ejemplo para user
func GetUserSeeds() []User {
	return []User{
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

// GetSQLUserSeeds retorna sentencias SQL INSERT para user
func GetSQLUserSeeds() string {
	return `-- Datos de ejemplo para tabla user
INSERT INTO users (name, email, age) VALUES ('Juan Pérez', 'juan@ejemplo.com', 25);\nINSERT INTO users (name, email, age) VALUES ('María García', 'maria@ejemplo.com', 30);\nINSERT INTO users (name, email, age) VALUES ('Carlos López', 'carlos@ejemplo.com', 35);\n`
}
