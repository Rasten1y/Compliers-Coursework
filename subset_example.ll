target triple = "x86_64-pc-windows-msvc"
target datalayout = "e-m:w-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"

declare void @gominic_memcpy(i8*, i8*, i64)
declare void @gominic_abort()
declare i8* @gominic_makeSlice(i64, i64, i64)
declare void @gominic_print(i8*, i64)
declare void @gominic_printInt(i64)
declare void @gominic_println()

declare i1 @gominic_str_eq(i8*, i64, i8*, i64)
declare void @gominic_str_concat({ i8*, i64 }*, { i8*, i64 }*, { i8*, i64 }*)

declare i8* @gominic_map_new(i64, i64, i32)
declare void @gominic_map_set(i8*, i8*, i8*)
declare i1 @gominic_map_get(i8*, i8*, i8*)
declare i64 @gominic_map_len(i8*)

declare i64 @gominic_argc()
declare void @gominic_argv({ i8*, i64 }* sret({ i8*, i64 }), i64)
declare i1 @gominic_write_file(i8*, i64, i8*, i64)

declare void @flag.Usage()
declare void @Usage()
declare void @flag.Parse()
declare { { i8*, i64 }*, i64, i64 } @flag.Args()
declare void @flag.BoolVar(i1*, { i8*, i64 }, i1, { i8*, i64 })
declare void @flag.StringVar({ i8*, i64 }*, { i8*, i64 }, { i8*, i64 }, { i8*, i64 })
declare { i8*, i64 } @error.Error(i8*)
declare i8* @errors.New({ i8*, i64 })
declare { { i8*, i64, i64 }, i8* } @exec.Cmd.CombinedOutput(i8*)
declare { { i8*, i64 }, { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }, { i8*, i64 }, i8*, i8*, i8*, { i8***, i64, i64 }, i8**, i8**, i8**, i8*, i8*, i8*, i8*, { i8**, i64, i64 }, { i8**, i64, i64 }, { i8**, i64, i64 }, i8*, i8*, { i8*, i64, i64 }, i8*, { { i8*, i64 }, { i8*, i64 } } }* @exec.Command({ i8*, i64 }, { { i8*, i64 }*, i64, i64 })
declare { { { i8**, { i8***, i64, i64 }, i8**, i8**, { i8*, i64 } }**, i64, i64 }, i8* } @frontend.ParseAndCheckAll({ { i8*, i64 }*, i64, i64 })
declare { i8**, i8* } @frontend.BuildModule({ { i8**, { i8***, i64, i64 }, i8**, i8**, { i8*, i64 } }**, i64, i64 })
declare void @frontend.SetSkipSubsetCheck(i1)
declare i8* @os.WriteFile({ i8*, i64 }, { i8*, i64, i64 }, i64)
declare void @os.File.WriteString(i8**, { i8*, i64 })
declare void @os.Exit(i64)
declare { i8*, i64 } @types.Package.Name(i8**)
@os.Stderr = external global i8*
@os.Stdout = external global i8*
@os.Stdin = external global i8*
@gominic.I64 = external global i8***
@gominic.I32 = external global i8***
@gominic.I8 = external global i8***
@gominic.I1 = external global i8***
@gominic.F64 = external global i8***
@gominic.Void = external global i8***
@gominic.PtrI8 = external global i8***
@gominic.false = external global i1
@gominic.true = external global i1
@gominic.nil = external global i8**
@ir.bb = external global i8***
@ir.m = external global i8***
@ir.t = external global i8***
@ir.v = external global i8***
@ir.fn = external global i8***
@ir.nil = external global i8***
@ir.false = external global i1
@ir.true = external global i1

@Stderr = external global i8*
@Stdout = external global i8*
@Stdin = external global i8*

