#
# Example Makefile that shows how a assembler code calls clang code and vice versa
#

# c code using an assembler provided function
clang_uses_asm_fun:
	# assemble assembler code
	pdp11-aout-as prov_fun.s -o prov_fun.o
	pdp11dasm prov_fun.o
	# compile c code
	pdp11-aout-gcc -c -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.s console.c use_asm_fun.c
	pdp11dasm use_asm_fun.o
	# link all object files
	pdp11-aout-gcc -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.o prov_fun.o use_asm_fun.o -o c_use_asm_fun
	# create PTAP file from a.out binary
	aout2lda --aout c_use_asm_fun --lda c_use_asm_fun.ptap --data-align 2 --text 0x200 --vector0
	# Create deposit file from PTAP file
	all2deposit --in c_use_asm_fun.ptap >c_use_asm_fun.deposit

# assembler code using a c provided function
asm_uses_clang_fun:
	# assemble assembler code
	pdp11-aout-as use_c_fun.s -o use_c_fun.o
	pdp11dasm use_c_fun.o
	# compile c code
	pdp11-aout-gcc -c -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.s console.c provide_asm_fun.c
	pdp11dasm provide_asm_fun.o
	# link all object files
	pdp11-aout-gcc -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.o  console.o use_c_fun.o provide_asm_fun.o -o asm_use_c_fun
	pdp11dasm asm_use_c_fun
	pdp11-aout-objdump -D asm_use_c_fun >asm_use_c_fun.objdump
	# create PTAP file from a.out binary
	aout2lda --aout asm_use_c_fun --lda asm_use_c_fun.ptap --data-align 2 --text 0x200 --vector0
	# Create deposit file from PTAP file
	all2deposit --in asm_use_c_fun.ptap >asm_use_c_fun.deposit

# assembler code using an assembler provided function
asm_prov_and_use:
	# assemble assembler code
	pdp11-aout-as prov_fun_and_use_fun.s -o prov_fun_and_use_fun.o
	pdp11dasm prov_fun_and_use_fun.o
	# Create a.out type file consumable by aout2lda
	pdp11-aout-ld -T ldaout.cmd prov_fun_and_use_fun.o -o prov_fun_and_use_fun.out
	# Create PTAP file from a.out file
	aout2lda --aout prov_fun_and_use_fun.out --lda prov_fun_and_use_fun.ptap --data-align 2 --text 0x200 --vector0
	# Create deposit file from PTAP
	all2deposit --in prov_fun_and_use_fun.ptap >prov_fun_and_use_fun.deposit

libgdd-test:
	# assemble assembler code
	pdp11-aout-as libgdd.s -o libgdd.o
	pdp11dasm libgdd.o
	# compile c code
	pdp11-aout-gcc -c -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.s console.c test-libgdd.c
	pdp11dasm test-libgdd.o
	# link all object files
	pdp11-aout-gcc -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.o console.o libgdd.o test-libgdd.o -o test-libgdd
	# create PTAP file from a.out binary
	aout2lda --aout test-libgdd --lda test-libgdd.ptap --data-align 2 --text 0x200 --vector0
	# Create deposit file from PTAP file
	all2deposit --in test-libgdd.ptap >test-libgdd.deposit

clean:
	rm *.o *.deposit *.das c_use_asm_fun asm_use_c_fun test-libgdd