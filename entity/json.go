package entity

import "time"

type EmbededJson map[string]interface{}

type ListJson []EmbededJson

type RequestObj struct {
	Time         time.Time `bson: time`
	Path         string    `bson: path`
	Method       string    `bson: method`
	Status       int       `bson: status`
	ResponseBody string    `bson: responsebody`
}

func (obj EmbededJson) AppendCustomKey() {
	obj["custom_key"] = true
}

func (s ListJson) AppendCustomKey() {
	for i := range s {
		s[i].AppendCustomKey()
	}
}
