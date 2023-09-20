                .title libdd contains small selection from libgcc functions for PDP11 with -nostdlib compile
                .ident "v00.00"

               #.globl ___mulhi3
               #.globl ___divhi3
               #.globl ___modhi3
               .global _udivmodhi4

                .text


# Multiplication, software version
# 2(sp)*4(sp) -> r2 -> r0; uses r0,r1,r2
___mulhi3_soft: mov     2(sp),r0
                mov     4(sp),r1
                clr     r2
mul1:           bit     $1,r1
                beq     mul2
                add     r0,r2
                bvs     mul3
mul2:           clc
                ror     r1
                asl     r0
                bvs     mul3
                tst     r1
                bne     mul1
mul3:           mov     r2,r0
                rts     pc


# Multiplication, hardware version
# a*b = 2(sp)*4(sp) -> r2 -> r0; uses r0,r1,r2
___mulhi3:      mov     2(sp),r1
                mov     4(sp),r3
                mul     r1,r3
                mov     r3,r0
                rts     pc

# Division, hardware version
# a/b = 4(sp)/2(sp) -> r2 -> r0 (remainder is in r3); uses r0,r1,r2,r3
# div r0,r2 = r2=r2.r3/r0; remainder goes to r3
# r2.r3 means HI-word.LO-word of 32 bit value in two registers
# registers must be even registers
___divhi3:      mov     r1,-(sp)
                mov     r2,-(sp)
                mov     r3,-(sp)
                mov     r4,-(sp)
                mov     r5,-(sp)
                mov     4(sp),r0
                mov     2(sp),r3
                clr     r2
                div     r0,r2
                mov     r2,r0
                mov     (sp)+,r5
                mov     (sp)+,r4
                mov     (sp)+,r3
                mov     (sp)+,r2
                mov     (sp)+,r1
                rts pc

# Modulo operation
# uses DIV opcode
___modhi3:      mov     r1,-(sp)
                mov     r2,-(sp)
                mov     r3,-(sp)
                mov     r4,-(sp)
                mov     r5,-(sp)
                mov     4(sp),r0
                mov     2(sp),r3
                clr     r2
                div     r0,r2
                mov     r3,r0   # remainder is in r3
                mov     (sp)+,r5
                mov     (sp)+,r4
                mov     (sp)+,r3
                mov     (sp)+,r2
                mov     (sp)+,r1
                rts pc


_udivmodhi4:    mov     r2,-(sp)
                mov     r3,-(sp)
                mov     r4,-(sp)
                mov     12(sp),r4 # modwanted
                mov     10(sp),r0 # a
                mov     8(sp),r3 # b
                clr     r2
                div     r0,r2
                tst     r4
                bne     modw
                mov     r2,r0   # res is in r2
                br      endfun
modw:           mov     r3,r0   # remainder is in r3
endfun:         mov     (sp)+,r4
                mov     (sp)+,r3
                mov     (sp)+,r2
                rts pc

                .data

                .end