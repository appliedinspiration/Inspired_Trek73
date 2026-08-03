#ident "@(#) TREK73 $Header: /home/ncvs/trek73/src/dist.c,v 1.2 1999/10/26 20:12:58 dillon Exp $"
/*
 * $Source: /home/ncvs/trek73/src/dist.c,v $
 *
 * $Header: /home/ncvs/trek73/src/dist.c,v 1.2 1999/10/26 20:12:58 dillon Exp $
 *
 * $Log: dist.c,v $
 * Revision 1.2  1999/10/26 20:12:58  dillon
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
 * Revision 1.1  87/10/09  11:04:52  11:04:52  okamoto (Jeff Okamoto)
 * Initial revision
 * 
 */
/*
 * TREK73: dist.c
 *
 * Power distribution routines
 *
 * distribute
 *
 */

#include "defs.h"

Prototype void distribute(struct ship *sp);

void
distribute(struct ship *sp)
{
	int i;
	float fuel;
	int load;
	int effload;
	int drain;
	int loop;
	float shield;
	struct ship *fed;

	fed = shiplist[0];
	/*
	 * Granularity of 1 second as far as this loop is concerned
	 */
	for (loop = 0; loop < (int)timeperturn; loop++) {

		fuel = sp->energy + sp->regen;	/* Slightly unrealistic */
		/*
		 * Calculate negative phaser drains
		 */
		for (i=0; i<sp->num_phasers; i++) {
			load = sp->phasers[i].load;
			drain = sp->phasers[i].drain;
			if ((sp->phasers[i].status & P_DAMAGED)
			    || (drain >= 0) || (load <= 0))
				continue;
			/*
			 * Drain the lesser of either the current load if the
			 * load is less than the drain, or the drain value
			 */
			effload = max(load + drain, 0);
			fuel += load - effload;
			sp->phasers[i].load = effload;
		}
		/*
		 * Calculate shield drains
		 */
		shield = 0.0;
		for (i=0; i<SHIELDS; i++)
			shield += sp->shields[i].attemp_drain;
		drain = ceil((double) shield);
		/*
		 * If all attempted drains are zero, or we have no
		 * fuel, our shields are down!
		 */
		if ((shield * fuel == 0) && !shutup[SHIELDSF]
		    && sp == shiplist[0]) {
			printf("%s: %s, our shields are down!\n",engineer, title);
			shutup[SHIELDSF]++;
		}
		/*
		 * If there's not enough fuel to sustain the drains, then
		 * ration it out in proportion to the attempted drains and
		 * say that shields are fluctuating.
		 */
		if (drain <= fuel) {
			fuel -= drain;
			for (i=0; i<SHIELDS; i++)
				sp->shields[i].drain = sp->shields[i].attemp_drain;
		} else {
			if (!shutup[SHIELDSF] && sp == shiplist[0]) {
				printf("%s: %s, our shields are fluctuating!\n",
				    engineer, title);
				shutup[SHIELDSF]++;
			}
			for (i=0; i<SHIELDS; i++)
				sp->shields[i].drain =
				    sp->shields[i].attemp_drain *
				    fuel / drain;
			fuel = 0.;
		}
		/*
		 * Calculate cloaking device drains.  If there is
		 * in sufficient energy to run the device, then
		 * it is turned off completely
		 */
		if (cantsee(sp)) {
			if (fuel < sp->cloak_energy) {
				if (sp == shiplist[0]) {
					sp->cloaking = C_OFF;
					printf("%s:  %s, there's not enough energy to",
					    engineer, title);
					puts("    keep our cloaking device activated.");
				} else
					(void) e_cloak_off(sp, fed);
			} else
				fuel -= sp->cloak_energy;
		}
		/*
		 * Calculate positive phaser drains
		 */
		for (i=0; i<sp->num_phasers && fuel > 0; i++) {
			if (fuel <=0.)
				break;
			load = sp->phasers[i].load;
			drain = sp->phasers[i].drain;
			if ((sp->phasers[i].status & P_DAMAGED)
			    || load >= MAX_PHASER_CHARGE || drain <= 0)
				continue;
			/*
			 * Load phasers either enough to top them off, or
			 * the full drain
			 */
			effload = min(MAX_PHASER_CHARGE,
			    load + min(drain, fuel));
			fuel -= effload - load;
			sp->phasers[i].load = effload;
		}
		/*
		 * Now balance the level of energy with the numer of pods
		 */
		sp->energy = min(fuel, sp->pods);
	}
}

