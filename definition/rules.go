package definition

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// RulesDef collection which can be referenced as Rules within a PortDef
type RulesDef struct {
	Name   string    `json:"name"`   // Identifier of the rule
	Type_  string    `json:"type"`   // Type, currently 'http' is the only supported type
	Status int       `json:"status"` // Optional, expected http status
	Rules  []RuleDef `json:"rules"`  // Optional, Rules which must match
}

// RuleDef which must be matched by a service
type RuleDef struct {
	Name     string `json:"name"`     // HTTP header with given name must exist
	Contains string `json:"contains"` // Optional, HTTP header must contain given value
}

// Collection of rules definitions from json file
type Rules []RulesDef

// Init loads the rule definitions form a json file
func (rules *Rules) Init(rulesFile string) error {

	jsonFile, err := os.Open(rulesFile)
	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	if json_error := json.Unmarshal(byteValue, rules); json_error != nil {
		return json_error
	}

	return nil
}
