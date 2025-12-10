# gominic — минимальный самоприменимый компилятор подмножества Go → LLVM IR

## Подмножество языка
Поддерживается только то, что нужно для сборки самого компилятора и простых демо:
- Типы: `bool`, `int`/`int32`/`int64`/`uint32`/`uint64`, `float64`, `string` (как `{ i8*, i64 }`), массивы `[N]T`, срезы `[]T`, указатели `*T`, структуры `struct{...}`.
- Операции: арифметика `+ - * / %` (целые), `+ - * /` (float64); сравнения `== != < <= > >=`; логика `&& || !`; присваивание `=`; вызовы функций; возврат `return` с 0/1 значением; индексирование массивов/срезов; доступ к полям.
- Управление: `if`/`if-else` без init; `for init; cond; post`, `for cond`, `for {}`; `break`/`continue`.
- Встроенные: `len` для массивов/срезов/map/строк; `make` для срезов и map; `append(slice, elem)` (один элемент); строковая конкатенация `+`.
- Рантайм: `gominic_print`, `gominic_printInt`, `gominic_println`, `gominic_memcpy`, `gominic_makeSlice`, `gominic_map_*`, `gominic_str_*`.

## Лексическая структура
- Идентификаторы: `[A-Za-z_][A-Za-z0-9_]*`
- Ключевые слова: `package import var const type func return if else for break continue true false len make append`
- Литералы: целые в десятичной форме (`42`), `float64` (`1.23`), строковые в двойных кавычках без escape-сложностей.
- Операторы/знаки: `+ - * / % == != < <= > >= && || ! = ; , . : ( ) { } [ ]`

## Грамматика (EBNF, упрощённая под наше подмножество)
```
Program        = "package" ident { ImportDecl } { TopDecl } .
ImportDecl     = "import" string_lit .
TopDecl        = FuncDecl | VarDecl | ConstDecl | TypeDecl .

FuncDecl       = "func" ident "(" [ ParamList ] ")" [ Result ] Block .
ParamList      = Param { "," Param } .
Param          = ident Type .
Result         = Type | "(" Type ")" .

VarDecl        = "var" ident [ Type ] [ "=" Expr ] .
ConstDecl      = "const" ident "=" Expr .
TypeDecl       = "type" ident Type .

Block          = "{" { Stmt } "}" .
Stmt           = DeclStmt | AssignStmt | IfStmt | ForStmt | ReturnStmt
               | ExprStmt | BreakStmt | ContinueStmt .
DeclStmt       = VarDecl .
AssignStmt     = ExprList "=" ExprList .
IfStmt         = "if" Expr Block [ "else" Block ] .
ForStmt        = "for" ( [ SimpleStmt ] ";" [ Expr ] ";" [ SimpleStmt ] | Expr ) Block .
SimpleStmt     = AssignStmt | ExprStmt .
ReturnStmt     = "return" [ ExprList ] .
BreakStmt      = "break" .
ContinueStmt   = "continue" .
ExprStmt       = CallExpr .

Expr           = BinaryExpr | UnaryExpr | PrimaryExpr .
BinaryExpr     = Expr BinOp Expr .
BinOp          = "+" | "-" | "*" | "/" | "%" | "==" | "!=" | "<" | "<=" | ">" | ">=" | "&&" | "||" .
UnaryExpr      = ("+" | "-" | "!") Expr .
PrimaryExpr    = ident | literal | "(" Expr ")" | Selector | Index | Call .
Selector       = PrimaryExpr "." ident .
Index          = PrimaryExpr "[" Expr "]" .
Call           = PrimaryExpr "(" [ ExprList ] ")" .
ExprList       = Expr { "," Expr } .

Type           = ident | "*" Type | "[" [ integer ] "]" Type | "[]" Type
               | "struct" "{" { Field } "}" .
Field          = ident Type .
literal        = int_lit | float_lit | string_lit | boolean_lit .
boolean_lit    = "true" | "false" .
```

### Сборка и запуск (.ll → .o → .exe)
Команды для сборки файла subset_example.go (Если имя файла другое - меняйте на соответствующее)
Сборка на Windows (через cmd):
```
go build -o gominic.exe ./cmd/gominic
```
```
.\gominic.exe -S subset_example.go > subset_example.ll
```
```
clang -c -o subset_example.o subset_example.ll
```
Может всплыть предупреждение, не повлияет на работу
```
clang -o subset_example.exe subset_example.o runtime\print.c runtime\map.c runtime\io.c
```
Аналогично с большим предупреждением
```
.\subset_example.exe
```
В демо/минимальном случае достаточно `runtime\print.c`; для map/len используется `runtime\map.c`, для файловых обёрток — `runtime\io.c`.

## Демо без фронтенда
Для проверки стадии синтеза без фронтенда есть флаг `-demo`, который генерирует IR из захардкоженного дерева в `cmd/synthdemo/main.go`. Структура демо-дерева: функция `buildSampleModule` создаёт `ir.Module` с глобальной строкой и одной функцией `main`, которая вызывает `gominic_print`, `gominic_printInt`, `gominic_println`. 
Быстрый запуск:
```
go build -o gominic.exe ./cmd/gominic
cmd /c ".\gominic.exe -demo -S > demo.ll"
clang -c -o demo.o demo.ll
clang -o demo.exe demo.o runtime\print.c
.\demo.exe    # печатает: hello from demo42
```
Само демо‑дерево находится в `cmd/synthdemo/main.go` (функции `buildSampleModule`, `addStringGlobal`).
