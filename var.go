package main

import "database/sql"

var (
	dbCon    *sql.DB
	payment  = Transaction{}
	acceptor = CardAcceptedData{}
)

const (
	serverError = "Internal server error, database disconnected. Contact Administrator"
)

type CardAcceptedData struct {
	TerminalId  string `json:"cardAcceptorTerminalID"`
	Name        string `json:"cardAcceptorName"`
	City        string `json:"cardAcceptorCity"`
	CountryCode string `json:"cardAcceptorCountryCode"`
}

type Transaction struct {
	Pan                           string           `json:"pan"`
	ProcessingCode                string           `json:"processingCode"`
	TotalAmount                   int              `json:"totalAmount"`
	AcquirerId                    string           `json:"acquirerId"`
	IssuerId                      string           `json:"issuerId"`
	TransmissionDateTime          string           `json:"transmissionDateTime"`
	LocalTransactionTime          string           `json:"localTransactionTime"`
	LocalTransactionDate          string           `json:"localTransactionDate"`
	CaptureDate                   string           `json:"captureDate"`
	AdditionalData                string           `json:"additionalData"`
	Stan                          string           `json:"stan"`
	Refnum                        string           `json:"refnum"`
	Currency                      string           `json:"currency"`
	TerminalId                    string           `json:"terminalId"`
	AccountFrom                   string           `json:"accountFrom"`
	AccountTo                     string           `json:"accountTo"`
	CategoryCode                  string           `json:"categoryCode"`
	SettlementAmount              string           `json:"settlementAmount"`
	CardholderBillingAmount       string           `json:"cardholderBillingAmount"`
	SettlementConversionrate      string           `json:"settlementConversionrate"`
	CardHolderBillingConvRate     string           `json:"cardHolderBillingConvRate"`
	PointOfServiceEntryMode       string           `json:"pointOfServiceEntryMode"`
	CardAcceptorData              CardAcceptedData `json:"cardAcceptorData"`
	SettlementCurrencyCode        string           `json:"settlementCurrencyCode"`
	CardHolderBillingCurrencyCode string           `json:"cardHolderBillingCurrencyCode"`
	AdditionalDataNational        string           `json:"additionalDataNational"`
}

type Response struct {
	ResponseCode        int    `json:"responseCode"`
	ReasonCode          int    `json:"reasonCode"`
	ResponseDescription string `json:"responseDescription"`
}

type PaymentsResponse struct {
	TransactionData []Transaction `json:"transactionData"`
	ResponseStatus  Response      `json:"responseStatus"`
}

type PaymentResponse struct {
	TransactionData Transaction `json:"response"`
	ResponseStatus  Response    `json:"responseStatus"`
}

type InsPaymentResponse struct {
	TransactionData Transaction `json:"processingCode"`
	ResponseStatus  Response    `json"responseStatus":`
}

type DelPaymentResponse struct {
	ResponseCode        int    `json:"responseCode"`
	ReasonCode          int    `json:"reasonCode"`
	ResponseDescription string `json:"responseDescription"`
}
