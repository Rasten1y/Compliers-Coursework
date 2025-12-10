package backend

import (
	"strconv"
	"strings"

	"gominic/ir"
)

// EmitModule renders an IR module into textual LLVM IR.
func EmitModule(mod *ir.Module) string {
	var e emitter
	emitTargetHeader(&e, mod)
	emitGlobals(&e, mod.Globals)
	for i := 0; i < len(mod.Functions); i++ {
		emitFunction(&e, mod.Functions[i])
	}
	return strBufString(&e.buf)
}

type strBuf struct {
	parts []string
}

// Plain helper functions instead of methods to simplify self-host lowering.
func strBufWriteString(b *strBuf, s string) {
	b.parts = append(b.parts, s)
}

func strBufString(b *strBuf) string {
	return strings.Join(b.parts, "")
}

type emitter struct {
	buf  strBuf
	msvc bool
}

func emitGlobals(e *emitter, globals []ir.Global) {
	// runtime declarations
	strBufWriteString(&e.buf, "declare void @gominic_memcpy(i8*, i8*, i64)\n")
	strBufWriteString(&e.buf, "declare void @gominic_abort()\n")
	strBufWriteString(&e.buf, "declare i8* @gominic_makeSlice(i64, i64, i64)\n")
	strBufWriteString(&e.buf, "declare void @gominic_print(i8*, i64)\n")
	strBufWriteString(&e.buf, "declare void @gominic_printInt(i64)\n")
	strBufWriteString(&e.buf, "declare void @gominic_println()\n\n")
	strBufWriteString(&e.buf, "declare i1 @gominic_str_eq(i8*, i64, i8*, i64)\n")
	strBufWriteString(&e.buf, "declare void @gominic_str_concat({ i8*, i64 }*, { i8*, i64 }*, { i8*, i64 }*)\n\n")
	strBufWriteString(&e.buf, "declare i8* @gominic_map_new(i64, i64, i32)\n")
	strBufWriteString(&e.buf, "declare void @gominic_map_set(i8*, i8*, i8*)\n")
	strBufWriteString(&e.buf, "declare i1 @gominic_map_get(i8*, i8*, i8*)\n")
	strBufWriteString(&e.buf, "declare i64 @gominic_map_len(i8*)\n\n")
	strBufWriteString(&e.buf, "declare i64 @gominic_argc()\n")
	if e.msvc {
		strBufWriteString(&e.buf, "declare void @gominic_argv({ i8*, i64 }* sret({ i8*, i64 }), i64)\n")
	} else {
		strBufWriteString(&e.buf, "declare { i8*, i64 } @gominic_argv(i64)\n")
	}
	strBufWriteString(&e.buf, "declare i1 @gominic_write_file(i8*, i64, i8*, i64)\n\n")
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
		strBufWriteString(&e.buf, "@"+g.Name+" = "+visibility+"constant "+llvmType(g.Type)+" "+g.Value+align+"\n")
	}
	if len(globals) > 0 {
		strBufWriteString(&e.buf, "\n")
	}
}

func emitFunction(e *emitter, fn *ir.Function) {
	params := make([]string, len(fn.Params))
	for i := 0; i < len(fn.Params); i++ {
		p := fn.Params[i]
		if i == 0 && fn.SretType != nil {
			params[i] = "sret(" + llvmType(fn.SretType) + ") " + llvmType(p.Type()) + " %" + p.Name()
		} else if p.ByVal() {
			params[i] = llvmType(p.Type()) + " byval(" + llvmType(p.ByValType()) + ") %" + p.Name()
		} else {
			params[i] = llvmType(p.Type()) + " %" + p.Name()
		}
	}

	var result string
	if fn.SretType != nil {
		result = llvmType(ir.Void)
	} else {
		resCount := len(fn.Results)
		if resCount == 0 {
			result = llvmType(ir.Void)
		} else if resCount == 1 {
			result = llvmType(fn.Results[0])
		} else {
			// TODO: support multiple return values via struct or sret lowering.
			result = llvmType(ir.Void)
		}
	}

	strBufWriteString(&e.buf, "define "+result+" @"+fn.Name+"("+strings.Join(params, ", ")+") {\n")
	for i := 0; i < len(fn.Blocks); i++ {
		emitBasicBlock(e, fn.Blocks[i])
	}
	strBufWriteString(&e.buf, "}\n\n")
}

func emitBasicBlock(e *emitter, bb *ir.BasicBlock) {
	strBufWriteString(&e.buf, bb.Name+":\n")
	for i := 0; i < len(bb.Instrs); i++ {
		strBufWriteString(&e.buf, "  "+bb.Instrs[i].String()+"\n")
	}
	if bb.Terminator != nil {
		strBufWriteString(&e.buf, "  "+bb.Terminator.String()+"\n")
	} else {
		// Fallback to avoid invalid IR.
		strBufWriteString(&e.buf, "  unreachable\n")
	}
}

func llvmType(t ir.Type) string {
	if t == nil {
		return "void"
	}
	return t.String()
}

func emitTargetHeader(e *emitter, mod *ir.Module) {
	triple := mod.TargetTriple
	if triple == "" {
		triple = "x86_64-pc-linux-gnu"
	}
	e.msvc = strings.Contains(triple, "windows-msvc")
	dl := mod.DataLayout
	if dl == "" {
		dl = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
	}
	strBufWriteString(&e.buf, "target triple = \""+triple+"\"\n")
	strBufWriteString(&e.buf, "target datalayout = \""+dl+"\"\n\n")
}
