package data

import "testing"

func TestProductValidation(t *testing.T) {
	p := Product{
		Name:  "as",
		Price: 3.0,
		SKU:   "wer-fsd-df33333",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
