package main

import (
	//"bufio"
	"fmt"
	"io"
	//"io/ioutil"
	"os"
)

func Check(e error) {
	if (e != nil) {
		panic(e)
	}
}

func main() {
	var fileName string
	var toBeReplaced string
	var newText string


	fmt.Printf("Please enter the file name of the file to edit: ")
	fmt.Scanln(&fileName)

	//ReadAndPrintFile(file);
	// Open is a wrapper OpenFile
	file, err := os.OpenFile(fileName, os.O_RDWR | os.O_APPEND, 0660)
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

		if numRead == 0 { break }

		fmt.Printf("%s", string(buffer))
	}

	fmt.Printf("\nEnter the text to be replaced: ")
	fmt.Scanln(&toBeReplaced)

	fmt.Printf("Enter the new text: ")
	fmt.Scanln(&newText)

	var replaceIndex int

	buffer2 := make([]byte, 1024)
	for {
		numRead, err := file.Read(buffer2)

		if err != nil && err != io.EOF {
			panic(err)
		}						

		fmt.Printf("%d", numRead)
		if numRead == 0 { break }

		for i := 0; i < numRead; i++ {
			fmt.Printf("for loop iterated\n")
			if buffer2[i] == byte(toBeReplaced[replaceIndex]) {
				fmt.Printf("Found a match\n")
				if replaceIndex == len(toBeReplaced) {
					fmt.Printf("Found word to be replaced!\n")
				} else {
					replaceIndex++
				}				
			} else {
				replaceIndex = 0
			}
		}
	}



	// Create a buffer of 1024 bytes
	// buffer := make([]byte, 1024)
	// for {
	// 	// read a chunk
	// 	numRead, err := file.Read(buffer)
	// 	if err != nil && err != io.EOF {
	// 		panic(err)
	// 	}
	// 	// If nothing was read from the file then break from loop
	// 	if numRead == 0 {
	// 		break
	// 	}

	// 	if _, err := 
	// }


}