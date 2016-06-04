package cache

import (
	"fmt"
	"github.com/madhusudhancs/redis/store"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	c := NewCache(2)
	if c == nil {
		t.Fatalf("Unable to create cache")
	}
}

func TestSetGet(t *testing.T) {
	c := NewCache(2)

	input := SetInput{
		Key:   "mykey",
		Value: "myvalue",
	}

	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to add key to cache. err: %v", err)
	}

	if _, err := c.Get("mykey"); err != nil {
		t.Fatalf("Failed to retrieve value from cache. err: %v", err)
	}
}

func TestSetWithExpiryTime(t *testing.T) {
	c := NewCache(2)
	key := "mykey"
	value := "myvalue"

	input := SetInput{
		Key:    key,
		Value:  value,
		Expiry: 5 * time.Second,
	}

	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to add key to cache. err: %v", err)
	}

	v, err := c.Get(key)
	if err != nil {
		t.Fatalf("Failed to retrieve value from cache. err: %v", err)
	}

	if v != value {
		t.Fatalf("expected: %s found: %s", value, v)
	}

	time.Sleep(5 * time.Second)
	v, err = c.Get(key)
	if err != nil {
		t.Fatalf("Failed to retrieve value from cache. err: %v", err)
	}

	if v != "" {
		t.Fatalf("expected empty value")
	}
}

func TestSetOverwrite(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	input := SetInput{
		Key:    key,
		Value:  "myvalue",
		Expiry: 5 * time.Second,
	}

	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to add key to cache. err: %v", err)
	}

	time.Sleep(3 * time.Second)

	value := "newvalue"
	input.Value = value
	input.Expiry = 0

	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to update to cache. err: %v", err)
	}

	time.Sleep(6 * time.Second)

	v, err := c.Get(key)
	if err != nil {
		t.Fatalf("Failed to retrieve value from cache. err: %v", err)
	}

	if v == "" {
		t.Fatalf("expected empty value")
	}

	if v != value {
		t.Fatalf("expected: %s found: %s", value, v)
	}
}

func TestSetNX(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	input := SetInput{
		Key:    key,
		Value:  "myvalue",
		Expiry: 5 * time.Second,
		NX:     true,
	}

	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to add key to cache. err: %v", err)
	}

	if err := c.Set(input); err == nil {
		t.Fatalf("expected to throw error")
	}
}

func TestSetXX(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	input := SetInput{
		Key:    key,
		Value:  "myvalue",
		Expiry: 5 * time.Second,
		XX:     true,
	}

	if err := c.Set(input); err == nil {
		t.Fatalf("expected to throw error")
	}

	input.XX = false
	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to add key to cache. err: %v", err)
	}

	input.XX = true
	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to update to cache. err: %v", err)
	}
}

func TestGetInvalidKey(t *testing.T) {
	c := NewCache(2)
	_, err := c.Get("myinvalidkey")
	if err != nil {
		t.Fatalf("unexpected error. err: %v", err)
	}
}

func TestSetGetBit(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	input := SetInput{
		Key:   key,
		Value: "myvalue",
	}

	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to add key to cache. err: %v", err)
	}

	bitValue, err := c.SetBit(key, 0, 1)
	if err != nil {
		t.Fatalf("Unable to set bit value. err: %v", err)
	}

	if bitValue != 0 {
		t.Fatalf("expected bitvalue: %d found: %d", 0, bitValue)
	}

	value, err := c.Get(key)
	if err != nil {
		t.Fatalf("unable to retrive value for cache. err: %v", err)
	}

	expectedValue := []byte{237, 121, 118, 97, 108, 117, 101}
	foundValue := []byte(value)
	for i, v := range expectedValue {
		if v != foundValue[i] {
			t.Fatalf("expected: %s found: %s", v, foundValue[i])
		}
	}

	bitValue = c.GetBit(key, 17)
	if bitValue != 1 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}
}

