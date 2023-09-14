/**
 * Example that uses an external function addasm()
 * this function is provided as assembler code. See file prov_fun.s
 */

extern int addasm(int a, int b);

int main() {
	int result = addasm(21,31);
}
