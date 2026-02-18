package backend

import (
	"gominic/ir"
)

func EmitModule(mod *ir.Module) string {
	var parts []string
	var msvc bool
	triple := ir.GetModuleTargetTriple(mod)
	if triple == "" {
		triple = "x86_64-pc-linux-gnu"
	}
	// Проверка на Windows/MSVC для использования sret в объявлениях функций
	substr := "windows-msvc"
	msvcVal := 0 == 1
	if len(substr) <= len(triple) {
		for iCheck := 0; iCheck <= len(triple)-len(substr); iCheck = iCheck + 1 {
			match := 1 == 1
			for jCheck := 0; jCheck < len(substr); jCheck = jCheck + 1 {
				if triple[iCheck+jCheck] != substr[jCheck] {
					match = 0 == 1
					break
				}
			}
			if match {
				msvcVal = 1 == 1
				break
			}
		}
	}
	msvc = msvcVal
	dl := ir.GetModuleDataLayout(mod)
	if dl == "" {
		dl = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
	}
	parts = strBufWriteString(parts, "target triple = \""+triple+"\"\n")
	parts = strBufWriteString(parts, "target datalayout = \""+dl+"\"\n\n")
	globals := ir.GetModuleGlobals(mod)
	parts = strBufWriteString(parts, "declare void @gominic_memcpy(i8*, i8*, i64)\n")
	parts = strBufWriteString(parts, "declare void @gominic_abort()\n")
	parts = strBufWriteString(parts, "declare i8* @gominic_makeSlice(i64, i64, i64)\n")
	parts = strBufWriteString(parts, "declare void @gominic_print(i8*, i64)\n")
	parts = strBufWriteString(parts, "declare void @gominic_printInt(i64)\n")
	parts = strBufWriteString(parts, "declare void @gominic_println()\n\n")
	parts = strBufWriteString(parts, "declare i1 @gominic_str_eq(i8*, i64, i8*, i64)\n")
	parts = strBufWriteString(parts, "declare void @gominic_str_concat({ i8*, i64 }*, { i8*, i64 }*, { i8*, i64 }*)\n\n")
	parts = strBufWriteString(parts, "declare i8* @gominic_map_new(i64, i64, i32)\n")
	parts = strBufWriteString(parts, "declare void @gominic_map_set(i8*, i8*, i8*)\n")
	parts = strBufWriteString(parts, "declare i1 @gominic_map_get(i8*, i8*, i8*)\n")
	parts = strBufWriteString(parts, "declare i64 @gominic_map_len(i8*)\n\n")
	parts = strBufWriteString(parts, "declare i64 @gominic_argc()\n")
	if msvc {
		parts = strBufWriteString(parts, "declare void @gominic_argv({ i8*, i64 }* sret({ i8*, i64 }), i64)\n")
	} else {
		parts = strBufWriteString(parts, "declare { i8*, i64 } @gominic_argv(i64)\n")
	}
	parts = strBufWriteString(parts, "declare i1 @gominic_write_file(i8*, i64, i8*, i64)\n\n")
	parts = strBufWriteString(parts, "declare void @flag.Usage()\n")
	parts = strBufWriteString(parts, "declare void @Usage()\n")
	parts = strBufWriteString(parts, "declare void @flag.Parse()\n")
	parts = strBufWriteString(parts, "declare { { i8*, i64 }*, i64, i64 } @flag.Args()\n")
	parts = strBufWriteString(parts, "declare void @flag.BoolVar(i1*, { i8*, i64 }, i1, { i8*, i64 })\n")
	parts = strBufWriteString(parts, "declare void @flag.StringVar({ i8*, i64 }*, { i8*, i64 }, { i8*, i64 }, { i8*, i64 })\n")
	parts = strBufWriteString(parts, "declare { i8*, i64 } @error.Error(i8*)\n")
	parts = strBufWriteString(parts, "declare i8* @errors.New({ i8*, i64 })\n")
	parts = strBufWriteString(parts, "declare { { i8*, i64, i64 }, i8* } @exec.Cmd.CombinedOutput(i8*)\n")
	parts = strBufWriteString(parts, "declare { { i8*, i64 }, { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }, { i8*, i64 }, i8*, i8*, i8*, { i8***, i64, i64 }, i8**, i8**, i8**, i8*, i8*, i8*, i8*, { i8**, i64, i64 }, { i8**, i64, i64 }, { i8**, i64, i64 }, i8*, i8*, { i8*, i64, i64 }, i8*, { { i8*, i64 }, { i8*, i64 } } }* @exec.Command({ i8*, i64 }, { { i8*, i64 }*, i64, i64 })\n")
	parts = strBufWriteString(parts, "declare { { { i8**, { i8***, i64, i64 }, i8**, i8**, { i8*, i64 } }**, i64, i64 }, i8* } @frontend.ParseAndCheckAll({ { i8*, i64 }*, i64, i64 })\n")
	parts = strBufWriteString(parts, "declare { i8**, i8* } @frontend.BuildModule({ { i8**, { i8***, i64, i64 }, i8**, i8**, { i8*, i64 } }**, i64, i64 })\n")
	parts = strBufWriteString(parts, "declare i8* @os.WriteFile({ i8*, i64 }, { i8*, i64, i64 }, i64)\n")
	parts = strBufWriteString(parts, "declare void @os.File.WriteString(i8**, { i8*, i64 })\n")
	parts = strBufWriteString(parts, "declare void @os.Exit(i64)\n")
	// Методы ir.* определяются при компиляции ir/ir.go в LLVM IR.
	// Не объявляем их здесь, чтобы избежать конфликта (redefinition)
	parts = strBufWriteString(parts, "declare { i8*, i64 } @types.Package.Name(i8**)\n")
	parts = strBufWriteString(parts, "@os.Stderr = external global i8*\n")
	parts = strBufWriteString(parts, "@os.Stdout = external global i8*\n")
	parts = strBufWriteString(parts, "@os.Stdin = external global i8*\n")
	parts = strBufWriteString(parts, "@gominic.I64 = external global i8***\n")
	parts = strBufWriteString(parts, "@gominic.I32 = external global i8***\n")
	parts = strBufWriteString(parts, "@gominic.I8 = external global i8***\n")
	parts = strBufWriteString(parts, "@gominic.I1 = external global i8***\n")
	parts = strBufWriteString(parts, "@gominic.F64 = external global i8***\n")
	parts = strBufWriteString(parts, "@gominic.Void = external global i8***\n")
	parts = strBufWriteString(parts, "@gominic.PtrI8 = external global i8***\n")
	parts = strBufWriteString(parts, "@gominic.false = external global i1\n")
	parts = strBufWriteString(parts, "@gominic.true = external global i1\n")
	parts = strBufWriteString(parts, "@gominic.nil = external global i8**\n")
	// Глобальные переменные для методов ir.* (как receiver).
	parts = strBufWriteString(parts, "@ir.bb = external global i8***\n")
	parts = strBufWriteString(parts, "@ir.m = external global i8***\n")
	parts = strBufWriteString(parts, "@ir.t = external global i8***\n")
	parts = strBufWriteString(parts, "@ir.v = external global i8***\n")
	parts = strBufWriteString(parts, "@ir.fn = external global i8***\n")
	parts = strBufWriteString(parts, "@ir.nil = external global i8***\n")
	// Константы булевых значений для пакета ir.
	parts = strBufWriteString(parts, "@ir.false = external global i1\n")
	parts = strBufWriteString(parts, "@ir.true = external global i1\n\n")
	// Те же символы без префикса для совместимости.
	parts = strBufWriteString(parts, "@Stderr = external global i8*\n")
	parts = strBufWriteString(parts, "@Stdout = external global i8*\n")
	parts = strBufWriteString(parts, "@Stdin = external global i8*\n\n")
	for iGlobal := 0; iGlobal < len(globals); iGlobal = iGlobal + 1 {
		g := globals[iGlobal]
		visibility := ""
		if ir.GetGlobalPrivate(&g) {
			visibility = "private unnamed_addr "
		}
		alignGlobal := ""
		alignValGlobal := ir.GetGlobalAlign(&g)
		if alignValGlobal > 0 {
			// Встроенная логика FormatInt64.
			alignStrGlobal := ""
			if alignValGlobal == 0 {
				alignStrGlobal = "0"
			} else {
				digitStrsGlobal1 := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
				var digitsGlobal1 []string
				negativeGlobal1 := 0 == 1
				nGlobal1 := alignValGlobal
				if nGlobal1 < 0 {
					negativeGlobal1 = 1 == 1
					nGlobal1 = -nGlobal1
				}
				for nGlobal1 > 0 {
					digitGlobal1 := nGlobal1 % 10
					digitsGlobal1 = append(digitsGlobal1, digitStrsGlobal1[digitGlobal1])
					nGlobal1 = nGlobal1 / 10
				}
				alignStrGlobal = ""
				if negativeGlobal1 {
					alignStrGlobal = "-"
				}
				for iDigitGlobal1 := len(digitsGlobal1) - 1; iDigitGlobal1 >= 0; iDigitGlobal1 = iDigitGlobal1 - 1 {
					alignStrGlobal = alignStrGlobal + digitsGlobal1[iDigitGlobal1]
				}
			}
			alignGlobal = ", align " + alignStrGlobal
		}
		parts = strBufWriteString(parts, "@"+ir.GetGlobalName(&g)+" = "+visibility+"constant "+llvmType(ir.GetGlobalType(&g))+" "+ir.GetGlobalValue(&g)+alignGlobal+"\n")
	}
	if len(globals) > 0 {
		parts = strBufWriteString(parts, "\n")
	}
	modFunctions := ir.GetModuleFunctions(mod)
	for iFunc := 0; iFunc < len(modFunctions); iFunc = iFunc + 1 {
		fn := modFunctions[iFunc]
		// Встроенная логика emitFunction.
		var params []string
		fnParams := ir.GetFunctionParams(fn)
		for iParam := 0; iParam < len(fnParams); iParam = iParam + 1 {
			p := fnParams[iParam]
			var paramStr string
			var nilTypeDesc3 *ir.TypeDesc
			if iParam == 0 && ir.GetFunctionSretType(fn) != nilTypeDesc3 {
				paramStr = "sret(" + llvmType(ir.GetFunctionSretType(fn)) + ") " + llvmType(valueType(p)) + " %" + valueName(p)
			} else if valueByVal(p) {
				paramStr = llvmType(valueType(p)) + " byval(" + llvmType(valueByValType(p)) + ") %" + valueName(p)
			} else {
				paramStr = llvmType(valueType(p)) + " %" + valueName(p)
			}
			params = append(params, paramStr)
		}

		var result string
		var nilTypeDesc2 *ir.TypeDesc
		if ir.GetFunctionSretType(fn) != nilTypeDesc2 {
			result = llvmType(ir.GetVoidType())
		} else {
			fnResults := ir.GetFunctionResults(fn)
			resCount := len(fnResults)
			if resCount == 0 {
				result = llvmType(ir.GetVoidType())
			} else if resCount == 1 {
				result = llvmType(fnResults[0])
			} else {
				// TODO: поддержать несколько возвращаемых значений через struct или sret lowering.
				result = llvmType(ir.GetVoidType())
			}
		}

		parts = strBufWriteString(parts, "define "+result+" @"+ir.GetFunctionName(fn)+"("+joinStrings(params, ", ")+") {\n")
		fnBlocks := ir.GetFunctionBlocks(fn)
		for iBlock := 0; iBlock < len(fnBlocks); iBlock = iBlock + 1 {
			bb := fnBlocks[iBlock]
			parts = strBufWriteString(parts, ir.GetBasicBlockName(bb)+":\n")
			bbInstrs := ir.GetBasicBlockInstrs(bb)
			for iInstr := 0; iInstr < len(bbInstrs); iInstr = iInstr + 1 {
				parts = strBufWriteString(parts, "  "+renderInstr(bbInstrs[iInstr])+"\n")
			}
			terminator := ir.GetBasicBlockTerminator(bb)
			var nilInstr1 *ir.Instruction
			if terminator != nilInstr1 {
				parts = strBufWriteString(parts, "  "+renderInstr(*terminator)+"\n")
			} else {
				// запасной путь, чтобы не получить невалидный IR.
				parts = strBufWriteString(parts, "  unreachable\n")
			}
		}
		parts = strBufWriteString(parts, "}\n\n")
	}
	return strBufString(parts)
}

