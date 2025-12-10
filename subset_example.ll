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

@g1 = constant i64 1, align 8
@g2 = constant i64 2, align 8
@.str.0 = private unnamed_addr constant [7 x i8] c"\67\6C\6F\62\61\6C\00", align 1
@gstr = constant { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.0, i32 0, i32 0), i64 6 }, align 8
@garr = constant [2 x i64] zeroinitializer, align 8
@gstruct = constant { i64, { i8*, i64 } } zeroinitializer, align 8
@.str.1 = private unnamed_addr constant [4 x i8] c"\68\69\20\00", align 1
@.str.2 = private unnamed_addr constant [6 x i8] c"\74\68\65\72\65\00", align 1
@.str.3 = private unnamed_addr constant [14 x i8] c"\61\66\74\65\72\20\73\74\72\75\63\74\0A\00", align 1
@.str.4 = private unnamed_addr constant [10 x i8] c"\66\6C\6F\61\74\20\6F\6B\0A\00", align 1
@.str.5 = private unnamed_addr constant [13 x i8] c"\61\66\74\65\72\20\66\6C\6F\61\74\0A\00", align 1
@.str.6 = private unnamed_addr constant [5 x i8] c"\65\6E\64\0A\00", align 1

define void @print({ i8*, i64 }* byval({ i8*, i64 }) %s) {
entry:
  %s.addr = alloca { i8*, i64 } , align 8
  call void @gominic_memcpy({ i8*, i64 }* %s.addr, { i8*, i64 }* %s, i64 16)
  %t0 = load { i8*, i64 }, { i8*, i64 }* %s.addr, align 8
  %taddr1 = alloca { i8*, i64, i64 } , align 8
  %t2 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %taddr1, i32 0, i32 0
  %taddr4 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t0, { i8*, i64 }* %taddr4, align 8
  %t3 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr4, i32 0, i32 0
  %t5 = load i8*, i8** %t3, align 8
  store i8* %t5, i8** %t2, align 8
  %t6 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %taddr1, i32 0, i32 1
  %t7 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %taddr1, i32 0, i32 2
  %taddr9 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t0, { i8*, i64 }* %taddr9, align 8
  %t8 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr9, i32 0, i32 1
  %t10 = load i64, i64* %t8, align 8
  store i64 %t10, i64* %t6, align 8
  store i64 %t10, i64* %t7, align 8
  %t11 = load { i8*, i64, i64 }, { i8*, i64, i64 }* %taddr1, align 8
  %b.addr12 = alloca { i8*, i64, i64 }
  %taddr13 = alloca { i8*, i64, i64 } , align 8
  store { i8*, i64, i64 } %t11, { i8*, i64, i64 }* %taddr13
  call void @gominic_memcpy({ i8*, i64, i64 }* %b.addr12, { i8*, i64, i64 }* %taddr13, i64 24)
  %t14 = load { i8*, i64, i64 }, { i8*, i64, i64 }* %b.addr12, align 8
  %taddr15 = alloca { i8*, i64, i64 } , align 8
  store { i8*, i64, i64 } %t14, { i8*, i64, i64 }* %taddr15, align 8
  %t16 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %taddr15, i32 0, i32 1
  %t17 = load i64, i64* %t16, align 8
  %t18 = icmp eq i64 %t17, 0
  %taddr29 = alloca { i8*, i64, i64 } , align 8
  br i1 %t18, label %then0, label %endif1
then0:
  %z.addr19 = alloca i8
  store i8 0, i8* %z.addr19
  call void @gominic_print(i8* %z.addr19, i64 0)
  ret void
endif1:
  %t20 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %b.addr12, i32 0, i32 1
  %t21 = load i64, i64* %t20, align 8
  %t22 = icmp sge i64 0, 0
  %t23 = icmp slt i64 0, %t21
  %t24 = and i1 %t22, %t23
  br i1 %t24, label %idx.ok2, label %idx.fail3
idx.ok2:
  %t25 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %b.addr12, i32 0, i32 0
  %t26 = load i8*, i8** %t25, align 8
  %t27 = getelementptr inbounds i8, i8* %t26, i64 0
  %t28 = load { i8*, i64, i64 }, { i8*, i64, i64 }* %b.addr12, align 8
  store { i8*, i64, i64 } %t28, { i8*, i64, i64 }* %taddr29, align 8
  %t30 = getelementptr inbounds { i8*, i64, i64 }, { i8*, i64, i64 }* %taddr29, i32 0, i32 1
  %t31 = load i64, i64* %t30, align 8
  call void @gominic_print(i8* %t27, i64 %t31)
  ret void
idx.fail3:
  call void @gominic_abort()
  br label %idx.ok2
}

