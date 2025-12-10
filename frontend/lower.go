package frontend

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"gominic/ir"
)

// BuildIR превращает проверенный Program в наш IR.
// Поддержка: функции без вызовов, локальные переменные (var и :=), присваивания,
// бинарные арифметические и сравнительные операции, if/else, return с 0 или 1 значением.
/*
func BuildIR(prog *Program) (*ir.Module, error) {
	l := &lowerer{
		prog: prog,
		mod: &ir.Module{
			TargetTriple: "x86_64-pc-linux-gnu",
			DataLayout:   "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128",
		},
	}

	for _, file := range prog.Files {
		if err := l.lowerGlobals(file.Decls); err != nil {
			return nil, err
		}
	}

	for _, file := range prog.Files {
		for _, decl := range file.Decls {
			fnDecl, ok := decl.(*ast.FuncDecl)
			if !ok || fnDecl.Recv != nil {
				continue
			}
			if err := l.lowerFunc(fnDecl); err != nil {
				return nil, err
			}
		}
	}

	return l.mod, nil
}
*/

func (l *lowerer) lowerGlobals(decls []ast.Decl) error {
	for _, decl := range decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gen.Tok == token.VAR || gen.Tok == token.CONST {
			for _, spec := range gen.Specs {
				vSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				if err := l.lowerGlobalSpec(vSpec, gen.Tok); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (l *lowerer) lowerGlobalSpec(spec *ast.ValueSpec, tok token.Token) error {
	for i, name := range spec.Names {
		obj := l.prog.TypesInfo.Defs[name]
		if obj == nil {
			return fmt.Errorf("no def info for global %s", name.Name)
		}
		irType, err := goTypeToIR(obj.Type())
		if err != nil {
			return err
		}
		fullName := name.Name
		if l.prefix != "" {
			fullName = l.prefix + name.Name
		}
		var init ir.Value
		if i < len(spec.Values) {
			if c, errConst := l.constValue(spec.Values[i]); errConst == nil {
				init = c
			}
		}
		if init.Kind() == ir.ValueInvalid {
			init, err = zeroValue(obj.Type())
			if err != nil {
				return err
			}
		}
		valStr := "zeroinitializer"
		if s, errStr := constantString(init); errStr == nil {
			valStr = s
		}
		align := alignOfIRType(irType)
		l.mod.Globals = append(l.mod.Globals, ir.Global{
			Name:    fullName,
			Type:    irType,
			Value:   valStr,
			Align:   align,
			Private: false,
		})
	}
	return nil
}

type lowerer struct {
	prog     *Program
	mod      *ir.Module
	fn       *ir.Function
	block    *ir.BasicBlock
	reg      int
	blockID  int
	locals   []map[string]ir.Value // имя -> адрес (alloca)
	loops    []loopContext
	strings  int
	retTypes []ir.Type
	prefix   string
}

var globalStringCounter int

type loopContext struct {
	breakTarget    string
	continueTarget string
}

func isRuntimeExtern(name string) bool {
	return name == "gominic_print" || name == "gominic_printInt" || name == "gominic_println" || name == "gominic_makeSlice"
}

func (l *lowerer) lowerFunc(fnDecl *ast.FuncDecl) error {
	if isRuntimeExtern(fnDecl.Name.Name) {
		return nil
	}
	// Skip declarations without bodies (externs).
	if fnDecl.Body == nil {
		return nil
	}

	obj, ok := l.prog.TypesInfo.Defs[fnDecl.Name]
	if !ok || obj == nil {
		return fmt.Errorf("no type info for function %s", fnDecl.Name.Name)
	}
	sig, ok := obj.Type().(*types.Signature)
	if !ok {
		return fmt.Errorf("unexpected signature type for %s", fnDecl.Name.Name)
	}

	params := l.lowerParams(sig)

	results := make([]ir.Type, sig.Results().Len())
	for i := 0; i < sig.Results().Len(); i++ {
		rType, err := goTypeToIR(sig.Results().At(i).Type())
		if err != nil {
			return fmt.Errorf("result %d: %v", i, err)
		}
		results[i] = rType
	}
	l.retTypes = results

	fnResults := results
	if len(results) > 1 {
		fnResults = []ir.Type{ir.StructType{Fields: results}}
	}

	funcName := fnDecl.Name.Name
	if sig.Recv() != nil {
		return fmt.Errorf("unsupported: methods with receiver")
	}
	if l.prefix != "" && !isRuntimeExtern(funcName) {
		funcName = l.prefix + funcName
	}

	l.fn = &ir.Function{
		Name:    funcName,
		Params:  params,
		Results: fnResults,
	}

	l.block = ir.NewBasicBlock("entry")
	l.addBlock(l.block)
	l.reg = 0
	l.blockID = 0
	l.resetScopes()
	l.pushScope()

	// Параметры кладем в стековые ячейки для возможности присваивания.
	for _, p := range l.fn.Params {
		if p.ByVal() {
			slot := ir.NewRegister(fmt.Sprintf("%s.addr", p.Name()), ir.PtrTo(p.ByValType()))
			l.block.Append(ir.Alloca{Dest: slot, AllocType: p.ByValType(), Align: alignOfIRType(p.ByValType())})
			size := ir.NewConstant(strconv.FormatInt(mustSizeofIRType(p.ByValType()), 10), ir.I64)
			l.block.Append(ir.Memcpy{Dest: slot, Src: p, Size: size})
			l.bindLocal(p.Name(), slot)
			continue
		}
		valTy := ir.ValueType(p)
		slot := ir.NewRegister(fmt.Sprintf("%s.addr", p.Name()), ir.PtrTo(valTy))
		l.block.Append(ir.Alloca{Dest: slot, AllocType: valTy, Align: alignOfIRType(valTy)})
		l.block.Append(ir.Store{Src: p, Dest: slot})
		l.bindLocal(p.Name(), slot)
	}

	if err := l.lowerStmtList(fnDecl.Body.List); err != nil {
		return err
	}

	// Добавляем ret по умолчанию, если нужно.
	if l.block != nil && l.block.Terminator == nil {
		if len(results) == 0 {
			l.block.Terminator = ir.Return{}
		} else {
			return fmt.Errorf("missing return in function %s", fnDecl.Name.Name)
		}
	}

	l.mod.AddFunction(l.fn)
	return nil
}

func (l *lowerer) lowerParams(sig *types.Signature) []ir.Param {
	return l.lowerParamTypes(sig)
}

func (l *lowerer) lowerParamTypes(sig *types.Signature) []ir.Param {
	offset := 0
	if sig.Recv() != nil {
		offset = 1
	}
	params := make([]ir.Param, sig.Params().Len()+offset)
	if sig.Recv() != nil {
		recvType, err := goTypeToIR(sig.Recv().Type())
		if err == nil {
			params[0] = ir.NewParam("recv", recvType)
		}
	}
	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		pType, err := goTypeToIR(p.Type())
		if err != nil {
			continue
		}
		name := p.Name()
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		if isAggregateType(pType) {
			params[i+offset] = ir.NewByValParam(name, ir.PtrTo(pType), pType)
		} else {
			params[i+offset] = ir.NewParam(name, pType)
		}
	}
	return params
}

func (l *lowerer) lowerStmtList(stmts []ast.Stmt) error {
	for _, s := range stmts {
		if l.block == nil || l.block.Terminator != nil {
			break
		}
		if err := l.lowerStmt(s); err != nil {
			return err
		}
	}
	return nil
}

func (l *lowerer) lowerStmt(stmt ast.Stmt) error {
	if s, ok := stmt.(*ast.DeclStmt); ok {
		return l.lowerDeclStmt(s)
	}
	if s, ok := stmt.(*ast.AssignStmt); ok {
		return l.lowerAssign(s)
	}
	if s, ok := stmt.(*ast.ReturnStmt); ok {
		return l.lowerReturn(s)
	}
	if s, ok := stmt.(*ast.IfStmt); ok {
		return l.lowerIf(s)
	}
	if s, ok := stmt.(*ast.ForStmt); ok {
		return l.lowerFor(s)
	}
	if _, ok := stmt.(*ast.SwitchStmt); ok {
		return fmt.Errorf("unsupported: switch statement")
	}
	if _, ok := stmt.(*ast.TypeSwitchStmt); ok {
		return fmt.Errorf("unsupported: type switch statement")
	}
	if s, ok := stmt.(*ast.BranchStmt); ok {
		return l.lowerBranch(s)
	}
	if _, ok := stmt.(*ast.RangeStmt); ok {
		return fmt.Errorf("unsupported: range statement")
	}
	if s, ok := stmt.(*ast.IncDecStmt); ok {
		return l.lowerIncDec(s)
	}
	if s, ok := stmt.(*ast.ExprStmt); ok {
		if _, ok := s.X.(*ast.CallExpr); !ok {
			return fmt.Errorf("unsupported expression statement %T", s.X)
		}
		_, err := l.lowerCall(s.X.(*ast.CallExpr), false)
		return err
	}
	return fmt.Errorf("unsupported statement %T", stmt)
}

func (l *lowerer) lowerDeclStmt(ds *ast.DeclStmt) error {
	gen, ok := ds.Decl.(*ast.GenDecl)
	if !ok {
		return fmt.Errorf("unsupported declaration %T", ds.Decl)
	}
	if gen.Tok != token.VAR {
		// ignore const/type declarations inside functions
		return nil
	}
	for _, spec := range gen.Specs {
		valSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			return fmt.Errorf("unsupported spec %T", spec)
		}
		for i, name := range valSpec.Names {
			obj := l.prog.TypesInfo.Defs[name]
			if obj == nil {
				return fmt.Errorf("no def info for var %s", name.Name)
			}
			slot, err := l.allocLocal(name.Name, obj.Type())
			if err != nil {
				return err
			}
			var init ir.Value
			if i < len(valSpec.Values) {
				init, err = l.lowerExpr(valSpec.Values[i])
				if err != nil {
					return err
				}
			} else {
				init, err = zeroValue(obj.Type())
				if err != nil {
					return err
				}
			}
			l.block.Append(ir.Store{Src: init, Dest: slot})
		}
	}
	return nil
}

func (l *lowerer) lowerAssign(as *ast.AssignStmt) error {
	if len(as.Lhs) != len(as.Rhs) {
		// special case: multiple LHS from single RHS returning struct (multiple values)
		if !(len(as.Lhs) > 1 && len(as.Rhs) == 1) {
			return fmt.Errorf("lhs/rhs arity mismatch")
		}
	}

	// handle map index with comma-ok assignment: x, ok := m[k]
	if len(as.Lhs) == 2 && len(as.Rhs) == 1 {
		if idx, ok := as.Rhs[0].(*ast.IndexExpr); ok {
			if tv, ok := l.prog.TypesInfo.Types[idx.X]; ok {
				if _, isMap := tv.Type.Underlying().(*types.Map); isMap {
					return l.lowerMapGetAssign(as.Lhs, idx, as.Tok == token.DEFINE)
				}
			}
		}
		if _, ok := as.Rhs[0].(*ast.TypeAssertExpr); ok {
			return fmt.Errorf("unsupported: type assertion in assignment")
		}
	}

	values := make([]ir.Value, len(as.Rhs))
	for i, expr := range as.Rhs {
		v, err := l.lowerExpr(expr)
		if err != nil {
			return err
		}
		values[i] = v
	}

	// handle tuple assignment from struct-valued RHS
	if len(as.Lhs) > 1 && len(values) == 1 {
		return l.destructureAssign(as, values[0])
	}

	for i, lhs := range as.Lhs {
		if id, ok := lhs.(*ast.Ident); ok && id.Name == "_" {
			continue
		}
		// map assignment: m[k] = v
		if idx, ok := lhs.(*ast.IndexExpr); ok {
			tv, ok := l.prog.TypesInfo.Types[idx.X]
			if ok {
				if _, isMap := tv.Type.Underlying().(*types.Map); isMap {
					if len(as.Rhs) != 1 {
						return fmt.Errorf("map assignment expects single rhs")
					}
					return l.lowerMapSet(idx, values[i])
				}
			}
		}
		if as.Tok == token.ASSIGN {
			slot, err := l.lowerAddr(lhs)
			if err != nil {
				return err
			}
			if err := l.storeValue(slot, values[i]); err != nil {
				return err
			}
		} else if as.Tok == token.DEFINE {
			if id, ok := lhs.(*ast.Ident); ok {
				if defObj := l.prog.TypesInfo.Defs[id]; defObj != nil {
					slot, err := l.allocLocal(id.Name, defObj.Type())
					if err != nil {
						return err
					}
					if err := l.storeValue(slot, values[i]); err != nil {
						return err
					}
					continue
				}
			}
				if slot, err := l.lowerAddr(lhs); err == nil {
					if err != nil {
						return err
					}
					if err := l.storeValue(slot, values[i]); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("short assign to unsupported lhs %T", lhs)
			} else {
				// compound assignments like +=, -=, etc.
				var baseTok token.Token
				switch as.Tok {
				case token.ADD_ASSIGN:
					baseTok = token.ADD
				case token.SUB_ASSIGN:
					baseTok = token.SUB
				case token.MUL_ASSIGN:
					baseTok = token.MUL
				case token.QUO_ASSIGN:
					baseTok = token.QUO
				case token.REM_ASSIGN:
					baseTok = token.REM
				default:
					return fmt.Errorf("unsupported assign token %s", as.Tok.String())
				}

				slot, err := l.lowerAddr(lhs)
				if err != nil {
					return err
				}
				cur, err := l.load(slot)
				if err != nil {
					return err
				}
				tv := l.prog.TypesInfo.TypeOf(lhs)
				if tv == nil {
					return fmt.Errorf("missing type for lhs in compound assign")
				}
				if baseTok == token.ADD {
					if basic, ok := tv.Underlying().(*types.Basic); ok && (basic.Kind() == types.String || basic.Kind() == types.UntypedString) {
						dest := l.newTemp(ir.ValueType(cur))
						l.block.Append(ir.Call{
							Dest: dest,
							Name: "gominic_str_concat",
							Args: []ir.Value{cur, values[i]},
							Ret:  ir.ValueType(cur),
						})
						if err := l.storeValue(slot, dest); err != nil {
							return err
						}
						continue
					}
				}
				op, err := selectBinOp(baseTok, tv)
				if err != nil {
					return err
				}
				ty, err := goTypeToIR(tv)
				if err != nil {
					return err
				}
				dest := l.newTemp(ty)
				l.block.Append(ir.BinOp{
					Dest: dest,
					Op:   op,
					X:    cur,
					Y:    values[i],
				})
				if err := l.storeValue(slot, dest); err != nil {
					return err
				}
			}
		}
		return nil
	}

func (l *lowerer) lowerReturn(rs *ast.ReturnStmt) error {
	if len(rs.Results) == 0 {
		if len(l.retTypes) != 0 {
			return fmt.Errorf("function %s must return a value", l.fn.Name)
		}
		l.block.Terminator = ir.Return{}
		return nil
	}

	if len(rs.Results) != len(l.retTypes) {
		// ignore mismatch; rely on caller expectations
	}

	if len(rs.Results) == 1 {
		val, err := l.lowerExpr(rs.Results[0])
		if err != nil {
			return err
		}
		l.block.Terminator = ir.Return{Results: []ir.Value{val}}
		return nil
	}

	// pack multiple return values into a struct.
	structTy := ir.StructType{Fields: l.retTypes}
	aggPtr := l.allocaTemp(structTy)
	zero := ir.NewConstant("0", ir.I32)
	for i := 0; i < len(rs.Results); i++ {
		val, err := l.lowerExpr(rs.Results[i])
		if err != nil {
			return err
		}
		idx := ir.NewConstant(strconv.Itoa(i), ir.I32)
	fieldPtr := l.newTemp(ir.PtrTo(ir.ValueType(val)))
		l.block.Append(ir.GetElementPtr{
			Dest:      fieldPtr,
			Src:       aggPtr,
			Pointee:   structTy,
			Indices:   []ir.Value{zero, idx},
		ResultTyp: ir.PtrTo(ir.ValueType(val)),
	})
	l.block.Append(ir.Store{Src: val, Dest: fieldPtr, Align: alignOfIRType(ir.ValueType(val))})
	}
	aggVal, err := l.load(aggPtr)
	if err != nil {
		return err
	}
	l.block.Terminator = ir.Return{Results: []ir.Value{aggVal}}
	return nil
}

func (l *lowerer) lowerIf(is *ast.IfStmt) error {
	l.pushScope()
	if is.Init != nil {
		if err := l.lowerStmt(is.Init); err != nil {
			l.popScope()
			return err
		}
	}

	cond, err := l.lowerExpr(is.Cond)
	if err != nil {
		l.popScope()
		return err
	}

	thenBlock := l.newBlock("then")
	var elseBlock *ir.BasicBlock
	joinBlock := l.newBlock("endif")

	if is.Else != nil {
		elseBlock = l.newBlock("else")
		l.block.Terminator = ir.CondBr{
			Cond:  cond,
			True:  thenBlock.Name,
			False: elseBlock.Name,
		}
	} else {
		// Пустой else ведет прямо в join.
		l.block.Terminator = ir.CondBr{
			Cond:  cond,
			True:  thenBlock.Name,
			False: joinBlock.Name,
		}
	}

	// then
	l.block = thenBlock
	l.pushScope()
	if err := l.lowerBlockStmt(is.Body); err != nil {
		return err
	}
	l.popScope()
	if l.block != nil && l.block.Terminator == nil {
		l.block.Terminator = ir.Br{Target: joinBlock.Name}
	}

	// else
	if is.Else != nil {
		l.block = elseBlock
		l.pushScope()
		if es, ok := is.Else.(*ast.BlockStmt); ok {
			if err := l.lowerBlockStmt(es); err != nil {
				return err
			}
		} else if esIf, ok := is.Else.(*ast.IfStmt); ok {
			if err := l.lowerStmt(esIf); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported else node %T", is.Else)
		}
		l.popScope()
		if l.block != nil && l.block.Terminator == nil {
			l.block.Terminator = ir.Br{Target: joinBlock.Name}
		}
	}

	// Устанавливаем текущий блок на join, если он достижим.
	l.block = joinBlock
	l.popScope()
	return nil
}

func (l *lowerer) lowerFor(fs *ast.ForStmt) error {
	if fs.Init != nil {
		if err := l.lowerStmt(fs.Init); err != nil {
			return err
		}
	}
	if l.block == nil || l.block.Terminator != nil {
		return nil
	}

	condBlock := l.newBlock("for.cond")
	bodyBlock := l.newBlock("for.body")
	var postBlock *ir.BasicBlock
	if fs.Post != nil {
		postBlock = l.newBlock("for.post")
	}
	endBlock := l.newBlock("for.end")

	// jump to condition
	l.block.Terminator = ir.Br{Target: condBlock.Name}

	// loop context
	contTarget := condBlock.Name
	if postBlock != nil {
		contTarget = postBlock.Name
	}
	l.pushLoop(loopContext{breakTarget: endBlock.Name, continueTarget: contTarget})

	// condition block
	l.block = condBlock
	if fs.Cond != nil {
		cond, err := l.lowerExpr(fs.Cond)
		if err != nil {
			return err
		}
		l.block.Terminator = ir.CondBr{
			Cond:  cond,
			True:  bodyBlock.Name,
			False: endBlock.Name,
		}
	} else {
		l.block.Terminator = ir.Br{Target: bodyBlock.Name}
	}

	// body
	l.block = bodyBlock
	l.pushScope()
	if err := l.lowerBlockStmt(fs.Body); err != nil {
		return err
	}
	l.popScope()
	if l.block != nil && l.block.Terminator == nil {
		if postBlock != nil {
			l.block.Terminator = ir.Br{Target: postBlock.Name}
		} else {
			l.block.Terminator = ir.Br{Target: condBlock.Name}
		}
	}

	// post
	if postBlock != nil {
		l.block = postBlock
		if err := l.lowerStmt(fs.Post); err != nil {
			return err
		}
		if l.block != nil && l.block.Terminator == nil {
			l.block.Terminator = ir.Br{Target: condBlock.Name}
		}
	}

	l.popLoop()

	// continue after loop
	l.block = endBlock
	return nil
}

func (l *lowerer) lowerBranch(bs *ast.BranchStmt) error {
	if len(l.loops) == 0 {
		return fmt.Errorf("%s outside of loop", bs.Tok.String())
	}
	ctx := l.loops[len(l.loops)-1]
	if bs.Tok == token.BREAK {
		l.block.Terminator = ir.Br{Target: ctx.breakTarget}
	} else if bs.Tok == token.CONTINUE {
		l.block.Terminator = ir.Br{Target: ctx.continueTarget}
	} else {
		return fmt.Errorf("unsupported branch token %s", bs.Tok.String())
	}
	return nil
}

// lowerSwitch handles expression switches (no fallthrough). Type switches are not supported.
func (l *lowerer) lowerSwitch(ss *ast.SwitchStmt) error {
	if ss.Init != nil {
		if err := l.lowerStmt(ss.Init); err != nil {
			return err
		}
		if l.block == nil || l.block.Terminator != nil {
			return nil
		}
	}

	var tagVal ir.Value
	var err error
	if ss.Tag != nil {
		tagVal, err = l.lowerExpr(ss.Tag)
		if err != nil {
			return err
		}
	}

	endBlock := l.newBlock("switch.end")
	nextBlock := endBlock

	// Iterate clauses in reverse to build if-else chain.
	for i := len(ss.Body.List) - 1; i >= 0; i-- {
		cc, ok := ss.Body.List[i].(*ast.CaseClause)
		if !ok {
			return fmt.Errorf("unexpected switch clause %T", ss.Body.List[i])
		}
		caseBlock := l.newBlock(fmt.Sprintf("case.%d", i))

		// Build condition for this case.
		var cond ir.Value
		if len(cc.List) == 0 {
			// default
			cond = ir.NewConstant("1", ir.I1)
		} else {
			for j := 0; j < len(cc.List); j++ {
				exprVal, err := l.lowerExpr(cc.List[j])
				if err != nil {
					return err
				}
				var eq ir.Value
				if tagVal.Kind() == ir.ValueInvalid {
					// switch without tag: case expressions are booleans.
					eq = exprVal
				} else {
					var cmpType types.Type
					if tv, okType := l.prog.TypesInfo.Types[cc.List[j]]; okType && tv.Type != nil {
						cmpType = tv.Type
					} else if tvTag := l.prog.TypesInfo.TypeOf(ss.Tag); tvTag != nil {
						cmpType = tvTag
					}
					if cmpType != nil {
						eq, err = l.lowerCompare(token.EQL, cmpType, tagVal, exprVal)
						if err != nil {
							return err
						}
					} else {
						eqReg := l.newTemp(ir.I1)
						l.block.Append(ir.ICmp{
							Dest: eqReg,
							Pred: ir.ICmpEq,
							X:    tagVal,
							Y:    exprVal,
						})
						eq = eqReg
					}
				}
				if cond.Kind() == ir.ValueInvalid {
					cond = eq
				} else {
					tmp := l.newTemp(ir.I1)
					l.block.Append(ir.BinOp{Dest: tmp, Op: ir.Or, X: cond, Y: eq})
					cond = tmp
				}
			}
		}

		// current block branches based on cond.
		l.block.Terminator = ir.CondBr{Cond: cond, True: caseBlock.Name, False: nextBlock.Name}

		// case body
		l.block = caseBlock
		l.pushScope()
		if err := l.lowerStmtList(cc.Body); err != nil {
			return err
		}
		l.popScope()
		if l.block != nil && l.block.Terminator == nil {
			l.block.Terminator = ir.Br{Target: endBlock.Name}
		}

		// Prepare for previous clause
		l.block = l.newBlock(fmt.Sprintf("case.entry.%d", i))
		nextBlock = l.block
	}

	// final jump into chain
	if l.block != nil && l.block.Terminator == nil {
		l.block.Terminator = ir.Br{Target: nextBlock.Name}
	}
	l.block = endBlock
	return nil
}

// lowerRange lowers range over slices and arrays: for k,v := range X { ... }.
func (l *lowerer) lowerRange(rs *ast.RangeStmt) error {
	tv, ok := l.prog.TypesInfo.Types[rs.X]
	if !ok || tv.Type == nil {
		return fmt.Errorf("no type info for range expr")
	}
	baseType := tv.Type.Underlying()

	var lenVal ir.Value
	var dataPtr ir.Value
	var elemIR ir.Type
	var arrayPointee ir.Type

	if t, ok := baseType.(*types.Slice); ok {
		sliceVal, err := l.lowerExpr(rs.X)
		if err != nil {
			return err
		}
		sliceTy, err := goTypeToIR(t)
		if err != nil {
			return err
		}
		slicePtr := l.allocaTemp(sliceTy)
		l.block.Append(ir.Store{Src: sliceVal, Dest: slicePtr, Align: alignOfIRType(sliceTy)})

		zero := ir.NewConstant("0", ir.I32)
		// len field at index 1
		lenField := ir.NewConstant("1", ir.I32)
		lenPtr := l.newTemp(ir.PtrTo(ir.I64))
		l.block.Append(ir.GetElementPtr{
			Dest:      lenPtr,
			Src:       slicePtr,
			Pointee:   sliceTy,
			Indices:   []ir.Value{zero, lenField},
			ResultTyp: ir.PtrTo(ir.I64),
		})
		val, err := l.load(lenPtr)
		if err != nil {
			return err
		}
		lenVal = val

		// data field at index 0
		dataField := ir.NewConstant("0", ir.I32)
		dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrTo(sliceTy.(ir.SliceType).Elem)))
		l.block.Append(ir.GetElementPtr{
			Dest:      dataPtrPtr,
			Src:       slicePtr,
			Pointee:   sliceTy,
			Indices:   []ir.Value{zero, dataField},
			ResultTyp: ir.PtrTo(ir.PtrTo(sliceTy.(ir.SliceType).Elem)),
		})
		dp, err := l.load(dataPtrPtr)
		if err != nil {
			return err
		}
		dataPtr = dp
		elemIR, err = goTypeToIR(t.Elem())
		if err != nil {
			return err
		}
	} else if t, ok := baseType.(*types.Array); ok {
		elem, err := goTypeToIR(t.Elem())
		if err != nil {
			return err
		}
		elemIR = elem
		lenVal = ir.NewConstant(strconv.FormatInt(t.Len(), 10), ir.I64)
		basePtr, err := l.lowerAddr(rs.X)
		if err != nil {
			return err
		}
		dataPtr = basePtr
		arrayPointee = ir.ArrayType{Len: int(t.Len()), Elem: elemIR}
	} else {
		return fmt.Errorf("range over unsupported type %T", baseType)
	}

	// index variable
	idxSlot := l.allocaTemp(ir.I64)
	l.block.Append(ir.Store{Src: ir.NewConstant("0", ir.I64), Dest: idxSlot})

	condBlock := l.newBlock("range.cond")
	bodyBlock := l.newBlock("range.body")
	postBlock := l.newBlock("range.post")
	endBlock := l.newBlock("range.end")

	l.block.Terminator = ir.Br{Target: condBlock.Name}
	l.pushLoop(loopContext{breakTarget: endBlock.Name, continueTarget: postBlock.Name})

	// cond
	l.block = condBlock
	curIdx, err := l.load(idxSlot)
	if err != nil {
		return err
	}
	cond := l.newTemp(ir.I1)
	l.block.Append(ir.ICmp{
		Dest: cond,
		Pred: ir.ICmpSlt,
		X:    curIdx,
		Y:    lenVal,
	})
	l.block.Terminator = ir.CondBr{Cond: cond, True: bodyBlock.Name, False: endBlock.Name}

	// body
	l.block = bodyBlock
	l.pushScope()
	if rs.Key != nil {
		if err := l.assignRangeKey(rs, curIdx); err != nil {
			return err
		}
	}
	if rs.Value != nil {
		var elemPtr ir.Register
		if _, ok := baseType.(*types.Slice); ok {
			elemPtr = l.newTemp(ir.PtrTo(elemIR))
			l.block.Append(ir.GetElementPtr{
				Dest:      elemPtr,
				Src:       dataPtr,
				Pointee:   elemIR,
				Indices:   []ir.Value{curIdx},
				ResultTyp: ir.PtrTo(elemIR),
			})
		} else {
			zero := ir.NewConstant("0", ir.I32)
			idx32 := l.newTemp(ir.I32)
			l.block.Append(ir.Conv{Dest: idx32, Op: ir.Trunc, Src: curIdx, To: ir.I32})
			elemPtr = l.newTemp(ir.PtrTo(elemIR))
			l.block.Append(ir.GetElementPtr{
				Dest:      elemPtr,
				Src:       dataPtr,
				Pointee:   arrayPointee,
				Indices:   []ir.Value{zero, idx32},
				ResultTyp: ir.PtrTo(elemIR),
			})
		}
		val, err := l.load(elemPtr)
		if err != nil {
			return err
		}
		if err := l.assignRangeValue(rs, val); err != nil {
			return err
		}
	}
	if err := l.lowerBlockStmt(rs.Body); err != nil {
		return err
	}
	l.popScope()
	if l.block != nil && l.block.Terminator == nil {
		l.block.Terminator = ir.Br{Target: postBlock.Name}
	}

	// post
	l.block = postBlock
	next := l.newTemp(ir.I64)
	l.block.Append(ir.BinOp{Dest: next, Op: ir.Add, X: curIdx, Y: ir.NewConstant("1", ir.I64)})
	l.block.Append(ir.Store{Src: next, Dest: idxSlot})
	l.block.Terminator = ir.Br{Target: condBlock.Name}

	l.popLoop()
	l.block = endBlock
	return nil
}

