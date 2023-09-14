# pdptools
Tools to communicate with and upload code to physical PDP11 connected via RS232.

* lst2deposit - creates a SIMH deposit file from a MACRO11 LST file
* pdpcom - allows to upload deposit file content to a serial attached PDP11 
* arthur-gill-examples - some nice examples
* examples - some simple examples

# How to Build
See sub projects readme files,  [lst2deposit/README.md](lst2deposit/README.md) 
and [pdpcom/README.md](pdpcom/README.md)

# How to use
See sub projects readme files,  [lst2deposit/README.md](lst2deposit/README.md)
and [pdpcom/README.md](pdpcom/README.md)

### Screencast demo session
Small video clip where I list some memory content, load an object file, show that is has been loaded to memory and
execute it (this simple object file waits for a single char and prints it to console):

[![https://youtu.be/ZbehDUC4euI](https://img.youtube.com/vi/ZbehDUC4euI/0.jpg)](https://www.youtube.com/watch?v=ZbehDUC4euI)

[https://youtu.be/ZbehDUC4euI](https://youtu.be/ZbehDUC4euI)

# Example use cases
See directory [arthur-gill-examples](arthur-gill-examples) for examples I typed in from
Arthur Gills great book.

See [arthur-gill-examples/Makefile](arthur-gill-examples/Makefile) on how to use it.

Some other random Macro11 examples are in [macro11-examples](macro11-examples).

# References
* MACRO-11 assembler for Linux - https://github.com/andpp/macro11
* MACRO-11 assembler manual - http://www.bitsavers.org/www.computer.museum.uq.edu.au/RSX/AA-V027A-TC%20PDP-11%20MACRO-11%20Language%20Reference%20Manual.pdf
* Octal/HEX/DEC/* online calculator - https://www.rapidtables.com/convert/number/hex-dec-bin-converter.html
* ASCII Code table - https://www.physics.udel.edu/~watson/scen103/ascii.html

* Machine and Assembly Language Programming of the PDP-11, Arthur Gill, Prentice-Hall, 1978.
  This book was written by Arthur who gave lessons on University of Berkeley. A very nice
  book for studying PDP-11 assembler. Used prints can still be found.

# Differences MACRO11 and Gnu Assembler dialect for PDP11

There are some differences in writiung down assembler code between the
original MACRO11 dialect and the Gnu Assembler dialect.

Some of them I found and did notice them here.

## Addressing

| MACRO11        | GAS             | Meaning                                     |
|----------------|-----------------|---------------------------------------------|
| MOVB ITIME,R0  | movb *$itime,r0 | Read content at address ITIME               |
| BIC #177760,R0 | bic $0177760,r0 | Literal (octal) value 0177760               |
| MOV #MESSG,R0  | mov $messg,r0   | Literal value (here: address of some label) |

## Numbers

Default number system (number without any modifier) for MACRO11 is *octal*. 
Default number system (number without any modifier) for GAS is *decimal*.

| MACRO11 | GAS  | Meaning           |
|--------|------|-------------------|
| #12.   | $12  | Decimal number 12 |
| #12    | $012 | Octal number 12   |

## Absolute addresses

in MACRO11 we can write:
```
        LC=.
        .=4+LC
        .WORD    6,0,12,0           ; INITIALIZE ERROR VECTORS
```
This would put 6,0,12,0 values to address 4,6,8,10 if "." points to address 0.
In GAS, the value "." comes from outside and may have any value. I do not
know how to do a valid set up there. For now, I usually change the MACRO11 code
to something like this:
``` 
vecs:   mov     $6,*$4              # initialize error vectors
        mov     $0,*$6
        mov     $012,*$010
        mov     $0,*$012
```
And that piece of code comes to the start section of the code.
So instead of initializing the memory during loading the binary code,
it will be initialized on startup of the code.