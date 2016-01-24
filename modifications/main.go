package modifications

import (
	"errors"
	"io"
	"strings"
)

// Definition is used to identify the modification and store appropriate options
type Definition struct {
	Name    string                 `json:"name"`
	Options map[string]interface{} `json:"options"`
}

// Modification is the interface any file modifier needs to implement
type Modification interface {
	Encode(io.Reader, chan<- error) io.Reader
	Decode(io.Reader, chan<- error) io.Reader
	Name() string
	Options() map[string]interface{}
	Configure(map[string]interface{}) error
}

// Pipeline is used to chain individual modifications together
type Pipeline struct {
	Head  io.Reader
	Tail  io.Reader
	Erroc chan<- error
}

// NewPipeline create a pipline connecting the provided modifications
func NewPipeline(head io.Reader, errc chan<- error, encode bool, mods ...Modification) *Pipeline {
	if mods == nil {
		panic("mods nil")
	}
	var next io.Reader
	for _, mod := range mods {
		if next == nil {
			if encode {
				next = mod.Encode(head, errc)
			} else {
				next = mod.Decode(head, errc)
			}
		} else {
			if encode {
				next = mod.Encode(next, errc)
			} else {
				next = mod.Decode(next, errc)
			}
		}
	}

	return &Pipeline{head, next, errc}
}

// GetModifications makes sure the agent can handle all the modifications
func GetModifications(m []Definition, reverse bool) ([]Modification, error) {
	var mods []Modification
	for _, modType := range m {
		switch strings.ToLower(modType.Name) {
		case "compress":
			compress := new(Compress)

			if err := compress.Configure(modType.Options); err != nil {
				return nil, err
			}

			mods = append(mods, compress)

		default:
			return nil, errors.New("I don't understand modification type: " + modType.Name)
		}

	}
	if reverse {
		for i, j := 0, len(mods)-1; i < j; i, j = i+1, j-1 {
			mods[i], mods[j] = mods[j], mods[i]
		}
	}
	return mods, nil
}
