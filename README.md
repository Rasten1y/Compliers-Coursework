# gominic — минимальный самоприменимый компилятор подмножества Go → LLVM IR

## Подмножество языка
Поддерживается только то, что нужно для сборки бэкенда самого компилятора и простых тестов:
- Типы: `bool`, `int`/`int32`/`int64`/`uint32`/`uint64`, `float64`, `string` (как `{ i8*, i64 }`), массивы `[N]T`, срезы `[]T`, указатели `*T`, структуры `struct{...}`, `map[K]V`.
- Операции: арифметика `+ - * / %` (целые), `+ - * /` (float64); сравнения `== != < <= > >=`; логика `&& || !`; побитовые `& | << >>` (целые); присваивание `=` и `:=`; вызовы функций и методов (методы с receiver поддерживаются); возврат `return` с 0/1 значением или паттерном `(T, error)`; индексирование массивов/срезов; доступ к полям; слайс-выражения `s[i:j]`, `s[i:]`, `s[:j]`, `s[:]`.
- Управление: `if`/`if-else` без init statement; `for init; cond; post`, `for cond`, `for {}`; `break`/`continue`.
- Встроенные: `len` для массивов/срезов/map/строк; `make` для срезов и map; `append(slice, elem)` (только один элемент, не более); строковая конкатенация `+`.
- Рантайм: `gominic_print`, `gominic_printInt`, `gominic_println`, `gominic_memcpy`, `gominic_makeSlice`, `gominic_map_*`, `gominic_str_*`.
- Не поддерживается: `go`, `defer`, `select`, `switch`, `type switch`, функциональные литералы, `range`, `goto`, метки, `if` с init statement, `append` с более чем одним элементом, интерфейсы, каналы.

## Лексическая структура
- Идентификаторы: `[A-Za-z_][A-Za-z0-9_]*`
- Ключевые слова: `package` | `import` | `var` | `const` | `type` | `func` | `return` | `if` | `else` | `for` | `break` | `continue` | `true` | `false` | `len` | `make` | `append`
- Целые числа: `[0-9][0-9]*`
- Числа с плавающей точкой: `[0-9][0-9]* \. [0-9][0-9]*`
- Строковые литералы: `"` { любой символ кроме `"` и `\n` } `"`
- Булевы литералы: `true` | `false`
- Операторы: `+` | `-` | `*` | `/` | `%` | `==` | `!=` | `<` | `<=` | `>` | `>=` | `&&` | `||` | `!` | `&` | `|` | `<<` | `>>` | `=` | `:=`
- Разделители: `;` | `,` | `.` | `:` | `(` | `)` | `{` | `}` | `[` | `]`
- Комментарии: `//` { любой символ кроме `\n` } `\n` | `/*` { любой символ } `*/`

## Грамматика (EBNF)
```
Program        = "package" ident { ImportDecl } { TopDecl } .
ImportDecl     = "import" string_lit .
TopDecl        = FuncDecl | VarDecl | ConstDecl | TypeDecl .

FuncDecl       = "func" ident "(" [ ParamList ] ")" [ Result ] Block .
ParamList      = Param { "," Param } .
Param          = ident Type .
Result         = SingleResult | TupleResult .
SingleResult   = Type .
TupleResult    = "(" Type [ "," Type ] ")" .

VarDecl        = "var" ident [ Type ] [ "=" Expr ] .
ConstDecl      = "const" ident "=" Expr .
TypeDecl       = "type" ident Type .

Block          = "{" { Stmt } "}" .
Stmt           = DeclStmt | AssignStmt | IfStmt | ForStmt | ReturnStmt
               | ExprStmt | BreakStmt | ContinueStmt .
DeclStmt       = VarDecl .
AssignStmt     = ExprList ( "=" | ":=" ) ExprList .
IfStmt         = "if" Expr Block [ "else" Block ] .
ForStmt        = "for" ( [ SimpleStmt ] ";" [ Expr ] ";" [ SimpleStmt ] | Expr | ) Block .
SimpleStmt     = AssignStmt | ExprStmt .
ReturnStmt     = "return" [ ReturnValues ] .
ReturnValues   = EmptyReturn | SingleReturn | TupleReturn .
EmptyReturn    = .
SingleReturn   = Expr .
TupleReturn    = Expr "," Expr .
BreakStmt      = "break" .
ContinueStmt   = "continue" .
ExprStmt       = CallExpr .

Expr           = BinaryExpr | UnaryExpr | PrimaryExpr .
BinaryExpr     = Expr BinOp Expr .
BinOp          = "+" | "-" | "*" | "/" | "%" | "==" | "!=" | "<" | "<=" | ">" | ">=" | "&&" | "||" | "&" | "|" | "<<" | ">>" .
UnaryExpr      = ("+" | "-" | "!") Expr .
PrimaryExpr    = ident | literal | "(" Expr ")" | Selector | Index | Slice | Call .
Selector       = PrimaryExpr "." ident .
Index          = PrimaryExpr "[" Expr "]" .
Slice          = PrimaryExpr "[" [ Expr ] ":" [ Expr ] "]" .
Call           = PrimaryExpr "(" [ CallArgs ] ")" .
CallArgs       = Expr { "," Expr } .
ExprList       = Expr { "," Expr } .

Type           = BasicType | PointerType | ArrayType | SliceType | StructType | MapType .
BasicType      = "bool" | "int" | "int32" | "int64" | "uint32" | "uint64" | "float64" | "string" | ident .
PointerType    = "*" Type .
ArrayType      = "[" integer "]" Type .
SliceType      = "[]" Type .
StructType     = "struct" "{" { FieldDecl } "}" .
FieldDecl      = ident Type .
MapType        = "map" "[" Type "]" Type .
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
Для теста вывод следующий:
    9
    hi there
    after struct
    float ok
    after float
    0
    end

Для map/len используется `runtime\map.c`, для файловых обёрток — `runtime\io.c`.

### Сборка бэкенда
Команды для сборки бэкенда (компиляция `backend/llvm.go` вместе с `ir/ir.go`):
Сборка на Windows (через cmd):
```
go build -o gominic.exe ./cmd/gominic
```
```
.\gominic.exe -skip-check -S -v backend\llvm.go ir\ir.go > llvm.ll
```
```
clang -c -o llvm.o llvm.ll
```
Всплывет предупреждение, но все работает
Собирать дальше не имеет смысла, нет точки входа
