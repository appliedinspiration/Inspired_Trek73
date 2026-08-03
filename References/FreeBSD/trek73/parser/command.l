%{
	double	atof();
	char	Input[BUFSIZ];
	char	*Inptr;
%}

WS	[ \t]*

%%
[+-]?[0-9.]+	{ yylval.dval = atof(yytext); return(NUMBER); }
abort		{ return(ABORT); }
all		{ return(ALL); }
control		{ return(CONTROL); }
corbomite{WS}bluff	|
corbomite	{ return(CORB); }
course		{ return(COURSE); }
damage{WS}report	|
damage		{ return(DAMAGE); }
play{WS}dead	|
dead		{ return(DEAD); }
demand		{ return(DEMAND); }
self-destruct		|
self{WS}destruct	|
destruct	{ return(DESTR); }
detonate	{ return(DET); }
display		{ return(DISPLAY); }
elude		{ return(ELUDE); }
engineering	|
eng		{ return(ENG); }
fire		{ return(FIRE); }
help		{ return(HELP); }
jettison	|
jett		{ return(JETT); }
launch		{ return(LAUNCH); }
load		{ return(LOAD); }
lock		{ return(LOCK); }
parameters	|
params		|
param		{ return(PARAM); }
phasers		|
phaser		{ return(PHASER); }
alter{WS}power	|
power		{ return(POWER); }
probe		{ return(PROBE); }
pursue		{ return(PURSUE); }
rotate		{ return(ROTATE); }
save		{ return(SAVE); }
scan		{ return(SCAN); }
spread		{ return(SPREAD); }
status		{ return(STATUS); }
surrender	{ return(SURREND); }
survivors{WS}report	|
survivors	|
surv		{ return(SURV); }
tactical	{ return(TACTIC); }
photons		|
photon		|
torpedos	|
torpedo		|
torp		|
tubes		|
tube		{ return(TUBE); }
unload		{ return(UNLOAD); }
version		{ return(VERSION); }
warp{WS}factor	|
warp		{ return(WARP); }

[A-Za-z#]+	{ (void) strcpy( yylval.sval, yytext); return(NAME); }
{WS}		;
\n		{ return(-1); }
.		;
%%

#undef input
#undef unput

input()
{
	int c;

	c = *Inptr++;
	return (c);
}

unput(c)
char c;
{
	Inptr--;
}
