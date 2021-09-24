package small

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"project_v3/connection"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/dgrijalva/jwt-go/v4/request"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

var (
	signKey = []byte("secret_key")
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	result, err := connection.Model.Order.GetOrders(id["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "loi roi", http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		http.Error(w, "loi roi", http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var Input struct {
		Revenue int64 `json:"revenue"`
	}
	stripe.Key = "sk_test_4eC39HqLyjWDarjtT1zdp7dc"
	err := json.NewDecoder(r.Body).Decode(&Input)
	if err != nil {
		log.Println(err)
	}
	log.Println(Input)
	log.Println(Input.Revenue)
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(Input.Revenue),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("pi.New: %v", err)
		return
	}

	writeJSON(w, struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: pi.ClientSecret,
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}

func AuthMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			id_new := mux.Vars(r)
			log.Println(id_new["id"])
			if id_new["id"] == "" {
				next(w, r)
				return
			}
			id_f, err := strconv.ParseFloat(id_new["id"], 64)
			if err != nil {
				http.Error(w, "nhighc daij a", http.StatusUnauthorized)
				return
			}
			if id_f != id {
				http.Error(w, "dung nghich", http.StatusUnauthorized)
				return
			} else {
				next(w, r)
			}

		} else {
			_ = json.NewEncoder(w).Encode("Invalid token")
		}
	}
}
