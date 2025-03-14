package entity

import "net/http"

type Cookies []*http.Cookie

func (c Cookies) String() string {
	var cookies string

	for _, cookie := range c {
		cookies += cookie.String() + "; "
	}

	return cookies
}
