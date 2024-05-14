#include <stdio.h>

int main() {
	char arr[5] ="12345";
	char na[5]="abc";
	printf("%.5s\n", arr);
	printf("%.*s, na %s \n", 5,arr, na);

	printf("%.3s\n", na);
	printf("%.5s, na %.3s\n", na, na);
	printf("%.*s, na %.3s\n", 3, na, na);
	printf("%.*s\n", 5, na);

	return 0;
}
