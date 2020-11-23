package main

import "database/sql"

// query for select all payments
// todo
// add limit
func selectPayments(db *sql.DB) ([]Transaction, error) {
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
		payment.CardAcceptorData = selectAcceptor(payment.ProcessingCode, db)
		payments = append(payments, payment)
	}
	return payments, e
}

// query select payment based on procCode
func selectPayment(procCode string, db *sql.DB) (Transaction, error) {
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
		payment.CardAcceptorData = selectAcceptor(procCode, db)
	}
	return payment, e
}

//query for make acceptor has it own field on response json
func selectAcceptor(procCode string, db *sql.DB) CardAcceptorData {
	acceptor = CardAcceptorData{}
	rowsAcceptor, e := db.Query(acceptorQuery, procCode)
	errorCheck(e)
	for rowsAcceptor.Next() {
		e = rowsAcceptor.Scan(
			&acceptor.CardAcceptorTerminalId,
			&acceptor.CardAcceptorName,
			&acceptor.CardAcceptorCity,
			&acceptor.CardAcceptorCountryCode,
		)
		errorCheck(e)
	}
	return acceptor
}

//query for insert payment
func insertPayment(data Transaction, db *sql.DB) (string, error) {
	stmt, e := db.Prepare(insertQuery)
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
		data.CardAcceptorData.CardAcceptorTerminalId,
		data.CardAcceptorData.CardAcceptorName,
		data.CardAcceptorData.CardAcceptorCity,
		data.CardAcceptorData.CardAcceptorCountryCode,
		data.SettlementCurrencyCode,
		data.CardHolderBillingCurrencyCode,
		data.AdditionalDataNational,
	)
	errorCheck(e)
	stmt.Close()
	return data.ProcessingCode, e
}

// query for update payment
func putPayment(updateQuery string, db *sql.DB) error {
	_, e := db.Query(updateQuery)
	return e
}

// query for select from proc code with less time than select all for check update and delete
// if data does not exist on this query the update or delete can't run
// and return response "data does not exist"
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

//query for delete payment data
func dropPayment(procCode string, db *sql.DB) error {
	_, e := db.Query(delTransactionQuery, procCode)
	return e
}
