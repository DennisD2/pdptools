# all2deposit

Creates a SIMH deposit file from an input file
input file format can be:
* PTAP (paper tape) format
* a.out format (not yet implemented)

## Conversion step
Convert ```hello.ptap```, deposit output file is 
```hello.deposit```:
```bash
all2deposit --in bare-metal-macro11-examples-as-clang/hello.ptap >hello.deposit
# or execute without go build step directly from sources
go run . --in bare-metal-macro11-examples-as-clang/hello.ptap >hello.deposit
```

For Makefile automation, see my [Makefile](../bare-metal-examples-clang/Makefile)

## Running deposit file (or papertape file) in SIMH
SIMH use for PTAP and DEPOSIT file:
```bash
% pdp11 hello.deposit
```
Then execute ```run 1000``` in simh console.

```bash
% pdp11
PDP-11 simulator V4.0-0 Current        git commit id: d5cc3406
sim> load bare-metal-macro11-examples-as-clang/hello.ptap
```
Then execute ```run``` in simh console.

## Important simh commands
```
// load DEPOSIT file
sim> do hello.deposit

// breakpoint
sim> br 1000
// run code, use PC content as starting point
sim> run

// run code 
sim> go 1000

Breakpoint, PC: 001000 (MOV #776,SP)
// Single step
sim> s
Step expired, PC: 001004 (CLR R0)
// examine some register
sim> ex r0
R0:     000010
// examine memory content
sim> ex 1010-1050
1010:   020027
1012:   000620
1014:   001374
// continue run
sim> cont
// show breakpoints
sim> show break
Supported Breakpoint Types: -E -P -R -S -W -X
1000:   E
1016:   E
1044:   E
// remove breakpoint#
sim> nobr  1044
// Examine with and w/o ASCII flag
sim> ex -a 1404-1422
1404:   H
1405:   e
1406:   l
1407:   l
1410:   o
1411:
1412:   W
1413:   o
1414:   r
1415:   l
1416:   d
1417:   !
1420:   <015>
1421:   <012>
1422:   <000>

sim> ex -m ... # will show mnemonics

sim> ex  1404-1422
1404:   062510
1406:   066154
1410:   020157
1412:   067527
1414:   066162
1416:   020544
1420:   005015
1422:   000000

// Examine all CPU details
sim> ex state
PC:     000104
R0:     000000
R1:     000000
R2:     000000
R3:     000000
R4:     000000
...
```

## Using second terminal to separate SIMH output and program output
Install a telnet client. E.g. putty, which has a Linux version.
Then define a new session, lets say named ```simh```, of
type ```telnet```, connecting to port 8023 on localhost.

Then start simh, which is able to open a telnet server port
where the terminal output of the running simulation is provided.
Onfigure simh to provide the telnet server port:
```shell
sim> set console telnet=localhost:8023
Listening on port LOCALHOST:8023
sim> set console pchar=37777777777
```

Then try the connect with putty. A window should open and
the program output should go there. simh commands can still
be entered in the first terminal session.

### Related
See also https://mdfs.net/Docs/Comp/PDP11/SIMHpuTTY

## Bare metal approach
I am following mainly these resources:

* DEC-11-XPTSA-B-D.pdf "PDP-11 Paper Tape Software Handbook"
* Bare metal compiling witch GCC
  https://www.teckelworks.com/2020/03/c-programming-on-a-bare-metal-pdp-11/ and https://xw.is/wiki/Bare_metal_PDP-11_GCC_9.3.0_cross_compiler_instructions
* Also info in bare metal compiling
  https://github.com/JamesHagerman/gcc-pdp11-aout/blob/master/README.md
* Insights on PDP11 programming
  https://www.learningpdp11.com/

## Open issues
* cleanup code by using the file creation and emit code, see aout.go.
* make aout part work

## Closed issues
* Fixed checksum calculation


## Related
* Bare metal intro https://www.teckelworks.com/2020/03/c-programming-on-a-bare-metal-pdp-11/
* Cross compiling - https://xw.is/wiki/Bare_metal_PDP-11_GCC_9.3.0_cross_compiler_instructions
* baremetal topic - https://stackoverflow.com/questions/38387155/beginner-bare-metal-pdp11-console-output

* Discussion Load binary PDP files https://groups.google.com/g/pidp-11/c/AZkNx5FUpRY?pli=1
* Paper tape format info http://avitech.com.au/?page_id=709
* Old UNIX on PDP11 https://pdos.csail.mit.edu/6.828/2005/homework/hw5.html

* Modern PDP-11 C Compilers - https://retrocomputingforum.com/t/modern-pdp-11-c-compilers/2329/4
* Search PDP11 bare metal - https://www.google.com/search?q=pdp11+bare+metal+clang+compile&oq=pdp11+bare+metal+clang+compile&aqs=chrome..69i57j33i10i160l3.8962j0j4&sourceid=chrome&ie=UTF-8
