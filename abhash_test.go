package abhash

import (
	"bytes"
	"testing"
)

func TestABHashConsistency(t *testing.T) {
	data := []byte("This is a test input for abhash!")
	sohp := 4
	hpa := 3

	hash1 := ABHash(data, sohp, hpa)
	hash2 := ABHash(data, sohp, hpa)

	if !bytes.Equal(hash1, hash2) {
		t.Errorf("Expected the same hash for identical input, got %x and %x", hash1, hash2)
	}
}

func TestABHashDifferentParts(t *testing.T) {
	// Verify that changing one part only affects its corresponding token.
	data := []byte("abcdefghij") // 10 bytes
	sohp := 3
	hpa := 3
	// Expected split:
	// Part1: data[0:3] -> "abc"
	// Part2: data[3:6] -> "def"
	// Part3: data[6:]  -> "ghij"
	hashOriginal := ABHash(data, sohp, hpa)

	// Modify a byte in part2 (e.g., change 'e' to 'x')
	modified := make([]byte, len(data))
	copy(modified, data)
	modified[4] = 'x' // Change in part2

	hashModified := ABHash(modified, sohp, hpa)

	// Tokens for parts 1 and 3 should remain unchanged.
	if !bytes.Equal(hashOriginal[0:sohp], hashModified[0:sohp]) {
		t.Errorf("Part1 token changed unexpectedly: original %x, modified %x",
			hashOriginal[0:sohp], hashModified[0:sohp])
	}
	if !bytes.Equal(hashOriginal[2*sohp:3*sohp], hashModified[2*sohp:3*sohp]) {
		t.Errorf("Part3 token changed unexpectedly: original %x, modified %x",
			hashOriginal[2*sohp:3*sohp], hashModified[2*sohp:3*sohp])
	}

	// Token for part2 should differ.
	if bytes.Equal(hashOriginal[sohp:2*sohp], hashModified[sohp:2*sohp]) {
		t.Errorf("Part2 token did not change as expected")
	}
}

func TestABHashShortData(t *testing.T) {
	// Test with input data shorter than required for hpa parts.
	data := []byte("abc")
	sohp := 2
	hpa := 3
	hash := ABHash(data, sohp, hpa)
	expectedLen := hpa * sohp
	if len(hash) != expectedLen {
		t.Errorf("Expected hash length %d, got %d", expectedLen, len(hash))
	}
}
