package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
}

func getPayments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := PaymentsResponse{}
	err := pingDb(dbCon)

	if err != nil {
		response.ResponseStatus.ReasonCode = 500
		response.ResponseStatus.ResponseDescription = serverError
	} else {
		w.WriteHeader(200)
		payments := selectPayments(dbCon)

		response.ResponseStatus.ResponseCode = 200
		response.ResponseStatus.ResponseDescription = "success"
		response.TransactionData = payments
	}

	json.NewEncoder(w).Encode(response)
}

func getPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := PaymentResponse{}
	err := pingDb(dbCon)
	processingCode := mux.Vars(r)["id"]

	if err != nil {
		response.ResponseStatus.ResponseCode = 500
		response.ResponseStatus.ResponseDescription = serverError
	} else {
		w.WriteHeader(200)
		payments := selectPayment(processingCode, dbCon)

		response.ResponseStatus.ResponseCode = 200
		response.ResponseStatus.ResponseDescription = "success"
		response.TransactionData = payments
	}
	json.NewEncoder(w).Encode(response)
}

func createPayment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	response := InsPaymentResponse{}

	errorCheck(err)
	var trs Transaction
	err = json.Unmarshal(b, &trs)

	errorCheck(err)

	if checkExistence(trs.ProcessingCode, dbCon) {
		processingCode, err := insertPayment(trs, dbCon)

		if err != nil {
			response.ResponseStatus.ReasonCode = 500
			response.ResponseStatus.ResponseDescription = err.Error()
		} else {
			response.ResponseStatus.ReasonCode = 200
			response.ResponseStatus.ResponseDescription = "success"
			response.TransactionData = selectPayment(processingCode, dbCon)
		}
	} else {
		response.ResponseStatus.ResponseCode = 500
		response.ResponseStatus.ResponseDescription = "duplicate processingCode"
		response.TransactionData = trs
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Print(w.Header().Values("Content-Type"))
	json.NewEncoder(w).Encode(response)
}

func updatePayment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)

	errorCheck(err)

	var trs Transaction
	var canQue []string
	var as string
	err = json.Unmarshal(b, &trs)

	s := reflect.ValueOf(&trs).Elem()
	t := reflect.ValueOf(&trs.CardAcceptorData).Elem()
	typeOfT := s.Type()
	typeOfU := t.Type()
	procCode := mux.Vars(r)["id"]

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		typeOfF := fmt.Sprintf(typeOfT.Field(i).Type.Name())

		if typeOfF == "CardAcceptorData" {
			for j := 0; j < t.NumField(); j++ {
				g := t.Field(j)
				as = fmt.Sprintf("%s='%v'", typeOfU.Field(j).Name, g.Interface())
				canQue = append(canQue, as)
			}
		} else {
			val := f.Interface()
			if val != "" {
				as = fmt.Sprintf("%s='%v'", typeOfT.Field(i).Name, f.Interface())
				canQue = append(canQue, as)
			}
		}
	}
	preQue := strings.Join(canQue, ", ")
	exeQue := fmt.Sprintln("UPDATE transaction SET " + preQue + " where processingCode =?")
	putPayment(exeQue, procCode, dbCon)
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
			response.ResponseCode = 400
			response.ResponseDescription = "data not exist"
		} else {
			dropPayment(procCode, dbCon)
			w.WriteHeader(200)
			response.ResponseCode = 200
			response.ResponseDescription = "Deleted"
		}
	}

	json.NewEncoder(w).Encode(response)
}
