package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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

func randomInt() int {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err) // Handle error appropriately
	}
	return int(binary.LittleEndian.Uint64(b[:]))
}

func SaveDataImproved(ctx context.Context, path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp := filepath.Join(dir, fmt.Sprintf(".%s.tmp.%d", filepath.Base(path), randomInt()))
	log.Printf("Writing data to temporary file: %s", tmp)

	fp, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("permission denied: %w", err)
		}
		return err
	}
	defer fp.Close()
	defer os.Remove(tmp)

	select {
	case <-ctx.Done():
		return ctx.Err() // Return if context is canceled
	default:
	}

	if _, err := io.WriteString(fp, string(data)); err != nil {
		return err
	}
	if err := fp.Sync(); err != nil {
		return err
	}

	log.Printf("Renaming temporary file to: %s", path)
	return os.Rename(tmp, path)
}

func main() {
	path := "example.txt"
	data := []byte("This is some example data.")

	err := SaveDataImproved(context.Background(), path, data)
	if err != nil {
		log.Printf("Error saving data: %v\n", err)
		return
	}

	log.Println("Data saved atomically!")
}
