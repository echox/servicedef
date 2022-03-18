package definition

import (
	"net/http"
	"strings"

	. "github.com/echox/servicedef/result"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

// RulesDef collection which can be referenced as Rules within a PortDef
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
		eval := false
		for _, r := range *(rules) {
			if r.Name == pRule {
				if r.Type_ == "http" {
					s := evalHTTP(r, port.Uri)
					portResult.RuleResults[r.Name] = s
					if s == false {
						color.Set(color.FgRed)
						log.Printf("! [%v] rule %v doesn't match %v", ip, r.Name, port.Uri)
						color.Unset()
					}
					eval = true
					continue
				}
			}
		}
		if !eval {
			log.Printf("! [%v] rule %v not found in rule definitions", ip, pRule)
		}
	}
}

func evalHTTP(rules RulesDef, uri string) bool {

	if uri == "" {
		log.Warnf("URI needed for checking %v", rules.Name)
		return false
	}

	hc := &http.Client{}

	if !rules.FollowRedirects {
		hc.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	r, err := hc.Get(uri)
	if err != nil {
		log.Error(err)
		return false
	}
	defer r.Body.Close()
	//body, err := io.ReadAll(r.Body)

	if rules.Status != 0 {

		if r.StatusCode != rules.Status {
			log.Printf("! [%v] Status doesn't match. Expected %v got %v", uri, rules.Status, r.Status)
			return false
		}
	}

	for _, rule := range rules.Rules {
		v := r.Header.Get(rule.Name)
		if v == "" {
			return false
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
					return false
				}

			}
		}

	}

	return true
}
