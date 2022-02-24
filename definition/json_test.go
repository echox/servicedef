package definition

import (
	"testing"
)

func TestParseJSONServiceDef(t *testing.T) {

	SERVICEDEFS_PATH := "../services.json.example"

	var services ServiceDefs

	if json_error := services.Init(SERVICEDEFS_PATH); json_error != nil {
		t.Errorf("%v", json_error)
	}

	if len(services) != 2 {
		t.Errorf("wrong count of defined services in example file")
	}
}

func TestParseJSONHostDef(t *testing.T) {

	HOSTDEFS_PATH := "../hosts.json.example"

	var hosts HostDefs

	if json_error := hosts.Init(HOSTDEFS_PATH); json_error != nil {
		t.Errorf("%v", json_error)
	}

	if len(hosts) != 1 {
		t.Errorf("wrong count of defined hosts in example file")
	}
}

func TestParseJSONRulesDef(t *testing.T) {

	RULES_PATH := "../rules.json.example"

	var rules Rules

	if json_error := rules.Init(RULES_PATH); json_error != nil {
		t.Errorf("%v", json_error)
	}

	if len(rules) != 1 {
		t.Errorf("wrong count of defined rules in example file")
	}
}

func TestPathDoesntExist(t *testing.T) {

	var dummy Rules
	if err := parseJSONFile("doesntExist", dummy); err == nil {
		t.Errorf("path does not exist, expected error")
	}
}

func TestNoValidJSON(t *testing.T) {

	var dummy Rules
	if err := parseJSONFile("../README.md", dummy); err == nil {
		t.Errorf("no valid json, expected error")
	}
}
