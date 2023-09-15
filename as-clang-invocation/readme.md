# Assembler calls to C and vice versa

## Assembler calls assembler
Assembler code calls an assembler function.

Compile with 
```
make asm_prov_and_use
```

Call code:
```
mov     $20,*-(sp)          
mov     $30,*-(sp)          
jsr     pc,_addasm          
sub     $4,sp
```

Function code:
```
_addasm:
        clr     r0
        add     *2(sp),r0
        add     *4(sp),r0
        rts     pc
```
**TODO: remove the * assignments to be compatible with C example!**

Invoked files: prov_fun_and_use_fun.s

## C calls assembler
Compile with
```
make clang_uses_asm_fun
```

Call code:
```
int result = addasm(21,31);

In assembler:
000024: 012746 000037       	mov	#37,-(r6)		; f...
000030: 012746 000025       	mov	#25,-(r6)		; f...
000034: 004767 177760       	call	20			; w.p.
000040: 062706 000004       	add	#4,r6			; Fe..
```
The octal values 037 (=3*8+7=31) and 025 (=2x8+5=21) are pushed to stack (r6 is sp register).
Stack is corrected after subroutine call by 4. 

Function code:
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

Invoked files: prov_fun.s, use_asm_fun.c, (and crt0.s, console.c)

## Assembler calls C
Compile with
```
make asm_uses_clang_fun
```

Call code:
```
        .extern _printLine      # clang function declared as extern
        
        ...
        
        jsr     pc,_printLine   # calling the clang function
```

Function code:
```
void printLine() {
    cons_puts("Hello World!\r\n");
}

Disassembled code for the C function:
000032: 012746 000060       	mov	#60,-(r6)		; f.0.
000036: 004767 177756       	call	20			; w.n.
000042: 062706 000002       	add	#2,r6			; Fe..
000046: 000207              	ret				; ..
```
In the disassembled code for the C function, it looks that #60 is the string address.

**TODO: For a better example, have a function with a single int parameter**

Invoked files: use_c_fun.s, provide_asm_fun.c, (and crt0.s , console.c)