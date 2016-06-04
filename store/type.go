package store

import (
	"time"
)

const (
	StringType = iota
	IntType    = iota
	SetType    = iota
)

const (
	MINSCORE = "-inf"
	MAXSCORE = "+inf"
)

var (
	NilTime = time.Time{}
)

type KeyValue struct {
	Key        string       `json:"key"`
	Value      []byte       `json:"value"`
	Expiration time.Time    `json:"expiration"`
	ValueType  int          `json:"valueType"`
	Scores     SortedScores `json:"scores"`
}

type DB struct {
	Files map[string]string `json:"files"`
}