@g1 = constant i64 1, align 8
@g2 = constant i64 2, align 8
@.str.1 = private unnamed_addr constant [7 x i8] c"\67\6C\6F\62\61\6C\00", align 1
@gstr = constant { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.1, i32 0, i32 0), i64 6 }, align 8
@garr = constant [2 x i64] zeroinitializer, align 8
@gstruct = constant { i64, { i8*, i64 } } zeroinitializer, align 8
@.str.2 = private unnamed_addr constant [4 x i8] c"\68\69\20\00", align 1
@.str.3 = private unnamed_addr constant [6 x i8] c"\74\68\65\72\65\00", align 1
@.str.4 = private unnamed_addr constant [14 x i8] c"\61\66\74\65\72\20\73\74\72\75\63\74\0A\00", align 1
@.str.5 = private unnamed_addr constant [10 x i8] c"\66\6C\6F\61\74\20\6F\6B\0A\00", align 1
@.str.6 = private unnamed_addr constant [13 x i8] c"\61\66\74\65\72\20\66\6C\6F\61\74\0A\00", align 1
@.str.7 = private unnamed_addr constant [5 x i8] c"\65\6E\64\0A\00", align 1

define void @print({ i8*, i64 } %s) {
entry:
  %p0.addr = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %s, { i8*, i64 }* %p0.addr, align 8
  %t1 = load { i8*, i64 }, { i8*, i64 }* %p0.addr, align 8
  %t2 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t1, { i8*, i64 }* %t2, align 8
  %t3 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %t2, i32 0, i32 0
  %t4 = load i8*, i8** %t3, align 8
  %t5 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %t2, i32 0, i32 1
  %t6 = load i64, i64* %t5, align 8
  %t7 = call i8* @gominic_makeSlice(i64 %t6, i64 %t6, i64 1)
  call void @gominic_memcpy(i8* %t7, i8* %t4, i64 %t6)
  %t8 = alloca { i8*, i64, i64 } , align 8
  %t9 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %t8, i32 0, i32 0
  store i8* %t7, i8** %t9, align 8
  %t10 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %t8, i32 0, i32 1
  store i64 %t6, i64* %t10, align 8
  %t11 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %t8, i32 0, i32 2
  store i64 %t6, i64* %t11, align 8
  %t12 = load { i8*, i64, i64 }, { i8*, i64, i64 }* %t8, align 8
  %b.addr1 = alloca { i8*, i64, i64 } , align 8
  store { i8*, i64, i64 } %t12, { i8*, i64, i64 }* %b.addr1, align 8
  %t13 = load { i8*, i64, i64 }, { i8*, i64, i64 }* %b.addr1, align 8
  %t14 = alloca { i8*, i64, i64 } , align 8
  store { i8*, i64, i64 } %t13, { i8*, i64, i64 }* %t14, align 8
  %t15 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %t14, i32 0, i32 1
  %t16 = load i64, i64* %t15, align 8
  %t17 = icmp eq i64 %t16, 0
  br i1 %t17, label %then1, label %else2
then1:
  %z.addr2 = alloca i8 , align 1
  call void @gominic_print(i8* %z.addr2, i64 0)
  ret void
else2:
  br label %endif3
endif3:
  %t18 = load { i8*, i64, i64 }, { i8*, i64, i64 }* %b.addr1, align 8
  %t19 = alloca { i8*, i64, i64 } , align 8
  store { i8*, i64, i64 } %t18, { i8*, i64, i64 }* %t19, align 8
  %t20 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %t19, i32 0, i32 0
  %t21 = load i8*, i8* %t20, align 8
  %t22 = getelementptr inbounds i8, i8* %t21, i64 0
  %t23 = load { i8*, i64, i64 }, { i8*, i64, i64 }* %b.addr1, align 8
  %t24 = alloca { i8*, i64, i64 } , align 8
  store { i8*, i64, i64 } %t23, { i8*, i64, i64 }* %t24, align 8
  %t25 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %t24, i32 0, i32 1
  %t26 = load i64, i64* %t25, align 8
  call void @gominic_print(i8* %t22, i64 %t26)
  ret void
}

define void @printInt(i64 %v) {
entry:
  %p0.addr = alloca i64 , align 8
  store i64 %v, i64* %p0.addr, align 8
  %t1 = load i64, i64* %p0.addr, align 8
  call void @gominic_printInt(i64 %t1)
  ret void
}

