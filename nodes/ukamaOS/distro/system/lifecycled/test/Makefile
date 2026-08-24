# lifecycle.d tests.
#
# Unit tests are self-contained and require only libc and pthreads. The
# component test requires a production lifecycle.d binary.

CC ?= gcc

UNIT_TARGET = test_lifecycle_fsm
UNIT_SRC = test_fsm.c ../src/fsm.c ../src/state_store.c

CFLAGS += -std=c11
CFLAGS += -g
CFLAGS += -O0
CFLAGS += -Wall
CFLAGS += -Wextra
CFLAGS += -Werror
CFLAGS += -D_GNU_SOURCE
CFLAGS += -I../inc

.PHONY: all unit syntax component run clean

all: unit

$(UNIT_TARGET): $(UNIT_SRC)
	$(CC) $(CFLAGS) -o $@ $(UNIT_SRC)

unit: $(UNIT_TARGET)
	./$(UNIT_TARGET)

syntax:
	$(CC) $(CFLAGS) -I./stubs -fsyntax-only ../src/*.c

component:
	@if [ -z "$(LIFECYCLED_BIN)" ]; then \
		echo "LIFECYCLED_BIN is required"; \
		exit 1; \
	fi
	LIFECYCLED_BIN=$(LIFECYCLED_BIN) python3 component_test.py

run: unit syntax component

clean:
	rm -f $(UNIT_TARGET)
