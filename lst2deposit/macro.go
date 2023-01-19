package main

import "C"
import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"unicode"
)

const LF byte = 0xa

// convertFile converts a LST file to DEPOSIT format
func convertFile(conversionInfo *ConversionInfo) {

	// Read in LST file
	lstBytes, err := ioutil.ReadFile(conversionInfo.lstFile)
	if err != nil {
		log.Printf("Error loading LST file %s: %s\n\r", conversionInfo.lstFile, err)
		return
	}

	conversionInfo.state = Start
	i := 0
	lineNumber := 1
	lineNumberRead := 0
	inLinePos := 0

	var address int
	var opcode int
	var operand1 int
	var operand2 int

	for i < len(lstBytes) {
		//fmt.Printf("%c", c)

		if conversionInfo.state == Start {
			inLinePos = 0
			address = -1
			opcode = -1
			operand1 = -1
			operand2 = -1
		}

		if conversionInfo.state == Start {
			// Overread leading spaces
			skipSpaces(lstBytes, &i, &inLinePos)
			conversionInfo.state = BeforeLineNumber
		}
		if conversionInfo.debug > 3 {
			fmt.Printf("%c", lstBytes[i])
		}

		if conversionInfo.state == BeforeLineNumber {
			// Most lines start with a line number
			// If .WORD has more than 3 args, the 4th and following are put to next line(s),
			// without a leading line number.
			//
			// Example:
			//9 000004 000006  000000  000012                  .WORD    6,0,12,0        ; INITIALIZE ERROR VECTORS
			//  000012 000000
			oldlineNumberRead := lineNumberRead
			readNumber(lstBytes, &i, &inLinePos, &lineNumberRead, lineNumber, conversionInfo.debug)
			if inLinePos == conversionInfo.lineNumberPos {
				// We have read a number and its last digit is at position of 'line number' field
				// This means we have actually read a line number
				conversionInfo.state = LineNumber

				if conversionInfo.directive != None {
					// We were up to now in a directive, and leave it
					finalizeDirective(conversionInfo, lineNumberRead, lineNumber, inLinePos)
				}
			} else {
				// We've read a number, but not at the position where we expect a line number.
				// So this number is NOT a line number, but an address
				if conversionInfo.debug > 1 {
					fmt.Printf(" <%d/%d/%d>: No line number given, assuming multiline directive\n\r", lineNumberRead, lineNumber, inLinePos)
				}
				address = lineNumberRead
				convertOctToDec(&address)
				// restore lineNumberRead
				lineNumberRead = oldlineNumberRead
				// Overread spaces
				skipSpaces(lstBytes, &i, &inLinePos)
				conversionInfo.state = AfterAddress
			}
			continue
		}

		if conversionInfo.state == LineNumber {
			convertOctToDec(&lineNumberRead)
			// Overread spaces after line number
			skipSpaces(lstBytes, &i, &inLinePos)
			conversionInfo.state = AfterLineNumber

			// Check for directive
			conversionInfo.directive = detectDirective(lstBytes, i, inLinePos, lineNumberRead, lineNumber,
				conversionInfo.debug)
		}

		if conversionInfo.state == AfterLineNumber {

			if lstBytes[i] == ';' {
				if conversionInfo.debug > 2 {
					fmt.Printf("<%d/%d: Comment line>", lineNumberRead, lineNumber)
				}
				// Overread comment until end of line
				proceedToChar(lstBytes, LF, &i, &inLinePos, &lineNumber)
				conversionInfo.state = Start
				continue
			}

			if lstBytes[i] == LF {
				if conversionInfo.debug > 2 {
					fmt.Printf("<%d/%d: Empty line>\n\r", lineNumberRead, lineNumber)
				}
				lineNumber++
				i++
				conversionInfo.state = Start
				continue
			}

			if lstBytes[i] == '.' {
				if conversionInfo.debug > 2 {
					fmt.Printf("<%d/%d: Directive line>", lineNumberRead, lineNumber)
				}
				// Overread comment until end of line
				proceedToChar(lstBytes, LF, &i, &inLinePos, &lineNumber)
				conversionInfo.state = Start
				continue
			}

			if lstBytes[i] == 'L' && lstBytes[i+1] == 'C' {
				if conversionInfo.debug > 2 {
					fmt.Printf("<%d/%d: LC Directive AfterLineNumber>",
						lineNumberRead, lineNumber)
				}
				// Overread comment until end of line
				proceedToChar(lstBytes, LF, &i, &inLinePos, &lineNumber)
				conversionInfo.state = Start
				continue
			}

			// If the above were not valid, we assume an address to follow
			if unicode.IsDigit(rune(lstBytes[i])) {
				readNumber(lstBytes, &i, &inLinePos, &address, lineNumber, conversionInfo.debug)
				convertOctToDec(&address)
				if conversionInfo.debug > 2 {
					fmt.Printf("%d/%d: Address read in: %06oo/%d.\n\r", lineNumberRead, lineNumber,
						address, address)
				}
				// Overread spaces
				skipSpaces(lstBytes, &i, &inLinePos)
				conversionInfo.state = AfterAddress
				continue
			}
		}

		if conversionInfo.state == AfterAddress {
			// If we get a number after address, this is an opcode or data word or byte
			if unicode.IsDigit(rune(lstBytes[i])) {
				readNumber(lstBytes, &i, &inLinePos, &opcode, lineNumber, conversionInfo.debug)
				convertOctToDec(&opcode)
				if conversionInfo.debug > 2 {
					fmt.Printf("%d/%d: Opcode read in: %06o/%d.\n\r", lineNumberRead, lineNumber,
						opcode, opcode)
				}
				// Overread spaces
				skipSpaces(lstBytes, &i, &inLinePos)
				conversionInfo.state = AfterOpcode
				continue
			}

			if lstBytes[i] == 'L' && lstBytes[i+1] == 'C' {
				if conversionInfo.debug > 2 {
					fmt.Printf("<%d/%d: LC Directive AfterAddress>", lineNumberRead, lineNumber)
				}
				proceedToChar(lstBytes, LF, &i, &inLinePos, &lineNumber)
				conversionInfo.state = Start
				continue
			}

			if lstBytes[i] == '.' {
				if conversionInfo.debug > 2 {
					fmt.Printf("<%d/%d: Directive AfterAddress>", lineNumberRead, lineNumber)
				}
				proceedToChar(lstBytes, LF, &i, &inLinePos, &lineNumber)
				conversionInfo.state = Start
				continue
			}

			//fmt.Printf("Line %d (phys: %d), unhandled char (1): '%c',0x%x.\n\r", lineNumberRead, lineNumber, lstBytes[i], lstBytes[i])
			conversionInfo.state = AfterOpcode
			continue
		}

		if conversionInfo.state == AfterOpcode {
			// After second number, this can follow
			// - spaces and LF
			// - A number = operand
			// - Some text (label or opcode-mnemnonic)

			// Overread spaces
			skipSpaces(lstBytes, &i, &inLinePos)

			if lstBytes[i] == LF {
				i++
				inLinePos = 0
				conversionInfo.state = EmitData
				continue
			}

			if unicode.IsDigit(rune(lstBytes[i])) {
				if conversionInfo.debug > 2 {
					fmt.Printf("%d/%d: Number follows\n\r", lineNumberRead, lineNumber)
				}
				conversionInfo.state = Operand1
				continue
			}

			if conversionInfo.debug > 2 {
				fmt.Printf("%d/%d: unhandled char (2): '%c',0x%x.", lineNumberRead,
					lineNumber, lstBytes[i], lstBytes[i])
			}
			// Overread until end of line
			proceedToChar(lstBytes, LF, &i, &inLinePos, &lineNumber)
			// TODO: EmitData or Start ?
			conversionInfo.state = EmitData
			continue
		}

		if conversionInfo.state == Operand1 {
			// first operand
			readNumber(lstBytes, &i, &inLinePos, &operand1, lineNumber, conversionInfo.debug)
			convertOctToDec(&operand1)
			if conversionInfo.debug > 2 {
				fmt.Printf("%d/%d: Operand1 read in: %06oo/%d.\n\r", lineNumberRead, lineNumber,
					operand1, operand1)
			}
			// Overread spaces
			skipSpaces(lstBytes, &i, &inLinePos)

			if unicode.IsDigit(rune(lstBytes[i])) {
				if conversionInfo.debug > 2 {
					fmt.Printf("%d/%d: Number follows\n\r", lineNumberRead, lineNumber)
				}
				conversionInfo.state = Operand2
				continue
			}

			// Overread until end of line
			proceedToChar(lstBytes, LF, &i, &inLinePos, &lineNumber)
			conversionInfo.state = EmitData
			continue
		}

		if conversionInfo.state == Operand2 {
			// We have a second operand
			readNumber(lstBytes, &i, &inLinePos, &operand2, lineNumber, conversionInfo.debug)
			convertOctToDec(&operand2)
			if conversionInfo.debug > 2 {
				fmt.Printf("%d/%d: Operand2 read in: %06o/%d.\n\r", lineNumberRead, lineNumber,
					operand2, operand2)
			}
			// Overread spaces
			skipSpaces(lstBytes, &i, &inLinePos)

			if unicode.IsDigit(rune(lstBytes[i])) {
				if conversionInfo.debug > 2 {
					fmt.Printf("%d/%d: number follows, not implemented\n\r", lineNumberRead,
						lineNumber)
				}
				conversionInfo.state = EmitData
				continue
			}

			if lstBytes[i] == '.' {
				if conversionInfo.debug > 2 {
					fmt.Printf("<%d/%d: Directive after Operand2>", lineNumberRead, lineNumber)
				}
				proceedToChar(lstBytes, LF, &i, &inLinePos, &lineNumber)
				conversionInfo.state = EmitData
				continue
			}

			if conversionInfo.debug > 2 {
				fmt.Printf("%d/%d: unhandled char (3): '%c',0x%x.\n\r", lineNumberRead,
					lineNumber, lstBytes[i], lstBytes[i])
			}
			// Overread until end of line
			proceedToChar(lstBytes, LF, &i, &inLinePos, &lineNumber)
			conversionInfo.state = EmitData
			continue
		}

		if conversionInfo.state == EmitData {
			emitData(conversionInfo, address, opcode, operand1, operand2)
			conversionInfo.state = Start
		}
	}
}

