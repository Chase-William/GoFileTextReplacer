package main

import (
	//"bufio"
	"fmt"
	"io"

	//"io/ioutil"
	"os"
)

/*

	This application is a designed to read a text file and replace a given target word with another.

	Files are read in 1024 byte chunks so large files won't cause issues.

	If a target word is starting to be identified and then the buffer is loaded with new data, the program will pickup where it left off with
		identifing the target word.

*/

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	const BUFFER_SIZE = 1024
	var fileName string = "test.txt"
	var toBeReplaced string = "Ring"
	var newText string = "Necklace"
	//matches := list.New()
	matches := []int{}

	fmt.Printf("Please enter the file name of the file to edit: ")
	//fmt.Scanln(&fileName)

	//ReadAndPrintFile(file);
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

		fmt.Printf("%s", string(buffer))
	}

	// fmt.Printf("\nEnter the text to be replaced: ")
	// fmt.Scanln(&toBeReplaced)

	// fmt.Printf("Enter the new text: ")
	// fmt.Scanln(&newText)

	var replaceIndex int
	var hasMatch bool
	//backBuffer := make([]byte, BUFFER_SIZE)		// Holds the last reading, important for matches that span buffers
	frontBuffer := make([]byte, BUFFER_SIZE) // Holds the most recent reading
	offset, err := file.Seek(0, 0)
	if err != nil || offset != 0 {
		panic(err)
	}
	fmt.Printf("\n=========================================================================\n")

	// Creating the output file
	oFile, err := os.Create("TestOutput.txt")
	Check(err)
	defer oFile.Close() // Close when leaving the scope of this function

	for {
		numRead, err := file.Read(frontBuffer)

		if err != nil && err != io.EOF {
			panic(err)
		}
		// How many characters were read
		fmt.Printf("Num Read: %d\n", numRead)

		if numRead == 0 {
			break
		}

		for i := 0; i < numRead; i++ {
			//fmt.Printf("%x ", frontBuffer[i])
			// If a character matches a byte we need to prepare to test the next character in the buffer against the next character in the toBeReplaced string
			if frontBuffer[i] == byte(toBeReplaced[replaceIndex]) {
				// Printing the matching character and its index value
				fmt.Printf("Character: %c | Index in file: %d\n", frontBuffer[i], i)
				// If the replace index matches the toBeReplaced string length we have a complete match
				if (replaceIndex + 1) == len(toBeReplaced) {
					// Print the characters for now
					fmt.Printf("Complete Match Made\n\n")
					// Updating state variable about a match being found
					hasMatch = true
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

		fmt.Printf("%s\n", frontBuffer)

		// A match was made so we want to replace the correct text as we write it to file
		// We cannot write to the buffer because we will not have room to replace small words with larger ones without overwriting pre-existing text
		if hasMatch {
			// Iterate through all matches

			/*
				Reading by filling in from previous match
			*/
			for i, v := range matches {
				fmt.Printf("Replace Starting Index:%2d\n", v)
				fmt.Printf("Current Index:%2d\n", i)
				if i == 0 {
					// Write the target area using a slice to file
					oFile.WriteString(string(frontBuffer[:v+1]) + newText)
					fmt.Printf("Write: %s%s\n\n", (frontBuffer[:v+1]), newText)
				} else {
					// Creating a string from a []byte slice that contains our the text before and then our keyword
					oFile.WriteString(string(frontBuffer[matches[i-1]+len(toBeReplaced)+1:v+1]) + newText)
					fmt.Printf("Write: %s%s\n\n", newText, (frontBuffer[matches[i-1]+len(toBeReplaced)+1 : v+1]))
				}
			}
		} else {
			// Write the frontBuffer directly to a new file
		}
		// Move the contents of the front
		// backBuffer = frontBuffer

		fmt.Printf("=========================================================================\n\n")
	}
}
