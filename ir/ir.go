package ir

type TypeKind int

const (
	KindInvalid TypeKind = iota
	KindBasic
	KindPointer
	KindStruct
	KindArray
	KindSlice
	KindString
)

type TypeDesc struct {
	Kind   TypeKind
	Basic  string      // for KindBasic
	Elem   *TypeDesc   // ptr/array/slice
	Fields []*TypeDesc // struct fields
	Len    int         // array length
}

var (
	I1    = &TypeDesc{Kind: KindBasic, Basic: "i1"}
	I8    = &TypeDesc{Kind: KindBasic, Basic: "i8"}
	I32   = &TypeDesc{Kind: KindBasic, Basic: "i32"}
	I64   = &TypeDesc{Kind: KindBasic, Basic: "i64"}
	F64   = &TypeDesc{Kind: KindBasic, Basic: "double"}
	Void  = &TypeDesc{Kind: KindBasic, Basic: "void"}
	PtrI8 = &TypeDesc{Kind: KindPointer, Elem: I8}
	strTy = &TypeDesc{Kind: KindString}
)

func PtrTo(elem *TypeDesc) *TypeDesc {
	return &TypeDesc{Kind: KindPointer, Elem: elem}
}

func Struct(fields []*TypeDesc) *TypeDesc {
	return &TypeDesc{Kind: KindStruct, Fields: fields}
}

func Array(length int, elem *TypeDesc) *TypeDesc {
	return &TypeDesc{Kind: KindArray, Len: length, Elem: elem}
}

func Slice(elem *TypeDesc) *TypeDesc {
	return &TypeDesc{Kind: KindSlice, Elem: elem}
}

func String() *TypeDesc { return strTy }

func (t *TypeDesc) String() string {
	if t == nil {
		return "void"
	}
	if t.Kind == KindBasic {
		return t.Basic
	} else if t.Kind == KindPointer {
		return t.Elem.String() + "*"
	} else if t.Kind == KindStruct {
		parts := make([]string, len(t.Fields))
		for i := 0; i < len(t.Fields); i = i + 1 {
			parts[i] = t.Fields[i].String()
		}
		// Собираем строку вручную (без strings.Join для самоприменимости)
		if len(parts) == 0 {
			return "{ }"
		}
		result := "{ " + parts[0]
		for i := 1; i < len(parts); i = i + 1 {
			result = result + ", " + parts[i]
		}
		return result + " }"
	} else if t.Kind == KindArray {
		lenStr := formatInt64Helper(int64(t.Len))
		return "[" + lenStr + " x " + t.Elem.String() + "]"
	} else if t.Kind == KindSlice {
		return "{ " + PtrTo(t.Elem).String() + ", i64, i64 }"
	} else if t.Kind == KindString {
		return "{ i8*, i64 }"
	} else {
		return "void"
	}
}

type ValueKind int

const (
	ValueInvalid ValueKind = iota
	ValueParam
	ValueRegister
	ValueConstant
)

type Value struct {
	name     string
	ty       *TypeDesc
	kind     ValueKind
	Raw      string // for constants
	byVal    bool
	byValTyp *TypeDesc
}

func (v *Value) Type() *TypeDesc {
	if v == nil {
		return nil
	}
	return v.ty
}

func (v *Value) Name() string {
	if v == nil {
		return ""
	}
	if v.kind == ValueConstant && v.Raw != "" {
		return v.Raw
	}
	return v.name
}

func (v *Value) ByVal() bool {
	if v == nil {
		return false
	}
	return v.byVal
}

func (v *Value) ByValType() *TypeDesc {
	if v == nil {
		return nil
	}
	return v.byValTyp
}

func (v *Value) Kind() ValueKind {
	if v == nil {
		return ValueInvalid
	}
	return v.kind
}

func NewParam(name string, ty *TypeDesc) *Value {
	return &Value{name: name, ty: ty, kind: ValueParam}
}

// Helper functions for self-hosted backend to access unexported fields
func ValueName(v *Value) string {
	if v == nil {
		return ""
	}
	if v.kind == ValueConstant && v.Raw != "" {
		return v.Raw
	}
	return v.name
}