func (l *lowerer) assignRangeKey(rs *ast.RangeStmt, idx ir.Value) error {
	if ident, ok := rs.Key.(*ast.Ident); ok {
		if rs.Tok == token.DEFINE {
			if def := l.prog.TypesInfo.Defs[ident]; def != nil {
				if _, err := l.allocLocal(ident.Name, def.Type()); err != nil {
					return err
				}
			}
		}
		ptr, err := l.lowerAddr(rs.Key)
		if err != nil {
			return err
		}
		return l.storeValue(ptr, idx)
	}
	return fmt.Errorf("unsupported range key %T", rs.Key)
}

func (l *lowerer) assignRangeValue(rs *ast.RangeStmt, val ir.Value) error {
	if ident, ok := rs.Value.(*ast.Ident); ok {
		if rs.Tok == token.DEFINE {
			if def := l.prog.TypesInfo.Defs[ident]; def != nil {
				if _, err := l.allocLocal(ident.Name, def.Type()); err != nil {
					return err
				}
			}
		}
		ptr, err := l.lowerAddr(rs.Value)
		if err != nil {
			return err
		}
		return l.storeValue(ptr, val)
	}
	return fmt.Errorf("unsupported range value %T", rs.Value)
}

func (l *lowerer) lowerIncDec(id *ast.IncDecStmt) error {
	slot, err := l.lowerAddr(id.X)
	if err != nil {
		return err
	}
	tv, ok := l.prog.TypesInfo.Types[id.X]
	if !ok || tv.Type == nil {
		return fmt.Errorf("no type info for inc/dec target")
	}
	if basic, ok := tv.Type.Underlying().(*types.Basic); !ok || basic.Info()&types.IsInteger == 0 {
		return fmt.Errorf("inc/dec allowed only on integer variables")
	}
	val, err := l.load(slot)
	if err != nil {
		return err
	}
	one := ir.NewConstant("1", ir.ValueType(val))
	dest := l.newTemp(ir.ValueType(val))
	op := ir.Add
	if id.Tok == token.DEC {
		op = ir.Sub
	}
	l.block.Append(ir.BinOp{
		Dest: dest,
		Op:   op,
		X:    val,
		Y:    one,
	})
	l.block.Append(ir.Store{Src: dest, Dest: slot, Align: alignOfIRType(ir.ValueType(dest))})
	return nil
}

func (l *lowerer) destructureAssign(as *ast.AssignStmt, rhs ir.Value) error {
	structTy, ok := ir.ValueType(rhs).(ir.StructType)
	if !ok {
		return fmt.Errorf("tuple assignment requires struct value")
	}
	if len(as.Lhs) != len(structTy.Fields) {
		return fmt.Errorf("tuple arity mismatch")
	}
	tmp := l.allocaTemp(structTy)
	l.block.Append(ir.Store{Src: rhs, Dest: tmp, Align: alignOfIRType(structTy)})
	for i, lhs := range as.Lhs {
		slot, err := l.lowerAddr(lhs)
		if err != nil {
			if as.Tok == token.DEFINE {
				if id, ok := lhs.(*ast.Ident); ok {
					if defObj := l.prog.TypesInfo.Defs[id]; defObj != nil {
						slot, err = l.allocLocal(id.Name, defObj.Type())
						if err != nil {
							return err
						}
					} else {
						return err
					}
				} else {
					return err
				}
			} else {
				return err
			}
		}
		zero := ir.NewConstant("0", ir.I32)
		idx := ir.NewConstant(strconv.Itoa(i), ir.I32)
		fieldPtr := l.newTemp(ir.PtrTo(structTy.Fields[i]))
		l.block.Append(ir.GetElementPtr{
			Dest:      fieldPtr,
			Src:       tmp,
			Pointee:   structTy,
			Indices:   []ir.Value{zero, idx},
			ResultTyp: ir.PtrTo(structTy.Fields[i]),
		})
		val, err := l.load(fieldPtr)
		if err != nil {
			return err
		}
		if err := l.storeValue(slot, val); err != nil {
			return err
		}
	}
	return nil
}

func (l *lowerer) lowerBlockStmt(b *ast.BlockStmt) error {
	return l.lowerStmtList(b.List)
}

