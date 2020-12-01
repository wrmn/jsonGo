package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

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
	iso := iso8583.NewISOStruct("spec1987.yml", false)

	processingCode := mux.Vars(r)["id"]
	transaction, _ := selectPayment(processingCode, dbCon)

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
			for len(val[id]) < ele.MaxLen {
				val[id] = val[id] + " "
			}
		}
		iso.AddField(int64(id), val[id])
	}

	result, _ := iso.ToString()
	fmt.Println(iso)
	w.Write([]byte(result))
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

	mti := req[:4]
	res := req[4:20]
	ele := req[20:]
	tlen := len(ele)
	mark := 0

	bitmap, _ := iso8583.HexToBitmapArray(res)

	nice.AddMTI(mti)
	nice.Bitmap = bitmap
	for idx := range bitmap {
		if bitmap[idx] == 1 {
			element := something.fields[idx+1]
			len := element.MaxLen
			if element.LenType == "llvar" {
				clen, _ := strconv.Atoi(ele[mark : mark+2])
				//fmt.Println(element.Label)
				//fmt.Println(ele[mark+2 : mark+clen+2])
				nice.AddField(int64(idx+1), ele[mark+2:mark+clen+2])
				tlen -= clen + 2
				mark += clen + 2
			} else if element.LenType == "lllvar" {
				clen, _ := strconv.Atoi(ele[mark : mark+3])
				//fmt.Println(element.Label)
				//fmt.Println(ele[mark+3 : mark+clen+3])
				nice.AddField(int64(idx+1), ele[mark+3:mark+clen+3])
				tlen -= clen + 3
				mark += clen + 3
			} else {
				//fmt.Println(element.Label)
				//fmt.Println(ele[mark : mark+len])
				nice.AddField(int64(idx+1), ele[mark:mark+len])
				tlen -= len
				mark += len
			}
		}
	}
	elm := nice.Elements.GetElements()

	reg, _ := regexp.Compile("[ ]")
	amnt := reg.ReplaceAllString(elm[4], "")
	amountTotal, _ := strconv.Atoi(amnt)

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
