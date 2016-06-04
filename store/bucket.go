package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"sync"
	"time"
)

const (
	fileName = "%s/bucket_%d.json"
)

//Holds keyvalue with synchronized for concurrent access
type Bucket struct {
	sync.Mutex
	id   int
	item map[string]KeyValue
}

func NewBucket(id int) *Bucket {
	return &Bucket{
		id:   id,
		item: make(map[string]KeyValue),
	}
}

func (b *Bucket) Set(kv KeyValue) {
	b.Lock()
	defer b.Unlock()

	b.item[kv.Key] = kv
}

//Return false if key is expired
func (b *Bucket) Get(key string) (KeyValue, bool) {
	b.Lock()
	defer b.Unlock()

	kv, ok := b.item[key]
	if !ok {
		return KeyValue{}, false
	}

	if kv.Expiration.Equal(NilTime) {
		return kv, true
	}

	if kv.Expiration.After(time.Now()) {
		return kv, true
	}

	return KeyValue{}, false
}

//Delete expired keyvalues from bucket
func (b *Bucket) Cleanup() {
	b.Lock()
	defer b.Unlock()

	currentTime := time.Now()
	for key, kv := range b.item {
		if kv.Expiration == NilTime {
			continue
		}

		if kv.Expiration.After(currentTime) {
			continue
		}

		delete(b.item, key)
	}
}

//Save the bucket to json file
func (b *Bucket) Save(filePath string) (string, error) {
	b.Lock()
	defer b.Unlock()

	if filePath == "" {
		filePath = "."
	}

	jsonBytes, err := json.Marshal(b.item)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal bucket items. err: %v", err)
	}

	bucket_file := fmt.Sprintf(fileName, filePath, b.id)
	bucket_file = path.Clean(bucket_file)

	if err := ioutil.WriteFile(bucket_file, jsonBytes, 0644); err != nil {
		return "", fmt.Errorf("Failed to write bucket items to file: %s. err: %v", bucket_file, err)
	}

	return bucket_file, nil
}

//Load the bucket from json file
func (b *Bucket) Load(fileName string) (bool, error) {
	b.Lock()
	defer b.Unlock()

	mBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return false, fmt.Errorf("Failed to read bucket items from file. err: %v", err)
	}

	b.item = make(map[string]KeyValue)
	if err := json.Unmarshal(mBytes, &b.item); err != nil {
		return false, fmt.Errorf("Failed to unmarshal bucket items. err: %v", err)
	}

	return true, nil
}