func (l *lowerer) lowerExpr(expr ast.Expr) (ir.Value, error) {
	if e, ok := expr.(*ast.BasicLit); ok {
		return l.constValue(e)
	}
	if e, ok := expr.(*ast.Ident); ok {
		if slot, ok := l.lookupLocal(e.Name); ok {
			return l.load(slot)
		}
		// Try package-level object.
		if obj := l.prog.TypesInfo.Uses[e]; obj != nil {
			if c, ok := obj.(*types.Const); ok {
				if v, err := l.constValue(e); err == nil {
					return v, nil
				}
				_ = c
			}
			if v, ok := obj.(*types.Var); ok {
				irType, err := goTypeToIR(v.Type())
				if err != nil {
					return nil, err
				}
				name := e.Name
				if l.prefix != "" {
					name = l.prefix + name
				}
				ptrConst := ir.NewConstant("@"+name, ir.PtrTo(irType))
				return l.load(ptrConst)
			}
		}
		return l.constValue(e)
	}
	if e, ok := expr.(*ast.ParenExpr); ok {
		return l.lowerExpr(e.X)
	}
	if _, ok := expr.(*ast.ArrayType); ok {
		if tv := l.prog.TypesInfo.TypeOf(expr); tv != nil {
			return zeroValue(tv)
		}
		return nil, fmt.Errorf("missing type for array expression")
	}
	if _, ok := expr.(*ast.MapType); ok {
		if tv := l.prog.TypesInfo.TypeOf(expr); tv != nil {
			return zeroValue(tv)
		}
		return nil, fmt.Errorf("missing type for map expression")
	}
	if e, ok := expr.(*ast.BinaryExpr); ok {
		return l.lowerBinary(e)
	}
	if e, ok := expr.(*ast.UnaryExpr); ok {
		return l.lowerUnary(e)
	}
	if e, ok := expr.(*ast.SelectorExpr); ok {
		// package selector or field
		if selInfo, okSel := l.prog.TypesInfo.Selections[e]; okSel && selInfo != nil {
			ptr, err := l.lowerAddr(e)
			if err != nil {
				return nil, err
			}
			return l.load(ptr)
		}
		// package-level object
		if obj, okUse := l.prog.TypesInfo.Uses[e.Sel]; okUse {
			if v, okVar := obj.(*types.Var); okVar {
				irType, err := goTypeToIR(v.Type())
				if err != nil {
					return nil, err
				}
				pkgName := ""
				if id, okID := e.X.(*ast.Ident); okID {
					pkgName = id.Name + "."
				}
				name := pkgName + e.Sel.Name
				ptrConst := ir.NewConstant("@"+name, ir.PtrTo(irType))
				return l.load(ptrConst)
			}
			if _, okConst := obj.(*types.Const); okConst {
				return l.constValue(e)
			}
		}
		return nil, fmt.Errorf("unsupported selector %s", e.Sel.Name)
	}
	if e, ok := expr.(*ast.IndexExpr); ok {
		return l.lowerIndexExpr(e)
	}
	if e, ok := expr.(*ast.SliceExpr); ok {
		return l.lowerSliceExpr(e)
	}
	if _, ok := expr.(*ast.TypeAssertExpr); ok {
		return nil, fmt.Errorf("unsupported: type assertion")
	}
	if e, ok := expr.(*ast.CallExpr); ok {
		return l.lowerCall(e, true)
	}
	if e, ok := expr.(*ast.CompositeLit); ok {
		return l.lowerCompositeLit(e)
	}
	return nil, fmt.Errorf("unsupported expression %T", expr)
}

func (l *lowerer) lowerBinary(e *ast.BinaryExpr) (ir.Value, error) {
	x, err := l.lowerExpr(e.X)
	if err != nil {
		return nil, err
	}
	y, err := l.lowerExpr(e.Y)
	if err != nil {
		return nil, err
	}

	tv, ok := l.prog.TypesInfo.Types[e]
	if !ok {
		return nil, fmt.Errorf("no type info for binary expr")
	}

	// Сравнения дают bool, остальное — арифметика.
	if e.Op == token.EQL || e.Op == token.NEQ || e.Op == token.LSS || e.Op == token.LEQ || e.Op == token.GTR || e.Op == token.GEQ {
		xType := l.prog.TypesInfo.Types[e.X]
		if !xType.IsValue() {
			return nil, fmt.Errorf("missing operand type for comparison")
		}
		return l.lowerCompare(e.Op, xType.Type, x, y)
	}
	if e.Op == token.LAND || e.Op == token.LOR {
		return l.lowerLogical(e.Op, e, x)
	}
	if e.Op == token.ADD || e.Op == token.SUB || e.Op == token.MUL || e.Op == token.QUO || e.Op == token.REM || e.Op == token.AND || e.Op == token.OR {
		if basic, ok := tv.Type.Underlying().(*types.Basic); ok && (basic.Kind() == types.String || basic.Kind() == types.UntypedString) && e.Op == token.ADD {
			dest := l.newTemp(ir.ValueType(x))
			l.block.Append(ir.Call{
				Dest: dest,
				Name: "gominic_str_concat",
				Args: []ir.Value{x, y},
				Ret:  ir.ValueType(x),
			})
			return dest, nil
		}
		ty, err := goTypeToIR(tv.Type)
		if err != nil {
			return nil, err
		}
		op, err := selectBinOp(e.Op, tv.Type)
		if err != nil {
			return nil, err
		}
		dest := l.newTemp(ty)
		l.block.Append(ir.BinOp{
			Dest: dest,
			Op:   op,
			X:    x,
			Y:    y,
		})
		return dest, nil
	}
	return nil, fmt.Errorf("unsupported binary op %s", e.Op.String())
}

func (l *lowerer) lowerLogical(tok token.Token, e *ast.BinaryExpr, lhs ir.Value) (ir.Value, error) {
	tv, ok := l.prog.TypesInfo.Types[e]
	if !ok {
		return nil, fmt.Errorf("no type info for logical expr")
	}
	if basic, ok := tv.Type.Underlying().(*types.Basic); !ok || basic.Info()&types.IsBoolean == 0 {
		return nil, fmt.Errorf("logical op requires bool type")
	}

	slot := l.allocaTemp(ir.I1)
	trueConst := ir.NewConstant("1", ir.I1)
	falseConst := ir.NewConstant("0", ir.I1)

	endBlock := l.newBlock("logic.end")

	if tok == token.LAND {
		rhsBlock := l.newBlock("logic.rhs")
		falseBlock := l.newBlock("logic.false")
		l.block.Terminator = ir.CondBr{
			Cond:  lhs,
			True:  rhsBlock.Name,
			False: falseBlock.Name,
		}

		l.block = rhsBlock
		rhs, err := l.lowerExpr(e.Y)
		if err != nil {
			return nil, err
		}
		l.block.Append(ir.Store{Src: rhs, Dest: slot})
		l.block.Terminator = ir.Br{Target: endBlock.Name}

		l.block = falseBlock
		l.block.Append(ir.Store{Src: falseConst, Dest: slot})
		l.block.Terminator = ir.Br{Target: endBlock.Name}

	} else if tok == token.LOR {
		trueBlock := l.newBlock("logic.true")
		rhsBlock := l.newBlock("logic.rhs")
		l.block.Terminator = ir.CondBr{
			Cond:  lhs,
			True:  trueBlock.Name,
			False: rhsBlock.Name,
		}

		l.block = trueBlock
		l.block.Append(ir.Store{Src: trueConst, Dest: slot})
		l.block.Terminator = ir.Br{Target: endBlock.Name}

		l.block = rhsBlock
		rhs, err := l.lowerExpr(e.Y)
		if err != nil {
			return nil, err
		}
		l.block.Append(ir.Store{Src: rhs, Dest: slot})
		l.block.Terminator = ir.Br{Target: endBlock.Name}
	} else {
		return nil, fmt.Errorf("unsupported logical op %s", tok.String())
	}

	l.block = endBlock
	result, err := l.load(slot)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// boundsCheckSlice emits a bounds check for slice index: 0 <= idx < len.
func (l *lowerer) boundsCheckSlice(slicePtr ir.Value, idx ir.Value) error {
	// len is at index 1 of slice struct {data*, len, cap}
	sliceTy, err := pointeeType(ir.ValueType(slicePtr))
	if err != nil {
		return err
	}
	sliceStruct, ok := sliceTy.(ir.SliceType)
	if !ok {
		return fmt.Errorf("bounds check on non-slice")
	}

	zero := ir.NewConstant("0", ir.I32)
	lenField := ir.NewConstant("1", ir.I32)
	lenPtr := l.newTemp(ir.PtrTo(ir.I64))
	l.block.Append(ir.GetElementPtr{
		Dest:      lenPtr,
		Src:       slicePtr,
		Pointee:   sliceStruct,
		Indices:   []ir.Value{zero, lenField},
		ResultTyp: ir.PtrTo(ir.I64),
	})
	lenVal, err := l.load(lenPtr)
	if err != nil {
		return err
	}

	idx64, err := l.ensureI64(idx)
	if err != nil {
		return err
	}

	// cond = idx < len
	condNonNeg := l.newTemp(ir.I1)
	l.block.Append(ir.ICmp{
		Dest: condNonNeg,
		Pred: ir.ICmpSge,
		X:    idx64,
		Y:    ir.NewConstant("0", ir.I64),
	})

	condLtLen := l.newTemp(ir.I1)
	l.block.Append(ir.ICmp{
		Dest: condLtLen,
		Pred: ir.ICmpSlt,
		X:    idx64,
		Y:    lenVal,
	})

	cond := l.newTemp(ir.I1)
	l.block.Append(ir.BinOp{
		Dest: cond,
		Op:   ir.And,
		X:    condNonNeg,
		Y:    condLtLen,
	})

	okBlock := l.newBlock("idx.ok")
	failBlock := l.newBlock("idx.fail")

	l.block.Terminator = ir.CondBr{
		Cond:  cond,
		True:  okBlock.Name,
		False: failBlock.Name,
	}

	// fail path: call abort
	l.block = failBlock
	l.block.Instrs = append(l.block.Instrs, ir.CallVoid{Name: "gominic_abort"})
	l.block.Terminator = ir.Br{Target: okBlock.Name}

	// continue in ok block
	l.block = okBlock
	return nil
}

// ensureIndexWithinBounds checks idx against bound (array length), assuming idx is int.
func (l *lowerer) ensureIndexWithinBounds(idx ir.Value, bound ir.Value) (ir.Value, error) {
	idx64, err := l.ensureI64(idx)
	if err != nil {
		return nil, err
	}
	bound64, err := l.ensureI64(bound)
	if err != nil {
		return nil, err
	}
	cond := l.newTemp(ir.I1)
	l.block.Append(ir.ICmp{
		Dest: cond,
		Pred: ir.ICmpSlt,
		X:    idx64,
		Y:    bound64,
	})

	okBlock := l.newBlock("idx.ok")
	failBlock := l.newBlock("idx.fail")
	l.block.Terminator = ir.CondBr{
		Cond:  cond,
		True:  okBlock.Name,
		False: failBlock.Name,
	}

	l.block = failBlock
	l.block.Instrs = append(l.block.Instrs, ir.CallVoid{Name: "gominic_abort"})
	l.block = okBlock
	return idx64, nil
}

func (l *lowerer) lowerCompare(tok token.Token, t types.Type, x, y ir.Value) (ir.Value, error) {
	if named, ok := t.(*types.Named); ok {
		return l.lowerCompare(tok, named.Underlying(), x, y)
	}
	u := t.Underlying()
	// Treat interfaces like opaque pointers; only == and != are allowed.
	if _, ok := u.(*types.Interface); ok {
		if tok != token.EQL && tok != token.NEQ {
			return nil, fmt.Errorf("unsupported interface compare op %s", tok.String())
		}
		pred := ir.ICmpEq
		if tok == token.NEQ {
			pred = ir.ICmpNe
		}
		dest := l.newTemp(ir.I1)
		l.block.Append(ir.ICmp{
			Dest: dest,
			Pred: pred,
			X:    x,
			Y:    y,
		})
		return dest, nil
	}
	if _, ok := u.(*types.Pointer); ok {
		if tok != token.EQL && tok != token.NEQ {
			return nil, fmt.Errorf("unsupported pointer compare op %s", tok.String())
		}
		pred := ir.ICmpEq
		if tok == token.NEQ {
			pred = ir.ICmpNe
		}
		dest := l.newTemp(ir.I1)
		// ICmp supports pointer types directly.
		l.block.Append(ir.ICmp{
			Dest: dest,
			Pred: pred,
			X:    x,
			Y:    y,
		})
		return dest, nil
	}
	if basic, ok := u.(*types.Basic); ok {
		if basic.Kind() == types.String || basic.Kind() == types.UntypedString {
			if tok != token.EQL && tok != token.NEQ {
				return nil, fmt.Errorf("unsupported string compare op %s", tok.String())
			}
			// If values are pointers/opaque handles, compare pointers directly.
			if _, ok := ir.ValueType(x).(ir.PointerType); ok || ir.ValueType(x).String() == ir.PtrI8.String() {
				pred := ir.ICmpEq
				if tok == token.NEQ {
					pred = ir.ICmpNe
				}
				yVal := y
				if ir.ValueType(x).String() != ir.ValueType(y).String() {
					cast := l.newTemp(ir.ValueType(x))
					l.block.Append(ir.Bitcast{Dest: cast, Src: y, Target: ir.ValueType(x)})
					yVal = cast
				}
				dest := l.newTemp(ir.I1)
				l.block.Append(ir.ICmp{Dest: dest, Pred: pred, X: x, Y: yVal})
				return dest, nil
			}
			// Ensure we have concrete string values (not pointers).
			xVal := x
			if p, ok := ir.ValueType(xVal).(ir.PointerType); ok {
				tmp := l.newTemp(p.Elem)
				l.block.Append(ir.Load{Dest: tmp, Src: xVal, Align: alignOfIRType(p.Elem)})
				xVal = tmp
			}
			yVal := y
			if p, ok := ir.ValueType(yVal).(ir.PointerType); ok {
				tmp := l.newTemp(p.Elem)
				l.block.Append(ir.Load{Dest: tmp, Src: yVal, Align: alignOfIRType(p.Elem)})
				yVal = tmp
			}
			strTy := ir.ValueType(xVal)
			baseX := l.allocaTemp(strTy)
			l.block.Append(ir.Store{Src: xVal, Dest: baseX, Align: alignOfIRType(strTy)})
			baseY := l.allocaTemp(strTy)
			l.block.Append(ir.Store{Src: yVal, Dest: baseY, Align: alignOfIRType(strTy)})

			zero32 := ir.NewConstant("0", ir.I32)
			// x data
			xDataPtr := l.newTemp(ir.PtrTo(ir.PtrI8))
			l.block.Append(ir.GetElementPtr{
				Dest:      xDataPtr,
				Src:       baseX,
				Pointee:   strTy,
				Indices:   []ir.Value{zero32, zero32},
				ResultTyp: ir.PtrTo(ir.PtrI8),
			})
			xData := l.newTemp(ir.PtrI8)
			l.block.Append(ir.Load{Dest: xData, Src: xDataPtr, Align: alignOfIRType(ir.PtrI8)})
			// x len
			xLenPtr := l.newTemp(ir.PtrTo(ir.I64))
			lenIdx := ir.NewConstant("1", ir.I32)
			l.block.Append(ir.GetElementPtr{
				Dest:      xLenPtr,
				Src:       baseX,
				Pointee:   strTy,
				Indices:   []ir.Value{zero32, lenIdx},
				ResultTyp: ir.PtrTo(ir.I64),
			})
			xLen := l.newTemp(ir.I64)
			l.block.Append(ir.Load{Dest: xLen, Src: xLenPtr, Align: alignOfIRType(ir.I64)})

			// y data
			yDataPtr := l.newTemp(ir.PtrTo(ir.PtrI8))
			l.block.Append(ir.GetElementPtr{
				Dest:      yDataPtr,
				Src:       baseY,
				Pointee:   strTy,
				Indices:   []ir.Value{zero32, zero32},
				ResultTyp: ir.PtrTo(ir.PtrI8),
			})
			yData := l.newTemp(ir.PtrI8)
			l.block.Append(ir.Load{Dest: yData, Src: yDataPtr, Align: alignOfIRType(ir.PtrI8)})
			// y len
			yLenPtr := l.newTemp(ir.PtrTo(ir.I64))
			l.block.Append(ir.GetElementPtr{
				Dest:      yLenPtr,
				Src:       baseY,
				Pointee:   strTy,
				Indices:   []ir.Value{zero32, lenIdx},
				ResultTyp: ir.PtrTo(ir.I64),
			})
			yLen := l.newTemp(ir.I64)
			l.block.Append(ir.Load{Dest: yLen, Src: yLenPtr, Align: alignOfIRType(ir.I64)})

			ret := l.newTemp(ir.I1)
			l.block.Append(ir.Call{
				Dest:    ret,
				Name:    "gominic_str_eq",
				Args:    []ir.Value{xData, xLen, yData, yLen},
				Ret:     ir.I1,
			})
			if tok == token.EQL {
				return ret, nil
			}
			// NEQ: invert
			neg := l.newTemp(ir.I1)
			l.block.Append(ir.ICmp{
				Dest: neg,
				Pred: ir.ICmpEq,
				X:    ret,
				Y:    ir.NewConstant("0", ir.I1),
			})
			return neg, nil
		}
		if basic.Info()&types.IsBoolean != 0 {
			if tok != token.EQL && tok != token.NEQ {
				return nil, fmt.Errorf("unsupported boolean compare op %s", tok.String())
			}
			pred := ir.ICmpEq
			if tok == token.NEQ {
				pred = ir.ICmpNe
			}
			dest := l.newTemp(ir.I1)
			l.block.Append(ir.ICmp{
				Dest: dest,
				Pred: pred,
				X:    x,
				Y:    y,
			})
			return dest, nil
		}
		if basic.Info()&types.IsFloat != 0 {
			pred, err := selectFCmpPred(tok)
			if err != nil {
				return nil, err
			}
			dest := l.newTemp(ir.I1)
			l.block.Append(ir.FCmp{
				Dest: dest,
				Pred: pred,
				X:    x,
				Y:    y,
			})
			return dest, nil
		}
		if basic.Info()&types.IsInteger != 0 {
			pred, err := selectICmpPred(tok, isUnsigned(t))
			if err != nil {
				return nil, err
			}
			dest := l.newTemp(ir.I1)
			l.block.Append(ir.ICmp{
				Dest: dest,
				Pred: pred,
				X:    x,
				Y:    y,
			})
			return dest, nil
		}
	}
	return nil, fmt.Errorf("unsupported compare type %T", t)
}

func (l *lowerer) constValue(expr ast.Expr) (ir.Value, error) {
	tv, ok := l.prog.TypesInfo.Types[expr]
	if !ok || !tv.IsValue() || tv.Value == nil {
		// Fallback: synthesize zero of the expression type if known.
		if tv.Type != nil {
			if z, err := zeroValue(tv.Type); err == nil {
				return z, nil
			}
		}
		// Ultimate fallback: i64 0.
		return ir.NewConstant("0", ir.I64), nil
	}

	ty, err := goTypeToIR(tv.Type)
	if err != nil {
		return nil, err
	}

	kind := tv.Value.Kind()
	if kind == constant.Bool {
		if constant.BoolVal(tv.Value) {
			return ir.NewConstant("1", ty), nil
		}
		return ir.NewConstant("0", ty), nil
	}
	if kind == constant.Int {
		if isUnsigned(tv.Type) {
			if v, ok := constant.Uint64Val(tv.Value); ok {
				return ir.NewConstant(strconv.FormatUint(v, 10), ty), nil
			}
		} else {
			if v, ok := constant.Int64Val(tv.Value); ok {
				return ir.NewConstant(strconv.FormatInt(v, 10), ty), nil
			}
		}
		return nil, fmt.Errorf("integer constant too large")
	}
	if kind == constant.Float {
		if v, ok := constant.Float64Val(tv.Value); ok {
			return ir.NewConstant(strconv.FormatFloat(v, 'g', -1, 64), ty), nil
		}
		return nil, fmt.Errorf("float constant out of range")
	}
	if kind == constant.String {
		val := constant.StringVal(tv.Value)
		return l.lowerStringLiteral(val)
	}
	return nil, fmt.Errorf("unsupported constant kind %v", tv.Value.Kind())
}

func constantString(v ir.Value) (string, error) {
	if v != nil && v.Kind() == ir.ValueConstant {
		return v.Raw, nil
	}
	return "", fmt.Errorf("not a constant")
}

func (l *lowerer) lowerMake(call *ast.CallExpr, wantValue bool) (ir.Value, error) {
	makeType := l.prog.TypesInfo.TypeOf(call)
	if makeType == nil {
		return nil, fmt.Errorf("missing type for make")
	}
	if sliceType, ok := makeType.Underlying().(*types.Slice); ok {
		if len(call.Args) < 2 || len(call.Args) > 3 {
			return nil, fmt.Errorf("make slice expects 2 or 3 arguments")
		}
		lenVal, err := l.lowerExpr(call.Args[1])
		if err != nil {
			return nil, err
		}
		var capVal ir.Value
		if len(call.Args) == 3 {
			capVal, err = l.lowerExpr(call.Args[2])
			if err != nil {
				return nil, err
			}
		} else {
			capVal = lenVal
		}

		elemTy, err := goTypeToIR(sliceType.Elem())
		if err != nil {
			return nil, err
		}
		elemSize, err := sizeofType(sliceType.Elem())
		if err != nil {
			return nil, err
		}

		allocCall := ir.Call{
			Name: "gominic_makeSlice",
			Args: []ir.Value{
				lenVal,
				capVal,
				ir.NewConstant(strconv.FormatInt(elemSize, 10), ir.I64),
			},
			Ret: ir.PtrI8,
		}
		rawPtr := l.newTemp(ir.PtrI8)
		allocCall.Dest = rawPtr
		l.block.Append(allocCall)

		dataPtr := l.newTemp(ir.PtrTo(elemTy))
		l.block.Append(ir.Bitcast{
			Dest:   dataPtr,
			Src:    rawPtr,
			Target: ir.PtrTo(elemTy),
		})

		return l.buildSliceValue(dataPtr, lenVal, capVal, elemTy)
	}
	if mapType, ok := makeType.Underlying().(*types.Map); ok {
		// make(map[K]V [, hint])
		if len(call.Args) < 1 || len(call.Args) > 2 {
			return nil, fmt.Errorf("make map expects 1 or 2 arguments")
		}
		keySize, err := sizeofType(mapType.Key())
		if err != nil {
			return nil, err
		}
		valSize, err := sizeofType(mapType.Elem())
		if err != nil {
			return nil, err
		}
		keyKind, err := mapKeyKind(mapType.Key())
		if err != nil {
			return nil, err
		}
		callInst := ir.Call{
			Name: "gominic_map_new",
			Args: []ir.Value{
				ir.NewConstant(strconv.FormatInt(keySize, 10), ir.I64),
				ir.NewConstant(strconv.FormatInt(valSize, 10), ir.I64),
				ir.NewConstant(strconv.FormatInt(int64(keyKind), 10), ir.I32),
			},
			Ret: ir.PtrI8,
		}
		dest := l.newTemp(ir.PtrI8)
		callInst.Dest = dest
		l.block.Append(callInst)
		return dest, nil
	}
	return nil, fmt.Errorf("make unsupported type %T", makeType.Underlying())
}

// lowerAppend supports append(slice, elem) for slices of any element type.
func (l *lowerer) lowerAppend(call *ast.CallExpr) (ir.Value, error) {
	if len(call.Args) != 2 {
		return nil, fmt.Errorf("append expects 2 arguments")
	}
	dstType := l.prog.TypesInfo.TypeOf(call)
	if dstType == nil {
		return nil, fmt.Errorf("missing type for append")
	}
	sliceType, ok := dstType.Underlying().(*types.Slice)
	if !ok {
		return nil, fmt.Errorf("append target is not a slice")
	}
	irSlice, err := goTypeToIR(sliceType)
	if err != nil {
		return nil, err
	}
	sliceTy, ok := irSlice.(ir.SliceType)
	if !ok {
		return nil, fmt.Errorf("append IR type is not slice")
	}
	elemTy := sliceTy.Elem
	elemSize, err := sizeofIRType(elemTy)
	if err != nil {
		return nil, err
	}

	sliceVal, err := l.lowerExpr(call.Args[0])
	if err != nil {
		return nil, err
	}
	if p, ok := ir.ValueType(sliceVal).(ir.PointerType); ok {
		tmp := l.newTemp(p.Elem)
		l.block.Append(ir.Load{Dest: tmp, Src: sliceVal, Align: alignOfIRType(p.Elem)})
		sliceVal = tmp
	}
	elemVal, err := l.lowerExpr(call.Args[1])
	if err != nil {
		return nil, err
	}

	// materialize slice fields
	slicePtr := l.allocaTemp(sliceTy)
	l.block.Append(ir.Store{Src: sliceVal, Dest: slicePtr, Align: alignOfIRType(sliceTy)})
	zero := ir.NewConstant("0", ir.I32)

	dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrTo(elemTy)))
	l.block.Append(ir.GetElementPtr{
		Dest:      dataPtrPtr,
		Src:       slicePtr,
		Pointee:   sliceTy,
		Indices:   []ir.Value{zero, zero},
		ResultTyp: ir.PtrTo(ir.PtrTo(elemTy)),
	})
	oldData := l.newTemp(ir.PtrTo(elemTy))
	l.block.Append(ir.Load{Dest: oldData, Src: dataPtrPtr, Align: alignOfIRType(ir.PtrTo(elemTy))})

	lenPtr := l.newTemp(ir.PtrTo(ir.I64))
	l.block.Append(ir.GetElementPtr{
		Dest:      lenPtr,
		Src:       slicePtr,
		Pointee:   sliceTy,
		Indices:   []ir.Value{zero, ir.NewConstant("1", ir.I32)},
		ResultTyp: ir.PtrTo(ir.I64),
	})
	oldLen, err := l.load(lenPtr)
	if err != nil {
		return nil, err
	}
	newLen := l.newTemp(ir.I64)
	l.block.Append(ir.BinOp{
		Dest: newLen,
		Op:   ir.Add,
		X:    oldLen,
		Y:    ir.NewConstant("1", ir.I64),
	})

	sizeConst := ir.NewConstant(strconv.FormatInt(elemSize, 10), ir.I64)
	copySize := l.newTemp(ir.I64)
	l.block.Append(ir.BinOp{
		Dest: copySize,
		Op:   ir.Mul,
		X:    oldLen,
		Y:    sizeConst,
	})

	rawPtr := l.newTemp(ir.PtrI8)
	l.block.Append(ir.Call{
		Dest: rawPtr,
		Name: "gominic_makeSlice",
		Args: []ir.Value{newLen, newLen, sizeConst},
		Ret:  ir.PtrI8,
	})
	newData := l.newTemp(ir.PtrTo(elemTy))
	l.block.Append(ir.Bitcast{
		Dest:   newData,
		Src:    rawPtr,
		Target: ir.PtrTo(elemTy),
	})

	l.block.Append(ir.Memcpy{Dest: newData, Src: oldData, Size: copySize})

	appendPtr := l.newTemp(ir.PtrTo(elemTy))
	l.block.Append(ir.GetElementPtr{
		Dest:      appendPtr,
		Src:       newData,
		Pointee:   elemTy,
		Indices:   []ir.Value{oldLen},
		ResultTyp: ir.PtrTo(elemTy),
	})
	if err := l.storeValue(appendPtr, elemVal); err != nil {
		return nil, err
	}

	return l.buildSliceValue(newData, newLen, newLen, elemTy)
}

