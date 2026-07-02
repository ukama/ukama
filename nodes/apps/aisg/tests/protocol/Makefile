# AISG protocol-only golden tests.
# These tests intentionally compile only pure protocol files and local stubs;
# they do not require the full Ukama platform/vendor tree.

ROOT ?= ../..
CC ?= gcc

CFLAGS += -std=c99 -Wall -Wextra -Werror -g
CFLAGS += -I./stubs
CFLAGS += -I$(ROOT)/ctrl/inc

PROTOCOL_SRCS = \
	$(ROOT)/ctrl/src/hdlc.c \
	$(ROOT)/ctrl/src/xid.c \
	$(ROOT)/ctrl/src/retap.c \
	$(ROOT)/ctrl/src/retap_ops.c

TARGET = test_protocol

.PHONY: all test clean

all: $(TARGET)

test: $(TARGET)
	./$(TARGET)

$(TARGET): test_protocol.c $(PROTOCOL_SRCS)
	$(CC) $(CFLAGS) -o $@ test_protocol.c $(PROTOCOL_SRCS)

clean:
	rm -f $(TARGET) *.o
