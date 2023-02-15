package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// ConversionInfo keeps config info for conversion
type ConversionInfo struct {
	aoutFile         string
	outFile          string
	dataAlign        int
	text             int
	vector           string
	dryMode          bool
	debug            int
	outContent       strings.Builder
	outFormat        OutFormat
	state            ConversionState
	directiveBuf     []int
	emitStartAddress int
	emitInfo         []EmitInfo
	directive        Directive
	lineNumberPos    int
}

// OutFormat number format to write
type OutFormat int

const (
	Octal OutFormat = 0
)

// OutFormat number format to write
type ConversionState int

const (
	Start            ConversionState = 0
	BeforeLineNumber                 = 1
	LineNumber                       = 2
	AfterLineNumber                  = 3
	AfterAddress                     = 4
	AfterOpcode                      = 5
	Operand1                         = 6
	Operand2                         = 7
	EmitData                         = 8
)

// Directive types
type Directive int

const (
	None  Directive = 0
	Word            = 1
	Byte            = 2
	Ascii           = 3
	Blkw            = 4
	Blkb            = 5
)

type EmitInfo struct {
	address int
	value   int
}

func main() {
	fmt.Println("Convert SREC file to PDP-11 LDA (papertape) format - Converts a PDP11 OBJ to SIMH deposit file")
	aoutFilePtr := flag.String("aout", "a.out",
		"Input a.out file")
	outFilePtr := flag.String("out", "deposit.out",
		"Output DEPOSIT file")
	dataAlignPtr := flag.Int("data-align", 256,
		"Debug level")
	textPtr := flag.String("text", "01000",
		"text")
	vectorPtr := flag.String("vector0", "",
		"vector, store JMP entry at vector 0")
	dryRunPtr := flag.Bool("dry-run", false,
		"dry run mode")
	debugPtr := flag.Int("debug", 0,
		"Debug level")
	flag.Parse()

	fmt.Printf("--aout, AOUT File: %s\n", *aoutFilePtr)
	fmt.Printf("--out, Output DEPOSIT File: %s\n", *outFilePtr)
	fmt.Printf("--data-align, Data align value: %d\n", *dataAlignPtr)
	fmt.Printf("--text, text pointer: %s\n", *textPtr)
	fmt.Printf("--vector, vector pointer: %s\n", *vectorPtr)
	fmt.Printf("--dry-run: %t\n", *dryRunPtr)
	fmt.Printf("--debug: %d\n", *debugPtr)

	var sb strings.Builder

	// Create type
	confInfo := ConversionInfo{
		aoutFile:      *aoutFilePtr,
		outFile:       *outFilePtr,
		dataAlign:     *dataAlignPtr,
		text:          convertNumberString(*textPtr),
		vector:        *vectorPtr,
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

	// Convert lst file to machine code array
	convertFile(&confInfo)

	// Convert machine code array to deposit codes
	formatCodesForDepositFile(&confInfo)

	// Write file
	fmt.Printf("%s\n\r", confInfo.outContent.String())
	err := os.WriteFile(confInfo.outFile, []byte(confInfo.outContent.String()), 0644)
	if err != nil {
		log.Printf("Error WritingD DEPOSIT file %s\n\r", err)
		return
	}
}

func convertNumberString(s string) int {
	var base int
	cleaned := s
	if s[0:2] == "0x" || s[0:2] == "0X" {
		base = 16
		cleaned = strings.Replace(s, "0x", "", -1)
		cleaned = strings.Replace(cleaned, "0X", "", -1)
	} else {
		if s[0] == '0' {
			base = 8
		} else {
			base = 10
		}
	}
	result, _ := strconv.ParseUint(cleaned, base, 32)
	return int(result)
}

// formatCodesForDepositFile puts machine codes to output format (simh deposit format)
func formatCodesForDepositFile(confInfo *ConversionInfo) {
	i := 0
	for i < len(confInfo.emitInfo) {
		adr := confInfo.emitInfo[i].address
		value := confInfo.emitInfo[i].value
		if confInfo.debug > 1 {
			fmt.Printf("D %s %s\n\r", toOctalString(adr), toOctalString(value))
		}
		confInfo.emit(toOctalString(adr), toOctalString(value), true)
		i++
	}
}

// emit() writes address and data words to output file. Flag endLine defines what command end is given (LF or CR)
// Output format is SIMH DEPOSIT
func (conversionInfo *ConversionInfo) emit(adr string, value string, endLine bool) {
	if conversionInfo.dryMode {
		fmt.Printf("dry mode, not writing: %s/%s\n\r", adr, value)
		return
	}

	if conversionInfo.debug > 0 {
		fmt.Printf("D %s %s\n\r", adr, value)
	}

	i := 0

	// Start DEPOSIT command
	conversionInfo.writeChar('D')
	conversionInfo.writeChar(' ')

	// Write address
	for i < len(adr) {
		conversionInfo.writeChar(adr[i])
		i++
	}
	// finish address with ' '
	conversionInfo.writeChar(' ')

	// Wait some time for output reaction of PDP
	//time.Sleep(time.Duration(20) * time.Millisecond)

	// Write value, if not ""
	i = 0
	for i < len(value) {
		conversionInfo.writeChar(value[i])
		i++
	}
	// finish line
	if endLine {
		conversionInfo.writeChar('\r')
		conversionInfo.writeChar('\n')
	}
}

// writeChar writes a single char
func (conversionInfo *ConversionInfo) writeChar(c byte) {
	if !conversionInfo.dryMode {
		conversionInfo.outContent.WriteByte(c)
	}
}
