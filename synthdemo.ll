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

declare i64 @gominic_argc()
declare { i8*, i64 } @gominic_argv(i64)
declare i1 @gominic_write_file(i8*, i64, i8*, i64)

@.str.demo = private unnamed_addr constant [16 x i8] c"\68\65\6C\6C\6F\20\66\72\6F\6D\20\64\65\6D\6F\00", align 1

define void @main() {
entry:
  call void @gominic_print(i8* getelementptr inbounds ([16 x i8], [16 x i8]* @.str.demo, i32 0, i32 0), i64 15)
  call void @gominic_printInt(i64 42)
  call void @gominic_println()
  ret void
}

