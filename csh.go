package payments

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UnmarshalCashFreeOrder(data []byte) (CashFreeOrder, error) {
	var r CashFreeOrder
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CashFreeOrder) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CashFreeOrder struct {
	CFOrderID       int64           `json:"cf_order_id"`
	CreatedAt       string          `json:"created_at"`
	CustomerDetails CustomerDetails `json:"customer_details"`
	Entity          string          `json:"entity"`
	OrderAmount     float64         `json:"order_amount"`
	OrderCurrency   string          `json:"order_currency"`
	OrderExpiryTime string          `json:"order_expiry_time"`
	OrderID         string          `json:"order_id"`
	OrderMeta       OrderMeta       `json:"order_meta"`
	OrderNote       interface{}     `json:"order_note"`
	OrderStatus     string          `json:"order_status"`
	OrderToken      string          `json:"order_token"`
	PaymentLink     string          `json:"payment_link"`
	Payments        Payments        `json:"payments"`
	Refunds         Payments        `json:"refunds"`
	Settlements     Payments        `json:"settlements"`
}

type CustomerDetails struct {
	CustomerID    string      `json:"customer_id"`
	CustomerName  interface{} `json:"customer_name"`
	CustomerEmail string      `json:"customer_email"`
	CustomerPhone string      `json:"customer_phone"`
}

type OrderMeta struct {
	ReturnURL      string      `json:"return_url"`
	NotifyURL      string      `json:"notify_url"`
	PaymentMethods interface{} `json:"payment_methods"`
}

type Payments struct {
	URL string `json:"url"`
}

func UnmarshalCSHToken(data []byte) (CSHToken, error) {
	var r CSHToken
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CSHToken) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CSHToken struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Cftoken string `json:"cftoken"`
}

func newCSHOrder(pi *PaymentInfo) (po *PaymentOrder, err error) {
	p, ok := registeredProvider[CSH]
	if !ok {
		return nil, errors.New("Provider not registred")
	}
	// reqBody, err := json.Marshal(map[string]interface{}{
	// 	"order_id":       pi.OrderID,
	// 	"order_amount":   pi.Amount,
	// 	"order_currency": pi.Currency,
	// 	"customer_details": map[string]string{
	// 		"customer_id":    pi.CustomerID,
	// 		"customer_email": pi.CustomerEmail,
	// 		"customer_phone": pi.CustomerMobile,
	// 	},
	// 	"order_meta": map[string]string{
	// 		"notify_url": os.Getenv("HOST") + "/payments/csh/webhook",
	// 	},
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }

	reqBody, err := json.Marshal(map[string]interface{}{
		"orderId":       pi.OrderID,
		"orderAmount":   pi.Amount,
		"orderCurrency": pi.Currency,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", p.Url+"/cftoken/order", bytes.NewBuffer(reqBody))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("x-client-id", p.Key)
	req.Header.Add("x-client-secret", p.Secret)
	req.Header.Add("x-api-version", "2021-05-21")
	req.Header.Add("Content-Type", `application/json`)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	body, berr := ioutil.ReadAll(res.Body)
	if berr != nil {
		fmt.Println(berr)
		return nil, berr
	}
	if res.StatusCode != 200 {
		fmt.Println(string(body))
		return nil, errors.New("error creating cashfree payment order")
	}

	// cashFreeOrder, err := UnmarshalCashFreeOrder(body)
	cSHToken, err := UnmarshalCSHToken(body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &PaymentOrder{p.Value, pi.OrderID, cSHToken.Cftoken, p.Key}, nil
}
