package payments

import (
	"errors"
	"fmt"

	"github.com/razorpay/razorpay-go"
)

func NewRzpOrder(pi *PaymentInfo) (po *PaymentOrder, err error) {
	p, ok := registeredProvider[RZP]
	if !ok {
		return nil, errors.New("Provider not registred")
	}
	rzp := razorpay.NewClient(p.Key, p.Secret)

	data := map[string]interface{}{
		"amount":   pi.Amount * 100,
		"currency": pi.Currency,
		"receipt":  pi.OrderID,
	}

	body, err := rzp.Order.Create(data, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	providerOrderId := body["id"].(string)
	fmt.Println(providerOrderId)
	return &PaymentOrder{p.Value, providerOrderId, "no-token", p.Key}, nil
}
