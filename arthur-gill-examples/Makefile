# macro11 options - Create absolute addresses (no relocable addresses, because we use no linker)
MOPTS		= -e AMA
LST2DEPOSITPATH=../lst2deposit
#MOPTS		=

%.lst : %.mac
	macro11 $(MOPTS) -l $@ $<

%.deposit : %.lst
	$(LST2DEPOSITPATH)/lst2deposit --lst $< --debug 1 --out $@

all: 039-echo.deposit 040-a-to-z-to-mem.deposit 054-multiecho.deposit 064-backward-echo.deposit \
	104-ascii-to-binary-conversion.deposit 117-intype.deposit 118-inbell.deposit \
	123-time.deposit

039-echo.lst: 039-echo.mac
040-a-to-z-to-mem.lst: 040-a-to-z-to-mem.mac
054-multiecho.lst: 054-multiecho.mac
064-backward-echo.lst: 064-backward-echo.mac
104-ascii-to-binary-conversion.lst: 104-ascii-to-binary-conversion.mac
117-intype.lst:  117-intype.mac
118-inbell.lst: 118-inbell.mac
123-time.lst: 123-time.mac

clean:
	rm *.lst *.deposit
