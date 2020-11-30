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

	field := []int64{2, 3, 4, 5, 6, 7, 9, 10, 11, 12, 13, 17, 18, 22, 37, 41, 43, 48, 49, 50, 51, 57}
	val := []string{transaction.Pan,
		transaction.ProcessingCode,
		strconv.Itoa(transaction.TotalAmount),
		transaction.SettlementAmount,
		transaction.CardholderBillingAmount,
		transaction.TransmissionDateTime,
		transaction.SettlementConversionrate,
		transaction.CardHolderBillingConvRate,
		transaction.Stan,
		transaction.LocalTransactionTime,
		transaction.LocalTransactionDate,
		transaction.CaptureDate,
		transaction.CategoryCode,
		transaction.PointOfServiceEntryMode,
		transaction.Refnum,
		transaction.CardAcceptorData.CardAcceptorTerminalId,
		cardAcceptor,
		transaction.AdditionalData,
		transaction.Currency,
		transaction.SettlementCurrencyCode,
		transaction.CardHolderBillingCurrencyCode,
		transaction.AdditionalDataNational,
	}
	iso.AddMTI("0200")

	for idx := range field {
		if val[idx] != "" && val[idx] != "0" {
			fmt.Println(field[idx])
			fmt.Println(val[idx])
			iso.AddField(field[idx], val[idx])
		}
	}

	result, _ := iso.ToString()
	fmt.Println(iso)
	w.Write([]byte(result))
}

func toJson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	errorCheck(err)
	req := string(b)
	//iso := iso8583.NewISOStruct("spec1987.yml", false)
	nice := iso8583.NewISOStruct("spec1987.yml", false)
	mti := req[:4]
	res := req[4:20]
	ele := req[21:]
	fmt.Println(ele)
	nice.AddMTI(mti)
	bitmap, _ := iso8583.HexToBitmapArray(res)
	nice.Bitmap = bitmap
	for idx := range bitmap {
		if bitmap[idx] == 1 {
			fmt.Println("oke")
		}
	}

	fmt.Println(nice)
	/* sum, e := nice.Parse(req)*/

	//if e != nil {
	//fmt.Println(e.Error())
	//}
	/*fmt.Println(sum)*/
}
