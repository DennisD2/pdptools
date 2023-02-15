all:
	go build

clean:
	rm lst2deposit

init:
	go mod init lst2deposit
	go mod tidy

test:
	pdp11-aout-gcc -nostdlib -Ttext 0x200 -m10 -Os -N -e _start crt0.s console.c hello.c -o hello
