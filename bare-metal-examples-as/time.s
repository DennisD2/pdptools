        .title time
        .ident "v00.00"

        .global start

#        stack = 0x1000

        kbstat=0177560
        kbdata=0177562
        prstat=0177564
        prdata=0177566
        clstat=0177546

        .text
/*
        .org 0x4
            .word    6,0,012,0           # initialize error vectors
        .org 060
            .word    kbint,0340          # initialize keyboard int. vec. (prior. 7)
        .org 0100
            .word    clint,0300          # initialize clock int. vec. (prior. 6)
*/
/*
        lc=.
        .=4+lc
                .word    6,0,012,0           # initialize error vectors
        .=060+lc
                .word    kbint,0340          # initialize keyboard int. vec. (prior. 7)
        .=0100+lc
                .word    clint,0300          # initialize clock int. vec. (prior. 6)
        .=0500+lc                            # allow for stack space
*/

start:
vecs:   mov     $6,*$4              # initialize error vectors
        mov     $0,*$6
        mov     $012,*$010
        mov     $0,*$012

# print query
        mov     pc,sp
        tst     -(sp)               # init sp to start
        mov     $query,r0           # set parameters
        mov     $endq,r1            #       for print subroutine
        jsr     pc,print            # print lf, cr, query text

# accept and echo initial time xxyy
         mov     $4,r2              # (r2)=digit count
         mov     $itime,r0          # set parameters
nextd:  mov     r0,r1               #       for print subroutine
         tstb    *$kbstat           # character entered?
         bpl     .-4                # if not, keep testing
         movb    *$kbdata,(r0)      # else, store digit in itime array
         bicb    $0200,(r0)         # remove check bit from digit
         jsr     pc,print           # print digit
         inc     r0                 # (r0)=(r0)+1
         dec     r2                 # (r2)=(r2)-1
         bne     nextd              # if (r2) not 0, accept next digit

# convert initial hours (xx) to binary
        movb    *$itime+1,r0        # set parameters
        movb    *$itime,r1          #       for incon subroutine
        jsr     pc,incon            # convert xx to binary
        mov     r2,*$hour           #     and store in hour
# convert initial minutes (yy) to binary
        movb    *$itime+3,r0        # set parameters
        movb    *$itime+2,r1        #       for incon subroutine
        jsr     pc,incon            # convert xx to binary
        mov     r2,*$min            #     and store in min

isrs:
        mov     $kbint,*$060        # initialize keyboard int. vec. (prior. 7)
        mov     $kbint+2,*$0340     #
        mov     $clint,*$0100       # initialize clock int. vec. (prior. 6)
        mov     $clint+2,*$0300     #

# set interrupt enable bits to 1 and wait
        mov     $0100,*$kbstat      # set keyboard intr. enble bit to 1
        mov     $0100,*$clstat      # set clock intr. enble bit to 1
loop:   br      loop                # wait for interrupts

# clock interrupt handler
# updates time every 1/60 second
clint:  mov     $tick,r4            # set parameter for update s.r.
        jsr     pc,update           # update clock count
        cmp     *$hour,$12          # is (hour)=12., or less?
        ble     exit3               # if so, time update is complete
        sub     $12,*$hour          # else, correct for 12-hour clock
exit3:  rti                         # return from interrupt
# update (recursive subroutine)
# updates tick, sec, min and hour. address of updated field is in r4.
update: inc     (r4)                # ((r4))=((r4))+1
        cmp     (r4),$60            # ((r4))=60.?
        bne     exit4               # if not, updating is complete
        clr     (r4)                # else, ((r4))=0 (reset count)
        tst     -(r4)               # (r4)=(r4)-2 (go to next field)
        jsr     pc,update           # update next field
exit4:  rts     pc                  # exit

# keyboard interrupt handler
# prints out time whenever a character is typed in.
kbint:  mov     $temp,r0            # save latest
        mov     *$hour,(r0)+        #    hour, min and sec
        mov     *$min,(r0)+         #       in temp array to
        mov     *$sec,(r0)          #          protect from clint
        clr     $0177776            # lower priority to accept clint
# print message
        mov     $messg,r0           # set parameters
        mov     $endm,r1            #       for print subroutine
        jsr     pc,print            # print lf, cr, message text
# convert hour, min and sec to ascii
        mov     $temp,r2            # set parameters
        mov     $output,r3          #       for outcon subroutine
        jsr     pc,outcon           # convert hour to ascii (hh)
        jsr     pc,outcon           # convert min to ascii (mm)
        jsr     pc,outcon           # convert sec to ascii (ss)
# print out hh:mm:ss and ring bell
        mov     $output,r0          # set parameters
        mov     $endo,r1            #       for print subroutine
        jsr     pc,print            # print output array
        tst     *$kbdata            # clear ready bit in kbstst
        rti                         # return from interrupt

# print
print:  mov     r0,r5               # (r5)=character array index
again:  cmp     r5,r1               # has string ended?
        bhi     exit1               # if so, exit
        tstb    *$prstat            # is printer ready?
        bpl     .-4                 # if not, keep testing
        movb    (r5)+,*$prdata      # else, print (r5). (r5)=(r5)+1
        br      again               # pick up next character
exit1:  rts     pc                  # exit

# incon
# converts a 2-digit decimal number stored in ascii in r0 (units) and
# r1 (tens) into binary, the result is placed in r2. r3, r4, r5 unchanged.
incon:  bic     $0177760,r0         # convert (r0) into binary
        mov     r0,r2               #     and store in r2
tens:   cmpb    r1,$'0              # (r1)='0? (any tens left?)
        beq     exit2               # if not, exit
        add     $10,r2              # else, (r2)=(r2)+10 decimal
        dec     r1                  # (r1)=(r1)-1 (1 ten less)
        br      tens                # check for tens again
exit2:  rts     pc                  # exit

# outcon
# converts a binary number n (from 0 to 60 decimal) into a 2-digit
# ascii number pq. address of in is (r2). addresses of p and q are (r3)
# and (r3)+1. before exit the contents of r1 is incremented by 2 and of
# r3 by 3. r4 and r5 are unchanged.
outcon: mov     (r2)+,r0            # (r0)=binary number (hour, min, sec)
        clr     r1                  # initialize tens
more:   cmp     r0,$10              # any tens left in r0?
        blt     units               # if none, process units
        inc     r1                  # else, (r1)=(r1)+1 (one more ten)
        sub     $10,r0              # (r0)=(r0)-10 decimal
        br      more                # check for more tens
units:  add     $'0,r1              # convert tens to ascii
        add     $'0,r0              # convert units to ascii
        movb    r1,(r3)+            # store tens in output array
        movb    r0,(r3)+            # store units in output array
        inc     r3                  # skip colon byte
        rts     pc                  # exit

        .data

query:  .byte   15,12               # cr, lf
        .ascii  "What time is it?"  # query text
endq:   .ascii  " "                 # end of query (space)
        #
messg:  .byte   15,12               # cr, lf
        .ascii  "At the bell the time will be:" # message text
endm:   .ascii  " "                             # end of message (space)
#
output: .ascii  "hh:mm:ss"          # storage for hh:mm:ss
endo:   .byte   7                   # end of output (bell)
#
itime:  .space   4                  # storage for initial time (xxyy)
#
        .even                       # adjust word boundary
hour:   .space   2                  # storage for hours (binary)
min:    .space   2                  # storage for minutes (binary)
sec:    .word   0                   # storage for seconds (binary)
tick:   .word   0                   # storage for tick count (binary)
temp:   .space   3*2                # temp. storage for hour, min, sec (binary)

        .end
