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
	fmt.Println("ok")
	err := pingDb(dbCon)
	fmt.Println(err)
}

// handler action from route with request get all payments
// todo
// return error from query
// get limit
func getPayments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := PaymentsResponse{}
	fresponse := FailPaymentsResponse{}
	err := pingDb(dbCon)

	if err != nil {
		w.WriteHeader(500)
		fresponse.ResponseStatus.ReasonCode, fresponse.ResponseStatus.ResponseDescription = 500, serverError
		json.NewEncoder(w).Encode(fresponse)
	} else {
		payments, err := selectPayments(dbCon)
		if err != nil {
			w.WriteHeader(500)
			fresponse.ResponseStatus.ReasonCode, fresponse.ResponseStatus.ResponseDescription = 500, err.Error()
			json.NewEncoder(w).Encode(fresponse)
		} else {
			w.WriteHeader(200)
			response.ResponseStatus.ResponseCode, response.ResponseStatus.ResponseDescription = 200, "success"
			response.TransactionData = payments
			json.NewEncoder(w).Encode(response)
		}
	}

}

// handler action from route with request get payment with
// procid
// todo
// return error from query
func getPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := PaymentResponse{}
	fresponse := FailPaymentsResponse{}
	err := pingDb(dbCon)
	processingCode := mux.Vars(r)["id"]

	if err != nil {
		w.WriteHeader(500)
		fresponse.ResponseStatus.ReasonCode, fresponse.ResponseStatus.ResponseDescription = 500, serverError
		json.NewEncoder(w).Encode(fresponse)
	} else {
		payment, err := selectPayment(processingCode, dbCon)
		if err != nil {
			w.WriteHeader(500)
			fresponse.ResponseStatus.ReasonCode, fresponse.ResponseStatus.ResponseDescription = 500, err.Error()
			json.NewEncoder(w).Encode(fresponse)
		} else if payment.ProcessingCode == "" {
			w.WriteHeader(404)
			fresponse.ResponseStatus.ReasonCode, fresponse.ResponseStatus.ResponseDescription = 404, "data not found"
			json.NewEncoder(w).Encode(fresponse)
		} else {
			w.WriteHeader(200)
			response.ResponseStatus.ResponseCode, response.ResponseStatus.ResponseDescription = 200, "success"
			response.TransactionData = payment
			json.NewEncoder(w).Encode(response)
		}
	}
}

//handler action from route with request post with json body required
//todo
//check if json in correct format
//return error from query
func createPayment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	errorCheck(err)
	response := PaymentResponse{}
	fresponse := FailPaymentsResponse{}

	err = pingDb(dbCon)

	if err != nil {
		response.ResponseStatus.ResponseCode, response.ResponseStatus.ResponseDescription = 500, serverError
	} else {
		var trs Transaction
		err = json.Unmarshal(b, &trs)

		errorCheck(err)

		if checkExistence(trs.ProcessingCode, dbCon) {
			processingCode, err := insertPayment(trs, dbCon)

			if err != nil {
				w.WriteHeader(500)
				fresponse.ResponseStatus.ResponseCode, fresponse.ResponseStatus.ResponseDescription = 500, err.Error()
				json.NewEncoder(w).Encode(fresponse)
			} else {
				w.WriteHeader(200)
				response.ResponseStatus.ResponseCode, response.ResponseStatus.ResponseDescription = 200, "success"
				response.TransactionData, _ = selectPayment(processingCode, dbCon)
				json.NewEncoder(w).Encode(response)
			}
		} else {
			w.WriteHeader(403)
			fresponse.ResponseStatus.ResponseCode, fresponse.ResponseStatus.ResponseDescription = 403, "duplicate processingCode"
			json.NewEncoder(w).Encode(fresponse)
		}
	}

}

//handler action from route with request put with json body required
//todo
//check if json in correct format
//return error from query
func updatePayment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)

	w.Header().Set("Content-Type", "application/json")
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
	response := PaymentResponse{}
	fresponse := FailPaymentsResponse{}

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		typeOfF := fmt.Sprintf(typeOfT.Field(i).Type.Name())

		if typeOfF == "CardAcceptorData" {
			for j := 0; j < t.NumField(); j++ {
				g := t.Field(j)
				val := g.Interface()
				if val != "" && val != "0" {
					as = fmt.Sprintf("%s='%v'", typeOfU.Field(j).Name, g.Interface())
					canQue = append(canQue, as)
				}
			}
		} else {
			val := f.Interface()
			if val != "" && val != 0 {
				as = fmt.Sprintf("%s='%v'", typeOfT.Field(i).Name, f.Interface())
				canQue = append(canQue, as)
			}
		}
	}

	preQue := strings.Join(canQue, ", ")
	exeQue := fmt.Sprintln("UPDATE transaction SET " + preQue + " where processingCode =" + procCode)
	payment, _ := selectPayment(procCode, dbCon)
	err = putPayment(exeQue, dbCon)

	if err != nil {
		w.WriteHeader(500)
		fresponse.ResponseStatus.ReasonCode, fresponse.ResponseStatus.ResponseDescription = 500, err.Error()
		json.NewEncoder(w).Encode(fresponse)
	} else {
		if payment.ProcessingCode == "" {
			w.WriteHeader(400)
			fresponse.ResponseStatus.ReasonCode, fresponse.ResponseStatus.ResponseDescription = 400, "data not exist"
			json.NewEncoder(w).Encode(fresponse)
		} else {
			w.WriteHeader(200)
			response.ResponseStatus.ResponseCode, response.ResponseStatus.ResponseDescription = 200, "updated"
			response.TransactionData, _ = selectPayment(procCode, dbCon)
			json.NewEncoder(w).Encode(response)
		}
	}

}

//handler action from route with request delte based on proc code that send as param
//todo
//return error from query
func deletePayment(w http.ResponseWriter, r *http.Request) {
	response := DelPaymentResponse{}
	w.Header().Set("Content-Type", "application/json")
	err := pingDb(dbCon)

	procCode := mux.Vars(r)["id"]
	if err != nil {
		w.WriteHeader(500)
		response.ResponseCode, response.ResponseDescription = 500, serverError
	} else {
		payment, _ := selectPayment(procCode, dbCon)
		if payment.ProcessingCode == "" {
			w.WriteHeader(404)
			response.ResponseCode, response.ResponseDescription = 404, "data not exist"
		} else {
			dropPayment(procCode, dbCon)
			w.WriteHeader(200)
			response.ResponseCode, response.ResponseDescription = 200, "Deleted"
		}
	}

	json.NewEncoder(w).Encode(response)
}
