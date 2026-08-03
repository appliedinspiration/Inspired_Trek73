#include <stdio.h>
#include <string.h>

main(argc, argv)
int argc;
char **argv;
{
	extern int yydebug;
	extern char Input[];
	extern char *Inptr;
	extern char parsed[];
	int i;
	char buf[BUFSIZ];

	if (argc > 1)
		yydebug = 1;
	else
		yydebug = 0;
	printf("\nCommand: ");
	while((i = fgets(buf, BUFSIZ, stdin) != NULL)) {
		(void) strcpy(Input, buf);
		Inptr = Input;
		parsed[0] = '\0';
		yyparse();
		printf("%s\n", parsed);
		printf("Command: ");
	}
}
