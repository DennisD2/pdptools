# pdpcom
This software can connect to a PDP11 and serve as a dumb terminal to it. It also allows to upload 
machine code in batch mode, created by e.g. MACRO-11 PDP11 assembler. 
The machine code need to be in SIMH deposit format. 
This can be created with my other tool, called lst2obj. TODO: give link

*Normal Mode* is like a dumb Terminal to PDP. There is a *Command Mode* to enter commands.

# How to build
```shell
make
```
Check Makefile for more infos regarding the simple build.

### Requisites
Uses os.unix packages of Go, so it will work only on Linux OS.

# How to use
Command line to start interactively:
```shell
pdpcom --device=/dev/ttyUSB1

# Or use arguments and sources
go run pdpcom.go pdpConnection.go \
   --device=/dev/ttyUSB1 --debug=1 --dry-run=false --baudrate=9600
```

Command line for a batch upload of a deposit file:
```shell
pdpcom --batch --upload 064-backward-echo.deposit 
```

All command line args and default values:
```bash
--device, TTY Device: /dev/ttyUSB0
--dry-run: false
--debug: 0
--baudrate: 9600
--batch: false
--upload: multiecho.deposit
```

Enter command mode with ```:```. Commands:
* ':' returns to "normal" ODT input mode
* 'q': quit
* 'r': read and upload PDP11 binary (in SIMH deposit format) to PDP. 
  After the ```r```, enter a ```SPACE``` and then path to file to upload.

Small video clip where I list some memory content, load an object file, show that is has been loaded to memory and
execute it (this simple object file waits for a single char and prints it to console):

https://user-images.githubusercontent.com/7112686/166147322-41d7fa8d-4714-4125-a60a-cdef343112e0.mp4
TODO; update video

# How it works
pdpcom handles keyboard input/output on a PC with a Go routine.
The in/output of the PDP machine is handled with a second Go routine.
