# wimc.d tests.
# Unit tests are self-contained and do not require Ukama platform/vendor libs.

CC ?= gcc

UNIT_TARGET = test_wimc_core

CFLAGS += -g
CFLAGS += -Wall
CFLAGS += -Wextra
CFLAGS += -O0
CFLAGS += -D_REENTRANT
CFLAGS += -I../test/stubs
CFLAGS += -I..
CFLAGS += -I../inc
CFLAGS += -I../inc/common
CFLAGS += -I../inc/agent
CFLAGS += -I../common

UNIT_SRC  = unit/test_wimc_core.c
UNIT_SRC += ../src/db.c
UNIT_SRC += ../src/tasks.c
UNIT_SRC += ../src/jserdes.c
UNIT_SRC += ../common/package_cache.c
UNIT_SRC += ../common/utils.c
UNIT_SRC += ../common/err.c

UNIT_OBJ = $(UNIT_SRC:.c=.o)

LIBS += -ljansson
LIBS += -lsqlite3
LIBS += -luuid
LIBS += -lpthread

all: unit

unit: fixtures $(UNIT_TARGET)
	./$(UNIT_TARGET)

fixtures:
	./fixtures/build-fixtures.sh

$(UNIT_TARGET): $(UNIT_OBJ)
	$(CC) -o $@ $(UNIT_OBJ) $(LIBS)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

integration: fixtures
	./integration/test_wimc_without_real_agent.sh

run: unit integration

clean:
	rm -f $(UNIT_TARGET) $(UNIT_OBJ)
	rm -rf fixtures/tarballs

.PHONY: all unit fixtures integration run clean
