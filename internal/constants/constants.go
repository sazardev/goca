package constants

const (
	// TestPerf constants
	MinTestPerfAge         = 0
	MaxTestPerfAge         = 150
	MinTestPerfNameLength  = 2
	MaxTestPerfNameLength  = 100
	TestPerfTableName      = "testperfs"
	TestPerfIDColumn       = "id"
	TestPerfNameColumn     = "name"
	TestPerfEmailColumn    = "email"
	TestPerfCachePrefix    = "testperf:"
	TestPerfCacheTTL       = 300 // 5 minutes
	TestPerfAPIVersion     = "v1"
	TestPerfEndpoint       = "/testperfs"
	MaxTestPerfPerPage     = 100
	DefaultTestPerfPerPage = 20
)

// Status constants
const (
	TestPerfStatusActive   = "active"
	TestPerfStatusInactive = "inactive"
	TestPerfStatusPending  = "pending"
	TestPerfStatusDeleted  = "deleted"
)
