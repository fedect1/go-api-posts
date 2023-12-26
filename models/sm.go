package models

type Secret struct {
	Host string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	JWTSing string `json:"jwtsign"`
	Database string `json:"database"`
}