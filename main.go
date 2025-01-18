package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// SaveData1 writes data directly to the file at the given path.
func SaveData1(path string, data []byte) error {
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(data)
	if err != nil {
		return err
	}
	return fp.Sync() // fsync
}

// SaveData2 performs an atomic write operation to ensure data integrity.
func SaveData2(path string, data []byte) error {
	// Create a temporary file name in the same directory as the target file
	tmp := fmt.Sprintf("%s.tmp.%d", path, randomInt())
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		return err
	}
	defer func() { // Cleanup: remove the temporary file if an error occurs
		fp.Close() // Close the file
		if err != nil {
			os.Remove(tmp) // Remove the temporary file
		}
	}()

	// Step 1: Write the data to the temporary file
	if _, err = fp.Write(data); err != nil {
		return err
	}

	// Step 2: Flush the data to the disk to ensure durability
	if err = fp.Sync(); err != nil {
		return err
	}

	// Step 3: Rename the temporary file to the target file atomically
	err = os.Rename(tmp, path)
	return err
}

// randomInt generates a random integer for creating unique temporary file names.
func randomInt() int {
	rand.NewSource(time.Now().UnixNano()) // Seed the random number generator
	return rand.Int()                     // Generate a random integer
}

func main() {
	// Example usage of SaveData2
	err := SaveData2("log.txt", []byte("Hello, Atomic World!"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Data saved atomically!")
}
