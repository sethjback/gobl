package job

import (
	"encoding/json"

	"github.com/sethjback/gobl/model"
)

// Jobber
type Jobber interface {
	Run(done chan<- string)
	Cancel()
	Status() model.JobMeta
}

type JobNotification struct {
	JF   *model.JobFile
	host string
	path string
}

func (jn *JobNotification) Host() string {
	return jn.host
}

func (jn *JobNotification) Path() string {
	return jn.path
}

func (jn *JobNotification) Body() []byte {
	var b []byte
	if jn.JF != nil {
		b, _ = json.Marshal(jn.JF)
	}
	return b
}