type strBuf struct {
	parts []string
}

// Простые вспомогательные функции вместо методов для самоприменимости.
func strBufWriteString(parts []string, s string) []string {
	return append(parts, s)
}

func strBufString(parts []string) string {
	var result string
	for iPart := 0; iPart < len(parts); iPart = iPart + 1 {
		result = result + parts[iPart]
	}
	return result
}

type emitter struct {
	buf strBuf
}

func emitGlobals(parts []string, globals []ir.Global, msvc bool) []string {
	// объявления runtime-функций
	parts = strBufWriteString(parts, "declare void @gominic_memcpy(i8*, i8*, i64)\n")
	parts = strBufWriteString(parts, "declare void @gominic_abort()\n")
	parts = strBufWriteString(parts, "declare i8* @gominic_makeSlice(i64, i64, i64)\n")
	parts = strBufWriteString(parts, "declare void @gominic_print(i8*, i64)\n")
	parts = strBufWriteString(parts, "declare void @gominic_printInt(i64)\n")
	parts = strBufWriteString(parts, "declare void @gominic_println()\n\n")
	parts = strBufWriteString(parts, "declare i1 @gominic_str_eq(i8*, i64, i8*, i64)\n")
	parts = strBufWriteString(parts, "declare void @gominic_str_concat({ i8*, i64 }*, { i8*, i64 }*, { i8*, i64 }*)\n\n")
	parts = strBufWriteString(parts, "declare i8* @gominic_map_new(i64, i64, i32)\n")
	parts = strBufWriteString(parts, "declare void @gominic_map_set(i8*, i8*, i8*)\n")
	parts = strBufWriteString(parts, "declare i1 @gominic_map_get(i8*, i8*, i8*)\n")
	parts = strBufWriteString(parts, "declare i64 @gominic_map_len(i8*)\n\n")
	parts = strBufWriteString(parts, "declare i64 @gominic_argc()\n")
	if msvc {
		parts = strBufWriteString(parts, "declare void @gominic_argv({ i8*, i64 }* sret({ i8*, i64 }), i64)\n")
	} else {
		parts = strBufWriteString(parts, "declare { i8*, i64 } @gominic_argv(i64)\n")
	}
	parts = strBufWriteString(parts, "declare i1 @gominic_write_file(i8*, i64, i8*, i64)\n\n")
	for iGlobal := 0; iGlobal < len(globals); iGlobal = iGlobal + 1 {
		g := globals[iGlobal]
		visibility := ""
		if ir.GetGlobalPrivate(&g) {
			visibility = "private unnamed_addr "
		}
		alignGlobal := ""
		alignValGlobal := ir.GetGlobalAlign(&g)
		if alignValGlobal > 0 {
			// Встроенная логика FormatInt64.
			alignStrGlobal := ""
			if alignValGlobal == 0 {
				alignStrGlobal = "0"
			} else {
				digitStrsGlobal2 := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
				var digitsGlobal2 []string
				negativeGlobal2 := 0 == 1
				nGlobal2 := alignValGlobal
				if nGlobal2 < 0 {
					negativeGlobal2 = 1 == 1
					nGlobal2 = -nGlobal2
				}
				for nGlobal2 > 0 {
					digitGlobal2 := nGlobal2 % 10
					digitsGlobal2 = append(digitsGlobal2, digitStrsGlobal2[digitGlobal2])
					nGlobal2 = nGlobal2 / 10
				}
				alignStrGlobal = ""
				if negativeGlobal2 {
					alignStrGlobal = "-"
				}
				for iDigitGlobal2 := len(digitsGlobal2) - 1; iDigitGlobal2 >= 0; iDigitGlobal2 = iDigitGlobal2 - 1 {
					alignStrGlobal = alignStrGlobal + digitsGlobal2[iDigitGlobal2]
				}
			}
			alignGlobal = ", align " + alignStrGlobal
		}
		parts = strBufWriteString(parts, "@"+ir.GetGlobalName(&g)+" = "+visibility+"constant "+llvmType(ir.GetGlobalType(&g))+" "+ir.GetGlobalValue(&g)+alignGlobal+"\n")
	}
	if len(globals) > 0 {
		parts = strBufWriteString(parts, "\n")
	}
	return parts
}

