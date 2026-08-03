
/*
 * $Source: /home/ncvs/trek73/src/lib.c,v $
 *
 * $Header: /home/ncvs/trek73/src/lib.c,v 1.2 1999/10/26 20:13:00 dillon Exp $
 *
 * $Log: lib.c,v $
 * Revision 1.2  1999/10/26 20:13:00  dillon
 *
 *     Ansize the code.  Convert everything to prototypes
 *
 * Revision 1.1  1999/10/26 18:01:10  dillon
 *
 *     Initial changes to make trek compile under FreeBSD.  Remove some externs
 *     for library functions, fixup command.l to operate with a modern lexer.
 *
 * Revision 1.1.1.1  1999/10/26 17:07:14  dillon
 *
 * Nov 1988 import of trek73 from comp.sources.usenet.  This is the trek73
 * that was originally ported from basic to C and then further modified by
 * Jeff, Peter, Matt, and a few others in 1985.  Jeff later made additional
 * changes which led to the 1988 posting.
 *
 *
 *
 * Revision 1.1  87/10/09  11:08:09  11:08:09  okamoto (Jeff Okamoto)
 * Initial revision
 * 
 */
/*
 * TREK73: lib.c
 *
 * Miscellaneous Self Contained Routines
 */

#include "defs.h"

Prototype char *gets2(char *buf);

char *
gets2(char *buf)
{
	char *p;

	if ((p = fgets(buf, 20, stdin)) != NULL) {
		int l;

		if ((l = strlen(p)) > 0 && p[l-1] == '\n')
			p[l-1] = 0;
	}
	return(p);
}

