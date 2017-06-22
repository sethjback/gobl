package config

type Store interface {
	Add(interface{}, interface{})
	Get(interface{}) (interface{}, bool)
	GetEnv(string) string
}

type baseCfg struct {
	options map[interface{}]interface{}
}

func New() Store {
	return &baseCfg{options: map[interface{}]interface{}{}}
}

func (b *baseCfg) Add(key interface{}, value interface{}) {
	b.options[key] = value
}

func (b *baseCfg) Get(key interface{}) (interface{}, bool) {
	i, ok := b.options[key]
	return i, ok
}

func (b *baseCfg) GetEnv(key string) string {
	if v, ok := b.options[key]; ok {
		vs, ok := v.(string)
		if ok {
			return vs
		}
	}

	return ""
}
