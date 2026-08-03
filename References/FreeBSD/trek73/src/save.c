#ident "@(#) TREK73 $Header: /home/ncvs/trek73/src/save.c,v 1.3 1999/10/26 20:13:01 dillon Exp $"
/*
 * $Source: /home/ncvs/trek73/src/save.c,v $
 *
 * $Header: /home/ncvs/trek73/src/save.c,v 1.3 1999/10/26 20:13:01 dillon Exp $
 *
 * $Log: save.c,v $
 * Revision 1.3  1999/10/26 20:13:01  dillon
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
 * Revision 1.1  87/10/09  11:10:52  11:10:52  okamoto (Jeff Okamoto)
 * Initial revision
 * 
 */
/*
 * TREK73: save.c
 *
 * save and restore routines
 *
 * @(#)save.c	4.15 (Berkeley) 5/10/82
 */

#include "defs.h"

Prototype int save_game(struct ship *dummy);
Prototype void set_save(void);
Prototype void restore(char *file, char **envp);

void save_file(FILE *savef);

void encwrite(char *start, unsigned int size, FILE *outf);
int encread(char *start, unsigned int size, int inf);

#define MAXSTR	256

typedef struct stat STAT;

STAT sbuf;

void
set_save(void)
{
	char		*env;
	struct passwd	*pw;

	if ((env = getenv("HOME")) != NULL)
		strcpy(home, env);
	else if ((pw = (struct passwd *)getpwuid((int)getuid()))
	    != NULL)
		strcpy(home, pw->pw_dir);
	else
		home[0] = '\0';
	strcat(home, "/");

	strcpy(savefile, home);
	strcat(savefile, "trek73.save");
}

/*
 * save_game:
 *	Implement the "save game" command
 */

int
save_game(struct ship *dummy)
{
    FILE *savef;
    int c = 0;
    char buf[MAXSTR];

    /*
     * get file name
     */

over:
    if (savefile[0] != '\0')
    {
	for (;;)
	{
	    printf("Save file (%s)? ", savefile);
	    c = getchar();
	    if (c == 'n' || c == 'N' || c == 'y' || c == 'Y')
		break;
	    else
		printf("\nPlease answer Yes or No\n");
	}
	if (c == 'y' || c == 'Y')
	{
	    strcpy(buf, savefile);
	    goto gotfile;
	}
    }

    do
    {
	fpurge(stdin);
	printf("File name: ");
	buf[0] = '\0';
	gets2(buf);
gotfile:
	/*
	 * test to see if the file exists
	 */
	if (stat(buf, &sbuf) >= 0)
	{
	    for (;;)
	    {
		fpurge(stdin);
		printf("\nFile exists.  Do you wish to overwrite it?");
		if (c == 'y' || c == 'Y')
		    break;
		else if (c == 'n' || c == 'N')
		    goto over;
		else
		    printf("\nPlease answer Y or N");
	    }
	}
	strcpy(savefile, buf);
	if ((savef = fopen(savefile, "w")) == NULL)
	    perror("Trek73: Save problems");
    } while (savef == NULL);

    /*
     * write out encrpyted file (after a stat)
     * The fwrite is to force allocation of the buffer before the write
     */
    save_file(savef);
    exit(0);
}

/*
 * save_file:
 *	Write the saved game on the file
 */

void
save_file(FILE *savef)
{
    /*
     * close any open save file
     */
    (void) fstat(fileno(savef), &sbuf);
    /*
     * DO NOT DELETE.  This forces stdio to allocate the output buffer
     * so that malloc doesn't get confused on restart
     */
    fwrite("junk", 1, 5, savef);

    fseek(savef, 0L, 0);
    encwrite(version, (unsigned int) (sbrk(0) - version), savef);
    fclose(savef);
    exit(0);
}

/*
 * restore:
 *	Restore a saved game from a file with elaborate checks for file
 *	integrity from cheaters
 */

void
restore(char *file, char **envp)
{
    int inf;
#ifdef PARANOID
    char syml;
#endif PARANOID
    extern char **environ;
    char buf[MAXSTR];
    STAT sbuf2;

#ifdef BSD
    (void) signal(SIGTSTP, SIG_IGN);
#endif
#ifdef SYSV
    (void) signal(SIGQUIT, SIG_IGN);
#endif

    if ((inf = open(file, 0)) < 0)
    {
	perror(file);
	return;
    }
    (void) fstat(inf, &sbuf2);
#ifdef PARANOID
#ifdef BSD
    syml = symlink(file);
#endif
#ifdef SYSV
    syml = link(file);
#endif
    if (unlink(file) < 0)
    {
	printf("Cannot unlink file\n");
	return;
    }
#endif PARANOID

    fflush(stdout);
    encread(buf, (unsigned int) (strlen(version) + 1), inf);
    if (strcmp(buf, version) != 0)
    {
	printf("Sorry, saved game is out of date.\n");
	return;
    }

    fflush(stdout);
    brk(version + sbuf2.st_size);
    lseek(inf, 0L, 0);
    encread(version, (unsigned int) sbuf2.st_size, inf);
#ifdef PARANOID
    /*
     * we do not close the file so that we will have a hold of the
     * inode for as long as possible
     */
	if (sbuf2.st_ino != sbuf.st_ino || sbuf2.st_dev != sbuf.st_dev)
	{
	    printf("Sorry, saved game is not in the same file.\n");
	    return;
	}
#ifdef NOTDEF
    /*
     * defeat multiple restarting from the same place
     */
	if (sbuf2.st_nlink != 1 || syml)
	{
	    printf("Cannot restore from a linked file %d %d\n", sbuf2.st_nlink, syml);
	    exit(1);
	}
#endif
#endif PARANOID
#ifdef BSD
    (void) signal(SIGTSTP, SIG_DFL);
#endif
#ifdef SYSV
    (void) signal(SIGQUIT, SIG_DFL);
#endif

    environ = envp;
    fpurge(stdin);
    playit();
#ifdef PARANOID
    syml = syml;		/* LINT */
#endif PARANOID
    /*NOTREACHED*/
}

/*
 * encwrite:
 *	Perform an encrypted write
 */

void
encwrite(char *start, unsigned int size, FILE *outf)
{
    char *ep;

    ep = encstr;

    while (size--)
    {
	putc(*start++ ^ *ep++, outf);
	if (*ep == '\0')
	    ep = encstr;
    }
}

/*
 * encread:
 *	Perform an encrypted read
 */

int
encread(char *start, unsigned int size, int inf)
{
    char *ep;
    int read_size;

    if ((read_size = read(inf, start, size)) == -1 || read_size == 0)
	return read_size;

    ep = encstr;

    while (size--)
    {
	*start++ ^= *ep++;
	if (*ep == '\0')
	    ep = encstr;
    }
    return read_size;
}
