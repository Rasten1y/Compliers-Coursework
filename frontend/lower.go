package frontend

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"gominic/ir"
)

type lowerer struct {
	prog            *Program
	mod             *ir.Module
	fn              *ir.Function
	block           *ir.BasicBlock
	reg             int
	blockID         int
	localID         int
	locals          []map[string]*ir.Value
	prefix          string
	strID           int
	breakTargets    []string
	continueTargets []string
}

func newLowerer(prog *Program, mod *ir.Module) *lowerer {
	return &lowerer{
		prog:   prog,
		mod:    mod,
		locals: []map[string]*ir.Value{},
	}
}

// pushScope/popScope: управление областями видимости для локальных переменных
func (l *lowerer) pushScope() { l.locals = append(l.locals, make(map[string]*ir.Value)) }
func (l *lowerer) popScope() {
	if len(l.locals) > 0 {
		l.locals = l.locals[:len(l.locals)-1]
	}
}

func (l *lowerer) setLocal(name string, v *ir.Value) {
	l.locals[len(l.locals)-1][name] = v
}

func (l *lowerer) getLocal(name string) *ir.Value {
	for i := len(l.locals) - 1; i >= 0; i-- {
		if v, ok := l.locals[i][name]; ok {
			return v
		}
	}
	return nil
}

func (l *lowerer) newTemp(ty *ir.TypeDesc) *ir.Value {
	l.reg++
	return ir.NewRegister(fmt.Sprintf("t%d", l.reg), ty)
}

func (l *lowerer) newBlockName(base string) string {
	l.blockID++
	return fmt.Sprintf("%s%d", base, l.blockID)
}

func (l *lowerer) newBlock(base string) *ir.BasicBlock {
	name := base
	if base == "" {
		name = l.newBlockName("bb")
	}
	bb := ir.NewBasicBlock(name)
	l.fn.Blocks = append(l.fn.Blocks, bb)
	return bb
}

// TopDecl = VarDecl | ConstDecl | TypeDecl
// VarDecl = "var" ident [ Type ] [ "=" Expr ]
// ConstDecl = "const" ident "=" Expr
func (l *lowerer) lowerGlobals(decls []ast.Decl) error {
	for _, d := range decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || (gd.Tok != token.VAR && gd.Tok != token.CONST) {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range vs.Names {
				obj := l.prog.TypesInfo.Defs[name]
				if obj == nil {
					return fmt.Errorf("no def for global %s", name.Name)
				}
				ty, err := goTypeToIR(obj.Type())
				if err != nil {
					return err
				}
				full := name.Name
				if l.prefix != "" {
					full = l.prefix + name.Name
				}
				init := "zeroinitializer"
				if i < len(vs.Values) {
					if c, err := l.constValue(vs.Values[i]); err == nil {
						if s, errStr := constantString(c); errStr == nil {
							init = s
						}
					}
				}
				l.mod.Globals = append(l.mod.Globals, ir.Global{
					Name:  full,
					Type:  ty,
					Value: init,
					Align: alignOfIRType(ty),
				})
			}
		}
	}
	return nil
}

// FuncDecl = "func" ident "(" [ ParamList ] ")" [ Result ] Block
// ParamList = Param { "," Param }
// Param = ident Type
// Result = SingleResult | TupleResult
// SingleResult = Type
// TupleResult = "(" Type [ "," Type ] ")"
func (l *lowerer) lowerFunc(fnDecl *ast.FuncDecl) error {
	sig, ok := l.prog.TypesInfo.Defs[fnDecl.Name].Type().(*types.Signature)
	if !ok {
		return fmt.Errorf("func %s: no signature", fnDecl.Name.Name)
	}
	// Для самоприменимости поддерживаем множественный возврат только формы (T, error).
	if sig.Results().Len() > 1 {
		// Сейчас поддерживаем только форму (T, error).
		if sig.Results().Len() == 2 {
			// Проверяем, что второй результат имеет тип error.
			if errType, ok := sig.Results().At(1).Type().(*types.Named); ok {
				if errType.Obj() != nil && errType.Obj().Name() == "error" {
					// Это форма (T, error), она поддерживается.
					// Возвращаем только первое значение, error обрабатывается отдельно.
				} else {
					return fmt.Errorf("unsupported: multiple returns (not (T, error) pattern)")
				}
			} else {
				return fmt.Errorf("unsupported: multiple returns (not (T, error) pattern)")
			}
		} else {
			return fmt.Errorf("unsupported: multiple returns (more than 2)")
		}
	}

	fnName := fnDecl.Name.Name
	// Спецслучай: функция main всегда называется "main" без префикса пакета.
	if fnName == "main" {
		fnName = "main"
	} else if sig.Recv() != nil {
		// Если это метод (есть receiver), добавляем префикс с именем типа.
		recvType := sig.Recv().Type()
		typeName := ""
		if named, ok := recvType.(*types.Named); ok {
			if named.Obj() != nil && named.Obj().Pkg() != nil {
				typeName = named.Obj().Pkg().Name() + "." + named.Obj().Name()
			} else if named.Obj() != nil {
				typeName = named.Obj().Name()
			}
		} else if ptr, ok := recvType.(*types.Pointer); ok {
			if named, ok := ptr.Elem().(*types.Named); ok {
				if named.Obj() != nil && named.Obj().Pkg() != nil {
					typeName = named.Obj().Pkg().Name() + "." + named.Obj().Name()
				} else if named.Obj() != nil {
					typeName = named.Obj().Name()
				}
			}
		}
		if typeName != "" {
			fnName = typeName + "." + fnName
		}
	} else {
		// Для функций пакета добавляем префикс, если это не main.
		// Если функция экспортируемая и пакет не main, добавляем префикс пакета.
		if l.prefix != "" {
			fnName = l.prefix + fnName
		} else if len(fnName) > 0 && fnName[0] >= 65 && fnName[0] <= 90 {
			// Функция экспортируемая.
			// Берём имя пакета из prog.
			if l.prog != nil && l.prog.PkgName != "" && l.prog.PkgName != "main" {
				fnName = l.prog.PkgName + "." + fnName
			}
		}
	}
	fn := &ir.Function{Name: fnName}
	for i := 0; i < sig.Params().Len(); i++ {
		pt, err := goTypeToIR(sig.Params().At(i).Type())
		if err != nil {
			return err
		}
		fn.Params = append(fn.Params, ir.NewParam(sig.Params().At(i).Name(), pt))
	}
	if sig.Results().Len() == 1 {
		rt, err := goTypeToIR(sig.Results().At(0).Type())
		if err != nil {
			return err
		}
		fn.Results = []*ir.TypeDesc{rt}
	}

	l.fn = fn
	entry := ir.NewBasicBlock("entry")
	l.fn.Blocks = append(l.fn.Blocks, entry)
	l.block = entry
	l.reg = 0
	l.blockID = 0
	l.locals = nil
	l.pushScope()

	// Параметры функции: выделяем память на стеке (alloca) и сохраняем значение
	// Это нужно, т.к. параметры передаются по значению, но мы работаем с указателями
	for i := 0; i < len(fn.Params); i++ {
		p := fn.Params[i]
		addr := ir.NewRegister(fmt.Sprintf("p%d.addr", i), ir.PtrTo(p.Type()))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: p.Type(), AllocaAlign: alignOfIRType(p.Type())})
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: p, StoreDst: addr, StoreAlign: alignOfIRType(p.Type())})
		l.setLocal(p.Name(), addr)
	}

	if err := l.lowerBlock(fnDecl.Body); err != nil {
		return err
	}
	if l.block.Terminator == nil {
		// Если функция не имеет явного return, добавляем неявный возврат нулевого значения
		if len(fn.Results) > 0 {
			retTy := fn.Results[0]
			var zeroVal *ir.Value
			if retTy.Kind == ir.KindPointer {
				zeroVal = ir.NewConstant("null", retTy)
			} else if retTy.Kind == ir.KindString {
				zeroVal = ir.NewConstant(`{ i8* null, i64 0 }`, retTy)
			} else if retTy.Basic != "" {
				zeroVal = ir.NewConstant("0", retTy)
			} else {
				zeroVal = ir.NewConstant("null", ir.PtrI8)
			}
			l.block.Terminator = &ir.Instruction{Kind: ir.InstrReturn, RetVals: []*ir.Value{zeroVal}}
		} else {
			l.block.Terminator = &ir.Instruction{Kind: ir.InstrReturn}
		}
	}
	l.mod.AddFunction(fn)
	return nil
}

// Block = "{" { Stmt } "}"
func (l *lowerer) lowerBlock(b *ast.BlockStmt) error {
	l.pushScope()
	for _, s := range b.List {
		if err := l.lowerStmt(s); err != nil {
			return err
		}
	}
	l.popScope()
	return nil
}

// Stmt = DeclStmt | AssignStmt | IfStmt | ForStmt | ReturnStmt
//
//	| ExprStmt | BreakStmt | ContinueStmt
func (l *lowerer) lowerStmt(s ast.Stmt) error {
	switch st := s.(type) {
	case *ast.BlockStmt:
		return l.lowerBlock(st)
	case *ast.DeclStmt:
		return l.lowerDeclStmt(st)
	case *ast.ExprStmt:
		_, err := l.lowerExpr(st.X)
		return err
	// ReturnStmt = "return" [ ReturnValues ]
	// ReturnValues = EmptyReturn | SingleReturn | TupleReturn
	// EmptyReturn = .
	// SingleReturn = Expr
	// TupleReturn = Expr "," Expr
	case *ast.ReturnStmt:
		var vals []*ir.Value
		for i := 0; i < len(st.Results); i++ {
			// Проверяем, что возвращаемое значение — это идентификатор nil.
			isNilIdent := false
			if ident, ok := st.Results[i].(*ast.Ident); ok && ident.Name == "nil" {
				isNilIdent = true
			}
			v, err := l.lowerExpr(st.Results[i])
			if err != nil {
				return err
			}
			// Если возвращается nil, а функция ждёт срез, создаём нулевой срез.
			if l.fn != nil && i < len(l.fn.Results) {
				expectedTy := l.fn.Results[i]
				if expectedTy.Kind == ir.KindSlice && isNilIdent {
					// Нулевой срез: { null, 0, 0 }.
					elemPtrTy := ir.PtrTo(expectedTy.Elem)
					zeroPtr := ir.NewConstant("null", elemPtrTy)
					zeroLen := ir.NewConstant("0", ir.I64)
					zeroCap := ir.NewConstant("0", ir.I64)
					zeroSlice := fmt.Sprintf("{ %s %s, i64 %s, i64 %s }", elemPtrTy.String(), zeroPtr.Name(), zeroLen.Name(), zeroCap.Name())
					v = ir.NewConstant(zeroSlice, expectedTy)
				}
			}
			vals = append(vals, v)
		}
		l.block.Terminator = &ir.Instruction{Kind: ir.InstrReturn, RetVals: vals}
	case *ast.AssignStmt:
		return l.lowerAssign(st)
	case *ast.IfStmt:
		return l.lowerIf(st)
	case *ast.ForStmt:
		return l.lowerFor(st)
	case *ast.BranchStmt:
		return l.lowerBranch(st)
	default:
		return fmt.Errorf("unsupported stmt %T", s)
	}
	return nil
}

// DeclStmt = VarDecl
// VarDecl = "var" ident [ Type ] [ "=" Expr ]
func (l *lowerer) lowerDeclStmt(ds *ast.DeclStmt) error {
	gd, ok := ds.Decl.(*ast.GenDecl)
	if !ok || gd.Tok != token.VAR {
		return fmt.Errorf("unsupported declstmt")
	}
	for _, spec := range gd.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for i, name := range vs.Names {
			if name.Name == "_" {
				continue
			}
			var rhsVal *ir.Value
			if i < len(vs.Values) {
				v, err := l.lowerExpr(vs.Values[i])
				if err != nil {
					return err
				}
				rhsVal = v
			}
			var ty *ir.TypeDesc
			if vs.Type != nil {
				if tGo := l.prog.TypesInfo.TypeOf(vs.Type); tGo != nil {
					if conv, err := goTypeToIR(tGo); err == nil {
						ty = conv
					}
				}
			}
			if ty == nil && rhsVal != nil {
				ty = rhsVal.Type()
			}
			if ty == nil {
				if tGo := l.prog.TypesInfo.TypeOf(name); tGo != nil {
					if conv, err := goTypeToIR(tGo); err == nil {
						ty = conv
					}
				}
			}
			if ty == nil {
				return fmt.Errorf("cannot infer type for %s", name.Name)
			}
			addr := l.ensureLocal(name.Name, ty)
			if rhsVal != nil {
				l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: rhsVal, StoreDst: addr, StoreAlign: alignOfIRType(ty)})
			}
		}
	}
	return nil
}

// BreakStmt = "break"
// ContinueStmt = "continue"
func (l *lowerer) lowerBranch(st *ast.BranchStmt) error {
	switch st.Tok {
	case token.BREAK:
		if len(l.breakTargets) == 0 {
			return fmt.Errorf("break not in loop")
		}
		target := l.breakTargets[len(l.breakTargets)-1]
		l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: target}
	case token.CONTINUE:
		if len(l.continueTargets) == 0 {
			return fmt.Errorf("continue not in loop")
		}
		target := l.continueTargets[len(l.continueTargets)-1]
		l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: target}
	default:
		return fmt.Errorf("unsupported branch %s", st.Tok.String())
	}
	return nil
}

