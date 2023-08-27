package types

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Cookies struct {
	Datr string `json:"datr,omitempty"`
	Sb string `json:"sb,omitempty"`
	AccountId string `json:"c_user,omitempty"`
	Xs string `json:"xs,omitempty"`
	Fr string `json:"fr,omitempty"`
	Wd string `json:"wd,omitempty"`
	Presence string `json:"presence,omitempty"`
}

func (c *Cookies) GetViewports() (string, string) {
	pxs := strings.Split(c.Wd, "x")
	return pxs[0], pxs[1]
}

func (c *Cookies) ToString() string {
	s := ""
	values := reflect.ValueOf(*c)
	for i := 0; i < values.NumField(); i++ {
		field := values.Type().Field(i)
		value := values.Field(i).Interface()
	
		zeroValue := reflect.Zero(field.Type).Interface()
		if value == zeroValue {
			continue
		}
	
		tagValue := field.Tag.Get("json")
		tagName := strings.Split(tagValue, ",")[0]
		s += fmt.Sprintf("%s=%v; ", tagName, value)
	}
	
	return s
}

// FROM JSON FILE.
func NewCookiesFromFile(path string) (*Cookies, error) {
	jsonBytes, jsonErr := os.ReadFile(path)
	if jsonErr != nil {
		return nil, jsonErr
	}

	session := &Cookies{}

	marshalErr := json.Unmarshal(jsonBytes, session)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return session, nil
}


func NewCookiesFromString(cookieStr string) *Cookies {
	datr := extractCookieValue(cookieStr, "datr")
	sb := extractCookieValue(cookieStr, "sb")
	accountId := extractCookieValue(cookieStr, "c_user")
	xs := extractCookieValue(cookieStr, "xs")
	fr := extractCookieValue(cookieStr, "fr")
	wd := extractCookieValue(cookieStr, "wd")
	presence := extractCookieValue(cookieStr, "presence")

	return &Cookies{
		Datr: datr,
		Sb: sb,
		AccountId: accountId,
		Xs: xs,
		Fr: fr,
		Wd: wd,
		Presence: presence,
	}
}



func extractCookieValue(cookieString, key string) string {
	startIndex := strings.Index(cookieString, key)
	if startIndex == -1 {
		return ""
	}

	startIndex += len(key) + 1
	endIndex := strings.IndexByte(cookieString[startIndex:], ';')
	if endIndex == -1 {
		return cookieString[startIndex:]
	}

	return cookieString[startIndex : startIndex+endIndex]
}