func emitFunction(parts []string, fn *ir.Function) []string {
	var params []string
	fnParams := ir.GetFunctionParams(fn)
	for iParam := 0; iParam < len(fnParams); iParam = iParam + 1 {
		p := fnParams[iParam]
		var paramStr string
		var nilTypeDesc1 *ir.TypeDesc
		if iParam == 0 && ir.GetFunctionSretType(fn) != nilTypeDesc1 {
			paramStr = "sret(" + llvmType(ir.GetFunctionSretType(fn)) + ") " + llvmType(valueType(p)) + " %" + valueName(p)
		} else if valueByVal(p) {
			paramStr = llvmType(valueType(p)) + " byval(" + llvmType(valueByValType(p)) + ") %" + valueName(p)
		} else {
			paramStr = llvmType(valueType(p)) + " %" + valueName(p)
		}
		params = append(params, paramStr)
	}

	var result string
	var nilTypeDesc4 *ir.TypeDesc
	if ir.GetFunctionSretType(fn) != nilTypeDesc4 {
		result = llvmType(ir.GetVoidType())
	} else {
		fnResults := ir.GetFunctionResults(fn)
		resCount := len(fnResults)
		if resCount == 0 {
			result = llvmType(ir.GetVoidType())
		} else if resCount == 1 {
			result = llvmType(fnResults[0])
		} else {
			// TODO: поддержать несколько возвращаемых значений через struct или sret lowering.
			result = llvmType(ir.GetVoidType())
		}
	}

	parts = strBufWriteString(parts, "define "+result+" @"+ir.GetFunctionName(fn)+"("+joinStrings(params, ", ")+") {\n")
	fnBlocks := ir.GetFunctionBlocks(fn)
	for iBlock := 0; iBlock < len(fnBlocks); iBlock = iBlock + 1 {
		bb := fnBlocks[iBlock]
		parts = strBufWriteString(parts, ir.GetBasicBlockName(bb)+":\n")
		bbInstrs := ir.GetBasicBlockInstrs(bb)
		for iInstr := 0; iInstr < len(bbInstrs); iInstr = iInstr + 1 {
			parts = strBufWriteString(parts, "  "+renderInstr(bbInstrs[iInstr])+"\n")
		}
		terminator := ir.GetBasicBlockTerminator(bb)
		var nilInstr1 *ir.Instruction
		if terminator != nilInstr1 {
			parts = strBufWriteString(parts, "  "+renderInstr(*terminator)+"\n")
		} else {
			// запасной путь, чтобы не получить невалидный IR.
			parts = strBufWriteString(parts, "  unreachable\n")
		}
	}
	parts = strBufWriteString(parts, "}\n\n")
	return parts
}

