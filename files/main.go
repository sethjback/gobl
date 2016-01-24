package files

// Signature contains all the data used to make an unique signature for a file
// this signature is used to determine wether the actual file needs to be sent
// to the engine again or not
type Signature struct {
	Name          string
	Path          string
	Hash          string
	Modifications []string
}

// FileMeta is information about the file as it was store on the drive
type FileMeta struct {
	Mode uint32
	UID  int
	GID  int
}
