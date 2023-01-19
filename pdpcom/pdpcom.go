package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"os"

	"golang.org/x/term"
)

// PDP11Connection connection to PDP machine
type PDP11Connection struct {
	continueLoop int
	state        ConnState
	dryMode      bool
	debug        int
	batch        bool
	uploadFile   string
	serial       *PDP11SerialConnection
}

func main() {
	fmt.Println("PDP Communicator")
	devicePtr := flag.String("device", "/dev/ttyUSB0",
		"TTY device used to access PDP")
	dryRunPtr := flag.Bool("dry-run", false,
		"Debug level")
	debugPtr := flag.Int("debug", 0,
		"Debug level")
	baudratePtr := flag.Int("baudrate", 9600,
		"Baudrate")
	batchPtr := flag.Bool("batch", false,
		"Non-interactive (batch) mode")
	uploadPtr := flag.String("upload", "multiecho.deposit",
		"Deposit file to upload")
	flag.Parse()

	fmt.Printf("--device, TTY Device: %s\n", *devicePtr)
	fmt.Printf("--dry-run: %t\n", *dryRunPtr)
	fmt.Printf("--debug: %d\n", *debugPtr)
	fmt.Printf("--baudrate: %d\n", *baudratePtr)
	fmt.Printf("--batch: %t\n", *batchPtr)
	fmt.Printf("--upload: %s\n", *uploadPtr)

	// Create serial connection
	pdpSerial := PDP11SerialConnection{
		nil, //priv
		*devicePtr,
		*baudratePtr,
		0,
	}

	// Create Connection
	pdp := PDP11Connection{
		1,         //priv
		ODTNormal, //priv
		*dryRunPtr,
		*debugPtr,
		*batchPtr,
		*uploadPtr,
		&pdpSerial,
	}

	if !pdp.dryMode {
		// open tty reader
		err := pdp.serial.openTTY()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer pdp.serial.tty.Close()
	}

	if !pdp.batch {
		// switch stdin into 'raw' mode
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)

		// Start local keyboard handler routine
		go localKeyboardReader(&pdp)

		// Start tty routine
		go ttyReader(&pdp)

		// stay in loop until end condition is met
		for pdp.continueLoop > 0 {
			time.Sleep(25 * time.Millisecond)
		}
	} else {
		// just upload
		uploadFile(&pdp, *uploadPtr, pdp.debug)
	}

	fmt.Println("\n\rQuitting PDP Communicator\n\r")
}

func ttyReader(pdp *PDP11Connection) {
	cbuf := make([]byte, 128)
	for pdp.continueLoop > 0 {
		if pdp.dryMode {
			continue
		}
		// check PDP tty
		_, err := pdp.serial.tty.Read(cbuf)
		if err != nil {
			fmt.Printf("Error in Read: %s\n", err)
			pdp.continueLoop = 0
		} else {
			c := cbuf[0]
			if pdp.debug > 0 {
				fmt.Printf("%c<%x>", c, c)
			} else {
				fmt.Printf("%c", c)
			}
		}
	}
}

