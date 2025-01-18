package main

import (
	"os"
)

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

func main() {
	// Example usage of SaveData1
	err := SaveData1("example.txt", []byte("Hello, World!"))
	if err != nil {
		panic(err)
	}
}
