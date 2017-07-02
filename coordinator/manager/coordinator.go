package manager

import (
	"github.com/sethjback/gobl/certificates"
	"github.com/sethjback/gobl/model"
)

func NewCAKey(key *model.Key) (err error) {
	if key == nil {
		key, err = certificates.NewCACertificate()
		if err != nil {
			return
		}
	}

	err = gDb.SaveKey("CA", *key)
	if err != nil {
		return err
	}

	grpcHup <- struct{}{}

	return nil
}

func GetCAKey() (*model.Key, error) {
	return gDb.GetKey("CA")
}
