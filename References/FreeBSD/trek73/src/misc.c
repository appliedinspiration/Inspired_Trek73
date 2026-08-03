#ident "@(#) TREK73 $Header: /home/ncvs/trek73/src/misc.c,v 1.4 1999/10/26 20:18:30 dillon Exp $"
/*
 * $Source: /home/ncvs/trek73/src/misc.c,v $
 *
 * $Header: /home/ncvs/trek73/src/misc.c,v 1.4 1999/10/26 20:18:30 dillon Exp $
 *
 * $Log: misc.c,v $
 * Revision 1.4  1999/10/26 20:18:30  dillon
 *
 *     help function must return a value
 *
 * Revision 1.3  1999/10/26 20:13:00  dillon
 *
 *     Ansize the code.  Convert everything to prototypes
 *
 * Revision 1.2  1999/10/26 18:01:11  dillon
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
 * TREK73: misc.c
 *
 * Miscellaneous Routines
 *
 * help, scancmd, new_slot, return_slot, vowelstr,
 * check_p_damage, check_t_damage, check_p_turn, check_t_turn,
 * ship_class
 *
 */

#include "defs.h"

Prototype int help(struct ship *dummy);
Prototype struct cmd *scancmd(char *buf);
Prototype int new_slot(void);
Prototype void return_slot(int i);
Prototype char *vowelstr(char *str);
Prototype void check_p_damage(int array[], struct ship *sp, char *string);
Prototype void check_t_damage(int array[], struct ship *sp, char *string);
Prototype void check_p_turn(int array[], struct ship *sp, int flag);
Prototype void check_t_turn(int array[], struct ship *sp, int flag);
Prototype struct ship_stat *ship_class(char *s);;

int
help(struct ship *dummy)
{
	struct cmd	*cpbegin, *cpmiddle;

	puts("\nTrek73 Commands:");
	puts(
"Code		Command			Code		Command");
	putchar('\n');
	cpbegin = &cmds[0];
	cpmiddle = &cmds[(cmdarraysize + 1) >> 1];
	while (cpmiddle->routine != NULL) {
		printf("%3s:  %c %-31s:%3s:  %c %-31s\n",
		    cpbegin->code_num, ((cpbegin->turns) ? ' ' : '*'),
		    cpbegin->explanation, cpmiddle->code_num,
		    ((cpmiddle->turns) ? ' ' : '*'),
		    cpmiddle->explanation);
		cpbegin++;
		cpmiddle++;
	}
	if (cmdarraysize & 1)
		printf("%3s:  %c %-31s", cpbegin->code_num,
		    ((cpbegin->turns) ? ' ' : '*'),
		    cpbegin->explanation);
	puts("\n\n * does not use a turn");
	dummy = dummy;				/* LINT */
	return(0);
}

struct cmd *
scancmd(char *buf)
{
	static char **argp = NULL;
	struct	cmd *cp;
	int	argnum;
	int	first;

	argnum = parsit(buf, &argp);
	first = argp[0] ? strlen(argp[0]) : 0;
	if (argnum && first)
		for (cp = &cmds[0]; cp->routine != NULL; cp++)
			if (!strncmp(argp[0], cp->code_num, first))
				return (cp);
	return (NULL);
}

/*
 * This routine handles getting unique identifier numbers for
 * all objects.
 */
int
new_slot(void)
{
	/*
	 * This is to make it appear that in a 2-ship duel, for
	 * instance, the first object to appear will be numbered
	 * as 3.
	 */
	int i = shipnum + 2;

	while ((slots[i] == 'X') && (i <= HIGHSLOT))
		i++;
	if (i > HIGHSLOT) {
		puts("The game will terminate now due to an inability to handle the number of");
		puts("objects in space (i.e. vessels, torpedoes, probes, etc).  Sorry!");
		exit(-1);
	}
	slots[i] = 'X';
	return i;
}

/* 
 * This routine handles returning identifiers
 */

void
return_slot(int i)
{
	if (slots[i] != 'X')
		printf("FATAL ERROR - Slot already empty!");
	slots[i] = ' ';
}


char *
vowelstr(char *str)
{
	switch(*str) {
	case 'a': case 'A':
	case 'e': case 'E':
	case 'i': case 'I':
	case 'o': case 'O':
	case 'u': case 'U':
		return "n";
		/* not reached */
	}
	return "";
}


