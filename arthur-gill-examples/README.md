# Example listings from Arthur Gills Book "Machine And Assembly Language Programming Of The PDP-11"
This directory contains some examples from that venerable book.

The examples are nearly 1:1 typed in from the book. I have only replaced some
absolute addresses with labels to make them more understandable.

The file names follow the syntax ```<book page number>-<example short name>.<extension>```.

* Simple Echo, page 39
* Deposit A-Z chars in memory, page 40
* Multiecho, page 54
* Backward echo, page 64
* ASCII to binary conversion, page 104
* Echoing using keyboard interrupt, page 117
* Ringing the bell each 10 seconds, using time line clock interrupt, page 118
* Time, a clock program, page 123

For examples using the timing (LTC) interrupt: To execute these files with simh PDP11 simulation, 
you need to set the cpu to 11/23. Otherwise, the LTC mechanism will not work.

A complete execution of one example looks like this:

```shell
# make all examples
% make
# execute machine code with pdp11 simh simulator
% pdp11 123-time.deposit 

PDP-11 simulator V4.0-0 Current        git commit id: d5cc3406
sim> set cpu 11/23
sim> run 500

WHAT TIME IS IT? 1230
AT THE BELL THE TIME WILL BE: 12:30:06
AT THE BELL THE TIME WILL BE: 12:30:09
...
```