func EmitBasicBlock(parts []string, bb *ir.BasicBlock) []string {
	parts = strBufWriteString(parts, ir.GetBasicBlockName(bb)+":\n")
	bbInstrs := ir.GetBasicBlockInstrs(bb)
	for iInstr := 0; iInstr < len(bbInstrs); iInstr = iInstr + 1 {
		parts = strBufWriteString(parts, "  "+renderInstr(bbInstrs[iInstr])+"\n")
	}
	terminator := ir.GetBasicBlockTerminator(bb)
	var nilInstr2 *ir.Instruction
	if terminator != nilInstr2 {
		parts = strBufWriteString(parts, "  "+renderInstr(*terminator)+"\n")
	} else {
		// запасной путь, чтобы не получить невалидный IR.
		parts = strBufWriteString(parts, "  unreachable\n")
	}
	return parts
}

func llvmType(t *ir.TypeDesc) string {
	var nilTypeDesc5 *ir.TypeDesc
	if t == nilTypeDesc5 {
		return "void"
	}
	return typeDescString(t)
}

func typeDescString(t *ir.TypeDesc) string {
	var nilTypeDesc6 *ir.TypeDesc
	if t == nilTypeDesc6 {
		return "void"
	}
	k := ir.GetTypeKind(t)
	if ir.IsKindBasic(k) {
		return ir.GetTypeDescBasic(t)
	}
	if ir.IsKindPointer(k) {
		return typeDescString(ir.GetTypeDescElem(t)) + "*"
	}
	if ir.IsKindStruct(k) {
		var parts []string
		fields := ir.GetTypeDescFields(t)
		for iField := 0; iField < len(fields); iField = iField + 1 {
			parts = append(parts, typeDescString(fields[iField]))
		}
		return "{ " + joinStrings(parts, ", ") + " }"
	}
	if ir.IsKindArray(k) {
		lenVal := int64(ir.GetTypeDescLen(t))
		lenStr := ""
		if lenVal == 0 {
			lenStr = "0"
		} else {
			digitStrsArray := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
			var digitsArray []string
			negativeArray := 0 == 1
			nArray := lenVal
			if nArray < 0 {
				negativeArray = 1 == 1
				nArray = -nArray
			}
			for nArray > 0 {
				digitArray := nArray % 10
				digitsArray = append(digitsArray, digitStrsArray[digitArray])
				nArray = nArray / 10
			}
			lenStr = ""
			if negativeArray {
				lenStr = "-"
			}
			for iDigitArray := len(digitsArray) - 1; iDigitArray >= 0; iDigitArray = iDigitArray - 1 {
				lenStr = lenStr + digitsArray[iDigitArray]
			}
		}
		return "[" + lenStr + " x " + typeDescString(ir.GetTypeDescElem(t)) + "]"
	}
	if ir.IsKindSlice(k) {
		return "{ " + typeDescString(ir.PtrTo(ir.GetTypeDescElem(t))) + ", i64, i64 }"
	}
	if ir.IsKindString(k) {
		return "{ i8*, i64 }"
	}
	return "void"
}

