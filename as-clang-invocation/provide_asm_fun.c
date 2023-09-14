/**
 * Example that uses an external function addasm()
 * this function is provided as assembler code. See file prov_fun.s
 *
 * The assembler function addasm() itself calls then a C function, named printLine(). This function is also
 * implemented below
 */

#include "console.h"

extern int addasm(int a, int b);

int add(int a, int b){
    return a+b;
}

// this function will be called by assembler code
void printLine() {
    cons_puts("Hello World!\r\n");
}

int main() {
	int result = addasm(21,31);
}
