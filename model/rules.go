package model

import (
	"encoding/json"
	"github.com/echox/servicedef/config"
	"io/ioutil"
	"log"
)

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
