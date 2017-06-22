package coordinator

import (
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/goblerr"
)

type key int

const (
	coordConfigKey key = 0

	ErrorConfigMissing = "ConfigMissing"
)

type Coordinator struct {
	Address string
}

type coordConfig struct {
	address string
}

func FromConfig(cs config.Store) (c *Coordinator) {
	cc := configFromStore(cs)
	if cc != nil {
		c = &Coordinator{Address: cc.address}
	}

	return c
}

func SaveConfig(cs config.Store, env map[string]string) error {
	cc := &coordConfig{}
	for k, v := range env {
		if k == "COORDINATOR_ADDRESS" {
			cc.address = v
		}
	}

	if cc.address == "" {
		return goblerr.New("Coordinator configuation required", ErrorConfigMissing, nil)
	}

	cs.Add(coordConfigKey, cc)

	return nil
}

func configFromStore(cs config.Store) *coordConfig {
	if cc, ok := cs.Get(coordConfigKey); ok {
		return cc.(*coordConfig)
	}
	return nil
}
