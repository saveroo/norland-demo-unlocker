package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var verbose = false                 // Set to true to see more details
var targetVersion = "0.2912.3227.0" // The version of the game that this patcher is compatible with
var targetName = "Norland.exe"      // The name of the game executable

func getCurrentWorkingDirFile(targetExecutable string) (string, error) {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	bin, err := filepath.Abs(filepath.Join(currentWorkingDir, targetExecutable))
	if err != nil {
		return "", err
	}
	return bin, nil
}

func main() {
	fmt.Println(`
/**
* ███╗|||██╗|██████╗|██████╗|██╗||||||█████╗|███╗|||██╗██████╗||||||||
* ████╗||██║██╔═══██╗██╔══██╗██║|||||██╔══██╗████╗||██║██╔══██╗|||||||
* ██╔██╗|██║██║|||██║██████╔╝██║|||||███████║██╔██╗|██║██║||██║|||||||
* ██║╚██╗██║██║|||██║██╔══██╗██║|||||██╔══██║██║╚██╗██║██║||██║|||||||
* ██║|╚████║╚██████╔╝██║||██║███████╗██║||██║██║|╚████║██████╔╝|||||||
* ██████╗═███████╗███╗═╝|███╗╚██████╗╚═╝||╚═╝╚═╝||╚═══╝╚═════╝||||||||
* ██╔══██╗██╔════╝████╗|████║██╔═══██╗||||||||||||||||||||||||||||||||
* ██║||██║█████╗||██╔████╔██║██║|||██║||||||||||||||||||||||||||||||||
* ██║||██║██╔══╝||██║╚██╔╝██║██║|||██║||||||||||||||||||||||||||||||||
* ██████╔╝███████╗██║|╚═╝|██║╚██████╔╝||||||||||||||||||||||||||||||||
* ██╗═══██╗███╗══╝██╗██╗||╚═╝|██████╗||██████╗██╗||██╗███████╗██████╗|
* ██║|||██║████╗||██║██║|||||██╔═══██╗██╔════╝██║|██╔╝██╔════╝██╔══██╗
* ██║|||██║██╔██╗|██║██║|||||██║|||██║██║|||||█████╔╝|█████╗||██████╔╝
* ██║|||██║██║╚██╗██║██║|||||██║|||██║██║|||||██╔═██╗|██╔══╝||██╔══██╗
* ╚██████╔╝██║|╚████║███████╗╚██████╔╝╚██████╗██║||██╗███████╗██║||██║
* |╚═════╝|╚═╝||╚═══╝╚══════╝|╚═════╝||╚═════╝╚═╝||╚═╝╚══════╝╚═╝||╚═╝
* ||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||
* || *For education purpose only, Buy the game once it comes out.
* ||
* || Title: Norland - Story Generating Strategy - Demo
* || Date: 7/21/23
* || Size: 61.5 MB
* || Version: 0.2912.3227.0
* || Dev: svr
* ||
* || *For education purpose only, Buy the game once it comes out.
* ||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||
**/
`)

	// Loc
	binPath, err := getCurrentWorkingDirFile(targetName)
	if err != nil {
		fmt.Println(">> Error getting installation path:", err)
		promptExit()
		return
	}

	printProgress("\r\nBinary path: " + binPath)
	backupPath := binPath + ".bak"

	// Essential checks
	f, err := fileExists(binPath)
	if err != nil {
		fmt.Println(">> Error checking file:", err)
		promptExit()
		return
	}

	if !f {
		fmt.Println(">> Make sure the game executable (Norland.exe) is in the same directory as this patcher.")
		promptExit()
		return
	}

	// Check Version
	binVersion, err := GetFileVersion(binPath)
	if err != nil {
		fmt.Println(">> Error getting file version:", err)
		promptExit()
		return
	}
	if binVersion != targetVersion {
		printProgress(">> This patcher is only compatible with version " + targetVersion + " of the game.")
		promptExit()
		return
	}
	printProgress(">> " + binVersion + " Binary version match!. Proceeding...\n")

	// Prep the bytes
	binaryDump, err := os.ReadFile(binPath)
	if err != nil {
		fmt.Println(">> Error reading binary:", err)
		promptExit()
	}

	targetCode := []byte{
		0xB0, 0x01, 0x83, 0xF8, 0xFE,
		0x75, 0x05, 0x40, 0x32, 0xFF,
		0xEB, 0x06, 0x85, 0xC0, 0x40,
		0x0F, 0x99, 0xC7, 0x8B, 0x4D, 0x8C,
	}
	shellCode := []byte{
		0xB0, 0x01, 0x83, 0xF8, 0xFE,
		0x74, 0x05, 0x40, 0x32, 0xFF,
		0xEB, 0x06, 0x85, 0xC0, 0x40,
		0x0F, 0x99, 0xC7, 0x8B, 0x4D, 0x8C,
	}

	// TODO: Another day... another sunset..
	// jmp := []byte{ // JMP addr+off
	//	0xB0, 0x01, 0x83, 0xF8, 0xFE,
	//	0xE9, 0x1D, 0x8B, 0x56, 0xFF,
	//	0xEB, 0x06, 0x85, 0xC0, 0x40,
	//	0x0F, 0x99, 0xC7, 0x8B, 0x4D, 0x8C,
	//}
	// tramp := []byte{ // jz, xor dil, dil, ret jmp, pad
	//	0x74, 0x05, 0x40, 0x32, 0xFF,
	//	0xEB, 0x06, 0x85, 0xC0, 0x40,
	//	0xE9, 0xD5, 0x74, 0x74, 0xA9,
	//	0x90, 0x00, 0x90, 0x00, 0x90,
	//	0x00, 0x90, 0x00, 0x00, 0x90,
	//	0x90, 0x90, 0x90, 0x00, 0x90,
	//}

	// Let's start
	printProgress("\nStep 1: Checking if the binary is patched...")
	isPatched := isBinaryPatched(binPath, shellCode)
	if isPatched {
		printProgress("-The binary is already patched.")
		printProgress("-[?] Do you want to reverse the patch?? (y/n):")

		var shouldReversePatch string
		_, err = fmt.Scan(&shouldReversePatch)
		if err != nil {
			return
		}
		if shouldReversePatch == "n" || shouldReversePatch == "N" {
			printProgress(">> Exiting..")
			return
		} else {
			verbosePrompt()

			// reverse
			reversed := patcherWrapper(binaryDump, shellCode, targetCode)

			// reverse shellcode patch
			//restoreCave := make([]byte, len(shellCode))
			//copy(restoreCave, bytes.Repeat([]byte{0x00}, len(shellCode)))
			//patcherWrapper(binaryDump, shellCode, restoreCave)

			if !reversed {
				fmt.Println("-Reversing the patch fail!")
				return
			}
			writePatchedBinary(binPath, binaryDump)
			fmt.Println("-Patch reversed to its original form successfully!")
			return
		}
	} else {
		fmt.Println("-The binary is not patched.")
	}

	//isShellCodeExists := isBinaryPatched(binPath, shellCode)
	//if isPatched && isShellCodeExists {
	//	a
	//}

	// Check if the backup file exists
	backupFileExists, err := fileExists(backupPath)
	if err != nil {
		fmt.Println("-Error checking backup:", err)
		return
	}

	// no need to restore if we could reverse it.
	//if backupFileExists && isPatched {
	//	// Ask the user whether to reverse the patch or patch again
	//	fmt.Println("Backup of the binary found.")
	//	fmt.Println("Do you want to restore from backup? (y/n):")
	//	var restoreBackupOption string
	//	fmt.Scan(&restoreBackupOption)
	//	if restoreBackupOption == "y" || restoreBackupOption == "Y" {
	//		restoreBackup(binPath, backupPath)
	//		return
	//	}
	//}

	// Step 2: Backup the Binary
	printProgress("\nStep 2: Creating a backup of the binary... if it doesn't exist already.")
	if !backupFileExists {
		err = copyFile(binPath, backupPath)
		if err != nil {
			fmt.Println("Error creating backup:", err)
			return
		}
		fmt.Println("Backup created successfully.")
	}

	// Step 3: Locate the Target Opcode and Step 4: Perform the Patch
	printProgress("\nStep 3: Patching the binary...")
	verbosePrompt()

	// TODO: let it be... this for another day
	//relativeOffset := (caveOffset) - (baseOffset + 5)
	//binary.LittleEndian.PutUint32(shellCode[6:10], uint32(relativeOffset))

	patchSuccessful := patcherWrapper(binaryDump, targetCode, shellCode)
	if !patchSuccessful {
		fmt.Println("-Patching seems fail.. exiting.")
		return
	}

	// Step 4: Pack the Patched Binary
	printProgress("\nStep 4: Packing the patched binary...")
	writePatchedBinary(binPath, binaryDump)

}

