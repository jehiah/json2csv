CFLAGS = -I. -I/usr/local/include -O2 -g
LIBS = -L. -L/usr/local/lib -ljson

all: json2csv

json2csv: json2csv.c 
	$(CC) $(CFLAGS) -o json2csv json2csv.c $(LIBS)

install:
	/usr/bin/install json2csv /usr/local/bin/

clean:
	rm -f json2csv
