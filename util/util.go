package util

import (
	"errors"
	"strings"
)

func UrlValidator(url string) error {
	//有一些url不合法,只打回  .. 好了
	if strings.Contains(url, "..") {
		return errors.New("not right url")
	}
	return nil
}
