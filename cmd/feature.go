package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var featureCmd = &cobra.Command{
	Use:   "feature <name>",
	Short: "Generar feature completo con Clean Architecture",
	Long: `Genera todas las capas necesarias para un feature completo, 
incluyendo dominio, casos de uso, repositorio y handlers en una sola operación.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		featureName := args[0]

		fields, _ := cmd.Flags().GetString("fields")
		database, _ := cmd.Flags().GetString("database")
		handlers, _ := cmd.Flags().GetString("handlers")
		validation, _ := cmd.Flags().GetBool("validation")
		businessRules, _ := cmd.Flags().GetBool("business-rules")

		if fields == "" {
			fmt.Println("Error: --fields flag es requerido")
			return
		}

		fmt.Printf("🚀 Generando feature completo '%s'\n", featureName)
		fmt.Printf("📋 Campos: %s\n", fields)
		fmt.Printf("🗄️  Base de datos: %s\n", database)
		fmt.Printf("🌐 Handlers: %s\n", handlers)

		if validation {
			fmt.Println("✅ Incluyendo validaciones")
		}
		if businessRules {
			fmt.Println("🧠 Incluyendo reglas de negocio")
		}

		generateCompleteFeature(featureName, fields, database, handlers, validation, businessRules)

		fmt.Printf("\n🎉 Feature '%s' generado exitosamente!\n", featureName)
		fmt.Println("\n📂 Estructura generada:")
		printFeatureStructure(featureName, handlers)

		fmt.Println("\n📝 Próximos pasos:")
		fmt.Println("1. Revisar y ajustar las entidades generadas")
		fmt.Println("2. Implementar lógica de negocio específica")
		fmt.Println("3. Configurar la inyección de dependencias")
		fmt.Printf("4. Registrar las rutas en tu servidor principal\n")
	},
}

func generateCompleteFeature(featureName, fields, database, handlers string, validation, businessRules bool) {
	fmt.Println("\n🔄 Generando capas...")

	// 1. Generate Entity (Domain layer)
	fmt.Println("1️⃣  Generando entidad de dominio...")
	generateEntity(featureName, fields, true, businessRules, false, false)

	// 2. Generate Use Case
	fmt.Println("2️⃣  Generando casos de uso...")
	generateUseCase(featureName+"UseCase", featureName, "create,read,update,delete,list", validation, false)

	// 3. Generate Repository
	fmt.Println("3️⃣  Generando repositorio...")
	generateRepository(featureName, database, false, true, false, false)

	// 4. Generate Handlers
	fmt.Println("4️⃣  Generando handlers...")
	handlerTypes := strings.Split(handlers, ",")
	for _, handlerType := range handlerTypes {
		handlerType = strings.TrimSpace(handlerType)
		fmt.Printf("   📡 Generando handler %s...\n", handlerType)
		generateHandler(featureName, handlerType, true, validation, handlerType == "http")
	}

	// 5. Generate Messages
	fmt.Println("5️⃣  Generando mensajes...")
	generateMessages(featureName, true, true, true)

	fmt.Println("✅ Todas las capas generadas exitosamente!")
}

func printFeatureStructure(featureName, handlers string) {
	featureLower := strings.ToLower(featureName)

	fmt.Printf(`%s/
├── domain/
│   ├── %s.go          # Entidad pura
│   ├── errors.go      # Errores de dominio
│   └── validations.go # Validaciones de negocio
├── usecase/
│   ├── dto.go              # DTOs de entrada/salida
│   ├── %s_usecase.go       # Interfaz de casos de uso
│   ├── %s_service.go       # Implementación de casos de uso
│   └── interfaces.go       # Contratos hacia otras capas
├── repository/
│   ├── interfaces.go       # Contratos de persistencia
│   └── postgres_%s_repo.go # Implementación PostgreSQL
├── handler/`, featureName, featureLower, featureLower, featureLower, featureLower)

	handlerTypes := strings.Split(handlers, ",")
	for _, handlerType := range handlerTypes {
		handlerType = strings.TrimSpace(handlerType)
		switch handlerType {
		case "http":
			fmt.Printf(`
│   ├── http/
│   │   ├── %s_handler.go   # Controlador HTTP
│   │   └── routes.go       # Rutas HTTP`, featureLower)
		case "grpc":
			fmt.Printf(`
│   ├── grpc/
│   │   ├── %s.proto        # Definición gRPC
│   │   └── %s_server.go    # Servidor gRPC`, featureLower, featureLower)
		case "cli":
			fmt.Printf(`
│   ├── cli/
│   │   └── %s_commands.go  # Comandos CLI`, featureLower)
		case "worker":
			fmt.Printf(`
│   ├── worker/
│   │   └── %s_worker.go    # Workers/Jobs`, featureLower)
		case "soap":
			fmt.Printf(`
│   ├── soap/
│   │   └── %s_client.go    # Cliente SOAP`, featureLower)
		}
	}

	fmt.Printf(`
└── messages/
    ├── errors.go       # Mensajes de error
    └── responses.go    # Mensajes de respuesta
`)
}

func init() {
	featureCmd.Flags().StringP("fields", "f", "", "Campos de la entidad \"field:type,field2:type\" (requerido)")
	featureCmd.Flags().StringP("database", "d", "postgres", "Tipo de base de datos (postgres, mysql, mongodb)")
	featureCmd.Flags().StringP("handlers", "", "http", "Tipos de handlers \"http,grpc,cli\"")
	featureCmd.Flags().BoolP("validation", "v", false, "Incluir validaciones en todas las capas")
	featureCmd.Flags().BoolP("business-rules", "b", false, "Incluir métodos de reglas de negocio")

	featureCmd.MarkFlagRequired("fields")
}
