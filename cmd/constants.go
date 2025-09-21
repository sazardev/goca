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

// Flag usage messages - Flag usage messages
const (
	DatabaseFlagUsage       = "Database type (postgres, mysql, mongodb)"
	FieldsFlagUsage         = "Comma-separated list of fields (ex: name:string,age:int)"
	InterfaceOnlyFlagUsage  = "Generate interfaces only"
	ImplementationFlagUsage = "Generate implementation only"
	CacheFlagUsage          = "Include cache layer"
	TransactionsFlagUsage   = "Include transaction support"
	HTTPFlagUsage           = "Include HTTP handlers"
	GRPCFlagUsage           = "Include gRPC handlers"
	GraphQLFlagUsage        = "Include GraphQL handlers"
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
	ErrInvalidDatabase    = "invalid database. Options: postgres, mysql, mongodb, sqlite"
	ErrInvalidHandler     = "invalid handler. Options: http, grpc, cli, worker"
	ErrInvalidOperation   = "invalid operation. Options: create, read, update, delete, list"
	ErrInvalidFieldType   = "invalid field type"
	ErrInvalidFieldSyntax = "invalid field syntax. Expected format: 'name:type'"
	ErrInvalidEntityName  = "invalid entity name"
	ErrEmptyFields        = "fields cannot be empty"
	ErrRequiredFlag       = "required flag not provided"
	ErrFileNotFound       = "file not found"
	ErrDirectoryNotFound  = "directory not found"
)

// Success messages
const (
	MsgEntityGenerated     = "Entity '%s' generated successfully!"
	MsgFeatureGenerated    = "Feature '%s' generated and integrated successfully!"
	MsgRepositoryGenerated = "Repository for '%s' generated successfully!"
	MsgHandlerGenerated    = "Handler '%s' for '%s' generated successfully!"
	MsgUseCaseGenerated    = "Use case '%s' generated successfully!"
	MsgProjectInitialized  = "Project '%s' created successfully!"
)

// Info messages
const (
	MsgGeneratingEntity     = "Generating entity '%s'"
	MsgGeneratingFeature    = "Generating complete feature '%s'"
	MsgGeneratingRepository = "Generating repository for entity '%s'"
	MsgGeneratingHandler    = "Generating handler '%s' for entity '%s'"
	MsgGeneratingUseCase    = "Generating use case '%s' for entity '%s'"
	MsgGeneratingLayers     = "Generating layers..."
)
