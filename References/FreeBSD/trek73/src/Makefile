#
# Makefile for TREK73
#

#
# Select which operating system you are using.
# Acceptable flags are BSD and SYSV.
OS =	-DBSD
#OS =	-DSYSV

#
# Select whether you wish to use symbolic debugger or not.
# NOTE: UNIX System V Release 2.0 cannot do profiling on programs
# compiled with -g.  Also, sdb will not be very useful if the
# symbols are stripped from the load module. (See STRIP)
SDB =	-g
#SDB =

#
# Select whether code optimization is to be done.
OPT =	-O2 -Wall -Wstrict-prototypes
#OPT =

#
# Select whether profiling is to be done.
# NOTE: In System V Relase 2.0, this is incompatible with
# both SDB and STRIP.
# PROF =	-p
PROF =

#
# Select whether or not the load module is to be stripped.
# This is incompatible with both SDB and PROF.
# Note: In BSD Unix, this option should always be blank
# STRIP =	-s
STRIP =

#
# Select the method by which the binaries are to be installed.
# Install is better if you have it.
INSTALL = install
#INSTALL = cp

#
# Select whether or not tracing mode is to be turned on.
# This is useful when testing new strategies.
#TRACE =	-DTRACE

#
# Select whether or not you want the command line parser.
# This may slow the game down, since lex/yacc is inefficient.
PARSER = -DPARSER

#
# Select the destination directory for the game.
DESTDIR = /usr/games

#
# Options that you should think about setting
#
# SHOWTORP: When set, the player decides whether or not photon
#     torpedos should show up on a position display (command 13)
#     If not set, torpedos will automatically be displayed.
#
# PARANOID: If you want ALL the checking of a save file that rogue
#     provides, set this.  If you want to be cool, don't set it.
#
# NOTDEF: This option only makes sense if you have also set PARANOID.
#     This option defeats making a link to a save file and then
#     restarting the game from the link.
#
# HISTORICAL: If this option is set, then any ship which takes more
#     than 43 points of actual damage would explode.  This led to
#     much abuse of high-pod probes.  For more details, cf damage.c
#
OPTIONS =

PROTOS= protos.h

DEFINES = $(OS) $(PARSER) $(TRACE)
CFLAGS = $(DEFINES) $(SDB) $(OPT) $(PROF) $(STRIP) $(OPTIONS)

T73OBJECTS = cmds1.o cmds2.o cmds3.o cmds4.o damage.o dist.o endgame.o\
	enemycom.o firing.o globals.o init.o main.o misc.o mission.o\
	moveships.o parseopts.o save.o ships.o special.o strat1.o\
	subs.o vers.o lib.o

T73CFILES = cmds1.c cmds2.c cmds3.c cmds4.c damage.c dist.c endgame.c\
	enemycom.c firing.c globals.c init.c main.c misc.c mission.c\
	moveships.c parseopts.c save.c ships.c special.c strat1.c\
	subs.c vers.c lib.c

PARSEFILES = command.l grammar.y

PRTOBJECTS = bpv.o ships.o printships.o

PRTFILES = bpv.c ships.c printships.c

DYOOBJECTS = shipyard.o bpv.o lib.o

DYOFILES = shipyard.c bpv.c lib.c

HEADS= defs.h defines.h structs.h union.h

#
# If you are not using the parser, you do not need the lex and yacc libs
#
# Note: Some systems have all the lex and yacc modules in either libl
# or liby.  If your linker complains, try removing one or the other
LIBS= -lm -ll

FLUFF = parsit.o lex.yy.c y.tab.c grammar.c grammar.o\
	make.out errs core lint.errs a.out tags

BINS = trek73 printships shipyard

all:	$(PROTOS) trek73 printships shipyard

$(PROTOS): $(T73CFILES) parsit.c bpv.c
	cat $(T73CFILES) parsit.c bpv.c | egrep '^Prototype' > protos.h.new
	mv -f protos.h.new protos.h

trek73:	$(T73OBJECTS) grammar.o parsit.o
	cc $(CFLAGS) -o $@ $(T73OBJECTS) grammar.o parsit.o $(LIBS)

printships:	$(PRTOBJECTS)
	cc $(CFLAGS) -o $@ $(PRTOBJECTS) -lm

shipyard:	$(DYOOBJECTS)
	cc $(CFLAGS) -o $@ $(DYOOBJECTS) -lm

install: trek73 printships shipyard
	$(INSTALL) trek73 $(DESTDIR)
	$(INSTALL) printships $(DESTDIR)
	$(INSTALL) shipyard $(DESTDIR)

$(T73OBJECTS): $(HEADS)

$(PRTOBJECTS): $(HEADS)

$(DYOOBJECTS): $(HEADS)

grammar.o: grammar.y union.h command.l
	lex command.l
	yacc grammar.y
	mv y.tab.c grammar.c
	cc $(CFLAGS) -c grammar.c

parsit.o: parsit.c
	cc $(CFLAGS) -c parsit.c

tags: $(T73CFILES)
	ctags $(T73CFILES) $(HEADS)

count:
	wc $(T73CFILES) $(HEADS)

#
# Flags for your lint may need to be changed
#
lint:
	lint -bx $(DEFINES) $(T73CFILES) > lint.errs

clean:
	\rm -f $(BINS) $(T73OBJECTS) $(DYOOBJECTS) $(PRTOBJECTS) $(FLUFF)
	\rm -f $(PROTOS) $(PROTOS).new

