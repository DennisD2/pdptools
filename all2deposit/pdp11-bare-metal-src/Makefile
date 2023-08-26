
test:
	pdp11-aout-gcc -c -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.s console.c hello.c
	pdp11-aout-gcc -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.s console.c hello.c -o hello
	aout2lda --aout hello --lda hello.ptap --data-align 2 --text 0x200 --vector0
	# ../aout2deposit -aout hello -data-align 2 --text 0x200

clean:
	rm *.o hello *.ptap