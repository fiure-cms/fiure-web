package fcore

var SearchTextLimiter int = 100

// Storage Bucket Type
const (
	LiveStore       = "livestore"  // sniper
	IndexStore      = "indexstore" // boltdb
	AdminStore      = "adminstore" // boltdb
	SonicSearchMode = "sonicsearch"
)

// Storage Bucket List
const (
	ItemsBucket = "items" // LiveStore, AdminStore
	PagesBucket = "pages" // LiveStore
)

// Item Status Type
const (
	StatusDelete   = "delete"
	StatusWaiting  = "waiting"
	StatusPassive  = "passive"
	StatusModerate = "moderate"
	StatusActive   = "active"
	StatusSchedule = "schedule"
)

// CDN File Type
const (
	CDNFileTypeItem  = "item"
	CDNFileTypeImage = "img"
)

// Search Collection List
const (
	COLLECTION_ALL               = "all"
	COLLECTION_POST_MODELS       = "postmodels"    // Full Text Search
	COLLECTION_SUGGESTION_MODELS = "suggestmodels" // Full Text Search
)

// Search Bucket List
const (
	BUCKET_DEFAULT = "default"
	BUCKET_POST    = "post"
)
