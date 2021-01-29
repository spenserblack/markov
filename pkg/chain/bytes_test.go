package chain

import (
	"crypto/sha1"
	"testing"
)

// TestNewPrefix1 calls generator.NewBytesChain with prefixLen=1 and confirms that the
// returned generator.BytesGenerator struct contains the expected contents.
func TestNewPrefix1(t *testing.T) {
	var feed [][][]byte = [][][]byte{
		{[]byte("Hello,"), []byte("World!")},
		{[]byte("Hello,"), []byte("Go!")},
	}
	generator, err := NewBytesChain(feed, 1)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	if len(generator.chainStarters) != 2 ||
		string(generator.chainStarters[0][0]) != "Hello," ||
		string(generator.chainStarters[1][0]) != "Hello," {
		t.Fatalf(`generator.chainStarters = %v, want ["Hello," "Hello,"]`, generator.chainStarters)
	}

	h := sha1.New()
	h.Write([]byte("Hello,"))
	helloHash := string(h.Sum(nil))

	if bytes := generator.chain[helloHash]; len(bytes) != 2 ||
		!chainContains(bytes, "World!") || !chainContains(bytes, "Go!") {
		t.Fatalf(`generator.chain[%q] = %v, want ["World!" "Go!"]`, helloHash, bytes)
	}

	h.Reset()
	h.Write([]byte("Go!"))
	goHash := string(h.Sum(nil))

	h.Reset()
	h.Write([]byte("World!"))
	worldHash := string(h.Sum(nil))

	for _, hash := range []string{goHash, worldHash} {
		if bytes := generator.chain[hash]; len(bytes) != 1 || bytes[0] != nil {
			t.Fatalf(`generator.chain[%q] = %v, want [nil]`, hash, bytes)
		}
	}
}

// TestNewPrefix2 calls generator.NewBytesChain with prefixLen=2 and confirms that the
// returned generator.BytesGenerator struct contains the expected contents.
func TestNewPrefix2(t *testing.T) {
	var feed [][][]byte = [][][]byte{
		{[]byte("Hello"), []byte(","), []byte("World!")},
		{[]byte("Hello"), []byte("."), []byte("Go!")},
	}
	generator, err := NewBytesChain(feed, 2)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	// TODO use contains for concurrency
	if len(generator.chainStarters) != 2 ||
		string(generator.chainStarters[0][0]) != "Hello" ||
		string(generator.chainStarters[1][0]) != "Hello" ||
		!chainContains(
			[][]byte{
				generator.chainStarters[0][1],
				generator.chainStarters[1][1],
			},
			",",
		) ||
		!chainContains(
			[][]byte{
				generator.chainStarters[0][1],
				generator.chainStarters[1][1],
			},
			".",
		) {
		t.Fatalf(`generator.chainStarters = %v, want ["Hello," "Hello."]`, generator.chainStarters)
	}

	h := sha1.New()
	h.Write([]byte("Hello"))
	h.Write([]byte(","))
	helloHashComma := string(h.Sum(nil))

	if bytes := generator.chain[helloHashComma]; len(bytes) != 1 || string(bytes[0]) != "World!" {
		t.Fatalf(`generator.chain[%q] = %v, want ["World!"]`, helloHashComma, bytes)
	}

	h.Reset()
	h.Write([]byte("Hello"))
	h.Write([]byte("."))
	helloHashPeriod := string(h.Sum(nil))

	if bytes := generator.chain[helloHashPeriod]; len(bytes) != 1 || string(bytes[0]) != "Go!" {
		t.Fatalf(`generator.chain[%q] = %v, want ["Go!"]`, helloHashComma, bytes)
	}

	h.Reset()
	h.Write([]byte("."))
	h.Write([]byte("Go!"))
	goHash := string(h.Sum(nil))

	h.Reset()
	h.Write([]byte(","))
	h.Write([]byte("World!"))
	worldHash := string(h.Sum(nil))

	for _, hash := range []string{goHash, worldHash} {
		if bytes := generator.chain[hash]; len(bytes) != 1 || bytes[0] != nil {
			t.Fatalf(`generator.chain[%q] = %v, want [nil]`, hash, bytes)
		}
	}
}

func chainContains(chainBranch [][]byte, value string) bool {
	for _, val := range chainBranch {
		if string(val) == value {
			return true
		}
	}
	return false
}
