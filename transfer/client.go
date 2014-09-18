// package transfer provides the /transfers APIs
package transfer

import (
	"fmt"
	"net/url"
	"strconv"

	stripe "github.com/stripe/stripe-go"
)

// Client is used to invoke /transfers APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New POSTs a new transfer.
// For more details see https://stripe.com/docs/api#create_transfer.
func New(params *stripe.TransferParams) (*stripe.Transfer, error) {
	return getC().New(params)
}

func (c Client) New(params *stripe.TransferParams) (*stripe.Transfer, error) {
	body := &url.Values{
		"amount":    {strconv.FormatInt(params.Amount, 10)},
		"currency":  {string(params.Currency)},
		"recipient": {params.Recipient},
	}

	if len(params.Bank) > 0 {
		body.Add("bank_account", params.Bank)
	} else if len(params.Card) > 0 {
		body.Add("card", params.Card)
	}

	if len(params.Desc) > 0 {
		body.Add("description", params.Desc)
	}

	if len(params.Statement) > 0 {
		body.Add("statement_description", params.Statement)
	}

	params.AppendTo(body)

	transfer := &stripe.Transfer{}
	err := c.B.Call("POST", "/transfers", c.Key, body, transfer)

	return transfer, err
}

// Get returns the details of a transfer.
// For more details see https://stripe.com/docs/api#retrieve_transfer.
func Get(id string, params *stripe.TransferParams) (*stripe.Transfer, error) {
	return getC().Get(id, params)
}

func (c Client) Get(id string, params *stripe.TransferParams) (*stripe.Transfer, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}
		params.AppendTo(body)
	}

	transfer := &stripe.Transfer{}
	err := c.B.Call("GET", "/transfers/"+id, c.Key, body, transfer)

	return transfer, err
}

// Update updates a transfer's properties.
// For more details see https://stripe.com/docs/api#update_transfer.
func Update(id string, params *stripe.TransferParams) (*stripe.Transfer, error) {
	return getC().Update(id, params)
}

func (c Client) Update(id string, params *stripe.TransferParams) (*stripe.Transfer, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}

		if len(params.Desc) > 0 {
			body.Add("description", params.Desc)
		}

		params.AppendTo(body)
	}

	transfer := &stripe.Transfer{}
	err := c.B.Call("POST", "/transfers/"+id, c.Key, body, transfer)

	return transfer, err
}

// Cancel cancels a pending transfer.
// For more details see https://stripe.com/docs/api#cancel_transfer.
func Cancel(id string, params *stripe.TransferParams) (*stripe.Transfer, error) {
	return getC().Cancel(id, params)
}

func (c Client) Cancel(id string, params *stripe.TransferParams) (*stripe.Transfer, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}
		params.AppendTo(body)
	}

	transfer := &stripe.Transfer{}
	err := c.B.Call("POST", fmt.Sprintf("/transfers/%v/cancel", id), c.Key, body, transfer)

	return transfer, err
}

// List returns a list of transfers.
// For more details see https://stripe.com/docs/api#list_transfers.
func List(params *stripe.TransferListParams) *stripe.TransferIter {
	return getC().List(params)
}

func (c Client) List(params *stripe.TransferListParams) *stripe.TransferIter {
	type transferList struct {
		stripe.ListMeta
		Values []*stripe.Transfer `json:"data"`
	}

	var body *url.Values
	var lp *stripe.ListParams

	if params != nil {
		body = &url.Values{}

		if params.Created > 0 {
			body.Add("created", strconv.FormatInt(params.Created, 10))
		}

		if params.Date > 0 {
			body.Add("date", strconv.FormatInt(params.Date, 10))
		}

		if len(params.Recipient) > 0 {
			body.Add("recipient", params.Recipient)
		}

		if len(params.Status) > 0 {
			body.Add("status", string(params.Status))
		}

		params.AppendTo(body)
		lp = &params.ListParams
	}

	return &stripe.TransferIter{stripe.GetIter(lp, body, func(b url.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &transferList{}
		err := c.B.Call("GET", "/transfers", c.Key, &b, list)

		ret := make([]interface{}, len(list.Values))
		for i, v := range list.Values {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

func getC() Client {
	return Client{stripe.GetBackend(), stripe.Key}
}