/*
 * This routine takes an array generated from commands 1, 3, and 5
 * to print out a list of those phasers damaged and unable to
 * either fire, lock, or turn.
 */

void
check_p_damage(int array[], struct ship *sp, char *string)
{
	int i, j;

	j = 0;
	for (i=0; i<sp->num_phasers; i++) {
		if (array[i] && (sp->phasers[i].status & P_DAMAGED)) {
			if (!j)
				printf("Computer: Phaser(s) %d", i+1);
			else
				printf(", %d", i+1);
			j++;
		}
	}
	if (j > 1)
		printf(" are damaged and unable to %s.\n", string);
	else if (j == 1)
		printf(" is damaged and unable to %s.\n", string);
}

/*
 * This routine takes an array generated from commands 2, 4, and 6
 * to print out a list of those tubes damaged and unable to either
 * fire, lock, or turn.
 */

void
check_t_damage(int array[], struct ship *sp, char *string)
{
	int i, j;

	j = 0;
	for (i=0; i<sp->num_tubes; i++) {
		if (array[i] && (sp->tubes[i].status & P_DAMAGED)) {
			if (!j)
				printf("Computer: Tube(s) %d", i+1);
			else
				printf(", %d", i+1);
			j++;
		}
	}
	if (j > 1)
		printf(" are damaged and unable to %s.\n", string);
	else if (j == 1)
		printf(" is damaged and unable to %s.\n", string);
}

/*
 * This routine checks to see if a phaser is pointing into our
 * blind side
 */

void
check_p_turn(int array[], struct ship *sp, int flag)
{
	int i;
	int j;
	float k;
	float bear;
	struct ship *target;

	j = 0;
	for (i=0; i<sp->num_phasers; i++) {
		if (!array[i])
			continue;
		if (flag && !(sp->phasers[i].status & P_FIRING))
			continue;
		target = sp->phasers[i].target;
		/*
		 * This hack is here since when the phaser is locked,
		 * the bearing points at the target, whereas when
		 * not locked, the bearing is relative to the ship.
		 */
		if (target == NULL) {
			bear = sp->phasers[i].bearing + sp->course;
			k = sp->phasers[i].bearing;
		} else {
			bear = bearing(sp->x,
			    (cantsee(target)) ? target->position.x : target->x,
			    sp->y,
			    (cantsee(target)) ? target->position.y : target->y);
			k = bear - sp->course;
		}
		k = rectify(k);

		if (betw(k, sp->p_blind_left, sp->p_blind_right)
		    && !(is_dead(sp, S_ENG))) {
			if (!j)
				printf("Computer: Phaser(s) %d", i + 1);
			else
				printf(", %d", i + 1);
			j++;
		}
	}
	if (j > 1)
		printf(" are pointing into our blind side.\n");
	else if (j == 1)
		printf(" is pointing into our blind side.\n");
}

/*
 * This routine checks to see if a tube is turned into
 * our blind side.
 */

void
check_t_turn(int array[], struct ship *sp, int flag)
{
	int i;
	int j;
	float k;
	float bear;
	struct ship *target;

	j = 0;
	for (i=0; i<sp->num_tubes; i++) {
		if (!array[i])
			continue;
		if (flag && !(sp->tubes[i].status & T_FIRING))
			continue;
		target = sp->tubes[i].target;
		/*
		 * This hack is here since when the tube is locked,
		 * the bearing points at the target, whereas when
		 * not locked, the bearing is relative to the ship.
		 */
		if (target == NULL) {
			bear = sp->tubes[i].bearing + sp->course;
			k = sp->tubes[i].bearing;
		} else {
			bear = bearing(sp->x,
			    (cantsee(target)) ? target->position.x : target->x,
			    sp->y,
			    (cantsee(target)) ? target->position.y : target->y);
			k = bear - sp->course;
		}
		k = rectify(k);
		if (betw(k, sp->t_blind_left, sp->t_blind_right) && !(is_dead(sp, S_ENG))) {
			if (!j)
				printf("Computer: Tubes(s) %d", i + 1);
			else
				printf(", %d", i + 1);
			j++;
		}
	}
	if (j > 1)
		printf(" are pointing into our blind side.\n");
	else if (j == 1)
		printf(" is pointing into our blind side.\n");
}

struct ship_stat *
ship_class(char *s)
{
	int i;

	for (i = 0; i< MAXSHIPCLASS; i++)
		if (!strcmp(stats[i].abbr, s)) {
			return(&stats[i]);
		}
	return(NULL);
}

