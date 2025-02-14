package auth

import (
	"testing"
)

func TestAuth(t *testing.T) {
	passwords := []string{"short", "loooooooooooooooooooooooooong", "n0n3s3n5e"}
	hashes := make([]string, len(passwords))
	for i, p := range passwords {
		h, err := HashPassword(p)
		if err != nil {
			t.Fatalf("%v test failed: could not hash %s", i+1, p)
		}
		hashes[i] = h
	}

	for i := 0; i < len(passwords); i++ {
		err := CheckPasswordHash(passwords[i], hashes[i])
		if err != nil {
			t.Fatalf("%v test failed: check hash failed %s", i+1, err)
		}
	}

	// expected to fail
	hash, err := HashPassword("ez")
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = CheckPasswordHash("beans", hash)
	if err == nil {
		t.Fatalf("%s", err)
	}
}
