package config

type Env struct {
	SMTP_KEY             string
	SMTP_ADDR            string
	SMTP_LOGIN           string
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