func renderInstr(inst ir.Instruction) string {
	k := ir.GetInstrKind(&inst)
	if ir.IsInstrBinOp(k) {
		return "%" + valueName(ir.GetInstrDest(&inst)) + " = " + string(ir.GetInstrBinOp(&inst)) + " " + llvmType(valueType(ir.GetInstrX(&inst))) + " " + valStr(ir.GetInstrX(&inst)) + ", " + valStr(ir.GetInstrY(&inst))
	}
	if ir.IsInstrReturn(k) {
		retVals := ir.GetInstrRetVals(&inst)
		if len(retVals) == 0 {
			return "ret void"
		}
		if len(retVals) == 1 {
			v := retVals[0]
			return "ret " + llvmType(valueType(v)) + " " + valStr(v)
		}
		return "ret void ; TODO multiple return values"
	}
	if ir.IsInstrCall(k) {
		var callArgs []string
		callArgsList := ir.GetInstrCallArgs(&inst)
		for iCallArg := 0; iCallArg < len(callArgsList); iCallArg = iCallArg + 1 {
			callArgStr := llvmType(valueType(callArgsList[iCallArg])) + " " + valStr(callArgsList[iCallArg])
			callArgs = append(callArgs, callArgStr)
		}
		ret := "void"
		var nilTypeDesc7 *ir.TypeDesc
		if ir.GetInstrCallRet(&inst) != nilTypeDesc7 {
			ret = llvmType(ir.GetInstrCallRet(&inst))
		}
		prefix := ""
		var nilValue1 *ir.Value
		if ir.GetInstrDest(&inst) != nilValue1 && valueName(ir.GetInstrDest(&inst)) != "" && ret != "void" {
			prefix = "%" + valueName(ir.GetInstrDest(&inst)) + " = "
		}
		return prefix + "call " + ret + " @" + ir.GetInstrCallName(&inst) + "(" + joinStrings(callArgs, ", ") + ")"
	}
	if ir.IsInstrConv(k) {
		return "%" + valueName(ir.GetInstrDest(&inst)) + " = " + string(ir.GetInstrConvOp(&inst)) + " " + llvmType(valueType(ir.GetInstrConvSrc(&inst))) + " " + valStr(ir.GetInstrConvSrc(&inst)) + " to " + llvmType(ir.GetInstrConvTo(&inst))
	}
	if ir.IsInstrAlloca(k) {
		alignAlloca := ""
		alignValAlloca := ir.GetInstrAllocaAlign(&inst)
		if alignValAlloca > 0 {
			// Встроенная логика FormatInt64.
			alignStrAlloca := ""
			if alignValAlloca == 0 {
				alignStrAlloca = "0"
			} else {
				digitStrsAlloca := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
				var digitsAlloca []string
				negativeAlloca := 0 == 1
				nAlloca := alignValAlloca
				if nAlloca < 0 {
					negativeAlloca = 1 == 1
					nAlloca = -nAlloca
				}
				for nAlloca > 0 {
					digitAlloca := nAlloca % 10
					digitsAlloca = append(digitsAlloca, digitStrsAlloca[digitAlloca])
					nAlloca = nAlloca / 10
				}
				alignStrAlloca = ""
				if negativeAlloca {
					alignStrAlloca = "-"
				}
				for iDigitAlloca := len(digitsAlloca) - 1; iDigitAlloca >= 0; iDigitAlloca = iDigitAlloca - 1 {
					alignStrAlloca = alignStrAlloca + digitsAlloca[iDigitAlloca]
				}
			}
			alignAlloca = " , align " + alignStrAlloca
		}
		return "%" + valueName(ir.GetInstrDest(&inst)) + " = alloca " + llvmType(ir.GetInstrAllocaType(&inst)) + alignAlloca
	}
	if ir.IsInstrLoad(k) {
		alignLoad := ""
		alignValLoad := ir.GetInstrLoadAlign(&inst)
		if alignValLoad > 0 {
			// Встроенная логика FormatInt64.
			alignStrLoad := ""
			if alignValLoad == 0 {
				alignStrLoad = "0"
			} else {
				digitStrsLoad := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
				var digitsLoad []string
				negativeLoad := 0 == 1
				nLoad := alignValLoad
				if nLoad < 0 {
					negativeLoad = 1 == 1
					nLoad = -nLoad
				}
				for nLoad > 0 {
					digitLoad := nLoad % 10
					digitsLoad = append(digitsLoad, digitStrsLoad[digitLoad])
					nLoad = nLoad / 10
				}
				alignStrLoad = ""
				if negativeLoad {
					alignStrLoad = "-"
				}
				for iDigitLoad := len(digitsLoad) - 1; iDigitLoad >= 0; iDigitLoad = iDigitLoad - 1 {
					alignStrLoad = alignStrLoad + digitsLoad[iDigitLoad]
				}
			}
			alignLoad = ", align " + alignStrLoad
		}
		return "%" + valueName(ir.GetInstrDest(&inst)) + " = load " + llvmType(valueType(ir.GetInstrDest(&inst))) + ", " + llvmType(valueType(ir.GetInstrLoadSrc(&inst))) + " " + valStr(ir.GetInstrLoadSrc(&inst)) + alignLoad
	}
	if ir.IsInstrStore(k) {
		srcType := valueType(ir.GetInstrStoreSrc(&inst))
		// Пропускаем store для void (это невалидно в LLVM IR).
		var nilTypeDescStore *ir.TypeDesc
		if srcType == nilTypeDescStore {
			return "; store void skipped (invalid in LLVM IR)"
		}
		// Проверяем, что тип действительно void.
		if ir.IsKindBasic(ir.GetTypeKind(srcType)) && ir.GetTypeDescBasic(srcType) == "void" {
			return "; store void skipped (invalid in LLVM IR)"
		}
		alignStore := ""
		alignValStore := ir.GetInstrStoreAlign(&inst)
		if alignValStore > 0 {
			// Встроенная логика FormatInt64.
			alignStrStore := ""
			if alignValStore == 0 {
				alignStrStore = "0"
			} else {
				digitStrsStore := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
				var digitsStore []string
				negativeStore := 0 == 1
				nStore := alignValStore
				if nStore < 0 {
					negativeStore = 1 == 1
					nStore = -nStore
				}
				for nStore > 0 {
					digitStore := nStore % 10
					digitsStore = append(digitsStore, digitStrsStore[digitStore])
					nStore = nStore / 10
				}
				alignStrStore = ""
				if negativeStore {
					alignStrStore = "-"
				}
				for iDigitStore := len(digitsStore) - 1; iDigitStore >= 0; iDigitStore = iDigitStore - 1 {
					alignStrStore = alignStrStore + digitsStore[iDigitStore]
				}
			}
			alignStore = ", align " + alignStrStore
		}
		return "store " + llvmType(srcType) + " " + valStr(ir.GetInstrStoreSrc(&inst)) + ", " + llvmType(valueType(ir.GetInstrStoreDst(&inst))) + " " + valStr(ir.GetInstrStoreDst(&inst)) + alignStore
	}
	if ir.IsInstrGEP(k) {
		var gepParts []string
		gepIndices := ir.GetInstrGepIndices(&inst)
		for iGep := 0; iGep < len(gepIndices); iGep = iGep + 1 {
			partStr := llvmType(valueType(gepIndices[iGep])) + " " + valStr(gepIndices[iGep])
			gepParts = append(gepParts, partStr)
		}
		return "%" + valueName(ir.GetInstrDest(&inst)) + " = getelementptr inbounds " + llvmType(ir.GetInstrGepPointee(&inst)) + ", " + llvmType(valueType(ir.GetInstrGepSrc(&inst))) + " " + valStr(ir.GetInstrGepSrc(&inst)) + ", " + joinStrings(gepParts, ", ")
	}
	if ir.IsInstrICmp(k) {
		return "%" + valueName(ir.GetInstrDest(&inst)) + " = icmp " + string(ir.GetInstrICmpPred(&inst)) + " " + llvmType(valueType(ir.GetInstrICmpX(&inst))) + " " + valStr(ir.GetInstrICmpX(&inst)) + ", " + valStr(ir.GetInstrICmpY(&inst))
	}
	if ir.IsInstrFCmp(k) {
		return "%" + valueName(ir.GetInstrDest(&inst)) + " = fcmp " + string(ir.GetInstrFCmpPred(&inst)) + " " + llvmType(valueType(ir.GetInstrFCmpX(&inst))) + " " + valStr(ir.GetInstrFCmpX(&inst)) + ", " + valStr(ir.GetInstrFCmpY(&inst))
	}
	if ir.IsInstrBr(k) {
		return "br label %" + ir.GetInstrBrTarget(&inst)
	}
	if ir.IsInstrCondBr(k) {
		return "br i1 " + valStr(ir.GetInstrCondCond(&inst)) + ", label %" + ir.GetInstrCondTrue(&inst) + ", label %" + ir.GetInstrCondFalse(&inst)
	}
	if ir.IsInstrCallVoid(k) {
		var callVoidArgs []string
		callVoidArgsList := ir.GetInstrCallArgs(&inst)
		for iCallVoidArg := 0; iCallVoidArg < len(callVoidArgsList); iCallVoidArg = iCallVoidArg + 1 {
			callVoidArgStr := llvmType(valueType(callVoidArgsList[iCallVoidArg])) + " " + valStr(callVoidArgsList[iCallVoidArg])
			callVoidArgs = append(callVoidArgs, callVoidArgStr)
		}
		return "call void @" + ir.GetInstrCallName(&inst) + "(" + joinStrings(callVoidArgs, ", ") + ")"
	}
	if ir.IsInstrBitcast(k) {
		return "%" + valueName(ir.GetInstrDest(&inst)) + " = bitcast " + llvmType(valueType(ir.GetInstrBitcastSrc(&inst))) + " " + valStr(ir.GetInstrBitcastSrc(&inst)) + " to " + llvmType(ir.GetInstrBitcastTarget(&inst))
	}
	if ir.IsInstrMemcpy(k) {
		return "call void @gominic_memcpy(" + llvmType(valueType(ir.GetInstrMemcpyDest(&inst))) + " " + valStr(ir.GetInstrMemcpyDest(&inst)) + ", " + llvmType(valueType(ir.GetInstrMemcpySrc(&inst))) + " " + valStr(ir.GetInstrMemcpySrc(&inst)) + ", " + llvmType(valueType(ir.GetInstrMemcpySize(&inst))) + " " + valStr(ir.GetInstrMemcpySize(&inst)) + ")"
	}
	return "unreachable"
}

