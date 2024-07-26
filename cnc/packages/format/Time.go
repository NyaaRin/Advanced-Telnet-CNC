package format

import (
	"fmt"
	"time"
)

func Until(t time.Time, now time.Time) string {
	d := t.Sub(now)
	if d < 0 {
		return "Time has already passed"
	}

	switch {
	case d < time.Minute:
		return fmt.Sprintf("%0.f second%s", d.Seconds(), Multiple(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%0.f minute%s", d.Minutes(), Multiple(d.Minutes()))
	case d < time.Hour*24:
		return fmt.Sprintf("%0.1f hour%s", d.Hours(), Multiple(d.Hours()))
	case d < time.Hour*24*7:
		return fmt.Sprintf("%0.1f day%s", d.Hours()/24, Multiple(d.Hours()/24))
	case d < time.Hour*24*30:
		return fmt.Sprintf("%0.1f week%s", d.Hours()/(24*7), Multiple(d.Hours()/(24*7)))
	case d < time.Hour*24*30*12:
		return fmt.Sprintf("%0.1f month%s", d.Hours()/(24*30.44), Multiple(d.Hours()/(24*30.44)))
	case d < time.Hour*24*365:
		return fmt.Sprintf("%0.1f year%s", d.Hours()/(24*365), Multiple(d.Hours()/(24*365)))
	default:
		return fmt.Sprintf("%0.1f decade%s", d.Hours()/(24*365*10), Multiple(d.Hours()/(24*365*10)))
	}
}

func Since(t time.Time, now time.Time) string {
	d := now.Sub(t)
	if d < 0 {
		return "In the future"
	}

	switch {
	case d < time.Minute:
		return fmt.Sprintf("%0.f second%s ago", d.Seconds(), Multiple(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%0.f minute%s ago", d.Minutes(), Multiple(d.Minutes()))
	case d < time.Hour*24:
		return fmt.Sprintf("%0.1f hour%s ago", d.Hours(), Multiple(d.Hours()))
	case d < time.Hour*24*7:
		return fmt.Sprintf("%0.1f day%s ago", d.Hours()/24, Multiple(d.Hours()/24))
	case d < time.Hour*24*30:
		return fmt.Sprintf("%0.1f week%s ago", d.Hours()/(24*7), Multiple(d.Hours()/(24*7)))
	case d < time.Hour*24*30*12:
		return fmt.Sprintf("%0.1f month%s ago", d.Hours()/(24*30.44), Multiple(d.Hours()/(24*30.44)))
	case d < time.Hour*24*365:
		return fmt.Sprintf("%0.1f year%s ago", d.Hours()/(24*365), Multiple(d.Hours()/(24*365)))
	default:
		return fmt.Sprintf("%0.1f decade%s ago", d.Hours()/(24*365*10), Multiple(d.Hours()/(24*365*10)))
	}
}

func Format(t time.Time) string {
	now := time.Now()
	if t.After(now) {
		return Until(t, now)
	}

	return Since(t, now)
}

func Bool(bool2 bool) string {
	if bool2 {
		return "\x1b[32mtrue\x1b[0m"
	}

	return "\x1b[31mfalse\x1b[0m"
}
func Multiple(val float64) string {
	if val > 1 {
		return "s"
	}

	return ""
}
