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

Some other random examples are in [examples](examples).

# References
* MACRO-11 assembler for Linux - https://github.com/andpp/macro11
* MACRO-11 assembler manual - http://www.bitsavers.org/www.computer.museum.uq.edu.au/RSX/AA-V027A-TC%20PDP-11%20MACRO-11%20Language%20Reference%20Manual.pdf
* Octal/HEX/DEC/* online calculator - https://www.rapidtables.com/convert/number/hex-dec-bin-converter.html
* ASCII Code table - https://www.physics.udel.edu/~watson/scen103/ascii.html

* Machine and Assembly Language Programming of the PDP-11, Arthur Gill, Prentice-Hall, 1978.
  This book was written by Arthur who gave lessons on University of Berkeley. A very nice
  book for studying PDP-11 assembler. Used prints can still be found.