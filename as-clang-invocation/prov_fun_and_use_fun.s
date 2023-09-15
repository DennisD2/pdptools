        .title add
        .ident "v00.00"

        .text

start:  mov     $6,*$4              # initialize error vectors
        mov     $0,*$6
        mov     $012,*$010
        mov     $0,*$012

# main
        mov     pc,sp
        tst     -(sp)               # init sp to start

        mov     $20,-(sp)          # push literal 20 dec. to stack
        mov     $30,-(sp)          # push literal 30 dec. to stack
        jsr     pc,_addasm         #
        add     $4,sp              # fix stack pointer, i.e. remove parameters 20 and 30 from stack
        halt

# add adds two int values from stack, result in r0
_addasm:
        clr     r0
        add     2(sp),r0
        add     4(sp),r0
        rts     pc

        .end
