     .TITLE Time
        .IDENT "V00.00"

        .GLOBAL start

#        STACK = 0x1000

        KBSTAT=0177560
        KBDATA=0177562
        PRSTAT=0177564
        PRDATA=0177566
        CLSTAT=0177546

        .text
/*
        .org 0x4
            .WORD    6,0,012,0           # INITIALIZE ERROR VECTORS
        .org 060
            .WORD    KBINT,0340          # INITIALIZE KEYBOARD INT. VEC. (PRIOR. 7)
        .org 0100
            .WORD    CLINT,0300          # INITIALIZE CLOCK INT. VEC. (PRIOR. 6)
*/
/*
        LC=.
        .=4+LC
                .WORD    6,0,012,0           # INITIALIZE ERROR VECTORS
        .=060+LC
                .WORD    KBINT,0340          # INITIALIZE KEYBOARD INT. VEC. (PRIOR. 7)
        .=0100+LC
                .WORD    CLINT,0300          # INITIALIZE CLOCK INT. VEC. (PRIOR. 6)
        .=0500+LC                            # ALLOW FOR STACK SPACE
*/

start:
VECS:   mov     $6,*$4                # INITIALIZE ERROR VECTORS
        mov     $0,*$6
        mov     $012,*$010
        mov     $0,*$012

# PRINT QUERY
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

# CONVERT INITIAL HOURS (XX) TO BINARY
        MOVB    *$ITIME+1,R0          # SET PARAMETERS
        MOVB    *$ITIME,R1            #       FOR INCON SUBROUTINE
        JSR     PC,INCON            # CONVERT XX TO BINARY
        MOV     R2,*$HOUR             #     AND STORE IN HOUR
# CONVERT INITIAL MINUTES (YY) TO BINARY
        MOVB    *$ITIME+3,R0          # SET PARAMETERS
        MOVB    *$ITIME+2,R1          #       FOR INCON SUBROUTINE
        JSR     PC,INCON            # CONVERT XX TO BINARY
        MOV     R2,*$MIN              #     AND STORE IN MIN

ISRS:
        mov     $KBINT,*$060          # INITIALIZE KEYBOARD INT. VEC. (PRIOR. 7)
        mov     $KBINT+2,*$0340       #
        mov     $CLINT,*$0100          # INITIALIZE CLOCK INT. VEC. (PRIOR. 6)
        mov     $CLINT+2,*$0300       #

# SET INTERRUPT ENABLE BITS TO 1 AND WAIT
        MOV     $0100,*$KBSTAT         # SET KEYBOARD INTR. ENBLE BIT TO 1
        MOV     $0100,*$CLSTAT         # SET CLOCK INTR. ENBLE BIT TO 1
LOOP:   BR      LOOP                # WAIT FOR INTERRUPTS

# CLOCK INTERRUPT HANDLER
# UPDATES TIME EVERY 1/60 SECOND
CLINT:  MOV     $TICK,R4            # SET PARAMETER FOR UPDATE S.R.
        JSR     PC,UPDATE           # UPDATE CLOCK COUNT
        CMP     *$HOUR,$12           # IS (HOUR)=12., OR LESS?
        BLE     EXIT3               # IF SO, TIME UPDATE IS COMPLETE
        SUB     $12,*$HOUR           # ELSE, CORRECT FOR 12-HOUR CLOCK
EXIT3:  RTI                         # RETURN FROM INTERRUPT
# UPDATE (RECURSIVE SUBROUTINE)
# UPDATES TICK, SEC, MIN AND HOUR. ADDRESS OF UPDATED FIELD IS IN R4.
UPDATE: INC     (R4)                # ((R4))=((R4))+1
        CMP     (R4),$60           # ((R4))=60.?
        BNE     EXIT4               # IF NOT, UPDATING IS COMPLETE
        CLR     (R4)                # ELSE, ((R4))=0 (RESET COUNT)
        TST     -(R4)               # (R4)=(R4)-2 (GO TO NEXT FIELD)
        JSR     PC,UPDATE           # UPDATE NEXT FIELD
EXIT4:  RTS     PC                  # EXIT
        