func (l *lowerer) allocLocal(name string, t types.Type) (ir.Value, error) {
	irType, err := goTypeToIR(t)
	if err != nil {
		return nil, err
	}
	slotName := fmt.Sprintf("%s.addr%d", name, l.reg)
	l.reg++
	slot := ir.NewRegister(slotName, ir.PtrTo(irType))
	l.block.Append(ir.Alloca{Dest: slot, AllocType: irType})
	l.bindLocal(name, slot)
	return slot, nil
}

func (l *lowerer) allocaTemp(ty ir.Type) ir.Register {
	name := fmt.Sprintf("taddr%d", l.reg)
	l.reg++
	reg := ir.NewRegister(name, ir.PtrTo(ty))
	entry := l.fn.Entry()
	entry.Instrs = append(entry.Instrs, ir.Alloca{Dest: reg, AllocType: ty, Align: alignOfIRType(ty)})
	return reg
}

// addrOfTemp stores v into a fresh alloca and returns its address.
func (l *lowerer) addrOfTemp(v ir.Value) ir.Register {
	ptr := l.allocaTemp(ir.ValueType(v))
	l.block.Append(ir.Store{Src: v, Dest: ptr, Align: alignOfIRType(ir.ValueType(v))})
	return ptr
}

func (l *lowerer) load(ptr ir.Value) (ir.Value, error) {
	elem, err := pointeeType(ir.ValueType(ptr))
	if err != nil {
		return nil, err
	}
	dest := l.newTemp(elem)
	l.block.Append(ir.Load{Dest: dest, Src: ptr, Align: alignOfIRType(elem)})
	return dest, nil
}

func (l *lowerer) newTemp(ty ir.Type) ir.Register {
	name := fmt.Sprintf("t%d", l.reg)
	l.reg++
	return ir.NewRegister(name, ty)
}

// ensureI64 converts an integer value to i64 if needed.
func (l *lowerer) ensureI64(v ir.Value) (ir.Value, error) {
	if ir.ValueType(v) != nil && ir.ValueType(v).String() == ir.I64.String() {
		return v, nil
	}
	if ir.ValueType(v) != nil && ir.ValueType(v).String() == ir.I32.String() {
		dest := l.newTemp(ir.I64)
		l.block.Append(ir.Conv{
			Dest: dest,
			Op:   ir.SExt,
			Src:  v,
			To:   ir.I64,
		})
		return dest, nil
	}
	return nil, fmt.Errorf("expected integer index, got %s", ir.ValueType(v).String())
}

func (l *lowerer) newBlock(prefix string) *ir.BasicBlock {
	name := fmt.Sprintf("%s%d", prefix, l.blockID)
	l.blockID++
	bb := ir.NewBasicBlock(name)
	l.addBlock(bb)
	return bb
}

func (l *lowerer) addBlock(bb *ir.BasicBlock) {
	l.fn.Blocks = append(l.fn.Blocks, bb)
}

func (l *lowerer) resetScopes() {
	l.locals = nil
}

func (l *lowerer) pushScope() {
	l.locals = append(l.locals, make(map[string]ir.Value))
}

func (l *lowerer) popScope() {
	if len(l.locals) == 0 {
		return
	}
	l.locals = l.locals[:len(l.locals)-1]
}

func (l *lowerer) bindLocal(name string, v ir.Value) {
	if len(l.locals) == 0 {
		l.pushScope()
	}
	l.locals[len(l.locals)-1][name] = v
}

func (l *lowerer) lookupLocal(name string) (ir.Value, bool) {
	for i := len(l.locals) - 1; i >= 0; i-- {
		if v, ok := l.locals[i][name]; ok {
			return v, true
		}
	}
	return nil, false
}

func (l *lowerer) lowerUnary(u *ast.UnaryExpr) (ir.Value, error) {
	if u.Op == token.NOT {
		tv, ok := l.prog.TypesInfo.Types[u.X]
		if !ok {
			return nil, fmt.Errorf("no type info for unary !")
		}
		basic, ok := tv.Type.Underlying().(*types.Basic)
		if !ok || basic.Info()&types.IsBoolean == 0 {
			return nil, fmt.Errorf("unary ! expects bool")
		}
		val, err := l.lowerExpr(u.X)
		if err != nil {
			return nil, err
		}
		dest := l.newTemp(ir.I1)
		l.block.Append(ir.ICmp{
			Dest: dest,
			Pred: ir.ICmpEq,
			X:    val,
			Y:    ir.NewConstant("0", ir.I1),
		})
		return dest, nil
	}
	if u.Op == token.SUB {
		tv, ok := l.prog.TypesInfo.Types[u.X]
		if !ok {
			return nil, fmt.Errorf("no type info for unary -")
		}
		val, err := l.lowerExpr(u.X)
		if err != nil {
			return nil, err
		}
		irType, err := goTypeToIR(tv.Type)
		if err != nil {
			return nil, err
		}
		dest := l.newTemp(irType)
		basic, _ := tv.Type.Underlying().(*types.Basic)
		var op ir.BinOpKind
		if basic.Info()&types.IsFloat != 0 {
			op = ir.FSub
		} else {
			op = ir.Sub
		}
		zeroLit := "0"
		if irType != nil && irType.String() == ir.F64.String() {
			zeroLit = "0.0"
		}
		zero := ir.NewConstant(zeroLit, irType)
		l.block.Append(ir.BinOp{
			Dest: dest,
			Op:   op,
			X:    zero,
			Y:    val,
		})
		return dest, nil
	}
	if u.Op == token.ADD {
		return l.lowerExpr(u.X)
	}
	if u.Op == token.AND {
		return l.lowerAddr(u.X)
	}
	return nil, fmt.Errorf("unsupported unary op %s", u.Op.String())
}

func (l *lowerer) lowerCompositeLit(lit *ast.CompositeLit) (ir.Value, error) {
	t := l.prog.TypesInfo.TypeOf(lit)
	if t == nil {
		return nil, fmt.Errorf("missing type for composite literal")
	}
	if u, ok := t.Underlying().(*types.Struct); ok {
		return l.lowerStructLit(lit, u)
	}
	if u, ok := t.Underlying().(*types.Array); ok {
		return l.lowerArrayLit(lit, u)
	}
	if u, ok := t.Underlying().(*types.Slice); ok {
		return l.lowerSliceLit(lit, u)
	}
	return nil, fmt.Errorf("unsupported composite literal type %T", t.Underlying())
}

func (l *lowerer) lowerStructLit(lit *ast.CompositeLit, st *types.Struct) (ir.Value, error) {
	irType, err := goTypeToIR(st)
	if err != nil {
		return nil, err
	}
	fields := make([]ir.Value, st.NumFields())
	for i := range fields {
		fields[i], err = zeroValue(st.Field(i).Type())
		if err != nil {
			return nil, err
		}
	}
	next := 0
	for _, elt := range lit.Elts {
		if e, ok := elt.(*ast.KeyValueExpr); ok {
			id, ok := e.Key.(*ast.Ident)
			if !ok {
				return nil, fmt.Errorf("struct literal key is not ident")
			}
			index := -1
			for i := 0; i < st.NumFields(); i++ {
				if st.Field(i).Name() == id.Name {
					index = i
					break
				}
			}
			if index < 0 {
				return nil, fmt.Errorf("unknown field %s in struct literal", id.Name)
			}
			val, err := l.lowerExpr(e.Value)
			if err != nil {
				return nil, err
			}
			fields[index] = val
			continue
		}
		if next >= len(fields) {
			return nil, fmt.Errorf("too many values in struct literal")
		}
		val, err := l.lowerExpr(elt)
		if err != nil {
			return nil, err
		}
		fields[next] = val
		next++
	}

	ptr := l.allocaTemp(irType)
	for i, fv := range fields {
		zero := ir.NewConstant("0", ir.I32)
		fieldIdx := ir.NewConstant(strconv.Itoa(i), ir.I32)
		fieldPtr := l.newTemp(ir.PtrTo(ir.ValueType(fv)))
		l.block.Append(ir.GetElementPtr{
			Dest:      fieldPtr,
			Src:       ptr,
			Pointee:   irType,
			Indices:   []ir.Value{zero, fieldIdx},
			ResultTyp: ir.PtrTo(ir.ValueType(fv)),
		})
		l.block.Append(ir.Store{Src: fv, Dest: fieldPtr, Align: alignOfIRType(ir.ValueType(fv))})
	}

	return l.load(ptr)
}

