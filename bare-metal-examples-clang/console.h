#ifndef CONSOLE_H
#define CONSOLE_H

void cons_putc(char c);
char cons_getc();
void cons_gets(char *buffer, int size);
void cons_puts(char *s);

#endif
