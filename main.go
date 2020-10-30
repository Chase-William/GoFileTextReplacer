package main

import (
	"fmt"
	"io"
	"os"
)

/*

	This application is a designed to read a text file and replace a given target word with another.

	Files are read in 1024 byte chunks so large files won't cause issues.

	If a target word is starting to be identified and then the buffer is loaded with new data, the program will pickup where it left off with
		identifing the target word.

	Caveats:

	The text to be replaced must be a single word and cannot contain spaces.

	The text to be replaced with must be a single word and connot contain spaces.

*/

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	const BUFFER_SIZE = 1024 // length of byte buffer
	var fileName string      // fileName given by the user
	var toBeReplaced string  // string given by the user to be replaced
	var newText string       // string to replace the toBeReplaced string
	matches := []int{}

	// Get filename
	fmt.Printf("Please enter the file name of the file to edit: ")
	fmt.Scanln(&fileName)

	// Open is a wrapper OpenFile
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0660)
	Check(err)

	// This function is called when the program is returning from this method
	// It will close the file and propogate an error to console window if one is to handle
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	buffer := make([]byte, 1024)
	for {
		numRead, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			panic(err)
		}

		if numRead == 0 {
			break
		}
	}

	fmt.Printf("Enter the text to be replaced: ")
	fmt.Scanln(&toBeReplaced)

	fmt.Printf("Enter the new text: ")
	fmt.Scanln(&newText)

	charBuffer := make([]byte, BUFFER_SIZE) // Holds the most recent reading
	offset, err := file.Seek(0, 0)
	if err != nil || offset != 0 {
		panic(err)
	}

	// Creating the output file
	oFile, err := os.Create("TestOutput.txt")
	Check(err)
	defer oFile.Close() // Close when leaving the scope of this function

	var replaceIndex int // Used to store the index when iterating for a target word
	for {
		numRead, err := file.Read(charBuffer)

		if err != nil && err != io.EOF {
			panic(err)
		}

		// If nothing was read exit the loop, eof basically
		if numRead == 0 {
			break
		}

		fmt.Printf("Read chunk from file, read characters: %d\n", numRead)

		// Iterating through the buffer just read and adding matches to our collection of matches
		for i := 0; i < numRead; i++ {
			// If a character matches a byte we need to prepare to test the next character in the buffer against the next character in the toBeReplaced string
			if charBuffer[i] == byte(toBeReplaced[replaceIndex]) {
				if i-1 != 0 { // Make sure we arn't indexing into -1
					if string(charBuffer[i-1]) != " " && replaceIndex == 0 { // When evaluating the first character of the word, the previous character must be a white space
						continue
					}
				}
				// If the replace index matches the toBeReplaced string length we have a complete match
				if (replaceIndex + 1) == len(toBeReplaced) {
					// Adding the replaceIndex which contains the index of the first character
					// The current replaceIndex contains the last index value of the matching
					matches = append(matches, i-len(toBeReplaced))
				} else {
					replaceIndex++
				}
			} else {
				replaceIndex = 0
			}
		}

		// Iterate through all matches
		/*
			Reading by filling in from previous match
			&
			Iterate through all matches
		*/
		for i, v := range matches {
			if i == 0 {
				// Write the target area using a slice to file
				oFile.WriteString(string(charBuffer[:v+1]) + newText)
			} else {
				if matches[i-1] > v {
					continue
				}
				// Creating a string from a []byte slice that contains our the text before and then our keyword
				oFile.WriteString(string(charBuffer[matches[i-1]+len(toBeReplaced)+1:v+1]) + newText)
			}

			// If this is the final iteration, check to make sure no more content needs to be written
			if i == len(matches)-1 {
				remaining := charBuffer[v+len(toBeReplaced)+1 : numRead] // Write from the end index of the last word to be replaced to the index of the last read index
				// If the slice is not empty, write it
				if remaining != nil {
					oFile.WriteString(string(remaining))
				}
				matches = matches[i+1:] // Slicing the array to be scoped on the new part
			}
		}
	}

	fmt.Printf("Finished Successfully!")
}

func PrintMatches(buf *[]byte, matches []int) {
	for _, v := range matches {
		fmt.Printf("%s  ", string((*buf)[v:v+5]))
	}
	fmt.Printf("\n")
}
