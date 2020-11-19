package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, time.Now().Format("15:04:05"))
}

func getPayments(w http.ResponseWriter, r *http.Request) {
	response := PaymentsResponse{}
	w.Header().Set("Content-Type", "application/json")
	err := pingDb(dbCon)

	if err != nil {
		response.ResponseCode = 500
		response.ResponseDescription = serverError
	} else {
		w.WriteHeader(200)
		response.ResponseCode = 200
		response.ResponseDescription = "success"
		payments := selectPayments(dbCon)

		response.Response = payments
	}

	json.NewEncoder(w).Encode(response)
}
func getPayment(w http.ResponseWriter, r *http.Request) {
	response := PaymentResponse{}
	w.Header().Set("Content-Type", "application/json")
	err := pingDb(dbCon)

	pan := mux.Vars(r)["id"]

	if err != nil {
		response.ResponseCode = 500
		response.ResponseDescription = serverError
	} else {
		w.WriteHeader(200)
		response.ResponseCode = 200
		response.ResponseDescription = "success"
		payments := selectPayment(pan, dbCon)

		response.Response = payments
	}

	json.NewEncoder(w).Encode(response)
}

func createPayment(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	response := InsPaymentResponse{}

	if !(errorCheck(err)) {
		return
	}

	var trs Transaction
	err = json.Unmarshal(b, &trs)

	if !(errorCheck(err)) {
		http.Error(w, err.Error(), 500)
		return
	}

	if !(errorCheck(err)) {
		http.Error(w, err.Error(), 500)
		return
	}

	processingCode, err := insertPayment(trs, dbCon)

	if err != nil {
		response.ResponseCode = 500
		response.ResponseDescription = err.Error()
	} else {
		w.WriteHeader(200)
		response.ResponseCode = 200
		response.ResponseDescription = "success"
		response.Response = processingCode
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updatePayment(w http.ResponseWriter, r *http.Request) {
}

func deletePayment(w http.ResponseWriter, r *http.Request) {
	response := DelPaymentResponse{}
	w.Header().Set("Content-Type", "application/json")
	err := pingDb(dbCon)

	procCode := mux.Vars(r)["id"]

	if err != nil {
		response.ResponseCode = 500
		response.ResponseDescription = serverError
	} else {
		payment := selectPayment(procCode, dbCon)
		if payment.ProcessingCode == "" {
			w.WriteHeader(400)
			response.ResponseDescription = "data not exist"
			response.ResponseCode = 400
		} else {
			dropPayment(procCode, dbCon)
			w.WriteHeader(200)
			response.ResponseCode = 200
			response.ResponseDescription = "Deleted"
		}
	}

	json.NewEncoder(w).Encode(response)
}
