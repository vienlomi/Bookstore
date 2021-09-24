package main

import (
	"encoding/json"
	"net/http"
	"project_v3/adminApis"
	"project_v3/connection"
	"project_v3/helperApis"
	"project_v3/mail"
	"project_v3/redis"
	"project_v3/small"
	"project_v3/userApis"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	connection.Start("./config/config.csv")
	redis.NewClient()
	mail.StartMQ()

	go func() {
		for {
			time.Sleep(time.Second * 30)
			if mail.Channel.Conn.IsClosed() == true {
				mail.StartMQ()
			}
		}
	}()

	defer mail.CloseMQ()
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.Methods(http.MethodGet).Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode("welcome")
	})

	router.Methods(http.MethodGet).Path("/products").HandlerFunc(userApis.GetAllProductHandler)
	router.Methods(http.MethodGet).Path("/products/{id}").HandlerFunc(userApis.GetDetailProductHandler)
	router.Methods(http.MethodPost).Path("/products").HandlerFunc(adminApis.CreateProductHandler)
	router.Methods(http.MethodPut).Path("/products").HandlerFunc(adminApis.UpdateProductHandler)
	router.Methods(http.MethodDelete).Path("/products/{id}").HandlerFunc(adminApis.DeleteProductHandler)
	router.Methods(http.MethodGet).Path("/products/list/{order-by}").HandlerFunc(userApis.GetListProductHandler)
	router.Methods(http.MethodPost).Path("/registry").HandlerFunc(userApis.RegisterHandler)
	router.Methods(http.MethodPost).Path("/login").HandlerFunc(userApis.LoginHandler)
	router.Methods(http.MethodPost).Path("/user").HandlerFunc(userApis.UpdateInfoHandler)
	//router.Methods(http.MethodGet).Path("/delete-redis").HandlerFunc(adminApis.DeleteAllKeyRedis)
	router.Methods(http.MethodGet).Path("/admin-user").HandlerFunc(adminApis.GetAllUserHandler)
	router.Methods(http.MethodDelete).Path("/admin-user/{id}").HandlerFunc(adminApis.DeleteUserHandler)
	router.HandleFunc("/users/orders/{id}", userApis.GetOrdersOfUserHandler).Methods("POST")
	router.Methods(http.MethodGet).Path("/admin/orders").HandlerFunc(adminApis.GetOrdersDESCHandler)
	router.Methods(http.MethodGet).Path("/orders/{id}").HandlerFunc(userApis.GetDetailOrderHandler)

	router.HandleFunc("/orders", small.AuthMiddleWare(userApis.CreateOrderHandler)).Methods("POST")
	router.HandleFunc("/create-payment-intent", small.AuthMiddleWare(small.CreatePaymentIntent))
	router.HandleFunc("/get-order/{id}", small.AuthMiddleWare(userApis.GetOrdersOfUserHandler)).Methods("POST")
	router.HandleFunc("/authen", userApis.AuthHandler).Methods("POST")
	router.HandleFunc("/reset", userApis.ResetPassHandler).Methods("POST")
	router.HandleFunc("/change-pass", userApis.ChangePassHandler).Methods("POST")
	router.HandleFunc("/reset-token/{token}", userApis.IssueJwtHandler).Methods("GET")

	router.Methods(http.MethodPost).Path("/home-list/{id}").HandlerFunc(adminApis.SetListHomeHandler)
	router.Methods(http.MethodGet).Path("/home-list/{id}").HandlerFunc(adminApis.GetListHomeHandler)
	router.Methods(http.MethodGet).Path("/home-list").HandlerFunc(adminApis.GetDataHomeHandler)
	router.NotFoundHandler = http.HandlerFunc(helperApis.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(helperApis.NotAllowHandler)
	err := http.ListenAndServe(":9000", handlers.CORS(headers, methods, origins)(router))
	if err != nil {
		panic(err)
	}
}