func TestSetBit(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	if bitValue, _ := c.SetBit(key, 0, 0); bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}

	if bitValue, _ := c.SetBit(key, 1, 1); bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}

	if bitValue, _ := c.SetBit(key, 2, 0); bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}

	if bitValue, _ := c.SetBit(key, 3, 1); bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}

	if bitValue, _ := c.SetBit(key, 4, 0); bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}

	if bitValue, _ := c.SetBit(key, 5, 0); bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}

	if bitValue, _ := c.SetBit(key, 6, 1); bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}

	if bitValue, _ := c.SetBit(key, 7, 0); bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}

	value, err := c.Get(key)
	if err != nil {
		t.Fatalf("failed to get value from cache")
	}

	if value != "R" {
		t.Fatalf("Expected: %s found: %s", "R", value)
	}
}

func TestInvalidSetBit(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	if _, err := c.SetBit(key, 0, 10); err == nil {
		t.Fatalf("expected to throw error")
	}
}

func TestEmptyKeyGetBit(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	bitValue := c.GetBit(key, 10)
	if bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}
}

func TestGetBitOffsetOutofRange(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	if bitValue, _ := c.SetBit(key, 7, 0); bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}

	bitValue := c.GetBit(key, 10)
	if bitValue != 0 {
		t.Fatalf("expected: %d found: %d", 0, bitValue)
	}
}

func TestZAdd(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	input := SetInput{
		Key:   key,
		Value: "myvalue",
	}

	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to add key to cache. err: %v", err)
	}

	zinput := ZSetInput{
		Key:    key,
		Scores: []store.ScoreMember{},
	}

	//Test invalid data type
	if _, err := c.ZAdd(zinput); err == nil {
		t.Fatalf("expected to throw error")
	}

	key = "myzkey"
	score := store.ScoreMember{
		Member: "GO",
		Score:  "1.5",
	}

	zinput.Key = key
	zinput.Scores = append(zinput.Scores, score)

	//Test Add
	value, err := c.ZAdd(zinput)
	if err != nil {
		t.Fatalf("Unable to add zset to cache. err: %v", err)
	}

	if value != 1 {
		t.Fatalf("Failed to add score to key")
	}

	zinput.NX = true

	//Test add when not present
	value, err = c.ZAdd(zinput)
	if err != nil {
		t.Fatalf("Unable to add zset to cache. err: %v", err)
	}

	if value != 0 {
		t.Fatalf("expected not to add score")
	}

	zinput.XX = true
	zinput.NX = false
	zinput.CH = true

	//Test update with CH
	value, err = c.ZAdd(zinput)
	if err != nil {
		t.Fatalf("Unable to add zset to cache. err: %v", err)
	}

	if value != 1 {
		t.Fatalf("unable to update zset")
	}

	zinput.INCR = true

	//Test Incr with update
	value, err = c.ZAdd(zinput)
	if err != nil {
		t.Fatalf("Unable to add zset to cache. err: %v", err)
	}

	if value != 3 {
		t.Fatalf("unable to increment zset")
	}

	zinput.CH = false
	zinput.XX = false
	zinput.NX = false

	//Test incr
	value, err = c.ZAdd(zinput)
	if err != nil {
		t.Fatalf("Unable to add zset to cache. err: %v", err)
	}

	if value != 4 {
		t.Fatalf("unable to increment zset")
	}

	score.Member = "Java"
	score.Score = "1"
	zinput.Scores = append(zinput.Scores, score)

	//Test INCR with mulitple scores
	value, err = c.ZAdd(zinput)
	if err == nil {
		t.Fatalf("expected to throw error")
	}

	zinput.INCR = false

	//Test with mulitple scores
	value, err = c.ZAdd(zinput)
	if err != nil {
		t.Fatalf("failed to add scores. err: %v", err)
	}

	if value != 1 {
		t.Fatalf("Expected to add 1 score.")
	}

	zinput.NX = true

	//Test only add multiple scores
	value, err = c.ZAdd(zinput)
	if err != nil {
		t.Fatalf("failed to add scores. err: %v", err)
	}

	if value != 0 {
		t.Fatalf("Expected to add 0 score.")
	}

	zinput.NX = false
	zinput.XX = true
	zinput.CH = true

	//Test update multiple scores
	value, err = c.ZAdd(zinput)
	if err != nil {
		t.Fatalf("failed to add scores. err: %v", err)
	}

	if value != 2 {
		t.Fatalf("Expected to update 2 score.")
	}

	if _, err := c.Get(key); err == nil {
		t.Fatalf("expected to throw error")
	}
}

