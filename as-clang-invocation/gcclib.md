# Adding missing gcclib functions

When adding a divion to my test C code, like
```
int remain = v / 10;
```
I got an error during linking:
```
:(.text+0x2c): undefined reference to `__divhi3'
```

Checking internet resources showed, that gcc creates function 
calls for divisions (and other math operations), when no
opcodes for that math operation are provided by CPU.

These intermediate layer of low level math functions are usually
part of libgcc.a, linked to all executables. 

For PDP11, at least the default CPU version, it seems that the 
compiler is not aware of the division opcodes present. So it
creates calls to some C functions. In my case a call to
_divhi3().

I noticed that another function, __modhi3 is missed for the line
```
int remain = v % 10;
```
But only when removing all division lines. This means that the
division function __modhi3 also provided the remain on 
execution.
```
:(.text+0x16): undefined reference to `__modhi3'
```

A line like ```res = a*c``` results in missing __mulhi3 .
```
:(.text+0x2c): undefined reference to `__mulhi3'
```

## Implementing _mulhi3(), _divhi3() and _modhi3()
It looks for my gcc compiler short=int=16 bit value.

PDP11: can cope with 16 bit values, but has opcodes that also cope with 32 bit values.

For some insight to gcc config for PDP11, see: https://github.com/gcc-mirror/gcc/tree/master/gcc/config/pdp11

Documentation on some required functions: http://www.mirbsd.org/htman/i386/manINFO/gccint.html#Libgcc
and https://gcc.gnu.org/onlinedocs/gccint/Integer-library-routines.html#Integer-library-routines

For some code see https://github.com/glitchub/arith64/blob/master/README (64 bit implementations for 32 bit compiler)


This code:
```cgo
int main() {
    volatile int a = 0123;
    volatile int b= 0456;
    volatile int res = a - b;
    ...
```

compiles to:
```
...
000046: 012716 000123       	mov	#123,(r6)		; local var #1/a on stack
000052: 012766 000456 000002	mov	#456,2(r6)		; local var #2/b on stack
000060: 011600              	mov	(r6),r0			; get value #1/a from stack
000062: 016601 000002       	mov	2(r6),r1		; get value #2/b from stack
000066: 160100              	sub	r1,r0			; "a-b", result in r0
000070: 010066 000004       	mov	r0,4(r6)		; push result to stack
000074: 005000              	clr	r0			    ; clean up r0
000076: 062706 000006       	add	#6,r6			; remove from stack: a,b,result
```

This code:
```cgo
int main() {
    volatile int a = 0123;
    volatile int b= 0456;
    volatile int res = a * b;
    ...
```

compiles to:
```
...
000046: 012716 000123       	mov	#123,(r6)		; local var #1/a on stack
000052: 012766 000456 000002	mov	#456,2(r6)		; local var #2/b on stack
000060: 011600              	mov	(r6),r0			; get value #1/a from stack
000062: 016601 000002       	mov	2(r6),r1		; get value #2/b from stack
000066: 010146              	mov	r1,-(r6)		; param a to stack
000070: 010046              	mov	r0,-(r6)		; param b to stack
000072: 004767 177722       	call	20			; call __mulhi3(a,b)

000076: 010066 000010       	mov	r0,10(r6)		; push result (in r0) to stack 
000102: 005000              	clr	r0			    ; clean up r0
000104: 062706 000012       	add	#12,r6			; remove from stack: 0-a,2-b,4-?,6-?,10-result,12-subr.ret.adr.
```

The line ```int res = a / b;``` results in same calling code assembly.
So div-function and multiplication function have same parameters and return value.

For multiplication: Result of multhi3(a,b) is in r0. 

If the result is used in later C code lines, the code after the subroutine call is
slightly different:
```
...
# code before call identical as in example above
000072: 004767 177722       	call	20			; call __mulhi3(a,b)

000076: 062706 000004       	add	#4,r6			; remove a,b from stack.
000102: 010066 000004       	mov	r0,4(r6)		; push result (in r0) to stack 
000106: 016600 000004       	mov	4(r6),r0		; read result from stack
000112: 010046              	mov	r0,-(r6)		; push result to stack, different position for call
000114: 004767 177700       	call	20			; do the call (which takes single int)
000120: 005000              	clr	r0			    ; clean up r0
000122: 062706 000010       	add	#10,r6			; Remove from stack: 0-?,2-?,4-result(local var),6-ubr.ret.adr.
```

So this analysis brings no changes. This means I can implement the multhi3() function taking
2 parameters, a and b. These are short/int parameters.
The result needs to be stored in r0.

Same goes for divhi3().

For the modulo function modhi3(), code is also the same.



