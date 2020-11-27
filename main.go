package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	dbCon = initDb()
	r.HandleFunc("/", mainHandler)
	r.HandleFunc("/payment", getPayments).Methods("GET")
	r.HandleFunc("/payment", createPayment).Methods("POST")
	r.HandleFunc("/payment/{id}", updatePayment).Methods("PUT")
	r.HandleFunc("/payment/{id}", deletePayment).Methods("DELETE")
	r.HandleFunc("/payment/{id}", getPayment).Methods("GET")
	r.HandleFunc("/payment/{id}/iso8583", getPaymentIso).Methods("GET")
	r.HandleFunc("/iso8583/toJson", toJson).Methods("GET")

	http.ListenAndServe(":5050", r)
}