// AssignStmt = ExprList ( "=" | ":=" ) ExprList
func (l *lowerer) lowerAssign(st *ast.AssignStmt) error {
	if st.Tok != token.ASSIGN && st.Tok != token.DEFINE {
		return fmt.Errorf("unsupported assign tok %s", st.Tok.String())
	}
	// Спецслучай: чтение из map с двумя результатами.
	if len(st.Lhs) == 2 && len(st.Rhs) == 1 {
		if idx, ok := st.Rhs[0].(*ast.IndexExpr); ok {
			if _, okm := l.prog.TypesInfo.TypeOf(idx.X).Underlying().(*types.Map); okm {
				return l.lowerMapGet(idx, st.Lhs[0], st.Lhs[1], st.Tok == token.DEFINE)
			}
		}
		// Также обрабатываем вызовы функций с двумя результатами.
		if call, ok := st.Rhs[0].(*ast.CallExpr); ok {
			// Пробуем получить сигнатуру функции из доступных источников.
			var sig *types.Signature
			callType := l.prog.TypesInfo.TypeOf(call)
			if sig2, ok := callType.(*types.Signature); ok {
				sig = sig2
			} else {
				// Пробуем получить сигнатуру из идентификатора функции.
				if fnIdent, ok := call.Fun.(*ast.Ident); ok {
					if obj, ok := l.prog.TypesInfo.Uses[fnIdent]; ok {
						if fn, ok := obj.(*types.Func); ok {
							sig = fn.Type().(*types.Signature)
						}
					}
				} else if selExpr, ok := call.Fun.(*ast.SelectorExpr); ok {
					if obj, ok := l.prog.TypesInfo.Uses[selExpr.Sel]; ok {
						if fn, ok := obj.(*types.Func); ok {
							sig = fn.Type().(*types.Signature)
						}
					}
					// Также пробуем Selections для вызовов методов.
					if sel, ok := l.prog.TypesInfo.Selections[selExpr]; ok {
						if fn := sel.Obj(); fn != nil {
							if fn2, ok := fn.(*types.Func); ok {
								sig = fn2.Type().(*types.Signature)
							}
						}
					}
				}
			}
			if sig != nil && sig.Results().Len() == 2 {
				// Если функция возвращает 2 значения, обрабатываем их отдельно.
				firstVal, err := l.lowerCall(call)
				if err != nil {
					return err
				}
				firstTy := firstVal.Type()
				if firstTy == nil {
					firstGoTy := sig.Results().At(0).Type()
					firstTy, err = goTypeToIR(firstGoTy)
					if err != nil {
						return err
					}
				}
				// Обрабатываем первый LHS.
				lhs0Val, err := l.lowerExpr(st.Lhs[0])
				if err != nil {
					if st.Tok == token.DEFINE {
						lhs0Ident, ok := st.Lhs[0].(*ast.Ident)
						if !ok {
							return fmt.Errorf("unsupported LHS type for define")
						}
						addr := l.ensureLocal(lhs0Ident.Name, firstTy)
						l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: addr, StoreAlign: alignOfIRType(firstTy)})
					} else {
						return err
					}
				} else {
					l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: lhs0Val, StoreAlign: alignOfIRType(firstTy)})
				}
				// Обрабатываем второй LHS (обычно error).
				returnTy, err := goTypeToIR(sig.Results().At(1).Type())
				if err != nil {
					return err
				}
				var zeroVal *ir.Value
				if returnTy.Kind == ir.KindPointer {
					zeroVal = ir.NewConstant("null", returnTy)
				} else if returnTy.Kind == ir.KindString {
					zeroVal = ir.NewConstant(`{ i8* null, i64 0 }`, returnTy)
				} else if returnTy.Basic != "" {
					zeroVal = ir.NewConstant("0", returnTy)
				} else {
					zeroVal = ir.NewConstant("null", ir.PtrI8)
				}
				lhs1Val, err := l.lowerExpr(st.Lhs[1])
				if err != nil {
					if st.Tok == token.DEFINE {
						lhs1Ident, ok := st.Lhs[1].(*ast.Ident)
						if !ok {
							return fmt.Errorf("unsupported LHS type for define")
						}
						addr := l.ensureLocal(lhs1Ident.Name, returnTy)
						l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: addr, StoreAlign: alignOfIRType(returnTy)})
					} else {
						return err
					}
				} else {
					l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: lhs1Val, StoreAlign: alignOfIRType(returnTy)})
				}
				return nil
			}
		}
	}
	// Поддержка множественного присваивания из функции с несколькими результатами.
	// Обрабатываем (T, error) и похожие случаи множественного возврата.
	if len(st.Rhs) == 1 {
		// Проверяем, что RHS — вызов функции с несколькими результатами.
		if call, ok := st.Rhs[0].(*ast.CallExpr); ok {
			callType := l.prog.TypesInfo.TypeOf(call)
			if sig, ok := callType.(*types.Signature); ok {
				numReturns := sig.Results().Len()
				// Обрабатываем случай, когда функция возвращает несколько значений.
				if numReturns > 1 {
					// Если число LHS равно числу результатов, присваиваем все.
					if numReturns == len(st.Lhs) {
						// Количество результатов совпадает с количеством LHS.
						// Понижаем вызов и получаем первое значение.
						firstVal, err := l.lowerCall(call)
						if err != nil {
							return err
						}
						// Записываем первое значение в первый LHS.
						firstTy := firstVal.Type()
						if firstTy == nil && numReturns > 0 {
							// Берём тип из сигнатуры.
							firstGoTy := sig.Results().At(0).Type()
							firstTy, err = goTypeToIR(firstGoTy)
							if err != nil {
								return err
							}
						}
						// Обрабатываем первый LHS.
						lhs0Val, err := l.lowerExpr(st.Lhs[0])
						if err != nil {
							// LHS ещё не существует, создаём его.
							if st.Tok == token.DEFINE {
								lhs0Ident, ok := st.Lhs[0].(*ast.Ident)
								if !ok {
									return fmt.Errorf("unsupported LHS type for define")
								}
								addr := l.ensureLocal(lhs0Ident.Name, firstTy)
								l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: addr, StoreAlign: alignOfIRType(firstTy)})
							} else {
								return err
							}
						} else {
							// LHS уже существует, выполняем присваивание.
							l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: lhs0Val, StoreAlign: alignOfIRType(firstTy)})
						}
						// Обрабатываем остальные LHS для множественного возврата.
						for i := 1; i < len(st.Lhs); i++ {
							returnTy, err := goTypeToIR(sig.Results().At(i).Type())
							if err != nil {
								return err
							}
							// Пока заполняем остальные переменные нулевыми значениями.
							// В полной реализации нужно извлекать все возвращаемые значения.
							var zeroVal *ir.Value
							if returnTy.Kind == ir.KindPointer {
								zeroVal = ir.NewConstant("null", returnTy)
							} else if returnTy.Kind == ir.KindString {
								zeroVal = ir.NewConstant(`{ i8* null, i64 0 }`, returnTy)
							} else if returnTy.Basic != "" {
								zeroVal = ir.NewConstant("0", returnTy)
							} else {
								zeroVal = ir.NewConstant("null", ir.PtrI8)
							}
							lhsVal, err := l.lowerExpr(st.Lhs[i])
							if err != nil {
								// LHS ещё не существует, создаём его.
								if st.Tok == token.DEFINE {
									lhsIdent, ok := st.Lhs[i].(*ast.Ident)
									if !ok {
										return fmt.Errorf("unsupported LHS type for define")
									}
									addr := l.ensureLocal(lhsIdent.Name, returnTy)
									l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: addr, StoreAlign: alignOfIRType(returnTy)})
								} else {
									return err
								}
							} else {
								// LHS уже существует, выполняем присваивание.
								l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: lhsVal, StoreAlign: alignOfIRType(returnTy)})
							}
						}
						return nil
					} else if len(st.Lhs) > 0 && len(st.Lhs) < numReturns {
						// Если LHS меньше, чем результатов, присваиваем только первые.
						// Функция вернула больше значений, используем первые N.
						for i := 0; i < len(st.Lhs); i++ {
							returnTy, err := goTypeToIR(sig.Results().At(i).Type())
							if err != nil {
								return err
							}
							// Пока вызываем функцию один раз и используем первое значение.
							// Затем заполняем остальные LHS нулевыми значениями.
							if i == 0 {
								firstVal, err := l.lowerCall(call)
								if err != nil {
									return err
								}
								firstTy := firstVal.Type()
								if firstTy == nil {
									firstTy = returnTy
								}
								lhsVal, err := l.lowerExpr(st.Lhs[i])
								if err != nil {
									if st.Tok == token.DEFINE {
										lhsIdent, ok := st.Lhs[i].(*ast.Ident)
										if !ok {
											return fmt.Errorf("unsupported LHS type for define")
										}
										addr := l.ensureLocal(lhsIdent.Name, firstTy)
										l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: addr, StoreAlign: alignOfIRType(firstTy)})
									} else {
										return err
									}
								} else {
									l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: lhsVal, StoreAlign: alignOfIRType(firstTy)})
								}
							} else {
								// Присваиваем нули оставшимся LHS.
								var zeroVal *ir.Value
								if returnTy.Kind == ir.KindPointer {
									zeroVal = ir.NewConstant("null", returnTy)
								} else if returnTy.Kind == ir.KindString {
									zeroVal = ir.NewConstant(`{ i8* null, i64 0 }`, returnTy)
								} else if returnTy.Basic != "" {
									zeroVal = ir.NewConstant("0", returnTy)
								} else {
									zeroVal = ir.NewConstant("null", ir.PtrI8)
								}
								lhsVal, err := l.lowerExpr(st.Lhs[i])
								if err != nil {
									if st.Tok == token.DEFINE {
										lhsIdent, ok := st.Lhs[i].(*ast.Ident)
										if !ok {
											return fmt.Errorf("unsupported LHS type for define")
										}
										addr := l.ensureLocal(lhsIdent.Name, returnTy)
										l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: addr, StoreAlign: alignOfIRType(returnTy)})
									} else {
										return err
									}
								} else {
									l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: lhsVal, StoreAlign: alignOfIRType(returnTy)})
								}
							}
						}
						return nil
					} else if len(st.Lhs) == 1 && numReturns > 1 {
						// Если LHS один, а значений несколько, берём первое.
						firstVal, err := l.lowerCall(call)
						if err != nil {
							return err
						}
						lhs, ok := st.Lhs[0].(*ast.Ident)
						if !ok {
							return fmt.Errorf("unsupported lhs")
						}
						if lhs.Name == "_" {
							// вычисляем только ради побочных эффектов
							return nil
						}
						rhsTy := firstVal.Type()
						if rhsTy == nil {
							firstGoTy := sig.Results().At(0).Type()
							rhsTy, err = goTypeToIR(firstGoTy)
							if err != nil {
								return err
							}
						}
						if st.Tok == token.DEFINE {
							addr := l.ensureLocal(lhs.Name, rhsTy)
							l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: addr, StoreAlign: alignOfIRType(rhsTy)})
						} else {
							addr := l.getLocal(lhs.Name)
							if addr == nil {
								return fmt.Errorf("unknown var %s", lhs.Name)
							}
							l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: addr, StoreAlign: alignOfIRType(rhsTy)})
						}
						return nil
					} else if len(st.Lhs) == 0 && numReturns > 0 {
						// Если LHS нет, вычисляем только ради побочных эффектов.
						_, err := l.lowerCall(call)
						return err
					}
					// Если возвращается одно значение и LHS один, идём в обычную ветку.
				}
			} else if callType == nil {
				// TypeOf вернул nil, пробуем взять тип из функции.
				// Такое бывает у некоторых вызовов.
				if fnIdent, ok := call.Fun.(*ast.Ident); ok {
					if obj, ok := l.prog.TypesInfo.Uses[fnIdent]; ok {
						if fn, ok := obj.(*types.Func); ok {
							sig := fn.Type().(*types.Signature)
							numReturns := sig.Results().Len()
							// Обрабатываем множественный возврат даже при nil в TypeOf(call).
							if numReturns > 1 && numReturns == len(st.Lhs) {
								// Дальше логика такая же.
								firstVal, err := l.lowerCall(call)
								if err != nil {
									return err
								}
								firstTy := firstVal.Type()
								if firstTy == nil {
									firstGoTy := sig.Results().At(0).Type()
									firstTy, err = goTypeToIR(firstGoTy)
									if err != nil {
										return err
									}
								}
								lhs0Val, err := l.lowerExpr(st.Lhs[0])
								if err != nil {
									if st.Tok == token.DEFINE {
										lhs0Ident, ok := st.Lhs[0].(*ast.Ident)
										if !ok {
											return fmt.Errorf("unsupported LHS type for define")
										}
										addr := l.ensureLocal(lhs0Ident.Name, firstTy)
										l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: addr, StoreAlign: alignOfIRType(firstTy)})
									} else {
										return err
									}
								} else {
									l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: lhs0Val, StoreAlign: alignOfIRType(firstTy)})
								}
								for i := 1; i < len(st.Lhs); i++ {
									returnTy, err := goTypeToIR(sig.Results().At(i).Type())
									if err != nil {
										return err
									}
									var zeroVal *ir.Value
									if returnTy.Kind == ir.KindPointer {
										zeroVal = ir.NewConstant("null", returnTy)
									} else if returnTy.Kind == ir.KindString {
										zeroVal = ir.NewConstant(`{ i8* null, i64 0 }`, returnTy)
									} else if returnTy.Basic != "" {
										zeroVal = ir.NewConstant("0", returnTy)
									} else {
										zeroVal = ir.NewConstant("null", ir.PtrI8)
									}
									lhsVal, err := l.lowerExpr(st.Lhs[i])
									if err != nil {
										if st.Tok == token.DEFINE {
											lhsIdent, ok := st.Lhs[i].(*ast.Ident)
											if !ok {
												return fmt.Errorf("unsupported LHS type for define")
											}
											addr := l.ensureLocal(lhsIdent.Name, returnTy)
											l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: addr, StoreAlign: alignOfIRType(returnTy)})
										} else {
											return err
										}
									} else {
										l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: lhsVal, StoreAlign: alignOfIRType(returnTy)})
									}
								}
								return nil
							}
						}
					}
				} else if selExpr, ok := call.Fun.(*ast.SelectorExpr); ok {
					// Пробуем получить тип из selector-выражения.
					if obj, ok := l.prog.TypesInfo.Uses[selExpr.Sel]; ok {
						if fn, ok := obj.(*types.Func); ok {
							sig := fn.Type().(*types.Signature)
							numReturns := sig.Results().Len()
							if numReturns > 1 && numReturns == len(st.Lhs) {
								// Дальше логика такая же.
								firstVal, err := l.lowerCall(call)
								if err != nil {
									return err
								}
								firstTy := firstVal.Type()
								if firstTy == nil {
									firstGoTy := sig.Results().At(0).Type()
									firstTy, err = goTypeToIR(firstGoTy)
									if err != nil {
										return err
									}
								}
								lhs0Val, err := l.lowerExpr(st.Lhs[0])
								if err != nil {
									if st.Tok == token.DEFINE {
										lhs0Ident, ok := st.Lhs[0].(*ast.Ident)
										if !ok {
											return fmt.Errorf("unsupported LHS type for define")
										}
										addr := l.ensureLocal(lhs0Ident.Name, firstTy)
										l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: addr, StoreAlign: alignOfIRType(firstTy)})
									} else {
										return err
									}
								} else {
									l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: lhs0Val, StoreAlign: alignOfIRType(firstTy)})
								}
								for i := 1; i < len(st.Lhs); i++ {
									returnTy, err := goTypeToIR(sig.Results().At(i).Type())
									if err != nil {
										return err
									}
									var zeroVal *ir.Value
									if returnTy.Kind == ir.KindPointer {
										zeroVal = ir.NewConstant("null", returnTy)
									} else if returnTy.Kind == ir.KindString {
										zeroVal = ir.NewConstant(`{ i8* null, i64 0 }`, returnTy)
									} else if returnTy.Basic != "" {
										zeroVal = ir.NewConstant("0", returnTy)
									} else {
										zeroVal = ir.NewConstant("null", ir.PtrI8)
									}
									lhsVal, err := l.lowerExpr(st.Lhs[i])
									if err != nil {
										if st.Tok == token.DEFINE {
											lhsIdent, ok := st.Lhs[i].(*ast.Ident)
											if !ok {
												return fmt.Errorf("unsupported LHS type for define")
											}
											addr := l.ensureLocal(lhsIdent.Name, returnTy)
											l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: addr, StoreAlign: alignOfIRType(returnTy)})
										} else {
											return err
										}
									} else {
										l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: lhsVal, StoreAlign: alignOfIRType(returnTy)})
									}
								}
								return nil
							}
						}
					}
				}
			}
		}
	}
	if len(st.Lhs) != len(st.Rhs) {
		// Добавляем больше контекста в ошибку
		lhsStr := fmt.Sprintf("%d", len(st.Lhs))
		rhsStr := fmt.Sprintf("%d", len(st.Rhs))
		if len(st.Rhs) == 1 {
			if call, ok := st.Rhs[0].(*ast.CallExpr); ok {
				callType := l.prog.TypesInfo.TypeOf(call)
				if sig, ok := callType.(*types.Signature); ok {
					rhsStr = rhsStr + fmt.Sprintf(" (function returns %d values)", sig.Results().Len())
				}
			}
		}
		return fmt.Errorf("unsupported: assign arity mismatch: LHS=%s, RHS=%s", lhsStr, rhsStr)
	}
	for i := 0; i < len(st.Lhs); i++ {
		// присваивание в map?
		if idx, ok := st.Lhs[i].(*ast.IndexExpr); ok {
			if _, okm := l.prog.TypesInfo.TypeOf(idx.X).Underlying().(*types.Map); okm {
				return l.lowerMapSet(idx, st.Rhs[i])
			}
			// Присваивание в элемент массива/среза: arr[i] = value.
			rhsVal, err := l.lowerExpr(st.Rhs[i])
			if err != nil {
				return err
			}
			elemPtr, err := l.lowerAddr(idx)
			if err != nil {
				return err
			}
			l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: rhsVal, StoreDst: elemPtr, StoreAlign: alignOfIRType(rhsVal.Type())})
			continue
		}
		// Присваивание в поле структуры: s.field = value.
		if selExpr, ok := st.Lhs[i].(*ast.SelectorExpr); ok {
			rhsVal, err := l.lowerExpr(st.Rhs[i])
			if err != nil {
				return err
			}
			fieldPtr, err := l.lowerAddr(selExpr)
			if err != nil {
				return err
			}
			l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: rhsVal, StoreDst: fieldPtr, StoreAlign: alignOfIRType(rhsVal.Type())})
			continue
		}
		// Присваивание по разыменованию: *ptr = value.
		if starExpr, ok := st.Lhs[i].(*ast.StarExpr); ok {
			rhsVal, err := l.lowerExpr(st.Rhs[i])
			if err != nil {
				return err
			}
			ptrVal, err := l.lowerExpr(starExpr.X)
			if err != nil {
				return err
			}
			if ptrVal.Type() == nil || ptrVal.Type().Kind != ir.KindPointer {
				return fmt.Errorf("cannot assign to non-pointer dereference")
			}
			l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: rhsVal, StoreDst: ptrVal, StoreAlign: alignOfIRType(rhsVal.Type())})
			continue
		}
		lhs, ok := st.Lhs[i].(*ast.Ident)
		if !ok {
			return fmt.Errorf("unsupported lhs type %T", st.Lhs[i])
		}
		if lhs.Name == "_" {
			// Вычисляем RHS только ради побочных эффектов.
			if _, err := l.lowerExpr(st.Rhs[i]); err != nil {
				return err
			}
			continue
		}
		rhsVal, err := l.lowerExpr(st.Rhs[i])
		if err != nil {
			return err
		}
		if st.Tok == token.DEFINE {
			addr := l.ensureLocal(lhs.Name, rhsVal.Type())
			l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: rhsVal, StoreDst: addr, StoreAlign: alignOfIRType(rhsVal.Type())})
		} else {
			addr := l.getLocal(lhs.Name)
			if addr == nil {
				return fmt.Errorf("unknown var %s", lhs.Name)
			}
			l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: rhsVal, StoreDst: addr, StoreAlign: alignOfIRType(rhsVal.Type())})
		}
	}
	return nil
}

