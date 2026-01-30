CC      = gcc
CFLAGS  = -O2 -Wall -Wextra
LDFLAGS = -lcurl

BIN     = service client

.PHONY: all clean

all: $(BIN)

service: service.c
	$(CC) $(CFLAGS) -o service service.c

client: client.c
	$(CC) $(CFLAGS) -o client client.c $(LDFLAGS)

clean:
	rm -f $(BIN)
