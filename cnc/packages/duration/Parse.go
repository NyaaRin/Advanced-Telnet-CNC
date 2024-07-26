package duration

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ModifiedParseDuration(duration string) (time.Duration, error) {
	regexMap := map[string]func(int) time.Duration{
		`(\d+)\s*d`:   func(d int) time.Duration { return time.Duration(d) * 24 * time.Hour },
		`(\d+)\s*h`:   func(h int) time.Duration { return time.Duration(h) * time.Hour },
		`(\d+)\s*m`:   func(m int) time.Duration { return time.Duration(m) * time.Minute },
		`(\d+)\s*s\b`: func(s int) time.Duration { return time.Duration(s) * time.Second }, // Updated regex
		`(\d+)\s*w`:   func(y int) time.Duration { return time.Duration(y) * 7 * 24 * time.Hour },
		`(\d+)\s*mo`:  func(y int) time.Duration { return time.Duration(y) * 30 * 24 * time.Hour },
		`(\d+)\s*y`:   func(y int) time.Duration { return time.Duration(y) * 365 * 24 * time.Hour },
		`(\d+)\s*dec`: func(dec int) time.Duration { return time.Duration(dec) * 10 * 365 * 24 * time.Hour },
	}

	for regex, durationFn := range regexMap {
		re := regexp.MustCompile(regex)
		matches := re.FindAllStringSubmatch(duration, -1)
		for _, match := range matches {
			value, err := strconv.Atoi(match[1])
			if err != nil {
				return 0, fmt.Errorf("invalid duration: %s", duration)
			}
			replaceStr := durationFn(value).String()
			duration = strings.Replace(duration, match[0], replaceStr, 1)
		}
	}

	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		return 0, err
	}

	return parsedDuration, nil
}

func FormatDuration(duration time.Duration) string {
	days := int(duration.Hours() / 24)
	weeks := days / 7
	months := days / 30
	years := days / 365
	decades := years / 10

	duration = duration % (24 * time.Hour)
	hours := int(duration.Hours())
	duration = duration % time.Hour

	var parts []string

	if decades > 0 {
		parts = append(parts, fmt.Sprintf("%d decade%s", decades, sOrNot(decades)))
	} else if years > 0 {
		parts = append(parts, fmt.Sprintf("%d year%s", years, sOrNot(years)))
	} else if weeks > 0 {
		parts = append(parts, fmt.Sprintf("%d week%s", weeks, sOrNot(weeks)))
	} else if months > 0 {
		parts = append(parts, fmt.Sprintf("%d month%s", months, sOrNot(months)))
	} else if days > 0 {
		parts = append(parts, fmt.Sprintf("%d day%s", days, sOrNot(days)))
	}

	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hour%s", hours, sOrNot(hours)))
	}

	return strings.Join(parts, ", ")
}

func sOrNot(value int) string {
	if value > 1 {
		return "s"
	}

	return ""
}
