package errorhandling

import (
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

// Note: `json` declared in errorhandling.go

func TestErrorEncoderFuncOmit(t *testing.T) {
	data, err := json.Marshal(struct {
		Err  error   `json:"err,omitempty"`
		Errs []error `json:"errs,omitempty"`
	}{})
	if err != nil {
		t.Fatal(err)
	}

	dataAsMap := make(map[string]interface{})
	err = json.Unmarshal(data, &dataAsMap)
	if err != nil {
		t.Fatal(err)
	}

	_, ok := dataAsMap["err"]
	if ok {
		t.Errorf("the `err` field should have been omitted")
	}
	_, ok = dataAsMap["errs"]
	if ok {
		t.Errorf("the `errs` field should have been omitted")
	}

	dataAsMap = make(map[string]interface{})
	data, err = json.Marshal(struct {
		Err  error   `json:"err"`
		Errs []error `json:"errs"`
	}{})
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(data, &dataAsMap)
	if err != nil {
		t.Fatal(err)
	}

	_, ok = dataAsMap["err"]
	if !ok {
		t.Errorf("the `err` field shouldn't have been omitted")
	}
	_, ok = dataAsMap["errs"]
	if !ok {
		t.Errorf("the `errs` field shouldn't have been omitted")
	}
}

func TestErrorEncoder(t *testing.T) {
	type payload struct {
		Err error
	}

	msg := "this is an error"
	body := payload{Err: fmt.Errorf(msg)}
	value, err := jsoniter.MarshalToString(body)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(`{"Err":%q}`, msg), value)
}

func TestErrorDecoder(t *testing.T) {
	type payload struct {
		Err error
	}

	var value payload
	msg := "this is an error"
	err := jsoniter.UnmarshalFromString(fmt.Sprintf(`{"Err":%q}`, msg), &value)
	assert.NoError(t, err)
	assert.Equal(t, msg, value.Err.Error())

	// This test is for corner case of Err changing value!
	err = jsoniter.UnmarshalFromString(`{"Err":null}`, &value)
	assert.NoError(t, err)
	assert.Nil(t, value.Err)
}

func TestErrorEncoderDecoder(t *testing.T) {
	type payload struct {
		Err error
	}

	msg := "this is an error"
	body := payload{Err: fmt.Errorf(msg)}
	value, err := jsoniter.Marshal(body)
	assert.NoError(t, err)

	var clientBody payload
	err = jsoniter.Unmarshal(value, &clientBody)
	assert.NoError(t, err)
	assert.Equal(t, body.Err.Error(), clientBody.Err.Error())
}

func TestErrorSliceEncoderDecoder(t *testing.T) {
	type payload struct {
		Err []error
	}

	tests := []struct {
		Err      []error
		expected []error
		len      int
	}{
		{[]error{}, []error{}, 0},
		{[]error{nil}, []error{nil}, 1},
		{[]error{fmt.Errorf("error")}, []error{fmt.Errorf("error")}, 1},
	}

	for _, testCase := range tests {
		body := payload{Err: testCase.Err}
		bytes, err := jsoniter.Marshal(body)
		assert.NoError(t, err)
		assert.Greater(t, len(bytes), 2)

		var clientBody payload
		err = jsoniter.Unmarshal(bytes, &clientBody)
		assert.NoError(t, err)
		assert.Len(t, clientBody.Err, testCase.len)
		assert.ElementsMatch(t, testCase.Err, clientBody.Err)
	}
}
