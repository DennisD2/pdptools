package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func convertAoutFile(conversionInfo *ConversionInfo, magic int, bytes []byte) {
	if magic == 0x0701 && conversionInfo.dataAlign != 2 {
		log.Printf("*** WARNING: impure executable, but alignment not 2")
	}

	text_addr := conversionInfo.text

	text_len := bytesToInt(bytes, 2) //  bytes[3]*256+bytes[2]
	data_addr := (conversionInfo.text + text_len + conversionInfo.dataAlign - 1) & ^(conversionInfo.dataAlign - 1)
	data_len := bytesToInt(bytes, 4) // bytes[5]*256+bytes[4]
	bss_len := bytesToInt(bytes, 6)  // bytes[7]*256+bytes[6]
	//sym_len = bytes[9]*256+bytes[8]
	entry := bytesToInt(bytes, 10) // bytes[11]*256+bytes[10]
	//unused = bytes[13]*256+bytes[12]
	//flags = bytes[15]*256+bytes[14]
	total := data_addr + data_len + bss_len
	log.Printf("      Magic    Text     Len    Data     Len     BSS = MaxMem   Entry")
	log.Printf("Hex:   %04x    %04x    %04x    %04x    %04x    %04x     %04x    %04x",
		magic, text_addr, text_len, data_addr, data_len, bss_len, total, entry)
	log.Printf("Dec:  %05d   %05d   %05d   %05d   %05d   %05d    %05d   %05d",
		magic, text_addr, text_len, data_addr, data_len, bss_len, total, entry)
	log.Printf("Oct: %06o  %06o  %06o  %06o  %06o  %06o   %06o  %06o",
		magic, text_addr, text_len, data_addr, data_len, bss_len, total, entry)

	// text section
	var buf1 = make([]byte, 16+text_len)
	writeSegmentHeader(buf1, text_addr)
	offset := 16
	i := 0
	for i < text_len {
		buf1[i+16] = bytes[i+offset]
		i++
	}
	adr := conversionInfo.text
	writeLdaRecordText(buf1, text_len, adr)
	offset += text_len

	// Data and BSS section
	buf1 = make([]byte, 16+data_len)
	writeSegmentHeader(buf1, data_addr)
	i = 0

	if data_len%2 != 0 {
		// align next BSS section
		fmt.Printf("TBD: Add one byte for alignment\n")
	}

	for i < data_len {
		buf1[i+16] = bytes[i+offset]
		i++
	}
	adr += offset
	writeLdaRecordData(buf1, data_len, adr)

	if conversionInfo.vector != "" {
		buf1 = make([]byte, 16)
		buf1[0] = 1
		buf1[1] = 0
		buf1[2] = 0
		buf1[3] = 0
		buf1[4] = 0
		buf1[5] = 0
		buf1[6] = 0x5f
		buf1[7] = 0
		// TBD convert to int
		//buf1[8] = conversionInfo.vector & 0xff
		//buf1[9] = (conversionInfo.vector >> 8) & 0xff#
		//write_lda_record(output, out_buffer)
	}
}

func writeSegmentHeader(buf []byte, addr int) {
	// 16 byte header
	buf[0] = 1
	buf[1] = 0
	buf[2] = 0
	buf[3] = 0
	buf[4] = byte(addr & 0xff)
	buf[5] = byte(addr / 0256)

	// 10 bytes to complete header
	i := 6
	for i < 16 {
		buf[i] = 0
		i++
	}
}

func writeLdaRecordText(buf []byte, num int, adr int) int {
	length := len(buf)
	//bytes[0] = 1
	//bytes[1] = 0
	buf[2] = byte(length & 0xff)
	buf[3] = byte((length >> 8) & 0xff)
	i := 0
	for i < num {
		e := toOctalString(bytesToInt(buf, i+16))
		fmt.Printf("D %s %s\n", toOctalString(adr), e)
		i += 2
		adr += 2
	}
	return adr
}

func writeLdaRecordData(buf []byte, num int, adr int) int {
	length := len(buf)
	//bytes[0] = 1
	//bytes[1] = 0
	buf[2] = byte(length & 0xff)
	buf[3] = byte((length >> 8) & 0xff)
	i := 0
	for i < num {
		e := toOctalString(bytesToInt(buf, i+16))
		fmt.Printf("D %s %s\n", toOctalString(adr), e)
		i += 2
		adr += 2
	}
	return adr
}

func bytesToInt(bytes []byte, i int) int {
	return int(bytes[i+1])*256 + int(bytes[i])
}

// toOctalString converts number to octal string representation
func toOctalString(n int) string {
	return fmt.Sprintf("%o", n)
}

func convertNumberStringAOUT(s string) int {
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
func formatCodesForDepositFileAOUT(confInfo *ConversionInfo) {
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
func (conversionInfo *ConversionInfo) emitAOUT(adr string, value string, endLine bool) {
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
func (conversionInfo *ConversionInfo) writeCharAOUT(c byte) {
	if !conversionInfo.dryMode {
		conversionInfo.outContent.WriteByte(c)
	}
}
