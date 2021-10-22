package database

// Config version 1.0 only support mysql
type Config struct {
	// db username
	Username string `json:"username"`
	// db password
	Password string `json:"password"`
	// db url
	Url string `json:"url"`
	// db port
	Port string `json:"port"`
	// database name
	DatabaseName string `json:"database_name"`
	// db MaxIdleConns
	MaxIdleConns int
	// db MaxOpenConns
	MaxOpenConns int
	// db ConnMaxLifeTime
	ConnMaxLifeTime string
}
