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
	dest string
}

func (jn *JobNotification) Destination() string {
	return jn.dest
}

func (jn *JobNotification) Body() []byte {
	b, _ := json.Marshal(jn.JF)
	return b
}
