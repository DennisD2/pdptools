# macro11 options - Create absolute addresses (no relocable addresses, because we use no linker)
MOPTS		= -e AMA
#MOPTS		=
LST2DEPOSITPATH=../lst2deposit

%.lst : %.mac
	macro11 $(MOPTS) -l $@ $<

%.deposit : %.lst
	$(LST2DEPOSITPATH)/lst2deposit --lst $< --debug 1 --out $@

all: echo-end-x.deposit bell-and-intype.deposit binary-division.deposit div-and-mul-tests.deposit \
     2processes.deposit

echo-end-x.lst: echo-end-x.mac
bell-and-intype.lst: bell-and-intype.mac
binary-division.lst: binary-division.mac
div-and-mul-tests.lst: div-and-mul-tests.mac
2processes.lst: 2processes.mac

clean:
	rm *.lst *.deposit
