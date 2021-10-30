#include "libcap.h"

static cap_value_t top;

static int cf(cap_value_t x) {
    return top - x - 1;
}

static int test_cap_bits(void) {
    static cap_value_t vs[] = {
	5, 6, 11, 12, 15, 16, 17, 38, 41, 63, 64, __CAP_MAXBITS+3, 0, -1
    };
    int failed = 0;
    cap_value_t i;
    for (i = 0; vs[i] >= 0; i++) {
	cap_value_t ans;

	top = i;
	_binary_search(ans, cf, 0, __CAP_MAXBITS, 0);
	if (ans != top) {
	    if (top > __CAP_MAXBITS && ans == __CAP_MAXBITS) {
	    } else {
		printf("test_cap_bits miscompared [%d] top=%d - got=%d\n",
		       i, top, ans);
		failed = -1;
	    }
	}
    }
    return failed;
}

int main(int argc, char **argv) {
    int result = 0;
    result = test_cap_bits() | result;
    if (result) {
	printf("test FAILED\n");
	exit(1);
    }
}
