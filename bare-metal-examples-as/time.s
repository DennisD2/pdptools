     .TITLE Time
        .IDENT "V00.00"

        .GLOBAL start

        STACK = 0x1000

        KBSTAT=0177560
        KBDATA=0177562
        PRSTAT=0177564
        PRDATA=0177566
        CLSTAT=0177546

        LC=.
        .=4+LC
                .WORD    6,0,12,0           # INITIALIZE ERROR VECTORS
#        .=60+LC
#                .WORD    KBINT,340          # INITIALIZE KEYBOARD INT. VEC. (PRIOR. 7)
#        .=100+LC
#                .WORD    CLINT,300          # INITIALIZE CLOCK INT. VEC. (PRIOR. 6)
        .=500+LC                            # ALLOW FOR STACK SPACE

        .text
start:
# print query
        MOV     PC,SP
        TST     -(SP)               # INIT SP TO START
        MOV     $QUERY,R0           # SET PARAMETERS
        MOV     $ENDQ,R1            #       FOR PRINT SUBROUTINE
        JSR     PC,PRINT            # PRINT LF, CR, QUERY TEXT

 # ACCEPT AND ECHO INITIAL TIME XXYY
         MOV     $4,R2               # (R2)=DIGIT COUNT
         MOV     $ITIME,R0           # SET PARAMETERS
 NEXTD:  MOV     R0,R1               #       FOR PRINT SUBROUTINE
         TSTB    *$KBSTAT            # CHARACTER ENTERED?
         BPL     .-4                 # IF NOT, KEEP TESTING
         MOVB    *$KBDATA,(R0)       # ELSE, STORE DIGIT IN ITIME ARRAY
         BICB    $0200,(R0)          # REMOVE CHECK BIT FROM DIGIT
         JSR     PC,PRINT            # PRINT DIGIT
         INC     R0                  # (R0)=(R0)+1
         DEC     R2                  # (R2)=(R2)-1
         BNE     NEXTD               # IF (R2) NOT 0, ACCEPT NEXT DIGIT

# print subroutine
PRINT:  MOV     R0,R5               # (R5)=CHARACTER ARRAY INDEX
AGAIN:  CMP     R5,R1               # HAS STRING ENDED?
        BHI     EXIT1               # IF SO, EXIT
        TSTB    *$PRSTAT              # IS PRINTER READY?
        BPL     .-4                 # IF NOT, KEEP TESTING
        MOVB    (R5)+,*$PRDATA        # ELSE, PRINT (R5). (R5)=(R5)+1
        BR      AGAIN               # PICK UP NEXT CHARACTER
EXIT1:  RTS     PC                  # EXIT


99$:    nop
        halt

        .data
QUERY:  .BYTE   15,12               # CR, LF
        .ASCII  "WHAT TIME IS IT?"  # QUERY TEXT
ENDQ:   .ASCII  " "                 # END OF QUERY (SPACE)
        #
MESSG:  .BYTE   15,12               # CR, LF
        .ASCII  "AT THE BELL THE TIME WILL BE:" # MESSAGE TEXT
ENDM:   .ASCII  " "                             # END OF MESSAGE (SPACE)
#
OUTPUT: .ASCII  "HH:MM:SS"          # STORAGE FOR HH:MM:SS
ENDO:   .BYTE   7                   # END OF OUTPUT (BELL)
#
ITIME:  .space   4                   # STORAGE FOR INITIAL TIME (XXYY)
#
        .EVEN                       # ADJUST WORD BOUNDARY
HOUR:   .space   2                   # STORAGE FOR HOURS (BINARY)
MIN:    .space   2                   # STORAGE FOR MINUTES (BINARY)
SEC:    .WORD   0                   # STORAGE FOR SECONDS (BINARY)
TICK:   .WORD   0                   # STORAGE FOR TICK COUNT (BINARY)
TEMP:   .space   3*2                   # TEMP. STORAGE FOR HOUR, MIN, SEC (BINARY)
        .end
