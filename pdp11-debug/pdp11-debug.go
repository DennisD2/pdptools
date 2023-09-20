package main

import "C"
import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// SubProcess connection to running pdp11-debug
type SubProcess struct {
	continueLoop int
	dryMode      bool
	debug        int
	batch        bool
	command      string
	args         string
	cmd          *exec.Cmd
	stdin        io.WriteCloser
	stdout       io.Reader
	stderr       io.Reader
	state        MachineState
	enterCommand bool
	windowSize   int
}

// MachineState state of running machine
type MachineState struct {
	sr int
	pc int
	sp int
	r  []int
}

func main() {
	var err error
	regs := make([]int, 6)

	state := MachineState{
		0,
		0,
		0,
		regs,
	}

	proc := SubProcess{
		1,
		true,
		0,
		false,
		"pdp11",
		".",
		nil,
		nil,
		nil,
		nil,
		state,
		false,
		80, /* 26 for 25 line, 80 for 80 lines terminal  */
	}

	// Create sub process structure
	proc.cmd = exec.Command(proc.command /*, proc.args*/)

	// Prepare handling pdp11-debug stdin; stdin will be served from localKeyboardReader()
	proc.stdin, err = proc.cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Handle local keyboard stdin: raw terminal input, handles by go routine
	if !proc.batch {
		// switch stdin into 'raw' mode
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}
	go proc.localKeyboardReader()

	// Handle pdp11-debug stdout in go routine
	proc.stdout, err = proc.cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go proc.stdoutReader()

	// Handle pdp11-debug stderr in go routine
	proc.stderr, err = proc.cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go proc.stderrReader()

	// Start sub process
	if err := proc.cmd.Start(); err != nil {
		return
	}

	// Init machine
	io.WriteString(proc.stdin, "do test-libgdd.deposit\n")
	log.Println("---------------------------\n")
	io.WriteString(proc.stdin, "break 1000\n")
	io.WriteString(proc.stdin, "run 1000\n")

	// Init machine state
	proc.dumpInfo(1500)

	// wait for terminate condition
	for proc.continueLoop > 0 {
		time.Sleep(25 * time.Millisecond)
	}
}

// stdoutReader Handle stdout
func (proc *SubProcess) stdoutReader() {
	scanner := bufio.NewScanner(proc.stdout)

	var cmd string
	var val = 0
	var changed = false

	for proc.continueLoop > 0 {

		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "PC:") {
				_, err := fmt.Sscanf(line, "%s %d", &cmd, &val)
				if err != nil {
					log.Printf("Error reading PC: ", err)
					val = 0
				}
				convertOctToDec(&val)
				if proc.debug > 1 {
					log.Printf("PC line, PC=%#o (%v dec)", val, val)
				}
				proc.state.pc = val
			}
			if strings.HasPrefix(line, "R") {
				_, err := fmt.Sscanf(line, "%s %d", &cmd, &val)
				if err != nil {
					log.Printf("Error reading register: ", err)
				}
				var reg, dummy, dummy2 rune
				_, err = fmt.Sscanf(cmd, "%c%c%c", &dummy, &reg, &dummy2)
				convertOctToDec(&val)
				registerIndex := reg - 48
				if proc.debug > 1 {
					log.Printf("Rx line, Rx=%d, val=%#o (%v dec)", registerIndex, val, val)
				}
				oldValue := proc.state.r[registerIndex]
				proc.state.r[registerIndex] = val
				if oldValue != val {
					changed = true
				} else {
					changed = false
				}
			}
			if strings.HasPrefix(line, "SP") {
				_, err := fmt.Sscanf(line, "%s %d", &cmd, &val)
				if err != nil {
					log.Printf("Error reading stackpointer: ", err)
				}
				convertOctToDec(&val)
				if proc.debug > 1 {
					log.Printf("SP line, SP=%#o (%v dec)", val, val)
				}
				oldValue := proc.state.sp
				proc.state.sp = val
				if oldValue != val {
					changed = true
				} else {
					changed = false
				}
			}
			if strings.HasPrefix(line, "SR") {
				_, err := fmt.Sscanf(line, "%s %d", &cmd, &val)
				if err != nil {
					log.Printf("Error reading statusregister: ", err)
				}
				convertOctToDec(&val)
				if proc.debug > 1 {
					log.Printf("SR line, SR=%#o (%v dec)", val, val)
				}
				oldValue := proc.state.sr
				proc.state.sr = val
				if oldValue != val {
					changed = true
				} else {
					changed = false
				}

			}
			printText := ""
			if changed {
				printText = fmt.Sprintf("%c[%c%c%s\r%c[%c%c\r", 0x1b, '7', 'm', line, 0x1b, '0', 'm')
			} else {
				if hasLeadLineNumber(line) && lineNumberEqualsPC(line, proc.state.pc) {
					printText = fmt.Sprintf("%c[%c%c%s\r%c[%c%c\r", 0x1b, '7', 'm', line, 0x1b, '0', 'm')
				} else {
					printText = line
				}
			}
			fmt.Println(printText + "\r")

		}
		if proc.continueLoop == 0 {
			return
		}
	}
}

