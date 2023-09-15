        .title creates sum from two int values
        .ident "v00.00"

        .globl _addasm
        .extern _printLine      # clang function declared as extern
        .extern _add      # clang function declared as extern

        .text

# add adds two int values from stack, result in r0
_addasm:
        clr     r0
        add     2(sp),r0
        add     4(sp),r0
        jsr     pc,_printLine      # calling the clang printLine() function

        mov     $20,-(sp)          # push literal 20 dec. to stack
        mov     $30,-(sp)          # push literal 30 dec. to stack
        jsr     pc,_add            # calling the clang add() function
        add     $4,sp              # fix stack pointer, i.e. remove parameters 20 and 30 from stack

        rts     pc

        .end