// finalizeDirective does all things required if a directive ends
func finalizeDirective(conversionInfo *ConversionInfo, lineNumberRead int, lineNumber int, inLinePos int) {
	if conversionInfo.debug > 0 {
		fmt.Printf(" <%d/%d/%d>: Leaving directive %d\n\r", lineNumberRead, lineNumber, inLinePos,
			conversionInfo.directive)
	}
	// Emit all collected values
	conversionInfo.emitDirective(conversionInfo.directive)

	// Set directive state back
	conversionInfo.directive = None
	conversionInfo.directiveBuf = nil
}

// detectDirective checks for a directive in current line
func detectDirective(lstBytes []byte, i int, inLinePos int, lineNumberRead int, lineNumber int, debug int) Directive {
	// get size off current line
	j := i
	for lstBytes[j] != LF {
		j++
	}
	l := make([]byte, j-i)
	// copy line to temp buffer
	j = i
	for lstBytes[j] != LF {
		l[j-i] = lstBytes[j]
		j++
	}
	line := strings.ToLower(string(l))
	if strings.Contains(line, ".word") {
		if debug > 0 {
			fmt.Printf("<%d/%d> Directive .WORD\n\r", lineNumber, lineNumberRead)
		}
		return Word
	}
	if strings.Contains(line, ".byte") {
		if debug > 0 {
			fmt.Printf("<%d/%d> Directive .BYTE\n\r", lineNumber, lineNumberRead)
		}
		return Byte
	}
	if strings.Contains(line, ".ascii") {
		if debug > 0 {
			fmt.Printf("<%d/%d> Directive .ASCII\n\r", lineNumber, lineNumberRead)
		}
		return Ascii
	}
	if strings.Contains(line, ".blkw") {
		if debug > 0 {
			fmt.Printf("<%d/%d> Directive .BLKW\n\r", lineNumber, lineNumberRead)
		}
		return Blkw
	}
	if strings.Contains(line, ".blkb") {
		if debug > 0 {
			fmt.Printf("<%d/%d> Directive .BLKB\n\r", lineNumber, lineNumberRead)
		}
		return Blkb
	}
	return None
}

