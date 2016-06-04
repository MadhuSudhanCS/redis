package store

import (
	"fmt"
	"testing"
	"time"
)

func TestNewBucket(t *testing.T) {
	b := NewBucket(1)
	if b == nil {
		t.Fatalf("Failed to create new bucket")
	}
}

func TestSetGet(t *testing.T) {
	b := NewBucket(1)
	kv := KeyValue{
		Key:       "mykey",
		Value:     []byte("myvalue"),
		ValueType: StringType,
	}

	b.Set(kv)

	if _, ok := b.Get("mykey"); !ok {
		t.Fatalf("Expected key to be found in bucket")
	}
}

func TestSetGetWithExpiration(t *testing.T) {
	b := NewBucket(1)

	kv := KeyValue{
		Key:        "mykey",
		Value:      []byte("myvalue"),
		ValueType:  StringType,
		Expiration: time.Now().Add(5 * time.Second),
	}

	b.Set(kv)

	if _, ok := b.Get("mykey"); !ok {
		t.Fatalf("Expected key to be found in bucket")
	}

	time.Sleep(5 * time.Second)
	if _, ok := b.Get("mykey"); ok {
		t.Fatalf("Expected nil value as the key is expired")
	}

	if _, ok := b.Get("myinvalidkey"); ok {
		t.Fatalf("Expected nil value as the key is invaid")
	}
}

func TestSetGetWithInvalidKey(t *testing.T) {
	b := NewBucket(1)

	if _, ok := b.Get("myinvalidkey"); ok {
		t.Fatalf("Expected nil value as the key is invaid")
	}
}

func TestCleanUp(t *testing.T) {
	b := NewBucket(1)

	kv := KeyValue{
		Key:        "mykey",
		Value:      []byte("myvalue"),
		ValueType:  StringType,
		Expiration: time.Now().Add(5 * time.Second),
	}

	b.Set(kv)

	if 1 != len(b.item) {
		t.Fatalf("Expected 1 item in bucket")
	}

	time.Sleep(5 * time.Second)

	b.Cleanup()

	if _, ok := b.Get("mykey"); ok {
		t.Fatalf("Expected item to be deleted in bucket")
	}
}

func generateItems(b *Bucket, count int) {
	key := "mykey_%d"

	for i := 0; i < count; i++ {
		kv := KeyValue{
			Key:       fmt.Sprintf(key, i),
			Value:     []byte("myvalue"),
			ValueType: StringType,
		}
		b.Set(kv)
	}

	s := SortedScores{
		Scores:  make(map[string]ScoreMember),
		Members: []string{},
	}

	s.Add(score1)
	s.Add(score4)

	for i := count; i < count*2; i++ {
		kv := KeyValue{
			Key:       fmt.Sprintf(key, i),
			ValueType: SetType,
			Scores:    s,
		}
		b.Set(kv)
	}
}

func TestSave(t *testing.T) {
	b := NewBucket(1)
	generateItems(b, 100)

	fileName, err := b.Save("/tmp/")
	if err != nil {
		t.Fatalf("failed to save bucket items. err: %v", err)
	}

	eFileName := "/tmp/bucket_1.json"
	if fileName != eFileName {
		t.Fatalf("expected: %s found: %s", eFileName, fileName)
	}
}

func TestLoad(t *testing.T) {
	b := NewBucket(1)
	generateItems(b, 100)

	fileName, err := b.Save("/tmp/")
	if err != nil {
		t.Fatalf("failed to save bucket items. err: %v", err)
	}

	eFileName := "/tmp/bucket_1.json"
	if fileName != eFileName {
		t.Fatalf("expected: %s found: %s", eFileName, fileName)
	}

	if _, err := b.Load(fileName); err != nil {
		t.Fatalf("failed to load bucket items. err: %v", err)
	}

	for key, _ := range b.item {
		if _, ok := b.Get(key); !ok {
			t.Fatalf("expected %s to be present in bucket", key)
		}
	}
}

func BenchmarkSave(b *testing.B) {
	bucket := NewBucket(1)
	generateItems(bucket, 1000000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fileName, err := bucket.Save("/tmp/")
		if err != nil {
			b.Fatalf("failed to save bucket items. err: %v", err)
		}

		eFileName := "/tmp/bucket_1.json"
		if fileName != eFileName {
			b.Fatalf("expected: %s found: %s", eFileName, fileName)
		}
	}
}

func BenchmarkLoad(b *testing.B) {
	bucket := NewBucket(1)
	generateItems(bucket, 1000000)

	fileName, err := bucket.Save("/tmp/")
	if err != nil {
		b.Fatalf("failed to save bucket items. err: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := bucket.Load(fileName); err != nil {
			b.Fatalf("failed to load bucket items. err: %v", err)
		}
	}
}
