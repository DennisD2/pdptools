all:
	go build

clean:
	rm all2deposit

init:
	go mod init aout2deposit
	go mod tidy

# Create an a.out format ´hello' and a PTAP file
test:
	pdp11-aout-gcc -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.s console.c hello.c -o hello
	aout2lda --aout hello --lda hello.ptap --data-align 2 --text 0x200 --vector0

install:
	go build
	cp -av all2deposit $(HOME)/bin