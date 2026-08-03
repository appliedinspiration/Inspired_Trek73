
/*
 * DEFS.H
 */

#include <sys/types.h>
#include <sys/stat.h>
#include <sys/time.h>
#include <stdio.h>
#include <stdlib.h>
#ifdef BSD
#include <strings.h>
#endif
#ifdef SYSV
#include <string.h>
#endif
#include <unistd.h>
#include <pwd.h>
#include <signal.h>
#include <ctype.h>
#include <fcntl.h>
#include <errno.h>
#include <math.h>
#include <setjmp.h>

#define Prototype extern

#include "defines.h"
#include "structs.h"
#include "union.h"

#include "protos.h"


