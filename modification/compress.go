package modification

import (
	"compress/gzip"
	"io"
	"strconv"

	"github.com/sethjback/gobl/goblerr"
)

const NameCompress = "compress"

// Compress modification takes a file and compresses it
type Compress struct {
	method    string
	level     int
	direction int
}

func (c *Compress) Process(input io.Reader, errc chan<- error) io.Reader {
	if c.direction == Backward {
		return decode(c.method, input, errc)
	}
	return encode(c.method, c.level, input, errc)
}

// Encode compresses the file using the defined method and options
func encode(method string, level int, input io.Reader, errc chan<- error) io.Reader {
	r, w := io.Pipe()

	go func() {
		defer w.Close()
		cw, err := getEncoder(method, level, w)
		if err != nil {
			errc <- err
			return
		}

		if _, err := io.Copy(cw, input); err != nil {
			errc <- err
			return
		}
		cw.Flush()
		cw.Close()
	}()

	return r
}

// Decode decodes an incoming stream
func decode(method string, reader io.Reader, errc chan<- error) io.Reader {
	r, err := getDecoder(method, reader)
	if err != nil {
		errc <- err
	}
	return r
}

func getDecoder(method string, stream io.Reader) (io.Reader, error) {
	switch method {
	case "gzip":
		return gzip.NewReader(stream)
	}
	return nil, goblerr.New("unrecognized decoder", ErrorInvalidOptionValue, nil)
}

func getEncoder(method string, level int, stream io.Writer) (*gzip.Writer, error) {
	switch method {
	case "gzip":
		return gzip.NewWriterLevel(stream, level)
	}
	return nil, goblerr.New("unrecognized encoder", ErrorInvalidOptionValue, nil)
}

// Name returns the modifications's name
func (c *Compress) Name() string {
	return NameCompress
}

// Options retun a list of possible options
func (c *Compress) Options() []Option {
	return []Option{
		Option{
			Name:        "method",
			Description: "compression method to use",
			Type:        "string",
			Default:     "gzip",
		},
		Option{
			Name:        "level",
			Description: "compression level to use",
			Type:        "int",
			Default:     5,
		}}
}

func (c *Compress) Direction(d int) {
	c.direction = d
}

// Configure configures the compressor
func (c *Compress) Configure(options map[string]string) error {
	//defaults

	c.level = 5
	c.method = "gzip"
	for k, v := range options {
		switch k {
		case "method":
			if v != "gzip" {
				return goblerr.New("method not supported", ErrorInvalidOptionValue, "acceptible options are: gzip")
			}

			c.method = v
		case "level":
			valI, err := strconv.Atoi(v)
			if err != nil {
				return goblerr.New("level must be int", ErrorInvalidOptionValue, "level must be between 1 and 9")
			}
			if valI < 1 || valI > 9 {
				return goblerr.New("level invalid", ErrorInvalidOptionValue, "level must be between 1 and 9")
			}
			c.level = valI
		}
	}

	return nil
}
