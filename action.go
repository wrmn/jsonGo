package main

import (
	"database/sql"
	"fmt"
)

func selectPayments(db *sql.DB) []Transaction {
	payments := []Transaction{}
	rowsPayment, e := db.Query(transactionsQuery)
	errorCheck(e)
	for rowsPayment.Next() {
		payment = Transaction{}
		e = rowsPayment.Scan(
			&payment.Pan,
			&payment.ProcessingCode,
			&payment.TotalAmount,
			&payment.AcquirerId,
			&payment.IssuerId,
			&payment.TransmissionDateTime,
			&payment.LocalTransactionTime,
			&payment.LocalTransactionDate,
			&payment.CaptureDate,
			&payment.AdditionalData,
			&payment.Stan,
			&payment.Refnum,
			&payment.Currency,
			&payment.TerminalId,
			&payment.AccountFrom,
			&payment.AccountTo,
			&payment.CategoryCode,
			&payment.SettlementAmount,
			&payment.CardholderBillingAmount,
			&payment.SettlementConversionrate,
			&payment.CardHolderBillingConvRate,
			&payment.PointOfServiceEntryMode,
			&payment.SettlementCurrencyCode,
			&payment.CardHolderBillingCurrencyCode,
			&payment.AdditionalDataNational,
		)
		errorCheck(e)
		payment.CardAcceptorData = selectAcceptor(payment.ProcessingCode, db)
		payments = append(payments, payment)
	}
	return payments
}

func selectPayment(procCode string, db *sql.DB) Transaction {
	payment = Transaction{}
	rowPayment, e := db.Query(transactionQuery, procCode)
	errorCheck(e)
	for rowPayment.Next() {
		e = rowPayment.Scan(
			&payment.Pan,
			&payment.ProcessingCode,
			&payment.TotalAmount,
			&payment.AcquirerId,
			&payment.IssuerId,
			&payment.TransmissionDateTime,
			&payment.LocalTransactionTime,
			&payment.LocalTransactionDate,
			&payment.CaptureDate,
			&payment.AdditionalData,
			&payment.Stan,
			&payment.Refnum,
			&payment.Currency,
			&payment.TerminalId,
			&payment.AccountFrom,
			&payment.AccountTo,
			&payment.CategoryCode,
			&payment.SettlementAmount,
			&payment.CardholderBillingAmount,
			&payment.SettlementConversionrate,
			&payment.CardHolderBillingConvRate,
			&payment.PointOfServiceEntryMode,
			&payment.SettlementCurrencyCode,
			&payment.CardHolderBillingCurrencyCode,
			&payment.AdditionalDataNational,
		)
		errorCheck(e)
		payment.CardAcceptorData = selectAcceptor(procCode, db)
	}
	return payment
}

func selectAcceptor(procCode string, db *sql.DB) CardAcceptedData {
	acceptor = CardAcceptedData{}
	rowsAcceptor, e := db.Query(acceptorQuery, procCode)
	errorCheck(e)
	for rowsAcceptor.Next() {
		e = rowsAcceptor.Scan(
			&acceptor.TerminalId,
			&acceptor.Name,
			&acceptor.City,
			&acceptor.CountryCode,
		)
		errorCheck(e)
	}
	return acceptor
}

func insertPayment(data Transaction, db *sql.DB) (string, error) {
	stmt, e := db.Prepare(insertQuery)
	fmt.Println(data.CardAcceptorData.TerminalId)
	errorCheck(e)
	_, e = stmt.Exec(
		data.Pan,
		data.ProcessingCode,
		data.TotalAmount,
		data.AcquirerId,
		data.IssuerId,
		data.TransmissionDateTime,
		data.LocalTransactionTime,
		data.LocalTransactionDate,
		data.CaptureDate,
		data.AdditionalData,
		data.Stan,
		data.Refnum,
		data.Currency,
		data.TerminalId,
		data.AccountFrom,
		data.AccountTo,
		data.CategoryCode,
		data.SettlementAmount,
		data.CardholderBillingAmount,
		data.SettlementConversionrate,
		data.CardHolderBillingConvRate,
		data.PointOfServiceEntryMode,
		data.CardAcceptorData.TerminalId,
		data.CardAcceptorData.Name,
		data.CardAcceptorData.City,
		data.CardAcceptorData.CountryCode,
		data.SettlementCurrencyCode,
		data.CardHolderBillingCurrencyCode,
		data.AdditionalDataNational,
	)
	errorCheck(e)

	return data.ProcessingCode, e

}

func checkExistence(procCode string, db *sql.DB) bool {
	var total int
	rowCheck, e := db.Query(checkQuery, procCode)
	errorCheck(e)
	for rowCheck.Next() {
		e = rowCheck.Scan(&total)
		errorCheck(e)
	}
	if total != 0 {
		return false
	}
	return true
}

func dropPayment(procCode string, db *sql.DB) error {
	_, e := db.Query(delTransactionQuery, procCode)
	return e
}
