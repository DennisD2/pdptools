# Assembler calls to C and vice versa

## Assembler calls assembler
Assembler code calls an assembler function.

Compile with 
```
make asm_prov_and_use
```

### Call code
```
mov     $20,-(sp)      # push literal 20 dec. to stack    
mov     $30,-(sp)      # push literal 30 dec. to stack   
jsr     pc,_addasm          
add     $4,sp          # fix stack pointer, i.e. remove parameters 20 and 30 from stack
```

### Function code
```
_addasm:
        clr     r0
        add     2(sp),r0  # reference literal 30 on stack
        add     4(sp),r0  # reference literal 20 on stack
        rts     pc
```

Invoked files: prov_fun_and_use_fun.s

## C calls assembler
Compile with
```
make clang_uses_asm_fun
```

### Call code
```
int result = addasm(21,31);
```

Disassembled C code:
```
000024: 012746 000037       	mov	#37,-(r6)		; f...
000030: 012746 000025       	mov	#25,-(r6)		; f...
000034: 004767 177760       	call	20			; w.p.
000040: 062706 000004       	add	#4,r6			; Fe..
```
The octal values 037 (=3*8+7=31) and 025 (=2x8+5=21) are pushed to stack (r6 is sp register).
Stack is corrected after subroutine call by 4. 

### Function code
```
        .globl _addasm

...

_addasm:
        clr     r0
        add     2(sp),r0
        add     4(sp),r0
        rts     pc
```
In assembler code, we just read the two values on stack and add to r0.
0(sp), the last value on stack, is the return address for the subroutione call.
So we need to access 2(sp) and 4(sp) to get the values.

Note that the C function is called ```addasm```, while the
assembler code exposes a global symbol ```_addasm```.

Invoked files: prov_fun.s, use_asm_fun.c, (and crt0.s, console.c)

## Assembler calls C
Compile with
```
make asm_uses_clang_fun
```

### Call code
```
        .extern _printLine      # clang function declared as extern
        
        ...
        
        jsr     pc,_printLine   # calling the clang function
```

### Function code
```
void printLine() {
    cons_puts("Hello World!\r\n");
}
```
Disassembled code for the C function:
```
000032: 012746 000060       	mov	#60,-(r6)		; f.0.
000036: 004767 177756       	call	20			; w.n.
000042: 062706 000002       	add	#2,r6			; Fe..
000046: 000207              	ret				; ..
```
In the disassembled code for the C function, it looks that #60 is the string address.

Again, note that the C function is called ```printLine```, while the
assembler code references an external symbol ```_printLine```.

**TODO: For a better example, have a function with a single int parameter**

Invoked files: use_c_fun.s, provide_asm_fun.c, (and crt0.s , console.c)

## Related

* How to Mix C and Assembly for X86 CPU, https://www.devdungeon.com/content/how-mix-c-and-assembly
* Another in depth explanation for X86 CPU https://en.wikibooks.org/wiki/X86_Assembly/GNU_assembly_syntax
* GNU Linker Manual https://ftp.gnu.org/old-gnu/Manuals/ld-2.9.1/, (and all manuals: https://ftp.gnu.org/old-gnu/Manuals/)
* GNU AS manual https://ftp.gnu.org/old-gnu/Manuals/gas-2.9.1/
* Best site for Assembler programming learning IMHO - https://www.chibiakumas.com/pdp11/
* SIMH Manual - http://simh.trailing-edge.com/pdf/simh_doc.pdf
* PDP11 assembler examples https://programmer209.wordpress.com/2011/08/03/example-pdp-11-programs/
* PDP11 assembler simulator (in JavaScript) https://programmer209.wordpress.com/2011/08/14/pdp-11-assembly-language-simulator/

* A guy who wrotes an own OS for PDP11 called MUXX, just for fun, with big bunch of knowlwdge http://ancientbits.blogspot.com/2012/09/writing-kernel-im-falling-in-love-with.html, https://ancientbits.blogspot.com/ ,
  also on Github https://github.com/jguillaumes

* PDP-11_Student_Workbook http://ancientbits.blogspot.com/2012/09/writing-kernel-im-falling-in-love-with.html
* GCC function attributes (e.g. interrupt service routine) https://gcc.gnu.org/onlinedocs/gcc/x86-Function-Attributes.html#x86-Function-Attributes
  and https://wiki.osdev.org/Interrupt_Service_Routines
* PDP11 Paper Tape Software Handbook -  http://www.bitsavers.org/www.computer.museum.uq.edu.au/pdf/DEC-11-XPTSA-B-D%20PDP-11%20Paper%20Tape%20Software%20Handbook.pdf
* https://www.learningpdp11.com/
* GNU GCC 9.2.0 for PDP-11 https://github.com/JamesHagerman/gcc-pdp11-aout

