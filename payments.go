package payments

import (
	"errors"
	"fmt"
	"os"
)

const RZP = "rzp"
const CSH = "csh"

type Provider struct {
	Value  string
	Url    string
	Key    string
	Secret string
}

// check if the provider is registred or not
var registeredProvider = make(map[string]Provider, 2)

func RegisterProvider(p Provider) {
	registeredProvider[p.Value] = p
}

type PaymentInfo struct {
	Amount         int
	Currency       string
	OrderID        string
	CustomerID     string
	CustomerEmail  string
	CustomerMobile string
}

type PaymentOrder struct {
	Provider   string
	PaymentUrl string
	Token      string
	Key        string
}

func NewOrder(pi PaymentInfo) (po *PaymentOrder, err error) {

	p := os.Getenv("DEFAULT_PAYMENT_PROVIDER")
	if p == "" {
		fmt.Println("Defalut payment provider not set, falling back csh")
		p = CSH
	}
	switch p {
	case CSH:
		return newCSHOrder(&pi)
	case RZP:
		return NewRzpOrder(&pi)
	}
	return nil, errors.New("Provider not supported")
}
