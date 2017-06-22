package leveldb

import (
	"encoding/json"
	"strings"

	"github.com/sethjback/gobl/gobldb/errors"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/model"
	"github.com/syndtr/goleveldb/leveldb"
)

func (l *Leveldb) SaveJobFile(jobID string, f model.JobFile) error {
	fbyte, err := json.Marshal(f)
	if err != nil {
		return goblerr.New("Unable to save file", errors.ErrCodeMarshal, err)
	}

	split := strings.Split(f.File.Path, "/")

	batch := indexBatchPut(
		index{indexTypeFileState, jobID + f.State + f.File.Path, f.File.Path},
		index{indexTypeFileParent, jobID + strings.Join(split[:len(split)-1], "/") + f.File.Path, f.File.Path},
	)

	batch.Put([]byte(keyTypeFile+jobID+f.File.Path), fbyte)

	for i, v := range split[1 : len(split)-1] {
		var parent string
		if i == 0 {
			parent = "/"
		} else {
			parent = strings.Join(split[:i+1], "/")
		}
		iv, err := l.Connection.Get([]byte(keyTypeFileDir+jobID+parent+"/"), nil)
		if err != nil {
			if err != leveldb.ErrNotFound {
				return goblerr.New("Unable to save file", errors.ErrCodeSave, err)
			}

			batch.Put([]byte(keyTypeFileDir+jobID+parent+"/"), []byte(v))
		} else {
			ivs := strings.Split(string(iv), ",")
			if !stringInSlice(ivs, v) {
				ivs = append(ivs, v)
				batch.Put([]byte(keyTypeFileDir+jobID+parent+"/"), []byte(strings.Join(ivs, ",")))
			}
		}

	}

	err = l.Connection.Write(batch, nil)
	if err != nil {
		return goblerr.New("Unable to save file", errors.ErrCodeSave, err)
	}

	return nil
}

func (l *Leveldb) getFile(jobID, filePath string) (*model.JobFile, error) {
	fbyte, err := l.Connection.Get([]byte(keyTypeFile+jobID+filePath), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, goblerr.New("Unable to find that file", errors.ErrCodeNotFound, err)
		}
		return nil, goblerr.New("Unable to get that file", errors.ErrCodeGet, err)
	}

	var f model.JobFile
	err = json.Unmarshal(fbyte, &f)
	if err != nil {
		return nil, goblerr.New("Unable to get that file", errors.ErrCodeUnMarshal, err)
	}

	return &f, nil
}

func (l *Leveldb) jobFileCount(jobID string, filters map[string]string) (int, error) {
	count := 0
	dir := ""
	state := ""
	for k, v := range filters {
		k = strings.ToLower(k)
		switch k {
		case "dir", "parent":
			if v == "*" {
				dir = jobID
			} else {
				dir = jobID + v + "/"
			}
		case "state":
			state = v
		default:
			return -1, goblerr.New("Invalid filter. Use state or parent", errors.ErrFilterOptions, nil)
		}
	}

	if dir != "" {
		vals, err := l.indexQuery(indexTypeFileParent, dir)
		if err != nil {
			return -1, goblerr.New("Unable to get file list", errors.ErrCodeGet, err)
		}

		if state != "" {
			for _, iv := range vals {
				f, err := l.getFile(jobID, iv.value)
				if err == nil {
					if state != "" && f.State != state {
						continue
					}
					count++
				} else {

				}
			}
		} else {
			count = len(vals)
		}

	} else if state != "" {
		vals, err := l.indexQuery(indexTypeFileState, jobID+state)
		if err != nil {
			return -1, goblerr.New("Unable to get file list", errors.ErrCodeGet, err)
		}

		count = len(vals)
	}

	return count, nil
}

func (l *Leveldb) JobFileList(jobID string, filters map[string]string) ([]model.JobFile, error) {
	dir := ""
	state := ""
	for k, v := range filters {
		k = strings.ToLower(k)
		switch k {
		case "dir", "parent":
			if v == "*" {
				dir = jobID
			} else {
				dir = jobID + v + "/"
			}
		case "state":
			state = v
		default:
			return nil, goblerr.New("Invalid filter. Use state or parent", errors.ErrFilterOptions, nil)
		}
	}

	jf := make([]model.JobFile, 0)

	if dir != "" {
		vals, err := l.indexQuery(indexTypeFileParent, dir)
		if err != nil {
			return nil, goblerr.New("Unable to get file list", errors.ErrCodeGet, err)
		}

		for _, iv := range vals {
			f, err := l.getFile(jobID, iv.value)
			if err == nil {
				if state != "" && f.State != state {
					continue
				}
				jf = append(jf, *f)
			} else {

			}
		}

	} else if state != "" {
		vals, err := l.indexQuery(indexTypeFileState, jobID+state)
		if err != nil {
			return nil, goblerr.New("Unable to get file list", errors.ErrCodeGet, err)
		}

		for _, iv := range vals {
			f, err := l.getFile(jobID, iv.value)
			if err == nil {
				jf = append(jf, *f)
			} else {

			}
		}
	}

	return jf, nil
}

func (l *Leveldb) JobDirectories(jobID, parent string) ([]string, error) {
	val, err := l.Connection.Get([]byte(keyTypeFileDir+jobID+parent+"/"), nil)
	if err != nil {
		return nil, goblerr.New("Unable to get directory list", errors.ErrCodeGet, err)
	}

	return strings.Split(string(val), ","), nil
}
