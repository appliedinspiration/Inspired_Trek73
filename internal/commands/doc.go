// Package commands will hold the Inspired Trek73 command table and the
// numeric/text command parser that replaces the original lex/yacc based
// parser (see References/FreeBSD/trek73/src/command.l and grammar.y).
//
// Per project decisions, all original commands must remain supported.
// Text input is handled with a simpler, hand-written parser rather than
// a full grammar, while numeric command codes remain fully compatible
// with the original.
package commands
