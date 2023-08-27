package main

import "strings"

// ConversionInfo keeps config info for conversion
type ConversionInfo struct {
	inFile           string
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

// ConversionState number format to write
type ConversionState int

const (
	Start ConversionState = 0
)

// Directive types
type Directive int

const (
	None Directive = 0
)

type EmitInfo struct {
	address int
	value   int
}

const LF byte = 0xa
