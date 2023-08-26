package main

import (
	"fmt"
	"strconv"
	"strings"
)

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

// formatCodesForDepositFileAOUT puts machine codes to output format (simh deposit format)
func formatCodesForDepositFile(confInfo *ConversionInfo) {
	i := 0
	for i < len(confInfo.emitInfo) {
		adr := confInfo.emitInfo[i].address
		value := confInfo.emitInfo[i].value
		if confInfo.debug > 1 {
			fmt.Printf("D %s %s\n\r", toOctalString(adr), toOctalString(value))
		}
		confInfo.emitAOUT(toOctalString(adr), toOctalString(value), true)
		i++
	}
}

// emitAOUT() writes address and data words to output file. Flag endLine defines what command end is given (LF or CR)
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
	conversionInfo.writeCharAOUT('D')
	conversionInfo.writeCharAOUT(' ')

	// Write address
	for i < len(adr) {
		conversionInfo.writeCharAOUT(adr[i])
		i++
	}
	// finish address with ' '
	conversionInfo.writeCharAOUT(' ')

	// Wait some time for output reaction of PDP
	//time.Sleep(time.Duration(20) * time.Millisecond)

	// Write value, if not ""
	i = 0
	for i < len(value) {
		conversionInfo.writeCharAOUT(value[i])
		i++
	}
	// finish line
	if endLine {
		conversionInfo.writeCharAOUT('\r')
		conversionInfo.writeCharAOUT('\n')
	}
}

// writeCharAOUT writes a single char
func (conversionInfo *ConversionInfo) writeChar(c byte) {
	if !conversionInfo.dryMode {
		conversionInfo.outContent.WriteByte(c)
	}
}
