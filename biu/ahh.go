package biu

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ahh struct {
	statusCode int
	body       []byte
	error      error
}

func newAhh(response *http.Response) *ahh {
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return newErrAhh(err)
	}

	return &ahh{
		statusCode: response.StatusCode,
		body:       b,
	}
}

func newErrAhh(err error) *ahh {
	return &ahh{
		error: err,
	}
}

func (a *ahh) Status() (statusCode int, err error) {
	return a.statusCode, a.error
}

func (a *ahh) String(strPtr *string) (statusCode int, err error) {
	if _, _err := a.Status(); _err != nil {
		return a.statusCode, _err
	}

	*strPtr = string(a.body)

	return a.statusCode, nil
}
func (a *ahh) Map(mapPtr *map[string]interface{}) (statusCode int, err error) {
	if _, _err := a.Status(); _err != nil {
		return a.statusCode, _err
	}

	decoder := json.NewDecoder(bytes.NewReader(a.body))
	decoder.UseNumber()

	err = decoder.Decode(mapPtr)
	return a.statusCode, err
}

func (a *ahh) Object(objectPtr interface{}) (statusCode int, err error) {
	if _, _err := a.Status(); _err != nil {
		return a.statusCode, _err
	}

	decoder := json.NewDecoder(bytes.NewReader(a.body))
	decoder.UseNumber()
	err = decoder.Decode(objectPtr)

	//err = json.Unmarshal(a.body, objectPtr)
	return a.statusCode, err
}
