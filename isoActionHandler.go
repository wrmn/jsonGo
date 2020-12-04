package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/gorilla/mux"
	"github.com/mofax/iso8583"
)

type FieldDescription struct {
	ContentType string `yaml:"ContentType"`
	MaxLen      int    `yaml:"MaxLen"`
	MinLen      int    `yaml:"MinLen"`
	LenType     string `yaml:"LenType"`
	Label       string `yaml:"Label"`
}

type Spec struct {
	fields map[int]FieldDescription
}

func (s *Spec) readFromFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	yaml.Unmarshal(content, &s.fields)
	return nil
}

func getPaymentIso(w http.ResponseWriter, r *http.Request) {

	processingCode := mux.Vars(r)["id"]
	transaction, _ := selectPayment(processingCode, dbCon)
	result := jsonToIso(transaction)

	w.Write([]byte(result))
}

func toIso(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	errorCheck(err)
	var transaction Transaction
	err = json.Unmarshal(b, &transaction)
	result := jsonToIso(transaction)
	w.Write([]byte(result))

}

func jsonToIso(transaction Transaction) string {

	logWriter("New request ISO:8583 to Json")
	logWriter("original : " + fmt.Sprint(transaction))
	iso := iso8583.NewISOStruct("spec1987.yml", false)
	for len(transaction.CardAcceptorData.CardAcceptorCity) < 13 {
		transaction.CardAcceptorData.CardAcceptorCity += " "
	}
	for len(transaction.CardAcceptorData.CardAcceptorName) < 25 {
		transaction.CardAcceptorData.CardAcceptorName += " "
	}

	cardAcceptor := transaction.CardAcceptorData.CardAcceptorName +
		transaction.CardAcceptorData.CardAcceptorCity +
		transaction.CardAcceptorData.CardAcceptorCountryCode

	amount := strconv.Itoa(transaction.TotalAmount)
	something := Spec{}
	e := something.readFromFile("spec1987.yml")
	if e != nil {
		fmt.Println(e.Error())
	}
	val := map[int]string{2: transaction.Pan,
		3:  transaction.ProcessingCode,
		4:  amount,
		5:  transaction.SettlementAmount,
		6:  transaction.CardholderBillingAmount,
		7:  transaction.TransmissionDateTime,
		9:  transaction.SettlementConversionrate,
		10: transaction.CardHolderBillingConvRate,
		11: transaction.Stan,
		12: transaction.LocalTransactionTime,
		13: transaction.LocalTransactionDate,
		17: transaction.CaptureDate,
		18: transaction.CategoryCode,
		22: transaction.PointOfServiceEntryMode,
		37: transaction.Refnum,
		41: transaction.CardAcceptorData.CardAcceptorTerminalId,
		43: cardAcceptor,
		48: transaction.AdditionalData,
		49: transaction.Currency,
		50: transaction.SettlementCurrencyCode,
		51: transaction.CardHolderBillingCurrencyCode,
		57: transaction.AdditionalDataNational,
	}
	iso.AddMTI("0200")

	for id := range something.fields {
		ele := something.fields[id]
		if ele.LenType == "fixed" {
			if id == 4 {
				for len(val[id]) < ele.MaxLen {
					val[id] = "0" + val[id]
				}
			} else {
				for len(val[id]) < ele.MaxLen {
					val[id] = val[id] + " "
				}
				if ele.LenType == "fixed" {
					logWriter(fmt.Sprintf("[%d] length %d = %s", id, ele.MaxLen, val[id]))
				} else {
					logWriter(fmt.Sprintf("[%d] length %d = %s", id, len(val[id]), val[id]))
				}
			}
			if len(val[id]) > ele.MaxLen {
				val[id] = val[id][:ele.MaxLen]
			}
		}
		iso.AddField(int64(id), val[id])
	}

	result, _ := iso.ToString()
	lnth := strconv.Itoa(len(result))
	for len(lnth) < 4 {
		lnth = "0" + lnth
	}

	mti := result[:4]
	res := result[4:20]
	ele := result[20:]
	bitmap, _ := iso8583.HexToBitmapArray(res)
	logWriter("Full message	: " + lnth + result)
	logWriter("Length		: " + lnth)
	logWriter("Msg Only		: " + result)
	logWriter("MTI			: " + mti)
	logWriter("Hexmap		: " + res)
	logWriter("Bitmap		: " + fmt.Sprintf("%d", bitmap))
	logWriter("Element		: " + ele)
	return lnth + result

}

