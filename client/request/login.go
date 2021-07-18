package request

import (
	"errors"
)

func Login(username, pwd string) error {
	data := PostData()
	data.Add("username", username)
	data.Add("password", pwd)
	rs, err := Post("/api/login", data)
	if err != nil {
		return err
	}
	message, err := ParseResult(rs)
	if err != nil {
		return err
	}
	if message.Code != int(OKStatus) {
		return errors.New(message.Msg)
	}
	return nil
}
