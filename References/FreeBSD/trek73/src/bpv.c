#ident "@(#) TREK73 $Header: /home/ncvs/trek73/src/bpv.c,v 1.2 1999/10/26 20:12:57 dillon Exp $"
/*
 * $Source: /home/ncvs/trek73/src/bpv.c,v $
 *
 * $Header: /home/ncvs/trek73/src/bpv.c,v 1.2 1999/10/26 20:12:57 dillon Exp $
 *
 * $Log: bpv.c,v $
 * Revision 1.2  1999/10/26 20:12:57  dillon
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
 * Revision 1.2  87/10/09  15:48:35  15:48:35  okamoto (Jeff Okamoto)
 * Added declaration of round as a function returning a double.
 * 
 * Revision 1.1  87/10/09  11:00:29  11:00:29  okamoto (Jeff Okamoto)
 * Initial revision
 * 
 */
/*
 * TREK73: bpv.c
 *
 * Calculate the Basic Point Value of a ships
 *
 */

#include "defs.h"

Prototype void calculate(double regen, float pods, float p_div, float t_div, int weapons, int crew, double *bpv, double *eff, int *turn, int *max);

double myround(double x);

void
calculate(double regen, float pods, float p_div, float t_div, int weapons, int crew, double *bpv, double *eff, int *turn, int *max)
{
	*bpv = 0.;
	*bpv += regen * 12;
	*bpv += pods / 2;
	*bpv += p_div * 30;
	*bpv += t_div * 40;
	*bpv += weapons * 10;
	*bpv += crew / 15;

	*eff = myround(4 * (0.0034 * *bpv - 0.78)) / 4.0;
	if (*eff< 0.25)
		*eff= 0.25;
	*turn = (int) (10 - floor(*bpv / 100.0));
	if (*turn < 1)
		*turn = 1;
	*max= (int) myround(-0.004 * *bpv + 11.0);
	if (*max < 1)
		*max = 1;
}

double
myround(double x)
{
	return(floor(x + 0.5));
}