# KEYBOARD INTERRUPT HANDLER
# PRINTS OUT TIME WHENEVER A CHARACTER IS TYPED IN.
KBINT:  MOV     $TEMP,R0            # SAVE LATEST
        MOV     *$HOUR,(R0)+          #    HOUR, MIN AND SEC
        MOV     *$MIN,(R0)+           #       IN TEMP ARRAY TO
        MOV     *$SEC,(R0)            #          PROTECT FROM CLINT
        CLR     $0177776             # LOWER PRIORITY TO ACCEPT CLINT
# PRINT MESSAGE
        MOV     $MESSG,R0           # SET PARAMETERS
        MOV     $ENDM,R1            #       FOR PRINT SUBROUTINE
        JSR     PC,PRINT            # PRINT LF, CR, MESSAGE TEXT
# CONVERT HOUR, MIN AND SEC TO ASCII
        MOV     $TEMP,R2            # SET PARAMETERS
        MOV     $OUTPUT,R3          #       FOR OUTCON SUBROUTINE
        JSR     PC,OUTCON           # CONVERT HOUR TO ASCII (HH)
        JSR     PC,OUTCON           # CONVERT MIN TO ASCII (MM)
        JSR     PC,OUTCON           # CONVERT SEC TO ASCII (SS)
# PRINT OUT HH:MM:SS AND RING BELL
        MOV     $OUTPUT,R0          # SET PARAMETERS
        MOV     $ENDO,R1            #       FOR PRINT SUBROUTINE
        JSR     PC,PRINT            # PRINT OUTPUT ARRAY
        TST     *$KBDATA              # CLEAR READY BIT IN KBSTST
        RTI                         # RETURN FROM INTERRUPT

# PRINT
PRINT:  MOV     R0,R5               # (R5)=CHARACTER ARRAY INDEX
AGAIN:  CMP     R5,R1               # HAS STRING ENDED?
        BHI     EXIT1               # IF SO, EXIT
        TSTB    *$PRSTAT              # IS PRINTER READY?
        BPL     .-4                 # IF NOT, KEEP TESTING
        MOVB    (R5)+,*$PRDATA        # ELSE, PRINT (R5). (R5)=(R5)+1
        BR      AGAIN               # PICK UP NEXT CHARACTER
EXIT1:  RTS     PC                  # EXIT

# INCON
# CONVERTS A 2-DIGIT DECIMAL NUMBER STORED IN ASCII IN R0 (UNITS) AND
# R1 (TENS) INTO BINARY, THE RESULT IS PLACED IN R2. R3, R4, R5 UNCHANGED.
INCON:  BIC     $0177760,R0          # CONVERT (R0) INTO BINARY
        MOV     R0,R2               #     AND STORE IN R2
TENS:   CMPB    R1,$'0              # (R1)='0? (ANY TENS LEFT?)
        BEQ     EXIT2               # IF NOT, EXIT
        ADD     $10,R2              # ELSE, (R2)=(R2)+10 DECIMAL
        DEC     R1                  # (R1)=(R1)-1 (1 TEN LESS)
        BR      TENS                # CHECK FOR TENS AGAIN
EXIT2:  RTS     PC                  # EXIT

# OUTCON
# CONVERTS A BINARY NUMBER N (FROM 0 TO 60 DECIMAL) INTO A 2-DIGIT
# ASCII NUMBER PQ. ADDRESS OF IN IS (R2). ADDRESSES OF P AND Q ARE (R3)
# AND (R3)+1. BEFORE EXIT THE CONTENTS OF R1 IS INCREMENTED BY 2 AND OF
# R3 BY 3. R4 AND R5 ARE UNCHANGED.
OUTCON: MOV     (R2)+,R0            # (R0)=BINARY NUMBER (HOUR, MIN, SEC)
        CLR     R1                  # INITIALIZE TENS
MORE:   CMP     R0,$10            # ANY TENS LEFT IN R0?
        BLT     UNITS               # IF NONE, PROCESS UNITS
        INC     R1                  # ELSE, (R1)=(R1)+1 (ONE MORE TEN)
        SUB     $10,R0             # (R0)=(R0)-10 DECIMAL
        BR      MORE                # CHECK FOR MORE TENS
UNITS:  ADD     $'0,R1              # CONVERT TENS TO ASCII
        ADD     $'0,R0              # CONVERT UNITS TO ASCII
        MOVB    R1,(R3)+            # STORE TENS IN OUTPUT ARRAY
        MOVB    R0,(R3)+            # STORE UNITS IN OUTPUT ARRAY
        INC     R3                  # SKIP COLON BYTE
        RTS     PC                  # EXIT

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
