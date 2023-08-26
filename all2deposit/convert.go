package main

import "C"
import (
	"io/ioutil"
	"log"
)

// convertFile converts a file to DEPOSIT format
func convertFile(conversionInfo *ConversionInfo) {

	// Read in file
	bytes, err := ioutil.ReadFile(conversionInfo.inFile)
	if err != nil {
		log.Printf("Error loading a.out file %s: %s\n\r", conversionInfo.inFile, err)
		return
	}
	magic := int(bytes[0])*256 + int(bytes[1])
	//magic = bytesToInt(bytes, 1)
	log.Printf("Magic: %x", magic)

	if magic == 0x701 {
		log.Printf("Converting from AOUT format")
		convertAoutFile(conversionInfo, magic, bytes)
	} else {
		log.Printf("Converting from PTAP format")
		convertPtapFile(conversionInfo, magic, bytes)
	}

}
