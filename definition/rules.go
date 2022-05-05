package definition

import (
	"fmt"
	"net/http"
	"strings"

	. "github.com/echox/servicedef/result"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

// RulesDef can be referenced as one of the desired matching rules within a PortDef of a service
type RulesDef struct {
	Name            string    `json:"name"`             // Identifier of the rule
	Type_           string    `json:"type"`             // Type, currently 'http' is the only supported type
	FollowRedirects bool      `json:"follow-redirects"` // Type, currently 'http' is the only supported type
	Status          int       `json:"status"`           // Optional, expected http status
	Rules           []RuleDef `json:"rules"`            // Optional, Rules which must match
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
	return parseJSONFile(rulesFile, rules)
}


// Check all rules against a specified host and port
func (rules *Rules) Check(port PortDef, portResult *Port, service ServiceDef, ip string) {

	for _, pRule := range port.Rules {
		evaluated := false
		for _, r := range *(rules) {
			if r.Name == pRule {
				match, err := r.eval(port)
				portResult.RuleResults[r.Name] = match
				if err != nil {
					color.Set(color.FgYellow)
					log.Printf("! [%s] rule %s couldn't evaluated: %v", ip, r.Name, err)
					color.Unset()
				} else if match == false {
					color.Set(color.FgRed)
					log.Printf("! [%s] rule %s doesn't match %s", ip, r.Name, port.Uri)
					color.Unset()
				}
				evaluated = true
				continue
			}
		}
		if !evaluated {
			log.Printf("! [%v] rule %v not found in rule definitions", ip, pRule)
		}
	}
}

func (rule RulesDef) eval(port PortDef) (bool, error) {

	switch rule.Type_ {

	case "http":
		return evalHTTP(rule, port.Uri)
	default:
		return false, fmt.Errorf("Couldn't find implementation for rule type %s (%s)", rule.Type_, rule.Name)
	}

}

func evalHTTP(rules RulesDef, uri string) (bool, error) {

	if uri == "" {
		return false, fmt.Errorf("URI needed for checking %s", rules.Name)
	}

	hc := &http.Client{}

	if !rules.FollowRedirects {
		hc.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	r, err := hc.Get(uri)
	if err != nil {
		return false, err
	}
	defer r.Body.Close()
	//body, err := io.ReadAll(r.Body)

	if rules.Status != 0 {

		if r.StatusCode != rules.Status {
			log.Printf("! [%v] Status doesn't match. Expected %v got %v", uri, rules.Status, r.Status)
			return false, nil
		}
	}

	for _, rule := range rules.Rules {
		v := r.Header.Get(rule.Name)
		if v == "" {
			return false, nil
		} else {
			if rule.Contains == "" {
				log.Printf("rule %v on %v matches header-rule '%v'", rules.Name, uri, rule.Name)
				continue
			} else {
				if strings.Contains(v, rule.Contains) {
					log.Printf("[%v] matches '%v'", uri, rule.Contains)
					continue
				} else {
					log.Printf("! [%v] Header mismatch for %v. Expected %v got %v", uri, rules.Name, rule.Contains, v)
					return false, nil
				}

			}
		}

	}

	return true, nil
}
