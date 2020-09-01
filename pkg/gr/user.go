// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-17

package gr

import (
	"encoding/xml"

	"github.com/pkg/errors"
)

const apiUser = "https://www.goodreads.com/api/auth_user"

// User is a Goodreads user.
type User struct {
	ID   int64  `xml:"id,attr"`
	Name string `xml:"name"`
}

// UserInfo retrieves user info from API.
func (c *Client) UserInfo() (User, error) {
	var (
		data []byte
		err  error
	)
	if data, err = c.apiRequest(apiUser); err != nil {
		return User{}, errors.Wrap(err, "contact user endpoint")
	}

	v := struct {
		User User `xml:"user"`
	}{}
	err = xml.Unmarshal(data, &v)
	return v.User, err
}