func toJson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	errorCheck(err)

	req := string(b)
	something := Spec{}
	nice := iso8583.NewISOStruct("spec1987.yml", false)
	e := something.readFromFile("spec1987.yml")

	if e != nil {
		fmt.Println(e.Error())
	}

	mti := req[4:8]
	res := req[8:24]
	ele := req[24:]
	bitmap, _ := iso8583.HexToBitmapArray(res)

	logWriter("New request ISO:8583 to Json")
	logWriter("Full message	: " + req)
	logWriter("Length		: " + req[:4])
	logWriter("Msg Only		: " + req[4:])
	logWriter("MTI			: " + mti)
	logWriter("Hexmap		: " + res)
	logWriter("Bitmap		: " + fmt.Sprintf("%d", bitmap))
	logWriter("Element		: " + ele)

	tlen := len(ele)
	mark := 0

	nice.AddMTI(mti)
	nice.Bitmap = bitmap
	for idx := range bitmap {
		if bitmap[idx] == 1 {
			element := something.fields[idx+1]
			len := element.MaxLen
			if element.LenType == "llvar" {
				clen, _ := strconv.Atoi(ele[mark : mark+2])
				msg := fmt.Sprintf("[%d] length %d = %s", idx, clen, ele[mark+2:mark+clen+2])
				logWriter(msg)
				nice.AddField(int64(idx+1), ele[mark+2:mark+clen+2])
				tlen -= clen + 2
				mark += clen + 2
			} else if element.LenType == "lllvar" {
				clen, _ := strconv.Atoi(ele[mark : mark+3])
				msg := fmt.Sprintf("[%d] length %d =  %s", idx, clen, ele[mark+3:mark+clen+3])
				logWriter(msg)
				nice.AddField(int64(idx+1), ele[mark+3:mark+clen+3])
				tlen -= clen + 3
				mark += clen + 3
			} else {
				msg := fmt.Sprintf("[%d] length %d = %s", idx, len, ele[mark:mark+len])
				logWriter(msg)
				nice.AddField(int64(idx+1), ele[mark:mark+len])
				tlen -= len
				mark += len
			}
		}
	}
	elm := nice.Elements.GetElements()

	amountTotal, _ := strconv.Atoi(elm[4])

	payment := PaymentResponse{}
	payment.TransactionData.Pan = elm[2]
	payment.TransactionData.ProcessingCode = elm[3]
	payment.TransactionData.TotalAmount = amountTotal
	payment.TransactionData.TransmissionDateTime = elm[7]
	payment.TransactionData.LocalTransactionTime = elm[12]
	payment.TransactionData.LocalTransactionDate = elm[13]
	payment.TransactionData.CaptureDate = elm[17]
	payment.TransactionData.AdditionalData = elm[48]
	payment.TransactionData.Stan = elm[11]
	payment.TransactionData.Refnum = elm[37]
	payment.TransactionData.Currency = elm[49]
	payment.TransactionData.CategoryCode = elm[18]
	payment.TransactionData.SettlementAmount = elm[5]
	payment.TransactionData.CardholderBillingAmount = elm[6]
	payment.TransactionData.SettlementConversionrate = elm[9]
	payment.TransactionData.CardHolderBillingConvRate = elm[10]
	payment.TransactionData.PointOfServiceEntryMode = elm[22]
	payment.TransactionData.SettlementCurrencyCode = elm[50]
	payment.TransactionData.CardHolderBillingCurrencyCode = elm[51]
	payment.TransactionData.AdditionalDataNational = elm[57]
	payment.TransactionData.CardAcceptorData.CardAcceptorTerminalId = elm[41]
	payment.TransactionData.CardAcceptorData.CardAcceptorName = elm[43][:24]
	payment.TransactionData.CardAcceptorData.CardAcceptorCity = elm[43][25:38]
	payment.TransactionData.CardAcceptorData.CardAcceptorCountryCode = elm[43][38:40]
	payment.ResponseStatus.ResponseCode = 200
	payment.ResponseStatus.ResponseDescription = "success"
	//fmt.Print(payment)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func logWriter(msg string) {
	log, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer log.Close()

	dt := time.Now()

	_, err := log.Write([]byte(dt.Format("01-02-2006 15:04:05 ") + msg + "\n"))
	errorCheck(err)
}
