package main

import (
	"crypto/sha256"
	"os"
	"testing"
)

func TestLineParse(t *testing.T){
	testcases := []struct {
		in      string
		outHash string
		outFile string
	}{
		{"486f25a109aae34cee78de7647047a289ea6a7c18d4f381bd20fbc9feb26da3c  with a  double     space.  ", "486f25a109aae34cee78de7647047a289ea6a7c18d4f381bd20fbc9feb26da3c", "with a  double     space.  "},
		{"498cf6bafdd6ced537e0433f28c06a78504c74f684148acfb51b5a04156c9e4f  trailing spaces  ", "498cf6bafdd6ced537e0433f28c06a78504c74f684148acfb51b5a04156c9e4f", "trailing spaces  "},
		{"d4c018b244054eb557b7bbb3011e859b705654253eafe52d5f0aa735ff739fb7   leading space.txt", "d4c018b244054eb557b7bbb3011e859b705654253eafe52d5f0aa735ff739fb7", " leading space.txt"},
		{"da8748d73e686de80325716885dc4c924fb34746fdd253f677009ebdeefcca01   with spaces.jpg", "da8748d73e686de80325716885dc4c924fb34746fdd253f677009ebdeefcca01", " with spaces.jpg"},
	}
	for _, tc := range testcases {
		hashhex, fname, err := parseLine(tc.in)
		if err != nil {
			t.Error(err)
		}
		if tc.outFile != fname {
			t.Errorf("hash file parse failed, expected filename '%s' got '%s'", tc.outFile, fname)
		}
		if tc.outHash != hashhex {
			t.Errorf("hash file parse failed, expected hash '%s' got '%s'", tc.outHash, hashhex)
		}
	}
}

const testFile ="testdata/testdata.dat"

func TestFileHash(t *testing.T) {
	testcases := []struct {
		algo    string
		outHash string
	}{
		{	"md5", "7b4db0bf2a2c1d2e24f7784601f65499"},
		{	"sha1", "b4fab5688a47de865b31b0de4a2c8b8c05d81279"},
		{	"sha224", "0da6da9aa95a6cc6aec8b2349db7bd570a115f7a1b712356dc63e24d"},
		{	"sha256", "e2540a96c617b4413167a284d83b1d9f7553866685b5395df250faf47328727a"},
		{	"sha384", "ec9634f5e5be1f232d7479db811d93668c08dfaed3b0b7870c8d04477faf47bb554edbb2938680da9748472840826a93"},
		{	"sha512", "7ff2d404e5170c63813a9067abbe2edadc64403b3d109977a6e9500bdacaaa9230843565537a1a8ed07202c962008ab79e94873fd7ccde01c30a39a9d4785ddb"},
		{	"sha512224", "d947e34afad4640c6a74b28b6b25ea84ca3ee0bc808e347bb2141a2c"},
		{	"sha512256", "8794c453b690b7e47d0a032b9aede84eaa5d988a148860b16a9dab0f0533500b"},
	}
	for _, tc := range testcases {
		hasher := NewHasher(tc.algo)
		if hasher == nil {
			t.Errorf("couldn't get a hasher for %s", tc.algo)
		} else {
			hexhash,err := HashFile(hasher, testFile)
			if err != nil {
				t.Errorf("error %s hashing %s: %v", tc.algo, testFile, err)
			}
			if hexhash != tc.outHash {
				t.Errorf("wrong %s hash: expected %s got %s", tc.algo, tc.outHash, hexhash)
			}
		}
	}
}

const checksumFile = "testdata/checksums.sha256"

func TestCheckFromFile(t *testing.T) {
	// Suppress stdout for the duration of this test
	stdout := os.Stdout
	defer func () {
		os.Stdout = stdout
	}()
	os.Stdout,_ = os.Open(os.DevNull)
	fails, err := CheckHashes(sha256.New(), checksumFile)
	if err != nil {
		t.Errorf("error checking checksums in %s: %v", checksumFile, err)
	}
	if fails != 0 {
		t.Errorf("%d error(s) checking checksums in %s, expected 0", fails, checksumFile)
	}
}

const badFile = "testdata/checksums.bad.sha256"

func TestFailFromFile(t *testing.T) {
	// Suppress stdout for the duration of this test
	stdout := os.Stdout
	defer func () {
		os.Stdout = stdout
	}()
	os.Stdout,_ = os.Open(os.DevNull)
	fails, err := CheckHashes(sha256.New(), badFile)
	if err != nil {
		t.Errorf("error checking checksums in %s: %v", badFile, err)
	}
	if fails != 4 {
		t.Errorf("%d error(s) checking checksums in %s, expected 4", fails, checksumFile)
	}
}

func ExampleMakeAllHashes() {
	MakeAllHashes(NewHasher("sha256"), []string{"testdata/multi"})
	// Output:
  // 2776461b5b9fc6baaae69f2e6367afeb2449c383a55de5ea5b936d6ecb9bbf84  testdata/multi/one.dat
	// 383f36c987fcadd6c033b0645a874db256755fd66fff9c486637123debdff6da  testdata/multi/subdir/zz.dat
	// 917ce65312b9ad00aaa4dce4fe4af5582ce1398fdefd33e787ad5db6dadd1844  testdata/multi/zz.dat
}

func TestCheckAllHashes(t *testing.T) {
	fails := CheckAllHashes(NewHasher("sha256"), []string{"testdata/multi.sha256"})
	if fails != 0 {
		t.Errorf("got %d failures expected 0", fails)
	}
}