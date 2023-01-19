package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// ConversionInfo keeps config info for conversion
type ConversionInfo struct {
	lstFile          string
	outFile          string
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
	fmt.Println("Converts a PDP11 OBJ to SIMH deposit file")
	lstFilePtr := flag.String("lst", "",
		"Input LST file")
	outFilePtr := flag.String("out", "deposit.out",
		"Output DEPOSIT file")
	dryRunPtr := flag.Bool("dry-run", false,
		"Debug level")
	debugPtr := flag.Int("debug", 0,
		"Debug level")
	flag.Parse()

	fmt.Printf("--lst, LST File: %s\n", *lstFilePtr)
	fmt.Printf("--out, Output DEPOSIT File: %s\n", *outFilePtr)
	fmt.Printf("--dry-run: %t\n", *dryRunPtr)
	fmt.Printf("--debug: %d\n", *debugPtr)

	var sb strings.Builder

	// Create type
	confInfo := ConversionInfo{
		lstFile:       *lstFilePtr,
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
	filePrefix := []byte("set cpu 11/23\n\r")
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

// appendToDirectiveBuf append directive data to directive buffer
func (conversionInfo *ConversionInfo) appendToDirectiveBuf(adr int, value int) {
	var size int
	var oldEmitBuf []int
	if conversionInfo.directiveBuf == nil {
		// Use first address in directive as start address
		conversionInfo.emitStartAddress = adr
		// directiveBuf needs initial size
		size = 0
	} else {
		size = len(conversionInfo.directiveBuf)
		oldEmitBuf = conversionInfo.directiveBuf
	}

	// Copy over
	size++
	var i = 0
	conversionInfo.directiveBuf = make([]int, size)
	for i < len(oldEmitBuf) {
		conversionInfo.directiveBuf[i] = oldEmitBuf[i]
		i++
	}
	// Append new value
	conversionInfo.directiveBuf[i] = value

	if conversionInfo.debug > 1 {
		i = 0
		fmt.Printf("directiveBuf[%d] = ", len(conversionInfo.directiveBuf))
		for i < len(conversionInfo.directiveBuf) {
			fmt.Printf("%o:%o, ", conversionInfo.emitStartAddress+i, conversionInfo.directiveBuf[i])
			i++
		}
		fmt.Printf("\n\r")
	}
}

// emitDirective emits the collected directive data
func (conversionInfo *ConversionInfo) emitDirective(directive Directive) {
	i := 0
	adr := conversionInfo.emitStartAddress

	offset, err := conversionInfo.getOffset(directive)
	if err {
		fmt.Printf("ERROR: Unknown directive value %d, not emitting data", directive)
		return
	}

	if offset == 2 {
		// data consists of words
		for i < len(conversionInfo.directiveBuf) {
			data := conversionInfo.directiveBuf[i]
			conversionInfo.appendToEmitBuf(adr, data)
			adr += offset
			i++
		}
	}
	if offset == 1 {
		// data consists of bytes; put two bytes from directiveBuf into a word
		if adr%2 != 0 {
			// Special handling for byte directives that start at an odd address
			// We need to put the first byte into the previous word
			fmt.Printf("Directive %d, start address at byte boundary %oo\n\r", directive, adr)
			prevEmitData, err := conversionInfo.getEmitData(adr - 1)
			if err {
				fmt.Printf("ERROR: No emit data found for address %oo\n\r", adr)
			} else {
				fmt.Printf("Emit data found for address %oo, value is %oo\n\r", prevEmitData.address, prevEmitData.value)
				// put first byte of new directive data into previous word
				if conversionInfo.directiveBuf != nil {
					prevEmitData.value += conversionInfo.directiveBuf[0] * 256
					// move i and adr one byte forward
					i++
					adr++
				}
			}
		}
		for i < len(conversionInfo.directiveBuf) {
			data := conversionInfo.directiveBuf[i]
			if i+1 < len(conversionInfo.directiveBuf) {
				data += conversionInfo.directiveBuf[i+1] * 256
			}
			conversionInfo.appendToEmitBuf(adr, data)
			adr += 2
			i += 2
		}
	}
}

// getOffset returns required address offset
func (conversionInfo *ConversionInfo) getOffset(directive Directive) (int, bool) {
	var offset int
	switch directive {
	case Word:
		offset = 2
		break
	case Blkw:
		offset = 2
		break
	case Byte:
		offset = 1
		break
	case Blkb:
		offset = 1
		break
	case Ascii:
		offset = 1
		break
	default:
		return 0, true
	}
	return offset, false
}

// appendToEmitBuf appends (non-directive) data to emit buffer
func (conversionInfo *ConversionInfo) appendToEmitBuf(adr int, value int) {
	var size int
	var oldEmitInfo []EmitInfo
	if conversionInfo.emitInfo == nil {
		// emitInfo needs initial size
		size = 0
	} else {
		size = len(conversionInfo.emitInfo)
		oldEmitInfo = conversionInfo.emitInfo
	}

	// Copy over
	size++
	var i = 0
	conversionInfo.emitInfo = make([]EmitInfo, size)
	for i < len(oldEmitInfo) {
		conversionInfo.emitInfo[i] = oldEmitInfo[i]
		i++
	}
	// Append new value
	newEmitInfo := EmitInfo{adr, value}
	conversionInfo.emitInfo[i] = newEmitInfo
}

// getEmitData returns emit data struct for address adr
func (conversionInfo *ConversionInfo) getEmitData(adr int) (*EmitInfo, bool) {
	i := 0
	max := len(conversionInfo.emitInfo)
	for i < max && conversionInfo.emitInfo[i].address != adr {
		i++
	}
	if i == max {
		return nil, true
	}
	return &conversionInfo.emitInfo[i], false
}