// IfStmt = "if" Expr Block [ "else" Block ]
func (l *lowerer) lowerIf(st *ast.IfStmt) error {
	if st.Init != nil {
		return fmt.Errorf("unsupported: if init")
	}
	cond, err := l.lowerExpr(st.Cond)
	if err != nil {
		return err
	}
	thenBB := l.newBlock(l.newBlockName("then"))
	elseBB := l.newBlock(l.newBlockName("else"))
	mergeBB := l.newBlock(l.newBlockName("endif"))

	l.block.Terminator = &ir.Instruction{Kind: ir.InstrCondBr, CondCond: cond, CondTrue: thenBB.Name, CondFalse: elseBB.Name}

	// ветка then
	l.block = thenBB
	if err := l.lowerStmt(st.Body); err != nil {
		return err
	}
	if l.block.Terminator == nil {
		l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: mergeBB.Name}
	}

	// ветка else
	l.block = elseBB
	if st.Else != nil {
		if err := l.lowerStmt(st.Else); err != nil {
			return err
		}
	}
	if l.block.Terminator == nil {
		l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: mergeBB.Name}
	}

	// продолжаем в merge-блоке
	l.block = mergeBB
	return nil
}

// ForStmt = "for" ( [ SimpleStmt ] ";" [ Expr ] ";" [ SimpleStmt ] | Expr ) Block
// SimpleStmt = AssignStmt | ExprStmt
func (l *lowerer) lowerFor(st *ast.ForStmt) error {
	// Поддерживаем формы: for cond { ... } и for init; cond; post { ... }.
	condBB := l.newBlock(l.newBlockName("for.cond"))
	bodyBB := l.newBlock(l.newBlockName("for.body"))
	postBB := l.newBlock(l.newBlockName("for.post"))
	endBB := l.newBlock(l.newBlockName("for.end"))

	// кладем loop-таргеты в стек
	l.breakTargets = append(l.breakTargets, endBB.Name)
	if st.Post != nil {
		l.continueTargets = append(l.continueTargets, postBB.Name)
	} else {
		l.continueTargets = append(l.continueTargets, condBB.Name)
	}

	// инициализация
	if st.Init != nil {
		if err := l.lowerStmt(st.Init); err != nil {
			return err
		}
	}

	l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: condBB.Name}

	// условие
	l.block = condBB
	var condVal *ir.Value
	if st.Cond != nil {
		v, err := l.lowerExpr(st.Cond)
		if err != nil {
			return err
		}
		condVal = v
	} else {
		condVal = ir.NewConstant("1", ir.I1)
	}
	l.block.Terminator = &ir.Instruction{Kind: ir.InstrCondBr, CondCond: condVal, CondTrue: bodyBB.Name, CondFalse: endBB.Name}

	// тело
	l.block = bodyBB
	if err := l.lowerBlock(st.Body); err != nil {
		return err
	}
	if l.block.Terminator == nil {
		if st.Post != nil {
			l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: postBB.Name}
		} else {
			l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: condBB.Name}
		}
	}

	if st.Post != nil {
		l.block = postBB
		if err := l.lowerStmt(st.Post); err != nil {
			return err
		}
		if l.block.Terminator == nil {
			l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: condBB.Name}
		}
	}

	// продолжение после цикла
	l.block = endBB
	// снимаем loop-таргеты со стека
	if len(l.breakTargets) > 0 {
		l.breakTargets = l.breakTargets[:len(l.breakTargets)-1]
	}
	if len(l.continueTargets) > 0 {
		l.continueTargets = l.continueTargets[:len(l.continueTargets)-1]
	}
	return nil
}

// Expr = BinaryExpr | UnaryExpr | PrimaryExpr
// PrimaryExpr = ident | literal | "(" Expr ")" | Selector | Index | Slice | Call
func (l *lowerer) lowerExpr(e ast.Expr) (*ir.Value, error) {
	switch ex := e.(type) {
	case *ast.BasicLit:
		return l.lowerBasicLit(ex)
	case *ast.ParenExpr:
		return l.lowerExpr(ex.X)
	case *ast.UnaryExpr:
		return l.lowerUnary(ex)
	case *ast.Ident:
		if v := l.getLocal(ex.Name); v != nil {
			loaded := l.newTemp(elemTypeFromPtr(v.Type()))
			l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: v, LoadAlign: alignOfIRType(elemTypeFromPtr(v.Type()))})
			return loaded, nil
		}
		// если локального нет, берем глобальную
		obj := l.prog.TypesInfo.Uses[ex]
		if obj == nil {
			return nil, fmt.Errorf("unknown ident %s", ex.Name)
		}
		ty, err := goTypeToIR(obj.Type())
		if err != nil {
			return nil, err
		}
		gname := ex.Name
		if l.prefix != "" {
			gname = l.prefix + gname
		}
		ptr := ir.NewConstant("@"+gname, ir.PtrTo(ty))
		loaded := l.newTemp(ty)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: ptr, LoadAlign: alignOfIRType(ty)})
		return loaded, nil
	case *ast.IndexExpr:
		return l.lowerIndex(ex)
	case *ast.CompositeLit:
		return l.lowerCompositeLit(ex)
	case *ast.BinaryExpr:
		return l.lowerBinary(ex)
	case *ast.CallExpr:
		return l.lowerCall(ex)
	case *ast.SelectorExpr:
		return l.lowerSelector(ex)
	case *ast.StarExpr:
		return l.lowerStar(ex)
	case *ast.SliceExpr:
		return l.lowerSliceExpr(ex)
	default:
		return nil, fmt.Errorf("unsupported expr %T", e)
	}
}

// Index = PrimaryExpr "[" Expr "]"
func (l *lowerer) lowerIndex(idx *ast.IndexExpr) (*ir.Value, error) {
	base, err := l.lowerExpr(idx.X)
	if err != nil {
		return nil, err
	}
	indexVal, err := l.lowerExpr(idx.Index)
	if err != nil {
		return nil, err
	}
	if _, ok := l.prog.TypesInfo.TypeOf(idx.X).Underlying().(*types.Map); ok {
		mapVal := base
		valGo := l.prog.TypesInfo.TypeOf(idx.X).Underlying().(*types.Map).Elem()
		valTy, err := goTypeToIR(valGo)
		if err != nil {
			return nil, err
		}
		keyTy, err := goTypeToIR(l.prog.TypesInfo.TypeOf(idx.Index))
		if err != nil {
			return nil, err
		}
		keyBuf := l.newTemp(ir.PtrTo(keyTy))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: keyBuf, AllocaType: keyTy, AllocaAlign: alignOfIRType(keyTy)})
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: indexVal, StoreDst: keyBuf, StoreAlign: alignOfIRType(keyTy)})
		valBuf := l.newTemp(ir.PtrTo(valTy))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: valBuf, AllocaType: valTy, AllocaAlign: alignOfIRType(valTy)})
		keyI8 := l.newTemp(ir.PtrI8)
		valI8 := l.newTemp(ir.PtrI8)
		l.block.Append(ir.Instruction{Kind: ir.InstrBitcast, Dest: keyI8, BitcastSrc: keyBuf, BitcastTarget: ir.PtrI8})
		l.block.Append(ir.Instruction{Kind: ir.InstrBitcast, Dest: valI8, BitcastSrc: valBuf, BitcastTarget: ir.PtrI8})
		// результат ok игнорируем
		l.block.Append(ir.Instruction{Kind: ir.InstrCall, CallName: "gominic_map_get", CallArgs: []*ir.Value{mapVal, keyI8, valI8}, CallRet: ir.I1})
		loaded := l.newTemp(valTy)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: valBuf, LoadAlign: alignOfIRType(valTy)})
		return loaded, nil
	}
	baseTy := base.Type()
	switch baseTy.Kind {
	case ir.KindArray:
		elemTy := baseTy.Elem
		arrPtr := l.valuePtr(base, baseTy)
		elemPtr := l.newTemp(ir.PtrTo(elemTy))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          elemPtr,
			GepSrc:        arrPtr,
			GepPointee:    baseTy,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), indexVal},
			GepResultType: ir.PtrTo(elemTy),
		})
		loaded := l.newTemp(elemTy)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: elemPtr, LoadAlign: alignOfIRType(elemTy)})
		return loaded, nil
	case ir.KindPointer:
		if baseTy.Elem != nil && baseTy.Elem.Kind == ir.KindArray {
			elemTy := baseTy.Elem.Elem
			elemPtr := l.newTemp(ir.PtrTo(elemTy))
			l.block.Append(ir.Instruction{
				Kind:          ir.InstrGEP,
				Dest:          elemPtr,
				GepSrc:        base,
				GepPointee:    baseTy.Elem,
				GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), indexVal},
				GepResultType: ir.PtrTo(elemTy),
			})
			loaded := l.newTemp(elemTy)
			l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: elemPtr, LoadAlign: alignOfIRType(elemTy)})
			return loaded, nil
		}
		pos := l.prog.Fset.Position(idx.Lbrack)
		return nil, fmt.Errorf("unsupported index on pointer %s at %s", baseTy.String(), pos.String())
	case ir.KindSlice:
		// Представление среза: { data *T, len i64, cap i64 }.
		slicePtr := l.valuePtr(base, baseTy)
		dataPtr := l.newTemp(ir.PtrTo(baseTy.Elem))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          dataPtr,
			GepSrc:        slicePtr,
			GepPointee:    baseTy,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("0", ir.I32)},
			GepResultType: ir.PtrTo(baseTy.Elem),
		})
		elemPtr := l.newTemp(ir.PtrTo(baseTy.Elem))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          elemPtr,
			GepSrc:        dataPtr,
			GepPointee:    baseTy.Elem,
			GepIndices:    []*ir.Value{indexVal},
			GepResultType: ir.PtrTo(baseTy.Elem),
		})
		loaded := l.newTemp(baseTy.Elem)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: elemPtr, LoadAlign: alignOfIRType(baseTy.Elem)})
		return loaded, nil
	case ir.KindString:
		// рассматриваем строку как { i8*, i64 }
		strPtr := l.valuePtr(base, baseTy)
		dataPtr := l.newTemp(ir.PtrTo(ir.I8))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          dataPtr,
			GepSrc:        strPtr,
			GepPointee:    baseTy,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("0", ir.I32)},
			GepResultType: ir.PtrTo(ir.I8),
		})
		elemPtr := l.newTemp(ir.PtrTo(ir.I8))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          elemPtr,
			GepSrc:        dataPtr,
			GepPointee:    ir.I8,
			GepIndices:    []*ir.Value{indexVal},
			GepResultType: ir.PtrTo(ir.I8),
		})
		byteVal := l.newTemp(ir.I8)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: byteVal, LoadSrc: elemPtr, LoadAlign: 1})
		return byteVal, nil
	default:
		return nil, fmt.Errorf("unsupported index on type %v", baseTy.Kind)
	}
}

// Slice = PrimaryExpr "[" [ Expr ] ":" [ Expr ] "]"
func (l *lowerer) lowerSliceExpr(slice *ast.SliceExpr) (*ir.Value, error) {
	baseVal, err := l.lowerExpr(slice.X)
	if err != nil {
		return nil, err
	}
	baseTy := baseVal.Type()
	if baseTy == nil {
		return nil, fmt.Errorf("slice expr: missing base type")
	}

	// Берём ожидаемый тип результата из Go type system.
	expectedTy, err := goTypeToIR(l.prog.TypesInfo.TypeOf(slice))
	if err != nil {
		// Если тип не удалось получить, используем базовый.
		expectedTy = nil
	}

	// Берём low-индекс (по умолчанию 0).
	var lowVal *ir.Value
	if slice.Low != nil {
		lowVal, err = l.lowerExpr(slice.Low)
		if err != nil {
			return nil, err
		}
		// Убеждаемся, что lowVal имеет тип i64.
		if lowVal.Type() == nil || lowVal.Type().Basic != "i64" {
			lowVal, err = l.lowerConvert(lowVal, ir.I64)
			if err != nil {
				return nil, err
			}
		}
	} else {
		lowVal = ir.NewConstant("0", ir.I64)
	}

	// Обрабатываем строки.
	if baseTy.Kind == ir.KindString {
		ptr, lenVal := l.materializeStringParts(baseVal)

		// Берём high-индекс (по умолчанию длина).
		var highVal *ir.Value
		if slice.High != nil {
			highVal, err = l.lowerExpr(slice.High)
			if err != nil {
				return nil, err
			}
			// Убеждаемся, что highVal имеет тип i64.
			if highVal.Type() == nil || highVal.Type().Basic != "i64" {
				highVal, err = l.lowerConvert(highVal, ir.I64)
				if err != nil {
					return nil, err
				}
			}
		} else {
			highVal = lenVal
		}

		// Вычисляем новый указатель: ptr + low через GEP.
		newPtr := l.newTemp(ir.PtrTo(ir.I8))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          newPtr,
			GepSrc:        ptr,
			GepPointee:    ir.I8,
			GepIndices:    []*ir.Value{lowVal},
			GepResultType: ir.PtrTo(ir.I8),
		})

		// Вычисляем новую длину: high - low.
		newLen := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{
			Kind:  ir.InstrBinOp,
			Dest:  newLen,
			BinOp: ir.Sub,
			X:     highVal,
			Y:     lowVal,
		})

		// Собираем строку или срез в зависимости от ожидаемого типа.
		if expectedTy != nil && expectedTy.Kind == ir.KindSlice && expectedTy.Elem == ir.I8 {
			// Ожидается []byte, возвращаем срез.
			return l.buildSliceValue(ir.I8, newPtr, newLen, newLen), nil
		}
		// Ожидается string, возвращаем строку.
		return l.buildStringValue(newPtr, newLen), nil
	}

	// Обрабатываем срезы.
	if baseTy.Kind == ir.KindSlice {
		elemTy := baseTy.Elem
		slicePtr := l.valuePtr(baseVal, baseTy)

		// читаем указатель на данные
		dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrTo(elemTy)))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          dataPtrPtr,
			GepSrc:        slicePtr,
			GepPointee:    baseTy,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("0", ir.I32)},
			GepResultType: ir.PtrTo(ir.PtrTo(elemTy)),
		})
		dataPtr := l.newTemp(ir.PtrTo(elemTy))
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: dataPtr, LoadSrc: dataPtrPtr, LoadAlign: 8})

		// читаем длину
		lenPtr := l.newTemp(ir.PtrTo(ir.I64))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          lenPtr,
			GepSrc:        slicePtr,
			GepPointee:    baseTy,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("1", ir.I32)},
			GepResultType: ir.PtrTo(ir.I64),
		})
		lenVal := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: lenVal, LoadSrc: lenPtr, LoadAlign: 8})

		// читаем емкость
		capPtr := l.newTemp(ir.PtrTo(ir.I64))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          capPtr,
			GepSrc:        slicePtr,
			GepPointee:    baseTy,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("2", ir.I32)},
			GepResultType: ir.PtrTo(ir.I64),
		})
		capVal := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: capVal, LoadSrc: capPtr, LoadAlign: 8})

		// Берём high-индекс (по умолчанию длина).
		var highVal *ir.Value
		if slice.High != nil {
			highVal, err = l.lowerExpr(slice.High)
			if err != nil {
				return nil, err
			}
			// Убеждаемся, что highVal имеет тип i64.
			if highVal.Type() == nil || highVal.Type().Basic != "i64" {
				highVal, err = l.lowerConvert(highVal, ir.I64)
				if err != nil {
					return nil, err
				}
			}
		} else {
			highVal = lenVal
		}

		// Вычисляем новый data: data + low через GEP.
		newDataPtr := l.newTemp(ir.PtrTo(elemTy))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          newDataPtr,
			GepSrc:        dataPtr,
			GepPointee:    elemTy,
			GepIndices:    []*ir.Value{lowVal},
			GepResultType: ir.PtrTo(elemTy),
		})

		// Вычисляем новую длину: high - low.
		newLen := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{
			Kind:  ir.InstrBinOp,
			Dest:  newLen,
			BinOp: ir.Sub,
			X:     highVal,
			Y:     lowVal,
		})

		// Вычисляем новую ёмкость: cap - low.
		newCap := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{
			Kind:  ir.InstrBinOp,
			Dest:  newCap,
			BinOp: ir.Sub,
			X:     capVal,
			Y:     lowVal,
		})

		// Собираем новый срез.
		return l.buildSliceValue(elemTy, newDataPtr, newLen, newCap), nil
	}

	return nil, fmt.Errorf("slice expr on unsupported type %v", baseTy.Kind)
}

func (l *lowerer) buildStringValue(ptr *ir.Value, lenVal *ir.Value) *ir.Value {
	strTy := ir.String()
	addr := l.newTemp(ir.PtrTo(strTy))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: strTy, AllocaAlign: alignOfIRType(strTy)})
	// указатель
	ptrField := l.newTemp(ir.PtrTo(ir.PtrTo(ir.I8)))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          ptrField,
		GepSrc:        addr,
		GepPointee:    strTy,
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("0", ir.I32)},
		GepResultType: ir.PtrTo(ir.PtrTo(ir.I8)),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: ptr, StoreDst: ptrField, StoreAlign: 8})
	// длина
	lenField := l.newTemp(ir.PtrTo(ir.I64))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          lenField,
		GepSrc:        addr,
		GepPointee:    strTy,
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("1", ir.I32)},
		GepResultType: ir.PtrTo(ir.I64),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: lenVal, StoreDst: lenField, StoreAlign: 8})
	// загружаем результат
	res := l.newTemp(strTy)
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: res, LoadSrc: addr, LoadAlign: alignOfIRType(strTy)})
	return res
}

