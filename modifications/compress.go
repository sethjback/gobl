package modifications

import (
	"compress/gzip"
	"errors"
	"io"
)

// Compress modification takes a file and compresses it
type Compress struct {
	compressMethod string
	compressLevel  int
}

// Encode compresses the file using the defined method and options
func (c *Compress) Encode(reader io.Reader, errc chan<- error) io.Reader {
	r, w := io.Pipe()

	go func() {
		defer w.Close()

		gw, err := gzip.NewWriterLevel(w, c.compressLevel)
		if err != nil {
			errc <- err
			return
		}
		defer gw.Close()

		if _, err := io.Copy(gw, reader); err != nil {
			errc <- err
			return
		}
	}()

	return r
}

// Decode decodes an incoming stream
func (c *Compress) Decode(reader io.Reader, errc chan<- error) io.Reader {
	r, err := gzip.NewReader(reader)
	if err != nil {
		errc <- err
	}
	return r
}

// Name returns the modifications's name
func (c *Compress) Name() string {
	return "Compress"
}

// Options retun a list of possible options
func (c *Compress) Options() map[string]interface{} {
	return nil
}

// Configure configures the compressor
func (c *Compress) Configure(options map[string]interface{}) error {
	val, ok := options["method"]
	if !ok {
		return errors.New("Must supply compression method")
	}

	vals, ok := val.(string)
	if !ok {
		return errors.New("Compression method must be string")
	}

	if vals != "gzip" {
		return errors.New("Only gzip supported at this time")
	}

	c.compressMethod = vals

	val, ok = options["compressionLevel"]
	if !ok {
		//default to 5
		c.compressLevel = 5
	} else {
		vali, ok := val.(int)
		if !ok {
			return errors.New("compressoinLevel must be an integer between 1 and 9")
		}
		c.compressLevel = vali
	}

	return nil
}