func valStr(v *ir.Value) string {
	var nilValue2 *ir.Value
	if v == nilValue2 {
		return ""
	}
	k := valueKind(v)
	if ir.IsValueRegister(k) || ir.IsValueParam(k) {
		return "%" + valueName(v)
	}
	if ir.IsValueConstant(k) {
		vt := valueType(v)
		var nilTypeDesc8 *ir.TypeDesc
		if vt != nilTypeDesc8 && ir.IsKindBasic(ir.GetTypeKind(vt)) && ir.GetTypeDescBasic(vt) == "double" {
			raw := valueName(v)
			// Встроенная проверка containsAny.
			chars := ".eE"
			hasAny := 0 == 1
			for iAny := 0; iAny < len(raw); iAny = iAny + 1 {
				for jAny := 0; jAny < len(chars); jAny = jAny + 1 {
					if raw[iAny] == chars[jAny] {
						hasAny = 1 == 1
						break
					}
				}
				if hasAny {
					break
				}
			}
			if !hasAny {
				return raw + ".0"
			}
			return raw
		}
		// Для целых переводим шестнадцатеричные константы (0x...) в десятичные.
		raw := valueName(v)
		if len(raw) > 2 && raw[0] == 48 && (raw[1] == 120 || raw[1] == 88) {
			// Разбираем hex и переводим в десятичный вид.
			hexStr := raw[2:]
			val := int64(0)
			for iHex := 0; iHex < len(hexStr); iHex = iHex + 1 {
				ch := hexStr[iHex]
				chVal := int64(0)
				if ch == 48 {
					chVal = 0
				} else if ch == 49 {
					chVal = 1
				} else if ch == 50 {
					chVal = 2
				} else if ch == 51 {
					chVal = 3
				} else if ch == 52 {
					chVal = 4
				} else if ch == 53 {
					chVal = 5
				} else if ch == 54 {
					chVal = 6
				} else if ch == 55 {
					chVal = 7
				} else if ch == 56 {
					chVal = 8
				} else if ch == 57 {
					chVal = 9
				} else if ch == 65 || ch == 97 {
					chVal = 10
				} else if ch == 66 || ch == 98 {
					chVal = 11
				} else if ch == 67 || ch == 99 {
					chVal = 12
				} else if ch == 68 || ch == 100 {
					chVal = 13
				} else if ch == 69 || ch == 101 {
					chVal = 14
				} else if ch == 70 || ch == 102 {
					chVal = 15
				}
				val = val*16 + chVal
			}
			// Встроенная логика FormatInt64.
			if val == 0 {
				return "0"
			}
			digitStrsVal := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
			var digitsVal []string
			negativeVal := 0 == 1
			nVal := val
			if nVal < 0 {
				negativeVal = 1 == 1
				nVal = -nVal
			}
			for nVal > 0 {
				digitVal := nVal % 10
				digitsVal = append(digitsVal, digitStrsVal[digitVal])
				nVal = nVal / 10
			}
			result := ""
			if negativeVal {
				result = "-"
			}
			for iDigitVal := len(digitsVal) - 1; iDigitVal >= 0; iDigitVal = iDigitVal - 1 {
				result = result + digitsVal[iDigitVal]
			}
			return result
		}
		return raw
	}
	return valueName(v)
}