// valuePtr: если значение не указатель, выделяем память и сохраняем его.
// Нужно для операций, требующих указатель (например, передача в runtime функции)
func (l *lowerer) valuePtr(v *ir.Value, ty *ir.TypeDesc) *ir.Value {
	if v.Type() != nil && v.Type().Kind == ir.KindPointer {
		return v
	}
	addr := l.newTemp(ir.PtrTo(ty))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: ty, AllocaAlign: alignOfIRType(ty)})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: v, StoreDst: addr, StoreAlign: alignOfIRType(ty)})
	return addr
}

// buildSliceValue: собирает значение среза { i8*, i64, i64 }.
func (l *lowerer) buildSliceValue(elemTy *ir.TypeDesc, dataPtr *ir.Value, lenVal *ir.Value, capVal *ir.Value) *ir.Value {
	sliceTy := ir.Slice(elemTy)
	addr := l.newTemp(ir.PtrTo(sliceTy))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: sliceTy, AllocaAlign: alignOfIRType(sliceTy)})
	// данные
	dataField := l.newTemp(ir.PtrTo(ir.PtrTo(elemTy)))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          dataField,
		GepSrc:        addr,
		GepPointee:    sliceTy,
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("0", ir.I32)},
		GepResultType: ir.PtrTo(ir.PtrTo(elemTy)),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: dataPtr, StoreDst: dataField, StoreAlign: 8})
	// длина
	lenField := l.newTemp(ir.PtrTo(ir.I64))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          lenField,
		GepSrc:        addr,
		GepPointee:    sliceTy,
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("1", ir.I32)},
		GepResultType: ir.PtrTo(ir.I64),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: lenVal, StoreDst: lenField, StoreAlign: 8})
	// емкость
	capField := l.newTemp(ir.PtrTo(ir.I64))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          capField,
		GepSrc:        addr,
		GepPointee:    sliceTy,
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("2", ir.I32)},
		GepResultType: ir.PtrTo(ir.I64),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: capVal, StoreDst: capField, StoreAlign: 8})
	// загружаем результат
	res := l.newTemp(sliceTy)
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: res, LoadSrc: addr, LoadAlign: alignOfIRType(sliceTy)})
	return res
}

func (l *lowerer) lowerCompositeLit(cl *ast.CompositeLit) (*ir.Value, error) {
	gt := l.prog.TypesInfo.TypeOf(cl)
	if gt == nil {
		return nil, fmt.Errorf("no type for composite literal")
	}
	irTy, err := goTypeToIR(gt)
	if err != nil {
		return nil, err
	}
	switch tt := gt.Underlying().(type) {
	case *types.Array:
		elemTy := irTy.Elem
		arrPtr := l.newTemp(ir.PtrTo(irTy))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: arrPtr, AllocaType: irTy, AllocaAlign: alignOfIRType(irTy)})
		for i, elt := range cl.Elts {
			idx := i
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				if bl, ok := kv.Key.(*ast.BasicLit); ok && bl.Kind == token.INT {
					if val, err := strconv.Atoi(bl.Value); err == nil {
						idx = val
					}
				}
				elt = kv.Value
			}
			val, err := l.lowerExpr(elt)
			if err != nil {
				return nil, err
			}
			ptr := l.newTemp(ir.PtrTo(elemTy))
			l.block.Append(ir.Instruction{
				Kind:          ir.InstrGEP,
				Dest:          ptr,
				GepSrc:        arrPtr,
				GepPointee:    irTy,
				GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant(strconv.Itoa(idx), ir.I32)},
				GepResultType: ir.PtrTo(elemTy),
			})
			l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: val, StoreDst: ptr, StoreAlign: alignOfIRType(elemTy)})
		}
		out := l.newTemp(irTy)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: out, LoadSrc: arrPtr, LoadAlign: alignOfIRType(irTy)})
		return out, nil
	case *types.Slice:
		elemTy := irTy.Elem
		ln := len(cl.Elts)
		lenConst := ir.NewConstant(strconv.Itoa(ln), ir.I64)
		elemSize := ir.NewConstant(strconv.Itoa(int(sizeOfIRType(elemTy))), ir.I64)
		ptr := l.newTemp(ir.PtrTo(elemTy))
		l.block.Append(ir.Instruction{Kind: ir.InstrCall, Dest: ptr, CallName: "gominic_makeSlice", CallArgs: []*ir.Value{lenConst, lenConst, elemSize}, CallRet: ir.PtrTo(elemTy)})
		for i, elt := range cl.Elts {
			val, err := l.lowerExpr(elt)
			if err != nil {
				return nil, err
			}
			elemPtr := l.newTemp(ir.PtrTo(elemTy))
			l.block.Append(ir.Instruction{
				Kind:          ir.InstrGEP,
				Dest:          elemPtr,
				GepSrc:        ptr,
				GepPointee:    elemTy,
				GepIndices:    []*ir.Value{ir.NewConstant(strconv.Itoa(i), ir.I64)},
				GepResultType: ir.PtrTo(elemTy),
			})
			l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: val, StoreDst: elemPtr, StoreAlign: alignOfIRType(elemTy)})
		}
		return l.buildSliceValue(elemTy, ptr, lenConst, lenConst), nil
	case *types.Struct:
		// Если рекурсивный тип заменён на PtrI8, поля читать нельзя.
		// В этом случае возвращаем неинициализированный указатель.
		if irTy.Kind == ir.KindPointer && irTy.Elem != nil && irTy.Elem.Basic == "i8" {
			// Это указательный тип (обычно PtrI8 после обработки рекурсии).
			// Создаём неинициализированное указательное значение.
			ptr := l.newTemp(irTy)
			l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: ptr, AllocaType: ir.PtrI8, AllocaAlign: 8})
			loaded := l.newTemp(irTy)
			l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: ptr, LoadAlign: 8})
			return loaded, nil
		}
		if irTy.Kind != ir.KindStruct {
			return nil, fmt.Errorf("expected struct type, got %v", irTy)
		}
		structPtr := l.newTemp(ir.PtrTo(irTy))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: structPtr, AllocaType: irTy, AllocaAlign: alignOfIRType(irTy)})
		fieldIndex := func(name string) int {
			for i := 0; i < tt.NumFields(); i++ {
				if tt.Field(i).Name() == name {
					return i
				}
			}
			return -1
		}
		if len(irTy.Fields) == 0 {
			return nil, fmt.Errorf("struct type has no fields")
		}
		for i, elt := range cl.Elts {
			idx := i
			expr := elt
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				if id, ok := kv.Key.(*ast.Ident); ok {
					if fi := fieldIndex(id.Name); fi >= 0 {
						idx = fi
					}
				}
				expr = kv.Value
			}
			if idx < 0 || idx >= len(irTy.Fields) {
				return nil, fmt.Errorf("field index %d out of range (struct has %d fields)", idx, len(irTy.Fields))
			}
			val, err := l.lowerExpr(expr)
			if err != nil {
				return nil, err
			}
			fieldPtr := l.newTemp(ir.PtrTo(irTy.Fields[idx]))
			l.block.Append(ir.Instruction{
				Kind:          ir.InstrGEP,
				Dest:          fieldPtr,
				GepSrc:        structPtr,
				GepPointee:    irTy,
				GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant(strconv.Itoa(idx), ir.I32)},
				GepResultType: ir.PtrTo(irTy.Fields[idx]),
			})
			l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: val, StoreDst: fieldPtr, StoreAlign: alignOfIRType(irTy.Fields[idx])})
		}
		out := l.newTemp(irTy)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: out, LoadSrc: structPtr, LoadAlign: alignOfIRType(irTy)})
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported composite literal")
	}
}

// Selector = PrimaryExpr "." ident
func (l *lowerer) lowerSelector(sel *ast.SelectorExpr) (*ir.Value, error) {
	// Сначала проверяем доступ к константе пакета (например, ir.KindBasic).
	// Для констант в первую очередь смотрим TypesInfo.Types.
	if tv, ok := l.prog.TypesInfo.Types[sel]; ok {
		if tv.IsValue() && tv.Value != nil {
			// Это константное значение.
			ty, err := goTypeToIR(tv.Type)
			if err != nil {
				return nil, err
			}
			valStr := tv.Value.String()
			return ir.NewConstant(valStr, ty), nil
		}
		// Также проверяем вариант с типовой константой (iota).
		if tv.IsType() {
			// Здесь может быть тип, а не константа.
		}
	}

	// Сначала проверяем доступ к полю структуры, потом Uses.
	// Иначе поле можно ошибочно принять за глобальную переменную.
	selection := l.prog.TypesInfo.Selections[sel]
	if selection != nil && selection.Kind() == types.FieldVal {
		// Это доступ к полю, обрабатываем ниже.
	} else {
		// Это не поле: проверяем константу или переменную пакета.
		if obj, ok := l.prog.TypesInfo.Uses[sel.Sel]; ok {
			if constVal, ok := obj.(*types.Const); ok {
				// Это константа.
				ty, err := goTypeToIR(constVal.Type())
				if err != nil {
					return nil, err
				}
				// Преобразуем значение константы в строковый вид.
				valStr := constVal.Val().String()
				return ir.NewConstant(valStr, ty), nil
			}
			// Обрабатываем переменные пакета (например, ir.Void, ir.I64).
			if varVal, ok := obj.(*types.Var); ok {
				// Это переменная пакета, считаем её глобальной.
				ty, err := goTypeToIR(varVal.Type())
				if err != nil {
					return nil, err
				}
				gname := varVal.Name()
				// Проверяем, что селектор действительно пакетный (через sel.X).
				if pkgIdent, ok := sel.X.(*ast.Ident); ok {
					if pkgObj, ok := l.prog.TypesInfo.Uses[pkgIdent]; ok {
						if pkg, ok := pkgObj.(*types.PkgName); ok {
							pkgName := pkg.Imported().Name()
							// Для переменных stdlib используем полное имя.
							if pkgName != "" && (pkgName == "os" || pkgName == "fmt" || pkgName == "io") {
								gname = pkgName + "." + gname
							} else if l.prefix != "" {
								gname = l.prefix + gname
							}
						} else if l.prefix != "" {
							gname = l.prefix + gname
						}
					} else if l.prefix != "" {
						gname = l.prefix + gname
					}
				} else if l.prefix != "" {
					gname = l.prefix + gname
				}
				ptr := ir.NewConstant("@"+gname, ir.PtrTo(ty))
				loaded := l.newTemp(ty)
				l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: ptr, LoadAlign: alignOfIRType(ty)})
				return loaded, nil
			}
		}
	}

	// Пробуем обработать как константу пакета через lookup.
	if pkgIdent, ok := sel.X.(*ast.Ident); ok {
		obj := l.prog.TypesInfo.Uses[pkgIdent]
		if obj != nil {
			if pkg, ok := obj.(*types.PkgName); ok {
				// Это пакетный селектор, ищем константу в scope пакета.
				pkgScope := pkg.Imported().Scope()
				pkgObj := pkgScope.Lookup(sel.Sel.Name)
				if pkgObj != nil {
					if constVal, ok := pkgObj.(*types.Const); ok {
						// Берем значение константы.
						ty, err := goTypeToIR(constVal.Type())
						if err != nil {
							return nil, err
						}
						// Преобразуем значение константы в строковый вид.
						valStr := constVal.Val().String()
						return ir.NewConstant(valStr, ty), nil
					}
					if varVal, ok := pkgObj.(*types.Var); ok {
						// Это переменная пакета, считаем её глобальной.
						ty, err := goTypeToIR(varVal.Type())
						if err != nil {
							return nil, err
						}
						// Для stdlib-переменных (типа os.Stderr) добавляем префикс пакета.
						pkgName := pkg.Imported().Name()
						gname := varVal.Name()
						// Для stdlib используем полное имя.
						if pkgName != "" && (pkgName == "os" || pkgName == "fmt" || pkgName == "io") {
							gname = pkgName + "." + gname
						} else if l.prefix != "" {
							gname = l.prefix + gname
						}
						ptr := ir.NewConstant("@"+gname, ir.PtrTo(ty))
						loaded := l.newTemp(ty)
						l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: ptr, LoadAlign: alignOfIRType(ty)})
						return loaded, nil
					}
				}
			}
		}
	}

	// Если выше не сработало, проверяем доступ к полю struct.
	if selection == nil {
		// Последняя попытка: ищем sel.Sel в Uses.
		if obj, ok := l.prog.TypesInfo.Uses[sel.Sel]; ok {
			if constVal, ok := obj.(*types.Const); ok {
				ty, err := goTypeToIR(constVal.Type())
				if err != nil {
					return nil, err
				}
				valStr := constVal.Val().String()
				return ir.NewConstant(valStr, ty), nil
			}
		}
		return nil, fmt.Errorf("unsupported selector (X=%T, Sel=%v)", sel.X, sel.Sel)
	}
	if selection.Kind() != types.FieldVal {
		return nil, fmt.Errorf("unsupported selector (kind=%v)", selection.Kind())
	}
	indices := selection.Index()
	if len(indices) == 0 {
		return nil, fmt.Errorf("selector has no field index")
	}
	fieldIdx := indices[0]
	baseVal, err := l.lowerExpr(sel.X)
	if err != nil {
		return nil, err
	}

	// Получаем IR-тип из Go-типа (включая именованные типы).
	baseGoType := l.prog.TypesInfo.TypeOf(sel.X)
	if baseGoType == nil {
		return nil, fmt.Errorf("no type info for selector base")
	}

	// Получаем IR-тип базы (указатель или значение), сначала проверяем рекурсивные случаи.
	baseType, err := goTypeToIR(baseGoType)
	if err != nil {
		return nil, err
	}

	// Извлекаем IR-тип struct из baseType.
	var structIRType *ir.TypeDesc
	if baseType.Kind == ir.KindPointer {
		structIRType = baseType.Elem
	} else {
		structIRType = baseType
	}

	// Сначала проверяем рекурсию, и только потом пытаемся взять Go-struct.
	if structIRType.Kind != ir.KindStruct {
		// Возможно, это рекурсивный тип: сначала смотрим исходный именованный тип.
		// Берем тип элемента из baseGoType до вызова Underlying().
		var elemType types.Type
		if ptrType, ok := baseGoType.(*types.Pointer); ok {
			elemType = ptrType.Elem()
		} else {
			elemType = baseGoType
		}
		// Проверяем известные рекурсивные типы (например, Global) со спец-обработкой.
		if named, ok := elemType.(*types.Named); ok {
			typeName := ""
			if named.Obj() != nil {
				typeName = named.Obj().Name()
			}
			fieldName := sel.Sel.Name
			// Поля Global обрабатываем отдельно (ir/main).
			// Сверяемся по имени типа, независимо от пакета.
			if typeName == "Global" {
				if fieldName == "Align" {
					return ir.NewConstant("0", ir.I64), nil
				} else if fieldName == "Private" {
					return ir.NewConstant("0", ir.I1), nil
				} else if fieldName == "Name" {
					return ir.NewConstant(`{ i8* null, i64 0 }`, ir.String()), nil
				} else if fieldName == "Type" {
					return ir.NewConstant("null", ir.PtrI8), nil
				} else if fieldName == "Value" {
					return ir.NewConstant(`{ i8* null, i64 0 }`, ir.String()), nil
				}
			}
		}
		// Если это не известный рекурсивный тип, пробуем получить Go-struct.
		// Получаем фактический struct-тип (с учетом именованных и указателей).
		underlying := baseGoType.Underlying()
		var structGoType types.Type
		if ptrType, ok := underlying.(*types.Pointer); ok {
			// Это указатель: берем тип элемента.
			structGoType = ptrType.Elem().Underlying()
		} else {
			// Это значение: берем underlying (для именованных тоже).
			structGoType = underlying
		}

		// Пробуем конвертировать Go-struct в IR.
		structIRTypeFromGo, err := goTypeToIR(structGoType)
		if err == nil && structIRTypeFromGo.Kind == ir.KindStruct {
			// В IR это действительно struct, используем его.
			structIRType = structIRTypeFromGo
		} else {
			// Это рекурсивный тип, который заменили на PtrI8 или аналог.
			// Возвращаем указатель, поля у него не читаем.
			return ir.NewConstant("null", ir.PtrI8), nil
		}
	}

	// Теперь получаем Go-struct для поиска индекса поля.
	underlying := baseGoType.Underlying()
	var structGoType types.Type
	if ptrType, ok := underlying.(*types.Pointer); ok {
		// Это указатель: берем тип элемента.
		structGoType = ptrType.Elem().Underlying()
	} else {
		// Это значение: берем underlying (для именованных тоже).
		structGoType = underlying
	}

	// Проверяем, что это действительно struct.
	_, ok := structGoType.(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("selector base not struct (got %T)", structGoType)
	}

	// Выбираем указатель, с которым будем работать.
	var ptr *ir.Value
	if baseVal.Type() != nil && baseVal.Type().Kind == ir.KindPointer {
		// baseVal уже указатель, используем как есть.
		ptr = baseVal
	} else if baseType.Kind == ir.KindPointer {
		// Go-тип — указатель, но baseVal значение (редко, но обрабатываем).
		ptr = baseVal
	} else {
		// Для значения нужно получить адрес.
		addr := l.newTemp(ir.PtrTo(structIRType))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: structIRType, AllocaAlign: alignOfIRType(structIRType)})
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: baseVal, StoreDst: addr, StoreAlign: alignOfIRType(structIRType)})
		ptr = addr
	}
	// Здесь ожидаем structIRType == struct, но проверим еще раз.
	if structIRType.Kind != ir.KindStruct {
		// Ситуация нештатная, но обрабатываем как рекурсивный тип.
		// Берем тип элемента из baseGoType.
		var elemType types.Type
		if ptrType, ok := baseGoType.(*types.Pointer); ok {
			elemType = ptrType.Elem()
		} else {
			elemType = baseGoType
		}
		if named, ok := elemType.(*types.Named); ok {
			typeName := ""
			if named.Obj() != nil {
				typeName = named.Obj().Name()
			}
			fieldName := sel.Sel.Name
			// Отдельная обработка полей Global.
			if typeName == "Global" {
				if fieldName == "Align" {
					return ir.NewConstant("0", ir.I64), nil
				} else if fieldName == "Private" {
					return ir.NewConstant("0", ir.I1), nil
				} else if fieldName == "Name" {
					return ir.NewConstant(`{ i8* null, i64 0 }`, ir.String()), nil
				} else if fieldName == "Type" {
					return ir.NewConstant("null", ir.PtrI8), nil
				} else if fieldName == "Value" {
					return ir.NewConstant(`{ i8* null, i64 0 }`, ir.String()), nil
				}
			}
		}
		return nil, fmt.Errorf("cannot access field %s on non-struct type %v", sel.Sel.Name, structIRType)
	}
	fieldTy, err := goTypeToIR(selection.Type())
	if err != nil {
		return nil, err
	}
	fieldPtr := l.newTemp(ir.PtrTo(fieldTy))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          fieldPtr,
		GepSrc:        ptr,
		GepPointee:    structIRType,
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant(strconv.Itoa(fieldIdx), ir.I32)},
		GepResultType: ir.PtrTo(fieldTy),
	})
	loaded := l.newTemp(fieldTy)
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: fieldPtr, LoadAlign: alignOfIRType(fieldTy)})
	return loaded, nil
}