func (l *lowerer) lowerArrayLit(lit *ast.CompositeLit, at *types.Array) (ir.Value, error) {
	irType, err := goTypeToIR(at)
	if err != nil {
		return nil, err
	}
	length := int(at.Len())
	elems := make([]ir.Value, length)
	for i := 0; i < length; i++ {
		elems[i], err = zeroValue(at.Elem())
		if err != nil {
			return nil, err
		}
	}
	next := 0
	for _, elt := range lit.Elts {
		if e, ok := elt.(*ast.KeyValueExpr); ok {
			idxLit, ok := e.Key.(*ast.BasicLit)
			if !ok || idxLit.Kind != token.INT {
				return nil, fmt.Errorf("array literal index must be int")
			}
			index, err := strconv.Atoi(idxLit.Value)
			if err != nil || index < 0 || index >= length {
				return nil, fmt.Errorf("array literal index out of range")
			}
			val, err := l.lowerExpr(e.Value)
			if err != nil {
				return nil, err
			}
			elems[index] = val
			continue
		}
		if next >= length {
			return nil, fmt.Errorf("too many values in array literal")
		}
		val, err := l.lowerExpr(elt)
		if err != nil {
			return nil, err
		}
		elems[next] = val
		next++
	}

	ptr := l.allocaTemp(irType)
	for i, ev := range elems {
		zero := ir.NewConstant("0", ir.I32)
		idx := ir.NewConstant(strconv.Itoa(i), ir.I32)
		elemPtr := l.newTemp(ir.PtrTo(ir.ValueType(ev)))
		l.block.Append(ir.GetElementPtr{
			Dest:      elemPtr,
			Src:       ptr,
			Pointee:   irType,
			Indices:   []ir.Value{zero, idx},
			ResultTyp: ir.PtrTo(ir.ValueType(ev)),
		})
		l.block.Append(ir.Store{Src: ev, Dest: elemPtr, Align: alignOfIRType(ir.ValueType(ev))})
	}

	return l.load(ptr)
}

func (l *lowerer) lowerSliceLit(lit *ast.CompositeLit, st *types.Slice) (ir.Value, error) {
	elemTy, err := goTypeToIR(st.Elem())
	if err != nil {
		return nil, err
	}
	count := len(lit.Elts)
	arrType := ir.ArrayType{Len: count, Elem: elemTy}
	arrPtr := l.allocaTemp(arrType)

	for i, elt := range lit.Elts {
		val, err := l.lowerExpr(elt)
		if err != nil {
			return nil, err
		}
		zero := ir.NewConstant("0", ir.I32)
		idx := ir.NewConstant(strconv.Itoa(i), ir.I32)
		elemPtr := l.newTemp(ir.PtrTo(elemTy))
		l.block.Append(ir.GetElementPtr{
			Dest:      elemPtr,
			Src:       arrPtr,
			Pointee:   arrType,
			Indices:   []ir.Value{zero, idx},
			ResultTyp: ir.PtrTo(elemTy),
		})
		l.block.Append(ir.Store{Src: val, Dest: elemPtr, Align: alignOfIRType(ir.ValueType(val))})
	}

	// data pointer to first element
	dataPtr := l.newTemp(ir.PtrTo(elemTy))
	zero := ir.NewConstant("0", ir.I32)
	l.block.Append(ir.GetElementPtr{
		Dest:      dataPtr,
		Src:       arrPtr,
		Pointee:   arrType,
		Indices:   []ir.Value{zero, zero},
		ResultTyp: ir.PtrTo(elemTy),
	})

	lenVal := ir.NewConstant(strconv.Itoa(count), ir.I64)
	return l.buildSliceValue(dataPtr, lenVal, lenVal, elemTy)
}

// lowerAddr возвращает адрес lvalue (идентификатор, селектор поля, индекс массива/среза).
func (l *lowerer) lowerAddr(expr ast.Expr) (ir.Value, error) {
	if e, ok := expr.(*ast.Ident); ok {
		if isRuntimeExtern(e.Name) {
			return nil, fmt.Errorf("cannot take address of runtime extern %s", e.Name)
		}
		if slot, ok := l.lookupLocal(e.Name); ok {
			return slot, nil
		}
		if obj := l.prog.TypesInfo.Uses[e]; obj != nil {
			if v, ok := obj.(*types.Var); ok {
				irType, err := goTypeToIR(v.Type())
				if err != nil {
					return nil, err
				}
				name := e.Name
				if l.prefix != "" {
					name = l.prefix + name
				}
				return ir.NewConstant("@"+name, ir.PtrTo(irType)), nil
			}
		}
		return nil, fmt.Errorf("address of unknown identifier %s", e.Name)
	}
	if e, ok := expr.(*ast.SelectorExpr); ok {
		return l.lowerSelectorAddr(e)
	}
	if e, ok := expr.(*ast.IndexExpr); ok {
		return l.lowerIndexAddr(e)
	}
	if e, ok := expr.(*ast.CompositeLit); ok {
		val, err := l.lowerCompositeLit(e)
		if err != nil {
			return nil, err
		}
		ptr := l.allocaTemp(ir.ValueType(val))
		l.block.Append(ir.Store{Src: val, Dest: ptr, Align: alignOfIRType(ir.ValueType(val))})
		return ptr, nil
	}
	if e, ok := expr.(*ast.TypeAssertExpr); ok {
		val, err := l.lowerExpr(e)
		if err != nil {
			return nil, err
		}
		ptr := l.allocaTemp(ir.ValueType(val))
		l.block.Append(ir.Store{Src: val, Dest: ptr, Align: alignOfIRType(ir.ValueType(val))})
		return ptr, nil
	}
	return nil, fmt.Errorf("unsupported address expression %T", expr)
}

func (l *lowerer) lowerSelectorAddr(sel *ast.SelectorExpr) (ir.Value, error) {
	// Package-level selector?
	if selInfo, okSel := l.prog.TypesInfo.Selections[sel]; !okSel || selInfo == nil {
		if obj, okUse := l.prog.TypesInfo.Uses[sel.Sel]; okUse {
			if v, okVar := obj.(*types.Var); okVar {
				irType, err := goTypeToIR(v.Type())
				if err != nil {
					return nil, err
				}
				pkg := ""
				if id, okID := sel.X.(*ast.Ident); okID {
					pkg = id.Name + "."
				}
				name := pkg + sel.Sel.Name
				return ir.NewConstant("@"+name, ir.PtrTo(irType)), nil
			}
		}
		return nil, fmt.Errorf("unsupported selector (no selection info)")
	}

	basePtr, err := l.lowerAddr(sel.X)
	if err != nil {
		return nil, err
	}

	selInfo := l.prog.TypesInfo.Selections[sel]
	if selInfo.Kind() != types.FieldVal {
		return nil, fmt.Errorf("unsupported selector (no selection info)")
	}
	path := selInfo.Index()
	if len(path) != 1 {
		return nil, fmt.Errorf("embedded fields not supported yet")
	}

	recv := selInfo.Recv()
	if ptr, ok := recv.(*types.Pointer); ok {
		recv = ptr.Elem()
	}
	irRecv, err := goTypeToIR(recv)
	if err != nil {
		return nil, err
	}
	irStruct, ok := irRecv.(ir.StructType)
	if !ok {
		// Fallback: produce zero temp for field type.
		fieldIR, errField := goTypeToIR(selInfo.Obj().Type())
		if errField != nil {
			return nil, fmt.Errorf("selector base is not a struct pointer")
		}
		tmp := l.allocaTemp(fieldIR)
		if z, errZero := zeroValue(selInfo.Obj().Type()); errZero == nil {
			l.block.Append(ir.Store{Src: z, Dest: tmp, Align: alignOfIRType(fieldIR)})
		}
		return tmp, nil
	}
	expectedPtr := ir.PtrTo(irStruct)
	// Insert bitcast to the expected struct pointer type to recover layout.
	casted := l.newTemp(expectedPtr)
	l.block.Append(ir.Bitcast{
		Dest:   casted,
		Src:    basePtr,
		Target: expectedPtr,
	})
	basePtr = casted
	if path[0] >= len(irStruct.Fields) {
		return nil, fmt.Errorf("field index out of range")
	}
	fieldTy := irStruct.Fields[path[0]]
	dest := l.newTemp(ir.PtrTo(fieldTy))

	zero := ir.NewConstant("0", ir.I32)
	fieldIdx := ir.NewConstant(strconv.Itoa(path[0]), ir.I32)
	l.block.Append(ir.GetElementPtr{
		Dest:      dest,
		Src:       basePtr,
		Pointee:   irStruct,
		Indices:   []ir.Value{zero, fieldIdx},
		ResultTyp: ir.PtrTo(fieldTy),
	})
	return dest, nil
}

// lowerIndexExpr handles indexing expressions, including maps.
func (l *lowerer) lowerIndexExpr(idx *ast.IndexExpr) (ir.Value, error) {
	tv := l.prog.TypesInfo.TypeOf(idx.X)
	if tv == nil {
		return nil, fmt.Errorf("missing type for index expression")
	}
	if _, ok := tv.Underlying().(*types.Map); ok {
		val, _, err := l.lowerMapGet(idx)
		return val, err
	}
	ptr, err := l.lowerIndexAddr(idx)
	if err != nil {
		return nil, err
	}
	return l.load(ptr)
}

func (l *lowerer) lowerIndexAddr(idx *ast.IndexExpr) (ir.Value, error) {
	baseType := l.prog.TypesInfo.TypeOf(idx.X)
	if baseType == nil {
		return nil, fmt.Errorf("missing type for index expression")
	}

	if t, ok := baseType.Underlying().(*types.Array); ok {
		basePtr, err := l.lowerAddr(idx.X)
		if err != nil {
			return nil, err
		}
		elemTy, err := goTypeToIR(t.Elem())
		if err != nil {
			return nil, err
		}
		indexVal, err := l.lowerExpr(idx.Index)
		if err != nil {
			return nil, err
		}
		indexVal, err = l.ensureIndexWithinBounds(indexVal, ir.NewConstant(strconv.FormatInt(t.Len(), 10), ir.I64))
		if err != nil {
			return nil, err
		}
		zero := ir.NewConstant("0", ir.I32)
		dest := l.newTemp(ir.PtrTo(elemTy))
		pointee, err := pointeeType(ir.ValueType(basePtr))
		if err != nil {
			return nil, err
		}
		l.block.Append(ir.GetElementPtr{
			Dest:      dest,
			Src:       basePtr,
			Pointee:   pointee,
			Indices:   []ir.Value{zero, indexVal},
			ResultTyp: ir.PtrTo(elemTy),
		})
		return dest, nil
	}
	if t, ok := baseType.Underlying().(*types.Slice); ok {
		basePtr, err := l.lowerAddr(idx.X)
		if err != nil {
			return nil, err
		}
		elemTy, err := goTypeToIR(t.Elem())
		if err != nil {
			return nil, err
		}
		indexVal, err := l.lowerExpr(idx.Index)
		if err != nil {
			return nil, err
		}
		if err := l.boundsCheckSlice(basePtr, indexVal); err != nil {
			return nil, err
		}

		zero := ir.NewConstant("0", ir.I32)
		dataField := ir.NewConstant("0", ir.I32)
		dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrTo(elemTy)))
		pointee, err := pointeeType(ir.ValueType(basePtr))
		if err != nil {
			return nil, err
		}
		l.block.Append(ir.GetElementPtr{
			Dest:      dataPtrPtr,
			Src:       basePtr,
			Pointee:   pointee,
			Indices:   []ir.Value{zero, dataField},
			ResultTyp: ir.PtrTo(ir.PtrTo(elemTy)),
		})
		dataPtr, err := l.load(dataPtrPtr)
		if err != nil {
			return nil, err
		}

		dest := l.newTemp(ir.PtrTo(elemTy))
		l.block.Append(ir.GetElementPtr{
			Dest:      dest,
			Src:       dataPtr,
			Pointee:   elemTy,
			Indices:   []ir.Value{indexVal},
			ResultTyp: ir.PtrTo(elemTy),
		})
		return dest, nil
	}
	if b, ok := baseType.Underlying().(*types.Basic); ok && (b.Kind() == types.String || b.Kind() == types.UntypedString) {
		strVal, err := l.lowerExpr(idx.X)
		if err != nil {
			return nil, err
		}
		if p, ok := ir.ValueType(strVal).(ir.PointerType); ok {
			tmp := l.newTemp(p.Elem)
			l.block.Append(ir.Load{Dest: tmp, Src: strVal, Align: alignOfIRType(p.Elem)})
			strVal = tmp
		}
		strPtr := l.allocaTemp(ir.ValueType(strVal))
		l.block.Append(ir.Store{Src: strVal, Dest: strPtr, Align: alignOfIRType(ir.ValueType(strVal))})
		zero := ir.NewConstant("0", ir.I32)
		// data ptr
		dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrI8))
		l.block.Append(ir.GetElementPtr{
			Dest:      dataPtrPtr,
			Src:       strPtr,
			Pointee:   ir.ValueType(strVal),
			Indices:   []ir.Value{zero, zero},
			ResultTyp: ir.PtrTo(ir.PtrI8),
		})
		dataPtr := l.newTemp(ir.PtrI8)
		l.block.Append(ir.Load{Dest: dataPtr, Src: dataPtrPtr, Align: alignOfIRType(ir.PtrI8)})
		indexVal, err := l.lowerExpr(idx.Index)
		if err != nil {
			return nil, err
		}
		index64, err := l.ensureI64(indexVal)
		if err != nil {
			return nil, err
		}
		dest := l.newTemp(ir.PtrTo(ir.I8))
	l.block.Append(ir.GetElementPtr{
		Dest:      dest,
		Src:       dataPtr,
		Pointee:   ir.BasicType("i8"),
		Indices:   []ir.Value{index64},
		ResultTyp: ir.PtrTo(ir.I8),
	})
	return dest, nil
}
	return nil, fmt.Errorf("index on unsupported type %T", baseType.Underlying())
}

// lowerMapSet emits runtime call to set m[k] = v.
func (l *lowerer) lowerMapSet(idx *ast.IndexExpr, val ir.Value) error {
	mapVal, err := l.lowerExpr(idx.X)
	if err != nil {
		return err
	}
	tv := l.prog.TypesInfo.TypeOf(idx.X)
	if tv == nil {
		return fmt.Errorf("missing type for map expression")
	}
	mapType, ok := tv.Underlying().(*types.Map)
	if !ok {
		return fmt.Errorf("map assignment on non-map")
	}
	keyIR, err := goTypeToIR(mapType.Key())
	if err != nil {
		return err
	}
	valIR, err := goTypeToIR(mapType.Elem())
	if err != nil {
		return err
	}
	keyTmp := l.allocaTemp(keyIR)
	keyVal, err := l.lowerExpr(idx.Index)
	if err != nil {
		return err
	}
	if err := l.storeValue(keyTmp, keyVal); err != nil {
		return err
	}
	valTmp := l.allocaTemp(valIR)
	if err := l.storeValue(valTmp, val); err != nil {
		return err
	}

	keyPtr := l.newTemp(ir.PtrI8)
	l.block.Append(ir.Bitcast{Dest: keyPtr, Src: keyTmp, Target: ir.PtrI8})
	valPtr := l.newTemp(ir.PtrI8)
	l.block.Append(ir.Bitcast{Dest: valPtr, Src: valTmp, Target: ir.PtrI8})

	l.block.Append(ir.Call{
		Name: "gominic_map_set",
		Args: []ir.Value{mapVal, keyPtr, valPtr},
		Ret:  ir.Void,
	})
	return nil
}

// lowerMapGet performs lookup and returns (value, ok).
func (l *lowerer) lowerMapGet(idx *ast.IndexExpr) (ir.Value, ir.Value, error) {
	mapVal, err := l.lowerExpr(idx.X)
	if err != nil {
		return nil, nil, err
	}
	tv := l.prog.TypesInfo.TypeOf(idx.X)
	if tv == nil {
		return nil, nil, fmt.Errorf("missing type for map expression")
	}
	mapType, ok := tv.Underlying().(*types.Map)
	if !ok {
		return nil, nil, fmt.Errorf("index on non-map")
	}
	keyIR, err := goTypeToIR(mapType.Key())
	if err != nil {
		return nil, nil, err
	}
	valIR, err := goTypeToIR(mapType.Elem())
	if err != nil {
		return nil, nil, err
	}

	keyTmp := l.allocaTemp(keyIR)
	keyVal, err := l.lowerExpr(idx.Index)
	if err != nil {
		return nil, nil, err
	}
	if err := l.storeValue(keyTmp, keyVal); err != nil {
		return nil, nil, err
	}
	valTmp := l.allocaTemp(valIR)
	if zero, err := zeroValue(mapType.Elem()); err == nil {
		if err := l.storeValue(valTmp, zero); err != nil {
			return nil, nil, err
		}
	}

	keyPtr := l.newTemp(ir.PtrI8)
	l.block.Append(ir.Bitcast{Dest: keyPtr, Src: keyTmp, Target: ir.PtrI8})
	valPtr := l.newTemp(ir.PtrI8)
	l.block.Append(ir.Bitcast{Dest: valPtr, Src: valTmp, Target: ir.PtrI8})
	okReg := l.newTemp(ir.I1)
	l.block.Append(ir.Call{
		Dest: okReg,
		Name: "gominic_map_get",
		Args: []ir.Value{mapVal, keyPtr, valPtr},
		Ret:  ir.I1,
	})
	valLoaded, err := l.load(valTmp)
	if err != nil {
		return nil, nil, err
	}
	return valLoaded, okReg, nil
}

