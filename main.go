package main

import (
	"os" // Import the "os" package for file operations
)

// SaveData1 is a function to save data to a specified file.
// Parameters:
// - path: the file path where data will be saved
// - data: the byte slice containing the data to write to the file
func SaveData1(path string, data []byte) error {
	// Open the file at the specified path for writing.
	// Flags used:
	// - os.O_WRONLY: Open the file for write-only access
	// - os.O_CREATE: Create the file if it doesn't exist
	// - os.O_TRUNC: Truncate (clear) the file's contents if it already exists
	// Permissions: 0664 (read/write for owner, read for group/others)
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		// Return an error if the file couldn't be opened or created
		return err
	}
	defer fp.Close() // Ensure the file is closed when the function exits

	// Write the provided data to the file
	_, err = fp.Write(data)
	if err != nil {
		// Return an error if the write operation fails
		return err
	}

	// Sync the file to ensure all data is flushed to disk
	return fp.Sync() // fsync ensures durability of the write operation
}

func main() {
	// Example usage of SaveData1 function

	// Define the file path where data will be saved
	filePath := "example.txt"

	// Define the data to be written to the file
	data := []byte("Hello, World!") // Convert a string to a byte slice

	// Call SaveData1 to save the data to the file
	err := SaveData1(filePath, data)
	if err != nil {
		// Handle any errors that occur during the save operation
		panic(err) // Panic is used here for simplicity
	}

	// If no error occurs, the program completes successfully
	// The file "example.txt" should now contain "Hello, World!"
}