func valueName(v *ir.Value) string {
	return ir.ValueName(v)
}

func valueType(v *ir.Value) *ir.TypeDesc {
	return ir.ValueType(v)
}

func valueKind(v *ir.Value) ir.ValueKind {
	return ir.GetValueKind(v)
}

func valueByVal(v *ir.Value) bool {
	return ir.ValueByVal(v)
}

func valueByValType(v *ir.Value) *ir.TypeDesc {
	return ir.ValueByValType(v)
}

func emitTargetHeader(parts []string, mod *ir.Module, msvc *bool) []string {
	triple := ir.GetModuleTargetTriple(mod)
	if triple == "" {
		triple = "x86_64-pc-linux-gnu"
	}
	// Встроенная проверка contains без отдельного вызова функции.
	substr := "windows-msvc"
	msvcVal := 0 == 1
	if len(substr) <= len(triple) {
		for iCheck := 0; iCheck <= len(triple)-len(substr); iCheck = iCheck + 1 {
			match := 1 == 1
			for jCheck := 0; jCheck < len(substr); jCheck = jCheck + 1 {
				if triple[iCheck+jCheck] != substr[jCheck] {
					match = 0 == 1
					break
				}
			}
			if match {
				msvcVal = 1 == 1
				break
			}
		}
	}
	*msvc = msvcVal
	dl := ir.GetModuleDataLayout(mod)
	if dl == "" {
		dl = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
	}
	parts = strBufWriteString(parts, "target triple = \""+triple+"\"\n")
	parts = strBufWriteString(parts, "target datalayout = \""+dl+"\"\n\n")
	return parts
}

