#ident "@(#) TREK73 $Header: /home/ncvs/trek73/src/ships.c,v 1.2 1999/10/26 20:13:01 dillon Exp $"
/*
 * $Source: /home/ncvs/trek73/src/ships.c,v $
 *
 * $Header: /home/ncvs/trek73/src/ships.c,v 1.2 1999/10/26 20:13:01 dillon Exp $
 *
 * $Log: ships.c,v $
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
 * Revision 1.1  87/10/09  15:55:36  15:55:36  okamoto (Jeff Okamoto)
 * Initial revision
 * 
 */
/*
 * TREK73: ships.c
 *
 * Global Ship Defs for TREK73
 *
 */

#include "defs.h"

Prototype struct ship_stat stats[];

/*
 * starting stats for the ships depending on the type
 */

struct ship_stat stats[] = {
{ "DN", 0, 6, 8, 8,  10, 2.,   1.5,  15., 200, 300, 600, 450, 5., 4.,  2, 4,
   135, 225, 125, 235, 4, 4},
{ "CA", 1, 4, 6, 9,  11, 1.,   0.75, 10., 150, 200, 450, 350, 3., 2.,  5, 2,
   135, 225, 125, 235, 4, 4},
{ "CL", 2, 4, 4, 9,  11, 0.75, 0.5,  10., 125, 175, 350, 250, 3., 2.,  6, 2,
   150, 210, 140, 220, 4, 4},
{ "DD", 3, 2, 4, 10, 12, 0.5,  0.5,  8.,  100, 150, 200, 150, 2., 1.5, 7, 1,
   160, 200, 150, 210, 4, 4}
};
