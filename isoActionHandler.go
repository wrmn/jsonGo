package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mofax/iso8583"
)

func getPaymentIso(w http.ResponseWriter, r *http.Request) {
	iso := iso8583.NewISOStruct("spec1987.yml", false)

	processingCode := mux.Vars(r)["id"]
	transaction, _ := selectPayment(processingCode, dbCon)

	cardAcceptor := transaction.CardAcceptorData.CardAcceptorName +
		transaction.CardAcceptorData.CardAcceptorCity +
		transaction.CardAcceptorData.CardAcceptorCountryCode

	iso.AddMTI("0200")
	iso.AddField(2, transaction.Pan)
	iso.AddField(3, transaction.ProcessingCode)
	iso.AddField(4, strconv.Itoa(transaction.TotalAmount))
	iso.AddField(5, transaction.SettlementAmount)
	iso.AddField(6, transaction.CardholderBillingAmount)
	iso.AddField(7, transaction.TransmissionDateTime)
	iso.AddField(9, transaction.SettlementConversionrate)
	iso.AddField(10, transaction.CardHolderBillingConvRate)
	iso.AddField(11, transaction.Stan)
	iso.AddField(12, transaction.LocalTransactionTime)
	iso.AddField(13, transaction.LocalTransactionDate)
	iso.AddField(17, transaction.CaptureDate)
	iso.AddField(18, transaction.CategoryCode)
	iso.AddField(22, transaction.PointOfServiceEntryMode)
	iso.AddField(37, transaction.Refnum)
	iso.AddField(41, transaction.CardAcceptorData.CardAcceptorTerminalId)
	iso.AddField(43, cardAcceptor)
	iso.AddField(48, transaction.AdditionalData)
	iso.AddField(49, transaction.Currency)
	iso.AddField(50, transaction.SettlementCurrencyCode)
	iso.AddField(51, transaction.CardHolderBillingCurrencyCode)
	iso.AddField(57, transaction.AdditionalDataNational)

	result, _ := iso.ToString()

	w.Write([]byte(result + "|"))
}

func toJson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	errorCheck(err)
	req := string(b)
	//iso := iso8583.NewISOStruct("spec1987.yml", false)
	res := req[4:20]
	fmt.Println(iso8583.HexToBitmapArray(res))

}
