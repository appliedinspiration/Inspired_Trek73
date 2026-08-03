%p 3000
%{
#ident "@(#) TREK73 $Header: /home/ncvs/trek73/src/command.l,v 1.3 1999/10/26 20:21:30 dillon Exp $"
/*
 * $Source: /home/ncvs/trek73/src/command.l,v $
 *
 * $Header: /home/ncvs/trek73/src/command.l,v 1.3 1999/10/26 20:21:30 dillon Exp $
 *
 * $Log: command.l,v $
 * Revision 1.3  1999/10/26 20:21:30  dillon
 *
 *     More cleanups
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
 * Revision 1.2  87/10/09  12:20:33  12:20:33  okamoto (Jeff Okamoto)
 * Added more positions for lex, added some null keywords
 * 
 * Revision 1.1  87/10/09  12:19:51  12:19:51  okamoto (Jeff Okamoto)
 * Initial revision
 */
	char	Input[BUFSIZ];
	char	*Inptr;

#define YY_INPUT(buf, result, max_size) \
	{ \
	    int c = *(unsigned char *)Inptr; \
 \
	    if (c == 0) { \
		result = YY_NULL; \
	    } else { \
		buf[0] = c; \
		result = 1; \
		++Inptr; \
	    } \
	}
%}

WS	[ \t]*

%%
[+-]?[0-9.]+		{ yylval.dval = atof(yytext); return(NUMBER); }
abort			{ return(ABORT); }
all			{ return(ALL); }
at			;
control			{ return(CONTROL); }
corbomite{WS}bluff	|
corbomite		{ return(CORB); }
turn			|
course			{ return(COURSE); }
damage{WS}report	|
damage			{ return(DAMAGE); }
play{WS}dead		|
dead			{ return(DEAD); }
delay			|
fuse			{ return(DELAY); }
demand			{ return(DEMAND); }
self-destruct		|
self{WS}destruct	|
destruct		{ return(DESTR); }
detonate		{ return(DET); }
display			{ return(DISPLAY); }
elude			{ return(ELUDE); }
engineering		|
eng			{ return(ENG); }
fire			{ return(FIRE); }
help			{ return(HELP); }
jettison		|
jett			{ return(JETT); }
launch			{ return(LAUNCH); }
load			{ return(LOAD); }
lock			{ return(LOCK); }
onto			;
on			;
parameters		|
params			|
param			{ return(PARAM); }
phasers			|
phaser			{ return(PHASER); }
pod			|
pods			;
alter{WS}power		|
power			{ return(POWER); }
probe			{ return(PROBE); }
prox			|
proximity		{ return(PROXIMITY); }
pursue			{ return(PURSUE); }
rotate			{ return(ROTATE); }
save			{ return(SAVE); }
scan			{ return(SCAN); }
spread			{ return(SPREAD); }
status			{ return(STATUS); }
surrender		{ return(SURREND); }
survivors{WS}report	|
survivors		|
surv			{ return(SURV); }
tactical		{ return(TACTIC); }
the			;
to			;
toward			|
towards			{ return(TOWARD); }
photons			|
photon			|
torpedos		|
torpedo			|
torp			|
torps			|
tubes			|
tube			{ return(TUBE); }
unload			{ return(UNLOAD); }
version			{ return(VERSION); }
warp{WS}factor		|
warp			{ return(WARP); }
with			;

[A-Za-z#]+		{ (void) strcpy( yylval.sval, yytext);
				return(NAME); }
{WS}			;
\n			{ return(-1); }
.			;
%%

