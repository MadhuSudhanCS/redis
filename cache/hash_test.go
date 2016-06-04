package cache

import (
	"testing"
)

func TestHash(t *testing.T) {
	hv := Hash("mykey")

	if hv != Hash("mykey") {
		t.Fatalf("Hash is not idemepotent")
	}

}
