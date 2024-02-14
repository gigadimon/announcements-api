package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func ReadAndUnmarshallInputBody(body io.Reader, destination interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("reading input data failed: %s", err.Error())

	}

	if err := json.Unmarshal(b, &destination); err != nil {
		return fmt.Errorf("unmarshalling input data failed: %s", err.Error())
	}

	return nil
}
