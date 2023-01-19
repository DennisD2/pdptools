package main

// #include <termios.h>
// #include <unistd.h>
import "C"

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
)

// PDP11SerialConnection TTY connection to PDP
type PDP11SerialConnection struct {
	tty      *os.File
	device   string
	baudrate int
	timeout  time.Duration
	//continueLoop int
	//state        ConnState
	//dryMode      bool
	//debug        int
	//batch        bool
	//uploadFile   string
}

// ConnState State of Connection
type ConnState int

const (
	ODTNormal      ConnState = 0
	CommandMode              = 1
	CommandCollect           = 2
)

// openTTY opens TTY connection
func (pdp *PDP11SerialConnection) openTTY() (err error) {
	tty, err := os.OpenFile(pdp.device, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
	if err != nil {
		return
	}

	fd := C.int(tty.Fd())
	if C.isatty(fd) != 1 {
		tty.Close()
		return errors.New("File is not a tty")
	}

	var st C.struct_termios
	_, err = C.tcgetattr(fd, &st)
	if err != nil {
		tty.Close()
		return err
	}
	var speed C.speed_t
	switch pdp.baudrate {
	case 115200:
		speed = C.B115200
	case 57600:
		speed = C.B57600
	case 38400:
		speed = C.B38400
	case 19200:
		speed = C.B19200
	case 9600:
		speed = C.B9600
	case 4800:
		speed = C.B4800
	case 2400:
		speed = C.B2400
	default:
		tty.Close()
		return fmt.Errorf("Unknown baud rate %s", pdp.baudrate)
	}

	_, err = C.cfsetispeed(&st, speed)
	if err != nil {
		tty.Close()
		return err
	}
	_, err = C.cfsetospeed(&st, speed)
	if err != nil {
		tty.Close()
		return err
	}

	// Turn off break interrupts, CR->NL, Parity checks, strip, and IXON
	st.c_iflag &= ^C.tcflag_t(C.BRKINT | C.ICRNL | C.INPCK | C.ISTRIP | C.IXOFF | C.IXON | C.PARMRK)

	// Select local mode, turn off parity, set to 8 bits
	st.c_cflag &= ^C.tcflag_t(C.CSIZE | C.PARENB)
	st.c_cflag |= (C.CLOCAL | C.CREAD | C.CS8)

	// Select raw mode
	st.c_lflag &= ^C.tcflag_t(C.ICANON | C.ECHO | C.ECHOE | C.ISIG)
	st.c_oflag &= ^C.tcflag_t(C.OPOST)

	// set blocking / non-blocking read
	/*
	*	http://man7.org/linux/man-pages/man3/termios.3.html
	* - Supports blocking read and read with timeout operations
	 */
	/*vmin, vtime := posixTimeoutValues(readTimeout)
	st.c_cc[C.VMIN] = C.cc_t(vmin)
	st.c_cc[C.VTIME] = C.cc_t(vtime)*/
	st.c_cc[C.VMIN] = C.cc_t(1)
	st.c_cc[C.VTIME] = C.cc_t(0)

	_, err = C.tcsetattr(fd, C.TCSANOW, &st)
	if err != nil {
		tty.Close()
		return err
	}

	//fmt.Println("Tweaking", name)
	r1, _, e := syscall.Syscall(syscall.SYS_FCNTL,
		uintptr(tty.Fd()),
		uintptr(syscall.F_SETFL),
		uintptr(0))
	if e != 0 || r1 != 0 {
		s := fmt.Sprint("Clearing NONBLOCK syscall error:", e, r1)
		tty.Close()
		return errors.New(s)
	}

	/*
		r1, _, e = syscall.Syscall(syscall.SYS_IOCTL,
					uintptr(tty.Fd()),
					uintptr(0x80045402), // IOSSIOSPEED
					uintptr(unsafe.Pointer(&baud)));
			if e != 0 || r1 != 0 {
					s := fmt.Sprint("Baudrate syscall error:", e, r1)
			tty.Close()
						return nil, os.NewError(s)
		}
	*/
	pdp.tty = tty
	return nil
}
