// antlr4 -Dlanguage=Go -o parser Expr.g4

// Expr.g4
grammar Expr;

// Tokens
AND: A N D;
OR: O R;
LOOP: L O O P;
ID: [a-zA-Z_0-9][a-zA-Z_0-9]*;
WHITESPACE: [ \r\n\t]+ -> skip;

fragment A : [aA];
fragment N : [nN];
fragment D : [dD];
fragment O : [oO];
fragment R : [rR];
fragment L : [lL];
fragment P : [pP];

// Rules
start : expression EOF;

expression
   : expression op=(AND|OR|LOOP) expression # ANDOR
   | ID                                    # ID
   | '(' expression ')'                    # Parenthesis
   ;