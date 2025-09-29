package constants

const (
	// Order constants
	MinOrderAge         = 0
	MaxOrderAge         = 150
	MinOrderNameLength  = 2
	MaxOrderNameLength  = 100
	OrderTableName      = "orders"
	OrderIDColumn       = "id"
	OrderNameColumn     = "name"
	OrderEmailColumn    = "email"
	OrderCachePrefix    = "order:"
	OrderCacheTTL       = 300 // 5 minutes
	OrderAPIVersion     = "v1"
	OrderEndpoint       = "/orders"
	MaxOrderPerPage     = 100
	DefaultOrderPerPage = 20
)

// Status constants
const (
	OrderStatusActive   = "active"
	OrderStatusInactive = "inactive"
	OrderStatusPending  = "pending"
	OrderStatusDeleted  = "deleted"
)
