package main

import (
	"fmt"
	"log"
	"os"
)

const BLOCK_HEADER_SIZE = 6

func convertPtapFile(info *ConversionInfo, magic int, bytes []byte) {

	fileLength := len(bytes)
	log.Printf("Total file length: %v (0x%x, octal: %o)", fileLength, fileLength, fileLength)
	offset := 0
	numBlocks := 0
	for offset < len(bytes) {
		log.Printf("\n\nBlock #%v, starts at %v (0x%x, octal: %o)", numBlocks, offset, offset, offset)
		blockLength := convertBlock(bytes, offset)
		offset += blockLength
		offset += 1 // checksum byte
		log.Printf("remaining bytes in file: %v", fileLength-offset)
		numBlocks++
	}
}

func convertBlock(bytes []byte, blockOffset int) int {
	validStartBlock := bytes[0] == 1 && bytes[1] == 0
	if !validStartBlock {
		log.Fatalf("Invalid start block, bytes 0 and 1: 0x%x, 0x%x (should be 0x1,0x0)", bytes[0], bytes[1])
		log.Fatalf("This is maybe not a PTAP file")
		os.Exit(-1)
	} else {
		log.Print("Valid block signature found")
	}

	dataLengthInBlock := bytesToInt(bytes, 2+blockOffset)
	log.Printf("Length of data inside block: dec: %v (0x%x, octal: %o(=brt. %o))", dataLengthInBlock, dataLengthInBlock,
		dataLengthInBlock, dataLengthInBlock-BLOCK_HEADER_SIZE)

	startAddress := bytesToInt(bytes, 4+blockOffset) //  bytes[3]*256+bytes[2]
	log.Printf("data startAddress block: octal: %o (0x%x, dec: %v)", startAddress, startAddress, startAddress)
	// full block = header + data bytes + 1 byte checksum
	lastByteOfBlock := blockOffset + dataLengthInBlock + 1
	log.Printf("Last byte position of block: octal: %o (0x%x, dec: %v)", lastByteOfBlock, lastByteOfBlock, lastByteOfBlock)

	if dataLengthInBlock == BLOCK_HEADER_SIZE {
		// no data at all
		// meaning:  the address is the place to start the program. If odd, usually 1, then stop.
		if startAddress%2 == 1 {
			// odd -> stop tape reader
			log.Printf("Read empty block with odd address value %v. Stopping reading tape.", startAddress)
		} else {
			log.Printf("Read empty block with valid start address  octal: %o (0x%x, dec: %v)", startAddress, startAddress, startAddress)
		}
		return BLOCK_HEADER_SIZE
	}

	// read in data for block
	// text section
	var buf1 = make([]byte, dataLengthInBlock-BLOCK_HEADER_SIZE)
	//writeSegmentHeader(buf1, text_addr)
	offset := BLOCK_HEADER_SIZE + blockOffset
	i := 0
	var checkSum byte = 0
	for i < 6 {
		checkSum += bytes[i+offset]
		i++
	}
	i = 0
	for i < dataLengthInBlock-BLOCK_HEADER_SIZE {
		buf1[i] = bytes[i+offset]
		checkSum += buf1[i]
		i++
	}
	checkSumTarget := bytes[offset+dataLengthInBlock]
	checkSum += checkSumTarget % 0xff
	log.Printf("Checksum: %0x, should be: %0x (diff ignored)\n", checkSum, checkSumTarget)
	i = 0
	for i < dataLengthInBlock-BLOCK_HEADER_SIZE {
		wordValue := bytesToInt(buf1, i)
		fmt.Printf("d %o %06o\n", i+startAddress, wordValue)
		i += 2
	}
	return dataLengthInBlock
}
