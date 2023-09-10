#include "console.h"

#define CLSTAT 177546
#define VECT_CLOCK 100

volatile int ticks = 0;
volatile char s=0102;

/*
struct interrupt_frame {
    int ps; // old status
    int pc; // old counter
};

__attribute__((interrupt)) static void ltc_ir(struct interrupt_frame* frame) {
    ticks++;
    if (ticks >= 50) {
        ticks=0;
        s++;
    }
} */

void ltc_ir(void) {
    ticks++;
    if (ticks >= 50) {
        ticks=0;
        s++;
    }
}

int main() {

/*
int src = 'A';
int dst;

asm ("mov %0, %1\n\t"
    "add $3, %0"
    : "=r" (dst)
    : "r" (src)
);

cons_putc(dst);
cons_putc('\r');
cons_putc('\n');
*/

/*
asm volatile ("mov %0,$0340\n"
    :
    : "i" (ltc_ir)
);
*/

// Disable IR
//unsigned char *clockStat = (unsigned char *)CLSTAT;
//*clockStat = 0000;

/*
	int i;
	for (i=0; i<10; i++) {
		cons_puts("Hello World!\r\n");
	}
*/

// Enable IR
//*clockStat = 0100;

/*
unsigned int *clockVect = (unsigned int *)VECT_CLOCK;
*clockVect = ltc_ir;
clockVect++;
*clockVect = 300;
*/

int old=0;
#define MAXLOOP 60000
unsigned int loop=0;
while (loop<MAXLOOP) {
    if (ticks!=old) {
        //cons_putc(s);
        //cons_putc('\r');
        cons_putc('.');
        old=ticks;
    }
    loop++;
}

}
