package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Test",
		Price: 1.00,
		SKU:   "abc-abc-abc",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