define void @println() {
entry:
  call void @gominic_println()
  ret void
}

define i64 @add(i64 %a, i64 %b) {
entry:
  %p0.addr = alloca i64 , align 8
  store i64 %a, i64* %p0.addr, align 8
  %p1.addr = alloca i64 , align 8
  store i64 %b, i64* %p1.addr, align 8
  %t1 = load i64, i64* %p0.addr, align 8
  %t2 = load i64, i64* %p1.addr, align 8
  %t3 = add i64 %t1, %t2
  ret i64 %t3
}

define void @main() {
entry:
  %a.addr3 = alloca i64 , align 8
  store i64 1, i64* %a.addr3, align 8
  %b.addr4 = alloca i64 , align 8
  store i64 2, i64* %b.addr4, align 8
  %t1 = load i64, i64* %a.addr3, align 8
  %t2 = load i64, i64* %b.addr4, align 8
  %t3 = call i64 @add(i64 %t1, i64 %t2)
  %c.addr5 = alloca i64 , align 8
  store i64 %t3, i64* %c.addr5, align 8
  %t4 = load i64, i64* %c.addr5, align 8
  %t5 = icmp slt i64 %t4, 10
  br i1 %t5, label %then1, label %else2
then1:
  %t6 = load i64, i64* %c.addr5, align 8
  %t7 = mul i64 %t6, 2
  store i64 %t7, i64* %c.addr5, align 8
  br label %endif3
else2:
  %t8 = load i64, i64* %c.addr5, align 8
  %t9 = sub i64 %t8, 1
  store i64 %t9, i64* %c.addr5, align 8
  br label %endif3
endif3:
  %i.addr6 = alloca i64 , align 8
  store i64 0, i64* %i.addr6, align 8
  br label %for.cond4
for.cond4:
  %t10 = load i64, i64* %i.addr6, align 8
  %t11 = icmp slt i64 %t10, 3
  br i1 %t11, label %for.body5, label %for.end7
for.body5:
  %t12 = load i64, i64* %c.addr5, align 8
  %t13 = load i64, i64* %i.addr6, align 8
  %t14 = add i64 %t12, %t13
  store i64 %t14, i64* %c.addr5, align 8
  br label %for.post6
for.post6:
  %t15 = load i64, i64* %i.addr6, align 8
  %t16 = add i64 %t15, 1
  store i64 %t16, i64* %i.addr6, align 8
  br label %for.cond4
for.end7:
  %t17 = load i64, i64* %c.addr5, align 8
  call void @printInt(i64 %t17)
  call void @println()
  br label %for.cond8
for.cond8:
  %t18 = load i64, i64* %c.addr5, align 8
  %t19 = icmp sgt i64 %t18, 0
  br i1 %t19, label %for.body9, label %for.end11
for.body9:
  %t20 = load i64, i64* %c.addr5, align 8
  %t21 = sub i64 %t20, 1
  store i64 %t21, i64* %c.addr5, align 8
  %t22 = load i64, i64* %c.addr5, align 8
  %t23 = icmp eq i64 %t22, 0
  br i1 %t23, label %then12, label %else13
for.post10:
  unreachable
for.end11:
  %t26 = alloca { i8*, i64 } , align 8
  %t27 = alloca { i8*, i64 } , align 8
  %t28 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.2, i32 0, i32 0), i64 3 }, { i8*, i64 }* %t27, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([6 x i8], [6 x i8]* @.str.3, i32 0, i32 0), i64 5 }, { i8*, i64 }* %t28, align 8
  call void @gominic_str_concat({ i8*, i64 }* %t26, { i8*, i64 }* %t27, { i8*, i64 }* %t28)
  %t29 = load { i8*, i64 }, { i8*, i64 }* %t26, align 8
  %msg.addr7 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t29, { i8*, i64 }* %msg.addr7, align 8
  %t30 = load { i8*, i64 }, { i8*, i64 }* %msg.addr7, align 8
  call void @print({ i8*, i64 } %t30)
  call void @println()
  %t31 = alloca [2 x i64] , align 8
  %t32 = getelementptr inbounds [2 x i64], [2 x i64]* %t31, i32 0, i32 0
  store i64 1, i64* %t32, align 8
  %t33 = getelementptr inbounds [2 x i64], [2 x i64]* %t31, i32 0, i32 1
  store i64 2, i64* %t33, align 8
  %t34 = load [2 x i64], [2 x i64]* %t31, align 8
  %arr.addr8 = alloca [2 x i64] , align 8
  store [2 x i64] %t34, [2 x i64]* %arr.addr8, align 8
  %t35 = load [2 x i64], [2 x i64]* %arr.addr8, align 8
  %t36 = alloca [2 x i64] , align 8
  store [2 x i64] %t35, [2 x i64]* %t36, align 8
  %t37 = getelementptr inbounds [2 x i64], [2 x i64]* %t36, i32 0, i64 1
  %t38 = load i64, i64* %t37, align 8
  %t39 = call i64* @gominic_makeSlice(i64 0, i64 0, i64 8)
  %t40 = alloca { i64*, i64, i64 } , align 8
  %t41 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %t40, i32 0, i32 0
  store i64* %t39, i64** %t41, align 8
  %t42 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %t40, i32 0, i32 1
  store i64 0, i64* %t42, align 8
  %t43 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %t40, i32 0, i32 2
  store i64 0, i64* %t43, align 8
  %t44 = load { i64*, i64, i64 }, { i64*, i64, i64 }* %t40, align 8
  %sl.addr9 = alloca { i64*, i64, i64 } , align 8
  store { i64*, i64, i64 } %t44, { i64*, i64, i64 }* %sl.addr9, align 8
  %t45 = load { i64*, i64, i64 }, { i64*, i64, i64 }* %sl.addr9, align 8
  %t46 = alloca { i64*, i64, i64 } , align 8
  store { i64*, i64, i64 } %t45, { i64*, i64, i64 }* %t46, align 8
  %t47 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %t46, i32 0, i32 1
  %t48 = load i64, i64* %t47, align 8
  %t49 = add i64 %t48, 1
  %t50 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %t46, i32 0, i32 2
  %t51 = load i64, i64* %t50, align 8
  %t52 = call i64* @gominic_makeSlice(i64 %t49, i64 %t49, i64 8)
  %t53 = getelementptr inbounds i64, i64* %t52, i64 0
  store i64 5, i64* %t53, align 8
  %t54 = alloca { i64*, i64, i64 } , align 8
  %t55 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %t54, i32 0, i32 0
  store i64* %t52, i64** %t55, align 8
  %t56 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %t54, i32 0, i32 1
  store i64 %t49, i64* %t56, align 8
  %t57 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %t54, i32 0, i32 2
  store i64 %t49, i64* %t57, align 8
  %t58 = load { i64*, i64, i64 }, { i64*, i64, i64 }* %t54, align 8
  store { i64*, i64, i64 } %t58, { i64*, i64, i64 }* %sl.addr9, align 8
  %t59 = load { i64*, i64, i64 }, { i64*, i64, i64 }* %sl.addr9, align 8
  %t60 = alloca { i64*, i64, i64 } , align 8
  store { i64*, i64, i64 } %t59, { i64*, i64, i64 }* %t60, align 8
  %t61 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %t60, i32 0, i32 1
  %t62 = load i64, i64* %t61, align 8
  %t63 = call i8* @gominic_map_new(i64 8, i64 8, i32 0)
  %m.addr10 = alloca i8* , align 8
  store i8* %t63, i8** %m.addr10, align 8
  %t64 = load i8*, i8** %m.addr10, align 8
  %t65 = load i64, i64* %c.addr5, align 8
  %t66 = alloca i64 , align 8
  store i64 1, i64* %t66, align 8
  %t67 = alloca i64 , align 8
  store i64 %t65, i64* %t67, align 8
  %t68 = bitcast i64* %t66 to i8*
  %t69 = bitcast i64* %t67 to i8*
  call void @gominic_map_set(i8* %t64, i8* %t68, i8* %t69)
  %t70 = load i8*, i8** %m.addr10, align 8
  %t71 = alloca i64 , align 8
  store i64 1, i64* %t71, align 8
  %t72 = alloca i64 , align 8
  %t73 = bitcast i64* %t71 to i8*
  %t74 = bitcast i64* %t72 to i8*
  %t75 = call i1 @gominic_map_get(i8* %t70, i8* %t73, i8* %t74)
  %t76 = load i64, i64* %t72, align 8
  %ok.addr11 = alloca i1 , align 1
  store i1 %t75, i1* %ok.addr11, align 1
  %t77 = load i1, i1* %ok.addr11, align 1
  br i1 %t77, label %then18, label %else19
