// antlr4 -Dlanguage=Go -o parser Expr.g4

// Expr.g4
grammar Expr;

// Tokens
AND: A N D;
OR: O R;
ID: [a-zA-Z_0-9][a-zA-Z_0-9]*;
WHITESPACE: [ \r\n\t]+ -> skip;

fragment A : [aA];
fragment N : [nN];
fragment D : [dD];
fragment O : [oO];
fragment R : [rR];

// Rules
start : expression EOF;

expression
   : expression op=(AND|OR) expression # ANDOR
   | ID                                    # ID
   | '(' expression ')'                    # Parenthesis
   ;