func verbosePrompt() {
	printProgress("-[?] verbose mode?? (y/n):")
	var isVerbose string
	_, err := fmt.Scan(&isVerbose)
	if err != nil {
		return
	}
	if isVerbose == "y" || isVerbose == "Y" {
		verbose = true
		printProgress(">> Verbose mode enabled.")
	}
}

func promptExit() {
	// gracefully exit
	printProgress(">> Press 'Enter' to exit..")
	_, err := fmt.Scanln()
	if err != nil {
		os.Exit(0)
	}
}

//func printProgress(caption string) {
//	rand.New(rand.NewSource(time.Now().UnixNano()))
//	for _, c := range caption {
//		fmt.Print(string(c))
//		time.Sleep(time.Duration(15) * time.Millisecond) // Simulate progress
//	}
//	fmt.Print("\n")
//}

//func assemblingTheTrampoline() {}

func writePatchedBinary(filePath string, originalBytes []byte) {
	time.Sleep(time.Duration(1) * time.Second) // Simulate progress
	err := os.WriteFile(filePath, originalBytes, 0644)
	if err != nil {
		fmt.Println("-Error writing patched binary:", err)
		os.Exit(0)
	}
	fmt.Println("-Binary patched successfully!")

	// quick protoype for testing
	printProgress(">> Verifying the patch...")
	// Compute hashes for the original file and the patched file
	originalHash, err := computeFileHash(filePath + ".bak")
	if err != nil {
		fmt.Println("Error computing hash for the original file:", err)
		return
	}

	patchedHash, err := computeFileHash(filePath)
	if err != nil {
		fmt.Println("Error computing hash for the patched file:", err)
		return
	}

	// Compare the hashes
	if originalHash == patchedHash {
		printProgress("Original SHA256: " + originalHash)
		printProgress("Patched SHA256: " + patchedHash)
		printProgress(">> Sum hashes are equal. Patch are reverted")
		return
	} else {
		printProgress("Original SHA256: " + originalHash)
		printProgress("Patched SHA256: " + patchedHash)
		fmt.Println(">> Sum hashes are not equal. Demo should be patched successfully.")
		return
	}
}

