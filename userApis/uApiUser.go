package userApis

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_v3/connection"
	"project_v3/data"
	"project_v3/mail"
	"project_v3/redis"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/dgrijalva/jwt-go/v4/request"
	"github.com/gorilla/mux"
)

//verifyKey, signKey []byte
var (
	signKey = []byte("secret_key")
)

type tokenres struct {
	ID    int64  `json:"id"`
	Token string `json:"token,omitempty"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	//Read body of request, hash password and create a new user
	var input data.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "info register not enough", http.StatusBadRequest)
		return
	}

	//Insert new user to db, if having error, username or email has exited
	err = connection.Model.User.Insert(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "username or email existed", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//read body of request, make new user
	var input struct {
		Username *string `json:"user_name"`
		Email    *string `json:"email"`
		Password string  `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not decode", http.StatusBadRequest)
		return
	}
	var username string
	var email string
	if input.Email != nil {
		email = *input.Email
	} else {
		email = ""
	}
	if input.Username != nil {
		username = *input.Username
	} else {
		username = ""
	}
	user, err := connection.Model.User.Login(username, email, input.Password)
	if user.UserId == 0 {
		log.Println(err)
		http.Error(w, "wrong username or password", http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "wrong username or password", http.StatusInternalServerError)
		return
	}
	var name string
	if (user.FirstName != "") || (user.LastName != "") {
		name = user.FirstName + " " + user.LastName
	} else {
		name = user.Username
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour * time.Duration(15)).Unix(),
		"iat":   time.Now().Unix(),
		"id":    user.UserId,
		"email": user.Email,
		"name":  name,
		"iss":   "admin",
	})
	ss, err := token.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "error while signing token")
		log.Printf("token signing err :%v\n", err)
		return
	}
	response := tokenres{user.UserId, ss}
	_ = json.NewEncoder(w).Encode(response)

}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "token expires", http.StatusUnauthorized)
		return
	}
	if token.Valid {
		tokenString := token.Raw
		claims := jwt.MapClaims{}
		token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return signKey, nil
		})
		// ... error handling
		if err != nil {
			log.Println(err)
			return
		}
		// do something with decoded claims
		var id float64
		for key, val := range claims {
			if key == "id" {
				id = val.(float64)
			}
		}
		user := &data.User{}
		err = connection.Model.User.GetUserById(id, user)
		if err != nil {
			http.Error(w, "Account deleted", http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(user)
	} else {
		_ = json.NewEncoder(w).Encode("Invalid token")
	}
}

func UpdateInfoHandler(w http.ResponseWriter, r *http.Request) {
	var input data.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not decode user", http.StatusBadRequest)
		return
	}
	err = connection.Model.User.Update(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not update user", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(http.StatusOK)
}

func ResetPassHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	ok, _ := connection.Model.User.CheckEmail(input.Email)
	if !ok {
		http.Error(w, "Did not find your email", http.StatusBadRequest)
		return
	}
	randomBytes := make([]byte, 16)
	_, err = rand.Read(randomBytes)
	if err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	err = redis.Client.Set(token, input.Email, time.Hour*24).Err()

	if err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	mail.Pre.PreSend2(token, input.Email)
}

func ChangePassHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "token expires", http.StatusUnauthorized)
		return
	}
	if token.Valid {
		tokenString := token.Raw
		claims := jwt.MapClaims{}
		token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return signKey, nil
		})
		// ... error handling
		if err != nil {
			log.Println(err)
			return
		}
		// do something with decoded claims
		var email string
		for key, val := range claims {
			if key == "email" {
				email = val.(string)
			}
		}
		ok, _ := connection.Model.User.ChangePass(input.Password, email)
		if !ok {
			http.Error(w, "error change pass", 400)
			return
		}
	}
}

func IssueJwtHandler(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)
	val, err := redis.Client.Get(token["token"]).Result()
	if err != nil {
		log.Println(err)
		http.Error(w, "token expired or invalid", 400)
		return
	}
	redis.Client.Del(token["token"]).Result()
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(time.Minute * time.Duration(15)).Unix(),
		"iat":   time.Now().Unix(),
		"email": val,
		"iss":   "admin",
	})
	ss, err := jwt_token.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "error while signing token")
		log.Printf("token signing err :%v\n", err)
		return
	}
	var response struct {
		Token string `json:"token"`
	}
	response.Token = ss
	json.NewEncoder(w).Encode(response)
}
