package definition

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func parseJSONFile(path string, target interface{}) error {

	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	if json_error := json.Unmarshal(byteValue, target); json_error != nil {
		return json_error
	}

	return nil
}