// lowerMapGetAssign handles x, ok := m[k].
func (l *lowerer) lowerMapGetAssign(lhs []ast.Expr, idx *ast.IndexExpr, isDefine bool) error {
	val, okVal, err := l.lowerMapGet(idx)
	if err != nil {
		return err
	}
	targets := []struct {
		expr ast.Expr
		val  ir.Value
	}{
		{lhs[0], val},
		{lhs[1], okVal},
	}
	for _, t := range targets {
		if isDefine {
			if ident, ok := t.expr.(*ast.Ident); ok {
				if def := l.prog.TypesInfo.Defs[ident]; def != nil {
					slot, err := l.allocLocal(ident.Name, def.Type())
					if err != nil {
						return err
					}
					if err := l.storeValue(slot, t.val); err != nil {
						return err
					}
					continue
				}
			}
		}
		slot, err := l.lowerAddr(t.expr)
		if err != nil {
			return err
		}
		if err := l.storeValue(slot, t.val); err != nil {
			return err
		}
	}
	return nil
}
func (l *lowerer) lowerStringLiteral(val string) (ir.Value, error) {
	dataName := fmt.Sprintf(".str.%d", globalStringCounter)
	globalStringCounter++

	bytes := []byte(val)
	encoded := escapeCString(bytes)
	dataType := ir.ArrayType{Len: len(bytes) + 1, Elem: ir.BasicType("i8")}
	length := int64(len(bytes))

	l.mod.Globals = append(l.mod.Globals, ir.Global{
		Name:    dataName,
		Type:    dataType,
		Value:   fmt.Sprintf("c\"%s\"", encoded),
		Align:   alignOfIRType(dataType),
		Private: true,
	})

	gep := fmt.Sprintf("getelementptr inbounds (%s, %s* @%s, i32 0, i32 0)", dataType.String(), dataType.String(), dataName)
	strConst := fmt.Sprintf("{ %s, %s }", ir.PtrTo(ir.BasicType("i8")).String()+" "+gep, ir.I64.String()+" "+strconv.FormatInt(length, 10))
	return ir.NewConstant(strConst, ir.StringType{}), nil
}

func (l *lowerer) buildSliceValue(dataPtr ir.Value, lenVal ir.Value, capVal ir.Value, elemTy ir.Type) (ir.Value, error) {
	sliceTy := ir.SliceType{Elem: elemTy}
	ptr := l.allocaTemp(sliceTy)

	zero := ir.NewConstant("0", ir.I32)

	// data field
	dataField := ir.NewConstant("0", ir.I32)
	dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrTo(elemTy)))
	l.block.Append(ir.GetElementPtr{
		Dest:      dataPtrPtr,
		Src:       ptr,
		Pointee:   sliceTy,
		Indices:   []ir.Value{zero, dataField},
		ResultTyp: ir.PtrTo(ir.PtrTo(elemTy)),
	})
	l.block.Append(ir.Store{Src: dataPtr, Dest: dataPtrPtr, Align: alignOfIRType(ir.PtrTo(elemTy))})

	// len field (index 1)
	lenField := ir.NewConstant("1", ir.I32)
	lenPtr := l.newTemp(ir.PtrTo(ir.I64))
	l.block.Append(ir.GetElementPtr{
		Dest:      lenPtr,
		Src:       ptr,
		Pointee:   sliceTy,
		Indices:   []ir.Value{zero, lenField},
		ResultTyp: ir.PtrTo(ir.I64),
	})
	lenVal64, err := l.ensureI64(lenVal)
	if err != nil {
		return nil, err
	}
	l.block.Append(ir.Store{Src: lenVal64, Dest: lenPtr, Align: alignOfIRType(ir.I64)})

	// cap field (index 2)
	capField := ir.NewConstant("2", ir.I32)
	capPtr := l.newTemp(ir.PtrTo(ir.I64))
	l.block.Append(ir.GetElementPtr{
		Dest:      capPtr,
		Src:       ptr,
		Pointee:   sliceTy,
		Indices:   []ir.Value{zero, capField},
		ResultTyp: ir.PtrTo(ir.I64),
	})
	capVal64, err := l.ensureI64(capVal)
	if err != nil {
		return nil, err
	}
	l.block.Append(ir.Store{Src: capVal64, Dest: capPtr, Align: alignOfIRType(ir.I64)})

	return l.load(ptr)
}

func (l *lowerer) lowerConversion(call *ast.CallExpr, wantValue bool) (ir.Value, error) {
	if len(call.Args) != 1 {
		return nil, fmt.Errorf("conversion requires one argument")
	}
	dstType := l.prog.TypesInfo.TypeOf(call)
	if dstType == nil {
		return nil, fmt.Errorf("missing type for conversion")
	}
	src, err := l.lowerExpr(call.Args[0])
	if err != nil {
		return nil, err
	}
	dstIR, err := goTypeToIR(dstType)
	if err != nil {
		return nil, err
	}
	srcIR := ir.ValueType(src)

	// If types already match, return as-is.
	if dstIR.String() == srcIR.String() {
		return src, nil
	}

	// String -> slice of bytes conversion: map {i8*, i64} to slice {i8*, len, len}
	if _, ok := dstIR.(ir.SliceType); ok {
		if _, ok := srcIR.(ir.StringType); ok {
			strVal := src
			if p, ok := ir.ValueType(strVal).(ir.PointerType); ok {
				tmp := l.newTemp(p.Elem)
				l.block.Append(ir.Load{Dest: tmp, Src: strVal, Align: alignOfIRType(p.Elem)})
				strVal = tmp
			}
			// Allocate slice descriptor
			sliceTy := dstIR
			dest := l.allocaTemp(sliceTy)
			zero32 := ir.NewConstant("0", ir.I32)
			// data from string
			dataPtr := l.newTemp(ir.PtrTo(ir.PtrI8))
			l.block.Append(ir.GetElementPtr{
				Dest:      dataPtr,
				Src:       dest,
				Pointee:   sliceTy,
				Indices:   []ir.Value{zero32, zero32},
				ResultTyp: ir.PtrTo(ir.PtrI8),
			})
			strDataPtr := l.newTemp(ir.PtrTo(ir.PtrI8))
			l.block.Append(ir.GetElementPtr{
				Dest:      strDataPtr,
				Src:       l.addrOfTemp(strVal),
				Pointee:   ir.ValueType(strVal),
				Indices:   []ir.Value{zero32, zero32},
				ResultTyp: ir.PtrTo(ir.PtrI8),
			})
			strData := l.newTemp(ir.PtrI8)
			l.block.Append(ir.Load{Dest: strData, Src: strDataPtr, Align: alignOfIRType(ir.PtrI8)})
			l.block.Append(ir.Store{Src: strData, Dest: dataPtr, Align: alignOfIRType(ir.PtrI8)})
			// len/cap from string len
			lenPtr := l.newTemp(ir.PtrTo(ir.I64))
			l.block.Append(ir.GetElementPtr{
				Dest:      lenPtr,
				Src:       dest,
				Pointee:   sliceTy,
				Indices:   []ir.Value{zero32, ir.NewConstant("1", ir.I32)},
				ResultTyp: ir.PtrTo(ir.I64),
			})
			capPtr := l.newTemp(ir.PtrTo(ir.I64))
			l.block.Append(ir.GetElementPtr{
				Dest:      capPtr,
				Src:       dest,
				Pointee:   sliceTy,
				Indices:   []ir.Value{zero32, ir.NewConstant("2", ir.I32)},
				ResultTyp: ir.PtrTo(ir.I64),
			})
			strLenPtr := l.newTemp(ir.PtrTo(ir.I64))
			l.block.Append(ir.GetElementPtr{
				Dest:      strLenPtr,
				Src:       l.addrOfTemp(strVal),
				Pointee:   ir.ValueType(strVal),
				Indices:   []ir.Value{zero32, ir.NewConstant("1", ir.I32)},
				ResultTyp: ir.PtrTo(ir.I64),
			})
			strLen := l.newTemp(ir.I64)
			l.block.Append(ir.Load{Dest: strLen, Src: strLenPtr, Align: alignOfIRType(ir.I64)})
			l.block.Append(ir.Store{Src: strLen, Dest: lenPtr, Align: alignOfIRType(ir.I64)})
			l.block.Append(ir.Store{Src: strLen, Dest: capPtr, Align: alignOfIRType(ir.I64)})

			// return loaded slice value
			out := l.newTemp(sliceTy)
			l.block.Append(ir.Load{Dest: out, Src: dest, Align: alignOfIRType(sliceTy)})
			return out, nil
		}
	}

	// Slice of bytes -> string
	if _, ok := dstIR.(ir.StringType); ok {
		if sTy, ok := srcIR.(ir.SliceType); ok && sTy.Elem.String() == ir.I8.String() {
			sliceVal := src
			if p, ok := ir.ValueType(sliceVal).(ir.PointerType); ok {
				tmp := l.newTemp(p.Elem)
				l.block.Append(ir.Load{Dest: tmp, Src: sliceVal, Align: alignOfIRType(p.Elem)})
				sliceVal = tmp
			}
			slicePtr := l.allocaTemp(ir.ValueType(sliceVal))
			l.block.Append(ir.Store{Src: sliceVal, Dest: slicePtr, Align: alignOfIRType(ir.ValueType(sliceVal))})

			zero := ir.NewConstant("0", ir.I32)
			lenIdx := ir.NewConstant("1", ir.I32)

			dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrI8))
			l.block.Append(ir.GetElementPtr{
				Dest:      dataPtrPtr,
				Src:       slicePtr,
				Pointee:   ir.ValueType(sliceVal),
				Indices:   []ir.Value{zero, zero},
				ResultTyp: ir.PtrTo(ir.PtrI8),
			})
			dataPtr := l.newTemp(ir.PtrI8)
			l.block.Append(ir.Load{Dest: dataPtr, Src: dataPtrPtr, Align: alignOfIRType(ir.PtrI8)})

			lenPtr := l.newTemp(ir.PtrTo(ir.I64))
			l.block.Append(ir.GetElementPtr{
				Dest:      lenPtr,
				Src:       slicePtr,
				Pointee:   ir.ValueType(sliceVal),
				Indices:   []ir.Value{zero, lenIdx},
				ResultTyp: ir.PtrTo(ir.I64),
			})
			lenVal, err := l.load(lenPtr)
			if err != nil {
				return nil, err
			}

			strTy := dstIR
			dest := l.allocaTemp(strTy)
			dataField := l.newTemp(ir.PtrTo(ir.PtrI8))
			l.block.Append(ir.GetElementPtr{
				Dest:      dataField,
				Src:       dest,
				Pointee:   strTy,
				Indices:   []ir.Value{zero, zero},
				ResultTyp: ir.PtrTo(ir.PtrI8),
			})
			l.block.Append(ir.Store{Src: dataPtr, Dest: dataField, Align: alignOfIRType(ir.PtrI8)})
			lenField := l.newTemp(ir.PtrTo(ir.I64))
			l.block.Append(ir.GetElementPtr{
				Dest:      lenField,
				Src:       dest,
				Pointee:   strTy,
				Indices:   []ir.Value{zero, lenIdx},
				ResultTyp: ir.PtrTo(ir.I64),
			})
			l.block.Append(ir.Store{Src: lenVal, Dest: lenField, Align: alignOfIRType(ir.I64)})

			out := l.newTemp(strTy)
			l.block.Append(ir.Load{Dest: out, Src: dest, Align: alignOfIRType(strTy)})
			return out, nil
		}
	}

	// string -> *i8 (data pointer)
	if dstIR.String() == ir.PtrI8.String() {
		if _, ok := srcIR.(ir.StringType); ok {
			strVal := src
			if p, ok := ir.ValueType(strVal).(ir.PointerType); ok {
				tmp := l.newTemp(p.Elem)
				l.block.Append(ir.Load{Dest: tmp, Src: strVal, Align: alignOfIRType(p.Elem)})
				strVal = tmp
			}
			strPtr := l.allocaTemp(ir.ValueType(strVal))
			l.block.Append(ir.Store{Src: strVal, Dest: strPtr, Align: alignOfIRType(ir.ValueType(strVal))})
			zero := ir.NewConstant("0", ir.I32)
			dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrI8))
			l.block.Append(ir.GetElementPtr{
				Dest:      dataPtrPtr,
				Src:       strPtr,
				Pointee:   ir.ValueType(strVal),
				Indices:   []ir.Value{zero, zero},
				ResultTyp: ir.PtrTo(ir.PtrI8),
			})
			dataPtr := l.newTemp(ir.PtrI8)
			l.block.Append(ir.Load{Dest: dataPtr, Src: dataPtrPtr, Align: alignOfIRType(ir.PtrI8)})
			return dataPtr, nil
		}
	}

	convOp, err := selectConvOp(srcIR, dstIR, isUnsigned(dstType))
	if err != nil {
		return nil, err
	}
	dest := l.newTemp(dstIR)
	l.block.Append(ir.Conv{
		Dest: dest,
		Op:   convOp,
		Src:  src,
		To:   dstIR,
	})
	return dest, nil
}

// storeValue writes rhs into lhs pointer, using memcpy for aggregates.
func (l *lowerer) storeValue(lhs ir.Value, rhs ir.Value) error {
	lhsPointee, err := pointeeType(ir.ValueType(lhs))
	if err != nil {
		return err
	}
	if isAggregateType(lhsPointee) {
		size, err := sizeofIRType(lhsPointee)
		if err != nil {
			return err
		}
		rhsPtr := rhs
		if _, ok := ir.ValueType(rhs).(ir.PointerType); !ok {
			tmp := l.allocaTemp(lhsPointee)
			l.block.Append(ir.Store{Src: rhs, Dest: tmp})
			rhsPtr = tmp
		}
		l.block.Append(ir.Memcpy{
			Dest: lhs,
			Src:  rhsPtr,
			Size: ir.NewConstant(strconv.FormatInt(size, 10), ir.I64),
		})
		return nil
	}
	l.block.Append(ir.Store{Src: rhs, Dest: lhs, Align: alignOfIRType(ir.ValueType(rhs))})
	return nil
}

func isAggregateType(t ir.Type) bool {
	if _, ok := t.(ir.StructType); ok {
		return true
	}
	if _, ok := t.(ir.ArrayType); ok {
		return true
	}
	if _, ok := t.(ir.SliceType); ok {
		return true
	}
	if _, ok := t.(ir.StringType); ok {
		return true
	}
	return false
}

func sizeofIRType(t ir.Type) (int64, error) {
		if tt, ok := t.(ir.BasicType); ok {
			if tt == ir.I1 {
				return 1, nil
			}
			if tt == ir.I8 {
				return 1, nil
			}
			if tt == ir.I32 {
				return 4, nil
			}
			if tt.String() == ir.I64.String() || tt.String() == ir.F64.String() {
				return 8, nil
			}
			if tt.String() == ir.PtrI8.String() {
				return 8, nil
			}
			return 0, fmt.Errorf("sizeof unsupported basic type %s", tt.String())
		}
		if _, ok := t.(ir.PointerType); ok {
			return 8, nil
		}
	if _, ok := t.(ir.StringType); ok {
		return 16, nil
	}
	if _, ok := t.(ir.SliceType); ok {
		return 24, nil
	}
	if tt, ok := t.(ir.ArrayType); ok {
		elemSize, err := sizeofIRType(tt.Elem)
		if err != nil {
			return 0, err
		}
		elemAlign := alignOfIRType(tt.Elem)
		stride := alignUp(elemSize, elemAlign)
		return int64(tt.Len) * stride, nil
	}
	if tt, ok := t.(ir.StructType); ok {
		var total int64
		var maxAlign int64 = 1
		for _, f := range tt.Fields {
			fs, err := sizeofIRType(f)
			if err != nil {
				return 0, err
			}
			a := alignOfIRType(f)
			if a > maxAlign {
				maxAlign = a
			}
			total = alignUp(total, a)
			total += fs
		}
		total = alignUp(total, maxAlign)
		return total, nil
	}
	return 0, fmt.Errorf("sizeof unsupported IR type %T", t)
}

func alignOfIRType(t ir.Type) int64 {
		if tt, ok := t.(ir.BasicType); ok {
			if tt == ir.I1 {
				return 1
			}
			if tt == ir.I8 {
				return 1
			}
			if tt == ir.I32 {
				return 4
			}
			if tt.String() == ir.I64.String() || tt.String() == ir.F64.String() {
				return 8
			}
			if tt.String() == ir.PtrI8.String() {
				return 8
			}
			return 1
		}
	if _, ok := t.(ir.PointerType); ok {
		return 8
	}
	if _, ok := t.(ir.StringType); ok {
		return 8
	}
	if _, ok := t.(ir.SliceType); ok {
		return 8
	}
	if tt, ok := t.(ir.ArrayType); ok {
		return alignOfIRType(tt.Elem)
	}
	if tt, ok := t.(ir.StructType); ok {
		var max int64 = 1
		for _, f := range tt.Fields {
			a := alignOfIRType(f)
			if a > max {
				max = a
			}
		}
		return max
	}
	return 1
}

func mustSizeofIRType(t ir.Type) int64 {
	n, err := sizeofIRType(t)
	if err != nil {
		return 0
	}
	return n
}

func selectConvOp(src ir.Type, dst ir.Type, srcUnsigned bool) (ir.ConvOp, error) {
	if s, ok := src.(ir.BasicType); ok {
		if d, ok := dst.(ir.BasicType); ok {
			if s == ir.I32 && d == ir.I64 {
				if srcUnsigned {
					return ir.ZExt, nil
				}
				return ir.SExt, nil
			}
			if s == ir.I64 && d == ir.I32 {
				return ir.Trunc, nil
			}
			if s == ir.I1 && (d == ir.I32 || d == ir.I64) {
				return ir.ZExt, nil
			}
			if (s == ir.I32 || s == ir.I64) && d == ir.F64 {
				if srcUnsigned {
					return ir.UIToFP, nil
				}
				return ir.SIToFP, nil
			}
			if s == ir.F64 && (d == ir.I32 || d == ir.I64) {
				return ir.FPToSI, nil
			}
		}
	}
	return "", fmt.Errorf("unsupported conversion from %s to %s", src.String(), dst.String())
}

func escapeCString(b []byte) string {
	var sb strings.Builder
	for _, ch := range b {
		fmt.Fprintf(&sb, "\\%02X", ch)
	}
	sb.WriteString("\\00")
	return sb.String()
}

