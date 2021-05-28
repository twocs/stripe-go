package customer

import (
	"testing"

	assert "github.com/stretchr/testify/require"
	stripe "github.com/stripe/stripe-go/v72"
	_ "github.com/stripe/stripe-go/v72/testing"
)

func TestCustomerDel(t *testing.T) {
	customer, err := Del("cus_123", nil)
	assert.Nil(t, err)
	assert.NotNil(t, customer)
}

func TestCustomerGet(t *testing.T) {
	customer, err := Get("cus_123", nil)
	assert.Nil(t, err)
	assert.NotNil(t, customer)
}

func TestCustomerList(t *testing.T) {
	i := List(&stripe.CustomerListParams{})

	// Verify that we can get at least one customer
	assert.True(t, i.Next())
	assert.Nil(t, i.Err())
	assert.NotNil(t, i.Customer())
	assert.NotNil(t, i.CustomerList())
}

func TestCustomerListLimit(t *testing.T) {
	params := &stripe.CustomerListParams{}
	params.Filters.AddFilter("limit", "", "3")
	
	// Verify the Limit example shown in https://stripe.com/docs/api/customers/list?lang=go

	i := List(params)

	customers := make([]stripe.Customer)
	for i.Next() {
		c := i.Customer()
		customers = append(customers, *c)
	}
	
	assert.Nil(t, i.Err())
	assert.Len(t, customers, 3, "The limit should return 3")

}

func TestCustomerNew(t *testing.T) {
	customer, err := New(&stripe.CustomerParams{
		Email: stripe.String("foo@example.com"),
		Shipping: &stripe.CustomerShippingDetailsParams{
			Address: &stripe.AddressParams{
				Line1: stripe.String("line1"),
				City:  stripe.String("city"),
			},
			Name: stripe.String("name"),
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, customer)
}

func TestCustomerNew_NilParams(t *testing.T) {
	customer, err := New(nil)
	assert.Nil(t, err)
	assert.NotNil(t, customer)
}

func TestCustomerUpdate(t *testing.T) {
	customer, err := Update("cus_123", &stripe.CustomerParams{
		Email: stripe.String("foo@example.com"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, customer)
}
