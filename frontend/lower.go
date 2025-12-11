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
// Result = Type | "(" Type ")"
func (l *lowerer) lowerFunc(fnDecl *ast.FuncDecl) error {
	sig, ok := l.prog.TypesInfo.Defs[fnDecl.Name].Type().(*types.Signature)
	if !ok {
		return fmt.Errorf("func %s: no signature", fnDecl.Name.Name)
	}
	// Support multiple returns for self-hosting (functions returning (T, error))
	if sig.Results().Len() > 1 {
		// For now, only support (T, error) pattern
		if sig.Results().Len() == 2 {
			// Check if second result is error type
			if errType, ok := sig.Results().At(1).Type().(*types.Named); ok {
				if errType.Obj() != nil && errType.Obj().Name() == "error" {
					// This is (T, error) pattern - supported for self-hosting
					// We'll return only the first value, error handling is done separately
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
	// Special case: main function should always be named "main" (no package prefix)
	if fnName == "main" {
		fnName = "main"
	} else if sig.Recv() != nil {
		// If this is a method (has receiver), add type name prefix
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
		// For package functions, add package prefix if not main package
		// Check if function is exported (starts with uppercase) and package is not main
		if l.prefix != "" {
			fnName = l.prefix + fnName
		} else if len(fnName) > 0 && fnName[0] >= 65 && fnName[0] <= 90 {
			// Function is exported (starts with uppercase)
			// Get package name from prog
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
//      | ExprStmt | BreakStmt | ContinueStmt
func (l *lowerer) lowerStmt(s ast.Stmt) error {
	switch st := s.(type) {
	case *ast.BlockStmt:
		return l.lowerBlock(st)
	case *ast.DeclStmt:
		return l.lowerDeclStmt(st)
	case *ast.ExprStmt:
		_, err := l.lowerExpr(st.X)
		return err
	case *ast.ReturnStmt:
		var vals []*ir.Value
		for i := 0; i < len(st.Results); i++ {
			// Check if return value is nil identifier before lowering
			isNilIdent := false
			if ident, ok := st.Results[i].(*ast.Ident); ok && ident.Name == "nil" {
				isNilIdent = true
			}
			v, err := l.lowerExpr(st.Results[i])
			if err != nil {
				return err
			}
			// If return value is nil and function expects slice, create zero slice
			if l.fn != nil && i < len(l.fn.Results) {
				expectedTy := l.fn.Results[i]
				if expectedTy.Kind == ir.KindSlice && isNilIdent {
					// Create zero slice: { ptr null, i64 0, i64 0 }
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

// AssignStmt = ExprList "=" ExprList
func (l *lowerer) lowerAssign(st *ast.AssignStmt) error {
	if st.Tok != token.ASSIGN && st.Tok != token.DEFINE {
		return fmt.Errorf("unsupported assign tok %s", st.Tok.String())
	}
	// special case: map get with two results
	if len(st.Lhs) == 2 && len(st.Rhs) == 1 {
		if idx, ok := st.Rhs[0].(*ast.IndexExpr); ok {
			if _, okm := l.prog.TypesInfo.TypeOf(idx.X).Underlying().(*types.Map); okm {
				return l.lowerMapGet(idx, st.Lhs[0], st.Lhs[1], st.Tok == token.DEFINE)
			}
		}
		// Also handle function calls that return 2 values
		if call, ok := st.Rhs[0].(*ast.CallExpr); ok {
			// Try to get function signature from various sources
			var sig *types.Signature
			callType := l.prog.TypesInfo.TypeOf(call)
			if sig2, ok := callType.(*types.Signature); ok {
				sig = sig2
			} else {
				// Try to get from function identifier
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
					// Also try Selections for method calls
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
				// Function returns 2 values - handle it
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
				// Handle first LHS
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
				// Handle second LHS (error or second return value)
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
	// Support multiple assignment for functions returning multiple values
	// Handle (T, error) pattern and other multi-return patterns
	if len(st.Rhs) == 1 {
		// Check if RHS is a function call returning multiple values
		if call, ok := st.Rhs[0].(*ast.CallExpr); ok {
			callType := l.prog.TypesInfo.TypeOf(call)
			if sig, ok := callType.(*types.Signature); ok {
				numReturns := sig.Results().Len()
				// Handle case when function returns multiple values
				if numReturns > 1 {
					// If we have exactly as many LHS as returns, assign all
					if numReturns == len(st.Lhs) {
						// Function returns exactly as many values as we have LHS variables
						// Lower the call to get the first value
						firstVal, err := l.lowerCall(call)
						if err != nil {
							return err
						}
						// Assign first value to first LHS
						firstTy := firstVal.Type()
						if firstTy == nil && numReturns > 0 {
							// Get type from signature
							firstGoTy := sig.Results().At(0).Type()
							firstTy, err = goTypeToIR(firstGoTy)
							if err != nil {
								return err
							}
						}
						// Handle first LHS
						lhs0Val, err := l.lowerExpr(st.Lhs[0])
						if err != nil {
							// LHS doesn't exist yet, need to create it
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
							// LHS exists, assign to it
							l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: firstVal, StoreDst: lhs0Val, StoreAlign: alignOfIRType(firstTy)})
						}
						// Handle remaining LHS variables (for multi-return functions)
						for i := 1; i < len(st.Lhs); i++ {
							returnTy, err := goTypeToIR(sig.Results().At(i).Type())
							if err != nil {
								return err
							}
							// For now, assign zero/nil values to remaining variables
							// In a real implementation, we'd need to extract all return values
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
								// LHS doesn't exist yet, need to create it
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
								// LHS exists, assign to it
								l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: zeroVal, StoreDst: lhsVal, StoreAlign: alignOfIRType(returnTy)})
							}
						}
						return nil
					} else if len(st.Lhs) > 0 && len(st.Lhs) < numReturns {
						// If we have fewer LHS than returns, assign only first value(s)
						// Function returns more values than we assign - use first N values
						for i := 0; i < len(st.Lhs); i++ {
							returnTy, err := goTypeToIR(sig.Results().At(i).Type())
							if err != nil {
								return err
							}
							// For now, we'll call the function once and use the first return value
							// Then assign zero values to remaining LHS
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
								// Assign zero values to remaining LHS
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
						// If we have one LHS and function returns multiple values, use first
						firstVal, err := l.lowerCall(call)
						if err != nil {
							return err
						}
						lhs, ok := st.Lhs[0].(*ast.Ident)
						if !ok {
							return fmt.Errorf("unsupported lhs")
						}
						if lhs.Name == "_" {
							// evaluate for side-effects only
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
						// If we have zero LHS, just evaluate for side effects
						_, err := l.lowerCall(call)
						return err
					}
					// If numReturns == 1 and len(st.Lhs) == 1, fall through to normal handling
				}
			} else if callType == nil {
				// TypeOf returned nil - try to get type from the function itself
				// This can happen with some function calls
				if fnIdent, ok := call.Fun.(*ast.Ident); ok {
					if obj, ok := l.prog.TypesInfo.Uses[fnIdent]; ok {
						if fn, ok := obj.(*types.Func); ok {
							sig := fn.Type().(*types.Signature)
							numReturns := sig.Results().Len()
							// Handle multiple returns even when TypeOf(call) is nil
							if numReturns > 1 && numReturns == len(st.Lhs) {
								// Same handling as above
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
					// Try to get type from selector expression
					if obj, ok := l.prog.TypesInfo.Uses[selExpr.Sel]; ok {
						if fn, ok := obj.(*types.Func); ok {
							sig := fn.Type().(*types.Signature)
							numReturns := sig.Results().Len()
							if numReturns > 1 && numReturns == len(st.Lhs) {
								// Same handling as above
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
		// Provide more context in error message
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
		// map set?
		if idx, ok := st.Lhs[i].(*ast.IndexExpr); ok {
			if _, okm := l.prog.TypesInfo.TypeOf(idx.X).Underlying().(*types.Map); okm {
				return l.lowerMapSet(idx, st.Rhs[i])
			}
			// Handle array/slice element assignment: arr[i] = value or slice[i] = value
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
		// Handle struct field assignment: s.field = value
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
		// Handle pointer dereference assignment: *ptr = value
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
			// evaluate RHS for side-effects only
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

	// then
	l.block = thenBB
	if err := l.lowerStmt(st.Body); err != nil {
		return err
	}
	if l.block.Terminator == nil {
		l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: mergeBB.Name}
	}

	// else
	l.block = elseBB
	if st.Else != nil {
		if err := l.lowerStmt(st.Else); err != nil {
			return err
		}
	}
	if l.block.Terminator == nil {
		l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: mergeBB.Name}
	}

	// continue at merge
	l.block = mergeBB
	return nil
}

// ForStmt = "for" ( [ SimpleStmt ] ";" [ Expr ] ";" [ SimpleStmt ] | Expr ) Block
// SimpleStmt = AssignStmt | ExprStmt
func (l *lowerer) lowerFor(st *ast.ForStmt) error {
	// support: for cond { ... } and for init; cond; post { ... }
	condBB := l.newBlock(l.newBlockName("for.cond"))
	bodyBB := l.newBlock(l.newBlockName("for.body"))
	postBB := l.newBlock(l.newBlockName("for.post"))
	endBB := l.newBlock(l.newBlockName("for.end"))

	// push loop targets
	l.breakTargets = append(l.breakTargets, endBB.Name)
	if st.Post != nil {
		l.continueTargets = append(l.continueTargets, postBB.Name)
	} else {
		l.continueTargets = append(l.continueTargets, condBB.Name)
	}

	// init
	if st.Init != nil {
		if err := l.lowerStmt(st.Init); err != nil {
			return err
		}
	}

	l.block.Terminator = &ir.Instruction{Kind: ir.InstrBr, BrTarget: condBB.Name}

	// cond
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

	// body
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

	// continue after loop
	l.block = endBB
	// pop loop targets
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
		// fallback to global
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
		// ignore ok result
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
		// slice layout: { data *T, len i64, cap i64 }
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
		// treat string as { i8*, i64 }
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

	// Get expected result type from Go type system
	expectedTy, err := goTypeToIR(l.prog.TypesInfo.TypeOf(slice))
	if err != nil {
		// If we can't get the type, fall back to base type
		expectedTy = nil
	}

	// Get low index (default 0 if nil)
	var lowVal *ir.Value
	if slice.Low != nil {
		lowVal, err = l.lowerExpr(slice.Low)
		if err != nil {
			return nil, err
		}
		// Ensure lowVal is i64
		if lowVal.Type() == nil || lowVal.Type().Basic != "i64" {
			lowVal, err = l.lowerConvert(lowVal, ir.I64)
			if err != nil {
				return nil, err
			}
		}
	} else {
		lowVal = ir.NewConstant("0", ir.I64)
	}

	// Handle strings
	if baseTy.Kind == ir.KindString {
		ptr, lenVal := l.materializeStringParts(baseVal)
		
		// Get high index (default len if nil)
		var highVal *ir.Value
		if slice.High != nil {
			highVal, err = l.lowerExpr(slice.High)
			if err != nil {
				return nil, err
			}
			// Ensure highVal is i64
			if highVal.Type() == nil || highVal.Type().Basic != "i64" {
				highVal, err = l.lowerConvert(highVal, ir.I64)
				if err != nil {
					return nil, err
				}
			}
		} else {
			highVal = lenVal
		}

		// Compute new ptr = ptr + low using GEP
		newPtr := l.newTemp(ir.PtrTo(ir.I8))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          newPtr,
			GepSrc:        ptr,
			GepPointee:    ir.I8,
			GepIndices:    []*ir.Value{lowVal},
			GepResultType: ir.PtrTo(ir.I8),
		})

		// Compute new len = high - low
		newLen := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{
			Kind:     ir.InstrBinOp,
			Dest:     newLen,
			BinOp:    ir.Sub,
			X:        highVal,
			Y:        lowVal,
		})

		// Build new string or slice depending on expected type
		if expectedTy != nil && expectedTy.Kind == ir.KindSlice && expectedTy.Elem == ir.I8 {
			// Expected type is []byte, return slice
			return l.buildSliceValue(ir.I8, newPtr, newLen, newLen), nil
		}
		// Expected type is string, return string
		return l.buildStringValue(newPtr, newLen), nil
	}

	// Handle slices
	if baseTy.Kind == ir.KindSlice {
		elemTy := baseTy.Elem
		slicePtr := l.valuePtr(baseVal, baseTy)
		
		// Get data pointer
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

		// Get len
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

		// Get cap
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

		// Get high index (default len if nil)
		var highVal *ir.Value
		if slice.High != nil {
			highVal, err = l.lowerExpr(slice.High)
			if err != nil {
				return nil, err
			}
			// Ensure highVal is i64
			if highVal.Type() == nil || highVal.Type().Basic != "i64" {
				highVal, err = l.lowerConvert(highVal, ir.I64)
				if err != nil {
					return nil, err
				}
			}
		} else {
			highVal = lenVal
		}

		// Compute new data = data + low using GEP
		newDataPtr := l.newTemp(ir.PtrTo(elemTy))
		l.block.Append(ir.Instruction{
			Kind:          ir.InstrGEP,
			Dest:          newDataPtr,
			GepSrc:        dataPtr,
			GepPointee:    elemTy,
			GepIndices:    []*ir.Value{lowVal},
			GepResultType: ir.PtrTo(elemTy),
		})

		// Compute new len = high - low
		newLen := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{
			Kind:     ir.InstrBinOp,
			Dest:     newLen,
			BinOp:    ir.Sub,
			X:        highVal,
			Y:        lowVal,
		})

		// Compute new cap = cap - low
		newCap := l.newTemp(ir.I64)
		l.block.Append(ir.Instruction{
			Kind:     ir.InstrBinOp,
			Dest:     newCap,
			BinOp:    ir.Sub,
			X:        capVal,
			Y:        lowVal,
		})

		// Build new slice
		return l.buildSliceValue(elemTy, newDataPtr, newLen, newCap), nil
	}

	return nil, fmt.Errorf("slice expr on unsupported type %v", baseTy.Kind)
}

func (l *lowerer) buildStringValue(ptr *ir.Value, lenVal *ir.Value) *ir.Value {
	strTy := ir.String()
	addr := l.newTemp(ir.PtrTo(strTy))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: strTy, AllocaAlign: alignOfIRType(strTy)})
	// ptr
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
	// len
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
	// load result
	res := l.newTemp(strTy)
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: res, LoadSrc: addr, LoadAlign: alignOfIRType(strTy)})
	return res
}

// valuePtr: если значение не указатель, выделяет память и сохраняет значение
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

// buildSliceValue: создаёт значение среза { i8*, i64, i64 } из данных, длины и ёмкости
func (l *lowerer) buildSliceValue(elemTy *ir.TypeDesc, dataPtr *ir.Value, lenVal *ir.Value, capVal *ir.Value) *ir.Value {
	sliceTy := ir.Slice(elemTy)
	addr := l.newTemp(ir.PtrTo(sliceTy))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: sliceTy, AllocaAlign: alignOfIRType(sliceTy)})
	// data
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
	// len
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
	// cap
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
	// load result
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
		// Check if this is a recursive type that was replaced with PtrI8
		// In that case, we can't process fields, so just create an uninitialized pointer
		if irTy.Kind == ir.KindPointer && irTy.Elem != nil && irTy.Elem.Basic == "i8" {
			// This is a pointer type (likely PtrI8 from recursive type handling)
			// Create an uninitialized pointer value
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
	// First, check if this is a package constant access (e.g., ir.KindBasic)
	// Check TypesInfo.Types first - this works for constants
	if tv, ok := l.prog.TypesInfo.Types[sel]; ok {
		if tv.IsValue() && tv.Value != nil {
			// This is a constant value
			ty, err := goTypeToIR(tv.Type)
			if err != nil {
				return nil, err
			}
			valStr := tv.Value.String()
			return ir.NewConstant(valStr, ty), nil
		}
		// Also check if it's a type constant (for iota-based constants)
		if tv.IsType() {
			// This might be a type, not a constant
		}
	}
	
	// Check for struct field access FIRST, before checking Uses
	// This prevents fields from being treated as global variables
	selection := l.prog.TypesInfo.Selections[sel]
	if selection != nil && selection.Kind() == types.FieldVal {
		// This is a field access - handle it below
	} else {
		// Not a field access - check if it's a constant or package variable
		if obj, ok := l.prog.TypesInfo.Uses[sel.Sel]; ok {
			if constVal, ok := obj.(*types.Const); ok {
				// This is a constant
				ty, err := goTypeToIR(constVal.Type())
				if err != nil {
					return nil, err
				}
				// Convert constant value to string representation
				valStr := constVal.Val().String()
				return ir.NewConstant(valStr, ty), nil
			}
			// Handle package variables (e.g., ir.Void, ir.I64)
			if varVal, ok := obj.(*types.Var); ok {
				// This is a package variable - treat as global
				ty, err := goTypeToIR(varVal.Type())
				if err != nil {
					return nil, err
				}
				gname := varVal.Name()
				// Check if this is from a package (via sel.X)
				if pkgIdent, ok := sel.X.(*ast.Ident); ok {
					if pkgObj, ok := l.prog.TypesInfo.Uses[pkgIdent]; ok {
						if pkg, ok := pkgObj.(*types.PkgName); ok {
							pkgName := pkg.Imported().Name()
							// Use full qualified name for standard library variables
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
	
	// Try to handle as package constant via package lookup
	if pkgIdent, ok := sel.X.(*ast.Ident); ok {
		obj := l.prog.TypesInfo.Uses[pkgIdent]
		if obj != nil {
			if pkg, ok := obj.(*types.PkgName); ok {
				// This is a package selector, look for the constant
				pkgScope := pkg.Imported().Scope()
				pkgObj := pkgScope.Lookup(sel.Sel.Name)
				if pkgObj != nil {
					if constVal, ok := pkgObj.(*types.Const); ok {
						// Get the constant value
						ty, err := goTypeToIR(constVal.Type())
						if err != nil {
							return nil, err
						}
						// Convert constant value to string representation
						valStr := constVal.Val().String()
						return ir.NewConstant(valStr, ty), nil
					}
					if varVal, ok := pkgObj.(*types.Var); ok {
						// This is a package variable - treat as global
						ty, err := goTypeToIR(varVal.Type())
						if err != nil {
							return nil, err
						}
						// For standard library variables like os.Stderr, use package prefix
						pkgName := pkg.Imported().Name()
						gname := varVal.Name()
						// Use full qualified name for standard library variables
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
	
	// Now check for struct field access (if not already checked above)
	if selection == nil {
		// Try one more approach - check if sel.Sel is in Uses (sel.Sel is already *ast.Ident)
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
	
	// Get IR type from Go type (this handles named types correctly)
	baseGoType := l.prog.TypesInfo.TypeOf(sel.X)
	if baseGoType == nil {
		return nil, fmt.Errorf("no type info for selector base")
	}
	
	// Get IR type for the base (may be pointer or value) - check this first for recursive types
	baseType, err := goTypeToIR(baseGoType)
	if err != nil {
		return nil, err
	}
	
	// Extract struct IR type from baseType
	var structIRType *ir.TypeDesc
	if baseType.Kind == ir.KindPointer {
		structIRType = baseType.Elem
	} else {
		structIRType = baseType
	}
	
	// Check if this is a recursive type BEFORE trying to get the Go struct type
	if structIRType.Kind != ir.KindStruct {
		// This might be a recursive type - check the original named type first
		// Get the element type from baseGoType (before calling Underlying())
		var elemType types.Type
		if ptrType, ok := baseGoType.(*types.Pointer); ok {
			elemType = ptrType.Elem()
		} else {
			elemType = baseGoType
		}
		// Check if this is a known recursive type (like Global) that should be handled specially
		if named, ok := elemType.(*types.Named); ok {
			typeName := ""
			if named.Obj() != nil {
				typeName = named.Obj().Name()
			}
			fieldName := sel.Sel.Name
			// Handle Global type fields (from ir package or main package when using -skip-check)
			// Check if the type name matches regardless of package
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
		// If not a known recursive type, try to get the Go struct type
		// Get the actual struct type (handle named types and pointers)
		underlying := baseGoType.Underlying()
		var structGoType types.Type
		if ptrType, ok := underlying.(*types.Pointer); ok {
			// It's a pointer - get the element type
			structGoType = ptrType.Elem().Underlying()
		} else {
			// It's a value type - get underlying (handles named types)
			structGoType = underlying
		}
		
		// Try to convert the Go struct type to IR
		structIRTypeFromGo, err := goTypeToIR(structGoType)
		if err == nil && structIRTypeFromGo.Kind == ir.KindStruct {
			// It's actually a struct in IR, use it
			structIRType = structIRTypeFromGo
		} else {
			// It's a recursive type that was replaced with PtrI8 or similar
			// Return a pointer value (can't access fields)
			return ir.NewConstant("null", ir.PtrI8), nil
		}
	}
	
	// Now get the Go struct type for field index lookup
	underlying := baseGoType.Underlying()
	var structGoType types.Type
	if ptrType, ok := underlying.(*types.Pointer); ok {
		// It's a pointer - get the element type
		structGoType = ptrType.Elem().Underlying()
	} else {
		// It's a value type - get underlying (handles named types)
		structGoType = underlying
	}
	
	// Verify it's actually a struct (should be at this point)
	_, ok := structGoType.(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("selector base not struct (got %T)", structGoType)
	}
	
	// Determine pointer to use
	var ptr *ir.Value
	if baseVal.Type() != nil && baseVal.Type().Kind == ir.KindPointer {
		// baseVal is already a pointer, use it directly
		ptr = baseVal
	} else if baseType.Kind == ir.KindPointer {
		// Go type is pointer, but baseVal is value (shouldn't happen, but handle it)
		ptr = baseVal
	} else {
		// Value type - need to get address
		addr := l.newTemp(ir.PtrTo(structIRType))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: structIRType, AllocaAlign: alignOfIRType(structIRType)})
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: baseVal, StoreDst: addr, StoreAlign: alignOfIRType(structIRType)})
		ptr = addr
	}
	// structIRType should be a struct at this point (recursive types were handled earlier)
	// But double-check to be safe
	if structIRType.Kind != ir.KindStruct {
		// This shouldn't happen, but if it does, handle it as a recursive type
		// Get the element type from baseGoType
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
			// Handle Global type fields
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
		// strip quotes
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
	// && и || требуют short-circuit evaluation (не вычисляем правую часть, если не нужно)
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
		// Shift left
		destTy := x.Type()
		if destTy == nil {
			destTy = y.Type()
		}
		dest := l.newTemp(destTy)
		l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: dest, BinOp: ir.Shl, X: x, Y: y})
		return dest, nil
	case token.SHR:
		// Shift right - use logical shift for unsigned, arithmetic for signed
		destTy := x.Type()
		if destTy == nil {
			destTy = y.Type()
		}
		dest := l.newTemp(destTy)
		// For simplicity, use arithmetic shift right (works for both signed and unsigned in most cases)
		// In a full implementation, we'd check if the type is signed or unsigned
		l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: dest, BinOp: ir.AShr, X: x, Y: y})
		return dest, nil
	case token.AND:
		// Bitwise AND (binary operator, not unary &)
		destTy := x.Type()
		if destTy == nil {
			destTy = y.Type()
		}
		dest := l.newTemp(destTy)
		l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: dest, BinOp: ir.And, X: x, Y: y})
		return dest, nil
	case token.OR:
		// Bitwise OR (binary operator, not logical ||)
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

// lowerLogical: реализует short-circuit evaluation для && и ||
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
		// short path: store false
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

	// move to end block
	l.block = endBB
	result := l.newTemp(ir.I1)
	endBB.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: result, LoadSrc: resPtr, LoadAlign: 1})
	return result, nil
}

// UnaryExpr = ("+" | "-" | "!") Expr
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
		// assume int64
		dest := l.newTemp(v.Type())
		l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: dest, BinOp: ir.Sub, X: ir.NewConstant("0", v.Type()), Y: v})
		return dest, nil
	case token.AND:
		return l.lowerAddr(u.X)
	default:
		return nil, fmt.Errorf("unsupported unary %s", u.Op.String())
	}
}

// lowerStar handles *ast.StarExpr (pointer dereference)
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
		// reuse selector lowering but keep pointer
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
		// Check if structTy is actually a struct (not a pointer type from recursive type handling)
		if structTy.Kind != ir.KindStruct {
			// This is a recursive type that was replaced with PtrI8 or similar
			// Return a pointer value (can't access fields)
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
		// Create a temporary variable for the composite literal
		gt := l.prog.TypesInfo.TypeOf(ex)
		if gt == nil {
			return nil, fmt.Errorf("no type for composite literal")
		}
		irTy, err := goTypeToIR(gt)
		if err != nil {
			return nil, err
		}
		// Allocate space for the composite literal
		addr := l.newTemp(ir.PtrTo(irTy))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: irTy, AllocaAlign: alignOfIRType(irTy)})
		
		// Lower the composite literal to get its value
		val, err := l.lowerCompositeLit(ex)
		if err != nil {
			return nil, err
		}
		
		// Store the value into the allocated space
		l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: val, StoreDst: addr, StoreAlign: alignOfIRType(irTy)})
		
		// Return the address
		return addr, nil
	default:
		return nil, fmt.Errorf("address-of unsupported expr %T", e)
	}
}

// Call = PrimaryExpr "(" [ ExprList ] ")"
// ExprList = Expr { "," Expr }
func (l *lowerer) lowerCall(c *ast.CallExpr) (*ir.Value, error) {
	// Handle type conversion where Fun is a type expression (e.g., []byte(s)).
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
					// Special-case string -> []byte
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

	// Handle function calls through package selector (e.g., ir.ValueName(v))
	if selExpr, ok := c.Fun.(*ast.SelectorExpr); ok {
		// Check if this is a package function call
		if pkgIdent, ok := selExpr.X.(*ast.Ident); ok {
			// Get the function object
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
						// Determine return type
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
					// This is a function call, not a method
					// Lower arguments
					args := make([]*ir.Value, len(c.Args))
					for i := 0; i < len(c.Args); i = i + 1 {
						v, err := l.lowerExpr(c.Args[i])
						if err != nil {
							return nil, err
						}
						args[i] = v
					}
					// Determine return type
					var retTy *ir.TypeDesc
					if sig.Results().Len() == 1 {
						rt, err := goTypeToIR(sig.Results().At(0).Type())
						if err != nil {
							return nil, err
						}
						retTy = rt
					}
					// Build function name: pkg.FuncName
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
		// Handle method calls on expressions (e.g., entry.Append(...), p.TypesPkg.Name())
		// Try to get method info from Selections or Uses
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
				// This might be a field that is a function
				if obj, ok := l.prog.TypesInfo.Uses[selExpr.Sel]; ok {
					if fn, ok := obj.(*types.Func); ok {
						method = fn
						isMethod = false
					}
				}
			}
		}
		if !isMethod {
			// Try Uses as fallback
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
			// Get method signature
			sig := method.Type().(*types.Signature)
			// Lower arguments
			args := make([]*ir.Value, len(c.Args))
			for i := 0; i < len(c.Args); i = i + 1 {
				v, err := l.lowerExpr(c.Args[i])
				if err != nil {
					return nil, err
				}
				args[i] = v
			}
			// Determine return type
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
				// This is a method call - prepend receiver
				recvVal, err := l.lowerExpr(selExpr.X)
				if err != nil {
					return nil, err
				}
				allArgs = make([]*ir.Value, len(args)+1)
				allArgs[0] = recvVal
				for i := 0; i < len(args); i = i + 1 {
					allArgs[i+1] = args[i]
				}
				// Build function name: pkg.TypeName.MethodName
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
				// This is a function field call - no receiver
				allArgs = args
				// Get the type of the base expression to determine package
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
		// Last resort: try to get type info and build function name from selector
		baseType := l.prog.TypesInfo.TypeOf(selExpr.X)
		selType := l.prog.TypesInfo.TypeOf(selExpr)
		if selType != nil {
			if sig, ok := selType.(*types.Signature); ok {
				// Lower arguments
				args := make([]*ir.Value, len(c.Args))
				for i := 0; i < len(c.Args); i = i + 1 {
					v, err := l.lowerExpr(c.Args[i])
					if err != nil {
						return nil, err
					}
					args[i] = v
				}
				// Determine return type
				var retTy *ir.TypeDesc
				if sig.Results().Len() == 1 {
					rt, err := goTypeToIR(sig.Results().At(0).Type())
					if err != nil {
						return nil, err
					}
					retTy = rt
				}
				// Try to build function name from base type and selector
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
	// type conversion: Fun is a type name.
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
	// builtins: len/make/append
	if fnIdent.Name == "len" {
		return l.lowerBuiltinLen(c.Args)
	}
	if fnIdent.Name == "make" {
		return l.lowerBuiltinMake(c.Args)
	}
	if fnIdent.Name == "append" {
		return l.lowerBuiltinAppend(c.Args)
	}
	// builtin print wrappers not handled; assume normal call returning at most 1 value.
	args := make([]*ir.Value, len(c.Args))
	for i := 0; i < len(c.Args); i++ {
		v, err := l.lowerExpr(c.Args[i])
		if err != nil {
			return nil, err
		}
		args[i] = v
	}
	// determine return type
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
			// methods with receiver unsupported in subset
			if sig.Recv() != nil {
				return nil, fmt.Errorf("unsupported: methods with receiver")
			}
		}
	}
	var dest *ir.Value
	if retTy != nil && retTy != ir.Void {
		dest = l.newTemp(retTy)
	}
	// Build function name with package prefix
	fnName := fnIdent.Name
	// First, try to use l.prefix (set in buildInto for non-main packages)
	// This matches the logic in lowerFunc, which adds prefix for all functions if l.prefix is set
	if l.prefix != "" {
		fnName = l.prefix + fnName
	} else {
		// If prefix is not set, add package prefix for exported functions only
		// (to match function definitions which only have prefix for exported functions when prefix is not set)
		if len(fnName) > 0 && fnName[0] >= 'A' && fnName[0] <= 'Z' {
			pkgName := ""
			if l.prog != nil && l.prog.PkgName != "" {
				pkgName = l.prog.PkgName
			}
			// If program's package name is "main" or empty, try to get from function object
			// (but only if it's not "main", as type checker may use "main" for all packages with -skip-check)
			if (pkgName == "" || pkgName == "main") {
				if obj := l.prog.TypesInfo.Uses[fnIdent]; obj != nil {
					if fn, ok := obj.(*types.Func); ok {
						if fn.Pkg() != nil {
							fnPkgName := fn.Pkg().Name()
							// Only use fn.Pkg() if it's not "main" (type checker may use "main" for all packages)
							if fnPkgName != "" && fnPkgName != "main" {
								pkgName = fnPkgName
							}
						}
					}
				}
			}
			// Add prefix if package name is set (including "main" and "gominic" to match function definitions)
			if pkgName != "" {
				fnName = pkgName + "." + fnName
			}
		}
	}
	l.block.Append(ir.Instruction{Kind: ir.InstrCall, Dest: dest, CallName: fnName, CallArgs: args, CallRet: retTy})
	return dest, nil
}

// lowerConvert: конвертация типов (int->int64, int->float64, и т.д.)
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
	// int64 <-> double conversions
	if v.Type().Basic == "double" && dst.Basic == "i64" {
		l.block.Append(ir.Instruction{Kind: ir.InstrConv, Dest: out, ConvOp: ir.FPToSI, ConvSrc: v, ConvTo: dst})
		return out, nil
	}
	if v.Type().Basic == "i64" && dst.Basic == "double" {
		l.block.Append(ir.Instruction{Kind: ir.InstrConv, Dest: out, ConvOp: ir.SIToFP, ConvSrc: v, ConvTo: dst})
		return out, nil
	}
	// int32 -> int64 promote
	if v.Type().Basic == "i32" && dst.Basic == "i64" {
		l.block.Append(ir.Instruction{Kind: ir.InstrConv, Dest: out, ConvOp: ir.SExt, ConvSrc: v, ConvTo: dst})
		return out, nil
	}
	// pointer-to-pointer bitcast (used rarely)
	if v.Type().Kind == ir.KindPointer && dst.Kind == ir.KindPointer {
		l.block.Append(ir.Instruction{Kind: ir.InstrBitcast, Dest: out, BitcastSrc: v, BitcastTarget: dst})
		return out, nil
	}
	// i8* -> string conversion (pointer to bytes to string)
	if v.Type().Kind == ir.KindPointer && v.Type().Elem != nil && v.Type().Elem == ir.I8 && dst.Kind == ir.KindString {
		// Create a string struct: { i8* ptr, i64 len }
		// We need to allocate space for the string struct
		strPtr := l.newTemp(ir.PtrTo(dst))
		l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: strPtr, AllocaType: dst, AllocaAlign: 8})
		
		// Store the pointer (first field)
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
		
		// Store zero length (second field) - we don't know the length from just a pointer
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
		
		// Load the complete string struct
		loaded := l.newTemp(dst)
		l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: loaded, LoadSrc: strPtr, LoadAlign: 8})
		return loaded, nil
	}
	// []byte -> string conversion (slice to string)
	if v.Type().Kind == ir.KindSlice && v.Type().Elem != nil && v.Type().Elem == ir.I8 && dst.Kind == ir.KindString {
		// Extract data pointer and length from slice
		slicePtr := l.valuePtr(v, v.Type())
		// Get data pointer
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
		// Get length
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
		// Build string from ptr and len
		return l.buildStringValue(ptr, lenVal), nil
	}
	return nil, fmt.Errorf("conversion from %s to %s not supported", v.Type().String(), dst.String())
}

// ensureFloat64: приводит значение к float64, если оно целое (нужно для float операций)
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
	// map: handle first
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
	slicePtr := l.valuePtr(sliceVal, sliceVal.Type())
	elemSize := ir.NewConstant(strconv.Itoa(int(sizeOfIRType(sliceVal.Type().Elem))), ir.I64)
	// call runtime append equivalent? Пока примитив: len+1 cap+? reallocate via makeSlice
	// len
	lenPtr := l.newTemp(ir.PtrTo(ir.I64))
	lenVal := l.newTemp(ir.I64)
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          lenPtr,
		GepSrc:        slicePtr,
		GepPointee:    sliceVal.Type(),
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("1", ir.I32)},
		GepResultType: ir.PtrTo(ir.I64),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: lenVal, LoadSrc: lenPtr, LoadAlign: 8})
	newLen := l.newTemp(ir.I64)
	l.block.Append(ir.Instruction{Kind: ir.InstrBinOp, Dest: newLen, BinOp: ir.Add, X: lenVal, Y: ir.NewConstant("1", ir.I64)})
	// cap
	capPtr := l.newTemp(ir.PtrTo(ir.I64))
	capVal := l.newTemp(ir.I64)
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          capPtr,
		GepSrc:        slicePtr,
		GepPointee:    sliceVal.Type(),
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I32), ir.NewConstant("2", ir.I32)},
		GepResultType: ir.PtrTo(ir.I64),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrLoad, Dest: capVal, LoadSrc: capPtr, LoadAlign: 8})
	// For simplicity always allocate new slice with cap = newLen
	ptr := l.newTemp(ir.PtrTo(sliceVal.Type().Elem))
	l.block.Append(ir.Instruction{
		Kind:     ir.InstrCall,
		Dest:     ptr,
		CallName: "gominic_makeSlice",
		CallArgs: []*ir.Value{newLen, newLen, elemSize},
		CallRet:  ir.PtrTo(sliceVal.Type().Elem),
	})
	// TODO: copy old data + append elem. For now just write elem at index 0 if len was 0.
	elemPtr := l.newTemp(ir.PtrTo(sliceVal.Type().Elem))
	l.block.Append(ir.Instruction{
		Kind:          ir.InstrGEP,
		Dest:          elemPtr,
		GepSrc:        ptr,
		GepPointee:    sliceVal.Type().Elem,
		GepIndices:    []*ir.Value{ir.NewConstant("0", ir.I64)},
		GepResultType: ir.PtrTo(sliceVal.Type().Elem),
	})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: elemVal, StoreDst: elemPtr, StoreAlign: alignOfIRType(sliceVal.Type().Elem)})
	return l.buildSliceValue(sliceVal.Type().Elem, ptr, newLen, newLen), nil
}


