package config

type Config struct {
	Server struct {
		Address string
		Port    string
	}
	Repository struct {
	}
}

func New() Config {
	return Config{Server: struct {
		Address string
		Port    string
	}{Address: "127.0.0.1", Port: "8080"}}
}