// ValueType: helper функция для доступа к полю type (неэкспортированному)
func ValueType(v *Value) *TypeDesc {
	if v == nil {
		return nil
	}
	return v.ty
}

func GetValueKind(v *Value) ValueKind {
	if v == nil {
		return ValueInvalid
	}
	return v.kind
}

func IsValueParam(k ValueKind) bool {
	return k == ValueParam
}

func IsValueRegister(k ValueKind) bool {
	return k == ValueRegister
}

func IsValueConstant(k ValueKind) bool {
	return k == ValueConstant
}

func GetTypeKind(t *TypeDesc) TypeKind {
	if t == nil {
		return KindInvalid
	}
	return t.Kind
}

func IsKindBasic(k TypeKind) bool {
	return k == KindBasic
}

func IsKindPointer(k TypeKind) bool {
	return k == KindPointer
}

func IsKindStruct(k TypeKind) bool {
	return k == KindStruct
}

func IsKindArray(k TypeKind) bool {
	return k == KindArray
}

func IsKindSlice(k TypeKind) bool {
	return k == KindSlice
}

func IsKindString(k TypeKind) bool {
	return k == KindString
}

func GetTypeDescBasic(t *TypeDesc) string {
	if t == nil {
		return ""
	}
	return t.Basic
}

func GetTypeDescElem(t *TypeDesc) *TypeDesc {
	if t == nil {
		return nil
	}
	return t.Elem
}

func GetTypeDescFields(t *TypeDesc) []*TypeDesc {
	if t == nil {
		return nil
	}
	return t.Fields
}

func GetTypeDescLen(t *TypeDesc) int {
	if t == nil {
		return 0
	}
	return t.Len
}

func GetVoidType() *TypeDesc {
	return Void
}

func GetI1Type() *TypeDesc {
	return I1
}

func GetI8Type() *TypeDesc {
	return I8
}

func GetI32Type() *TypeDesc {
	return I32
}

func GetI64Type() *TypeDesc {
	return I64
}

func GetF64Type() *TypeDesc {
	return F64
}

func GetPtrI8Type() *TypeDesc {
	return PtrI8
}

func GetInstrKind(inst *Instruction) InstrKind {
	if inst == nil {
		return InstrInvalid
	}
	return inst.Kind
}

func IsInstrBinOp(k InstrKind) bool {
	return k == InstrBinOp
}

func IsInstrReturn(k InstrKind) bool {
	return k == InstrReturn
}

func IsInstrCall(k InstrKind) bool {
	return k == InstrCall
}

func IsInstrConv(k InstrKind) bool {
	return k == InstrConv
}

func IsInstrAlloca(k InstrKind) bool {
	return k == InstrAlloca
}

func IsInstrLoad(k InstrKind) bool {
	return k == InstrLoad
}

func IsInstrStore(k InstrKind) bool {
	return k == InstrStore
}

func IsInstrGEP(k InstrKind) bool {
	return k == InstrGEP
}

func IsInstrICmp(k InstrKind) bool {
	return k == InstrICmp
}

func IsInstrFCmp(k InstrKind) bool {
	return k == InstrFCmp
}

func IsInstrBr(k InstrKind) bool {
	return k == InstrBr
}

func IsInstrCondBr(k InstrKind) bool {
	return k == InstrCondBr
}

func IsInstrCallVoid(k InstrKind) bool {
	return k == InstrCallVoid
}

func IsInstrBitcast(k InstrKind) bool {
	return k == InstrBitcast
}

func IsInstrMemcpy(k InstrKind) bool {
	return k == InstrMemcpy
}

func GetInstrDest(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.Dest
}

func GetInstrBinOp(inst *Instruction) BinOpKind {
	if inst == nil {
		return ""
	}
	return inst.BinOp
}

func GetInstrX(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.X
}

func GetInstrY(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.Y
}

func GetInstrRetVals(inst *Instruction) []*Value {
	if inst == nil {
		return nil
	}
	return inst.RetVals
}