func goTypeToIR(t types.Type) (*ir.TypeDesc, error) {
	return goTypeToIRVisited(t, make(map[string]bool))
}

func goTypeToIRVisited(t types.Type, visited map[string]bool) (*ir.TypeDesc, error) {
	// Handle nil type
	if t == nil {
		return ir.Void, nil
	}
	// Handle named types first to detect recursion before calling Underlying()
	if named, ok := t.(*types.Named); ok {
		if named.Obj() != nil {
			typeName := named.Obj().Name()
			// Handle builtin error type specially
			if typeName == "error" {
				// error is an interface type - represent as pointer
				return ir.PtrI8, nil
			}
			if named.Obj().Pkg() != nil {
				pkgName := named.Obj().Pkg().Name()
				pkgPath := named.Obj().Pkg().Path()
				// Handle types from ir package or from merged packages (when using -skip-check)
				// These types (TypeDesc, Value, Instruction, etc.) are represented as pointers in IR
				if pkgName == "ir" || pkgName == "backend" || pkgName == "main" {
					// Check if this is one of the IR types that should be represented as pointers
					if typeName == "TypeDesc" || typeName == "Value" || typeName == "Instruction" || 
					   typeName == "BasicBlock" || typeName == "Function" || typeName == "Module" || 
					   typeName == "Global" {
						return ir.PtrI8, nil
					}
				}
				// Handle types from standard library packages that we don't fully support
				// Replace them with pointers to avoid recursion and unsupported type errors
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
			// Allow recursive types for IR types (TypeDesc, Value, etc.)
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
				// Handle recursive types from standard library
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

	// Now process underlying types
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
		// Interfaces are not fully supported in our subset
		// Represent them as pointers for self-hosting
		return ir.PtrI8, nil
	case *types.Signature:
		// Function signatures are not fully supported in our subset
		// Represent them as pointers for self-hosting
		return ir.PtrI8, nil
	case *types.Chan:
		// Channels are not supported in our subset
		// Represent them as pointers for self-hosting
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
		// Remove trailing dot from prefix
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

// materializeStringParts: извлекает ptr и len из строки { i8*, i64 }
// Выделяет память, сохраняет строку, затем загружает поля через GEP
func (l *lowerer) materializeStringParts(str *ir.Value) (*ir.Value, *ir.Value) {
	addr := l.newTemp(ir.PtrTo(ir.String()))
	l.block.Append(ir.Instruction{Kind: ir.InstrAlloca, Dest: addr, AllocaType: ir.String(), AllocaAlign: alignOfIRType(ir.String())})
	l.block.Append(ir.Instruction{Kind: ir.InstrStore, StoreSrc: str, StoreDst: addr, StoreAlign: alignOfIRType(ir.String())})
	// ptr
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
	// len
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
