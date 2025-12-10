package ir

import (
	"strconv"
	"strings"
)

// Type is a minimal LLVM-ish type description used by our IR.
type Type interface {
	String() string
}

// BasicType holds a raw LLVM IR type name (e.g. "i64", "double").
type BasicType string

// StructType is an aggregate of fields.
type StructType struct {
	Fields []Type
}

// ArrayType is a fixed-length array.
type ArrayType struct {
	Len  int
	Elem Type
}

// SliceType is a runtime slice descriptor { data *T, len i64, cap i64 }.
type SliceType struct {
	Elem Type
}

// StringType is lowered as { i8*, i64 } (pointer + length).
type StringType struct{}

// PointerType models T*.
type PointerType struct {
	Elem Type
}

// Predefined primitive types for convenience.
var (
	I1    = BasicType("i1")
	I8    = BasicType("i8")
	I32   = BasicType("i32")
	I64   = BasicType("i64")
	F64   = BasicType("double")
	Void  = BasicType("void")
	PtrI8 = BasicType("i8*")
)

func (t BasicType) String() string { return string(t) }

func (p PointerType) String() string { return p.Elem.String() + "*" }

// PtrTo constructs a pointer type to elem.
func PtrTo(elem Type) PointerType {
	return PointerType{Elem: elem}
}

func (s StructType) String() string {
	parts := make([]string, len(s.Fields))
	for i := 0; i < len(s.Fields); i++ {
		parts[i] = s.Fields[i].String()
	}
	return "{ " + strings.Join(parts, ", ") + " }"
}

func (a ArrayType) String() string {
	return "[" + strconv.Itoa(a.Len) + " x " + a.Elem.String() + "]"
}

func (s SliceType) String() string {
	// Avoid calling methods on package-level BasicType singletons to keep self-hosting simpler.
	return "{ " + PtrTo(s.Elem).String() + ", i64, i64 }"
}

func (StringType) String() string {
	// String is lowered as data pointer + length.
	return "{ i8*, i64 }"
}

// Value represents a typed value in the IR (register, parameter or constant).
type valueData struct {
	name     string
	ty       Type
	kind     ValueKind
	Raw      string // for constants
	byVal    bool
	byValTyp Type
}

// Value is a pointer to the underlying valueData so nil can be used where helpful.
type Value = *valueData

// ValueKind distinguishes concrete value categories.
type ValueKind int

const (
	ValueInvalid ValueKind = iota
	ValueParam
	ValueRegister
	ValueConstant
)

func (v Value) Type() Type {
	if v == nil {
		return nil
	}
	return v.ty
}
func (v Value) Name() string {
	if v == nil {
		return ""
	}
	if v.kind == ValueConstant && v.Raw != "" {
		return v.Raw
	}
	return v.name
}
func (v Value) ByVal() bool {
	if v == nil {
		return false
	}
	return v.byVal
}
func (v Value) ByValType() Type {
	if v == nil {
		return nil
	}
	return v.byValTyp
}
func (v Value) Kind() ValueKind {
	if v == nil {
		return ValueInvalid
	}
	return v.kind
}

// Param represents a function argument.
type Param = Value

func NewParam(name string, ty Type) Param {
	return &valueData{name: name, ty: ty, kind: ValueParam}
}

func NewByValParam(name string, ptr Type, byValTyp Type) Param {
	return &valueData{name: name, ty: ptr, byVal: true, byValTyp: byValTyp, kind: ValueParam}
}

// Register is a named SSA-like value.
type Register = Value

func NewRegister(name string, ty Type) Register {
	return &valueData{name: name, ty: ty, kind: ValueRegister}
}

// Constant encodes a literal value with its type.
type Constant = Value

func NewConstant(raw string, ty Type) Constant {
	return &valueData{name: raw, Raw: raw, ty: ty, kind: ValueConstant}
}

// Instruction is any IR instruction. Most also produce a value.
type Instruction interface {
	isInstruction()
	String() string
}

// BinOpKind enumerates basic arithmetic and comparison ops.
type BinOpKind string

const (
	Add  BinOpKind = "add"
	Sub  BinOpKind = "sub"
	Mul  BinOpKind = "mul"
	SDiv BinOpKind = "sdiv"
	UDiv BinOpKind = "udiv"
	SRem BinOpKind = "srem"
	URem BinOpKind = "urem"

	FAdd BinOpKind = "fadd"
	FSub BinOpKind = "fsub"
	FMul BinOpKind = "fmul"
	FDiv BinOpKind = "fdiv"

	And BinOpKind = "and"
	Or  BinOpKind = "or"
)

// BinOp computes X op Y and stores the result in Dest.
type BinOp struct {
	Dest Register
	Op   BinOpKind
	X    Value
	Y    Value
}

func (BinOp) isInstruction() {}
func (i BinOp) String() string {
	return "%" + i.Dest.name + " = " + string(i.Op) + " " + typeString(i.X) + " " + formatValue(i.X) + ", " + formatValue(i.Y)
}

// Return terminates a basic block.
type Return struct {
	Results []Value
}

