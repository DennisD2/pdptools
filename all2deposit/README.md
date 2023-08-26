# aout2deposit

Creates a SIMH deposit file from an ...

# todos
* make it work: 
  * find issue: what goes wrong; what is the jmp thing in crt0.s? does it harm? 
  * compile smallest possible thing like a=1;b=2c=a+b check code and if it runs
  * check if generated example runs on SIMH (ptap and my deposit file)
* create file (not stdout)
* handle vector0 arg
* find out how to do single step in ODT

Original:
```shell
aout2lda --aout hello --lda hello.ptap --data-align 2     --text 0x200 --vector0
      Magic    Text +    Len    Data +    Len     BSS = MaxMem   Entry
Hex:   0701    0200 +   0104    0304 +   0010    0000     0314    0200
Dec:   1793     512 +    260     772 +     16       0      788     512
Oct: 003401  001000 + 000404  001404 + 000020  000000   001424  001000
```

## Related
* Bare metal intro https://www.teckelworks.com/2020/03/c-programming-on-a-bare-metal-pdp-11/
* Cross compiling - https://xw.is/wiki/Bare_metal_PDP-11_GCC_9.3.0_cross_compiler_instructions
* baremetal topic - https://stackoverflow.com/questions/38387155/beginner-bare-metal-pdp11-console-output

* Discussion Load binary PDP files https://groups.google.com/g/pidp-11/c/AZkNx5FUpRY?pli=1
* Paper tape format info http://avitech.com.au/?page_id=709
* Old UNIX on PDP11 https://pdos.csail.mit.edu/6.828/2005/homework/hw5.html

* Modern PDP-11 C Compilers - https://retrocomputingforum.com/t/modern-pdp-11-c-compilers/2329/4
* Search PDP11 bare metal - https://www.google.com/search?q=pdp11+bare+metal+clang+compile&oq=pdp11+bare+metal+clang+compile&aqs=chrome..69i57j33i10i160l3.8962j0j4&sourceid=chrome&ie=UTF-8