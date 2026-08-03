#ident "@(#) TREK73 $Header: /home/ncvs/trek73/src/printships.c,v 1.2 1999/10/26 20:13:00 dillon Exp $"
/*
 * $Source: /home/ncvs/trek73/src/printships.c,v $
 *
 * $Header: /home/ncvs/trek73/src/printships.c,v 1.2 1999/10/26 20:13:00 dillon Exp $
 *
 * $Log: printships.c,v $
 * Revision 1.2  1999/10/26 20:13:00  dillon
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
 * Revision 1.1  87/10/09  11:10:37  11:10:37  okamoto (Jeff Okamoto)
 * Initial revision
 * 
 */
/*
 * TREK73: printships.c
 *
 * Print Basic Point Values for all standard ships.
 *
 */

#include "defs.h"

int
main(int ac, char **av)
{
	double bpv, efficiency;
	int turn, max_speed;
	int i;

	for(i=0; i<MAXSHIPCLASS; i++) {
		calculate(
		    stats[i].regen,
		    stats[i].pods,
		    stats[i].ph_shield,
		    stats[i].tp_shield,
		    stats[i].num_phaser + stats[i].num_torp,
		    stats[i].o_crew,
		    &bpv,
		    &efficiency,
		    &turn,
		    &max_speed
		);
		printf("%s:", stats[i].abbr);
		printf("\tBPV =\t\t%.2f\n", bpv);
		printf("\tEfficiency =\t%.2f\n", efficiency);
		printf("\tTurn =\t\t%d\n", turn);
		printf("\tMax speed =\t%d\n\n", max_speed);
	}
	return(0);
}