func (Return) isInstruction() {}
func (r Return) String() string {
	if len(r.Results) == 0 {
		return "ret void"
	}
	if len(r.Results) == 1 {
		val := r.Results[0]
		return "ret " + typeString(val) + " " + formatValue(val)
	}
	// TODO: support multiple return values via structs or sret.
	return "ret void ; TODO multiple return values"
}

// Call invokes a function. Dest is optional: if Name is empty or Dest.Type is void, a void call is emitted.
type Call struct {
	Dest Register
	Name string
	Args []Value
	Ret  Type
	// SretType marks the call as sret; first argument must be the sret pointer.
	SretType Type
}

func (Call) isInstruction() {}
func (c Call) String() string {
	args := make([]string, len(c.Args))
	for i := 0; i < len(c.Args); i++ {
		a := c.Args[i]
		prefix := ""
		if i == 0 && c.SretType != nil {
			prefix = "sret(" + c.SretType.String() + ") "
		}
		tyStr := "i8*"
		valStr := "null"
		if a != nil {
			if a.ty != nil {
				tyStr = a.ty.String()
			}
			if a.kind == ValueRegister || a.kind == ValueParam {
				valStr = "%" + a.name
			} else {
				if a.kind == ValueConstant && a.Raw != "" {
					valStr = a.Raw
				} else {
					valStr = a.name
				}
			}
		}
		args[i] = prefix + tyStr + " " + valStr
	}
	ret := "void"
	if c.Ret != nil {
		ret = c.Ret.String()
	}
	prefix := ""
	if c.Dest != nil && c.Dest.name != "" && ret != "void" {
		prefix = "%" + c.Dest.name + " = "
	}
	return prefix + "call " + ret + " @" + c.Name + "(" + strings.Join(args, ", ") + ")"
}

// Conv converts a value to another type using a specific opcode.
type Conv struct {
	Dest Register
	Op   ConvOp
	Src  Value
	To   Type
}

// ConvOp enumerates conversion operations.
type ConvOp string

const (
	Trunc  ConvOp = "trunc"
	ZExt   ConvOp = "zext"
	SExt   ConvOp = "sext"
	SIToFP ConvOp = "sitofp"
	UIToFP ConvOp = "uitofp"
	FPToSI ConvOp = "fptosi"
	FPToUI ConvOp = "fptoui"
)

func (Conv) isInstruction() {}
func (c Conv) String() string {
	return "%" + c.Dest.name + " = " + string(c.Op) + " " + typeString(c.Src) + " " + formatValue(c.Src) + " to " + c.To.String()
}

// Alloca allocates stack memory for a given type.
type Alloca struct {
	Dest      Register
	AllocType Type
	Align     int64
}

func (Alloca) isInstruction() {}
func (a Alloca) String() string {
	align := ""
	if a.Align > 0 {
		align = " , align " + strconv.FormatInt(a.Align, 10)
	}
	return "%" + a.Dest.name + " = alloca " + a.AllocType.String() + align
}

// Load reads from pointer Src into Dest.
type Load struct {
	Dest  Register
	Src   Value
	Align int64
}

func (Load) isInstruction() {}
func (l Load) String() string {
	align := ""
	if l.Align > 0 {
		align = ", align " + strconv.FormatInt(l.Align, 10)
	}
	return "%" + l.Dest.name + " = load " + typeString(l.Dest) + ", " + typeString(l.Src) + " " + formatValue(l.Src) + align
}

// Store writes Src into Dest pointer.
type Store struct {
	Src   Value
	Dest  Value
	Align int64
}

func (Store) isInstruction() {}
func (s Store) String() string {
	align := ""
	if s.Align > 0 {
		align = ", align " + strconv.FormatInt(s.Align, 10)
	}
	return "store " + typeString(s.Src) + " " + formatValue(s.Src) + ", " + typeString(s.Dest) + " " + formatValue(s.Dest) + align
}

// GetElementPtr computes an address inside an aggregate.
type GetElementPtr struct {
	Dest      Register
	Src       Value
	Pointee   Type // type pointed to by Src
	Indices   []Value
	ResultTyp Type // pointer to element
}

func (GetElementPtr) isInstruction() {}
func (g GetElementPtr) String() string {
	parts := make([]string, len(g.Indices))
	for i := 0; i < len(g.Indices); i++ {
		idx := g.Indices[i]
		parts[i] = typeString(idx) + " " + formatValue(idx)
	}
	return "%" + g.Dest.name + " = getelementptr inbounds " + g.Pointee.String() + ", " + typeString(g.Src) + " " + formatValue(g.Src) + ", " + strings.Join(parts, ", ")
}

// ICmp compares integers.
type ICmp struct {
	Dest Register
	Pred ICmpPred
	X    Value
	Y    Value
}

// ICmpPred lists integer comparison predicates.
type ICmpPred string

const (
	ICmpEq  ICmpPred = "eq"
	ICmpNe  ICmpPred = "ne"
	ICmpSlt ICmpPred = "slt"
	ICmpSle ICmpPred = "sle"
	ICmpSgt ICmpPred = "sgt"
	ICmpSge ICmpPred = "sge"
	ICmpUlt ICmpPred = "ult"
	ICmpUle ICmpPred = "ule"
	ICmpUgt ICmpPred = "ugt"
	ICmpUge ICmpPred = "uge"
)

