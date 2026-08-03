#ident "@(#) TREK73 $Header: /home/ncvs/trek73/src/main.c,v 1.3 1999/10/26 20:13:00 dillon Exp $"
/*
 * $Source: /home/ncvs/trek73/src/main.c,v $
 *
 * $Header: /home/ncvs/trek73/src/main.c,v 1.3 1999/10/26 20:13:00 dillon Exp $
 *
 * $Log: main.c,v $
 * Revision 1.3  1999/10/26 20:13:00  dillon
 *
 *     Ansize the code.  Convert everything to prototypes
 *
 * Revision 1.2  1999/10/26 18:01:10  dillon
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
 * Revision 1.1  87/10/09  11:08:01  11:08:01  okamoto (Jeff Okamoto)
 * Initial revision
 * 
 */
/*
 * TREK73: main.c
 *
 * Originally written (in HP-2000 BASIC) by
 *	William K. Char, Perry Lee, and Dan Gee
 *
 * Rewritten in C by
 *	Dave Pare (sdcsvax!sdamos!mr-frog)
 * 		  (mr-frog@amos.ling.ucsd.edu)
 *		and
 *	Christopher Williams (ucbvax!ucbmerlin!williams)
 *			     (williams@merlin.berkeley.edu)
 *
 * Corrected, Completed, and Enhanced by
 *	Jeff Okamoto	(hpccc!okamoto)
 *			(okamoto%hpccc@hplabs.hp.com)
 *	Peter Yee	(ames!yee)
 *			(yee@ames.arc.nasa.gov)
 *	Matt Dillon	(ucbvax!dillon)
 *			(dillon@ucbvax.berkeley.edu)
 *	Dave Sharnoff	(ucbvax!ucbcory!muir)
 *			(muir@cogsci.berkeley.edu)
 *	Joel Duisman
 *	    and
 *	Roger J. Noe    (hplabs!uniq!uniqvax!rjnoe)
 *
 * Main Loop
 *
 * main, alarmtrap, quitgame, buffering
 *
 */

#include "defs.h"

Prototype void playit(void);

void alarmtrap(int);
void quitgame(int);
int buffering(FILE *stream);

extern int yyparse(void);

static jmp_buf	jumpbuf;

int
main(int argc, char *argv[], char *envp[])
{
	if (buffering(stdout) < 0)
		perror("cannot fstat stdout");
	(void) signal(SIGALRM, alarmtrap);
	(void) signal(SIGINT, quitgame);
	srandom(time(0));
	time_delay = DEFAULT_TIME;
	set_save();
	options = getenv("TREK73OPTS");
	if (options != NULL)
		parse_opts(options);
	get_comlineopts(argc, argv);
	if (restart && savefile[0] != '\0') {
		restore(savefile, envp);	/* Will not return */
		exit(1);
	}
	name_crew();
	init_ships();
	(void) mission();
	(void) alert();
	playit();
	/*NOTREACHED*/
	return(0);
}

/*
 * Main loop
 */

void
playit(void)
{
	struct ship	*sp;
#ifndef PARSER
	char			buf1[30];
#endif PARSER
	struct cmd		*cp;
	int			loop;
	char			*ch;

	(void) setjmp(jumpbuf);
	sp = shiplist[0];
	if (!(is_dead(sp, S_DEAD))) {
next:
		for (loop = 0; loop < HIGHSHUTUP; loop++)
			shutup[loop] = 0;
		(void) fflush(stdin);
		printf("\n%s: Code [1-%d] ", captain, high_command);
		(void) fflush(stdout);
		Input[0] = '\0';
		(void) alarm((unsigned) time_delay);
#ifdef PARSER
		(void) Gets(Input, sizeof(Input));
		if (Input[0] != NULL) {
#else
		(void) Gets(buf1, sizeof(buf1));
		if (buf1[0] != NULL) {
#endif PARSER
			(void) alarm(0);
#ifdef PARSER
			Inptr = Input;
			parsed[0] = '\0';
			ch = Inptr;
			while (isspace(*ch))
				ch++;
			if (isalpha(*ch))
				yyparse();
			else
				strcpy(parsed, Input);
#endif PARSER
#ifdef PARSER
			cp = scancmd(parsed);
#else
			cp = scancmd(buf1);
#endif PARSER
			if (cp != NULL) {
				(*cp->routine)(sp);
				if (cp->turns == FREE)
					goto next;
			} else
				printf("\n%s: %s, I am unable to interpret your last utterance.\n", science, title);
		} else
			(void) alarm(0);
	}
	ch = ch;		/* LINT */
	alarmtrap(0);
	/* This point is never reached since alarmtrap() always concludes
	   with a longjmp() back to the setjmp() above the next: label */
	/*NOTREACHED*/
}

void
alarmtrap(int sig)
{
	int i;

	if (sig) {
		puts("\n** TIME **");
		(void) signal(sig, alarmtrap);
		fpurge(stdin);
	}
	for (i = 1; i <= shipnum; i++)
		(*shiplist[i]->strategy)(shiplist[i]);
	if (!(is_dead(shiplist[0], S_DEAD)))
		printf("\n");
	move_ships();
	check_targets();
	misc_timers();
	disposition();
	longjmp(jumpbuf, 1);
}

void
quitgame(int signo)
{
	char answer[20];
	unsigned timeleft;

	timeleft = alarm(0);
	(void) signal(SIGINT, SIG_IGN);
	printf("\n\nDo you really wish to stop now?  [yes or no] ");
	(void) Gets(answer, sizeof(answer));
	if(answer[0] == NULL || answer[0] == 'y' || answer[0] == 'Y')
		exit(0);
	(void) signal(SIGINT, quitgame);
	if(timeleft)
		(void) alarm((unsigned)timeleft);
	return;
}


/* buffering: Determine whether or not stream is to be buffered.  If
   it's a character-special device, any buffering in effect will remain.
   If it's not a character-special device, then stream will be
   unbuffered.  There are many ways to decide what to do here.  One
   would have been to make it unbuffered if and only if
   !isatty(fileno(stream)).  This is usually implemented as a single
   ioctl() system call which returns true if the ioctl() succeeds, false
   if it fails.  But there are ways it could fail and still be a tty.
   Then there's also examination of stream->_flag.  UNIX is supposed to
   make any stream attached to a terminal line-buffered and all others
   fully buffered by default.  But sometimes even when isatty()
   succeeds, stream->_flag indicates _IOFBF, not _IOLBF.  And even if it
   is determined that the stream should be line buffered, setvbuf(3S)
   doesn't work right (in UNIX 5.2) to make it _IOLBF.  So about the
   only choice is to do a straightforward fstat() and ascertain
   definitely to what the stream is attached.  Then go with old reliable
   setbuf(stream, NULL) to make it _IONBF.  The whole reason this is
   being done is because the user may be using a pipefitting program to
   collect a "transcript" of a session (e.g. tee(1)), or redirecting to
   a regular file and then keeping a tail(1) going forever to actually
   play the game.  This assures that the output will keep pace with the
   execution with no sacrifice in efficiency for normal execution. [RJN]
*/

int
buffering(FILE *stream)
{
	struct stat	st;
	if (fstat(fileno(stream), &st) < 0)
		return -1;
	if ((st.st_mode & S_IFMT) != S_IFCHR)
		setbuf(stream, NULL);
	return 0;
}
