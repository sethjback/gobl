package modification

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	assert := assert.New(t)

	defs := []Definition{Definition{Name: NameCompress}}

	mods, err := Build(defs, Forward)
	assert.Nil(err)
	assert.Len(mods, 1)

	mods, err = Build(defs, Backward)
	assert.Nil(err)
	assert.Len(mods, 1)

	defs = append(defs, Definition{Name: "asdf"})

	mods, err = Build(defs, Forward)
	assert.NotNil(err)
}

func TestPipe(t *testing.T) {
	assert := assert.New(t)

	defs := []Definition{Definition{Name: NameCompress}}
	mods, err := Build(defs, Forward)
	assert.Nil(err)

	outC := make(chan []byte)

	input := []byte("This is the input string")
	r := bytes.NewReader(input)
	pipe := Pipeline(r, mods...)

	go func() {
		gOut, ioErr := ioutil.ReadAll(pipe.Tail)
		if ioErr != nil {
			pipe.Erroc <- ioErr
		}
		outC <- gOut
	}()

	var fOut []byte
	select {
	case err = <-pipe.Erroc:
		assert.Nil(err)
		return
	case fOut = <-outC:
		assert.NotEmpty(fOut)
	}

	mods, err = Build(defs, Backward)
	assert.Nil(err)

	pipe = Pipeline(bytes.NewReader(fOut), mods...)

	go func() {
		gOut, ioErr := ioutil.ReadAll(pipe.Tail)
		if ioErr != nil {
			pipe.Erroc <- ioErr
		}
		outC <- gOut
	}()

	var bOut []byte
	select {
	case err = <-pipe.Erroc:
		assert.Nil(err)
		return
	case bOut = <-outC:
		assert.Equal(input, bOut)
	}

	assert.NotEqual(fOut, bOut)
}