// literal = int_lit | float_lit | string_lit | boolean_lit
// boolean_lit = "true" | "false"
func (l *lowerer) lowerBasicLit(lit *ast.BasicLit) (*ir.Value, error) {
	switch lit.Kind {
	case token.INT:
		return ir.NewConstant(lit.Value, ir.I64), nil
	case token.FLOAT:
		return ir.NewConstant(lit.Value, ir.F64), nil
	case token.STRING:
		// убираем кавычки
		s, err := strconv.Unquote(lit.Value)
		if err != nil {
			return nil, err
		}
		return l.stringConstant(s), nil
	default:
		return nil, fmt.Errorf("unsupported literal %s", lit.Value)
	}
}

// BinaryExpr = Expr BinOp Expr
// BinOp = "+" | "-" | "*" | "/" | "%" | "==" | "!=" | "<" | "<=" | ">" | ">=" | "&&" | "||" | "&" | "|" | "<<" | ">>"
func (l *lowerer) lowerBinary(b *ast.BinaryExpr) (*ir.Value, error) {
	op := b.Op
	// && и || требуют короткого замыкания: правую часть считаем только при необходимости.
	if op == token.LAND || op == token.LOR {
		x, err := l.lowerExpr(b.X)
		if err != nil {
			return nil, err
		}
		return l.lowerLogical(op, x, b.Y)
	}

	x, err := l.lowerExpr(b.X)
	if err != nil {
		return nil, err
	}
	y, err := l.lowerExpr(b.Y)
	if err != nil {
		return nil, err
	}

	// Конкатенация строк: выделяем память для всех строк и вызываем runtime функцию
	if op == token.ADD && x.Type() != nil && x.Type().Kind == ir.KindString && y.Type() != nil && y.Type().Kind == ir.KindString {
		outPtr := l.newTemp(ir.PtrTo(ir.String()))
		xPtr := l.newTemp(ir.PtrTo(ir.String()))
		yPtr := l.newTemp(ir.PtrTo(ir.String()))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: outPtr, AllocaType: ir.String(), AllocaAlign: alignOfIRType(ir.String())})
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: xPtr, AllocaType: ir.String(), AllocaAlign: alignOfIRType(ir.String())})
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: yPtr, AllocaType: ir.String(), AllocaAlign: alignOfIRType(ir.String())})
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: x, StoreDst: xPtr, StoreAlign: alignOfIRType(ir.String())})
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: y, StoreDst: yPtr, StoreAlign: alignOfIRType(ir.String())})
		l.block.Append(ir.Instruction{Kind: ir.InstrCallVoid, CallName: "gominic_str_concat", CallArgs: []*ir.Value{outPtr, xPtr, yPtr}})
		res := l.newTemp(ir.String())
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: res, LoadSrc: outPtr, LoadAlign: alignOfIRType(ir.String())})
		return res, nil
	}

	if (op == token.EQL || op == token.NEQ) && x.Type() != nil && x.Type().Kind == ir.KindString && y.Type() != nil && y.Type().Kind == ir.KindString {
		dst := l.newTemp(ir.I1)
		p1, len1 := l.materializeStringParts(x)
		p2, len2 := l.materializeStringParts(y)
		l.block.Append(ir.Instruction{Kind: ir.InstrCall, Dest: dst, CallName: "gominic_str_eq", CallArgs: []*ir.Value{p1, len1, p2, len2}, CallRet: ir.I1})
		if op == token.NEQ {
			inv := l.newTemp(ir.I1)
			l.block.Append(ir.Instruction{Kind: ir.InstrICmp, Dest: inv, ICmpPred: ir.ICmpNe, ICmpX: dst, ICmpY: ir.NewConstant("1", ir.I1)})
			return inv, nil
		}
		return dst, nil
	}

	isFloat := (x.Type() != nil && x.Type().Basic == "double") || (y.Type() != nil && y.Type().Basic == "double")

	switch op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		if isFloat {
			x = l.ensureFloat64(x)
			y = l.ensureFloat64(y)
		}
		destTy := x.Type()
		if destTy == nil {
			destTy = y.Type()
		}
		dest := l.newTemp(destTy)
		var bop ir.BinOpKind
		if isFloat {
			switch op {
			case token.ADD:
				bop = ir.FAdd
			case token.SUB:
				bop = ir.FSub
			case token.MUL:
				bop = ir.FMul
			case token.QUO:
				bop = ir.FDiv
			default:
				return nil, fmt.Errorf("unsupported float op %s", op)
			}
		} else {
			switch op {
			case token.ADD:
				bop = ir.Add
			case token.SUB:
				bop = ir.Sub
			case token.MUL:
				bop = ir.Mul
			case token.QUO:
				bop = ir.SDiv
			case token.REM:
				bop = ir.SRem
			}
		}
		l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: dest, BinOp: bop, X: x, Y: y})
		return dest, nil
	case token.EQL, token.NEQ, token.LSS, token.GTR, token.LEQ, token.GEQ:
		if isFloat {
			x = l.ensureFloat64(x)
			y = l.ensureFloat64(y)
			dest := l.newTemp(ir.I1)
			var pred ir.FCmpPred
			switch op {
			case token.EQL:
				pred = ir.FCmpOeq
			case token.NEQ:
				pred = ir.FCmpOne
			case token.LSS:
				pred = ir.FCmpOlt
			case token.GTR:
				pred = ir.FCmpOgt
			case token.LEQ:
				pred = ir.FCmpOle
			case token.GEQ:
				pred = ir.FCmpOge
			}
			l.block.Append(ir.Instruction{Kind: ir.InstrFCmp, Dest: dest, FCmpPred: pred, FCmpX: x, FCmpY: y})
			return dest, nil
		}
		dest := l.newTemp(ir.I1)
		var pred ir.ICmpPred
		switch op {
		case token.EQL:
			pred = ir.ICmpEq
		case token.NEQ:
			pred = ir.ICmpNe
		case token.LSS:
			pred = ir.ICmpSlt
		case token.GTR:
			pred = ir.ICmpSgt
		case token.LEQ:
			pred = ir.ICmpSle
		case token.GEQ:
			pred = ir.ICmpSge
		}
		l.block.Append(ir.Instruction{Kind: ir.InstrICmp, Dest: dest, ICmpPred: pred, ICmpX: x, ICmpY: y})
		return dest, nil
	case token.SHL:
		// сдвиг влево
		destTy := x.Type()
		if destTy == nil {
			destTy = y.Type()
		}
		dest := l.newTemp(destTy)
		l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: dest, BinOp: ir.Shl, X: x, Y: y})
		return dest, nil
	case token.SHR:
		// Сдвиг вправо: пока используем арифметический вариант.
		destTy := x.Type()
		if destTy == nil {
			destTy = y.Type()
		}
		dest := l.newTemp(destTy)
		// Для простоты используем арифметический сдвиг вправо.
		// В полной версии тут нужна проверка знаковости типа.
		l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: dest, BinOp: ir.AShr, X: x, Y: y})
		return dest, nil
	case token.AND:
		// Побитовое AND (бинарный оператор, не unary &).
		destTy := x.Type()
		if destTy == nil {
			destTy = y.Type()
		}
		dest := l.newTemp(destTy)
		l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: dest, BinOp: ir.And, X: x, Y: y})
		return dest, nil
	case token.OR:
		// Побитовое OR (бинарный оператор, не логический ||).
		destTy := x.Type()
		if destTy == nil {
			destTy = y.Type()
		}
		dest := l.newTemp(destTy)
		l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: dest, BinOp: ir.Or, X: x, Y: y})
		return dest, nil
	default:
		return nil, fmt.Errorf("unsupported binop %s", op.String())
	}
}

// lowerLogical: реализует короткое замыкание для && и ||.
// Создаёт условные переходы, чтобы не вычислять правую часть, если результат уже известен
func (l *lowerer) lowerLogical(op token.Token, x *ir.Value, rhsExpr ast.Expr) (*ir.Value, error) {
	resPtr := ir.NewRegister(l.newTempName()+".bool", ir.PtrTo(ir.I1))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: resPtr, AllocaType: ir.I1, AllocaAlign: 1})

	rhsBB := l.newBlock(l.newBlockName("logic.rhs"))
	endBB := l.newBlock(l.newBlockName("logic.end"))
	shortBB := l.newBlock(l.newBlockName("logic.short"))

	if op == token.LAND {
		// if x == false -> short with false, else evaluate y
		l.block.Terminator = &ir.Instruction{Kind: ir.InstrCondBr, CondCond: x, CondTrue: rhsBB.Name, CondFalse: shortBB.Name}
		// короткая ветка: запись в память false
		shortBB.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: ir.NewConstant("0", ir.I1), StoreDst: resPtr, StoreAlign: 1})
		shortBB.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: endBB.Name}
	} else {
		// OR: if x == true -> short with true, else evaluate y
		l.block.Terminator = &ir.Instruction{Kind: ir.InstrCondBr, CondCond: x, CondTrue: shortBB.Name, CondFalse: rhsBB.Name}
		shortBB.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: ir.NewConstant("1", ir.I1), StoreDst: resPtr, StoreAlign: 1})
		shortBB.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: endBB.Name}
	}

	// RHS evaluation block
	l.block = rhsBB
	rhsVal, err := l.lowerExpr(rhsExpr)
	if err != nil {
		return nil, err
	}
	rhsBB.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: rhsVal, StoreDst: resPtr, StoreAlign: 1})
	if rhsBB.Terminator == nil {
		rhsBB.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: endBB.Name}
	}

	// переход в конечный блок
	l.block = endBB
	result := l.newTemp(ir.I1)
	endBB.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: result, LoadSrc: resPtr, LoadAlign: 1})
	return result, nil
}

// UnaryExpr = ("&" | "-" | "!") Expr
func (l *lowerer) lowerUnary(u *ast.UnaryExpr) (*ir.Value, error) {
	switch u.Op {
	case token.NOT:
		v, err := l.lowerExpr(u.X)
		if err != nil {
			return nil, err
		}
		dest := l.newTemp(ir.I1)
		zero := ir.NewConstant("0", ir.I1)
		l.block.Append(ir.Instruction{Kind: ir.InstrICmp, Dest: dest, ICmpPred: ir.ICmpEq, ICmpX: v, ICmpY: zero})
		return dest, nil
	case token.SUB:
		v, err := l.lowerExpr(u.X)
		if err != nil {
			return nil, err
		}
		// предполагаем int64
		dest := l.newTemp(v.Type())
		l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: dest, BinOp: ir.Sub, X: ir.NewConstant("0", v.Type()), Y: v})
		return dest, nil
	case token.AND:
		return l.lowerAddr(u.X)
	default:
		return nil, fmt.Errorf("unsupported unary %s", u.Op.String())
	}
}

// lowerStar: обрабатывает *ast.StarExpr (разыменование указателя).
func (l *lowerer) lowerStar(star *ast.StarExpr) (*ir.Value, error) {
	ptrVal, err := l.lowerExpr(star.X)
	if err != nil {
		return nil, err
	}
	if ptrVal.Type() == nil || ptrVal.Type().Kind != ir.KindPointer {
		return nil, fmt.Errorf("cannot dereference non-pointer")
	}
	elemTy := ptrVal.Type().Elem
	loaded := l.newTemp(elemTy)
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: ptrVal, LoadAlign: alignOfIRType(elemTy)})
	return loaded, nil
}

// lowerAddr: вычисляет адрес выражения (для &x, &arr[i], &s.field)
// Поддерживает только идентификаторы, индексацию и селекторы
func (l *lowerer) lowerAddr(e ast.Expr) (*ir.Value, error) {
	switch ex := e.(type) {
	case *ast.Ident:
		// Локальная переменная уже хранится как указатель
		if v := l.getLocal(ex.Name); v != nil {
			return v, nil
		}
		obj := l.prog.TypesInfo.Uses[ex]
		if obj == nil {
			return nil, fmt.Errorf("unknown ident %s", ex.Name)
		}
		ty, err := goTypeToIR(obj.Type())
		if err != nil {
			return nil, err
		}
		gname := ex.Name
		if l.prefix != "" {
			gname = l.prefix + gname
		}
		return ir.NewConstant("@"+gname, ir.PtrTo(ty)), nil
	case *ast.IndexExpr:
		baseVal, err := l.lowerExpr(ex.X)
		if err != nil {
			return nil, err
		}
		idxVal, err := l.lowerExpr(ex.Index)
		if err != nil {
			return nil, err
		}
		baseTy := baseVal.Type()
		switch baseTy.Kind {
		case ir.KindArray:
			elemTy := baseTy.Elem
			arrPtr := l.valuePtr(baseVal, baseTy)
			elemPtr := l.newTemp(ir.PtrTo(elemTy))
			l.block.Append(ir.Instruction{
				Kind:          ir.InstrGEP,
				Dest:          elemPtr,
				GepSrc:        arrPtr,
				GepPointee:    baseTy,
				GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), idxVal},
				GepResultType: ir.PtrTo(elemTy),
			})
			return elemPtr, nil
		case ir.KindSlice:
			elemTy := baseTy.Elem
			slicePtr := l.valuePtr(baseVal, baseTy)
			dataPtr := l.newTemp(ir.PtrTo(elemTy))
			l.block.Append(ir.Instruction{
				Kind:          ir.InstrGEP,
				Dest:          dataPtr,
				GepSrc:        slicePtr,
				GepPointee:    baseTy,
				GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("0", ir.I32)},
				GepResultType: ir.PtrTo(ir.PtrTo(elemTy)),
			})
			loadedData := l.newTemp(ir.PtrTo(elemTy))
			l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loadedData, LoadSrc: dataPtr, LoadAlign: 8})
			elemPtr := l.newTemp(ir.PtrTo(elemTy))
			l.block.Append(ir.Instruction{
				Kind:          ir.InstrGEP,
				Dest:          elemPtr,
				GepSrc:        loadedData,
				GepPointee:    elemTy,
				GepIndices:    []*ir.Value{idxVal},
				GepResultType: ir.PtrTo(elemTy),
			})
			return elemPtr, nil
		default:
			return nil, fmt.Errorf("address-of index on unsupported base")
		}
	case *ast.SelectorExpr:
		// переиспользуем lowering селектора, но сохраняем указатель
		baseVal, err := l.lowerExpr(ex.X)
		if err != nil {
			return nil, err
		}
		baseType, err := goTypeToIR(l.prog.TypesInfo.TypeOf(ex.X))
		if err != nil {
			return nil, err
		}
		structTy := baseType
		ptr := baseVal
		if baseType.Kind == ir.KindPointer {
			structTy = baseType.Elem
		} else {
			addr := l.newTemp(ir.PtrTo(baseType))
			l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: baseType, AllocaAlign: alignOfIRType(baseType)})
			l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: baseVal, StoreDst: addr, StoreAlign: alignOfIRType(baseType)})
			ptr = addr
		}
		sel := l.prog.TypesInfo.Selections[ex]
		if sel == nil {
			return nil, fmt.Errorf("no selection info")
		}
		// Проверяем, что structTy действительно структура.
		if structTy.Kind != ir.KindStruct {
			// Это рекурсивный тип, заменённый на указатель.
			// Возвращаем указатель: поля такого типа недоступны.
			return ir.NewConstant("null", ir.PtrI8), nil
		}
		indices := sel.Index()
		if len(indices) == 0 {
			return nil, fmt.Errorf("selector has no field index")
		}
		fieldIdx := indices[0]
		fieldTy, err := goTypeToIR(sel.Type())
		if err != nil {
			return nil, err
		}
		fieldPtr := l.newTemp(ir.PtrTo(fieldTy))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          fieldPtr,
			GepSrc:        ptr,
			GepPointee:    structTy,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant(strconv.Itoa(fieldIdx), ir.I32)},
			GepResultType: ir.PtrTo(fieldTy),
		})
		return fieldPtr, nil
	case *ast.CompositeLit:
		// Создаём временную переменную для composite literal.
		gt := l.prog.TypesInfo.TypeOf(ex)
		if gt == nil {
			return nil, fmt.Errorf("no type for composite literal")
		}
		irTy, err := goTypeToIR(gt)
		if err != nil {
			return nil, err
		}
		// Выделяем память под composite literal.
		addr := l.newTemp(ir.PtrTo(irTy))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: irTy, AllocaAlign: alignOfIRType(irTy)})

		// Понижаем composite literal и получаем его значение.
		val, err := l.lowerCompositeLit(ex)
		if err != nil {
			return nil, err
		}

		// Сохраняем значение в выделенную память.
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: val, StoreDst: addr, StoreAlign: alignOfIRType(irTy)})

		// возвращаем адрес
		return addr, nil
	default:
		return nil, fmt.Errorf("address-of unsupported expr %T", e)
	}
}

