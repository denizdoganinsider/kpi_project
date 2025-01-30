package mysql

type Config struct {
	Host                  string
	Port                  string
	UserName              string
	Password              string
	DbName                string
	MaxConnections        int
	MaxConnectionIdleTime int
}
