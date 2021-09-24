package authen

import (
	"connection"
	"data"
	"fmt"
	"testing"
)

func TestAuth(t *testing.T) {
	connection.Start("/home/binh/Desktop/project_v3/src/config/config.csv")

	var user data.User
	user.UserId = 1
	user.FirstName = "binh"
	user.Email = "mail"
	token, err := CreateToken(user)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(token)
	}

}
