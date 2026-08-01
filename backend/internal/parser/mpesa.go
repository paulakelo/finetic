package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MpesaTransaction holds the strictly typed data extracted from a raw SMS.
type MpesaTransaction struct {
	ReceiptNumber string
	Amount        int64  // Stored in the lowest denomination (cents)
	Type          string // "EXPENSE" or "INCOME"
	Counterparty  string // The sender or recipient's name
	Date          time.Time
}

var (
	// Matches the 10-character alphanumeric receipt at the start of the SMS
	receiptRegex = regexp.MustCompile(`^([A-Z0-9]{10}\s+Confirmed\.?)`)

	// Matches the currency amount, capturing the numbers and commas
	amountRegex = regexp.MustCompile(`Ksh([0-9,.]+)`)

	// Matches who the money was sent to (Expense)
	sentRegex = regexp.MustCompile(`sent to\s+(.+?)\s+on\s+`)
	// Matches who the money was paid to (e.g., Paybills/Tills - Expense)
	paidRegex = regexp.MustCompile(`paid to \s+(.+?)\s+on\s+`)
	// Matches who the money was received from (Income)
	receivedRegex = regexp.MustCompile(`received\s+from\s+(.+?)\s+on\s+`)

	// Matches the Safaricom date format: "on 12/9/26 at 10:14 AM"
	dateRegex = regexp.MustCompile(`on\s+(\d{1,2})/(\d{1,2})/(\d{2})\s+at\s+(\d{1,2}):(\d{2})\s+(AM|PM)`)
)

// ParseMpesaMessage takes a raw Safaricom SMS string and extracts the financial data.
// It returns a pointer to an MpesaTransaction, or an error if the SMS format is unrecognized.
func ParseMpesaMessage(rawSMS string) (*MpesaTransaction, error) {
	// 1. Extract Receipt Number
	receiptMatch := receiptRegex.FindStringSubmatch(rawSMS)
	if len(receiptMatch) < 2 {
		return nil, errors.New("failed to parse receipt number: invalid SMS format")
	}
	receipt := receiptMatch[1]

	// 2. Extract Amount
	amountMatch := amountRegex.FindStringSubmatch(rawSMS)
	if len(amountMatch) < 2 {
		return nil, fmt.Errorf("failed to parse amount for receipt %s", receipt)
	}

	// Convert "1,500.50" -> "1500.50" -> float -> int64 (cents)
	cleanAmountStr := strings.ReplaceAll(amountMatch[1], ",", "")
	parsedFloat, err := strconv.ParseFloat(cleanAmountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert amount string to number: %v", err)
	}
	amountInCents := int64(parsedFloat * 100)

	// 3. Determine Transaction Type and Counterparty
	var txType, counterparty string

	if sentMatch := sentRegex.FindStringSubmatch(rawSMS); len(sentMatch) >= 2 {
		txType = "EXPENSE"
		counterparty = sentMatch[1]
	} else if receivedMatch := receivedRegex.FindStringSubmatch(rawSMS); len(receivedMatch) >= 2 {
		txType = "INCOME"
		counterparty = receivedMatch[1]
	} else {
		return nil, fmt.Errorf("failed to determine transaction direction for receipt %s", receipt)
	}

	// 4. Extract and Parse the Date
	dateMatch := dateRegex.FindStringSubmatch(rawSMS)
	if len(dateMatch) < 3 {
		return nil, fmt.Errorf("failed to parse date for receipt %s", receipt)
	}

	// Combine the extracted date and time into a single string for the Go time parser
	dateStr := fmt.Sprintf("%s %s", dateMatch[1], dateMatch[2])

	// We map the Safaricom format "D/M/YY 3:04 PM" to the Go reference layout.
	location, _ := time.LoadLocation("Africa/Nairobi")
	parsedDate, err := time.ParseInLocation("2/1/06 3:04 PM", dateStr, location)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time string %v", err)
	}

	return &MpesaTransaction{
		ReceiptNumber: receipt,
		Amount:        amountInCents,
		Type:          txType,
		Counterparty:  strings.TrimSpace(counterparty),
		Date:          parsedDate,
	}, nil
}
