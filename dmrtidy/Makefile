SRC = main.c
OBJ = $(SRC:.c=.o)

.c.o: clean
	$(CC) -c $(CFLAGS) $<

all: dmr

dmr: dmr.o
	$(CC) -o dmrtidy dmr.o  $(LDFLAGS)

get:
	curl -s -o- https://ham-digital.org/user_by_call.php | pv | ./dmrtidy > contacts.csv

clean:
	rm -f dmr.o
