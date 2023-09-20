        .title creates sum from two int values
        .ident "v00.00"

        .globl _addasm
       # .globl ___divhi3
       # .globl ___mulhi3
       # .extern _printLine      # clang function declared as extern
       # .extern _add            # clang function declared as extern
       # .extern printInt

        .text

# add adds two int values from stack, result in r0
_addasm:clr     r0
        add     2(sp),r0
        add     4(sp),r0

        #mov     $messg,-(sp)
        #jsr     pc,_printLine      # calling the clang printLine() function
        #add     $2,sp

        #mov     $20,-(sp)          # push literal 20 dec. to stack
        #mov     $30,-(sp)          # push literal 30 dec. to stack
        #jsr     pc,_add            # calling the clang add() function
        #add     $4,sp              # fix stack pointer, i.e. remove parameters 20 and 30 from stack

        #mov     $04567,-(sp)
        #jsr     pc,_printInt
        rts     pc

# Multiplication, software version
# 2(sp)*4(sp) -> r2 -> r0; uses r0,r1,r2
___mulhi3_soft:
        mov     2(sp),r0
        mov     4(sp),r1
        clr     r2
mul1:   bit     $1,r1
        beq     mul2
        add     r0,r2
        bvs     mul3
mul2:   clc
        ror     r1
        asl     r0
        bvs     mul3
        tst     r1
        bne     mul1
mul3:   mov     r2,r0
        rts     pc

# Multiplication, hardware version
# 2(sp)*4(sp) -> r2 -> r0; uses r0,r1,r2
___mulhi3:
        mov     2(sp),r1
        mov     4(sp),r3
        mul     r1,r3
        mov     r3,r0
        rts     pc

___divhi3:  mov r0,r0
        mov r0,r0
        mov r0,r0
        rts pc

        .data

messg:  .ascii  "Hello there!"     # message text
        .byte   15,12              # cr, lf
endm:   .byte  0                  # end of message (space)

        .end
