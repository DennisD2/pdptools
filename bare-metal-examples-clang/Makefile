# Makefile for testing bae metal compile for PDP11 with gcc cross compiler
# run "make test"

test:
	# -nostdlib - do not link a standard c library (libc.a) to result
	# -Ttext 0x200 - start text (code) section at 0x200 = 1000octal
	# -m10 - limit the assembly instructions generated to those of a PDP 11/10 (or 11/05)
	# -Os - optimize for space
	# -N - tells the linker put the data section immediately after the code section and not on its own page
	# -e _start -  tells the linker to start execution at the _start function in crt0.s.
	pdp11-aout-gcc -c -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.s console.c hello.c
	# next line should use '.o instead of *.c
	pdp11-aout-gcc -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.o console.o hello.o -o hello
	# create PTAP file from a.out binary
	aout2lda --aout hello --lda hello.ptap --data-align 2 --text 0x200 --vector0
	# Create deposit file from PTAP file
	../all2deposit --in hello.ptap >hello.deposit

# disassemble created object files
disassemble:
	pdp11dasm crt0.o
	pdp11dasm console.o
	pdp11dasm hello.o

clean:
	rm *.o hello *.ptap *.deposit
