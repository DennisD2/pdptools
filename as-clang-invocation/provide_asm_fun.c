/**
 * Example that uses an external function addasm()
 * this function is provided as assembler code. See file prov_fun.s
 *
 * The assembler function addasm() itself calls then a C function, named printLine(). This function is also
 * implemented below
 */

#include "console.h"

extern int addasm(int a, int b);
//extern int _divhi3( int a, int b);

/*
int add(int a, int b){
    return a+b;
}
*/
/*
void printInt(int v) {
    char s[8];
    int divider=10000;
    int remain,si;
    si=0;
    while (v>10) {
        remain = v % 10;
        s[si] = remain+'0';
        si++;
        //v = v/10;
    }
    s[si] = remain+'0';
    cons_puts(s);

}
*/

// this function will be called by assembler code
void printLine(char *message) {
    cons_puts(message);
}

int main() {
    // int ret = _divhi3(0123,0456);

    volatile int a = 012;
    volatile int b= 045;
    volatile int res = a * b;
    //volatile int res1 = a % b;
    //printInt(res);
    //printLine((char*)res);
	//int result = addasm(21,31);
}

/*
int __divhi3( int a, int b) {
    return a+b+0123;
}*/
