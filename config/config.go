package config

type Env struct {
	GOOGLE_GMAIL_KEY     string
	ACCESS_TOKEN_SECRET  string
	REFRESH_TOKEN_SECRET string
	PGHOST               string
	PGPORT               string
	PGUSER               string
	PGPASSWORD           string
	PGDATABASE           string
	PORT                 string
}

var EnvData Env
