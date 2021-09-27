package definition

import (
	"encoding/json"
	"github.com/echox/servicedef/config"
	"io/ioutil"
	"log"
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

type Rules []RulesDef

func (rules *Rules) Init(cfg config.Config) error {

	log.Printf("loading rule file %v", cfg.Rules)

	byteValue, err := ioutil.ReadAll(cfg.Rules_File)
	if err != nil {
		return err
	}

	if json_error := json.Unmarshal(byteValue, rules); json_error != nil {
		return json_error
	}
	log.Printf("Rules #: %v", len(*rules))

	return nil
}
