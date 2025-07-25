package constants

const (
	// TestUserFeature constants
	MinTestUserFeatureAge         = 0
	MaxTestUserFeatureAge         = 150
	MinTestUserFeatureNameLength  = 2
	MaxTestUserFeatureNameLength  = 100
	TestUserFeatureTableName      = "testuserfeatures"
	TestUserFeatureIDColumn       = "id"
	TestUserFeatureNameColumn     = "name"
	TestUserFeatureEmailColumn    = "email"
	TestUserFeatureCachePrefix    = "testuserfeature:"
	TestUserFeatureCacheTTL       = 300 // 5 minutes
	TestUserFeatureAPIVersion     = "v1"
	TestUserFeatureEndpoint       = "/testuserfeatures"
	MaxTestUserFeaturePerPage     = 100
	DefaultTestUserFeaturePerPage = 20
)

// Status constants
const (
	TestUserFeatureStatusActive   = "active"
	TestUserFeatureStatusInactive = "inactive"
	TestUserFeatureStatusPending  = "pending"
	TestUserFeatureStatusDeleted  = "deleted"
)
