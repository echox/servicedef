package export

import (
	log "github.com/sirupsen/logrus"

	"encoding/json"
	"io/ioutil"

	. "github.com/echox/servicedef/result"
)

func WriteJSON(results ResultHosts, path string) {

	log.Println("exporting results as json...")
	resultFile, _ := json.MarshalIndent(results, "", " ")
	if err := ioutil.WriteFile(path, resultFile, 0644); err != nil {
		log.Error(err)
	}
	log.Tracef("finished exporting results as json")
}