func GetInstrCallName(inst *Instruction) string {
	if inst == nil {
		return ""
	}
	return inst.CallName
}

func GetInstrCallArgs(inst *Instruction) []*Value {
	if inst == nil {
		return nil
	}
	return inst.CallArgs
}

func GetInstrCallRet(inst *Instruction) *TypeDesc {
	if inst == nil {
		return nil
	}
	return inst.CallRet
}

func GetInstrConvOp(inst *Instruction) ConvOp {
	if inst == nil {
		return ""
	}
	return inst.ConvOp
}

func GetInstrConvSrc(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.ConvSrc
}

func GetInstrConvTo(inst *Instruction) *TypeDesc {
	if inst == nil {
		return nil
	}
	return inst.ConvTo
}

func GetInstrAllocaType(inst *Instruction) *TypeDesc {
	if inst == nil {
		return nil
	}
	return inst.AllocaType
}

func GetInstrAllocaAlign(inst *Instruction) int64 {
	if inst == nil {
		return 0
	}
	return inst.AllocaAlign
}

func GetInstrLoadSrc(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.LoadSrc
}

func GetInstrLoadAlign(inst *Instruction) int64 {
	if inst == nil {
		return 0
	}
	return inst.LoadAlign
}

func GetInstrStoreSrc(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.StoreSrc
}

func GetInstrStoreDst(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.StoreDst
}

func GetInstrStoreAlign(inst *Instruction) int64 {
	if inst == nil {
		return 0
	}
	return inst.StoreAlign
}

func GetInstrGepSrc(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.GepSrc
}

func GetInstrGepPointee(inst *Instruction) *TypeDesc {
	if inst == nil {
		return nil
	}
	return inst.GepPointee
}

func GetInstrGepIndices(inst *Instruction) []*Value {
	if inst == nil {
		return nil
	}
	return inst.GepIndices
}

func GetInstrICmpPred(inst *Instruction) ICmpPred {
	if inst == nil {
		return ""
	}
	return inst.ICmpPred
}

func GetInstrICmpX(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.ICmpX
}

func GetInstrICmpY(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.ICmpY
}

func GetInstrFCmpPred(inst *Instruction) FCmpPred {
	if inst == nil {
		return ""
	}
	return inst.FCmpPred
}

func GetInstrFCmpX(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.FCmpX
}

func GetInstrFCmpY(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.FCmpY
}

func GetInstrBrTarget(inst *Instruction) string {
	if inst == nil {
		return ""
	}
	return inst.BrTarget
}

func GetInstrCondCond(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.CondCond
}

func GetInstrCondTrue(inst *Instruction) string {
	if inst == nil {
		return ""
	}
	return inst.CondTrue
}

func GetInstrCondFalse(inst *Instruction) string {
	if inst == nil {
		return ""
	}
	return inst.CondFalse
}

func GetInstrBitcastSrc(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.BitcastSrc
}

func GetInstrBitcastTarget(inst *Instruction) *TypeDesc {
	if inst == nil {
		return nil
	}
	return inst.BitcastTarget
}

func GetInstrMemcpyDest(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.MemcpyDest
}

func GetInstrMemcpySrc(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.MemcpySrc
}

func GetInstrMemcpySize(inst *Instruction) *Value {
	if inst == nil {
		return nil
	}
	return inst.MemcpySize
}

func ValueByVal(v *Value) bool {
	if v == nil {
		return false
	}
	return v.byVal
}

func ValueByValType(v *Value) *TypeDesc {
	if v == nil {
		return nil
	}
	return v.byValTyp
}