// Call = PrimaryExpr "(" [ ExprList ] ")"
// ExprList = Expr { "," Expr }
func (l *lowerer) lowerCall(c *ast.CallExpr) (*ir.Value, error) {
	// Обрабатываем конверсию, когда Fun — выражение типа (например, []byte(s)).
	if len(c.Args) == 1 {
		dstType := l.prog.TypesInfo.TypeOf(c.Fun)
		if dstType == nil {
			dstType = l.prog.TypesInfo.TypeOf(c)
		}
		if dstType != nil {
			if _, isSig := dstType.(*types.Signature); !isSig {
				dstTy, err := goTypeToIR(dstType)
				if err == nil {
					srcVal, err := l.lowerExpr(c.Args[0])
					if err != nil {
						return nil, err
					}
					// Отдельный путь для string -> []byte.
					// Конвертация string -> []byte: создаём срез и копируем данные через memcpy
					if dstTy.Kind == ir.KindSlice && dstTy.Elem == ir.I8 && srcVal.Type() != nil && srcVal.Type().Kind == ir.KindString {
						ptr, ln := l.materializeStringParts(srcVal)
						elemSize := ir.NewConstant("1", ir.I64)
						dataPtr := l.newTemp(ir.PtrTo(dstTy.Elem))
						l.block.Append(ir.Instruction{
							Kind:     ir.InstrCall,
							Dest:     dataPtr,
							CallName: "gominic_makeSlice",
							CallArgs: []*ir.Value{ln, ln, elemSize},
							CallRet:  ir.PtrTo(dstTy.Elem),
						})
						l.block.Append(ir.Instruction{Kind: ir.InstrCallVoid, CallName: "gominic_memcpy", CallArgs: []*ir.Value{dataPtr, ptr, ln}})
						return l.buildSliceValue(dstTy.Elem, dataPtr, ln, ln), nil
					}
					return l.lowerConvert(srcVal, dstTy)
				}
			}
		}
	}

	// Обрабатываем вызовы через селектор пакета (например, ir.ValueName(v)).
	if selExpr, ok := c.Fun.(*ast.SelectorExpr); ok {
		// Проверяем, что это вызов функции пакета.
		if pkgIdent, ok := selExpr.X.(*ast.Ident); ok {
			// Получаем объект функции.
			if obj, ok := l.prog.TypesInfo.Uses[selExpr.Sel]; ok {
				if fn, ok := obj.(*types.Func); ok {
					sig := fn.Type().(*types.Signature)
					// Методы: receiver передаётся как первый аргумент функции
					if sig.Recv() != nil {
						recvVal, err := l.lowerExpr(selExpr.X)
						if err != nil {
							return nil, err
						}
						args := make([]*ir.Value, len(c.Args))
						for i := 0; i < len(c.Args); i = i + 1 {
							v, err := l.lowerExpr(c.Args[i])
							if err != nil {
								return nil, err
							}
							args[i] = v
						}
						// Добавляем receiver в начало списка аргументов
						allArgs := make([]*ir.Value, len(args)+1)
						allArgs[0] = recvVal
						for i := 0; i < len(args); i = i + 1 {
							allArgs[i+1] = args[i]
						}
						// Определяем тип возвращаемого значения.
						var retTy *ir.TypeDesc
						if sig.Results().Len() == 1 {
							rt, err := goTypeToIR(sig.Results().At(0).Type())
							if err != nil {
								return nil, err
							}
							retTy = rt
						}
						// Имя метода: pkg.TypeName.MethodName (для самоприменимости)
						recvType := sig.Recv().Type()
						typeName := ""
						if named, ok := recvType.(*types.Named); ok {
							if named.Obj() != nil && named.Obj().Pkg() != nil {
								typeName = named.Obj().Pkg().Name() + "." + named.Obj().Name()
							} else {
								typeName = named.Obj().Name()
							}
						} else if ptr, ok := recvType.(*types.Pointer); ok {
							if named, ok := ptr.Elem().(*types.Named); ok {
								if named.Obj() != nil && named.Obj().Pkg() != nil {
									typeName = named.Obj().Pkg().Name() + "." + named.Obj().Name()
								} else {
									typeName = named.Obj().Name()
								}
							}
						}
						fnName := typeName + "." + selExpr.Sel.Name
						var dest *ir.Value
						if retTy != nil && retTy != ir.Void {
							dest = l.newTemp(retTy)
						}
						l.block.Append(ir.Instruction{Kind: ir.InstrCall, Dest: dest, CallName: fnName, CallArgs: allArgs, CallRet: retTy})
						return dest, nil
					}
					// Это вызов функции, не метода.
					// Понижаем аргументы.
					args := make([]*ir.Value, len(c.Args))
					for i := 0; i < len(c.Args); i = i + 1 {
						v, err := l.lowerExpr(c.Args[i])
						if err != nil {
							return nil, err
						}
						args[i] = v
					}
					// Определяем тип возвращаемого значения.
					var retTy *ir.TypeDesc
					if sig.Results().Len() == 1 {
						rt, err := goTypeToIR(sig.Results().At(0).Type())
						if err != nil {
							return nil, err
						}
						retTy = rt
					}
					// Формируем имя функции: pkg.FuncName.
					fnName := pkgIdent.Name + "." + selExpr.Sel.Name
					var dest *ir.Value
					if retTy != nil && retTy != ir.Void {
						dest = l.newTemp(retTy)
					}
					l.block.Append(ir.Instruction{Kind: ir.InstrCall, Dest: dest, CallName: fnName, CallArgs: args, CallRet: retTy})
					return dest, nil
				}
			}
		}
		// Обрабатываем вызовы методов на выражениях (например, entry.Append(...)).
		// Информацию о методе берем из Selections или Uses.
		sel := l.prog.TypesInfo.Selections[selExpr]
		var method *types.Func
		var isMethod bool
		if sel != nil {
			if sel.Kind() == types.MethodVal {
				if fn, ok := sel.Obj().(*types.Func); ok {
					method = fn
					isMethod = true
				}
			} else if sel.Kind() == types.FieldVal {
				// Это может быть поле-функция.
				if obj, ok := l.prog.TypesInfo.Uses[selExpr.Sel]; ok {
					if fn, ok := obj.(*types.Func); ok {
						method = fn
						isMethod = false
					}
				}
			}
		}
		if !isMethod {
			// Запасной путь: смотрим Uses.
			if obj, ok := l.prog.TypesInfo.Uses[selExpr.Sel]; ok {
				if fn, ok := obj.(*types.Func); ok {
					sig := fn.Type().(*types.Signature)
					if sig.Recv() != nil {
						method = fn
						isMethod = true
					} else {
						method = fn
						isMethod = false
					}
				}
			}
		}
		if method != nil {
			// Получаем сигнатуру метода.
			sig := method.Type().(*types.Signature)
			// Понижаем аргументы.
			args := make([]*ir.Value, len(c.Args))
			for i := 0; i < len(c.Args); i = i + 1 {
				v, err := l.lowerExpr(c.Args[i])
				if err != nil {
					return nil, err
				}
				args[i] = v
			}
			// Определяем тип возвращаемого значения.
			var retTy *ir.TypeDesc
			if sig.Results().Len() == 1 {
				rt, err := goTypeToIR(sig.Results().At(0).Type())
				if err != nil {
					return nil, err
				}
				retTy = rt
			}
			var fnName string
			var allArgs []*ir.Value
			if isMethod {
				// Это вызов метода: receiver добавляем первым аргументом.
				recvVal, err := l.lowerExpr(selExpr.X)
				if err != nil {
					return nil, err
				}
				allArgs = make([]*ir.Value, len(args)+1)
				allArgs[0] = recvVal
				for i := 0; i < len(args); i = i + 1 {
					allArgs[i+1] = args[i]
				}
				// Формируем имя: pkg.TypeName.MethodName.
				recvType := sig.Recv().Type()
				typeName := ""
				if named, ok := recvType.(*types.Named); ok {
					if named.Obj() != nil && named.Obj().Pkg() != nil {
						typeName = named.Obj().Pkg().Name() + "." + named.Obj().Name()
					} else if named.Obj() != nil {
						typeName = named.Obj().Name()
					}
				} else if ptr, ok := recvType.(*types.Pointer); ok {
					if named, ok := ptr.Elem().(*types.Named); ok {
						if named.Obj() != nil && named.Obj().Pkg() != nil {
							typeName = named.Obj().Pkg().Name() + "." + named.Obj().Name()
						} else if named.Obj() != nil {
							typeName = named.Obj().Name()
						}
					}
				}
				fnName = typeName + "." + selExpr.Sel.Name
			} else {
				// Это вызов поля-функции, receiver не нужен.
				allArgs = args
				// По типу базового выражения определяем пакет.
				baseType := l.prog.TypesInfo.TypeOf(selExpr.X)
				pkgName := ""
				if named, ok := baseType.(*types.Named); ok {
					if named.Obj() != nil && named.Obj().Pkg() != nil {
						pkgName = named.Obj().Pkg().Name()
					}
				} else if ptr, ok := baseType.(*types.Pointer); ok {
					if named, ok := ptr.Elem().(*types.Named); ok {
						if named.Obj() != nil && named.Obj().Pkg() != nil {
							pkgName = named.Obj().Pkg().Name()
						}
					}
				}
				if pkgName != "" {
					fnName = pkgName + "." + selExpr.Sel.Name
				} else {
					fnName = selExpr.Sel.Name
				}
			}
			var dest *ir.Value
			if retTy != nil && retTy != ir.Void {
				dest = l.newTemp(retTy)
			}
			l.block.Append(ir.Instruction{Kind: ir.InstrCall, Dest: dest, CallName: fnName, CallArgs: allArgs, CallRet: retTy})
			return dest, nil
		}
		// Последняя попытка: взять тип селектора и собрать имя функции.
		baseType := l.prog.TypesInfo.TypeOf(selExpr.X)
		selType := l.prog.TypesInfo.TypeOf(selExpr)
		if selType != nil {
			if sig, ok := selType.(*types.Signature); ok {
				// Понижаем аргументы.
				args := make([]*ir.Value, len(c.Args))
				for i := 0; i < len(c.Args); i = i + 1 {
					v, err := l.lowerExpr(c.Args[i])
					if err != nil {
						return nil, err
					}
					args[i] = v
				}
				// Определяем тип возвращаемого значения.
				var retTy *ir.TypeDesc
				if sig.Results().Len() == 1 {
					rt, err := goTypeToIR(sig.Results().At(0).Type())
					if err != nil {
						return nil, err
					}
					retTy = rt
				}
				// Пытаемся собрать имя функции из базового типа и селектора.
				fnName := selExpr.Sel.Name
				if baseType != nil {
					if named, ok := baseType.(*types.Named); ok {
						if named.Obj() != nil && named.Obj().Pkg() != nil {
							fnName = named.Obj().Pkg().Name() + "." + fnName
						}
					} else if ptr, ok := baseType.(*types.Pointer); ok {
						if named, ok := ptr.Elem().(*types.Named); ok {
							if named.Obj() != nil && named.Obj().Pkg() != nil {
								fnName = named.Obj().Pkg().Name() + "." + fnName
							}
						}
					}
				}
				var dest *ir.Value
				if retTy != nil && retTy != ir.Void {
					dest = l.newTemp(retTy)
				}
				l.block.Append(ir.Instruction{Kind: ir.InstrCall, Dest: dest, CallName: fnName, CallArgs: args, CallRet: retTy})
				return dest, nil
			}
		}
		return nil, fmt.Errorf("unsupported call target *ast.SelectorExpr (X=%T, Sel=%s)", selExpr.X, selExpr.Sel.Name)
	}

	fnIdent, ok := c.Fun.(*ast.Ident)
	if !ok {
		return nil, fmt.Errorf("unsupported call target %T", c.Fun)
	}
	// Конверсия типа: Fun является именем типа.
	if obj, ok := l.prog.TypesInfo.Uses[fnIdent]; ok {
		if _, isType := obj.(*types.TypeName); isType {
			if len(c.Args) != 1 {
				return nil, fmt.Errorf("conversion needs 1 arg")
			}
			dstTy, err := goTypeToIR(obj.Type())
			if err != nil {
				return nil, err
			}
			val, err := l.lowerExpr(c.Args[0])
			if err != nil {
				return nil, err
			}
			return l.lowerConvert(val, dstTy)
		}
	}
	// Встроенные функции: len/make/append.
	// Поддерживается только форма append(slice, elem).
	if fnIdent.Name == "len" {
		return l.lowerBuiltinLen(c.Args)
	}
	if fnIdent.Name == "make" {
		return l.lowerBuiltinMake(c.Args)
	}
	if fnIdent.Name == "append" {
		return l.lowerBuiltinAppend(c.Args)
	}
	// Обёртки print отдельно не разбираем: считаем обычным вызовом.
	args := make([]*ir.Value, len(c.Args))
	for i := 0; i < len(c.Args); i++ {
		v, err := l.lowerExpr(c.Args[i])
		if err != nil {
			return nil, err
		}
		args[i] = v
	}
	// Определяем тип возвращаемого значения.
	var retTy *ir.TypeDesc
	if obj := l.prog.TypesInfo.Uses[fnIdent]; obj != nil {
		if sig, ok := obj.Type().(*types.Signature); ok {
			if sig.Results().Len() == 1 {
				rt, err := goTypeToIR(sig.Results().At(0).Type())
				if err != nil {
					return nil, err
				}
				retTy = rt
			}
			// Методы с receiver в этом пути не поддерживаем.
			if sig.Recv() != nil {
				return nil, fmt.Errorf("unsupported: methods with receiver")
			}
		}
	}
	var dest *ir.Value
	if retTy != nil && retTy != ir.Void {
		dest = l.newTemp(retTy)
	}
	// Собираем имя функции с префиксом пакета.
	fnName := fnIdent.Name
	// Сначала пробуем l.prefix (он задаётся в buildInto для non-main).
	// Это совпадает с логикой lowerFunc.
	if l.prefix != "" {
		fnName = l.prefix + fnName
	} else {
		// Если prefix пуст, добавляем префикс только для экспортируемых функций.
		if len(fnName) > 0 && fnName[0] >= 'A' && fnName[0] <= 'Z' {
			pkgName := ""
			if l.prog != nil && l.prog.PkgName != "" {
				pkgName = l.prog.PkgName
			}
			// Если имя пакета пустое или main, пробуем взять его из объекта функции.
			// main здесь не используем: type checker иногда подставляет его всем пакетам.
			if pkgName == "" || pkgName == "main" {
				if obj := l.prog.TypesInfo.Uses[fnIdent]; obj != nil {
					if fn, ok := obj.(*types.Func); ok {
						if fn.Pkg() != nil {
							fnPkgName := fn.Pkg().Name()
							// Используем fn.Pkg() только если это не main.
							if fnPkgName != "" && fnPkgName != "main" {
								pkgName = fnPkgName
							}
						}
					}
				}
			}
			// Добавляем префикс, если имя пакета определено.
			if pkgName != "" {
				fnName = pkgName + "." + fnName
			}
		}
	}
	l.block.Append(ir.Instruction{Kind: ir.InstrCall, Dest: dest, CallName: fnName, CallArgs: args, CallRet: retTy})
	return dest, nil
}

