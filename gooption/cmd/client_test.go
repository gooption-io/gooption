package cmd

import (
	"testing"
)

func Test_insertImpliedVolRequest(t *testing.T) {
	insertImpliedVolRequest()
}

func Test_insertPriceRequest(t *testing.T) {
	insertPriceRequest()
}

func Test_insertGreekRequest(t *testing.T) {
	insertGreekRequest()
}

func Test_priceRequest(t *testing.T) {
	client{}.priceRequest()
}

func Test_ivRequest(t *testing.T) {
	client{}.ivRequest()
}

func Test_greekRequest(t *testing.T) {
	client{}.greekRequest()
}
