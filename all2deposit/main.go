package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("Convert PDP11 PTAP file (papertape) format to SIMH deposit file")
	inFilePtr := flag.String("in", "a.ptap",
		"Input PTAP file")
	outFilePtr := flag.String("out", "deposit.out",
		"Output DEPOSIT file")
	dryRunPtr := flag.Bool("dry-run", false,
		"dry run mode")
	debugPtr := flag.Int("debug", 0,
		"Debug level")
	flag.Parse()

	fmt.Printf("--in, in File (PTAP or AOUT): %s\n", *inFilePtr)
	fmt.Printf("--out, Output DEPOSIT File: %s\n", *outFilePtr)
	fmt.Printf("--dry-run: %t\n", *dryRunPtr)
	fmt.Printf("--debug: %d\n", *debugPtr)

	var sb strings.Builder

	// Create type
	confInfo := ConversionInfo{
		inFile:        *inFilePtr,
		outFile:       *outFilePtr,
		dryMode:       *dryRunPtr,
		debug:         *debugPtr,
		outContent:    sb,
		outFormat:     Octal,
		state:         Start,
		directive:     None,
		lineNumberPos: 8, // line number last digit position, counted from 1 (not 0)
	}

	// Write prefix lines
	filePrefix := []byte("")
	confInfo.outContent.Write(filePrefix)

	// Convert PTAP to Deposit
	convertFile(&confInfo)

	//
	formatCodesForDepositFile(&confInfo)

	// Write file
	fmt.Printf("%s\n\r", confInfo.outContent.String())
	err := os.WriteFile(confInfo.outFile, []byte(confInfo.outContent.String()), 0644)
	if err != nil {
		log.Printf("Error WritingD DEPOSIT file %s\n\r", err)
		return
	}
}
