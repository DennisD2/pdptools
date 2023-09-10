	.text
    	.even
    	.globl	_main
    	.globl	___main
    	.globl	_start
   	    .globl	_ltc_isr2

#############################################################################*
##### _start: initialize stack pointer,
#####         clear vector memory area,
#####         save program entry in vector 0
#####         call C main() function
#############################################################################*
_start:
	mov	$00776,sp
	clr	r0
L_0:
	clr	(r0)+
	cmp	r0, $400
	bne	L_0
    mov	$000137,*$0     # Store JMP _start in vector 0
    mov	$_start,*$2

    mov $_ltc_isr2,*$0100 # LTC values
    mov $00300,*$0102

	jsr 	pc,_main
	halt
    br	_start

#############################################################################*
##### ___main: called by C main() function. Currently does nothing
#############################################################################*
___main:
	rts	pc


#############################################################################*
##### _ltc_isr2: call ltc_ir function in C code
#############################################################################*
_ltc_isr2:
	jsr 	pc,_ltc_ir
	rti
	halt