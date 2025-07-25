package messages

const (
	// TestUserFeature errors
	ErrTestUserFeatureNotFound      = "testuserfeature not found"
	ErrTestUserFeatureAlreadyExists = "testuserfeature already exists"
	ErrInvalidTestUserFeatureData   = "invalid testuserfeature data"
	ErrTestUserFeatureEmailRequired = "testuserfeature email is required"
	ErrTestUserFeatureNameRequired  = "testuserfeature name is required"
	ErrTestUserFeatureAgeInvalid    = "testuserfeature age must be positive"
	ErrTestUserFeatureAccessDenied  = "access denied to testuserfeature"
	ErrTestUserFeatureUpdateFailed  = "failed to update testuserfeature"
	ErrTestUserFeatureDeleteFailed  = "failed to delete testuserfeature"
)