func TestZCardInvalidType(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	input := SetInput{
		Key:   key,
		Value: "myvalue",
	}

	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to add key to cache. err: %v", err)
	}

	if _, err := c.ZCard(key); err == nil {
		t.Fatalf("expected to throw error")
	}
}

func TestZCardEmptyKey(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	value, err := c.ZCard(key)
	if err != nil {
		t.Fatalf("unable to get cardinals for set. err: %v", err)
	}

	if value != 0 {
		t.Fatalf("expected 0 cardinal")
	}
}

func TestZCardValidKey(t *testing.T) {
	c := NewCache(2)
	key := "myzkey"
	score := store.ScoreMember{
		Member: "GO",
		Score:  "1.5",
	}

	zinput := ZSetInput{
		Key:    key,
		Scores: []store.ScoreMember{score},
	}

	value, err := c.ZAdd(zinput)
	if err != nil {
		t.Fatalf("Unable to add zset to cache. err: %v", err)
	}

	value, err = c.ZCard(key)
	if err != nil {
		t.Fatalf("unable to get cardinals for set. err: %v", err)
	}

	if value != 1 {
		t.Fatalf("expected 1 cardinal")
	}
}

func TestZCount(t *testing.T) {
	c := NewCache(2)
	key := "myzkey"
	score := store.ScoreMember{
		Member: "GO",
		Score:  "1.5",
	}

	zinput := ZSetInput{
		Key:    key,
		Scores: []store.ScoreMember{score},
	}

	score.Member = "Java"
	score.Score = "2"
	zinput.Scores = append(zinput.Scores, score)

	score.Member = "Python"
	score.Score = "2.2"
	zinput.Scores = append(zinput.Scores, score)

	score.Member = "C"
	score.Score = "2.2"
	zinput.Scores = append(zinput.Scores, score)

	if _, err := c.ZAdd(zinput); err != nil {
		t.Fatalf("Unable to add zset to cache. err: %v", err)
	}

	if count, _ := c.ZCount(key, store.MINSCORE, store.MAXSCORE); count != 4 {
		t.Fatalf("Expected: %d found: %d", 4, count)
	}

	if count, _ := c.ZCount(key, store.MINSCORE, "2.0"); count != 2 {
		t.Fatalf("Expected: %d found: %d", 2, count)
	}

	if count, _ := c.ZCount(key, "-1.0", "2.0"); count != 2 {
		t.Fatalf("Expected: %d found: %d", 2, count)
	}

	if count, _ := c.ZCount(key, "-1.0", store.MAXSCORE); count != 4 {
		t.Fatalf("Expected: %d found: %d", 4, count)
	}

	if count, _ := c.ZCount(key, "2.0", "2.2"); count != 3 {
		t.Fatalf("Expected: %d found: %d", 3, count)
	}

	if count, _ := c.ZCount(key, "2", "-2"); count != 0 {
		t.Fatalf("Expected: %d found: %d", 0, count)
	}
}

func TestZCountEmptyKey(t *testing.T) {
	c := NewCache(2)
	key := "myzkey"

	if count, _ := c.ZCount(key, "2.0", "2.2"); count != 0 {
		t.Fatalf("Expected: %d found: %d", 0, count)
	}
}

func TestZCountInvalidType(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	input := SetInput{
		Key:   key,
		Value: "myvalue",
	}

	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to add key to cache. err: %v", err)
	}

	if _, err := c.ZCount(key, store.MINSCORE, store.MAXSCORE); err == nil {
		t.Fatalf("expected to throw error")
	}
}