func TypeDescString(t *TypeDesc) string {
	if t == nil {
		return "void"
	}
	if t.Kind == KindBasic {
		return t.Basic
	} else if t.Kind == KindPointer {
		return TypeDescString(t.Elem) + "*"
	} else if t.Kind == KindStruct {
		// Build string manually without strings.Join
		if len(t.Fields) == 0 {
			return "{ }"
		}
		result := "{ " + TypeDescString(t.Fields[0])
		for i := 1; i < len(t.Fields); i = i + 1 {
			result = result + ", " + TypeDescString(t.Fields[i])
		}
		return result + " }"
	} else if t.Kind == KindArray {
		// Convert int to string manually
		lenStr := formatInt64Helper(int64(t.Len))
		return "[" + lenStr + " x " + TypeDescString(t.Elem) + "]"
	} else if t.Kind == KindSlice {
		return "{ " + TypeDescString(PtrTo(t.Elem)) + ", i64, i64 }"
	} else if t.Kind == KindString {
		return "{ i8*, i64 }"
	} else {
		return "void"
	}
}

func formatInt64Helper(n int64) string {
	if n == 0 {
		return "0"
	}
	digitStrs := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var digits []string
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}
	for n > 0 {
		digit := n % 10
		digits = append(digits, digitStrs[digit])
		n = n / 10
	}
	result := ""
	if negative {
		result = "-"
	}
	for i := len(digits) - 1; i >= 0; i = i - 1 {
		result = result + digits[i]
	}
	return result
}

func NewByValParam(name string, ptr *TypeDesc, byValTyp *TypeDesc) *Value {
	return &Value{name: name, ty: ptr, byVal: true, byValTyp: byValTyp, kind: ValueParam}
}

func NewRegister(name string, ty *TypeDesc) *Value {
	return &Value{name: name, ty: ty, kind: ValueRegister}
}

func NewConstant(raw string, ty *TypeDesc) *Value {
	return &Value{name: raw, Raw: raw, ty: ty, kind: ValueConstant}
}


type InstrKind int

const (
	InstrInvalid InstrKind = iota
	InstrBinOp
	InstrReturn
	InstrCall
	InstrConv
	InstrAlloca
	InstrLoad
	InstrStore
	InstrGEP
	InstrICmp
	InstrFCmp
	InstrBr
	InstrCondBr
	InstrCallVoid
	InstrBitcast
	InstrMemcpy
)

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
	
	Shl  BinOpKind = "shl"
	LShr BinOpKind = "lshr"
	AShr BinOpKind = "ashr"
)

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

type FCmpPred string

const (
	FCmpOeq FCmpPred = "oeq"
	FCmpOne FCmpPred = "one"
	FCmpOlt FCmpPred = "olt"
	FCmpOle FCmpPred = "ole"
	FCmpOgt FCmpPred = "ogt"
	FCmpOge FCmpPred = "oge"
)

type Instruction struct {
	Kind InstrKind

	// Common destination (for value-producing instructions).
	Dest *Value

	// BinOp
	BinOp  BinOpKind
	X, Y   *Value

	// Return
	RetVals []*Value

	// Call / CallVoid
	CallName string
	CallArgs []*Value
	CallRet  *TypeDesc
	SretType *TypeDesc

	// Conv
	ConvOp  ConvOp
	ConvSrc *Value
	ConvTo  *TypeDesc

	// Alloca
	AllocaType *TypeDesc
	AllocaAlign int64

	// Load
	LoadSrc   *Value
	LoadAlign int64

	// Store
	StoreSrc   *Value
	StoreDst   *Value
	StoreAlign int64

	// GEP
	GepSrc        *Value
	GepPointee    *TypeDesc
	GepIndices    []*Value
	GepResultType *TypeDesc

	// ICmp
	ICmpPred ICmpPred
	ICmpX    *Value
	ICmpY    *Value

	// FCmp
	FCmpPred FCmpPred
	FCmpX    *Value
	FCmpY    *Value

	// Br
	BrTarget string

	// CondBr
	CondCond  *Value
	CondTrue  string
	CondFalse string

	// Bitcast
	BitcastSrc    *Value
	BitcastTarget *TypeDesc

	// Memcpy
	MemcpyDest *Value
	MemcpySrc  *Value
	MemcpySize *Value
}


type BasicBlock struct {
	Name       string
	Instrs     []Instruction
	Terminator *Instruction
}

