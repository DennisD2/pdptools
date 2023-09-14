        .title creates sum from two int values
        .ident "v00.00"

        .globl _addasm
        .extern _printLine      # clang function declared as extern

        .text

# add adds two int values from stack, result in r0
_addasm:
        clr     r0
        add     2(sp),r0
        add     4(sp),r0
        jsr     pc,_printLine   # calling the clang function
        rts     pc

        .end