func hasLeadLineNumber(line string) bool {
	// line starts with number after ^
	// number ends with ':'
	match, _ := regexp.MatchString("^[0-9]+:", line)
	return match
}

func lineNumberEqualsPC(line string, pc int) bool {
	var dummy string
	var val = 0

	line2 := strings.Replace(line, "\t", " ", -1)
	_, err := fmt.Sscanf(line2, "%d%s", &val, &dummy)
	if err != nil {
		log.Printf("Error reading address: ", err)
		val = 0
	}
	convertOctToDec(&val)
	return pc == val
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

// stdoutReader Handle stderr
func (proc *SubProcess) stderrReader() {
	scanner := bufio.NewScanner(proc.stderr)

	for proc.continueLoop > 0 {

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line + "\r")
		}
		if proc.continueLoop == 0 {
			return
		}
	}
}

// localKeyboardReader handles all local keyboard input and interaction
func (proc *SubProcess) localKeyboardReader() {
	cmdLine := make([]byte, 128)
	cmdLinePtr := 0

	cbuf := make([]byte, 128)
	consoleReader := bufio.NewReader(os.Stdin)
	b := make([]byte, 1)

	for proc.continueLoop > 0 {
		num, err := consoleReader.Read(cbuf)
		if err != nil {
			fmt.Println(err)
		}
		if num == 0 {
			continue
		}
		if num > 1 {
			// We currently cannot handle multiple chars at once
			fmt.Println("Multiple chars!")
		}
		b[0] = cbuf[0]
		if proc.debug > 0 {
			fmt.Printf("<%d:%s:%x>", num, b, b)
		} else {
			//fmt.Printf("%s", b)
		}
		if !proc.enterCommand {
			if b[0] == 's' {
				// 's' - Single Step
				log.Println("----- (s) Single Step --------------------------\n")
				io.WriteString(proc.stdin, "s\n")

				fmt.Printf("%c", 0x1b)
				fmt.Printf("%c", '[')
				fmt.Printf("%c", '2')
				fmt.Printf("%c", 'J')

				fmt.Printf("%c", 0x1b)
				fmt.Printf("%c", '[')
				fmt.Printf("%c", 'H')
				proc.dumpInfo(1)
			}

			if cbuf[0] == 'q' {
				log.Println("----- (q) Quit --------------------------\n")
				proc.continueLoop = 0
			}
		}

		if proc.enterCommand {
			if cbuf[0] == '\r' {
				proc.enterCommand = false
				fmt.Print(string(cbuf[0]))
				cmdString := fmt.Sprintf("%s\n", string(cmdLine))
				io.WriteString(proc.stdin, cmdString)
				proc.dumpInfo(1)
			} else {
				fmt.Print(string(cbuf[0]))
				cmdLine[cmdLinePtr] = cbuf[0]
				cmdLinePtr++
			}
		}

		if cbuf[0] == '>' {
			log.Println("----- (q) Enter SIMH Command-------------------\n")
			proc.enterCommand = true
			i := 0
			for i < len(cmdLine) {
				cmdLine[i] = 0
				i++
			}
			cmdLinePtr = 0
			fmt.Print(">")
		}

	}
}

// dumpInfo dumps info regarding PC, registers and assembler codes in PC vicinity
func (proc *SubProcess) dumpInfo(waitMillis int) {
	io.WriteString(proc.stdin, "ex pc\n")
	io.WriteString(proc.stdin, "ex r0-r5,sp,sr\n")
	time.Sleep(time.Duration(waitMillis) * time.Millisecond)
	var windowStart int
	var windowEnd int
	calcCodeWindowValues(proc.state.pc, proc.windowSize, &windowStart, &windowEnd, proc.debug)
	exLine := fmt.Sprintf("ex -m %#o-%#o\n", windowStart, windowEnd)
	io.WriteString(proc.stdin, exLine)
}

// calcCodeWindowValues calculate useful window (startaddress,endaddress) for code to dump
func calcCodeWindowValues(pc int, windowSize int, windowStart *int, windowEnd *int, debug int) {
	if pc-windowSize <= 0 {
		*windowStart = 0
	} else {
		*windowStart = pc - windowSize
	}
	*windowEnd = *windowStart + 2*windowSize
	if debug > 0 {
		log.Printf("Window %#o to %#o\n", *windowStart, *windowEnd)
	}
}