// localKeyboardReader handles all local keyboard input and interaction
func localKeyboardReader(pdp *PDP11Connection) {
	cbuf := make([]byte, 128)
	cmdLine := make([]byte, 128)
	cmdLinePtr := 0
	consoleReader := bufio.NewReader(os.Stdin)
	b := make([]byte, 1)
	for pdp.continueLoop > 0 {
		num, err := consoleReader.Read(cbuf)
		if err != nil {
			fmt.Println(err)
		} else {
			if num == 0 {
				continue
			}
			if num > 1 {
				// We currently cannot handle multiple chars at once
				fmt.Println("Multiple chars!")
			}

			if pdp.state == CommandCollect {
				b[0] = cbuf[0]
				if pdp.debug > 0 {
					fmt.Printf("<%d:%s:%x>", num, b, b)
				} else {
					fmt.Printf("%s", b)
				}
				if cbuf[0] == ' ' {
					// ignore spaces
					continue
				}
				if cbuf[0] == 0x0d {
					fmt.Printf("\r\n%s\r\n", cmdLine)
					uploadFile(pdp, string(cmdLine), pdp.debug)
					//uploadFile(pdp, "macro11-examples/char-in.obj")
					pdp.state = ODTNormal
					fmt.Println("\n\rBack to ODT\n\r")
				} else {
					cmdLine[cmdLinePtr] = cbuf[0]
					cmdLinePtr++
				}
				continue
			}
			//fmt.Printf("%d - %s", num, cbuf)
			if pdp.state == CommandMode {
				fmt.Printf("%s", cbuf)
				// In command mode, execute command based on key input
				if cbuf[0] == ':' {
					pdp.state = ODTNormal
					fmt.Println(" Back to ODT\n\r")
					continue
				}
				if cbuf[0] == 'q' {
					pdp.continueLoop = 0
					pdp.state = ODTNormal
				}
				if cbuf[0] == 't' {
					fmt.Println("Poke test...\r")
					pdp.testPoke("000000", "000123")
					pdp.state = ODTNormal
				}
				if cbuf[0] == 'r' {
					// collect complete command line (until <CR>)
					pdp.state = CommandCollect
					cmdLine[0] = 0
					cmdLinePtr = 0
				}
				continue
			}

			if pdp.state == ODTNormal {
				if cbuf[0] == ':' {
					// If ':' is selected, check next char for command to execute
					// We switch state to CommandMode for that
					pdp.state = CommandMode
					fmt.Print("Command (:qrtu):")
					continue
				}
			}
			// Normal ODT input, forward it to tty
			b[0] = cbuf[0]
			if !pdp.dryMode {
				if pdp.debug > 0 {
					fmt.Printf("<%d:%s:%x>", num, b, b)
				} else {
					pdp.serial.tty.Write(b)
				}
			}
		}
	}
}

// writeChar writes a single char to PDP
func (pdp PDP11Connection) writeChar(c byte) {
	b := make([]byte, 1)
	b[0] = c
	if !pdp.dryMode {
		pdp.serial.tty.Write(b)
	}
}

// pdpPoke writes a value word to a PDP address. Flag endLine defines what command end is given (LF or CR)
func (pdp PDP11Connection) pdpPoke(adr string, value string, endLine bool) {
	if pdp.dryMode {
		fmt.Printf("dry mode, not writing: %s/%s\n\r", adr, value)
		return
	}

	i := 0
	// Write address
	for i < len(adr) {
		pdp.writeChar(adr[i])
		i++
	}
	// finish address with '/'
	pdp.writeChar('/')

	// Wait some time for output reaction of PDP
	time.Sleep(time.Duration(20) * time.Millisecond)

	// Write value, if not ""
	i = 0
	for i < len(value) {
		pdp.writeChar(value[i])
		i++
	}
	// finish value with '\n' (move to next address) or '\r' (end of input)
	if endLine {
		pdp.writeChar('\r')
	} else {
		pdp.writeChar('\n')
	}
}

// uploadFile uploads machine code to PDP11
func uploadFile(pdp *PDP11Connection, file string, debug int) {
	// create useful file argument, is not allowed to contain 0 in its whole length
	// See here: https://github.com/golang/go/issues/24195
	i := 0
	var sb strings.Builder
	for i < len(file) && file[i] != 0 {
		sb.WriteByte(file[i])
		i++
	}

	// Read on object file
	var sb1 strings.Builder
	sb1.WriteString(sb.String())
	bytes, err := ioutil.ReadFile(sb1.String())
	if err != nil {
		log.Printf("Error loading file %s\n\r", err)
		return
	}

	i = 0
	var line strings.Builder
	for i < len(bytes) {
		// Read one line
		for bytes[i] != 0xa {
			line.WriteByte(bytes[i])
			i++
		}
		i++

		fmt.Printf("%s\n\r", line.String())

		var adr, value int
		fmt.Sscanf(line.String(), "D %d %d", &adr, &value)
		fmt.Printf("Read values: adr=%d value=%d\n\r", adr, value)

		aStr := fmt.Sprintf("%d", adr)
		vStr := fmt.Sprintf("%d", value)
		pdp.pdpPoke(aStr, vStr, true)
		time.Sleep(10 * time.Millisecond)

		line.Reset()
	}
}

func (pdp PDP11Connection) testPoke(adr string, value string) {
	pdp.pdpPoke(adr, value, false)
}
