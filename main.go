package main

import (
	"fmt"
	"math/rand"
	"os"
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
	tmp := fmt.Sprintf("%s.%d", path, rand.Int())
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}

func SaveData3(path string, data []byte) error {
	tmp := fmt.Sprintf("%s.tmp.%d", path, rand.Int())
	fp, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		/// Remove the tmp file if the operation failed
		os.Remove(tmp)
		return err
	}
	err = fp.Sync()
	if err != nil {
		/// Remove the tmp file if the operation failed
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}

func main() {
	// Example usage of SaveData2
	path := "example.txt"
	data := []byte("This is some example data.")

	err := SaveData3(path, data)
	if err != nil {
		fmt.Printf("Error saving data: %v\n", err)
		return
	}

	fmt.Println("Data saved atomically!")
}