func (l *lowerer) lowerCall(call *ast.CallExpr, wantValue bool) (ir.Value, error) {
	var fnName string
	if f, ok := call.Fun.(*ast.Ident); ok {
		fnName = f.Name
	} else if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if pkgIdent, ok := sel.X.(*ast.Ident); ok {
			fnName = pkgIdent.Name + "." + sel.Sel.Name
		} else {
			return l.lowerMethodCall(sel, call, wantValue)
		}
	} else {
		return l.lowerConversion(call, wantValue)
	}
	if l.prefix != "" {
		if !strings.Contains(fnName, ".") && !isRuntimeExtern(fnName) {
			fnName = l.prefix + fnName
		}
	}

	if fnName == "make" {
		return l.lowerMake(call, wantValue)
	}
	if fnName == "len" {
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("len expects 1 argument")
		}
		arg := call.Args[0]
		tv := l.prog.TypesInfo.TypeOf(arg)
		if tv == nil {
			return nil, fmt.Errorf("missing type for len arg")
		}
			if u, ok := tv.Underlying().(*types.Array); ok {
				return ir.NewConstant(strconv.FormatInt(u.Len(), 10), ir.I64), nil
			}
			if _, ok := tv.Underlying().(*types.Slice); ok {
				sliceVal, err := l.lowerExpr(arg)
				if err != nil {
					return nil, err
				}
				slicePtr := l.allocaTemp(ir.ValueType(sliceVal))
			l.block.Append(ir.Store{Src: sliceVal, Dest: slicePtr, Align: alignOfIRType(ir.ValueType(sliceVal))})
			zero := ir.NewConstant("0", ir.I32)
			lenIdx := ir.NewConstant("1", ir.I32)
				lenPtr := l.newTemp(ir.PtrTo(ir.I64))
				l.block.Append(ir.GetElementPtr{
					Dest:      lenPtr,
					Src:       slicePtr,
					Pointee:   ir.ValueType(sliceVal),
					Indices:   []ir.Value{zero, lenIdx},
					ResultTyp: ir.PtrTo(ir.I64),
				})
				return l.load(lenPtr)
			}
			if _, ok := tv.Underlying().(*types.Map); ok {
				mapVal, err := l.lowerExpr(arg)
				if err != nil {
					return nil, err
				}
				dest := l.newTemp(ir.I64)
			l.block.Append(ir.Call{
				Dest: dest,
				Name: "gominic_map_len",
				Args: []ir.Value{mapVal},
				Ret:  ir.I64,
			})
			return dest, nil
			}
			if b, ok := tv.Underlying().(*types.Basic); ok {
				if b.Kind() == types.String || b.Kind() == types.UntypedString {
					strVal, err := l.lowerExpr(arg)
					if err != nil {
						return nil, err
					}
					strPtr := l.allocaTemp(ir.ValueType(strVal))
					l.block.Append(ir.Store{Src: strVal, Dest: strPtr, Align: alignOfIRType(ir.ValueType(strVal))})
					zero := ir.NewConstant("0", ir.I32)
					lenIdx := ir.NewConstant("1", ir.I32)
					lenPtr := l.newTemp(ir.PtrTo(ir.I64))
					l.block.Append(ir.GetElementPtr{
						Dest:      lenPtr,
						Src:       strPtr,
						Pointee:   ir.ValueType(strVal),
						Indices:   []ir.Value{zero, lenIdx},
						ResultTyp: ir.PtrTo(ir.I64),
					})
					return l.load(lenPtr)
				}
			}
			return nil, fmt.Errorf("len unsupported for type %T", tv.Underlying())
		}
	
	if fnName == "append" {
		return l.lowerAppend(call)
	}

	sigType := l.prog.TypesInfo.TypeOf(call.Fun)
	if sigType != nil {
		if _, isSig := sigType.(*types.Signature); !isSig {
			// Treat as conversion if the callee is a type, not a function.
			return l.lowerConversion(call, wantValue)
		}
	}
	var sig *types.Signature
	if s, ok := sigType.(*types.Signature); ok {
		sig = s
	}
	if sig == nil {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if obj := l.prog.TypesInfo.Uses[sel.Sel]; obj != nil {
				if s, okSig := obj.Type().(*types.Signature); okSig {
					sig = s
				}
			}
		}
	}
	if sig == nil {
		if id, okID := call.Fun.(*ast.Ident); okID {
			if obj := l.prog.TypesInfo.Uses[id]; obj != nil {
				if s, okSig := obj.Type().(*types.Signature); okSig {
					sig = s
				}
			}
			if sig == nil {
				if obj := l.prog.TypesInfo.Defs[id]; obj != nil {
					if s, okSig := obj.Type().(*types.Signature); okSig {
						sig = s
					}
				}
			}
		}
	}
	if sig == nil {
		// Synthetic signatures for known helper wrappers without type info.
		if fnName == "print" && len(call.Args) == 1 {
			strTyp := types.Typ[types.String]
			p := types.NewTuple(types.NewParam(token.NoPos, nil, "s", strTyp))
			sig = types.NewSignature(nil, p, types.NewTuple(), false)
		} else if fnName == "printInt" && len(call.Args) == 1 {
			p := types.NewTuple(types.NewParam(token.NoPos, nil, "v", types.Typ[types.Int64]))
			sig = types.NewSignature(nil, p, types.NewTuple(), false)
		} else if fnName == "println" && len(call.Args) == 0 {
			sig = types.NewSignature(nil, types.NewTuple(), types.NewTuple(), false)
		}
	}
	if sig == nil {
		// Allow known runtime externs even if type info is missing.
		if sigType == nil && isRuntimeExtern(fnName) {
			// handle below via manual lowering
		} else {
			// fallback: treat as void function with no signature
			sig = nil
		}
	}
	if sig != nil {
		if wantValue && sig.Results().Len() == 0 {
			return nil, fmt.Errorf("void call used in expression context")
		}
		paramsLen := sig.Params().Len()
		if sig.Variadic() {
			if len(call.Args) < paramsLen-1 {
				return nil, fmt.Errorf("argument count mismatch in call to %s", fnName)
			}
		} else {
			if paramsLen != len(call.Args) {
				return nil, fmt.Errorf("argument count mismatch in call to %s", fnName)
			}
		}
	}

	if sig == nil && isRuntimeExtern(fnName) {
		// Manual lowering for runtime externs when type info is missing.
		args := make([]ir.Value, 0, len(call.Args))
		for i := 0; i < len(call.Args); i++ {
			v, err := l.lowerExpr(call.Args[i])
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		}
		var retTy ir.Type = ir.Void
		switch fnName {
		case "gominic_print":
			retTy = ir.Void
		case "gominic_printInt":
			retTy = ir.Void
		case "gominic_println":
			retTy = ir.Void
		case "gominic_makeSlice":
			retTy = ir.PtrI8
		case "gominic_map_new":
			retTy = ir.PtrI8
		case "gominic_map_set":
			retTy = ir.Void
		case "gominic_map_get":
			retTy = ir.I1
		case "gominic_map_len":
			retTy = ir.I64
		case "gominic_str_eq":
			retTy = ir.I1
		}
		callInst := ir.Call{Name: fnName, Args: args, Ret: retTy}
		if retTy != nil && retTy.String() != ir.Void.String() && wantValue {
			dest := l.newTemp(retTy)
			callInst.Dest = dest
			l.block.Append(callInst)
			return dest, nil
		}
		l.block.Append(callInst)
		return nil, nil
	}

	argVals := make([]ir.Value, len(call.Args))
	for i, a := range call.Args {
		v, err := l.lowerExpr(a)
		if err != nil {
			return nil, err
		}
		argVals[i] = v
	}

	// Fallback: no signature, assume void call with lowered args.
	if sig == nil {
		callInst := ir.Call{Name: fnName, Args: argVals, Ret: ir.Void}
		l.block.Append(callInst)
		if wantValue {
			return nil, fmt.Errorf("cannot use call result of %s in expression", fnName)
		}
		return nil, nil
	}

	paramIR := l.lowerParamTypes(sig)
	if sig != nil && sig.Variadic() && len(call.Args) > len(paramIR) {
		if len(paramIR) > 0 {
			last := paramIR[len(paramIR)-1]
			for i := len(paramIR); i < len(call.Args); i++ {
				paramIR = append(paramIR, last)
			}
		}
	}
	args := make([]ir.Value, 0, len(paramIR))
	for i, p := range paramIR {
		if i >= len(argVals) {
			break
		}
		val := argVals[i]
		if p.ByVal() {
			if ptr, ok := ir.ValueType(val).(ir.PointerType); ok && ptr.Elem.String() == p.ByValType().String() {
				args = append(args, val)
			} else {
				tmp := l.allocaTemp(p.ByValType())
				l.block.Append(ir.Store{Src: val, Dest: tmp})
				args = append(args, tmp)
			}
		} else {
			args = append(args, val)
		}
	}

	var retTy ir.Type = ir.Void
	if sig.Results().Len() == 0 {
		retTy = ir.Void
	} else if sig.Results().Len() == 1 {
		var err error
		retTy, err = goTypeToIR(sig.Results().At(0).Type())
		if err != nil {
			return nil, err
		}
	} else {
		fields := make([]ir.Type, sig.Results().Len())
		for i := 0; i < sig.Results().Len(); i++ {
			t, err := goTypeToIR(sig.Results().At(i).Type())
			if err != nil {
				return nil, err
			}
			fields[i] = t
		}
		retTy = ir.StructType{Fields: fields}
	}

	callInst := ir.Call{
		Name: fnName,
		Args: args,
		Ret:  retTy,
	}

	if retTy != nil && retTy.String() != ir.Void.String() && wantValue {
		dest := l.newTemp(retTy)
		callInst.Dest = dest
		l.block.Append(callInst)
		return dest, nil
	}

	l.block.Append(callInst)
	return nil, nil
}

// lowerMethodCall lowers a method invocation X.M(...), treating it as a function with receiver as first arg.
func (l *lowerer) lowerMethodCall(sel *ast.SelectorExpr, call *ast.CallExpr, wantValue bool) (ir.Value, error) {
	selInfo, ok := l.prog.TypesInfo.Selections[sel]
	if !ok {
		return nil, fmt.Errorf("missing selection info for method call")
	}
	methodSig, ok := selInfo.Obj().Type().(*types.Signature)
	if !ok {
		return nil, fmt.Errorf("method has unexpected type %T", selInfo.Obj().Type())
	}

	recvVal, err := l.lowerExpr(sel.X)
	if err != nil {
		return nil, err
	}

	argVals := make([]ir.Value, 0, len(call.Args)+1)
	argVals = append(argVals, recvVal)
	for i := 0; i < len(call.Args); i++ {
		v, err := l.lowerExpr(call.Args[i])
		if err != nil {
			return nil, err
		}
		argVals = append(argVals, v)
	}

	paramIR := l.lowerParamTypes(methodSig)
	args := make([]ir.Value, 0, len(paramIR))
	for i := 0; i < len(paramIR); i++ {
		p := paramIR[i]
		val := argVals[i]
		if p.ByVal() {
			if ptr, ok := ir.ValueType(val).(ir.PointerType); ok && ptr.Elem.String() == p.ByValType().String() {
				args = append(args, val)
			} else {
				tmp := l.allocaTemp(p.ByValType())
				l.block.Append(ir.Store{Src: val, Dest: tmp})
				args = append(args, tmp)
			}
		} else {
			args = append(args, val)
		}
	}

	var retTy ir.Type = ir.Void
	if methodSig.Results().Len() == 0 {
		retTy = ir.Void
	} else if methodSig.Results().Len() == 1 {
		retTy, err = goTypeToIR(methodSig.Results().At(0).Type())
		if err != nil {
			return nil, err
		}
	} else {
		fields := make([]ir.Type, methodSig.Results().Len())
		for i := 0; i < methodSig.Results().Len(); i++ {
			t, err := goTypeToIR(methodSig.Results().At(i).Type())
			if err != nil {
				return nil, err
			}
			fields[i] = t
		}
		retTy = ir.StructType{Fields: fields}
	}

	name := recvBaseName(methodSig.Recv().Type()) + "." + sel.Sel.Name
	if l.prefix != "" {
		name = l.prefix + name
	}
	callInst := ir.Call{
		Name: name,
		Args: args,
		Ret:  retTy,
	}

	if retTy != nil && retTy.String() != ir.Void.String() && wantValue {
		dest := l.newTemp(retTy)
		callInst.Dest = dest
		l.block.Append(callInst)
		return dest, nil
	}

	l.block.Append(callInst)
	return nil, nil
}

func (l *lowerer) pushLoop(ctx loopContext) {
	l.loops = append(l.loops, ctx)
}

func (l *lowerer) popLoop() {
	if len(l.loops) == 0 {
		return
	}
	l.loops = l.loops[:len(l.loops)-1]
}

func goTypeToIR(t types.Type) (ir.Type, error) {
	return goTypeToIRRec(t, make(map[types.Type]bool))
}

// goTypeToIRRec converts Go types to our IR types while guarding against recursion.
func goTypeToIRRec(t types.Type, seen map[types.Type]bool) (ir.Type, error) {
	if u, ok := t.(*types.Named); ok {
		if obj := u.Obj(); obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "gominic/ir" {
			// For ir package helper enums/aliases, try underlying when it's basic to keep correct ABI.
			if _, ok := u.Underlying().(*types.Basic); ok {
				return goTypeToIRRec(u.Underlying(), seen)
			}
			// Treat other ir package type descriptors as opaque pointers at runtime.
			return ir.PtrI8, nil
		}
		if seen[u] {
			return ir.PtrI8, nil
		}
		seen[u] = true
		irType, err := goTypeToIRRec(u.Underlying(), seen)
		delete(seen, u)
		return irType, err
	}
	if u, ok := t.(*types.Basic); ok {
		kind := u.Kind()
		if kind == types.UntypedNil {
			return ir.PtrI8, nil
		}
		if kind == types.Bool || kind == types.UntypedBool {
			return ir.I1, nil
		}
		if kind == types.Int8 || kind == types.Uint8 {
			return ir.I8, nil
		}
		if kind == types.Int16 || kind == types.Uint16 {
			return ir.I32, nil
		}
		if kind == types.Int || kind == types.Int64 || kind == types.UntypedInt || kind == types.Uint || kind == types.Uintptr {
			return ir.I64, nil
		}
		if kind == types.Int32 || kind == types.Uint32 {
			return ir.I32, nil
		}
		if kind == types.Uint64 {
			return ir.I64, nil
		}
		if kind == types.Float64 || kind == types.UntypedFloat {
			return ir.F64, nil
		}
		if kind == types.String || kind == types.UntypedString {
			return ir.StringType{}, nil
		}
		if kind == types.UnsafePointer {
			return ir.PtrI8, nil
		}
		return nil, fmt.Errorf("unsupported basic type %v", u.Kind())
	}
	if u, ok := t.(*types.Pointer); ok {
		elem, err := goTypeToIRRec(u.Elem(), seen)
		if err != nil {
			return nil, err
		}
		return ir.PtrTo(elem), nil
	}
	if u, ok := t.(*types.Array); ok {
		elem, err := goTypeToIRRec(u.Elem(), seen)
		if err != nil {
			return nil, err
		}
		return ir.ArrayType{Len: int(u.Len()), Elem: elem}, nil
	}
	if u, ok := t.(*types.Slice); ok {
		elem, err := goTypeToIRRec(u.Elem(), seen)
		if err != nil {
			return nil, err
		}
		return ir.SliceType{Elem: elem}, nil
	}
	if u, ok := t.(*types.Struct); ok {
		fields := make([]ir.Type, u.NumFields())
		for i := 0; i < u.NumFields(); i++ {
			ft, err := goTypeToIRRec(u.Field(i).Type(), seen)
			if err != nil {
				return nil, err
			}
			fields[i] = ft
		}
		return ir.StructType{Fields: fields}, nil
	}
	if u, ok := t.(*types.Tuple); ok {
		if u.Len() == 0 {
			return ir.Void, nil
		}
		if u.Len() == 1 {
			return goTypeToIRRec(u.At(0).Type(), seen)
		}
		fields := make([]ir.Type, u.Len())
		for i := 0; i < u.Len(); i++ {
			t0, err := goTypeToIRRec(u.At(i).Type(), seen)
			if err != nil {
				return nil, err
			}
			fields[i] = t0
		}
		return ir.StructType{Fields: fields}, nil
	}
	if _, ok := t.(*types.Map); ok {
		return ir.PtrI8, nil
	}
	if _, ok := t.(*types.Interface); ok {
		return nil, fmt.Errorf("unsupported interface type")
	}
	if _, ok := t.(*types.Chan); ok {
		return ir.PtrI8, nil
	}
	if _, ok := t.(*types.Signature); ok {
		return ir.PtrI8, nil
	}
	return nil, fmt.Errorf("unsupported type %T", t)
}