// KMP Algorithm
func computePrefixTable(needle []byte) []int {
	prefixTable := make([]int, len(needle))
	k := 0

	for i := 1; i < len(needle); i++ {
		for k > 0 && needle[k] != needle[i] {
			k = prefixTable[k-1]
		}
		if needle[k] == needle[i] {
			k++
		}
		prefixTable[i] = k
	}

	return prefixTable
}

func printProgress(caption string) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, c := range caption {
		fmt.Print(string(c))
		time.Sleep(time.Duration(15) * time.Millisecond) // Simulate progress
	}
	fmt.Print("\n")
}

func patcherWrapper(originalBytes, targetBytes, replacementBytes []byte) bool {
	fmt.Print("\x1b[2J") // ANSI escape code for clearing the screen
	fmt.Print("\x1b[H")  // ANSI escape code for moving the cursor to the top-left corner

	printProgress(">> Target: " + hex.EncodeToString(targetBytes))
	startTime := time.Now()
	isPatched := patchBinaryKMP(originalBytes, targetBytes, replacementBytes)
	endTime := time.Now()
	printProgress("\n--- Elapsed time: " + endTime.Sub(startTime).String() + " ---\n")

	return isPatched
}

func patchBinaryKMP(dump, target, patch []byte) bool {
	idx := findBytes(dump, target)
	n, m := len(dump), len(target)
	if idx == -1 {
		return false
	}

	if verbose {
		fmt.Printf("\r\n>> % x [%d/%d] Found!", dump[idx:idx+m], idx, n)
		fmt.Printf("\r\n>> % x [%d/%d] Assembling", patch, idx, n)
		time.Sleep(time.Duration(2) * time.Second) // Simulate progress
		copy(dump[idx:], patch)
		//fmt.Printf("\r\n>> % x [%d/%d] Should be assembled", dump[idx:idx+m], idx, n)
		fmt.Printf("\r\n>> now it should be patched")
		fmt.Println("\r\n>> [==================================================>] 100%")
		return true
	}

	copy(dump[idx:], patch)
	return true
}

func computeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func fileExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(srcFile)

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(destFile)

	_, err = io.Copy(destFile, srcFile)
	return err
}

func isBinaryPatched(filePath string, replacementBytes []byte) bool {
	currentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	for i := 0; i < len(currentBytes)-len(replacementBytes); i++ {
		match := true
		for j, b := range replacementBytes {
			if currentBytes[i+j] != b {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
func findBytes(originalBytes, targetBytes []byte) int {
	prefixTable := computePrefixTable(targetBytes)
	n, m := len(originalBytes), len(targetBytes)
	j := 0

	startTime := time.Now()
	for i := 0; i < n; i++ {
		elapsedTime := time.Since(startTime)

		// These progressbar will block the stack by abt 30s from supposedly 2ms but it looks cooler than nothing :))
		if verbose {
			if i+m >= n {
				m = n - i - 1
			}
			// to mitigate the blocking stack, only print progressbar for even i
			if i%20 == 0 {
				fmt.Printf("\r>> % x [%d/%d] %s", originalBytes[i+1:i+1+m], i, n, elapsedTime)
			}
		}

		for j > 0 && originalBytes[i] != targetBytes[j] {
			j = prefixTable[j-1]
		}
		if originalBytes[i] == targetBytes[j] {
			j++
			if j == m {
				return i - m + 1
			}
		}
	}
	return -1
}

func GetFileVersion(filePath string) (fileVersion string, err error) {
	ps := exec.Command("powershell.exe",
		"-Command", fmt.Sprintf("(Get-Command \"%s\").FileVersionInfo.FileVersion", filePath))

	var stdout, stderr bytes.Buffer
	ps.Stdout = &stdout
	ps.Stderr = &stderr

	err = ps.Run()
	if err != nil {
		return "", fmt.Errorf(">> Error running powershell command: %v", err)
	}

	trimmedVersion := strings.TrimSpace(stdout.String())
	if trimmedVersion == "" {
		return "", fmt.Errorf(">> Empty version string")
	}
	//
	//if trimmedVersion == targetVersion {
	//	printProgress(">> " + trimmedVersion + " Binary version match!. Proceeding...\n")
	//	return trimmedVersion, err
	//}

	//printProgress(">> Current version: " + trimmedVersion)
	//printProgress(">> Compatible version: " + targetVersion)
	//printProgress(">> Incompatible version. Please use the compatible version.")
	return trimmedVersion, err
}

// Original generic nested, bloated, naive less concerns approach
//func patchBinary(originalBytes, targetBytes, replacementBytes []byte) bool {
//	found := false
//	for i := 0; i < len(originalBytes)-len(targetBytes); i++ {
//		match := true
//		for j, b := range targetBytes {
//			fmt.Printf("\rSearching for target bytes [%d/%d]", i+j, len(originalBytes))
//			if originalBytes[i+j] != b {
//				match = false
//				break
//			}
//		}
//		if match {
//			copy(originalBytes[i:], replacementBytes)
//			found = true
//			break
//		}
//	}
//	return found
//}

//func digCave(startAddr int, binary []byte, shellcodeSize int) (int, bool) {
//	minCaveSize := shellcodeSize
//
//	// Iterate through the binary to find a suitable cave
//	for i := startAddr; i < len(binary)-minCaveSize; i++ {
//
//		//fmt.Printf("\rSearching for code cave [%d/%d]", i, len(binary))
//		// Check if this region can serve as a cave
//		// For simplicity, we'll just check if there are enough consecutive zero bytes (nop instructions) in the region.
//		isCave := true
//		for j := i; j < i+minCaveSize; j++ {
//			if binary[j] != 0x00 {
//				isCave = false
//				break
//			}
//		}
//
//		if isCave {
//			fmt.Printf("\rfound for cave [%d/%d]", i, len(binary))
//			return i, true
//		}
//	}
//	return 0, false
//}
