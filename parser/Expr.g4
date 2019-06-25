grammar Expr;

start
    : e=expr EOF
    ;

expr
    : '.' name=Identifier                       # ClosureMemberDot
    | expr '[' index=expr ']'                   # MemberIndex
    | expr '[' a=expr? ':' b=expr? ']'          # Slice
    | expr '.' name=Identifier                  # MemberDot
    | builtins                                  # BuiltinsList
    | expr '(' args=arguments? ')'              # Call
    | op=( '+' | '-' | '!' | 'not' ) expr       # Unary
    | expr op='..' expr                         # Binary
    | expr op=( '*' | '**' | '/' | '%' ) expr   # Binary
    | expr op=( '+' | '-' ) expr                # Binary
    | expr op=( '<' | '>' | '<=' | '>=' ) expr  # Binary
    | expr op='startsWith' expr                 # Binary
    | expr op='endsWith' expr                   # Binary
    | expr op='contains' expr                   # Binary
    | expr op='matches' pattern=expr            # Matches
    | expr op=( 'in' | 'not in' ) expr          # Binary
    | expr op=( '==' | '!=' ) expr              # Binary
    | expr op=And expr                          # Binary
    | expr op=Or expr                           # Binary
    | expr '?' e1=expr ':' e2=expr              # Ternary
    | Nil                                       # Nil
    | BooleanLiteral                            # Boolean
    | StringLiteral                             # String
    | IntegerLiteral                            # Integer
    | HexIntegerLiteral                         # Integer
    | FloatLiteral                              # Float
    | Identifier                                # Identifier
    | Pointer                                   # Pointer
    | arrayLiteral                              # Array
    | mapLiteral                                # Map
    | '(' expr ')'                              # Parenthesized
    ;

builtins
    : name='len'    '(' e=expr ')'                # BuiltinLen
    | name='all'    '(' e=expr ',' c=closure ')'  # Builtin
    | name='none'   '(' e=expr ',' c=closure ')'  # Builtin
    | name='any'    '(' e=expr ',' c=closure ')'  # Builtin
    | name='one'    '(' e=expr ',' c=closure ')'  # Builtin
    | name='filter' '(' e=expr ',' c=closure ')'  # Builtin
    | name='map'    '(' e=expr ',' c=closure ')'  # Builtin
    ;

closure
    : '{' body=expr '}'
    ;

arguments
    : list+=expr ( ',' list+=expr )*
    ;

arrayLiteral
    : '[' ']'
    | '[' list+=expr ( ',' list+=expr )* ','? ']'
    ;

mapLiteral
    : '{' '}'
    | '{' e=propertyNameAndValueList ','? '}'
    ;

propertyNameAndValueList
    : list+=propertyAssignment ( ',' list+=propertyAssignment )*
    ;

propertyAssignment
    : name=propertyName ':' value=expr
    ;

propertyName
    : Identifier
    | StringLiteral
    ;

/*****************************/
/*                           */
/*           LEXER           */
/*                           */
/*****************************/

OpenBracket                : '[';
CloseBracket               : ']';
OpenParen                  : '(';
CloseParen                 : ')';
OpenBrace                  : '{';
CloseBrace                 : '}';
SemiColon                  : ';';
Comma                      : ',';
Assign                     : '=';
QuestionMark               : '?';
Colon                      : ':';
Dot                        : '.';
Plus                       : '+';
Minus                      : '-';
Negate                     : '!';
Not                        : 'not';
Nil                        : 'nil';
Multiply                   : '*';
Exponent                   : '**';
Divide                     : '/';
Modulus                    : '%';
RightShiftArithmetic       : '>>';
LeftShiftArithmetic        : '<<';
LessThan                   : '<';
MoreThan                   : '>';
LessThanEquals             : '<=';
GreaterThanEquals          : '>=';
Equals                     : '==';
NotEquals                  : '!=';
Pointer                    : '#';
And                        : ( '&&' | 'and' );
Or                         : ( '||' | 'or' );
Builtins                   : ( 'all' | 'none' | 'any' | 'one' | 'filter' | 'map' );
Ops                        : ( 'startsWith' | 'endsWith' | 'contains' | 'matches' | 'in' | 'not in' );

BooleanLiteral
    : 'true'
    | 'false'
    ;

IntegerLiteral
    : '0'
    | [1-9] [0-9_]*
    ;

FloatLiteral
    : DecimalLiteral '.' Digit+
    | '.' Digit+
    ;

HexIntegerLiteral
    : '0' [xX] HexDigit+
    ;

Identifier
    : IdentifierStart IdentifierPart*
    ;

StringLiteral
    : '"' DoubleStringCharacter* '"'
    | '\'' SingleStringCharacter* '\''
    ;

WhiteSpaces
    : [\t\u000B\u000C\u0020\u00A0]+ -> channel(HIDDEN)
    ;

MultiLineComment
    : '/*' .*? '*/' -> channel(HIDDEN)
    ;

SingleLineComment
    : '//' ~[\r\n\u2028\u2029]* -> channel(HIDDEN)
    ;

LineTerminator
    : [\r\n\u2028\u2029] -> channel(HIDDEN)
    ;

UnexpectedCharacter
    : .
    ;

fragment DoubleStringCharacter
    : ~["\\\r\n]
    | '\\' EscapeSequence
    | LineContinuation
    ;

fragment SingleStringCharacter
    : ~['\\\r\n]
    | '\\' EscapeSequence
    | LineContinuation
    ;

fragment EscapeSequence
    : CharacterEscapeSequence
    | '0' // no digit ahead! TODO
    | HexEscapeSequence
    | UnicodeEscapeSequence
    ;

fragment CharacterEscapeSequence
    : SingleEscapeCharacter
    | NonEscapeCharacter
    ;

fragment HexEscapeSequence
    : 'x' HexDigit HexDigit
    ;

fragment UnicodeEscapeSequence
    : 'u' HexDigit HexDigit HexDigit HexDigit
    ;

fragment SingleEscapeCharacter
    : ['"\\bfnrtv]
    ;

fragment NonEscapeCharacter
    : ~['"\\bfnrtv0-9xu\r\n]
    ;

fragment EscapeCharacter
    : SingleEscapeCharacter
    | Digit
    | [xu]
    ;

fragment LineContinuation
    : '\\' LineTerminatorSequence
    ;

fragment LineTerminatorSequence
    : '\r\n'
    | LineTerminator
    ;

fragment IdentifierStart
    : Letter
    | [$_]
    ;

fragment IdentifierPart
    : Letter
    | Digit
    | [_]
    ;

fragment Letter
    : 'A'..'Z'
    | 'a'..'z'
    ;

fragment Digit
    : '0'..'9'
    ;

fragment DecimalLiteral
    : '0'
    | [1-9] Digit*
    ;

fragment HexDigit
    : [0-9a-fA-F]
    ;