// emitData dumps collected data to output file
func emitData(conversionInfo *ConversionInfo, address int, opcode int, operand1 int, operand2 int) {
	if address == -1 || opcode == -1 {
		fmt.Printf("NOT EMITTING D %06o %06o %06o %06o\n\r", address, opcode, operand1, operand2)
	} else {
		adr := address
		emitOrKeep(conversionInfo, adr, opcode)
		if operand1 != -1 {
			adr += 2
			emitOrKeep(conversionInfo, adr, operand1)
		}
		if operand2 != -1 {
			adr += 2
			emitOrKeep(conversionInfo, adr, operand2)
		}
	}
}

// emitOrKeep routes data to handle to processing functions
func emitOrKeep(conversionInfo *ConversionInfo, adr int, data int) {
	if conversionInfo.directive == None {
		// Write directly to output buffer
		conversionInfo.appendToEmitBuf(adr, data)
	} else {
		// Write directive data to directive buffer for later processing
		conversionInfo.appendToDirectiveBuf(adr, data)
	}
}

// toOctalString converts number to octal string representation
func toOctalString(n int) string {
	return fmt.Sprintf("%o", n)
}

// convertOctToDec creates a octal number (stored as decimal value) to its decimal equivalent
func convertOctToDec(n *int) {
	pot := 1
	m := *n
	result := 0
	for m > 0 {
		lastDigit := m % 10
		result += lastDigit * pot

		pot *= 8
		m /= 10
	}
	*n = result
}

// proceedToChar proceeds in byte stream to a character code
func proceedToChar(lstBytes []byte, searchChar byte, i *int, inLinePos *int, lineNumber *int) {
	for lstBytes[*i] != searchChar {
		fmt.Printf("%c", lstBytes[*i])
		*i++
		*inLinePos++
	}
	fmt.Print("\n\r")
	if searchChar == LF {
		*lineNumber++
		*inLinePos = 0
	}
	*i++
}

// skipSpaces skips spaces in a byte stream
func skipSpaces(lstBytes []byte, i *int, inLinePos *int) {
	for lstBytes[*i] == ' ' {
		*i++
		*inLinePos++
	}
}

// readNumber reads in an octal number
func readNumber(lstBytes []byte, i *int, inLinePos *int, lineNumberRead *int, lineNumber int, debug int) {
	*lineNumberRead = 0
	for unicode.IsDigit(rune(lstBytes[*i])) {
		*lineNumberRead = *lineNumberRead*10 + int(lstBytes[*i]) - int('0')
		*i++
		*inLinePos++
	}
	if debug > 3 {
		fmt.Printf("%d: number read in: %d.\n\r", lineNumber, *lineNumberRead)
	}
}
