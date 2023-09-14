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

Invoked files: prov_fun_and_use_fun.s

## C calls assembler
Compile with
```
make clang_uses_asm_fun
```

Call code:
```
int result = addasm(21,31);
```

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
```

Invoked files: use_c_fun.s, provide_asm_fun.c, (and crt0.s , console.c)