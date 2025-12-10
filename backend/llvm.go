package backend

import (
	"bytes"
	"strconv"
	"strings"

	"gominic/ir"
)

// EmitModule renders an IR module into textual LLVM IR.
func EmitModule(mod *ir.Module) string {
	var e emitter
	e.emitTargetHeader(mod)
	e.emitGlobals(mod.Globals)
	for i := 0; i < len(mod.Functions); i++ {
		e.emitFunction(mod.Functions[i])
	}
	return e.buf.String()
}

type emitter struct {
	buf bytes.Buffer
}

func (e *emitter) emitGlobals(globals []ir.Global) {
	// runtime declarations
	e.buf.WriteString("declare void @gominic_memcpy(i8*, i8*, i64)\n")
	e.buf.WriteString("declare void @gominic_abort()\n")
	e.buf.WriteString("declare i8* @gominic_makeSlice(i64, i64, i64)\n")
	e.buf.WriteString("declare void @gominic_print(i8*, i64)\n")
	e.buf.WriteString("declare void @gominic_printInt(i64)\n")
	e.buf.WriteString("declare void @gominic_println()\n\n")
	e.buf.WriteString("declare i1 @gominic_str_eq(i8*, i64, i8*, i64)\n")
	e.buf.WriteString("declare { i8*, i64 } @gominic_str_concat({ i8*, i64 }, { i8*, i64 })\n\n")
	e.buf.WriteString("declare i8* @gominic_map_new(i64, i64, i32)\n")
	e.buf.WriteString("declare void @gominic_map_set(i8*, i8*, i8*)\n")
	e.buf.WriteString("declare i1 @gominic_map_get(i8*, i8*, i8*)\n")
	e.buf.WriteString("declare i64 @gominic_map_len(i8*)\n\n")
	e.buf.WriteString("declare i64 @gominic_argc()\n")
	e.buf.WriteString("declare { i8*, i64 } @gominic_argv(i64)\n")
	e.buf.WriteString("declare i1 @gominic_write_file(i8*, i64, i8*, i64)\n\n")
	for i := 0; i < len(globals); i++ {
		g := globals[i]
		visibility := ""
		if g.Private {
			visibility = "private unnamed_addr "
		}
		align := ""
		if g.Align > 0 {
			align = ", align " + strconv.FormatInt(g.Align, 10)
		}
		e.buf.WriteString("@" + g.Name + " = " + visibility + "constant " + e.llvmType(g.Type) + " " + g.Value + align + "\n")
	}
	if len(globals) > 0 {
		e.buf.WriteString("\n")
	}
}

func (e *emitter) emitFunction(fn *ir.Function) {
	params := make([]string, len(fn.Params))
	for i := 0; i < len(fn.Params); i++ {
		p := fn.Params[i]
		if i == 0 && fn.SretType != nil {
			params[i] = "sret(" + e.llvmType(fn.SretType) + ") " + e.llvmType(p.Type()) + " %" + p.Name()
		} else if p.ByVal() {
			params[i] = e.llvmType(p.Type()) + " byval(" + e.llvmType(p.ByValType()) + ") %" + p.Name()
		} else {
			params[i] = e.llvmType(p.Type()) + " %" + p.Name()
		}
	}

	var result string
	if fn.SretType != nil {
		result = ir.Void.String()
	} else {
		resCount := len(fn.Results)
		if resCount == 0 {
			result = ir.Void.String()
		} else if resCount == 1 {
			result = e.llvmType(fn.Results[0])
		} else {
			// TODO: support multiple return values via struct or sret lowering.
			result = ir.Void.String()
		}
	}

	e.buf.WriteString("define " + result + " @" + fn.Name + "(" + strings.Join(params, ", ") + ") {\n")
	for i := 0; i < len(fn.Blocks); i++ {
		e.emitBasicBlock(fn.Blocks[i])
	}
	e.buf.WriteString("}\n\n")
}

func (e *emitter) emitBasicBlock(bb *ir.BasicBlock) {
	e.buf.WriteString(bb.Name + ":\n")
	for i := 0; i < len(bb.Instrs); i++ {
		e.buf.WriteString("  " + bb.Instrs[i].String() + "\n")
	}
	if bb.Terminator != nil {
		e.buf.WriteString("  " + bb.Terminator.String() + "\n")
	} else {
		// TODO: enforce presence of a terminator in earlier phases.
		e.buf.WriteString("  ; missing terminator\n")
	}
}

func (e *emitter) llvmType(t ir.Type) string {
	if t == nil {
		return "void"
	}
	return t.String()
}

func (e *emitter) emitTargetHeader(mod *ir.Module) {
	triple := mod.TargetTriple
	if triple == "" {
		triple = "x86_64-pc-linux-gnu"
	}
	dl := mod.DataLayout
	if dl == "" {
		dl = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
	}
	e.buf.WriteString("target triple = \"" + triple + "\"\n")
	e.buf.WriteString("target datalayout = \"" + dl + "\"\n\n")
}
