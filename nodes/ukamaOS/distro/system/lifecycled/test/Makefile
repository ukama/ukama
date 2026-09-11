CC ?= cc
CFLAGS += -std=c11
CFLAGS += -D_GNU_SOURCE
CFLAGS += -g
CFLAGS += -Wall
CFLAGS += -Wextra
CFLAGS += -Werror
CFLAGS += -pthread
CFLAGS += -I../inc
CFLAGS += $(TEST_CPPFLAGS)
LDFLAGS += $(TEST_LDFLAGS)

.PHONY: all unit syntax component component-build run clean
all: unit

unit: test_lifecycle_fsm
	./test_lifecycle_fsm

test_lifecycle_fsm: test_fsm.c ../src/fsm.c ../src/state_store.c ../inc/fsm.h ../inc/state_store.h
	$(CC) $(CFLAGS) test_fsm.c ../src/fsm.c ../src/state_store.c $(LDFLAGS) -ljansson -o $@

syntax:
	$(CC) $(CFLAGS) -Iplatform -fsyntax-only ../src/*.c

component-build:
	$(CC) $(CFLAGS) -Iplatform ../src/*.c platform/platform.c $(LDFLAGS) \
		-lulfius -lorcania -ljansson -lcurl -o lifecycle_component

component:
	@test -n "$(LIFECYCLED_BIN)" || (echo "LIFECYCLED_BIN is required"; exit 1)
	LIFECYCLED_BIN=$(LIFECYCLED_BIN) python3 component_test.py

clean:
	rm -f test_lifecycle_fsm lifecycle_component

run: unit syntax component
