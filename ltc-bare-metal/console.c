#ifndef CONSOLE_H
#include "console.h"
#endif

#define DL11_RCSR	0177560
#define DL11_RCSR_DONE	0x80
#define DL11_RBUF	0177562
#define DL11_XCSR	0177564
#define DL11_XCSR_READY	0x80
#define DL11_XBUF	0177566

#define KBSTAT 177560
#define KBDATA 177562
#define PRSTAT 177564
#define PRDATA 177566
#define CLSTAT 177546

void cons_putc(char c) {
	volatile unsigned int *xcsr = (unsigned int *)DL11_XCSR;
	unsigned char *xbuf = (unsigned char *)DL11_XBUF;
	while (!(*xcsr & DL11_XCSR_READY)) ;
	*xbuf = c;
}

char cons_getc() {
	volatile unsigned int *rcsr = (unsigned int *)DL11_RCSR;
	unsigned char *rbuf = (unsigned char *)DL11_RBUF;
	while (! (*rcsr & DL11_RCSR_DONE)) ;
	return *rbuf & 0x7F;
}

void cons_gets(char *buffer, int size) {
	char c, *p = buffer;
	while (1) {
		c = cons_getc();
		if ((c == '\b') || (c == 0x7F)) {
		    // Backspace
			if (p > buffer) {
			    // go back in buffer
				p--;
				cons_putc('#');
			} else {
    			// cannot go further back,  Ring Bell
				cons_putc(7);
			}
		} else if (c >= ' ') {
			if (p < buffer + size - 2) {
			    // Add to buffer
				*(p++) = c;
				cons_putc(c);
			}
		} else if (c == '\r') {
		    // return -> line input end
			cons_putc(c);
			cons_putc('\n');
			return;
		}
		*p = 0;
	}
}

void cons_puts(char *s) {
	for (;*s;s++) cons_putc(*s);
}