func (ICmp) isInstruction() {}
func (c ICmp) String() string {
	return "%" + c.Dest.name + " = icmp " + string(c.Pred) + " " + typeString(c.X) + " " + formatValue(c.X) + ", " + formatValue(c.Y)
}

// FCmp compares floats.
type FCmp struct {
	Dest Register
	Pred FCmpPred
	X    Value
	Y    Value
}

// FCmpPred lists float comparison predicates.
type FCmpPred string

const (
	FCmpOeq FCmpPred = "oeq"
	FCmpOne FCmpPred = "one"
	FCmpOlt FCmpPred = "olt"
	FCmpOle FCmpPred = "ole"
	FCmpOgt FCmpPred = "ogt"
	FCmpOge FCmpPred = "oge"
)

func (FCmp) isInstruction() {}
func (c FCmp) String() string {
	return "%" + c.Dest.name + " = fcmp " + string(c.Pred) + " " + typeString(c.X) + " " + formatValue(c.X) + ", " + formatValue(c.Y)
}

// Br performs an unconditional branch.
type Br struct {
	Target string
}

func (Br) isInstruction() {}
func (b Br) String() string {
	return "br label %" + b.Target
}

// CondBr branches to True or False depending on Cond.
type CondBr struct {
	Cond  Value
	True  string
	False string
}

func (CondBr) isInstruction() {}
func (b CondBr) String() string {
	return "br i1 " + formatValue(b.Cond) + ", label %" + b.True + ", label %" + b.False
}

// CallVoid is a helper to emit direct call without result.
type CallVoid struct {
	Name string
	Args []Value
}

func (CallVoid) isInstruction() {}
func (c CallVoid) String() string {
	args := make([]string, len(c.Args))
	for i := 0; i < len(c.Args); i++ {
		a := c.Args[i]
		args[i] = typeString(a) + " " + formatValue(a)
	}
	return "call void @" + c.Name + "(" + strings.Join(args, ", ") + ")"
}

// BasicBlock is a straight-line sequence terminated by a control instruction.
type BasicBlock struct {
	Name       string
	Instrs     []Instruction
	Terminator Instruction
}

func NewBasicBlock(name string) *BasicBlock {
	return &BasicBlock{Name: name}
}

func (bb *BasicBlock) Append(inst Instruction) {
	bb.Instrs = append(bb.Instrs, inst)
}

// Function groups blocks into a callable unit.
type Function struct {
	Name    string
	Params  []Param
	Results []Type
	Blocks  []*BasicBlock
	// If SretType is set, the function returns aggregate via sret pointer (first param).
	SretType Type
}

func (fn *Function) Entry() *BasicBlock {
	if len(fn.Blocks) == 0 {
		entry := NewBasicBlock("entry")
		fn.Blocks = append(fn.Blocks, entry)
	}
	return fn.Blocks[0]
}

// Module is a collection of functions.
type Module struct {
	Functions []*Function
	Globals   []Global
	// TargetTriple and DataLayout describe the LLVM target; if empty, defaults are used.
	TargetTriple string
	DataLayout   string
}

func (m *Module) AddFunction(fn *Function) {
	m.Functions = append(m.Functions, fn)
}

// Global represents a global variable/constant.
type Global struct {
	Name    string
	Type    Type
	Value   string
	Align   int64
	Private bool
}

// Bitcast converts a pointer to another pointer type.
type Bitcast struct {
	Dest   Register
	Src    Value
	Target Type
}

func (Bitcast) isInstruction() {}
func (b Bitcast) String() string {
	return "%" + b.Dest.name + " = bitcast " + typeString(b.Src) + " " + formatValue(b.Src) + " to " + b.Target.String()
}

// Memcpy copies n bytes from Src to Dest.
type Memcpy struct {
	Dest Value
	Src  Value
	Size Value
}

func (Memcpy) isInstruction() {}
func (m Memcpy) String() string {
	return "call void @gominic_memcpy(" + typeString(m.Dest) + " " + formatValue(m.Dest) + ", " + typeString(m.Src) + " " + formatValue(m.Src) + ", " + typeString(m.Size) + " " + formatValue(m.Size) + ")"
}

// formatValue renders a Value reference for use in instruction text.
func formatValue(v Value) string {
	if v == nil {
		return ""
	}
	if v.kind == ValueRegister || v.kind == ValueParam {
		return "%" + v.name
	}
	if v.kind == ValueConstant && v.Raw != "" {
		if v.ty != nil && v.ty.String() == "double" {
			if !strings.ContainsAny(v.Raw, ".eE") {
				return v.Raw + ".0"
			}
		}
		return v.Raw
	}
	return v.name
}

func typeString(v Value) string {
	if v != nil && v.ty != nil {
		return v.ty.String()
	}
	return "void"
}

// ValueType returns the static IR type stored in v (nil if unset).
func ValueType(v Value) Type {
	if v == nil {
		return nil
	}
	return v.ty
}
