# Makefile for testing base metal compile for PDP11 with gcc cross compiler
# run "make test"

all: compile disassemble

compile:
	# -nostdlib - do not link a standard c library (libc.a) to result
	# -Ttext 0x200 - start text (code) section at 0x200 = 1000octal
	# -m10 - limit the assembly instructions generated to those of a PDP 11/10 (or 11/05)
	# -Os - optimize for space
	# -N - tells the linker put the data section immediately after the code section and not on its own page
	# -e _start -  tells the linker to start execution at the _start function in crt0.s.
	pdp11-aout-gcc -c -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.s console.c ltc.c
	# Link step
	pdp11-aout-gcc -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.o console.o ltc.o -o ltc
	# create PTAP file from a.out binary
	aout2lda --aout ltc --lda ltc.ptap --data-align 2 --text 0x200 --vector0
	# Create deposit file from PTAP file
	all2deposit --in ltc.ptap >ltc.deposit

# disassemble created object files
disassemble:
	pdp11dasm crt0.o
	pdp11dasm console.o
	pdp11dasm ltc.o

clean:
	rm *.o ltc *.ptap *.deposit *.das
