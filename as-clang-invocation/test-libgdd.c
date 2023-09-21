/**
 * Example that uses an external function addasm()
 * this function is provided as assembler code. See file prov_fun.s
 *
 * The assembler function addasm() itself calls then a C function, named printLine(). This function is also
 * implemented below
 */

#include "console.h"

#define BTOA_MAXDIGITS 6 /* 5+VZ */

extern unsigned short udivmodhi4 (unsigned short num, unsigned short den, short modwanted);

/*
unsigned short udivmodhi4 (unsigned short num, unsigned short den, short modwanted) {
  unsigned short bit = 1;
  unsigned short res = 0;

  while (den < num && bit && !(den & (1 << 15))) {
      den <<= 1;
      bit <<= 1;
  }
  while (bit) {
      if (num >= den) {
          num -= den;
          res |= bit;
      }
      bit >>= 1;
      den >>= 1;
  }

  if (modwanted)
    return num;
  return res;
}
*/

short __divhi3 (short a, short b) {
  short neg = 0;
  short res;

  if (a < 0) {
      a = -a;
      //neg = !neg;
      neg=1;
  }

  if (b < 0) {
      b = -b;
      //neg = !neg;
      if (neg!=0)
        neg=0;
      else
        neg=1;
  }

  res = udivmodhi4 (a, b, 0);

  if (neg)
    res = -res;

  return res;
}

short __modhi3 (short a, short b) {
  short neg = 0;
  short res;

  if (a < 0) {
      a = -a;
      neg = 1;
  }

  if (b < 0)
    b = -b;

  res = udivmodhi4 (a, b, 1);

  if (neg)
    res = -res;

  return res;
}

short __udivhi3 (short a, short b) {
  return udivmodhi4 (a, b, 0);
}

short __umodhi3 (short a, short b) {
  return udivmodhi4 (a, b, 1);
}

void printInt16(int v) {
    int neg=0, zero=1;
    if (v<0) {
        neg = 1;
        v = -v;
    }
    if (v==0)
       zero=1;
    char s[BTOA_MAXDIGITS+2]; // with CRLF
    int remainder,si;
    si=BTOA_MAXDIGITS-1;
    while (v > 0) {
        remainder = v % 10;
        s[si--] = (char)remainder + '0';
        v = v / 10;
    }
    if (zero) {
        s[si--] = '0';
    }
    if (neg) {
        s[si--] = '-';
    }
    si++;
    s[BTOA_MAXDIGITS] = '\r';
    s[BTOA_MAXDIGITS+1] = '\n';
    cons_puts( &(s[si]) );
}


int main() {
    // int ret = _divhi3(0123,0456);

    //volatile int a = 100;
    //volatile int b= 330;
    //volatile int res = a * b;
    // volatile int res1 = b / a;
    //volatile int res = b % a;
    printInt16(12456);
    //printInt(res);
    //printLine((char*)res);
	//int result = addasm(21,31);
}