func GetBasicBlockName(bb *BasicBlock) string {
	if bb == nil {
		return ""
	}
	return bb.Name
}

func GetBasicBlockInstrs(bb *BasicBlock) []Instruction {
	if bb == nil {
		return nil
	}
	return bb.Instrs
}

func GetBasicBlockTerminator(bb *BasicBlock) *Instruction {
	if bb == nil {
		return nil
	}
	return bb.Terminator
}

func NewBasicBlock(name string) *BasicBlock {
	return &BasicBlock{Name: name}
}

func (bb *BasicBlock) Append(inst Instruction) {
	bb.Instrs = append(bb.Instrs, inst)
}

type Function struct {
	Name    string
	Params  []*Value
	Results []*TypeDesc
	Blocks  []*BasicBlock
	SretType *TypeDesc
}

func GetFunctionName(fn *Function) string {
	if fn == nil {
		return ""
	}
	return fn.Name
}

func GetFunctionParams(fn *Function) []*Value {
	if fn == nil {
		return nil
	}
	return fn.Params
}

func GetFunctionResults(fn *Function) []*TypeDesc {
	if fn == nil {
		return nil
	}
	return fn.Results
}

func GetFunctionBlocks(fn *Function) []*BasicBlock {
	if fn == nil {
		return nil
	}
	return fn.Blocks
}

func GetFunctionSretType(fn *Function) *TypeDesc {
	if fn == nil {
		return nil
	}
	return fn.SretType
}

func (fn *Function) Entry() *BasicBlock {
	if len(fn.Blocks) == 0 {
		entry := NewBasicBlock("entry")
		fn.Blocks = append(fn.Blocks, entry)
	}
	return fn.Blocks[0]
}

type Module struct {
	Functions    []*Function
	Globals      []Global
	TargetTriple string
	DataLayout   string
}

func GetModuleFunctions(mod *Module) []*Function {
	if mod == nil {
		return nil
	}
	return mod.Functions
}

func GetModuleGlobals(mod *Module) []Global {
	if mod == nil {
		return nil
	}
	return mod.Globals
}

func GetModuleTargetTriple(mod *Module) string {
	if mod == nil {
		return ""
	}
	return mod.TargetTriple
}

func GetModuleDataLayout(mod *Module) string {
	if mod == nil {
		return ""
	}
	return mod.DataLayout
}

func (m *Module) AddFunction(fn *Function) {
	m.Functions = append(m.Functions, fn)
}

type Global struct {
	Name    string
	Type    *TypeDesc
	Value   string
	Align   int64
	Private bool
}

func GetGlobalAlign(g *Global) int64 {
	if g == nil {
		return 0
	}
	return g.Align
}

func GetGlobalPrivate(g *Global) bool {
	if g == nil {
		return false
	}
	return g.Private
}

func GetGlobalName(g *Global) string {
	if g == nil {
		return ""
	}
	return g.Name
}

func GetGlobalType(g *Global) *TypeDesc {
	if g == nil {
		return nil
	}
	return g.Type
}

func GetGlobalValue(g *Global) string {
	if g == nil {
		return ""
	}
	return g.Value
}


func formatValue(v *Value) string {
	if v == nil {
		return ""
	}
	if v.kind == ValueRegister || v.kind == ValueParam {
		return "%" + v.name
	}
	if v.kind == ValueConstant && v.Raw != "" {
		if v.ty != nil && v.ty.Kind == KindBasic && v.ty.Basic == "double" {
			// Check manually if v.Raw contains any of ".eE"
			hasAny := false
			chars := ".eE"
			for i := 0; i < len(v.Raw); i = i + 1 {
				for j := 0; j < len(chars); j = j + 1 {
					if v.Raw[i] == chars[j] {
						hasAny = true
						break
					}
				}
				if hasAny {
					break
				}
			}
			if !hasAny {
				return v.Raw + ".0"
			}
		}
		return v.Raw
	}
	return v.name
}

func typeString(v *Value) string {
	if v != nil && v.ty != nil {
		return v.ty.String()
	}
	return "void"
}
