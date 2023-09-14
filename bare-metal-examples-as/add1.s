        .title add
        .ident "v00.00"

        .global start

        kbstat=0177560
        kbdata=0177562
        prstat=0177564
        prdata=0177566
        clstat=0177546

        .text

        v_illadr=4
        v_illins=010

start:
vecs:   mov     $6,*$4              # initialize error vectors
        mov     $0,*$6
        mov     $012,*$010
        mov     $0,*$012

# main
        mov     pc,sp
        tst     -(sp)               # init sp to start

        mov     $20,r0           # set parameters
        mov     $30,r1            #       for print subroutine
        jsr     pc,addasm1            # print lf, cr, query text
        halt

# add adds r0+r1, result in r2
addasm1: add     r0,r2
        add     r1,r2
        rts     pc

        .end
