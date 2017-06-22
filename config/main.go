package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func Map(filePath string) (map[string]string, error) {
	amap := map[string]string{}
	for _, v := range os.Environ() {
		vs := strings.Split(v, "=")
		if len(vs) == 2 {
			amap[vs[0]] = vs[1]
		}
	}

	if filePath != "" {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error opening config file: %v", err))
		}
		s := bufio.NewScanner(io.LimitReader(f, 10000000))
		for s.Scan() {
			line := s.Text()
			if len(line) == 0 || len(line) < 2 || line[:2] == "//" {
				continue
			}
			split := strings.Split(line, "=")
			if len(split) != 2 {
				continue
			}
			amap[strings.TrimSpace(split[0])] = strings.TrimSpace(split[1])
		}
		err = f.Close()
		if err != nil {
			return nil, errors.New(fmt.Sprint("unable to close config file"))
		}
	}

	return amap, nil
}

func SetEnvMap(env map[string]string, store Store) Store {
	for k, v := range env {
		store.Add(k, v)
	}
	return store
}
