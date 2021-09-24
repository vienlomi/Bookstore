package authen

import (
	"config"
	"data"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type Token struct {
	Id        int64         `json:"id"`
	FirstName string        `json:"first_name"`
	Email     string        `json:"email"`
	Admin     bool          `json:"admin"`
	Exp       time.Duration `exp`
}

func CreateToken(user data.User) (string, error) {
	var err error
	//Creating Access Token
	_ = os.Setenv("ACCESS_SECRET", config.DataConfig["key_token"]) //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["id"] = user.UserId
	atClaims["first_name"] = user.FirstName
	atClaims["email"] = user.Email
	if user.Role == "admin" {
		atClaims["admin"] = true
	} else {
		atClaims["admin"] = false
	}
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