// lowerConvert: конвертация типов (int->int64, int->float64 и т.д.).
// Использует инструкции SIToFP, FPToSI, PtrToInt, IntToPtr
func (l *lowerer) lowerConvert(v *ir.Value, dst *ir.TypeDesc) (*ir.Value, error) {
	if v.Type() == nil || dst == nil {
		return nil, fmt.Errorf("convert: missing type")
	}
	// Если типы совпадают, конвертация не нужна
	if v.Type().String() == dst.String() {
		return v, nil
	}
	out := l.newTemp(dst)
	// Конверсии между int64 и double.
	if v.Type().Basic == "double" && dst.Basic == "i64" {
		l.block.Append(ir.Instruction{Kind: ir.InstrConv, Dest: out, ConvOp: ir.FPToSI, ConvSrc: v, ConvTo: dst})
		return out, nil
	}
	if v.Type().Basic == "i64" && dst.Basic == "double" {
		l.block.Append(ir.Instruction{Kind: ir.InstrConv, Dest: out, ConvOp: ir.SIToFP, ConvSrc: v, ConvTo: dst})
		return out, nil
	}
	// Расширение int32 -> int64.
	if v.Type().Basic == "i32" && dst.Basic == "i64" {
		l.block.Append(ir.Instruction{Kind: ir.InstrConv, Dest: out, ConvOp: ir.SExt, ConvSrc: v, ConvTo: dst})
		return out, nil
	}
	// Bitcast указатель->указатель (редкий случай).
	if v.Type().Kind == ir.KindPointer && dst.Kind == ir.KindPointer {
		l.block.Append(ir.Instruction{Kind: ir.InstrBitcast, Dest: out, BitcastSrc: v, BitcastTarget: dst})
		return out, nil
	}
	// Конверсия i8* -> string.
	if v.Type().Kind == ir.KindPointer && v.Type().Elem != nil && v.Type().Elem == ir.I8 && dst.Kind == ir.KindString {
		// Создаем string-структуру: { i8* ptr, i64 len }.
		// Выделяем под нее память.
		strPtr := l.newTemp(ir.PtrTo(dst))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: strPtr, AllocaType: dst, AllocaAlign: 8})

		// Записываем указатель (первое поле).
		ptrFieldPtr := l.newTemp(ir.PtrTo(ir.PtrTo(ir.I8)))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          ptrFieldPtr,
			GepSrc:        strPtr,
			GepPointee:    dst,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("0", ir.I32)},
			GepResultType: ir.PtrTo(ir.PtrTo(ir.I8)),
		})
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: v, StoreDst: ptrFieldPtr, StoreAlign: 8})

		// Длину ставим 0 (второе поле), потому что по одному указателю длина неизвестна.
		lenFieldPtr := l.newTemp(ir.PtrTo(ir.I64))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          lenFieldPtr,
			GepSrc:        strPtr,
			GepPointee:    dst,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("1", ir.I32)},
			GepResultType: ir.PtrTo(ir.I64),
		})
		zeroLen := ir.NewConstant("0", ir.I64)
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroLen, StoreDst: lenFieldPtr, StoreAlign: 8})

		// Загружаем готовую строковую структуру.
		loaded := l.newTemp(dst)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: strPtr, LoadAlign: 8})
		return loaded, nil
	}
	// Конверсия []byte -> string.
	if v.Type().Kind == ir.KindSlice && v.Type().Elem != nil && v.Type().Elem == ir.I8 && dst.Kind == ir.KindString {
		// Извлекаем указатель на данные и длину из среза.
		slicePtr := l.valuePtr(v, v.Type())
		// Читаем указатель на данные.
		dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrTo(ir.I8)))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          dataPtrPtr,
			GepSrc:        slicePtr,
			GepPointee:    v.Type(),
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("0", ir.I32)},
			GepResultType: ir.PtrTo(ir.PtrTo(ir.I8)),
		})
		ptr := l.newTemp(ir.PtrTo(ir.I8))
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: ptr, LoadSrc: dataPtrPtr, LoadAlign: 8})
		// Читаем длину.
		lenPtr := l.newTemp(ir.PtrTo(ir.I64))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          lenPtr,
			GepSrc:        slicePtr,
			GepPointee:    v.Type(),
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("1", ir.I32)},
			GepResultType: ir.PtrTo(ir.I64),
		})
		lenVal := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: lenVal, LoadSrc: lenPtr, LoadAlign: 8})
		// Собираем string из указателя и длины.
		return l.buildStringValue(ptr, lenVal), nil
	}
	return nil, fmt.Errorf("conversion from %s to %s not supported", v.Type().String(), dst.String())
}

// ensureFloat64: приводит значение к float64 для float-операций.
func (l *lowerer) ensureFloat64(v *ir.Value) *ir.Value {
	if v.Type() != nil && v.Type().Basic == "double" {
		return v
	}
	dest := l.newTemp(ir.F64)
	l.block.Append(ir.Instruction{Kind: ir.InstrConv, Dest: dest, ConvOp: ir.SIToFP, ConvSrc: v, ConvTo: ir.F64})
	return dest
}

func (l *lowerer) lowerBuiltinLen(args []ast.Expr) (*ir.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len: want 1 arg")
	}
	v, err := l.lowerExpr(args[0])
	if err != nil {
		return nil, err
	}
	ty := v.Type()
	switch ty.Kind {
	case ir.KindArray:
		return ir.NewConstant(strconv.Itoa(ty.Len), ir.I64), nil
	case ir.KindSlice:
		ptr := l.valuePtr(v, ty)
		lenPtr := l.newTemp(ir.PtrTo(ir.I64))
		lenVal := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          lenPtr,
			GepSrc:        ptr,
			GepPointee:    ty,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("1", ir.I32)},
			GepResultType: ir.PtrTo(ir.I64),
		})
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: lenVal, LoadSrc: lenPtr, LoadAlign: 8})
		return lenVal, nil
	case ir.KindString:
		ptr := l.valuePtr(v, ty)
		lenPtr := l.newTemp(ir.PtrTo(ir.I64))
		lenVal := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          lenPtr,
			GepSrc:        ptr,
			GepPointee:    ty,
			GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("1", ir.I32)},
			GepResultType: ir.PtrTo(ir.I64),
		})
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: lenVal, LoadSrc: lenPtr, LoadAlign: 8})
		return lenVal, nil
	case ir.KindBasic: // map treated as PtrI8
		if ty == ir.PtrI8 {
			dst := l.newTemp(ir.I64)
			l.block.Append(ir.Instruction{Kind: ir.InstrCall, Dest: dst, CallName: "gominic_map_len", CallArgs: []*ir.Value{v}, CallRet: ir.I64})
			return dst, nil
		}
	}
	return nil, fmt.Errorf("len unsupported type")
}

func (l *lowerer) lowerBuiltinMake(args []ast.Expr) (*ir.Value, error) {
	makeTy, err := goTypeToIR(l.prog.TypesInfo.TypeOf(args[0]))
	if err != nil {
		return nil, err
	}
	// Сначала обрабатываем map.
	if mapType, ok := l.prog.TypesInfo.TypeOf(args[0]).Underlying().(*types.Map); ok {
		keyTy, err := goTypeToIR(mapType.Key())
		if err != nil {
			return nil, err
		}
		valTy, err := goTypeToIR(mapType.Elem())
		if err != nil {
			return nil, err
		}
		keySize := ir.NewConstant(strconv.FormatInt(sizeOfIRType(keyTy), 10), ir.I64)
		valSize := ir.NewConstant(strconv.FormatInt(sizeOfIRType(valTy), 10), ir.I64)
		keyKind := ir.NewConstant(strconv.Itoa(keyKindFromType(mapType.Key())), ir.I32)
		dst := l.newTemp(ir.PtrI8)
		l.block.Append(ir.Instruction{
			Kind:     ir.InstrCall,
			Dest:     dst,
			CallName: "gominic_map_new",
			CallArgs: []*ir.Value{keySize, valSize, keyKind},
			CallRet:  ir.PtrI8,
		})
		return dst, nil
	}

	switch makeTy.Kind {
	case ir.KindSlice:
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("make: want 2 or 3 args")
		}
		lenVal, err := l.lowerExpr(args[1])
		if err != nil {
			return nil, err
		}
		capVal := lenVal
		if len(args) == 3 {
			capVal, err = l.lowerExpr(args[2])
			if err != nil {
				return nil, err
			}
		}
		elemSize := ir.NewConstant(strconv.Itoa(int(sizeOfIRType(makeTy.Elem))), ir.I64)
		ptr := l.newTemp(ir.PtrTo(makeTy.Elem))
		l.block.Append(ir.Instruction{
			Kind:     ir.InstrCall,
			Dest:     ptr,
			CallName: "gominic_makeSlice",
			CallArgs: []*ir.Value{lenVal, capVal, elemSize},
			CallRet:  ir.PtrTo(makeTy.Elem),
		})
		return l.buildSliceValue(makeTy.Elem, ptr, lenVal, capVal), nil
	}
	return nil, fmt.Errorf("make unsupported type")
}

// Поддерживается только форма append(slice, elem).
func (l *lowerer) lowerBuiltinAppend(args []ast.Expr) (*ir.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("append: want 2 args")
	}
	sliceVal, err := l.lowerExpr(args[0])
	if err != nil {
		return nil, err
	}
	elemVal, err := l.lowerExpr(args[1])
	if err != nil {
		return nil, err
	}
	if sliceVal.Type() == nil || sliceVal.Type().Kind != ir.KindSlice {
		return nil, fmt.Errorf("append: first arg not slice")
	}
	sliceTy := sliceVal.Type()
	elemTy := sliceTy.Elem
	slicePtr := l.valuePtr(sliceVal, sliceTy)
	elemSize := ir.NewConstant(strconv.FormatInt(sizeOfIRType(elemTy), 10), ir.I64)

	// Читаем data/len/cap из среза.
	dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrTo(elemTy)))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          dataPtrPtr,
		GepSrc:        slicePtr,
		GepPointee:    sliceTy,
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("0", ir.I32)},
		GepResultType: ir.PtrTo(ir.PtrTo(elemTy)),
	})
	dataPtr := l.newTemp(ir.PtrTo(elemTy))
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: dataPtr, LoadSrc: dataPtrPtr, LoadAlign: 8})

	lenPtr := l.newTemp(ir.PtrTo(ir.I64))
	lenVal := l.newTemp(ir.I64)
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          lenPtr,
		GepSrc:        slicePtr,
		GepPointee:    sliceTy,
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("1", ir.I32)},
		GepResultType: ir.PtrTo(ir.I64),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: lenVal, LoadSrc: lenPtr, LoadAlign: 8})

	newLen := l.newTemp(ir.I64)
	l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: newLen, BinOp: ir.Add, X: lenVal, Y: ir.NewConstant("1", ir.I64)})

	capPtr := l.newTemp(ir.PtrTo(ir.I64))
	capVal := l.newTemp(ir.I64)
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          capPtr,
		GepSrc:        slicePtr,
		GepPointee:    sliceTy,
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("2", ir.I32)},
		GepResultType: ir.PtrTo(ir.I64),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: capVal, LoadSrc: capPtr, LoadAlign: 8})

	// Решаем, нужно ли расширение: newLen > cap.
	needGrow := l.newTemp(ir.I1)
	l.block.Append(ir.Instruction{Kind: ir.InstrICmp, Dest: needGrow, ICmpPred: ir.ICmpSgt, ICmpX: newLen, ICmpY: capVal})

	// Через стек храним итоговые data/cap, чтобы обойтись без phi.
	outDataAddr := l.newTemp(ir.PtrTo(ir.PtrTo(elemTy)))
	outCapAddr := l.newTemp(ir.PtrTo(ir.I64))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: outDataAddr, AllocaType: ir.PtrTo(elemTy), AllocaAlign: 8})
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: outCapAddr, AllocaType: ir.I64, AllocaAlign: 8})

	growBB := l.newBlock(l.newBlockName("append.grow"))
	noGrowBB := l.newBlock(l.newBlockName("append.nogrow"))
	contBB := l.newBlock(l.newBlockName("append.cont"))
	l.block.Terminator = &ir.Instruction{Kind: ir.InstrCondBr, CondCond: needGrow, CondTrue: growBB.Name, CondFalse: noGrowBB.Name}
	l.block = growBB

	// Ветка роста: выделяем новый буфер, копируем старые элементы.
	newDataPtr := l.newTemp(ir.PtrTo(elemTy))
	l.block.Append(ir.Instruction{
		Kind:     ir.InstrCall,
		Dest:     newDataPtr,
		CallName: "gominic_makeSlice",
		CallArgs: []*ir.Value{newLen, newLen, elemSize},
		CallRet:  ir.PtrTo(elemTy),
	})
	oldBytes := l.newTemp(ir.I64)
	l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: oldBytes, BinOp: ir.Mul, X: lenVal, Y: elemSize})
	newDataI8 := l.newTemp(ir.PtrI8)
	oldDataI8 := l.newTemp(ir.PtrI8)
	l.block.Append(ir.Instruction{Kind: ir.InstrBitcast, Dest: newDataI8, BitcastSrc: newDataPtr, BitcastTarget: ir.PtrI8})
	l.block.Append(ir.Instruction{Kind: ir.InstrBitcast, Dest: oldDataI8, BitcastSrc: dataPtr, BitcastTarget: ir.PtrI8})
	l.block.Append(ir.Instruction{Kind: ir.InstrCallVoid, CallName: "gominic_memcpy", CallArgs: []*ir.Value{newDataI8, oldDataI8, oldBytes}})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: newDataPtr, StoreDst: outDataAddr, StoreAlign: 8})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: newLen, StoreDst: outCapAddr, StoreAlign: 8})
	l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: contBB.Name}

	// Без роста: пишем в текущий буфер.
	l.block = noGrowBB
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: dataPtr, StoreDst: outDataAddr, StoreAlign: 8})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: capVal, StoreDst: outCapAddr, StoreAlign: 8})
	l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: contBB.Name}

	// Общая часть: записываем элемент в индекс len и собираем новый slice.
	l.block = contBB
	finalDataPtr := l.newTemp(ir.PtrTo(elemTy))
	finalCap := l.newTemp(ir.I64)
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: finalDataPtr, LoadSrc: outDataAddr, LoadAlign: 8})
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: finalCap, LoadSrc: outCapAddr, LoadAlign: 8})

	elemPtr := l.newTemp(ir.PtrTo(elemTy))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          elemPtr,
		GepSrc:        finalDataPtr,
		GepPointee:    elemTy,
		GepIndices:    []*ir.Value{lenVal},
		GepResultType: ir.PtrTo(elemTy),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: elemVal, StoreDst: elemPtr, StoreAlign: alignOfIRType(elemTy)})

	return l.buildSliceValue(elemTy, finalDataPtr, newLen, finalCap), nil
}

func goTypeToIR(t types.Type) (*ir.TypeDesc, error) {
	return goTypeToIRVisited(t, make(map[string]bool))
}

