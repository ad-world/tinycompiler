# Tiny Compiler

Building a very basic compiler in Go. It will compile BASIC to C code. 
Using this experience to learn Golang.

## Parts
1. [Lexer](https://github.com/ad-world/tinycompiler/blob/master/lexer/lexer.go)
2. [Parser](https://github.com/ad-world/tinycompiler/blob/master/parser/parser.go)
3. [Emitter](https://github.com/ad-world/tinycompiler/blob/master/emitter/emitter.go)

## Formal Language

```
program ::= {statement}
statement ::= "PRINT" (expression | string) nl
    | "IF" comparison "THEN" nl {statement} "ENDIF" nl
    | "WHILE" comparison "REPEAT" nl {statement} "ENDWHILE" nl
    | "LABEL" ident nl
    | "GOTO" ident nl
    | "LET" ident "=" expression nl
    | "INPUT" ident nl
comparison ::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
expression ::= term {( "-" | "+" ) term}
term ::= unary {( "/" | "*" ) unary}
unary ::= ["+" | "-"] primary
primary ::= number | ident
nl ::= '\n'+
```

## Usage

This compiler works by compiling your BASIC code -> C code -> Executable.
You may use `compile.sh` which takes an input file and will create the necessary executable.

Use it like this
```
bash compile.sh {sample}.teeny
```

Then, check the `exec/` folder and you should have an executable named `{sample}`.

I'm currently working on adding more features to this compiler, such as functions and return types, logical operators, and more.