then12:
  br label %for.end11
else13:
  br label %endif14
endif14:
  %t24 = load i64, i64* %c.addr5, align 8
  %t25 = icmp slt i64 %t24, 0
  br i1 %t25, label %then15, label %else16
then15:
  br label %for.cond8
else16:
  br label %endif17
endif17:
  br label %for.cond8
then18:
  %t78 = load i64, i64* %c.addr5, align 8
  %t79 = load i8*, i8** %m.addr10, align 8
  %t80 = alloca i64 , align 8
  store i64 1, i64* %t80, align 8
  %t81 = alloca i64 , align 8
  %t82 = bitcast i64* %t80 to i8*
  %t83 = bitcast i64* %t81 to i8*
  call i1 @gominic_map_get(i8* %t79, i8* %t82, i8* %t83)
  %t84 = load i64, i64* %t81, align 8
  %t85 = add i64 %t78, %t84
  store i64 %t85, i64* %c.addr5, align 8
  br label %endif20
else19:
  br label %endif20
endif20:
  %t86 = alloca { i64, { i8*, i64 } } , align 8
  %t87 = load i64, i64* %c.addr5, align 8
  %t88 = getelementptr inbounds { i64, { i8*, i64 } }, { i64, { i8*, i64 } }* %t86, i32 0, i32 0
  store i64 %t87, i64* %t88, align 8
  %t89 = load { i8*, i64 }, { i8*, i64 }* %msg.addr7, align 8
  %t90 = getelementptr inbounds { i64, { i8*, i64 } }, { i64, { i8*, i64 } }* %t86, i32 0, i32 1
  store { i8*, i64 } %t89, { i8*, i64 }* %t90, align 8
  %t91 = load { i64, { i8*, i64 } }, { i64, { i8*, i64 } }* %t86, align 8
  %st.addr12 = alloca { i64, { i8*, i64 } } , align 8
  store { i64, { i8*, i64 } } %t91, { i64, { i8*, i64 } }* %st.addr12, align 8
  %t92 = load { i64, { i8*, i64 } }, { i64, { i8*, i64 } }* %st.addr12, align 8
  call void @print({ i8*, i64 } { i8* getelementptr inbounds ([14 x i8], [14 x i8]* @.str.4, i32 0, i32 0), i64 13 })
  %t93 = sitofp i64 1 to double
  %t94 = fadd double %t93, 2.5
  %f.addr13 = alloca double , align 8
  store double %t94, double* %f.addr13, align 8
  %t95 = load double, double* %f.addr13, align 8
  %t96 = sitofp i64 0 to double
  %t97 = fcmp ogt double %t95, %t96
  br i1 %t97, label %then21, label %else22
then21:
  call void @print({ i8*, i64 } { i8* getelementptr inbounds ([10 x i8], [10 x i8]* @.str.5, i32 0, i32 0), i64 9 })
  br label %endif23
else22:
  br label %endif23
endif23:
  call void @print({ i8*, i64 } { i8* getelementptr inbounds ([13 x i8], [13 x i8]* @.str.6, i32 0, i32 0), i64 12 })
  %t98 = load i64, i64* %c.addr5, align 8
  call void @printInt(i64 %t98)
  call void @println()
  call void @print({ i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.7, i32 0, i32 0), i64 4 })
  ret void
}

