target triple = "x86_64-pc-windows-msvc"
target datalayout = "e-m:w-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"

declare void @gominic_memcpy(i8*, i8*, i64)
declare void @gominic_abort()
declare i8* @gominic_makeSlice(i64, i64, i64)
declare void @gominic_print(i8*, i64)
declare void @gominic_printInt(i64)
declare void @gominic_println()

declare i1 @gominic_str_eq(i8*, i64, i8*, i64)
declare { i8*, i64 } @gominic_str_concat({ i8*, i64 }, { i8*, i64 })

declare i8* @gominic_map_new(i64, i64, i32)
declare void @gominic_map_set(i8*, i8*, i8*)
declare i1 @gominic_map_get(i8*, i8*, i8*)
declare i64 @gominic_map_len(i8*)

@.str.0 = private unnamed_addr constant [20 x i8] c"\68\65\6C\6C\6F\20\66\72\6F\6D\20\67\6F\6D\69\6E\69\63\0A\00", align 1
@.str.1 = private unnamed_addr constant [8 x i8] c"\65\78\61\6D\70\6C\65\00", align 1

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

define void @main() {
entry:
  %taddr0 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([20 x i8], [20 x i8]* @.str.0, i32 0, i32 0), i64 19 }, { i8*, i64 }* %taddr0
  call void @print({ i8*, i64 }* %taddr0)
  call void @printInt(i64 42)
  call void @println()
  %msg.addr1 = alloca { i8*, i64 }
  %taddr2 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([8 x i8], [8 x i8]* @.str.1, i32 0, i32 0), i64 7 }, { i8*, i64 }* %taddr2
  call void @gominic_memcpy({ i8*, i64 }* %msg.addr1, { i8*, i64 }* %taddr2, i64 16)
  %t3 = load { i8*, i64 }, { i8*, i64 }* %msg.addr1, align 8
  %taddr4 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t3, { i8*, i64 }* %taddr4
  call void @print({ i8*, i64 }* %taddr4)
  call void @println()
  ret void
}