// joinStrings: замена strings.Join для самоприменимости.
func joinStrings(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	result := parts[0]
	for iJoin := 1; iJoin < len(parts); iJoin = iJoin + 1 {
		result = result + sep + parts[iJoin]
	}
	return result
}

func FormatInt64(n int64) string {
	if n == 0 {
		return "0"
	}
	digitStrsFormat := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var digitsFormat []string
	negativeFormat := 0 == 1
	nFormat := n
	if nFormat < 0 {
		negativeFormat = 1 == 1
		nFormat = -nFormat
	}
	for nFormat > 0 {
		digitFormat := nFormat % 10
		digitsFormat = append(digitsFormat, digitStrsFormat[digitFormat])
		nFormat = nFormat / 10
	}
	result := ""
	if negativeFormat {
		result = "-"
	}
	for iDigitFormat := len(digitsFormat) - 1; iDigitFormat >= 0; iDigitFormat = iDigitFormat - 1 {
		result = result + digitsFormat[iDigitFormat]
	}
	return result
}

func Contains(s string, substr string) bool {
	if len(substr) == 0 {
		return 1 == 1
	}
	if len(substr) > len(s) {
		return 0 == 1
	}
	for iContains := 0; iContains <= len(s)-len(substr); iContains = iContains + 1 {
		match := 1 == 1
		for jContains := 0; jContains < len(substr); jContains = jContains + 1 {
			if s[iContains+jContains] != substr[jContains] {
				match = 0 == 1
				break
			}
		}
		if match {
			return 1 == 1
		}
	}
	return 0 == 1
}

func containsAny(s string, chars string) bool {
	for iAny := 0; iAny < len(s); iAny = iAny + 1 {
		for jAny := 0; jAny < len(chars); jAny = jAny + 1 {
			if s[iAny] == chars[jAny] {
				return 1 == 1
			}
		}
	}
	return 0 == 1
}
