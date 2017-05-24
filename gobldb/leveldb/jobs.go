package leveldb

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/sethjback/gobl/gobldb/errors"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/model"
	"github.com/syndtr/goleveldb/leveldb"
)

func (l *Leveldb) SaveJob(j model.Job) error {
	lj := &job{
		ID:      j.ID,
		AgentId: j.Agent.ID,
		Def:     j.Definition,
		Meta:    j.Meta,
	}

	jbyte, err := json.Marshal(lj)
	if err != nil {
		return goblerr.New("Unable to save job", errors.ErrCodeMarshal, err)
	}
	err = l.Connection.Put([]byte(keyTypeJob+j.ID), jbyte, nil)
	if err != nil {
		return goblerr.New("Unable to save job", errors.ErrCodeSave, err)
	}

	l.NewIndex(index{itype: indexTypeJobDate, key: "start-" + strconv.Itoa(int(j.Meta.Start.UnixNano())) + j.ID, value: j.ID})
	l.NewIndex(index{itype: indexTypeJobDate, key: "end-" + strconv.Itoa(int(j.Meta.End.UnixNano())) + j.ID, value: j.ID})
	if i, _ := l.GetIndex(indexTypeJobState, j.ID); i != nil {
		if i.value != j.Meta.State {
			l.NewIndex(index{itype: indexTypeJobState, key: j.Meta.State + j.ID, value: j.ID})
			l.NewIndex(index{itype: indexTypeJobState, key: j.ID, value: j.Meta.State})
			l.DeleteIndex(index{itype: indexTypeJobState, key: i.value + j.ID, value: j.ID})
		}
	} else {
		l.NewIndex(index{itype: indexTypeJobState, key: j.Meta.State + j.ID, value: j.ID})
		l.NewIndex(index{itype: indexTypeJobState, key: j.ID, value: j.Meta.State})
	}

	l.NewIndex(index{itype: indexTypeJobAgent, key: j.Agent.ID + j.ID, value: j.ID})

	return nil
}

func (l *Leveldb) GetJob(id string) (*model.Job, error) {
	jbyte, err := l.Connection.Get([]byte(keyTypeJob+id), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, goblerr.New("No job with that ID", errors.ErrCodeNotFound, err)
		}
		return nil, goblerr.New("Unable to get job", errors.ErrCodeGet, err)
	}

	var lj job
	err = json.Unmarshal(jbyte, &lj)
	if err != nil {
		return nil, goblerr.New("Unable to get job", errors.ErrCodeUnMarshal, err)
	}

	j := &model.Job{
		ID:         lj.ID,
		Definition: lj.Def,
		Meta:       lj.Meta,
	}

	j.Agent, err = l.GetAgent(lj.AgentId)
	if err != nil {
		return nil, goblerr.New("Unable to get job", errors.ErrCodeUnMarshal, err)
	}

	j.Meta.Errors, _ = l.jobFileCount(j.ID, map[string]string{"state": model.StateFailed})
	j.Meta.Complete, _ = l.jobFileCount(j.ID, map[string]string{"state": model.StateFinished})

	return j, nil
}

func (l *Leveldb) JobList(filters map[string]string) ([]model.Job, error) {
	limit := 10
	offset := 0

	var ids []string
	var start, end, state, agent string

	for k, v := range filters {
		switch k {
		case "state":
			state = v
		case "start", "end":
			ts, err := parseDate(v)
			if err != nil {
				return nil, goblerr.New("date stamp invalid", errors.ErrFilterOptions, err)
			}
			if k == "start" {
				start = strconv.Itoa(ts)
			} else {
				end = strconv.Itoa(ts)
			}
		case "agent":
			agent = v
		case "limit", "offset":
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, goblerr.New(k+" invalid", errors.ErrFilterOptions, err)
			}
			if k == "limit" {
				limit = i
			} else {
				offset = i
			}
		default:
			return nil, goblerr.New("unrecognized filter: "+k+"", errors.ErrFilterOptions, "Valid options are: state, start, end, agent, limit, and offset")
		}
	}

	if start != "" || end != "" {
		startIds := make([]string, 0)
		if start != "" {
			is, err := l.indexRange(indexTypeJobDate, "start-"+start, "start-"+strconv.Itoa(int(time.Now().UTC().Unix())+500))
			if err != nil {
				return nil, goblerr.New("Unable to get job list", errors.ErrCodeGet, err)
			}
			for _, i := range is {
				startIds = append(startIds, i.value)
			}
		}

		endIds := make([]string, 0)
		if end != "" {
			is, err := l.indexRange(indexTypeJobDate, "end-", "end-"+end)
			if err != nil {
				return nil, goblerr.New("Unable to get job list", errors.ErrCodeGet, err)
			}

			// if we set a start date, we only want the ids intersection of the slices
			for _, i := range is {
				endIds = append(endIds, i.value)
			}
		}

		if start != "" && end == "" {
			ids = startIds
		}

		if end != "" && start == "" {
			ids = endIds
		}

		if end != "" && start != "" {
			ids = intersectSlice(startIds, endIds)
		}

	}

	if agent != "" {
		if ids == nil {
			is, err := l.indexQuery(indexTypeJobAgent, agent)
			if err != nil {
				return nil, goblerr.New("Unable to get job list", errors.ErrCodeGet, err)
			}
			for _, i := range is {
				if !stringInSlice(ids, i.value) {
					ids = append(ids, i.value)
				}
			}
		} else {
			for i := len(ids) - 1; i > 0; i-- {
				if ids[i] != agent {
					ids = append(ids[:i], ids[i+1:]...)
				}
			}
		}

		if ids == nil {
			ids = make([]string, 0)
		}
	}

	if state != "" {
		if ids == nil {
			is, err := l.indexQuery(indexTypeJobState, state)
			if err != nil {
				return nil, goblerr.New("Unable to get job list", errors.ErrCodeGet, err)
			}
			for _, i := range is {
				if !stringInSlice(ids, i.value) {
					ids = append(ids, i.value)
				}
			}
		} else {
			for i := len(ids) - 1; i > 0; i-- {
				if ids[i] != agent {
					ids = append(ids[:i], ids[i+1:]...)
				}
			}
		}
	}

	// there were not filters (apart from limit or count)
	if start == "" && end == "" && state == "" && agent == "" {
		is, err := l.indexQuery(indexTypeJobDate, "start-")
		if err != nil {
			return nil, goblerr.New("Unable to get job list", errors.ErrCodeGet, err)
		}
		for _, i := range is {
			ids = append(ids, i.value)
		}
	}

	if len(ids) < offset {
		return []model.Job{}, nil
	}

	if offset+limit > len(ids) {
		limit = len(ids) - offset
	}

	jobs := []model.Job{}

	for _, i := range ids[offset:limit] {
		j, err := l.GetJob(i)
		if err != nil {
			return nil, goblerr.New("Unable to get job list", errors.ErrCodeGet, err)
		}
		jobs = append(jobs, *j)
	}

	return jobs, nil
}
