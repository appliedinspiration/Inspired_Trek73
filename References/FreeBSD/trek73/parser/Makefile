CFLAGS = -g -DYYDEBUG
#CFLAGS = -g

all: main.o y.tab.o
	cc $(CFLAGS) main.o y.tab.o -ll -ly

y.tab.o: y.tab.h union.h
	cc -c $(CFLAGS) y.tab.c

y.tab.h: grammar.y lex.yy.c union.h
	yacc -d grammar.y

lex.yy.c: command.l union.h
	lex command.l

clean:
	/bin/rm -f a.out y.tab.o y.tab.h y.tab.c lex.yy.c main.o make.out y.output