define void @printInt(i64 %v) {
entry:
  %v.addr = alloca i64 , align 8
  store i64 %v, i64* %v.addr
  %t0 = load i64, i64* %v.addr, align 8
  call void @gominic_printInt(i64 %t0)
  ret void
}

define void @println() {
entry:
  call void @gominic_println()
  ret void
}

define i64 @add(i64 %a, i64 %b) {
entry:
  %a.addr = alloca i64 , align 8
  store i64 %a, i64* %a.addr
  %b.addr = alloca i64 , align 8
  store i64 %b, i64* %b.addr
  %t0 = load i64, i64* %a.addr, align 8
  %t1 = load i64, i64* %b.addr, align 8
  %t2 = add i64 %t0, %t1
  ret i64 %t2
}

define void @main() {
entry:
  %a.addr0 = alloca i64
  store i64 1, i64* %a.addr0, align 8
  %b.addr1 = alloca i64
  store i64 2, i64* %b.addr1
  %t2 = load i64, i64* %a.addr0, align 8
  %t3 = load i64, i64* %b.addr1, align 8
  %t4 = call i64 @add(i64 %t2, i64 %t3)
  %c.addr5 = alloca i64
  store i64 %t4, i64* %c.addr5, align 8
  %t6 = load i64, i64* %c.addr5, align 8
  %t7 = icmp slt i64 %t6, 10
  %taddr29 = alloca { i8*, i64 } , align 8
  %taddr30 = alloca { i8*, i64 } , align 8
  %taddr31 = alloca { i8*, i64 } , align 8
  %taddr34 = alloca { i8*, i64 } , align 8
  %taddr36 = alloca { i8*, i64 } , align 8
  %taddr37 = alloca [2 x i64] , align 8
  %taddr42 = alloca [2 x i64] , align 8
  %taddr46 = alloca [0 x i64] , align 8
  %taddr48 = alloca { i64*, i64, i64 } , align 8
  %taddr54 = alloca { i64*, i64, i64 } , align 8
  %taddr56 = alloca { i64*, i64, i64 } , align 8
  %taddr66 = alloca { i64*, i64, i64 } , align 8
  %taddr71 = alloca { i64*, i64, i64 } , align 8
  %taddr73 = alloca { i64*, i64, i64 } , align 8
  %taddr80 = alloca i64 , align 8
  %taddr81 = alloca i64 , align 8
  %taddr85 = alloca i64 , align 8
  %taddr86 = alloca i64 , align 8
  %taddr96 = alloca i64 , align 8
  %taddr97 = alloca i64 , align 8
  %taddr105 = alloca { i64, { i8*, i64 } } , align 8
  %taddr110 = alloca { i64, { i8*, i64 } } , align 8
  %taddr112 = alloca { i8*, i64 } , align 8
  %taddr117 = alloca { i8*, i64 } , align 8
  %taddr118 = alloca { i8*, i64 } , align 8
  %taddr120 = alloca { i8*, i64 } , align 8
  br i1 %t7, label %then0, label %else2
then0:
  %t8 = load i64, i64* %c.addr5, align 8
  %t9 = mul i64 %t8, 2
  store i64 %t9, i64* %c.addr5, align 8
  br label %endif1
endif1:
  %i.addr12 = alloca i64
  store i64 0, i64* %i.addr12, align 8
  br label %for.cond3
else2:
  %t10 = load i64, i64* %c.addr5, align 8
  %t11 = sub i64 %t10, 1
  store i64 %t11, i64* %c.addr5, align 8
  br label %endif1
for.cond3:
  %t13 = load i64, i64* %i.addr12, align 8
  %t14 = icmp slt i64 %t13, 3
  br i1 %t14, label %for.body4, label %for.end6
for.body4:
  %t15 = load i64, i64* %c.addr5, align 8
  %t16 = load i64, i64* %i.addr12, align 8
  %t17 = add i64 %t15, %t16
  store i64 %t17, i64* %c.addr5, align 8
  br label %for.post5
for.post5:
  %t18 = load i64, i64* %i.addr12, align 8
  %t19 = add i64 %t18, 1
  store i64 %t19, i64* %i.addr12, align 8
  br label %for.cond3
for.end6:
  %t20 = load i64, i64* %c.addr5, align 8
  call void @printInt(i64 %t20)
  call void @println()
  br label %for.cond7
for.cond7:
  %t21 = load i64, i64* %c.addr5, align 8
  %t22 = icmp sgt i64 %t21, 0
  br i1 %t22, label %for.body8, label %for.end9
for.body8:
  %t23 = load i64, i64* %c.addr5, align 8
  %t24 = sub i64 %t23, 1
  store i64 %t24, i64* %c.addr5, align 8
  %t25 = load i64, i64* %c.addr5, align 8
  %t26 = icmp eq i64 %t25, 0
  br i1 %t26, label %then10, label %endif11
for.end9:
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.1, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr30, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([6 x i8], [6 x i8]* @.str.2, i32 0, i32 0), i64 5 }, { i8*, i64 }* %taddr31, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr29, { i8*, i64 }* %taddr30, { i8*, i64 }* %taddr31)
  %t32 = load { i8*, i64 }, { i8*, i64 }* %taddr29, align 8
  %msg.addr33 = alloca { i8*, i64 }
  store { i8*, i64 } %t32, { i8*, i64 }* %taddr34
  call void @gominic_memcpy({ i8*, i64 }* %msg.addr33, { i8*, i64 }* %taddr34, i64 16)
  %t35 = load { i8*, i64 }, { i8*, i64 }* %msg.addr33, align 8
  store { i8*, i64 } %t35, { i8*, i64 }* %taddr36
  call void @print({ i8*, i64 }* %taddr36)
  call void @println()
  %t38 = getelementptr inbounds [2 x i64], [2 x i64]* %taddr37, i32 0, i32 0
  store i64 1, i64* %t38, align 8
  %t39 = getelementptr inbounds [2 x i64], [2 x i64]* %taddr37, i32 0, i32 1
  store i64 2, i64* %t39, align 8
  %t40 = load [2 x i64], [2 x i64]* %taddr37, align 8
  %arr.addr41 = alloca [2 x i64]
  store [2 x i64] %t40, [2 x i64]* %taddr42
  call void @gominic_memcpy([2 x i64]* %arr.addr41, [2 x i64]* %taddr42, i64 16)
  %t43 = icmp slt i64 1, 2
  br i1 %t43, label %idx.ok14, label %idx.fail15
