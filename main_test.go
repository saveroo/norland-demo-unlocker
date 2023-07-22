package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestComputePrefixTable(t *testing.T) {
	// Test case with needle: "abcabc"
	needle := []byte("abcabc")
	expectedPrefixTable := []int{0, 0, 0, 1, 2, 3}
	result := computePrefixTable(needle)
	if !reflect.DeepEqual(result, expectedPrefixTable) {
		t.Errorf("computePrefixTable() returned incorrect result for needle 'abcabc'.\nExpected: %v\nGot: %v", expectedPrefixTable, result)
	}

	// Test case with needle: "aaa"
	needle = []byte("aaa")
	expectedPrefixTable = []int{0, 1, 2}
	result = computePrefixTable(needle)
	if !reflect.DeepEqual(result, expectedPrefixTable) {
		t.Errorf("computePrefixTable() returned incorrect result for needle 'aaa'.\nExpected: %v\nGot: %v", expectedPrefixTable, result)
	}

	// Test case with needle: "abcdabcdabc"
	needle = []byte("abcdabcdabc")
	expectedPrefixTable = []int{0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7}
	result = computePrefixTable(needle)
	if !reflect.DeepEqual(result, expectedPrefixTable) {
		t.Errorf("computePrefixTable() returned incorrect result for needle 'abcdabcdabc'.\nExpected: %v\nGot: %v", expectedPrefixTable, result)
	}
}

func TestFindBytes(t *testing.T) {
	originalBytes := []byte("0123456")
	targetBytes := []byte("6")

	// Test scenario where the targetBytes are present in the originalBytes
	expectedIndex := 6
	result := findBytes(originalBytes, targetBytes)
	if result != expectedIndex {
		t.Errorf("findBytes() did not return the correct index. Expected: %d, Got: %d", expectedIndex, result)
	}

	// Test scenario where the targetBytes are not present in the originalBytes
	targetBytes = []byte("missing")
	expectedIndex = -1
	result = findBytes(originalBytes, targetBytes)
	if result != expectedIndex {
		t.Errorf("findBytes() did not return the correct index. Expected: %d, Got: %d", expectedIndex, result)
	}
}

func TestPatchBinaryKMP(t *testing.T) {
	dump := []byte("0123456789")
	target := []byte("5678")
	patch := []byte("0000")
	expected := []byte("0123400009")

	_ = patchBinaryKMP(dump, target, patch)
	if string(dump) != string(expected) {
		t.Errorf("patchBinaryKMP() did not return the correct result. Expected: %s, Got: %s", string(expected), dump)
	}
}

func TestPatcherWrapper(t *testing.T) {
	dump := []byte("0123456789")
	target := []byte("5678")
	patch := []byte("0000")
	expected := []byte("0123400009")

	_ = patcherWrapper(dump, target, patch)
	if string(dump) != string(expected) {
		t.Errorf("patchBinaryKMP() did not return the correct result. Expected: %s, Got: %s", string(expected), dump)
	}
}

func TestIsBinaryPatched(t *testing.T) {
	// Create a temporary file to test with
	tmpfile, err := os.CreateTemp("./", "test_file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	fmt.Print(tmpfile.Name())
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpfile.Name())

	// Write some content to the temporary file
	content := []byte("This is some content with a pattern to be matched.")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	defer func(tmpfile *os.File) {
		err := tmpfile.Close()
		if err != nil {
			log.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tmpfile)

	// Test scenario where the replacementBytes are not present in the file
	replacementBytes := []byte("xyz")
	result := isBinaryPatched(tmpfile.Name(), replacementBytes)
	if result != false {
		t.Errorf("isBinaryPatched() did not return the correct result. Expected: false, Got: %v", result)
	}

	// Test scenario where the replacementBytes are present in the file
	replacementBytes = []byte("pattern")
	result = isBinaryPatched(tmpfile.Name(), replacementBytes)
	if result != true {
		t.Errorf("isBinaryPatched() did not return the correct result. Expected: true, Got: %v", result)
	}
}

func TestComputeFileHash(t *testing.T) {
	// Create a temporary file to test with
	tmpfile, err := os.CreateTemp("./", "test_file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpfile.Name())

	// Write some content to the temporary file
	content := []byte("This is some content to be hashed.")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	defer func(tmpfile *os.File) {
		err := tmpfile.Close()
		if err != nil {
			log.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tmpfile)

	// Get the expected SHA-256 hash of the content
	expectedHash := sha256.Sum256(content)
	expectedHashString := hex.EncodeToString(expectedHash[:])

	// Call the function under test
	hash, err := computeFileHash(tmpfile.Name())
	if err != nil {
		t.Fatalf("computeFileHash() returned an error: %v", err)
	}

	// Verify if the computed hash matches the expected hash
	if hash != expectedHashString {
		t.Errorf("computeFileHash() returned incorrect hash. Expected: %s, Got: %s", expectedHashString, hash)
	}
}

func TestGetFileVersion(t *testing.T) {
	// Get the current working directory
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current working directory: %v", err)
		return
	}

	// Assuming your file is in the current working directory
	fileName := "./Norland.exe"
	filePath := filepath.Join(currentWorkingDir, fileName)

	fileVersion, err := GetFileVersion(filePath)
	if err != nil {
		t.Errorf("Error getting file version: %v", err)
		return
	}

	// Add assertions as needed
	if fileVersion == "" {
		t.Errorf("File version should not be empty")
	}
}
