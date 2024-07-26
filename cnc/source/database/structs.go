package database

import "time"

type UserProfile struct {
	Id         int
	Name       string
	Password   string
	Methods    []int
	Cooldown   int
	MaxTime    int
	MaxAttacks int `lufx:"Attacks"`
	Devices    int
	Expiry     time.Time
	Reseller   bool
	Admin      bool
	Role       string
}
