package util

import (
	"fmt"
	"strconv"
	"strings"
)

// Version of the program
type Version struct {
	Major int
	Minor int
	Patch int
}

// NewVersion Returns a new version struct
func NewVersion(version string) Version {
	split := strings.Split(version, ".")

	M, err := strconv.Atoi(split[0])
	if err != nil {
		panic("invalid version string")
	}

	m, err := strconv.Atoi(split[1])
	if err != nil {
		panic("invalid version string")
	}

	p, err := strconv.Atoi(split[2])
	if err != nil {
		panic("invalid version string")
	}

	return Version{M, m, p}
}

func (version Version) String() string {
	return fmt.Sprintf("%v.%v.%v", version.Major, version.Minor, version.Patch)
}