func goTypeToIRVisited(t types.Type, visited map[string]bool) (*ir.TypeDesc, error) {
	// Обрабатываем nil-тип.
	if t == nil {
		return ir.Void, nil
	}
	// Сначала обрабатываем именованные типы, чтобы поймать рекурсию до Underlying().
	if named, ok := t.(*types.Named); ok {
		if named.Obj() != nil {
			typeName := named.Obj().Name()
			// Отдельно обрабатываем встроенный error.
			if typeName == "error" {
				// error в IR представляем как указатель.
				return ir.PtrI8, nil
			}
			if named.Obj().Pkg() != nil {
				pkgName := named.Obj().Pkg().Name()
				pkgPath := named.Obj().Pkg().Path()
				// Типы из ir/merged-пакетов представляем указателями.
				if pkgName == "ir" || pkgName == "backend" || pkgName == "main" {
					// Проверяем, что это один из внутренних IR-типов.
					if typeName == "TypeDesc" || typeName == "Value" || typeName == "Instruction" ||
						typeName == "BasicBlock" || typeName == "Function" || typeName == "Module" ||
						typeName == "Global" {
						return ir.PtrI8, nil
					}
				}
				// Неподдержанные типы stdlib заменяем на указатели.
				if pkgPath == "sync/atomic" || pkgPath == "sync" || pkgPath == "internal/race" ||
					pkgPath == "syscall" || pkgPath == "internal/syscall/windows" ||
					pkgPath == "internal/syscall/windows/sysdll" || pkgPath == "internal/poll" ||
					pkgPath == "internal/syscall/unix" || pkgPath == "time" ||
					pkgPath == "go/token" || pkgPath == "go/ast" || pkgPath == "go/types" ||
					pkgPath == "go/parser" || pkgPath == "go/importer" || pkgPath == "go/format" ||
					pkgPath == "go/printer" || pkgPath == "go/doc" || pkgPath == "go/build" ||
					pkgPath == "go/scanner" || pkgPath == "go/constant" ||
					pkgPath == "io" || pkgPath == "os" || pkgPath == "path/filepath" ||
					pkgPath == "fmt" || pkgPath == "strconv" || pkgPath == "strings" ||
					pkgPath == "bytes" || pkgPath == "unicode" || pkgPath == "unicode/utf8" ||
					pkgPath == "unicode/utf16" || pkgPath == "path" || pkgPath == "errors" ||
					pkgPath == "reflect" || pkgPath == "runtime" || pkgPath == "internal/bytealg" ||
					pkgPath == "internal/cpu" || pkgPath == "internal/reflectlite" ||
					pkgPath == "internal/unsafeheader" || pkgPath == "unsafe" {
					return ir.PtrI8, nil
				}
			}
		}
		key := named.String()
		if named.Obj() != nil && named.Obj().Pkg() != nil {
			key = named.Obj().Pkg().Path() + "." + named.Obj().Name()
		}
		if visited[key] {
			// Для внутренних IR-типов допускаем рекурсию.
			if named.Obj() != nil {
				typeName := named.Obj().Name()
				pkgPath := ""
				if named.Obj().Pkg() != nil {
					pkgPath = named.Obj().Pkg().Path()
				}
				if typeName == "TypeDesc" || typeName == "Value" || typeName == "Instruction" ||
					typeName == "BasicBlock" || typeName == "Function" || typeName == "Module" ||
					typeName == "Global" {
					return ir.PtrI8, nil
				}
				// Для рекурсивных stdlib-типов тоже используем указатели.
				if pkgPath == "sync/atomic" || pkgPath == "sync" || pkgPath == "internal/race" ||
					pkgPath == "syscall" || pkgPath == "internal/syscall/windows" ||
					pkgPath == "internal/syscall/windows/sysdll" || pkgPath == "internal/poll" ||
					pkgPath == "internal/syscall/unix" || pkgPath == "time" ||
					pkgPath == "go/token" || pkgPath == "go/ast" || pkgPath == "go/types" ||
					pkgPath == "go/parser" || pkgPath == "go/importer" || pkgPath == "go/format" ||
					pkgPath == "go/printer" || pkgPath == "go/doc" || pkgPath == "go/build" ||
					pkgPath == "go/scanner" || pkgPath == "go/constant" ||
					pkgPath == "io" || pkgPath == "os" || pkgPath == "path/filepath" ||
					pkgPath == "fmt" || pkgPath == "strconv" || pkgPath == "strings" ||
					pkgPath == "bytes" || pkgPath == "unicode" || pkgPath == "unicode/utf8" ||
					pkgPath == "unicode/utf16" || pkgPath == "path" || pkgPath == "errors" ||
					pkgPath == "reflect" || pkgPath == "runtime" || pkgPath == "internal/bytealg" ||
					pkgPath == "internal/cpu" || pkgPath == "internal/reflectlite" ||
					pkgPath == "internal/unsafeheader" || pkgPath == "unsafe" {
					return ir.PtrI8, nil
				}
			}
			return nil, fmt.Errorf("unsupported recursive type %v", named)
		}
		visited[key] = true
		return goTypeToIRVisited(named.Underlying(), visited)
	}

	// Теперь обрабатываем underlying-типы.
	underlying := t.Underlying()
	if underlying == nil {
		return ir.Void, nil
	}
	switch tt := underlying.(type) {
	case *types.Basic:
		switch tt.Kind() {
		case types.Bool:
			return ir.I1, nil
		case types.UntypedBool:
			return ir.I1, nil
		case types.Int8, types.Uint8:
			return ir.I8, nil
		case types.Int, types.Int64, types.UntypedInt:
			return ir.I64, nil
		case types.Uint64, types.Uint:
			return ir.I64, nil
		case types.Int32, types.UntypedRune:
			return ir.I32, nil
		case types.Uint32:
			return ir.I32, nil
		case types.Float64, types.UntypedFloat:
			return ir.F64, nil
		case types.String, types.UntypedString:
			return ir.String(), nil
		case types.UntypedNil:
			return ir.PtrI8, nil
		default:
			return nil, fmt.Errorf("unsupported basic type %v", tt)
		}
	case *types.Pointer:
		elem, err := goTypeToIRVisited(tt.Elem(), visited)
		if err != nil {
			return nil, err
		}
		return ir.PtrTo(elem), nil
	case *types.Array:
		elem, err := goTypeToIRVisited(tt.Elem(), visited)
		if err != nil {
			return nil, err
		}
		return ir.Array(int(tt.Len()), elem), nil
	case *types.Slice:
		elem, err := goTypeToIRVisited(tt.Elem(), visited)
		if err != nil {
			return nil, err
		}
		return ir.Slice(elem), nil
	case *types.Map:
		return ir.PtrI8, nil
	case *types.Interface:
		// Интерфейсы в этом подмножестве не моделируем, используем указатель.
		return ir.PtrI8, nil
	case *types.Signature:
		// Сигнатуры функций не моделируем, используем указатель.
		return ir.PtrI8, nil
	case *types.Chan:
		// Каналы не поддержаны, используем указатель.
		return ir.PtrI8, nil
	case *types.Struct:
		fields := make([]*ir.TypeDesc, tt.NumFields())
		for i := 0; i < tt.NumFields(); i++ {
			ft, err := goTypeToIRVisited(tt.Field(i).Type(), visited)
			if err != nil {
				return nil, err
			}
			fields[i] = ft
		}
		return ir.Struct(fields), nil
	default:
		return nil, fmt.Errorf("unsupported go type %T", t)
	}
}

func keyKindFromType(t types.Type) int {
	switch tt := t.Underlying().(type) {
	case *types.Basic:
		switch tt.Kind() {
		case types.Bool, types.Uint, types.Uint64, types.Uint32, types.UntypedBool:
			return 1
		case types.String, types.UntypedString:
			return 2
		default:
			return 0 // signed ints by default
		}
	case *types.Pointer:
		return 3
	default:
		return 0
	}
}

func elemTypeFromPtr(t *ir.TypeDesc) *ir.TypeDesc {
	if t != nil && t.Kind == ir.KindPointer {
		return t.Elem
	}
	return ir.I64
}

func alignOfIRType(t *ir.TypeDesc) int64 {
	if t == nil {
		return 8
	}
	switch t.Kind {
	case ir.KindBasic:
		switch t.Basic {
		case "i1", "i8":
			return 1
		case "i32":
			return 4
		default:
			return 8
		}
	case ir.KindPointer:
		return 8
	case ir.KindString, ir.KindSlice:
		return 8
	default:
		return 8
	}
}

func sizeOfIRType(t *ir.TypeDesc) int64 {
	if t == nil {
		return 0
	}
	switch t.Kind {
	case ir.KindBasic:
		switch t.Basic {
		case "i1", "i8":
			return 1
		case "i32":
			return 4
		case "double":
			return 8
		default:
			return 8
		}
	case ir.KindPointer:
		return 8
	case ir.KindString:
		return 16
	case ir.KindSlice:
		return 24
	case ir.KindStruct:
		var off int64 = 0
		for i := 0; i < len(t.Fields); i++ {
			a := alignOfIRType(t.Fields[i])
			off = alignUp(off, a)
			off += sizeOfIRType(t.Fields[i])
		}
		return alignUp(off, alignOfIRType(t))
	case ir.KindArray:
		return int64(t.Len) * sizeOfIRType(t.Elem)
	default:
		return 8
	}
}

func alignUp(x, a int64) int64 {
	if a <= 1 {
		return x
	}
	r := x % a
	if r == 0 {
		return x
	}
	return x + (a - r)
}

func constantString(v *ir.Value) (string, error) {
	if v == nil || v.Kind() != ir.ValueConstant {
		return "", fmt.Errorf("not a const")
	}
	return v.Name(), nil
}

func (l *lowerer) stringConstant(s string) *ir.Value {
	// Добавляем глобал с именем .str.N (с префиксом пакета, если есть)
	l.strID++
	gname := fmt.Sprintf(".str.%d", l.strID)
	if l.prefix != "" {
		// убираем завершающую точку у префикса
		prefix := l.prefix
		if len(prefix) > 0 && prefix[len(prefix)-1] == '.' {
			prefix = prefix[:len(prefix)-1]
		}
		gname = prefix + gname
	} else if l.prog != nil && l.prog.PkgName != "" && l.prog.PkgName != "main" {
		gname = l.prog.PkgName + gname
	}
	bytes := []byte(s)
	arr := ir.Array(len(bytes)+1, ir.I8)
	enc := escapeCString(bytes)
	l.mod.Globals = append(l.mod.Globals, ir.Global{
		Name:    gname,
		Type:    arr,
		Value:   fmt.Sprintf("c\"%s\"", enc),
		Align:   1,
		Private: true,
	})
	gep := fmt.Sprintf("getelementptr inbounds ([%d x i8], [%d x i8]* @%s, i32 0, i32 0)", len(bytes)+1, len(bytes)+1, gname)
	ptr := ir.NewConstant(gep, ir.PtrTo(ir.I8))
	raw := fmt.Sprintf("{ i8* %s, i64 %d }", ptr.Name(), len(bytes))
	return ir.NewConstant(raw, ir.String())
}

func (l *lowerer) constValue(e ast.Expr) (*ir.Value, error) {
	if bl, ok := e.(*ast.BasicLit); ok {
		return l.lowerBasicLit(bl)
	}
	return nil, fmt.Errorf("not const")
}

func (l *lowerer) newTempName() string {
	l.reg++
	return fmt.Sprintf("t%d", l.reg)
}

func escapeCString(b []byte) string {
	out := make([]byte, 0, len(b)*4+4)
	for _, ch := range b {
		out = append(out, '\\')
		out = append(out, []byte(fmt.Sprintf("%02X", ch))...)
	}
	out = append(out, '\\', '0', '0')
	return string(out)
}

// materializeStringParts: извлекает указатель и длину из строки { i8*, i64 }.
// Выделяет память, сохраняет строку, затем загружает поля через GEP
func (l *lowerer) materializeStringParts(str *ir.Value) (*ir.Value, *ir.Value) {
	addr := l.newTemp(ir.PtrTo(ir.String()))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: ir.String(), AllocaAlign: alignOfIRType(ir.String())})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: str, StoreDst: addr, StoreAlign: alignOfIRType(ir.String())})
	// указатель
	ptrPtr := l.newTemp(ir.PtrTo(ir.PtrTo(ir.I8)))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          ptrPtr,
		GepSrc:        addr,
		GepPointee:    ir.String(),
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("0", ir.I32)},
		GepResultType: ir.PtrTo(ir.PtrTo(ir.I8)),
	})
	ptrVal := l.newTemp(ir.PtrTo(ir.I8))
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: ptrVal, LoadSrc: ptrPtr, LoadAlign: 8})
	// длина
	lenPtr := l.newTemp(ir.PtrTo(ir.I64))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          lenPtr,
		GepSrc:        addr,
		GepPointee:    ir.String(),
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("1", ir.I32)},
		GepResultType: ir.PtrTo(ir.I64),
	})
	lenVal := l.newTemp(ir.I64)
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: lenVal, LoadSrc: lenPtr, LoadAlign: 8})
	return ptrVal, lenVal
}

// lowerMapSet: установка значения в map через runtime функцию
// Ключ и значение сохраняются в буферы, затем приводятся к i8* для передачи в runtime
func (l *lowerer) lowerMapSet(idx *ast.IndexExpr, rhs ast.Expr) error {
	mapVal, err := l.lowerExpr(idx.X)
	if err != nil {
		return err
	}
	mapType, _ := l.prog.TypesInfo.TypeOf(idx.X).Underlying().(*types.Map)
	keyVal, err := l.lowerExpr(idx.Index)
	if err != nil {
		return err
	}
	valVal, err := l.lowerExpr(rhs)
	if err != nil {
		return err
	}
	keyTy, err := goTypeToIR(l.prog.TypesInfo.TypeOf(idx.Index))
	if err != nil {
		return err
	}
	valGo := l.prog.TypesInfo.TypeOf(idx)
	if mapType != nil {
		valGo = mapType.Elem()
	}
	valTy, err := goTypeToIR(valGo)
	if err != nil {
		return err
	}
	// Сохраняем ключ и значение в буферы на стеке
	keyBuf := l.newTemp(ir.PtrTo(keyTy))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: keyBuf, AllocaType: keyTy, AllocaAlign: alignOfIRType(keyTy)})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: keyVal, StoreDst: keyBuf, StoreAlign: alignOfIRType(keyTy)})
	valBuf := l.newTemp(ir.PtrTo(valTy))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: valBuf, AllocaType: valTy, AllocaAlign: alignOfIRType(valTy)})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: valVal, StoreDst: valBuf, StoreAlign: alignOfIRType(valTy)})
	// Приводим к i8* для передачи в runtime (runtime работает с байтами)
	keyI8 := l.newTemp(ir.PtrI8)
	valI8 := l.newTemp(ir.PtrI8)
	l.block.Append(ir.Instruction{Kind: ir.InstrBitcast, Dest: keyI8, BitcastSrc: keyBuf, BitcastTarget: ir.PtrI8})
	l.block.Append(ir.Instruction{Kind: ir.InstrBitcast, Dest: valI8, BitcastSrc: valBuf, BitcastTarget: ir.PtrI8})
	l.block.Append(ir.Instruction{
		Kind:     ir.InstrCall,
		CallName: "gominic_map_set",
		CallArgs: []*ir.Value{mapVal, keyI8, valI8},
	})
	return nil
}

// lowerMapGet: получение значения из map с проверкой наличия (возвращает value, ok)
// Использует runtime функцию gominic_map_get
func (l *lowerer) lowerMapGet(idx *ast.IndexExpr, lhs0 ast.Expr, lhs1 ast.Expr, define bool) error {
	mapVal, err := l.lowerExpr(idx.X)
	if err != nil {
		return err
	}
	mapType, _ := l.prog.TypesInfo.TypeOf(idx.X).Underlying().(*types.Map)
	keyVal, err := l.lowerExpr(idx.Index)
	if err != nil {
		return err
	}
	valGo := l.prog.TypesInfo.TypeOf(idx)
	if mapType != nil {
		valGo = mapType.Elem()
	}
	valTy, err := goTypeToIR(valGo)
	if err != nil {
		return err
	}
	keyTy, err := goTypeToIR(l.prog.TypesInfo.TypeOf(idx.Index))
	if err != nil {
		return err
	}
	keyBuf := l.newTemp(ir.PtrTo(keyTy))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: keyBuf, AllocaType: keyTy, AllocaAlign: alignOfIRType(keyTy)})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: keyVal, StoreDst: keyBuf, StoreAlign: alignOfIRType(keyTy)})
	valBuf := l.newTemp(ir.PtrTo(valTy))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: valBuf, AllocaType: valTy, AllocaAlign: alignOfIRType(valTy)})
	keyI8 := l.newTemp(ir.PtrI8)
	valI8 := l.newTemp(ir.PtrI8)
	l.block.Append(ir.Instruction{Kind: ir.InstrBitcast, Dest: keyI8, BitcastSrc: keyBuf, BitcastTarget: ir.PtrI8})
	l.block.Append(ir.Instruction{Kind: ir.InstrBitcast, Dest: valI8, BitcastSrc: valBuf, BitcastTarget: ir.PtrI8})
	okVal := l.newTemp(ir.I1)
	l.block.Append(ir.Instruction{
		Kind:     ir.InstrCall,
		Dest:     okVal,
		CallName: "gominic_map_get",
		CallArgs: []*ir.Value{mapVal, keyI8, valI8},
		CallRet:  ir.I1,
	})
	loaded := l.newTemp(valTy)
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: valBuf, LoadAlign: alignOfIRType(valTy)})
	l.assignToLHS(lhs0, loaded, define)
	l.assignToLHS(lhs1, okVal, define)
	return nil
}

func (l *lowerer) assignToLHS(lhs ast.Expr, val *ir.Value, define bool) {
	id, ok := lhs.(*ast.Ident)
	if !ok {
		return
	}
	if id.Name == "_" {
		return
	}
	addr := l.getLocal(id.Name)
	if addr == nil && define {
		addr = l.ensureLocal(id.Name, val.Type())
	}
	if addr != nil {
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: val, StoreDst: addr, StoreAlign: alignOfIRType(val.Type())})
	}
}

func (l *lowerer) ensureLocal(name string, ty *ir.TypeDesc) *ir.Value {
	if v := l.getLocal(name); v != nil {
		return v
	}
	l.localID++
	addr := ir.NewRegister(fmt.Sprintf("%s.addr%d", name, l.localID), ir.PtrTo(ty))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: ty, AllocaAlign: alignOfIRType(ty)})
	l.setLocal(name, addr)
	return addr
}
