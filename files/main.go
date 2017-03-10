package files

import "github.com/sethjback/gobl/modification"

// Signature contains all the data used to make an unique signature for a file
// this signature is used to determine wether the actual file needs to be sent
// to the engine again or not
type Signature struct {
	Path          string
	Hash          string
	Modifications []string
}

// Meta is information about the file as it was store on the drive
type Meta struct {
	Mode uint32
	UID  int
	GID  int
}

type File struct {
	Signature
	Meta
}

func NewSignature(path string, mods []modification.Definition) Signature {
	s := Signature{Path: path}
	for i := 0; i < len(mods); i++ {
		s.Modifications = append(s.Modifications, mods[i].Name)
	}

	return s
}
