package models

import "time"

type User struct {
	UserName         string
	Password         string
	VerificationCode byte
	ExpirePeriod     time.Time
}

type SSHUser struct {
	UserName string
	Password string
}
