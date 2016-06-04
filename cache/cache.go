package cache

import (
	"encoding/json"
	"fmt"
	"github.com/madhusudhancs/redis/store"
	"github.com/madhusudhancs/redis/utils"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// Cache holds buckets
type Cache struct {
	noOfBuckets uint64
	buckets     []Bucket
}

// Input holds all the args passed to Set API
type SetInput struct {
	Key    string
	Value  string
	Expiry time.Duration
	NX     bool
	XX     bool
}

// Input holds all the args passed to Set API
type ZSetInput struct {
	Key    string
	Scores []store.ScoreMember
	NX     bool
	XX     bool
	CH     bool
	INCR   bool
}

// Bucket interface provides API to acess bucket
type Bucket interface {
	//Store key value in bucket
	Set(kv store.KeyValue)

	//Returns key value from bucket. Return nil if key does not exist
	Get(key string) (store.KeyValue, bool)

	//Cleanup deletes all the expired key values
	Cleanup()

	//Save the entiries to file
	Save(path string) (string, error)

	//Load the entries from file
	Load(fileName string) (bool, error)
}

func NewCache(noOfBuckets int) *Cache {
	c := Cache{
		noOfBuckets: uint64(noOfBuckets),
		buckets:     make([]Bucket, noOfBuckets),
	}

	for i := 0; i < noOfBuckets; i++ {
		c.buckets[i] = store.NewBucket(i)
	}

	return &c
}

func requiredFieldsPresentSet(input SetInput) bool {
	if input.Key == "" {
		return false
	}

	if input.Value == "" {
		return false
	}

	return true
}

func (c *Cache) Set(input SetInput) error {
	if !requiredFieldsPresentSet(input) {
		return fmt.Errorf("mandatory key and value not present in the input")
	}

	bucketID := c.bucketID(input.Key)

	_, ok := c.buckets[bucketID].Get(input.Key)

	if input.NX && ok {
		return fmt.Errorf("Key: %s exists", input.Key)
	}

	if input.XX && !ok {
		return fmt.Errorf("Key: %s does not exist", input.Key)
	}

	kv := store.KeyValue{
		Key:       input.Key,
		Value:     []byte(input.Value),
		ValueType: store.StringType,
	}

	if input.Expiry != 0 {
		kv.Expiration = time.Now().Add(input.Expiry)
	}

	c.buckets[bucketID].Set(kv)

	return nil
}

func (c *Cache) bucketID(key string) int {
	return int(Hash(key) % c.noOfBuckets)
}

func (c *Cache) Get(key string) (string, error) {
	bucketID := c.bucketID(key)

	kv, ok := c.buckets[bucketID].Get(key)
	if !ok {
		return "", nil
	}

	if kv.ValueType != store.StringType {
		return "", fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return string(kv.Value), nil
}

func (c *Cache) SetBit(key string, offset, value int) (uint8, error) {
	if value != 0 && value != 1 {
		return 0, fmt.Errorf("ERR bit is not an integer or out of range")
	}

	bucketID := c.bucketID(key)

	kv, ok := c.buckets[bucketID].Get(key)
	if !ok {
		kv = store.KeyValue{
			Key:       key,
			Value:     []byte{},
			ValueType: store.StringType,
		}
	}

	byteValue, oldBitValue := utils.PutBit(kv.Value, offset, value)

	kv.Value = byteValue
	c.buckets[bucketID].Set(kv)

	return oldBitValue, nil
}

func (c *Cache) GetBit(key string, offset int) uint8 {
	bucketID := c.bucketID(key)

	kv, ok := c.buckets[bucketID].Get(key)
	if !ok {
		return 0
	}

	return utils.GetBit(kv.Value, offset)
}

func (c *Cache) ZAdd(input ZSetInput) (int, error) {
	if len(input.Scores) > 1 && input.INCR {
		return 0, fmt.Errorf("ERR INCR option supports a single increment-element pair")
	}

	bucketID := c.bucketID(input.Key)

	kv, ok := c.buckets[bucketID].Get(input.Key)
	if !ok {
		kv = store.KeyValue{
			Key:       input.Key,
			Scores:    store.NewSortedScores(),
			ValueType: store.SetType,
		}
	}

	if kv.ValueType != store.SetType {
		return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	add := 0
	update := 0
	newScore := 0.0

	for _, scoreMember := range input.Scores {
		sMember, ok := kv.Scores.Get(scoreMember.Member)

		if input.XX && !ok {
			continue
		}

		if input.NX && ok {
			continue
		}

		if ok && input.INCR {
			newScore = store.ToScoreF(sMember.Score) + store.ToScoreF(scoreMember.Score)
			scoreMember.Score = store.ToScoreS(newScore)
		}

		if !ok && input.INCR {
			newScore = store.ToScoreF(scoreMember.Score)
		}

		if !ok {
			add = add + 1
		}

		if ok {
			update = update + 1
		}

		kv.Scores.Add(scoreMember)
	}

	if (add+update) == 0 && input.INCR {
		return 0, fmt.Errorf("Failed to increment the score")
	}

	if (add + update) == 0 {
		return 0, nil
	}

	sort.Sort(kv.Scores)
	kv.Scores.BuildRank()

	c.buckets[bucketID].Set(kv)

	if input.INCR {
		return int(newScore), nil
	}

	if input.CH {
		return add + update, nil
	}

	return add, nil
}

func (c *Cache) ZCard(key string) (int, error) {
	bucketID := c.bucketID(key)

	kv, ok := c.buckets[bucketID].Get(key)
	if !ok {
		return 0, nil
	}

	if kv.ValueType != store.SetType {
		return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return kv.Scores.Len(), nil
}

func (c *Cache) ZCount(key, min, max string) (int, error) {
	bucketID := c.bucketID(key)

	kv, ok := c.buckets[bucketID].Get(key)
	if !ok {
		return 0, nil
	}

	if kv.ValueType != store.SetType {
		return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	count := kv.Scores.Count(min, max)

	return count, nil
}

func (c *Cache) ZRange(key string, start, stop int) ([]store.ScoreMember, error) {
	bucketID := c.bucketID(key)

	kv, ok := c.buckets[bucketID].Get(key)
	if !ok {
		return []store.ScoreMember{}, nil
	}

	if kv.ValueType != store.SetType {
		return []store.ScoreMember{}, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return kv.Scores.Range(start, stop), nil
}

func (c *Cache) Save(fileName string) (string, error) {
	dir, _ := filepath.Split(fileName)
	db := store.DB{
		Files: make(map[string]string),
	}

	for id, bucket := range c.buckets {
		name, err := bucket.Save(dir)
		if err != nil {
			return "", fmt.Errorf("Failed to save bucket: %d entries to file. err: %v", id, err)
		}

		db.Files[strconv.Itoa(id)] = name
	}

	jsonBytes, err := json.Marshal(db)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal db. err: %v", err)
	}

	if err := ioutil.WriteFile(fileName, jsonBytes, 0644); err != nil {
		return "", fmt.Errorf("Failed to write cache items to file: %s. err: %v", fileName, err)
	}

	return "OK", nil
}

func (c *Cache) Load(fileName string) (string, error) {
	if utils.FileExists(fileName) {
		return "!OK", nil
	}

	mBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("Failed to read db items from file. err: %v", err)
	}

	db := store.DB{
		Files: make(map[string]string),
	}

	if err := json.Unmarshal(mBytes, &db); err != nil {
		return "", fmt.Errorf("Failed to unmarshal db items. err: %v", err)
	}

	for id, file := range db.Files {
		bucketID, err := strconv.Atoi(id)
		if err != nil {
			return "", fmt.Errorf("Unable to identify bucket id. err: %v", err)
		}

		if _, err := c.buckets[bucketID].Load(file); err != nil {
			return "", fmt.Errorf("Failed to load bucket: %d err: %v", bucketID, err)
		}
	}

	return "OK", nil
}

func (c *Cache) Cleanup() {
	for _, bucket := range c.buckets {
		bucket.Cleanup()
	}
}
