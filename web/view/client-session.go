package view

import (
	"encoding/base64"
	"encoding/json"
	"flashcard"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

const cookie_max_size = 2048
const cookie_prefix = "memorizing-session"

func EncodeMemorizingSessionToCookies(session *flashcard.MemorizingSessionDTO) map[string]string {
	sessionJson, _ := json.Marshal(session)
	sessionEncoded := base64.StdEncoding.EncodeToString(sessionJson)

	cookies := make(map[string]string, 0)

	for i := 0; i < len(sessionEncoded); i += cookie_max_size {
		cookieName := cookie_prefix + "-p-" + fmt.Sprint(len(cookies))
		cookieValue := sessionEncoded[i:int(math.Min(float64(len(sessionEncoded)), float64(i+2048)))]
		cookies[cookieName] = cookieValue
	}

	noOfPages := len(cookies)

	cookies[cookie_prefix+"-len"] = fmt.Sprint(noOfPages)

	return cookies
}

func DecodeCookiesToMemorizingSession(r *http.Request) *flashcard.MemorizingSessionDTO {
	cookies := make(map[string]string, 0)

	for _, str := range r.Cookies() {
		cookies[str.Name] = str.Value
	}

	noOfPages, _ := strconv.ParseInt(cookies[cookie_prefix+"-len"], 10, 64)

	memorizingSessionEncoded := ""
	for i := 0; i < int(noOfPages); i++ {
		memorizingSessionEncoded += cookies[cookie_prefix+"-p-"+fmt.Sprint(i)]
	}

	memorizingSessionCookieJson, _ := base64.StdEncoding.DecodeString(memorizingSessionEncoded)
	memorizingSession := flashcard.MemorizingSessionDTO{}
	_ = json.Unmarshal(memorizingSessionCookieJson, &memorizingSession)

	return &memorizingSession
}
