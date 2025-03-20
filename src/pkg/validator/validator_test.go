package validator

import "testing"

type testStruct struct {
	Name              string `json:"name" validate:"required,len=11"`
	Description       string `json:"description" validate:"required,min=3"`
	Address           string `json:"address" validate:"required,min=3,max=5"`
	Age               int    `json:"age" validate:"required"`
	NumberOfQuestions int    `json:"numberOfQuestions" validate:"required,gt=10"`
}

func Test01(t *testing.T) {
	ts := testStruct{
		Name:              "1234567",
		Description:       "12",
		Address:           "123456",
		NumberOfQuestions: 5,
	}
	err := ValidateStruct(ts)
	if err == nil {
		t.Error("Expected error")
	}
}
