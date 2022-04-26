package model

// JSONData wrapper struct
type JSONData map[string]interface{}

// Config the main config structure
type Config struct {
	MongoDBAuth  string
	MongoDBName  string
	MongoTimeout string

	AuthKeys       string
	InternalAPIKey string
	CoreBBHost     string
	PollServiceURL string
	UiucOrgID      string
}
