package cookies

import (
	"strings"
)

type FacebookCookies struct {
	Datr 	  string `cookie:"datr,omitempty"`
	Sb   	  string `cookie:"sb,omitempty"`
	AccountId string `cookie:"c_user,omitempty"`
	Xs 		  string `cookie:"xs,omitempty"`
	Fr 		  string `cookie:"fr,omitempty"`
	Wd 		  string `cookie:"wd,omitempty"`
	Presence  string `cookie:"presence,omitempty"`
}

func (fb *FacebookCookies) GetValue(name string) string {
	return getCookieValue(name, fb)
}

func (fb *FacebookCookies) IsLoggedIn() bool {
	return fb.Xs != ""
}

func (fb *FacebookCookies) GetViewports() (string, string) {
	pxs := strings.Split(fb.Wd, "x")
	if len(pxs) != 2 {
		return "2276", "1156"
	}
	return pxs[0], pxs[1]
}