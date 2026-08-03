#
# Makefile for TREK73 documentation
#

#
# Select destination directory for manuals and documents
MANDIR = /usr/local/man6
CATDIR = /usr/local/cat6
DOCDIR = /usr/local/doc

TEXT = header 00.t 01.t 02.t 03.t 04.t 04.5.t 05.t 06.t 07.t 08dd.t\
       08cl.t 08ca.t 08dn.t 09dd.t 09cl.t 09ca.t 09dn.t 10.t
TREKDOC = trek73.doc

TREKSRC= trek73.6
TREKMAN = trek73.man
PRTSRC = printships.6
PRTMAN = printships.man
SHIPSRC = shipyard.6
SHIPMAN = shipyard.man

all: $(TREKDOC) $(TREKMAN) $(PRTMAN) $(SHIPMAN)

install: $(TREKDOC) $(TREKMAN) $(PRTMAN) $(SHIPMAN)
	cp $(TREKSRC) $(MANDIR)
	cp $(PRTSRC) $(MANDIR)
	cp $(SHIPSRC) $(MANDIR)
	cp $(TREKDOC) $(DOCDIR)
	cp $(TREKMAN) $(CATDIR)
	cp $(PRTMAN) $(CATDIR)
	cp $(SHIPMAN) $(CATDIR)

$(TREKDOC): $(TEXT)
	nroff -ms $(TEXT) > $@

$(TREKMAN): $(TREKSRC)
	nroff -man $(TREKSRC) > $@

$(PRTMAN): $(PRTSRC)
	nroff -man $(PRTSRC) > $@

$(SHIPMAN): $(SHIPSRC)
	nroff -man $(SHIPSRC) > $@

clean:
	\rm -f make.out $(TREKMAN) $(TREKDOC) $(PRTMAN) $(SHIPMAN)