func TestZRangeEmptyKey(t *testing.T) {
	c := NewCache(2)
	key := "myzkey"

	if scores, _ := c.ZRange(key, 0, -1); len(scores) != 0 {
		t.Fatalf("Expected: %d found: %d", 0, len(scores))
	}
}

func TestZRangeInvalidType(t *testing.T) {
	c := NewCache(2)
	key := "mykey"

	input := SetInput{
		Key:   key,
		Value: "myvalue",
	}

	if err := c.Set(input); err != nil {
		t.Fatalf("Unable to add key to cache. err: %v", err)
	}

	if _, err := c.ZRange(key, 0, -1); err == nil {
		t.Fatalf("expected to throw error")
	}
}

func TestZRange(t *testing.T) {
	c := NewCache(2)
	key := "myzkey"
	score := store.ScoreMember{
		Member: "GO",
		Score:  "1.5",
	}

	zinput := ZSetInput{
		Key:    key,
		Scores: []store.ScoreMember{score},
	}

	score.Member = "Java"
	score.Score = "2"
	zinput.Scores = append(zinput.Scores, score)

	score.Member = "Python"
	score.Score = "2.2"
	zinput.Scores = append(zinput.Scores, score)

	score.Member = "C"
	score.Score = "2.2"
	zinput.Scores = append(zinput.Scores, score)

	if _, err := c.ZAdd(zinput); err != nil {
		t.Fatalf("Unable to add zset to cache. err: %v", err)
	}

	if scores, _ := c.ZRange(key, 0, -1); len(scores) != 4 {
		t.Fatalf("Expected: %d found: %d", 4, len(scores))
	}

	if scores, _ := c.ZRange(key, 0, 1); len(scores) != 2 {
		t.Fatalf("Expected: %d found: %d", 2, len(scores))
	}

	if scores, _ := c.ZRange(key, 0, 10); len(scores) != 4 {
		t.Fatalf("Expected: %d found: %d", 4, len(scores))
	}

	if scores, _ := c.ZRange(key, 2, 3); len(scores) != 2 {
		t.Fatalf("Expected: %d found: %d", 2, len(scores))
	}

	if scores, _ := c.ZRange(key, -2, 3); len(scores) != 2 {
		t.Fatalf("Expected: %d found: %d", 0, len(scores))
	}

	if scores, _ := c.ZRange(key, -2, -1); len(scores) != 2 {
		t.Fatalf("Expected: %d found: %d", 0, len(scores))
	}

	if scores, _ := c.ZRange(key, 0, 3); len(scores) != 4 {
		t.Fatalf("Expected: %d found: %d", 4, len(scores))
	}
}

func generateItems(c *Cache, count int) {
	key := "mykey_%d"
	value := "myvalue_%d"

	input := SetInput{}

	for i := 0; i < count; i++ {
		input.Key = fmt.Sprintf(key, i)
		input.Value = fmt.Sprintf(value, i)

		c.Set(input)
	}
}

func TestSave(t *testing.T) {
	c := NewCache(2)
	generateItems(c, 100)

	if _, err := c.Save("/tmp/dump.json"); err != nil {
		t.Fatalf("failed to save db to file. err: %v", err)
	}
}

func TestLoad(t *testing.T) {
	c := NewCache(2)
	key := "mykey_%d"
	value := "myvalue_%d"

	generateItems(c, 100)

	if _, err := c.Save("/tmp/dump.json"); err != nil {
		t.Fatalf("failed to save db to file. err: %v", err)
	}

	c = NewCache(2)
	if _, err := c.Load("/tmp/dump.json"); err != nil {
		t.Fatalf("failed to load db from file. err: %v", err)
	}

	v, err := c.Get(fmt.Sprintf(key, 0))
	if err != nil {
		t.Fatalf("failed to read key. err: %v", err)
	}

	eValue := fmt.Sprintf(value, 0)
	if strings.Compare(v, eValue) != 0 {
		t.Fatalf("expected: %s found: %s", eValue, v)
	}
}

