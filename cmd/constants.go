package cmd

// Flag names - Nombres de flags
const (
	DatabaseFlag       = "database"
	FieldsFlag         = "fields"
	InterfaceOnlyFlag  = "interface-only"
	ImplementationFlag = "implementation"
	CacheFlag          = "cache"
	TransactionsFlag   = "transactions"
	HTTPFlag           = "http"
	GRPCFlag           = "grpc"
	GraphQLFlag        = "graphql"
)

// Flag usage messages - Mensajes de uso de flags
const (
	DatabaseFlagUsage       = "Tipo de base de datos (postgres, mysql, mongodb)"
	FieldsFlagUsage         = "Lista de campos separados por comas (ej: name:string,age:int)"
	InterfaceOnlyFlagUsage  = "Solo generar interfaces"
	ImplementationFlagUsage = "Solo generar implementación"
	CacheFlagUsage          = "Incluir capa de caché"
	TransactionsFlagUsage   = "Incluir soporte para transacciones"
	HTTPFlagUsage           = "Incluir handlers HTTP"
	GRPCFlagUsage           = "Incluir handlers gRPC"
	GraphQLFlagUsage        = "Incluir handlers GraphQL"
)

// Database constants
const (
	DBPostgres = "postgres"
	DBMySQL    = "mysql"
	DBMongoDB  = "mongodb"
	DBSQLite   = "sqlite"
)

// Valid database types
var ValidDatabases = []string{DBPostgres, DBMySQL, DBMongoDB, DBSQLite}

// Handler/Protocol constants
const (
	HandlerHTTP   = "http"
	HandlerGRPC   = "grpc"
	HandlerCLI    = "cli"
	HandlerWorker = "worker"
)

// Valid handler types
var ValidHandlers = []string{HandlerHTTP, HandlerGRPC, HandlerCLI, HandlerWorker}

// CRUD Operations constants
const (
	OpCreate = "create"
	OpRead   = "read"
	OpUpdate = "update"
	OpDelete = "delete"
	OpList   = "list"
)

// Default operation combinations
const (
	DefaultOperations = "create,read,update,delete,list"
	BasicOperations   = "create,read"
	CRUDOperations    = "create,read,update,delete"
)

// Valid operations
var ValidOperations = []string{OpCreate, OpRead, OpUpdate, OpDelete, OpList}

// API Types constants
const (
	APITypeRest    = "rest"
	APITypeGraphQL = "graphql"
	APITypeGRPC    = "grpc"
)

// Valid API types
var ValidAPITypes = []string{APITypeRest, APITypeGraphQL, APITypeGRPC}

// Field Types constants
const (
	FieldString    = "string"
	FieldInt       = "int"
	FieldInt64     = "int64"
	FieldUint      = "uint"
	FieldUint64    = "uint64"
	FieldFloat32   = "float32"
	FieldFloat64   = "float64"
	FieldBool      = "bool"
	FieldTime      = "time.Time"
	FieldBytes     = "[]byte"
	FieldInterface = "interface{}"
)

// Valid field types
var ValidFieldTypes = []string{
	FieldString, FieldInt, FieldInt64, FieldUint, FieldUint64,
	FieldFloat32, FieldFloat64, FieldBool, FieldTime, FieldBytes, FieldInterface,
}

// Template constants
const (
	TemplateEntity     = "entity"
	TemplateUseCase    = "usecase"
	TemplateRepository = "repository"
	TemplateHandler    = "handler"
	TemplateDI         = "di"
)

// File extensions
const (
	ExtGo   = ".go"
	ExtYAML = ".yaml"
	ExtYML  = ".yml"
	ExtJSON = ".json"
	ExtSQL  = ".sql"
)

// Directory names
const (
	DirInternal   = "internal"
	DirDomain     = "domain"
	DirUseCase    = "usecase"
	DirRepository = "repository"
	DirHandler    = "handler"
	DirHTTP       = "http"
	DirGRPC       = "grpc"
	DirCLI        = "cli"
	DirWorker     = "worker"
	DirSOAP       = "soap"
	DirMessages   = "messages"
	DirInterfaces = "interfaces"
	DirPkg        = "pkg"
	DirConfig     = "config"
	DirLogger     = "logger"
	DirAuth       = "auth"
	DirCmd        = "cmd"
	DirServer     = "server"
	DirMigrations = "migrations"
)

// Common field names that might be used for queries
var CommonQueryFields = map[string][]string{
	"user":     {"email", "username", "id"},
	"product":  {"name", "sku", "id", "category"},
	"order":    {"id", "user_id", "status"},
	"customer": {"email", "id", "phone"},
	"article":  {"title", "slug", "id"},
	"post":     {"title", "slug", "id", "author_id"},
	"category": {"name", "slug", "id"},
}

// Validation constants
const (
	MinFieldNameLength  = 1
	MaxFieldNameLength  = 50
	MinEntityNameLength = 1
	MaxEntityNameLength = 50
)

// Error messages
const (
	ErrInvalidDatabase    = "base de datos no válida. Opciones: postgres, mysql, mongodb, sqlite"
	ErrInvalidHandler     = "handler no válido. Opciones: http, grpc, cli, worker"
	ErrInvalidOperation   = "operación no válida. Opciones: create, read, update, delete, list"
	ErrInvalidFieldType   = "tipo de campo no válido"
	ErrInvalidFieldSyntax = "sintaxis de campo no válida. Formato esperado: 'nombre:tipo'"
	ErrInvalidEntityName  = "nombre de entidad no válido"
	ErrEmptyFields        = "campos no pueden estar vacíos"
)
