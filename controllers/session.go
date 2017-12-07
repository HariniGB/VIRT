package controllers

import (
	"log"
	"net/http"
	"github.com/sony/sonyflake"
	"github.com/gorilla/securecookie"
)

// cookie handling
var hashKey = []byte("12345")
var blockKey = []byte("1234567890123456")

var flake = sonyflake.NewSonyflake(sonyflake.Settings{})
var cookieHandler = securecookie.New(hashKey, blockKey)

func getSession(request *http.Request) (string, string){
	c, err := request.Cookie("session")
	if err != nil {
		log.Println(err)
		return "", ""
	}
	values := make(map[string]string)
	err = cookieHandler.Decode("session", c.Value, &values)
	if err != nil {
		log.Println(err)
		return "", ""
	}

	user, _ := values["name"]
	pass, _ := values["token"]
	return user, pass
}

// Clear the current session. Called in Logout handler
func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

//Create  a session for the successful login user
func setSession(userName, token string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
		"token": token,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	} else {
		log.Println(err)
	}
}