func TestEmptyFileLoad(t *testing.T) {
	c := NewCache(2)

	if _, err := c.Load("/tmp/dump.json"); err != nil {
		t.Fatalf("failed to load db from file. err: %v", err)
	}
}

func BenchmarkSave(b *testing.B) {
	c := NewCache(10)
	generateItems(c, 1000000)

	for i := 0; i < b.N; i++ {
		if _, err := c.Save("/tmp/dump.json"); err != nil {
			b.Fatalf("failed to save db to file. err: %v", err)
		}
	}
}

func BenchmarkLoad(b *testing.B) {
	c := NewCache(10)
	generateItems(c, 1000000)

	if _, err := c.Save("/tmp/dump.json"); err != nil {
		b.Fatalf("failed to save db to file. err: %v", err)
	}

	for i := 0; i < b.N; i++ {
		c = NewCache(10)
		if _, err := c.Load("/tmp/dump.json"); err != nil {
			b.Fatalf("failed to load db from file. err: %v", err)
		}
	}
}

func BenchmarkCacheSet(b *testing.B) {
	c := NewCache(10)
	key := "mykey_%d"
	value := "myvalue_%d"

	input := SetInput{}

	for i := 0; i < b.N; i++ {
		input.Key = fmt.Sprintf(key, i)
		input.Value = fmt.Sprintf(value, i)

		if err := c.Set(input); err != nil {
			b.Fatalf("Unable to add key to cache. err: %v", err)
		}
	}
}

func BenchmarkCacheGet(b *testing.B) {
	c := NewCache(10)
	key := "mykey_%d"
	value := "myvalue_%d"

	input := SetInput{}

	for i := 0; i < b.N; i++ {
		input.Key = fmt.Sprintf(key, i)
		input.Value = fmt.Sprintf(value, i)

		if err := c.Set(input); err != nil {
			b.Fatalf("Unable to add key to cache. err: %v", err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v, err := c.Get(fmt.Sprintf(key, i))
		if err != nil {
			b.Fatalf("Failed to retrieve value from cache. err: %v", err)
		}

		iValue := fmt.Sprintf(value, i)
		if strings.Compare(v, iValue) != 0 {
			b.Fatalf("Expected: %s found: %s", iValue, v)
		}
	}
}

func BenchmarkSetBit(b *testing.B) {
	c := NewCache(10)
	key := "mykey_%d"

	for i := 0; i < b.N; i++ {
		k := rand.Intn(b.N)
		keyS := fmt.Sprintf(key, k)
		value := rand.Intn(2)
		offset := rand.Intn(100)

		if _, err := c.SetBit(keyS, offset, value); err != nil {
			b.Fatalf("Unable to set bit value. err: %v", err)
		}
	}
}

func BenchmarkGetBit(b *testing.B) {
	c := NewCache(10)
	key := "mykey_%d"

	for i := 0; i < b.N; i++ {
		keyS := fmt.Sprintf(key, i)
		value := rand.Intn(2)
		offset := rand.Intn(100)

		if _, err := c.SetBit(keyS, offset, value); err != nil {
			b.Fatalf("Unable to set bit value. err: %v", err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		keyS := fmt.Sprintf(key, i)
		offset := rand.Intn(100)

		c.GetBit(keyS, offset)
	}
}

func BenchmarkZAdd(b *testing.B) {
	c := NewCache(10)
	key := "mykey_%d"
	member := "member_%d"

	zinput := ZSetInput{
		Scores: []store.ScoreMember{},
	}

	score := store.ScoreMember{}

	for i := 0; i < b.N; i++ {
		zinput.Key = fmt.Sprintf(key, i)
		score.Member = fmt.Sprintf(member, i)
		score.Score = fmt.Sprintf("%f", rand.Float64())

		zinput.Scores = append([]store.ScoreMember{}, score)

		if _, err := c.ZAdd(zinput); err != nil {
			b.Fatalf("Unable to add zset to cache. err: %v", err)
		}
	}
}
