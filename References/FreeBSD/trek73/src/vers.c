#ident "@(#) TREK73 $Header: /home/ncvs/trek73/src/vers.c,v 1.3 1999/10/26 20:39:35 dillon Exp $"
/*
 * $Source: /home/ncvs/trek73/src/vers.c,v $
 *
 * $Header: /home/ncvs/trek73/src/vers.c,v 1.3 1999/10/26 20:39:35 dillon Exp $
 *
 * $Log: vers.c,v $
 * Revision 1.3  1999/10/26 20:39:35  dillon
 *
 *     Update documentation
 *
 * Revision 1.2  1999/10/26 20:13:01  dillon
 *
 *     Ansize the code.  Convert everything to prototypes
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
 * Revision 1.1  87/10/09  11:13:48  11:13:48  okamoto (Jeff Okamoto)
 * Initial revision
 * 
 */
/*
 * TREK73: vers.c
 *
 * version
 */

#include "defs.h"

Prototype int vers(struct ship *dummy);

int
vers(struct ship *dummy)
{
	printf("%s\n", version);
	printf(
	    "Originally written (in HP-2000 BASIC) by\n"
	    "     William K. Char, Perry Lee, and Dan Gee\n"
	    "\n"
	    "Rewritten in C by\n"
	    "     Dave Pare (sdcsvax!sdamos!mr-frog)\n"
	    "             and\n"
	    "     Christopher Williams (ucbvax!ucbmerlin!williams)\n"
	    "\n"
	    "Corrected, Completed, and Enhanced circa 1985 by\n"
	    "     Jeff Okamoto    (ucbvax!okamoto)\n"
	    "     Peter Yee       (ucbvax!yee)\n"
	    "     Matt Dillon     (ucbvax!dillon)\n"
	    "     Dave Sharnoff   (ucbvax!ucbcory!muir)\n"
	    "             and\n"
	    "     Joel Duisman    (ucbvax!duisman)\n"
	    "\n"
	    "Further enhanced circa 1987 by\n"
	    "     Jeff, Peter, and Roger Noe\n"
	    "\n"
	    "Modernization Port by Matt Dillon\n"
	    "\n"
	);
	return 1;
}