func (l *lowerer) lowerSliceExpr(se *ast.SliceExpr) (ir.Value, error) {
	tv := l.prog.TypesInfo.TypeOf(se.X)
	if tv == nil {
		return nil, fmt.Errorf("missing type for slice expression")
	}
	u := tv.Underlying()
	if basic, ok := u.(*types.Basic); ok && (basic.Kind() == types.String || basic.Kind() == types.UntypedString) {
		// Slice of string -> string
		strVal, err := l.lowerExpr(se.X)
		if err != nil {
			return nil, err
		}
		// Ensure we have the value, not a pointer
		if p, ok := ir.ValueType(strVal).(ir.PointerType); ok {
			tmp := l.newTemp(p.Elem)
			l.block.Append(ir.Load{Dest: tmp, Src: strVal, Align: alignOfIRType(p.Elem)})
			strVal = tmp
		}
		strPtr := l.allocaTemp(ir.ValueType(strVal))
		l.block.Append(ir.Store{Src: strVal, Dest: strPtr, Align: alignOfIRType(ir.ValueType(strVal))})

		zero := ir.NewConstant("0", ir.I32)
		lenIdx := ir.NewConstant("1", ir.I32)
		// load data ptr
		dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrI8))
		l.block.Append(ir.GetElementPtr{
			Dest:      dataPtrPtr,
			Src:       strPtr,
			Pointee:   ir.ValueType(strVal),
			Indices:   []ir.Value{zero, zero},
			ResultTyp: ir.PtrTo(ir.PtrI8),
		})
		dataPtr := l.newTemp(ir.PtrI8)
		l.block.Append(ir.Load{Dest: dataPtr, Src: dataPtrPtr, Align: alignOfIRType(ir.PtrI8)})
		// load len
		lenPtr := l.newTemp(ir.PtrTo(ir.I64))
		l.block.Append(ir.GetElementPtr{
			Dest:      lenPtr,
			Src:       strPtr,
			Pointee:   ir.ValueType(strVal),
			Indices:   []ir.Value{zero, lenIdx},
			ResultTyp: ir.PtrTo(ir.I64),
		})
		lenVal, err := l.load(lenPtr)
		if err != nil {
			return nil, err
		}

		var low ir.Value = ir.NewConstant("0", ir.I64)
		if se.Low != nil {
			v, err := l.lowerExpr(se.Low)
			if err != nil {
				return nil, err
			}
			low, err = l.ensureI64(v)
			if err != nil {
				return nil, err
			}
		}
		high := lenVal
		if se.High != nil {
			v, err := l.lowerExpr(se.High)
			if err != nil {
				return nil, err
			}
			high, err = l.ensureI64(v)
			if err != nil {
				return nil, err
			}
		}
		// new len = high - low
		newLen := l.newTemp(ir.I64)
		l.block.Append(ir.BinOp{
			Dest: newLen,
			Op:   ir.Sub,
			X:    high,
			Y:    low,
		})
		// new data = dataPtr + low
		newData := l.newTemp(ir.PtrI8)
		l.block.Append(ir.GetElementPtr{
			Dest:      newData,
			Src:       dataPtr,
			Pointee:   ir.BasicType("i8"),
			Indices:   []ir.Value{low},
			ResultTyp: ir.PtrI8,
		})
		// assemble result string
		strTy := ir.ValueType(strVal)
		dest := l.allocaTemp(strTy)
		// store data
		dataFieldPtr := l.newTemp(ir.PtrTo(ir.PtrI8))
		l.block.Append(ir.GetElementPtr{
			Dest:      dataFieldPtr,
			Src:       dest,
			Pointee:   strTy,
			Indices:   []ir.Value{zero, zero},
			ResultTyp: ir.PtrTo(ir.PtrI8),
		})
		l.block.Append(ir.Store{Src: newData, Dest: dataFieldPtr, Align: alignOfIRType(ir.PtrI8)})
		// store len
		lenFieldPtr := l.newTemp(ir.PtrTo(ir.I64))
		l.block.Append(ir.GetElementPtr{
			Dest:      lenFieldPtr,
			Src:       dest,
			Pointee:   strTy,
			Indices:   []ir.Value{zero, lenIdx},
			ResultTyp: ir.PtrTo(ir.I64),
		})
		l.block.Append(ir.Store{Src: newLen, Dest: lenFieldPtr, Align: alignOfIRType(ir.I64)})

		result := l.newTemp(strTy)
		l.block.Append(ir.Load{Dest: result, Src: dest, Align: alignOfIRType(strTy)})
		return result, nil
	}
	if sl, ok := u.(*types.Slice); ok {
		sliceVal, err := l.lowerExpr(se.X)
		if err != nil {
			return nil, err
		}
		// Ensure value not pointer
		if p, ok := ir.ValueType(sliceVal).(ir.PointerType); ok {
			tmp := l.newTemp(p.Elem)
			l.block.Append(ir.Load{Dest: tmp, Src: sliceVal, Align: alignOfIRType(p.Elem)})
			sliceVal = tmp
		}
		slicePtr := l.allocaTemp(ir.ValueType(sliceVal))
		l.block.Append(ir.Store{Src: sliceVal, Dest: slicePtr, Align: alignOfIRType(ir.ValueType(sliceVal))})

		zero32 := ir.NewConstant("0", ir.I32)
		// data
		dataPtrPtr := l.newTemp(ir.PtrTo(ir.PtrTo(ir.I8))) // temporary type; will bitcast later
		l.block.Append(ir.GetElementPtr{
			Dest:      dataPtrPtr,
			Src:       slicePtr,
			Pointee:   ir.ValueType(sliceVal),
			Indices:   []ir.Value{zero32, zero32},
			ResultTyp: ir.PtrTo(ir.PtrTo(ir.I8)),
		})
		rawData := l.newTemp(ir.PtrI8)
		l.block.Append(ir.Load{Dest: rawData, Src: dataPtrPtr, Align: alignOfIRType(ir.PtrI8)})
		elemIR, err := goTypeToIR(sl.Elem())
		if err != nil {
			return nil, err
		}
		dataPtr := l.newTemp(ir.PtrTo(elemIR))
		l.block.Append(ir.Bitcast{Dest: dataPtr, Src: rawData, Target: ir.PtrTo(elemIR)})

		// len
		lenIdx := ir.NewConstant("1", ir.I32)
		lenPtr := l.newTemp(ir.PtrTo(ir.I64))
		l.block.Append(ir.GetElementPtr{
			Dest:      lenPtr,
			Src:       slicePtr,
			Pointee:   ir.ValueType(sliceVal),
			Indices:   []ir.Value{zero32, lenIdx},
			ResultTyp: ir.PtrTo(ir.I64),
		})
		lenVal, err := l.load(lenPtr)
		if err != nil {
			return nil, err
		}
		// cap
		capIdx := ir.NewConstant("2", ir.I32)
		capPtr := l.newTemp(ir.PtrTo(ir.I64))
		l.block.Append(ir.GetElementPtr{
			Dest:      capPtr,
			Src:       slicePtr,
			Pointee:   ir.ValueType(sliceVal),
			Indices:   []ir.Value{zero32, capIdx},
			ResultTyp: ir.PtrTo(ir.I64),
		})
		capVal, err := l.load(capPtr)
		if err != nil {
			return nil, err
		}

		var low ir.Value = ir.NewConstant("0", ir.I64)
		if se.Low != nil {
			v, err := l.lowerExpr(se.Low)
			if err != nil {
				return nil, err
			}
			low, err = l.ensureI64(v)
			if err != nil {
				return nil, err
			}
		}
		high := lenVal
		if se.High != nil {
			v, err := l.lowerExpr(se.High)
			if err != nil {
				return nil, err
			}
			high, err = l.ensureI64(v)
			if err != nil {
				return nil, err
			}
		}
		newLen := l.newTemp(ir.I64)
		l.block.Append(ir.BinOp{Dest: newLen, Op: ir.Sub, X: high, Y: low})
		newCap := l.newTemp(ir.I64)
		l.block.Append(ir.BinOp{Dest: newCap, Op: ir.Sub, X: capVal, Y: low})

		// adjust data pointer
		adjData := l.newTemp(ir.PtrTo(elemIR))
		l.block.Append(ir.GetElementPtr{
			Dest:      adjData,
			Src:       dataPtr,
			Pointee:   elemIR,
			Indices:   []ir.Value{low},
			ResultTyp: ir.PtrTo(elemIR),
		})

		return l.buildSliceValue(adjData, newLen, newCap, elemIR)
	}
	return nil, fmt.Errorf("unsupported slice expr on type %T", u)
}
func selectBinOp(tok token.Token, t types.Type) (ir.BinOpKind, error) {
	b, ok := t.Underlying().(*types.Basic)
	if !ok {
		return "", fmt.Errorf("unsupported binary expr type %T", t)
	}
	if tok == token.ADD {
		if b.Info()&types.IsFloat != 0 {
			return ir.FAdd, nil
		}
		return ir.Add, nil
	}
	if tok == token.SUB {
		if b.Info()&types.IsFloat != 0 {
			return ir.FSub, nil
		}
		return ir.Sub, nil
	}
	if tok == token.MUL {
		if b.Info()&types.IsFloat != 0 {
			return ir.FMul, nil
		}
		return ir.Mul, nil
	}
	if tok == token.QUO {
		if b.Info()&types.IsFloat != 0 {
			return ir.FDiv, nil
		}
		if isUnsigned(t) {
			return ir.UDiv, nil
		}
		return ir.SDiv, nil
	}
	if tok == token.REM {
		if b.Info()&types.IsFloat != 0 {
			return "", fmt.Errorf("remainder not supported for float")
		}
		if isUnsigned(t) {
			return ir.URem, nil
		}
		return ir.SRem, nil
	}
	if tok == token.AND {
		return ir.And, nil
	}
	if tok == token.OR {
		return ir.Or, nil
	}
	if tok == token.LAND {
		return ir.And, nil
	}
	if tok == token.LOR {
		return ir.Or, nil
	}
	return "", fmt.Errorf("unsupported binary op %s", tok.String())
}

func selectICmpPred(tok token.Token, unsigned bool) (ir.ICmpPred, error) {
	if tok == token.EQL {
		return ir.ICmpEq, nil
	}
	if tok == token.NEQ {
		return ir.ICmpNe, nil
	}
	if tok == token.LSS {
		if unsigned {
			return ir.ICmpUlt, nil
		}
		return ir.ICmpSlt, nil
	}
	if tok == token.LEQ {
		if unsigned {
			return ir.ICmpUle, nil
		}
		return ir.ICmpSle, nil
	}
	if tok == token.GTR {
		if unsigned {
			return ir.ICmpUgt, nil
		}
		return ir.ICmpSgt, nil
	}
	if tok == token.GEQ {
		if unsigned {
			return ir.ICmpUge, nil
		}
		return ir.ICmpSge, nil
	}
	return "", fmt.Errorf("unsupported compare op %s", tok.String())
}

func selectFCmpPred(tok token.Token) (ir.FCmpPred, error) {
	if tok == token.EQL {
		return ir.FCmpOeq, nil
	}
	if tok == token.NEQ {
		return ir.FCmpOne, nil
	}
	if tok == token.LSS {
		return ir.FCmpOlt, nil
	}
	if tok == token.LEQ {
		return ir.FCmpOle, nil
	}
	if tok == token.GTR {
		return ir.FCmpOgt, nil
	}
	if tok == token.GEQ {
		return ir.FCmpOge, nil
	}
	return "", fmt.Errorf("unsupported float compare op %s", tok.String())
}

func isUnsigned(t types.Type) bool {
	b, ok := t.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return b.Info()&types.IsUnsigned != 0
}

// mapKeyKind classifies key type for runtime hashing/comparison.
func mapKeyKind(t types.Type) (int, error) {
	u := t.Underlying()
	if b, ok := u.(*types.Basic); ok {
		switch b.Kind() {
		case types.String, types.UntypedString:
			return 2, nil
		case types.Int, types.Int64, types.Int32, types.UntypedInt:
			return 0, nil
		case types.Uint, types.Uint64, types.Uint32, types.Uintptr:
			return 1, nil
		case types.Bool:
			return 1, nil
		}
	}
	if _, ok := u.(*types.Pointer); ok {
		return 3, nil
	}
	if _, ok := u.(*types.Interface); ok {
		return 3, nil
	}
	return 0, fmt.Errorf("unsupported map key type %T", u)
}

// recvBaseName returns a stable string for receiver types for naming methods.
func recvBaseName(t types.Type) string {
	if p, ok := t.(*types.Pointer); ok {
		if named, ok := p.Elem().(*types.Named); ok {
			return named.Obj().Name()
		}
		return "recvptr"
	}
	if named, ok := t.(*types.Named); ok {
		return named.Obj().Name()
	}
	return "recv"
}

func zeroValue(t types.Type) (ir.Value, error) {
	irType, err := goTypeToIR(t)
	if err != nil {
		return nil, err
	}
		if b, ok := t.Underlying().(*types.Basic); ok {
			kind := b.Kind()
			if kind == types.UntypedNil {
				return ir.NewConstant("null", ir.PtrI8), nil
			}
			if kind == types.Bool || kind == types.UntypedBool || kind == types.Int || kind == types.Int32 || kind == types.Int64 || kind == types.Uint32 || kind == types.Uint64 || kind == types.UntypedInt || kind == types.Uint || kind == types.Uintptr || kind == types.UnsafePointer || kind == types.Int8 || kind == types.Uint8 {
				return ir.NewConstant("0", irType), nil
			}
		if kind == types.Float64 || kind == types.UntypedFloat {
			return ir.NewConstant("0.0", irType), nil
		}
		if kind == types.String || kind == types.UntypedString {
			return ir.NewConstant("zeroinitializer", irType), nil
		}
		return nil, fmt.Errorf("zero value for unsupported basic type %v", b.Kind())
	}
	if _, ok := t.Underlying().(*types.Struct); ok {
		return ir.NewConstant("zeroinitializer", irType), nil
	}
	if _, ok := t.Underlying().(*types.Array); ok {
		return ir.NewConstant("zeroinitializer", irType), nil
	}
	if _, ok := t.Underlying().(*types.Slice); ok {
		return ir.NewConstant("zeroinitializer", irType), nil
	}
	if _, ok := t.Underlying().(*types.Map); ok {
		return ir.NewConstant("null", irType), nil
	}
	if _, ok := t.Underlying().(*types.Interface); ok {
		return ir.NewConstant("null", irType), nil
	}
	if _, ok := t.Underlying().(*types.Pointer); ok {
		return ir.NewConstant("null", irType), nil
	}
	if _, ok := t.Underlying().(*types.Chan); ok {
		return ir.NewConstant("null", irType), nil
	}
	if _, ok := t.Underlying().(*types.Signature); ok {
		return ir.NewConstant("null", irType), nil
	}
	return nil, fmt.Errorf("zero value for unsupported type %T", t)
}

func sizeofType(t types.Type) (int64, error) {
	if b, ok := t.Underlying().(*types.Basic); ok {
		kind := b.Kind()
		if kind == types.Bool || kind == types.UntypedBool {
			return 1, nil
		}
		if kind == types.Int8 || kind == types.Uint8 {
			return 1, nil
		}
		if kind == types.Int32 || kind == types.Uint32 {
			return 4, nil
		}
		if kind == types.Int || kind == types.Int64 || kind == types.Uint64 || kind == types.UntypedInt || kind == types.Uint || kind == types.Uintptr || kind == types.UnsafePointer {
			return 8, nil
		}
		if kind == types.Float64 || kind == types.UntypedFloat {
			return 8, nil
		}
		if kind == types.String || kind == types.UntypedString {
			return 16, nil
		}
		return 0, fmt.Errorf("sizeof unsupported basic type %v", b.Kind())
	}
	if _, ok := t.Underlying().(*types.Pointer); ok {
		return 8, nil
	}
	if arr, ok := t.Underlying().(*types.Array); ok {
		elemSize, err := sizeofType(arr.Elem())
		if err != nil {
			return 0, err
		}
		elemAlign := alignOfGoType(arr.Elem())
		stride := alignUp(elemSize, elemAlign)
		return int64(arr.Len()) * stride, nil
	}
	if _, ok := t.Underlying().(*types.Slice); ok {
		return alignUp(24, 8), nil
	}
	if _, ok := t.Underlying().(*types.Map); ok {
		return 8, nil
	}
	if st, ok := t.Underlying().(*types.Struct); ok {
		var off int64
		var maxAlign int64
		for i := 0; i < st.NumFields(); i++ {
			fs, err := sizeofType(st.Field(i).Type())
			if err != nil {
				return 0, err
			}
			align := alignOfGoType(st.Field(i).Type())
			if align > maxAlign {
				maxAlign = align
			}
			off = alignUp(off, align)
			off += fs
		}
		if maxAlign == 0 {
			maxAlign = 1
		}
		off = alignUp(off, maxAlign)
		return off, nil
	}
		return 8, nil
	}

func alignOfGoType(t types.Type) int64 {
	if b, ok := t.Underlying().(*types.Basic); ok {
		kind := b.Kind()
		if kind == types.Bool || kind == types.UntypedBool {
			return 1
		}
		if kind == types.Int8 || kind == types.Uint8 {
			return 1
		}
		if kind == types.Int32 || kind == types.Uint32 {
			return 4
		}
		if kind == types.Int || kind == types.Int64 || kind == types.Uint64 || kind == types.UntypedInt || kind == types.Float64 || kind == types.UntypedFloat || kind == types.String || kind == types.UntypedString || kind == types.Uint || kind == types.Uintptr || kind == types.UnsafePointer {
			return 8
		}
		return 8
	}
	if _, ok := t.Underlying().(*types.Pointer); ok {
		return 8
	}
	if arr, ok := t.Underlying().(*types.Array); ok {
		return alignOfGoType(arr.Elem())
	}
	if _, ok := t.Underlying().(*types.Slice); ok {
		return 8
	}
	if _, ok := t.Underlying().(*types.Map); ok {
		return 8
	}
	if st, ok := t.Underlying().(*types.Struct); ok {
		var max int64 = 1
		for i := 0; i < st.NumFields(); i++ {
			a := alignOfGoType(st.Field(i).Type())
			if a > max {
				max = a
			}
		}
		return max
	}
	return 1
}

func alignUp(x int64, align int64) int64 {
	if align <= 1 {
		return x
	}
	r := x % align
	if r == 0 {
		return x
	}
	return x + (align - r)
}

func pointeeType(t ir.Type) (ir.Type, error) {
	if p, ok := t.(ir.PointerType); ok {
		return p.Elem, nil
	}
	if p, ok := t.(ir.BasicType); ok {
		if s := p.String(); len(s) > 0 && s[len(s)-1] == '*' {
			return ir.BasicType(s[:len(s)-1]), nil
		}
	}
	return nil, fmt.Errorf("type %s is not a pointer", t.String())
}
