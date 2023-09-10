all: helloworld time

helloworld:
	pdp11-aout-as putconch.s -o putconch.o
	pdp11-aout-as hellopdp.s -o hellopdp.o
	pdp11-aout-ld -T ldaout.cmd hellopdp.o putconch.o -o hellopdp.out
	aout2lda --aout hellopdp.out --lda hellopdp.ptap --data-align 2 --text 0x200 --vector0
	all2deposit --in hellopdp.ptap >hellopdp.deposit

time:
	pdp11-aout-as time.s -o time.o
	pdp11-aout-ld -T ldaout.cmd time.o -o time.out
	aout2lda --aout time.out --lda time.ptap --data-align 2 --text 0x200 --vector0
	pdp11dasm time.o
	all2deposit --in time.ptap >time.deposit