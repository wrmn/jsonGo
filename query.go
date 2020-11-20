package main

const (
	transactionQuery = `
		SELECT
			pan,
			processingCode,
			totalAmount,
			acquirerId,
			issuerId,
			transmissionDateTime,
			localTransactionTime,
			localTransactionDate,
			captureDate,
			additionalData,
			stan,
			refnum,
			currency,
			terminalId,
			accountFrom,
			accountTo,
			categoryCode,
			settlementAmount,
			cardholderBillingAmount,
			settlementConversionrate,
			cardHolderBillingConvRate,
			pointOfServiceEntryMode,
			settlementCurrencyCode,
			cardHolderBillingCurrencyCode,
			additionalDataNational
		FROM
			transaction
		WHERE processingCode =?
	`
	transactionsQuery = `
		SELECT
			pan,
			processingCode,
			totalAmount,
			acquirerId,
			issuerId,
			transmissionDateTime,
			localTransactionTime,
			localTransactionDate,
			captureDate,
			additionalData,
			stan,
			refnum,
			currency,
			terminalId,
			accountFrom,
			accountTo,
			categoryCode,
			settlementAmount,
			cardholderBillingAmount,
			settlementConversionrate,
			cardHolderBillingConvRate,
			pointOfServiceEntryMode,
			settlementCurrencyCode,
			cardHolderBillingCurrencyCode,
			additionalDataNational
		FROM
			transaction
	`
	acceptorQuery = `
		SELECT
			cardAcceptorTerminalID,
			cardAcceptorName,
			cardAcceptorCity,
			cardAcceptorCountryCode
		FROM
			transaction
		WHERE processingCode = ?

	`

	insertQuery = `
		INSERT INTO transaction (
			pan, 
			processingCode, 
			totalAmount, 
			acquirerId, 
			issuerId, 
			transmissionDateTime, 
			localTransactionTime, 
			localTransactionDate, 
			captureDate, 
			additionalData, 
			stan, 
			refnum, 
			currency, 
			terminalId, 
			accountFrom, 
			accountTo, 
			categoryCode, 
			settlementAmount, 
			cardholderBillingAmount, 
			settlementConversionRate, 
			cardholderBillingConvRate, 
			pointOfServiceEntryMode, 
			cardAcceptorTerminalID, 
			cardAcceptorName, 
			cardAcceptorCity, 
			cardAcceptorCountryCode, 
			settlementCurrencyCode, 
			cardholderBillingCurrencyCode, 
			additionalDataNational) 
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	delTransactionQuery = `
		DELETE FROM transaction WHERE processingCode = ?		
	`

	checkQuery = `
		SELECT count(*) FROM transaction WHERE processingCode = ?
	`
)
