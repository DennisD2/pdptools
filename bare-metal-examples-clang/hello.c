#include "console.h"


extern int addasm(int a, int b);

int add(int a, int b){
    return a+b;
}

void printLine() {
    cons_puts("Hello World!\r\n");
}

int main() {
	/*int i;
	for (i=0; i<10; i++) {
		cons_puts("Hello World!\r\n");
	}*/
	int result = addasm(21,31);

}
