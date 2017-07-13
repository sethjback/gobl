package modification

import (
	"io"
	"strings"

	"github.com/sethjback/gobl/goblerr"
)

const (
	ErrorUnrecognizedModifyer = "ModifyerNotRecognized"
	ErrorInvalidOptionValue   = "InvalidOptionValue"
	Forward                   = 1
	Backward                  = 2
)

// Definition is used to identify the modification and store appropriate options
type Definition struct {
	// Name of the modifyer to use
	Name string `json:"name"`
	// Options to configure the modifyer
	Options map[string]string `json:"options"`
}

// Option represents a value that can be used to configure a modifyer
type Option struct {
	// Name of the option
	Name string `json:"name"`
	// Description of what the option does
	Description string `json:"description"`
	// Type of vale, e.g string, bool, int
	Type string `json:"type"`
	// Default value for the option, if any
	Default interface{} `json:"default"`
}

// Modifyer is the required interface for modifying files while being backed up
type Modifyer interface {
	// Process takes an input and returns an output reader to pass on
	// if errors are encountered processing shoud stop and a message sent over the error chan
	// action specifies what should be done to the stream, essentially encoding or decoding
	Process(input io.Reader, errc chan<- error) io.Reader
	// Name of the modifyer
	Name() string
	// Options are the avaialble configuration option definitions
	Options() []Option
	// Configure sets the options for the modifyer
	Configure(options map[string]string) error
	// Direction sets the direction either forward or back for the modification
	Direction(direction int)
}

// Pipeline is used to chain individual modifications together
type Pipe struct {
	// Head of the pipe, e.g. the file stream of the file to be modified
	Head io.Reader
	// Tail output after all modifyers have been applied
	Tail io.Reader
	// Errorc is a channel over which modifyers will send any errors encountered
	Erroc chan error
}

// NewPipeline create a pipline connecting the provided modifications
func Pipeline(head io.Reader, mods ...Modifyer) *Pipe {
	errc := make(chan error)
	next := head
	for _, mod := range mods {
		next = mod.Process(next, errc)
	}

	return &Pipe{head, next, errc}
}

// Build takes defitions and configures the modifyers
// Build always reqiures the definitions to be passed in the same order:
// if the direction is forward, the modifications will be returned in this order
// if direction is backward, the modifications will be reversed
// If a modifyer is not recognized, an error will be returned
func Build(m []Definition, direction int) ([]Modifyer, error) {
	var mods []Modifyer
	for _, modType := range m {
		switch strings.ToLower(modType.Name) {
		case NameCompress:
			compress := &Compress{}

			if err := compress.Configure(modType.Options); err != nil {
				return nil, err
			}

			compress.Direction(direction)

			mods = append(mods, compress)

		default:
			return nil, goblerr.New("Invalid modifyer", ErrorUnrecognizedModifyer, "I don't understand modifyer type: "+modType.Name)
		}

	}
	if direction == Backward {
		for i, j := 0, len(mods)-1; i < j; i, j = i+1, j-1 {
			mods[i], mods[j] = mods[j], mods[i]
		}
	}
	return mods, nil
}

// Available retuns the list of available modifiers for an agent
func Available() []Modifyer {
	return []Modifyer{&Compress{}}
}