then10:
  br label %for.end9
endif11:
  %t27 = load i64, i64* %c.addr5, align 8
  %t28 = icmp slt i64 %t27, 0
  br i1 %t28, label %then12, label %endif13
then12:
  br label %for.cond7
endif13:
  br label %for.cond7
idx.ok14:
  %t44 = getelementptr inbounds [2 x i64], [2 x i64]* %arr.addr41, i32 0, i64 1
  %t45 = load i64, i64* %t44, align 8
  %t47 = getelementptr inbounds [0 x i64], [0 x i64]* %taddr46, i32 0, i32 0
  %t49 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr48, i32 0, i32 0
  store i64* %t47, i64** %t49, align 8
  %t50 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr48, i32 0, i32 1
  store i64 0, i64* %t50, align 8
  %t51 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr48, i32 0, i32 2
  store i64 0, i64* %t51, align 8
  %t52 = load { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr48, align 8
  %sl.addr53 = alloca { i64*, i64, i64 }
  store { i64*, i64, i64 } %t52, { i64*, i64, i64 }* %taddr54
  call void @gominic_memcpy({ i64*, i64, i64 }* %sl.addr53, { i64*, i64, i64 }* %taddr54, i64 24)
  %t55 = load { i64*, i64, i64 }, { i64*, i64, i64 }* %sl.addr53, align 8
  store { i64*, i64, i64 } %t55, { i64*, i64, i64 }* %taddr56, align 8
  %t57 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr56, i32 0, i32 0
  %t58 = load i64*, i64** %t57, align 8
  %t59 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr56, i32 0, i32 1
  %t60 = load i64, i64* %t59, align 8
  %t61 = add i64 %t60, 1
  %t62 = mul i64 %t60, 8
  %t63 = call i8* @gominic_makeSlice(i64 %t61, i64 %t61, i64 8)
  %t64 = bitcast i8* %t63 to i64*
  call void @gominic_memcpy(i64* %t64, i64* %t58, i64 %t62)
  %t65 = getelementptr inbounds i64, i64* %t64, i64 %t60
  store i64 5, i64* %t65, align 8
  %t67 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr66, i32 0, i32 0
  store i64* %t64, i64** %t67, align 8
  %t68 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr66, i32 0, i32 1
  store i64 %t61, i64* %t68, align 8
  %t69 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr66, i32 0, i32 2
  store i64 %t61, i64* %t69, align 8
  %t70 = load { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr66, align 8
  store { i64*, i64, i64 } %t70, { i64*, i64, i64 }* %taddr71
  call void @gominic_memcpy({ i64*, i64, i64 }* %sl.addr53, { i64*, i64, i64 }* %taddr71, i64 24)
  %t72 = load { i64*, i64, i64 }, { i64*, i64, i64 }* %sl.addr53, align 8
  store { i64*, i64, i64 } %t72, { i64*, i64, i64 }* %taddr73, align 8
  %t74 = getelementptr inbounds { i64*, i64, i64 }, { i64*, i64, i64 }* %taddr73, i32 0, i32 1
  %t75 = load i64, i64* %t74, align 8
  %t76 = call i8* @gominic_map_new(i64 8, i64 8, i32 0)
  %m.addr77 = alloca i8*
  store i8* %t76, i8** %m.addr77, align 8
  %t78 = load i64, i64* %c.addr5, align 8
  %t79 = load i8*, i8** %m.addr77, align 8
  store i64 1, i64* %taddr80, align 8
  store i64 %t78, i64* %taddr81, align 8
  %t82 = bitcast i64* %taddr80 to i8*
  %t83 = bitcast i64* %taddr81 to i8*
  call void @gominic_map_set(i8* %t79, i8* %t82, i8* %t83)
  %t84 = load i8*, i8** %m.addr77, align 8
  store i64 1, i64* %taddr85, align 8
  store i64 0, i64* %taddr86, align 8
  %t87 = bitcast i64* %taddr85 to i8*
  %t88 = bitcast i64* %taddr86 to i8*
  %t89 = call i1 @gominic_map_get(i8* %t84, i8* %t87, i8* %t88)
  %t90 = load i64, i64* %taddr86, align 8
  %_.addr91 = alloca i64
  store i64 %t90, i64* %_.addr91, align 8
  %ok.addr92 = alloca i1
  store i1 %t89, i1* %ok.addr92, align 1
  %t93 = load i1, i1* %ok.addr92, align 1
  br i1 %t93, label %then16, label %endif17
idx.fail15:
  call void @gominic_abort()
  unreachable
then16:
  %t94 = load i64, i64* %c.addr5, align 8
  %t95 = load i8*, i8** %m.addr77, align 8
  store i64 1, i64* %taddr96, align 8
  store i64 0, i64* %taddr97, align 8
  %t98 = bitcast i64* %taddr96 to i8*
  %t99 = bitcast i64* %taddr97 to i8*
  %t100 = call i1 @gominic_map_get(i8* %t95, i8* %t98, i8* %t99)
  %t101 = load i64, i64* %taddr97, align 8
  %t102 = add i64 %t94, %t101
  store i64 %t102, i64* %c.addr5, align 8
  br label %endif17
endif17:
  %t103 = load i64, i64* %c.addr5, align 8
  %t104 = load { i8*, i64 }, { i8*, i64 }* %msg.addr33, align 8
  %t106 = getelementptr inbounds { i64, { i8*, i64 } }, { i64, { i8*, i64 } }* %taddr105, i32 0, i32 0
  store i64 %t103, i64* %t106, align 8
  %t107 = getelementptr inbounds { i64, { i8*, i64 } }, { i64, { i8*, i64 } }* %taddr105, i32 0, i32 1
  store { i8*, i64 } %t104, { i8*, i64 }* %t107, align 8
  %t108 = load { i64, { i8*, i64 } }, { i64, { i8*, i64 } }* %taddr105, align 8
  %st.addr109 = alloca { i64, { i8*, i64 } }
  store { i64, { i8*, i64 } } %t108, { i64, { i8*, i64 } }* %taddr110
  call void @gominic_memcpy({ i64, { i8*, i64 } }* %st.addr109, { i64, { i8*, i64 } }* %taddr110, i64 24)
  %t111 = load { i64, { i8*, i64 } }, { i64, { i8*, i64 } }* %st.addr109, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([14 x i8], [14 x i8]* @.str.3, i32 0, i32 0), i64 13 }, { i8*, i64 }* %taddr112
  call void @print({ i8*, i64 }* %taddr112)
  %t113 = fadd double 1.0, 2.5
  %f.addr114 = alloca double
  store double %t113, double* %f.addr114, align 8
  %t115 = load double, double* %f.addr114, align 8
  %t116 = fcmp ogt double %t115, 0.0
  br i1 %t116, label %then18, label %endif19
then18:
  store { i8*, i64 } { i8* getelementptr inbounds ([10 x i8], [10 x i8]* @.str.4, i32 0, i32 0), i64 9 }, { i8*, i64 }* %taddr117
  call void @print({ i8*, i64 }* %taddr117)
  br label %endif19
endif19:
  store { i8*, i64 } { i8* getelementptr inbounds ([13 x i8], [13 x i8]* @.str.5, i32 0, i32 0), i64 12 }, { i8*, i64 }* %taddr118
  call void @print({ i8*, i64 }* %taddr118)
  %t119 = load i64, i64* %c.addr5, align 8
  call void @printInt(i64 %t119)
  call void @println()
  store { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.6, i32 0, i32 0), i64 4 }, { i8*, i64 }* %taddr120
  call void @print({ i8*, i64 }* %taddr120)
  ret void
}

