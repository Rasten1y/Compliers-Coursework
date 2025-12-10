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

@.str.0 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.1 = private unnamed_addr constant [45 x i8] c"\64\65\63\6C\61\72\65\20\76\6F\69\64\20\40\67\6F\6D\69\6E\69\63\5F\6D\65\6D\63\70\79\28\69\38\2A\2C\20\69\38\2A\2C\20\69\36\34\29\0A\00", align 1
@.str.2 = private unnamed_addr constant [31 x i8] c"\64\65\63\6C\61\72\65\20\76\6F\69\64\20\40\67\6F\6D\69\6E\69\63\5F\61\62\6F\72\74\28\29\0A\00", align 1
@.str.3 = private unnamed_addr constant [47 x i8] c"\64\65\63\6C\61\72\65\20\69\38\2A\20\40\67\6F\6D\69\6E\69\63\5F\6D\61\6B\65\53\6C\69\63\65\28\69\36\34\2C\20\69\36\34\2C\20\69\36\34\29\0A\00", align 1
@.str.4 = private unnamed_addr constant [39 x i8] c"\64\65\63\6C\61\72\65\20\76\6F\69\64\20\40\67\6F\6D\69\6E\69\63\5F\70\72\69\6E\74\28\69\38\2A\2C\20\69\36\34\29\0A\00", align 1
@.str.5 = private unnamed_addr constant [37 x i8] c"\64\65\63\6C\61\72\65\20\76\6F\69\64\20\40\67\6F\6D\69\6E\69\63\5F\70\72\69\6E\74\49\6E\74\28\69\36\34\29\0A\00", align 1
@.str.6 = private unnamed_addr constant [34 x i8] c"\64\65\63\6C\61\72\65\20\76\6F\69\64\20\40\67\6F\6D\69\6E\69\63\5F\70\72\69\6E\74\6C\6E\28\29\0A\0A\00", align 1
@.str.7 = private unnamed_addr constant [48 x i8] c"\64\65\63\6C\61\72\65\20\69\31\20\40\67\6F\6D\69\6E\69\63\5F\73\74\72\5F\65\71\28\69\38\2A\2C\20\69\36\34\2C\20\69\38\2A\2C\20\69\36\34\29\0A\00", align 1
@.str.8 = private unnamed_addr constant [80 x i8] c"\64\65\63\6C\61\72\65\20\76\6F\69\64\20\40\67\6F\6D\69\6E\69\63\5F\73\74\72\5F\63\6F\6E\63\61\74\28\7B\20\69\38\2A\2C\20\69\36\34\20\7D\2A\2C\20\7B\20\69\38\2A\2C\20\69\36\34\20\7D\2A\2C\20\7B\20\69\38\2A\2C\20\69\36\34\20\7D\2A\29\0A\0A\00", align 1
@.str.9 = private unnamed_addr constant [45 x i8] c"\64\65\63\6C\61\72\65\20\69\38\2A\20\40\67\6F\6D\69\6E\69\63\5F\6D\61\70\5F\6E\65\77\28\69\36\34\2C\20\69\36\34\2C\20\69\33\32\29\0A\00", align 1
@.str.10 = private unnamed_addr constant [46 x i8] c"\64\65\63\6C\61\72\65\20\76\6F\69\64\20\40\67\6F\6D\69\6E\69\63\5F\6D\61\70\5F\73\65\74\28\69\38\2A\2C\20\69\38\2A\2C\20\69\38\2A\29\0A\00", align 1
@.str.11 = private unnamed_addr constant [44 x i8] c"\64\65\63\6C\61\72\65\20\69\31\20\40\67\6F\6D\69\6E\69\63\5F\6D\61\70\5F\67\65\74\28\69\38\2A\2C\20\69\38\2A\2C\20\69\38\2A\29\0A\00", align 1
@.str.12 = private unnamed_addr constant [36 x i8] c"\64\65\63\6C\61\72\65\20\69\36\34\20\40\67\6F\6D\69\6E\69\63\5F\6D\61\70\5F\6C\65\6E\28\69\38\2A\29\0A\0A\00", align 1
@.str.13 = private unnamed_addr constant [29 x i8] c"\64\65\63\6C\61\72\65\20\69\36\34\20\40\67\6F\6D\69\6E\69\63\5F\61\72\67\63\28\29\0A\00", align 1
@.str.14 = private unnamed_addr constant [67 x i8] c"\64\65\63\6C\61\72\65\20\76\6F\69\64\20\40\67\6F\6D\69\6E\69\63\5F\61\72\67\76\28\7B\20\69\38\2A\2C\20\69\36\34\20\7D\2A\20\73\72\65\74\28\7B\20\69\38\2A\2C\20\69\36\34\20\7D\29\2C\20\69\36\34\29\0A\00", align 1
@.str.15 = private unnamed_addr constant [41 x i8] c"\64\65\63\6C\61\72\65\20\7B\20\69\38\2A\2C\20\69\36\34\20\7D\20\40\67\6F\6D\69\6E\69\63\5F\61\72\67\76\28\69\36\34\29\0A\00", align 1
@.str.16 = private unnamed_addr constant [53 x i8] c"\64\65\63\6C\61\72\65\20\69\31\20\40\67\6F\6D\69\6E\69\63\5F\77\72\69\74\65\5F\66\69\6C\65\28\69\38\2A\2C\20\69\36\34\2C\20\69\38\2A\2C\20\69\36\34\29\0A\0A\00", align 1
@.str.17 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.18 = private unnamed_addr constant [22 x i8] c"\70\72\69\76\61\74\65\20\75\6E\6E\61\6D\65\64\5F\61\64\64\72\20\00", align 1
@.str.19 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.20 = private unnamed_addr constant [9 x i8] c"\2C\20\61\6C\69\67\6E\20\00", align 1
@.str.21 = private unnamed_addr constant [2 x i8] c"\40\00", align 1
@.str.22 = private unnamed_addr constant [4 x i8] c"\20\3D\20\00", align 1
@.str.23 = private unnamed_addr constant [10 x i8] c"\63\6F\6E\73\74\61\6E\74\20\00", align 1
@.str.24 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.25 = private unnamed_addr constant [2 x i8] c"\0A\00", align 1
@.str.26 = private unnamed_addr constant [2 x i8] c"\0A\00", align 1
@.str.27 = private unnamed_addr constant [6 x i8] c"\73\72\65\74\28\00", align 1
@.str.28 = private unnamed_addr constant [3 x i8] c"\29\20\00", align 1
@.str.29 = private unnamed_addr constant [3 x i8] c"\20\25\00", align 1
@.str.30 = private unnamed_addr constant [8 x i8] c"\20\62\79\76\61\6C\28\00", align 1
@.str.31 = private unnamed_addr constant [4 x i8] c"\29\20\25\00", align 1
@.str.32 = private unnamed_addr constant [3 x i8] c"\20\25\00", align 1
@.str.33 = private unnamed_addr constant [8 x i8] c"\64\65\66\69\6E\65\20\00", align 1
@.str.34 = private unnamed_addr constant [3 x i8] c"\20\40\00", align 1
@.str.35 = private unnamed_addr constant [2 x i8] c"\28\00", align 1
@.str.36 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.37 = private unnamed_addr constant [5 x i8] c"\29\20\7B\0A\00", align 1
@.str.38 = private unnamed_addr constant [4 x i8] c"\7D\0A\0A\00", align 1
@.str.39 = private unnamed_addr constant [3 x i8] c"\3A\0A\00", align 1
@.str.40 = private unnamed_addr constant [3 x i8] c"\20\20\00", align 1
@.str.41 = private unnamed_addr constant [2 x i8] c"\0A\00", align 1
@.str.42 = private unnamed_addr constant [3 x i8] c"\20\20\00", align 1
@.str.43 = private unnamed_addr constant [2 x i8] c"\0A\00", align 1
@.str.44 = private unnamed_addr constant [15 x i8] c"\20\20\75\6E\72\65\61\63\68\61\62\6C\65\0A\00", align 1
@.str.45 = private unnamed_addr constant [5 x i8] c"\76\6F\69\64\00", align 1
@.str.46 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.47 = private unnamed_addr constant [20 x i8] c"\78\38\36\5F\36\34\2D\70\63\2D\6C\69\6E\75\78\2D\67\6E\75\00", align 1
@.str.48 = private unnamed_addr constant [13 x i8] c"\77\69\6E\64\6F\77\73\2D\6D\73\76\63\00", align 1
@.str.49 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.50 = private unnamed_addr constant [71 x i8] c"\65\2D\6D\3A\65\2D\70\32\37\30\3A\33\32\3A\33\32\2D\70\32\37\31\3A\33\32\3A\33\32\2D\70\32\37\32\3A\36\34\3A\36\34\2D\69\36\34\3A\36\34\2D\66\38\30\3A\31\32\38\2D\6E\38\3A\31\36\3A\33\32\3A\36\34\2D\53\31\32\38\00", align 1
@.str.51 = private unnamed_addr constant [18 x i8] c"\74\61\72\67\65\74\20\74\72\69\70\6C\65\20\3D\20\22\00", align 1
@.str.52 = private unnamed_addr constant [3 x i8] c"\22\0A\00", align 1
@.str.53 = private unnamed_addr constant [22 x i8] c"\74\61\72\67\65\74\20\64\61\74\61\6C\61\79\6F\75\74\20\3D\20\22\00", align 1
@.str.54 = private unnamed_addr constant [4 x i8] c"\22\0A\0A\00", align 1
@.str.55 = private unnamed_addr constant [3 x i8] c"\69\31\00", align 1
@ir.I1 = constant { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.55, i32 0, i32 0), i64 2 }, align 8
@.str.56 = private unnamed_addr constant [3 x i8] c"\69\38\00", align 1
@ir.I8 = constant { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.56, i32 0, i32 0), i64 2 }, align 8
@.str.57 = private unnamed_addr constant [4 x i8] c"\69\33\32\00", align 1
@ir.I32 = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.57, i32 0, i32 0), i64 3 }, align 8
@.str.58 = private unnamed_addr constant [4 x i8] c"\69\36\34\00", align 1
@ir.I64 = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.58, i32 0, i32 0), i64 3 }, align 8
@.str.59 = private unnamed_addr constant [7 x i8] c"\64\6F\75\62\6C\65\00", align 1
@ir.F64 = constant { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.59, i32 0, i32 0), i64 6 }, align 8
@.str.60 = private unnamed_addr constant [5 x i8] c"\76\6F\69\64\00", align 1
@ir.Void = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.60, i32 0, i32 0), i64 4 }, align 8
@.str.61 = private unnamed_addr constant [4 x i8] c"\69\38\2A\00", align 1
@ir.PtrI8 = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.61, i32 0, i32 0), i64 3 }, align 8
@ir.ValueInvalid = constant i64 3, align 8
@ir.ValueParam = constant i64 0, align 8
@ir.ValueRegister = constant i64 0, align 8
@ir.ValueConstant = constant i64 0, align 8
@.str.62 = private unnamed_addr constant [4 x i8] c"\61\64\64\00", align 1
@ir.Add = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.62, i32 0, i32 0), i64 3 }, align 8
@.str.63 = private unnamed_addr constant [4 x i8] c"\73\75\62\00", align 1
@ir.Sub = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.63, i32 0, i32 0), i64 3 }, align 8
@.str.64 = private unnamed_addr constant [4 x i8] c"\6D\75\6C\00", align 1
@ir.Mul = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.64, i32 0, i32 0), i64 3 }, align 8
@.str.65 = private unnamed_addr constant [5 x i8] c"\73\64\69\76\00", align 1
@ir.SDiv = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.65, i32 0, i32 0), i64 4 }, align 8
@.str.66 = private unnamed_addr constant [5 x i8] c"\75\64\69\76\00", align 1
@ir.UDiv = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.66, i32 0, i32 0), i64 4 }, align 8
@.str.67 = private unnamed_addr constant [5 x i8] c"\73\72\65\6D\00", align 1
@ir.SRem = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.67, i32 0, i32 0), i64 4 }, align 8
@.str.68 = private unnamed_addr constant [5 x i8] c"\75\72\65\6D\00", align 1
@ir.URem = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.68, i32 0, i32 0), i64 4 }, align 8
@.str.69 = private unnamed_addr constant [5 x i8] c"\66\61\64\64\00", align 1
@ir.FAdd = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.69, i32 0, i32 0), i64 4 }, align 8
@.str.70 = private unnamed_addr constant [5 x i8] c"\66\73\75\62\00", align 1
@ir.FSub = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.70, i32 0, i32 0), i64 4 }, align 8
@.str.71 = private unnamed_addr constant [5 x i8] c"\66\6D\75\6C\00", align 1
@ir.FMul = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.71, i32 0, i32 0), i64 4 }, align 8
@.str.72 = private unnamed_addr constant [5 x i8] c"\66\64\69\76\00", align 1
@ir.FDiv = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.72, i32 0, i32 0), i64 4 }, align 8
@.str.73 = private unnamed_addr constant [4 x i8] c"\61\6E\64\00", align 1
@ir.And = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.73, i32 0, i32 0), i64 3 }, align 8
@.str.74 = private unnamed_addr constant [3 x i8] c"\6F\72\00", align 1
@ir.Or = constant { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.74, i32 0, i32 0), i64 2 }, align 8
@.str.75 = private unnamed_addr constant [6 x i8] c"\74\72\75\6E\63\00", align 1
@ir.Trunc = constant { i8*, i64 } { i8* getelementptr inbounds ([6 x i8], [6 x i8]* @.str.75, i32 0, i32 0), i64 5 }, align 8
@.str.76 = private unnamed_addr constant [5 x i8] c"\7A\65\78\74\00", align 1
@ir.ZExt = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.76, i32 0, i32 0), i64 4 }, align 8
@.str.77 = private unnamed_addr constant [5 x i8] c"\73\65\78\74\00", align 1
@ir.SExt = constant { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.77, i32 0, i32 0), i64 4 }, align 8
@.str.78 = private unnamed_addr constant [7 x i8] c"\73\69\74\6F\66\70\00", align 1
@ir.SIToFP = constant { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.78, i32 0, i32 0), i64 6 }, align 8
@.str.79 = private unnamed_addr constant [7 x i8] c"\75\69\74\6F\66\70\00", align 1
@ir.UIToFP = constant { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.79, i32 0, i32 0), i64 6 }, align 8
@.str.80 = private unnamed_addr constant [7 x i8] c"\66\70\74\6F\73\69\00", align 1
@ir.FPToSI = constant { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.80, i32 0, i32 0), i64 6 }, align 8
@.str.81 = private unnamed_addr constant [7 x i8] c"\66\70\74\6F\75\69\00", align 1
@ir.FPToUI = constant { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.81, i32 0, i32 0), i64 6 }, align 8
@.str.82 = private unnamed_addr constant [3 x i8] c"\65\71\00", align 1
@ir.ICmpEq = constant { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.82, i32 0, i32 0), i64 2 }, align 8
@.str.83 = private unnamed_addr constant [3 x i8] c"\6E\65\00", align 1
@ir.ICmpNe = constant { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.83, i32 0, i32 0), i64 2 }, align 8
@.str.84 = private unnamed_addr constant [4 x i8] c"\73\6C\74\00", align 1
@ir.ICmpSlt = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.84, i32 0, i32 0), i64 3 }, align 8
@.str.85 = private unnamed_addr constant [4 x i8] c"\73\6C\65\00", align 1
@ir.ICmpSle = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.85, i32 0, i32 0), i64 3 }, align 8
@.str.86 = private unnamed_addr constant [4 x i8] c"\73\67\74\00", align 1
@ir.ICmpSgt = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.86, i32 0, i32 0), i64 3 }, align 8
@.str.87 = private unnamed_addr constant [4 x i8] c"\73\67\65\00", align 1
@ir.ICmpSge = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.87, i32 0, i32 0), i64 3 }, align 8
@.str.88 = private unnamed_addr constant [4 x i8] c"\75\6C\74\00", align 1
@ir.ICmpUlt = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.88, i32 0, i32 0), i64 3 }, align 8
@.str.89 = private unnamed_addr constant [4 x i8] c"\75\6C\65\00", align 1
@ir.ICmpUle = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.89, i32 0, i32 0), i64 3 }, align 8
@.str.90 = private unnamed_addr constant [4 x i8] c"\75\67\74\00", align 1
@ir.ICmpUgt = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.90, i32 0, i32 0), i64 3 }, align 8
@.str.91 = private unnamed_addr constant [4 x i8] c"\75\67\65\00", align 1
@ir.ICmpUge = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.91, i32 0, i32 0), i64 3 }, align 8
@.str.92 = private unnamed_addr constant [4 x i8] c"\6F\65\71\00", align 1
@ir.FCmpOeq = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.92, i32 0, i32 0), i64 3 }, align 8
@.str.93 = private unnamed_addr constant [4 x i8] c"\6F\6E\65\00", align 1
@ir.FCmpOne = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.93, i32 0, i32 0), i64 3 }, align 8
@.str.94 = private unnamed_addr constant [4 x i8] c"\6F\6C\74\00", align 1
@ir.FCmpOlt = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.94, i32 0, i32 0), i64 3 }, align 8
@.str.95 = private unnamed_addr constant [4 x i8] c"\6F\6C\65\00", align 1
@ir.FCmpOle = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.95, i32 0, i32 0), i64 3 }, align 8
@.str.96 = private unnamed_addr constant [4 x i8] c"\6F\67\74\00", align 1
@ir.FCmpOgt = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.96, i32 0, i32 0), i64 3 }, align 8
@.str.97 = private unnamed_addr constant [4 x i8] c"\6F\67\65\00", align 1
@ir.FCmpOge = constant { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.97, i32 0, i32 0), i64 3 }, align 8
@.str.98 = private unnamed_addr constant [2 x i8] c"\2A\00", align 1
@.str.99 = private unnamed_addr constant [3 x i8] c"\7B\20\00", align 1
@.str.100 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.101 = private unnamed_addr constant [3 x i8] c"\20\7D\00", align 1
@.str.102 = private unnamed_addr constant [2 x i8] c"\5B\00", align 1
@.str.103 = private unnamed_addr constant [4 x i8] c"\20\78\20\00", align 1
@.str.104 = private unnamed_addr constant [2 x i8] c"\5D\00", align 1
@.str.105 = private unnamed_addr constant [3 x i8] c"\7B\20\00", align 1
@.str.106 = private unnamed_addr constant [13 x i8] c"\2C\20\69\36\34\2C\20\69\36\34\20\7D\00", align 1
@.str.107 = private unnamed_addr constant [13 x i8] c"\7B\20\69\38\2A\2C\20\69\36\34\20\7D\00", align 1
@.str.108 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.109 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.110 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.111 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.112 = private unnamed_addr constant [4 x i8] c"\20\3D\20\00", align 1
@.str.113 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.114 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.115 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.116 = private unnamed_addr constant [9 x i8] c"\72\65\74\20\76\6F\69\64\00", align 1
@.str.117 = private unnamed_addr constant [5 x i8] c"\72\65\74\20\00", align 1
@.str.118 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.119 = private unnamed_addr constant [39 x i8] c"\72\65\74\20\76\6F\69\64\20\3B\20\54\4F\44\4F\20\6D\75\6C\74\69\70\6C\65\20\72\65\74\75\72\6E\20\76\61\6C\75\65\73\00", align 1
@.str.120 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.121 = private unnamed_addr constant [6 x i8] c"\73\72\65\74\28\00", align 1
@.str.122 = private unnamed_addr constant [3 x i8] c"\29\20\00", align 1
@.str.123 = private unnamed_addr constant [4 x i8] c"\69\38\2A\00", align 1
@.str.124 = private unnamed_addr constant [5 x i8] c"\6E\75\6C\6C\00", align 1
@.str.125 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.126 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.127 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.128 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.129 = private unnamed_addr constant [5 x i8] c"\76\6F\69\64\00", align 1
@.str.130 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.131 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.132 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.133 = private unnamed_addr constant [5 x i8] c"\76\6F\69\64\00", align 1
@.str.134 = private unnamed_addr constant [5 x i8] c"\76\6F\69\64\00", align 1
@.str.135 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.136 = private unnamed_addr constant [4 x i8] c"\20\3D\20\00", align 1
@.str.137 = private unnamed_addr constant [6 x i8] c"\63\61\6C\6C\20\00", align 1
@.str.138 = private unnamed_addr constant [3 x i8] c"\20\40\00", align 1
@.str.139 = private unnamed_addr constant [2 x i8] c"\28\00", align 1
@.str.140 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.141 = private unnamed_addr constant [2 x i8] c"\29\00", align 1
@.str.142 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.143 = private unnamed_addr constant [4 x i8] c"\20\3D\20\00", align 1
@.str.144 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.145 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.146 = private unnamed_addr constant [5 x i8] c"\20\74\6F\20\00", align 1
@.str.147 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.148 = private unnamed_addr constant [10 x i8] c"\20\2C\20\61\6C\69\67\6E\20\00", align 1
@.str.149 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.150 = private unnamed_addr constant [11 x i8] c"\20\3D\20\61\6C\6C\6F\63\61\20\00", align 1
@.str.151 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.152 = private unnamed_addr constant [9 x i8] c"\2C\20\61\6C\69\67\6E\20\00", align 1
@.str.153 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.154 = private unnamed_addr constant [9 x i8] c"\20\3D\20\6C\6F\61\64\20\00", align 1
@.str.155 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.156 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.157 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.158 = private unnamed_addr constant [9 x i8] c"\2C\20\61\6C\69\67\6E\20\00", align 1
@.str.159 = private unnamed_addr constant [7 x i8] c"\73\74\6F\72\65\20\00", align 1
@.str.160 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.161 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.162 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.163 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.164 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.165 = private unnamed_addr constant [27 x i8] c"\20\3D\20\67\65\74\65\6C\65\6D\65\6E\74\70\74\72\20\69\6E\62\6F\75\6E\64\73\20\00", align 1
@.str.166 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.167 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.168 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.169 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.170 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.171 = private unnamed_addr constant [9 x i8] c"\20\3D\20\69\63\6D\70\20\00", align 1
@.str.172 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.173 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.174 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.175 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.176 = private unnamed_addr constant [9 x i8] c"\20\3D\20\66\63\6D\70\20\00", align 1
@.str.177 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.178 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.179 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.180 = private unnamed_addr constant [11 x i8] c"\62\72\20\6C\61\62\65\6C\20\25\00", align 1
@.str.181 = private unnamed_addr constant [7 x i8] c"\62\72\20\69\31\20\00", align 1
@.str.182 = private unnamed_addr constant [10 x i8] c"\2C\20\6C\61\62\65\6C\20\25\00", align 1
@.str.183 = private unnamed_addr constant [10 x i8] c"\2C\20\6C\61\62\65\6C\20\25\00", align 1
@.str.184 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.185 = private unnamed_addr constant [12 x i8] c"\63\61\6C\6C\20\76\6F\69\64\20\40\00", align 1
@.str.186 = private unnamed_addr constant [2 x i8] c"\28\00", align 1
@.str.187 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.188 = private unnamed_addr constant [2 x i8] c"\29\00", align 1
@.str.189 = private unnamed_addr constant [6 x i8] c"\65\6E\74\72\79\00", align 1
@.str.190 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.191 = private unnamed_addr constant [12 x i8] c"\20\3D\20\62\69\74\63\61\73\74\20\00", align 1
@.str.192 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.193 = private unnamed_addr constant [5 x i8] c"\20\74\6F\20\00", align 1
@.str.194 = private unnamed_addr constant [27 x i8] c"\63\61\6C\6C\20\76\6F\69\64\20\40\67\6F\6D\69\6E\69\63\5F\6D\65\6D\63\70\79\28\00", align 1
@.str.195 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.196 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.197 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.198 = private unnamed_addr constant [3 x i8] c"\2C\20\00", align 1
@.str.199 = private unnamed_addr constant [2 x i8] c"\20\00", align 1
@.str.200 = private unnamed_addr constant [2 x i8] c"\29\00", align 1
@.str.201 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.202 = private unnamed_addr constant [2 x i8] c"\25\00", align 1
@.str.203 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.204 = private unnamed_addr constant [1 x i8] c"\00", align 1
@.str.205 = private unnamed_addr constant [7 x i8] c"\64\6F\75\62\6C\65\00", align 1
@.str.206 = private unnamed_addr constant [7 x i8] c"\64\6F\75\62\6C\65\00", align 1
@.str.207 = private unnamed_addr constant [4 x i8] c"\2E\65\45\00", align 1
@.str.208 = private unnamed_addr constant [3 x i8] c"\2E\30\00", align 1
@.str.209 = private unnamed_addr constant [5 x i8] c"\76\6F\69\64\00", align 1

define { i8*, i64 } @backend.EmitModule(i8** %mod) {
entry:
  %mod.addr = alloca i8** , align 8
  store i8** %mod, i8*** %mod.addr
  %e.addr0 = alloca { { { { i8*, i64 }*, i64, i64 } }, i1 }
  store { { { { i8*, i64 }*, i64, i64 } }, i1 } zeroinitializer, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %e.addr0
  %t1 = load i8**, i8*** %mod.addr, align 8
  call void @backend.emitTargetHeader({ { { { i8*, i64 }*, i64, i64 } }, i1 }* %e.addr0, i8** %t1)
  %taddr2 = alloca { i8**, i64, i64 } , align 8
  store { i8**, i64, i64 } zeroinitializer, { i8**, i64, i64 }* %taddr2, align 8
  %t3 = load { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr2, align 8
  %taddr4 = alloca { i8**, i64, i64 } , align 8
  store { i8**, i64, i64 } %t3, { i8**, i64, i64 }* %taddr4
  call void @backend.emitGlobals({ { { { i8*, i64 }*, i64, i64 } }, i1 }* %e.addr0, { i8**, i64, i64 }* %taddr4)
  %i.addr5 = alloca i64
  store i64 0, i64* %i.addr5, align 8
  %taddr7 = alloca { i8***, i64, i64 } , align 8
  %taddr9 = alloca { i8***, i64, i64 } , align 8
  %taddr13 = alloca { i8***, i64, i64 } , align 8
  br label %for.cond0
for.cond0:
  %t6 = load i64, i64* %i.addr5, align 8
  store { i8***, i64, i64 } zeroinitializer, { i8***, i64, i64 }* %taddr7, align 8
  %t8 = load { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr7, align 8
  store { i8***, i64, i64 } %t8, { i8***, i64, i64 }* %taddr9, align 8
  %t10 = getelementptr inbounds { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr9, i32 0, i32 1
  %t11 = load i64, i64* %t10, align 8
  %t12 = icmp slt i64 %t6, %t11
  br i1 %t12, label %for.body1, label %for.end3
for.body1:
  store { i8***, i64, i64 } zeroinitializer, { i8***, i64, i64 }* %taddr13, align 8
  %t14 = load i64, i64* %i.addr5, align 8
  %t15 = getelementptr inbounds { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr13, i32 0, i32 1
  %t16 = load i64, i64* %t15, align 8
  %t17 = icmp sge i64 %t14, 0
  %t18 = icmp slt i64 %t14, %t16
  %t19 = and i1 %t17, %t18
  br i1 %t19, label %idx.ok4, label %idx.fail5
for.post2:
  %t24 = load i64, i64* %i.addr5, align 8
  %t25 = add i64 %t24, 1
  store i64 %t25, i64* %i.addr5, align 8
  br label %for.cond0
for.end3:
  %t26 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }* %e.addr0 to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t27 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t26, i32 0, i32 0
  %t28 = call { i8*, i64 } @backend.strBufString({ { { i8*, i64 }*, i64, i64 } }* %t27)
  ret { i8*, i64 } %t28
idx.ok4:
  %t20 = getelementptr inbounds { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr13, i32 0, i32 0
  %t21 = load i8***, i8**** %t20, align 8
  %t22 = getelementptr inbounds i8**, i8*** %t21, i64 %t14
  %t23 = load i8**, i8*** %t22, align 8
  call void @backend.emitFunction({ { { { i8*, i64 }*, i64, i64 } }, i1 }* %e.addr0, i8** %t23)
  br label %for.post2
idx.fail5:
  call void @gominic_abort()
  br label %idx.ok4
}

define void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %b, { i8*, i64 }* byval({ i8*, i64 }) %s) {
entry:
  %b.addr = alloca { { { i8*, i64 }*, i64, i64 } }* , align 8
  store { { { i8*, i64 }*, i64, i64 } }* %b, { { { i8*, i64 }*, i64, i64 } }** %b.addr
  %s.addr = alloca { i8*, i64 } , align 8
  call void @gominic_memcpy({ i8*, i64 }* %s.addr, { i8*, i64 }* %s, i64 16)
  %t0 = bitcast { { { i8*, i64 }*, i64, i64 } }** %b.addr to { { { i8*, i64 }*, i64, i64 } }*
  %t1 = getelementptr inbounds { { { i8*, i64 }*, i64, i64 } }, { { { i8*, i64 }*, i64, i64 } }* %t0, i32 0, i32 0
  %t2 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %t1, align 8
  %t3 = load { i8*, i64 }, { i8*, i64 }* %s.addr, align 8
  %taddr4 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  store { { i8*, i64 }*, i64, i64 } %t2, { { i8*, i64 }*, i64, i64 }* %taddr4, align 8
  %t5 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr4, i32 0, i32 0
  %t6 = load { i8*, i64 }*, { i8*, i64 }** %t5, align 8
  %t7 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr4, i32 0, i32 1
  %t8 = load i64, i64* %t7, align 8
  %t9 = add i64 %t8, 1
  %t10 = mul i64 %t8, 16
  %t11 = call i8* @gominic_makeSlice(i64 %t9, i64 %t9, i64 16)
  %t12 = bitcast i8* %t11 to { i8*, i64 }*
  call void @gominic_memcpy({ i8*, i64 }* %t12, { i8*, i64 }* %t6, i64 %t10)
  %t13 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %t12, i64 %t8
  %taddr14 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t3, { i8*, i64 }* %taddr14
  call void @gominic_memcpy({ i8*, i64 }* %t13, { i8*, i64 }* %taddr14, i64 16)
  %taddr15 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %t16 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr15, i32 0, i32 0
  store { i8*, i64 }* %t12, { i8*, i64 }** %t16, align 8
  %t17 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr15, i32 0, i32 1
  store i64 %t9, i64* %t17, align 8
  %t18 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr15, i32 0, i32 2
  store i64 %t9, i64* %t18, align 8
  %t19 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr15, align 8
  %t20 = bitcast { { { i8*, i64 }*, i64, i64 } }** %b.addr to { { { i8*, i64 }*, i64, i64 } }*
  %t21 = getelementptr inbounds { { { i8*, i64 }*, i64, i64 } }, { { { i8*, i64 }*, i64, i64 } }* %t20, i32 0, i32 0
  %taddr22 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  store { { i8*, i64 }*, i64, i64 } %t19, { { i8*, i64 }*, i64, i64 }* %taddr22
  call void @gominic_memcpy({ { i8*, i64 }*, i64, i64 }* %t21, { { i8*, i64 }*, i64, i64 }* %taddr22, i64 24)
  ret void
}

define { i8*, i64 } @backend.strBufString({ { { i8*, i64 }*, i64, i64 } }* %b) {
entry:
  %b.addr = alloca { { { i8*, i64 }*, i64, i64 } }* , align 8
  store { { { i8*, i64 }*, i64, i64 } }* %b, { { { i8*, i64 }*, i64, i64 } }** %b.addr
  %t0 = bitcast { { { i8*, i64 }*, i64, i64 } }** %b.addr to { { { i8*, i64 }*, i64, i64 } }*
  %t1 = getelementptr inbounds { { { i8*, i64 }*, i64, i64 } }, { { { i8*, i64 }*, i64, i64 } }* %t0, i32 0, i32 0
  %t2 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %t1, align 8
  %taddr3 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  store { { i8*, i64 }*, i64, i64 } %t2, { { i8*, i64 }*, i64, i64 }* %taddr3
  %taddr4 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.0, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr4
  %t5 = call { i8*, i64 } @strings.Join({ { i8*, i64 }*, i64, i64 }* %taddr3, { i8*, i64 }* %taddr4)
  ret { i8*, i64 } %t5
}

define void @backend.emitGlobals({ { { { i8*, i64 }*, i64, i64 } }, i1 }* %e, { i8**, i64, i64 }* byval({ i8**, i64, i64 }) %globals) {
entry:
  %e.addr = alloca { { { { i8*, i64 }*, i64, i64 } }, i1 }* , align 8
  store { { { { i8*, i64 }*, i64, i64 } }, i1 }* %e, { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr
  %globals.addr = alloca { i8**, i64, i64 } , align 8
  call void @gominic_memcpy({ i8**, i64, i64 }* %globals.addr, { i8**, i64, i64 }* %globals, i64 24)
  %t0 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t1 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t0, i32 0, i32 0
  %taddr2 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([45 x i8], [45 x i8]* @.str.1, i32 0, i32 0), i64 44 }, { i8*, i64 }* %taddr2
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t1, { i8*, i64 }* %taddr2)
  %t3 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t4 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t3, i32 0, i32 0
  %taddr5 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([31 x i8], [31 x i8]* @.str.2, i32 0, i32 0), i64 30 }, { i8*, i64 }* %taddr5
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t4, { i8*, i64 }* %taddr5)
  %t6 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t7 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t6, i32 0, i32 0
  %taddr8 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([47 x i8], [47 x i8]* @.str.3, i32 0, i32 0), i64 46 }, { i8*, i64 }* %taddr8
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t7, { i8*, i64 }* %taddr8)
  %t9 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t10 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t9, i32 0, i32 0
  %taddr11 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([39 x i8], [39 x i8]* @.str.4, i32 0, i32 0), i64 38 }, { i8*, i64 }* %taddr11
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t10, { i8*, i64 }* %taddr11)
  %t12 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t13 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t12, i32 0, i32 0
  %taddr14 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([37 x i8], [37 x i8]* @.str.5, i32 0, i32 0), i64 36 }, { i8*, i64 }* %taddr14
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t13, { i8*, i64 }* %taddr14)
  %t15 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t16 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t15, i32 0, i32 0
  %taddr17 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([34 x i8], [34 x i8]* @.str.6, i32 0, i32 0), i64 33 }, { i8*, i64 }* %taddr17
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t16, { i8*, i64 }* %taddr17)
  %t18 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t19 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t18, i32 0, i32 0
  %taddr20 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([48 x i8], [48 x i8]* @.str.7, i32 0, i32 0), i64 47 }, { i8*, i64 }* %taddr20
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t19, { i8*, i64 }* %taddr20)
  %t21 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t22 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t21, i32 0, i32 0
  %taddr23 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([80 x i8], [80 x i8]* @.str.8, i32 0, i32 0), i64 79 }, { i8*, i64 }* %taddr23
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t22, { i8*, i64 }* %taddr23)
  %t24 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t25 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t24, i32 0, i32 0
  %taddr26 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([45 x i8], [45 x i8]* @.str.9, i32 0, i32 0), i64 44 }, { i8*, i64 }* %taddr26
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t25, { i8*, i64 }* %taddr26)
  %t27 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t28 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t27, i32 0, i32 0
  %taddr29 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([46 x i8], [46 x i8]* @.str.10, i32 0, i32 0), i64 45 }, { i8*, i64 }* %taddr29
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t28, { i8*, i64 }* %taddr29)
  %t30 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t31 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t30, i32 0, i32 0
  %taddr32 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([44 x i8], [44 x i8]* @.str.11, i32 0, i32 0), i64 43 }, { i8*, i64 }* %taddr32
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t31, { i8*, i64 }* %taddr32)
  %t33 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t34 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t33, i32 0, i32 0
  %taddr35 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([36 x i8], [36 x i8]* @.str.12, i32 0, i32 0), i64 35 }, { i8*, i64 }* %taddr35
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t34, { i8*, i64 }* %taddr35)
  %t36 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t37 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t36, i32 0, i32 0
  %taddr38 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([29 x i8], [29 x i8]* @.str.13, i32 0, i32 0), i64 28 }, { i8*, i64 }* %taddr38
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t37, { i8*, i64 }* %taddr38)
  %t39 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t40 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t39, i32 0, i32 1
  %t41 = load i1, i1* %t40, align 1
  %taddr44 = alloca { i8*, i64 } , align 8
  %taddr47 = alloca { i8*, i64 } , align 8
  %taddr50 = alloca { i8*, i64 } , align 8
  %taddr54 = alloca { i8**, i64, i64 } , align 8
  %taddr70 = alloca { i8*, i64 } , align 8
  %taddr71 = alloca i1 , align 1
  %taddr73 = alloca { i8*, i64 } , align 8
  %taddr75 = alloca { i8*, i64 } , align 8
  %taddr76 = alloca i64 , align 8
  %taddr79 = alloca i64 , align 8
  %taddr82 = alloca { i8*, i64 } , align 8
  %taddr83 = alloca { i8*, i64 } , align 8
  %taddr84 = alloca { i8*, i64 } , align 8
  %taddr86 = alloca { i8*, i64 } , align 8
  %taddr89 = alloca { i8*, i64 } , align 8
  %taddr91 = alloca { i8*, i64 } , align 8
  %taddr92 = alloca { i8*, i64 } , align 8
  %taddr93 = alloca { i8*, i64 } , align 8
  %taddr95 = alloca { i8*, i64 } , align 8
  %taddr96 = alloca { i8*, i64 } , align 8
  %taddr97 = alloca { i8*, i64 } , align 8
  %taddr100 = alloca { i8*, i64 } , align 8
  %taddr101 = alloca { i8*, i64 } , align 8
  %taddr102 = alloca { i8*, i64 } , align 8
  %taddr104 = alloca { i8*, i64 } , align 8
  %taddr105 = alloca { i8*, i64 } , align 8
  %taddr106 = alloca { i8*, i64 } , align 8
  %taddr108 = alloca i8* , align 8
  %taddr111 = alloca { i8*, i64 } , align 8
  %taddr112 = alloca { i8*, i64 } , align 8
  %taddr113 = alloca { i8*, i64 } , align 8
  %taddr115 = alloca { i8*, i64 } , align 8
  %taddr116 = alloca { i8*, i64 } , align 8
  %taddr117 = alloca { i8*, i64 } , align 8
  %taddr119 = alloca { i8*, i64 } , align 8
  %taddr121 = alloca { i8*, i64 } , align 8
  %taddr122 = alloca { i8*, i64 } , align 8
  %taddr123 = alloca { i8*, i64 } , align 8
  %taddr126 = alloca { i8*, i64 } , align 8
  %taddr127 = alloca { i8*, i64 } , align 8
  %taddr128 = alloca { i8*, i64 } , align 8
  %taddr130 = alloca { i8*, i64 } , align 8
  %taddr131 = alloca { i8*, i64 } , align 8
  %taddr132 = alloca { i8*, i64 } , align 8
  %taddr134 = alloca { i8*, i64 } , align 8
  %taddr138 = alloca { i8**, i64, i64 } , align 8
  %taddr144 = alloca { i8*, i64 } , align 8
  br i1 %t41, label %then0, label %else2
then0:
  %t42 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t43 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t42, i32 0, i32 0
  store { i8*, i64 } { i8* getelementptr inbounds ([67 x i8], [67 x i8]* @.str.14, i32 0, i32 0), i64 66 }, { i8*, i64 }* %taddr44
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t43, { i8*, i64 }* %taddr44)
  br label %endif1
endif1:
  %t48 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t49 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t48, i32 0, i32 0
  store { i8*, i64 } { i8* getelementptr inbounds ([53 x i8], [53 x i8]* @.str.16, i32 0, i32 0), i64 52 }, { i8*, i64 }* %taddr50
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t49, { i8*, i64 }* %taddr50)
  %i.addr51 = alloca i64
  store i64 0, i64* %i.addr51, align 8
  br label %for.cond3
else2:
  %t45 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t46 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t45, i32 0, i32 0
  store { i8*, i64 } { i8* getelementptr inbounds ([41 x i8], [41 x i8]* @.str.15, i32 0, i32 0), i64 40 }, { i8*, i64 }* %taddr47
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t46, { i8*, i64 }* %taddr47)
  br label %endif1
for.cond3:
  %t52 = load i64, i64* %i.addr51, align 8
  %t53 = load { i8**, i64, i64 }, { i8**, i64, i64 }* %globals.addr, align 8
  store { i8**, i64, i64 } %t53, { i8**, i64, i64 }* %taddr54, align 8
  %t55 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr54, i32 0, i32 1
  %t56 = load i64, i64* %t55, align 8
  %t57 = icmp slt i64 %t52, %t56
  br i1 %t57, label %for.body4, label %for.end6
for.body4:
  %t58 = load i64, i64* %i.addr51, align 8
  %t59 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %globals.addr, i32 0, i32 1
  %t60 = load i64, i64* %t59, align 8
  %t61 = icmp sge i64 %t58, 0
  %t62 = icmp slt i64 %t58, %t60
  %t63 = and i1 %t61, %t62
  br i1 %t63, label %idx.ok7, label %idx.fail8
for.post5:
  %t135 = load i64, i64* %i.addr51, align 8
  %t136 = add i64 %t135, 1
  store i64 %t136, i64* %i.addr51, align 8
  br label %for.cond3
for.end6:
  %t137 = load { i8**, i64, i64 }, { i8**, i64, i64 }* %globals.addr, align 8
  store { i8**, i64, i64 } %t137, { i8**, i64, i64 }* %taddr138, align 8
  %t139 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr138, i32 0, i32 1
  %t140 = load i64, i64* %t139, align 8
  %t141 = icmp sgt i64 %t140, 0
  br i1 %t141, label %then13, label %endif14
idx.ok7:
  %t64 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %globals.addr, i32 0, i32 0
  %t65 = load i8**, i8*** %t64, align 8
  %t66 = getelementptr inbounds i8*, i8** %t65, i64 %t58
  %t67 = load i8*, i8** %t66, align 8
  %g.addr68 = alloca i8*
  store i8* %t67, i8** %g.addr68, align 8
  %visibility.addr69 = alloca { i8*, i64 }
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.17, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr70
  call void @gominic_memcpy({ i8*, i64 }* %visibility.addr69, { i8*, i64 }* %taddr70, i64 16)
  store i1 0, i1* %taddr71, align 1
  %t72 = load i1, i1* %taddr71, align 1
  br i1 %t72, label %then9, label %endif10
idx.fail8:
  call void @gominic_abort()
  br label %idx.ok7
then9:
  store { i8*, i64 } { i8* getelementptr inbounds ([22 x i8], [22 x i8]* @.str.18, i32 0, i32 0), i64 21 }, { i8*, i64 }* %taddr73
  call void @gominic_memcpy({ i8*, i64 }* %visibility.addr69, { i8*, i64 }* %taddr73, i64 16)
  br label %endif10
endif10:
  %align.addr74 = alloca { i8*, i64 }
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.19, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr75
  call void @gominic_memcpy({ i8*, i64 }* %align.addr74, { i8*, i64 }* %taddr75, i64 16)
  store i64 0, i64* %taddr76, align 8
  %t77 = load i64, i64* %taddr76, align 8
  %t78 = icmp sgt i64 %t77, 0
  br i1 %t78, label %then11, label %endif12
then11:
  store i64 0, i64* %taddr79, align 8
  %t80 = load i64, i64* %taddr79, align 8
  %t81 = call { i8*, i64 } @strconv.FormatInt(i64 %t80, i64 10)
  store { i8*, i64 } { i8* getelementptr inbounds ([9 x i8], [9 x i8]* @.str.20, i32 0, i32 0), i64 8 }, { i8*, i64 }* %taddr83, align 8
  store { i8*, i64 } %t81, { i8*, i64 }* %taddr84, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr82, { i8*, i64 }* %taddr83, { i8*, i64 }* %taddr84)
  %t85 = load { i8*, i64 }, { i8*, i64 }* %taddr82, align 8
  store { i8*, i64 } %t85, { i8*, i64 }* %taddr86
  call void @gominic_memcpy({ i8*, i64 }* %align.addr74, { i8*, i64 }* %taddr86, i64 16)
  br label %endif12
endif12:
  %t87 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t88 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t87, i32 0, i32 0
  store { i8*, i64 } zeroinitializer, { i8*, i64 }* %taddr89, align 8
  %t90 = load { i8*, i64 }, { i8*, i64 }* %taddr89, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.21, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr92, align 8
  store { i8*, i64 } %t90, { i8*, i64 }* %taddr93, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr91, { i8*, i64 }* %taddr92, { i8*, i64 }* %taddr93)
  %t94 = load { i8*, i64 }, { i8*, i64 }* %taddr91, align 8
  store { i8*, i64 } %t94, { i8*, i64 }* %taddr96, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.22, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr97, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr95, { i8*, i64 }* %taddr96, { i8*, i64 }* %taddr97)
  %t98 = load { i8*, i64 }, { i8*, i64 }* %taddr95, align 8
  %t99 = load { i8*, i64 }, { i8*, i64 }* %visibility.addr69, align 8
  store { i8*, i64 } %t98, { i8*, i64 }* %taddr101, align 8
  store { i8*, i64 } %t99, { i8*, i64 }* %taddr102, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr100, { i8*, i64 }* %taddr101, { i8*, i64 }* %taddr102)
  %t103 = load { i8*, i64 }, { i8*, i64 }* %taddr100, align 8
  store { i8*, i64 } %t103, { i8*, i64 }* %taddr105, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([10 x i8], [10 x i8]* @.str.23, i32 0, i32 0), i64 9 }, { i8*, i64 }* %taddr106, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr104, { i8*, i64 }* %taddr105, { i8*, i64 }* %taddr106)
  %t107 = load { i8*, i64 }, { i8*, i64 }* %taddr104, align 8
  store i8* null, i8** %taddr108, align 8
  %t109 = load i8*, i8** %taddr108, align 8
  %t110 = call { i8*, i64 } @backend.llvmType(i8* %t109)
  store { i8*, i64 } %t107, { i8*, i64 }* %taddr112, align 8
  store { i8*, i64 } %t110, { i8*, i64 }* %taddr113, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr111, { i8*, i64 }* %taddr112, { i8*, i64 }* %taddr113)
  %t114 = load { i8*, i64 }, { i8*, i64 }* %taddr111, align 8
  store { i8*, i64 } %t114, { i8*, i64 }* %taddr116, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.24, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr117, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr115, { i8*, i64 }* %taddr116, { i8*, i64 }* %taddr117)
  %t118 = load { i8*, i64 }, { i8*, i64 }* %taddr115, align 8
  store { i8*, i64 } zeroinitializer, { i8*, i64 }* %taddr119, align 8
  %t120 = load { i8*, i64 }, { i8*, i64 }* %taddr119, align 8
  store { i8*, i64 } %t118, { i8*, i64 }* %taddr122, align 8
  store { i8*, i64 } %t120, { i8*, i64 }* %taddr123, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr121, { i8*, i64 }* %taddr122, { i8*, i64 }* %taddr123)
  %t124 = load { i8*, i64 }, { i8*, i64 }* %taddr121, align 8
  %t125 = load { i8*, i64 }, { i8*, i64 }* %align.addr74, align 8
  store { i8*, i64 } %t124, { i8*, i64 }* %taddr127, align 8
  store { i8*, i64 } %t125, { i8*, i64 }* %taddr128, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr126, { i8*, i64 }* %taddr127, { i8*, i64 }* %taddr128)
  %t129 = load { i8*, i64 }, { i8*, i64 }* %taddr126, align 8
  store { i8*, i64 } %t129, { i8*, i64 }* %taddr131, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.25, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr132, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr130, { i8*, i64 }* %taddr131, { i8*, i64 }* %taddr132)
  %t133 = load { i8*, i64 }, { i8*, i64 }* %taddr130, align 8
  store { i8*, i64 } %t133, { i8*, i64 }* %taddr134
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t88, { i8*, i64 }* %taddr134)
  br label %for.post5
then13:
  %t142 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t143 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t142, i32 0, i32 0
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.26, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr144
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t143, { i8*, i64 }* %taddr144)
  br label %endif14
endif14:
  ret void
}

define void @backend.emitFunction({ { { { i8*, i64 }*, i64, i64 } }, i1 }* %e, i8** %fn) {
entry:
  %e.addr = alloca { { { { i8*, i64 }*, i64, i64 } }, i1 }* , align 8
  store { { { { i8*, i64 }*, i64, i64 } }, i1 }* %e, { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr
  %fn.addr = alloca i8** , align 8
  store i8** %fn, i8*** %fn.addr
  %taddr0 = alloca { i8***, i64, i64 } , align 8
  store { i8***, i64, i64 } zeroinitializer, { i8***, i64, i64 }* %taddr0, align 8
  %t1 = load { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr0, align 8
  %taddr2 = alloca { i8***, i64, i64 } , align 8
  store { i8***, i64, i64 } %t1, { i8***, i64, i64 }* %taddr2, align 8
  %t3 = getelementptr inbounds { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr2, i32 0, i32 1
  %t4 = load i64, i64* %t3, align 8
  %t5 = call i8* @gominic_makeSlice(i64 %t4, i64 %t4, i64 16)
  %t6 = bitcast i8* %t5 to { i8*, i64 }*
  %taddr7 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %t8 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr7, i32 0, i32 0
  store { i8*, i64 }* %t6, { i8*, i64 }** %t8, align 8
  %t9 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr7, i32 0, i32 1
  store i64 %t4, i64* %t9, align 8
  %t10 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr7, i32 0, i32 2
  store i64 %t4, i64* %t10, align 8
  %t11 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr7, align 8
  %params.addr12 = alloca { { i8*, i64 }*, i64, i64 }
  %taddr13 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  store { { i8*, i64 }*, i64, i64 } %t11, { { i8*, i64 }*, i64, i64 }* %taddr13
  call void @gominic_memcpy({ { i8*, i64 }*, i64, i64 }* %params.addr12, { { i8*, i64 }*, i64, i64 }* %taddr13, i64 24)
  %i.addr14 = alloca i64
  store i64 0, i64* %i.addr14, align 8
  %taddr16 = alloca { i8***, i64, i64 } , align 8
  %taddr18 = alloca { i8***, i64, i64 } , align 8
  %taddr22 = alloca { i8***, i64, i64 } , align 8
  %taddr36 = alloca i8* , align 8
  %taddr39 = alloca i1 , align 1
  %taddr40 = alloca i8* , align 8
  %taddr44 = alloca i8* , align 8
  %taddr47 = alloca { i8*, i64 } , align 8
  %taddr48 = alloca { i8*, i64 } , align 8
  %taddr49 = alloca { i8*, i64 } , align 8
  %taddr51 = alloca { i8*, i64 } , align 8
  %taddr52 = alloca { i8*, i64 } , align 8
  %taddr53 = alloca { i8*, i64 } , align 8
  %taddr57 = alloca { i8*, i64 } , align 8
  %taddr58 = alloca { i8*, i64 } , align 8
  %taddr59 = alloca { i8*, i64 } , align 8
  %taddr61 = alloca { i8*, i64 } , align 8
  %taddr62 = alloca { i8*, i64 } , align 8
  %taddr63 = alloca { i8*, i64 } , align 8
  %taddr66 = alloca { i8*, i64 } , align 8
  %taddr67 = alloca { i8*, i64 } , align 8
  %taddr68 = alloca { i8*, i64 } , align 8
  %taddr79 = alloca { i8*, i64 } , align 8
  %taddr83 = alloca { i8*, i64 } , align 8
  %taddr84 = alloca { i8*, i64 } , align 8
  %taddr85 = alloca { i8*, i64 } , align 8
  %taddr89 = alloca { i8*, i64 } , align 8
  %taddr90 = alloca { i8*, i64 } , align 8
  %taddr91 = alloca { i8*, i64 } , align 8
  %taddr93 = alloca { i8*, i64 } , align 8
  %taddr94 = alloca { i8*, i64 } , align 8
  %taddr95 = alloca { i8*, i64 } , align 8
  %taddr98 = alloca { i8*, i64 } , align 8
  %taddr99 = alloca { i8*, i64 } , align 8
  %taddr100 = alloca { i8*, i64 } , align 8
  %taddr111 = alloca { i8*, i64 } , align 8
  %taddr114 = alloca { i8*, i64 } , align 8
  %taddr115 = alloca { i8*, i64 } , align 8
  %taddr116 = alloca { i8*, i64 } , align 8
  %taddr119 = alloca { i8*, i64 } , align 8
  %taddr120 = alloca { i8*, i64 } , align 8
  %taddr121 = alloca { i8*, i64 } , align 8
  %taddr132 = alloca { i8*, i64 } , align 8
  %taddr136 = alloca i8* , align 8
  %taddr141 = alloca { i8*, i64 } , align 8
  %taddr142 = alloca { i8**, i64, i64 } , align 8
  %taddr144 = alloca { i8**, i64, i64 } , align 8
  %taddr152 = alloca { i8*, i64 } , align 8
  %taddr155 = alloca { i8**, i64, i64 } , align 8
  %taddr166 = alloca { i8*, i64 } , align 8
  %taddr169 = alloca { i8*, i64 } , align 8
  %taddr173 = alloca { i8*, i64 } , align 8
  %taddr174 = alloca { i8*, i64 } , align 8
  %taddr175 = alloca { i8*, i64 } , align 8
  %taddr177 = alloca { i8*, i64 } , align 8
  %taddr178 = alloca { i8*, i64 } , align 8
  %taddr179 = alloca { i8*, i64 } , align 8
  %taddr181 = alloca { i8*, i64 } , align 8
  %taddr183 = alloca { i8*, i64 } , align 8
  %taddr184 = alloca { i8*, i64 } , align 8
  %taddr185 = alloca { i8*, i64 } , align 8
  %taddr187 = alloca { i8*, i64 } , align 8
  %taddr188 = alloca { i8*, i64 } , align 8
  %taddr189 = alloca { i8*, i64 } , align 8
  %taddr192 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %taddr193 = alloca { i8*, i64 } , align 8
  %taddr195 = alloca { i8*, i64 } , align 8
  %taddr196 = alloca { i8*, i64 } , align 8
  %taddr197 = alloca { i8*, i64 } , align 8
  %taddr199 = alloca { i8*, i64 } , align 8
  %taddr200 = alloca { i8*, i64 } , align 8
  %taddr201 = alloca { i8*, i64 } , align 8
  %taddr203 = alloca { i8*, i64 } , align 8
  %taddr206 = alloca { i8***, i64, i64 } , align 8
  %taddr208 = alloca { i8***, i64, i64 } , align 8
  %taddr213 = alloca { i8***, i64, i64 } , align 8
  %taddr228 = alloca { i8*, i64 } , align 8
  br label %for.cond0
for.cond0:
  %t15 = load i64, i64* %i.addr14, align 8
  store { i8***, i64, i64 } zeroinitializer, { i8***, i64, i64 }* %taddr16, align 8
  %t17 = load { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr16, align 8
  store { i8***, i64, i64 } %t17, { i8***, i64, i64 }* %taddr18, align 8
  %t19 = getelementptr inbounds { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr18, i32 0, i32 1
  %t20 = load i64, i64* %t19, align 8
  %t21 = icmp slt i64 %t15, %t20
  br i1 %t21, label %for.body1, label %for.end3
for.body1:
  store { i8***, i64, i64 } zeroinitializer, { i8***, i64, i64 }* %taddr22, align 8
  %t23 = load i64, i64* %i.addr14, align 8
  %t24 = getelementptr inbounds { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr22, i32 0, i32 1
  %t25 = load i64, i64* %t24, align 8
  %t26 = icmp sge i64 %t23, 0
  %t27 = icmp slt i64 %t23, %t25
  %t28 = and i1 %t26, %t27
  br i1 %t28, label %idx.ok4, label %idx.fail5
for.post2:
  %t133 = load i64, i64* %i.addr14, align 8
  %t134 = add i64 %t133, 1
  store i64 %t134, i64* %i.addr14, align 8
  br label %for.cond0
for.end3:
  %result.addr135 = alloca { i8*, i64 }
  store { i8*, i64 } zeroinitializer, { i8*, i64 }* %result.addr135
  store i8* null, i8** %taddr136, align 8
  %t137 = load i8*, i8** %taddr136, align 8
  %t138 = icmp ne i8* %t137, null
  br i1 %t138, label %then21, label %else23
idx.ok4:
  %t29 = getelementptr inbounds { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr22, i32 0, i32 0
  %t30 = load i8***, i8**** %t29, align 8
  %t31 = getelementptr inbounds i8**, i8*** %t30, i64 %t23
  %t32 = load i8**, i8*** %t31, align 8
  %p.addr33 = alloca i8**
  store i8** %t32, i8*** %p.addr33, align 8
  %t34 = load i64, i64* %i.addr14, align 8
  %t35 = icmp eq i64 %t34, 0
  store i8* null, i8** %taddr36, align 8
  %t37 = load i8*, i8** %taddr36, align 8
  %t38 = icmp ne i8* %t37, null
  br i1 %t35, label %logic.rhs7, label %logic.false8
idx.fail5:
  call void @gominic_abort()
  br label %idx.ok4
logic.end6:
  %t43 = load i1, i1* %taddr39, align 1
  br i1 %t43, label %then9, label %else11
logic.rhs7:
  store i8* null, i8** %taddr40, align 8
  %t41 = load i8*, i8** %taddr40, align 8
  %t42 = icmp ne i8* %t41, null
  store i1 %t42, i1* %taddr39
  br label %logic.end6
logic.false8:
  store i1 0, i1* %taddr39
  br label %logic.end6
then9:
  store i8* null, i8** %taddr44, align 8
  %t45 = load i8*, i8** %taddr44, align 8
  %t46 = call { i8*, i64 } @backend.llvmType(i8* %t45)
  store { i8*, i64 } { i8* getelementptr inbounds ([6 x i8], [6 x i8]* @.str.27, i32 0, i32 0), i64 5 }, { i8*, i64 }* %taddr48, align 8
  store { i8*, i64 } %t46, { i8*, i64 }* %taddr49, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr47, { i8*, i64 }* %taddr48, { i8*, i64 }* %taddr49)
  %t50 = load { i8*, i64 }, { i8*, i64 }* %taddr47, align 8
  store { i8*, i64 } %t50, { i8*, i64 }* %taddr52, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.28, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr53, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr51, { i8*, i64 }* %taddr52, { i8*, i64 }* %taddr53)
  %t54 = load { i8*, i64 }, { i8*, i64 }* %taddr51, align 8
  %t55 = call i8* @p.Type()
  %t56 = call { i8*, i64 } @backend.llvmType(i8* %t55)
  store { i8*, i64 } %t54, { i8*, i64 }* %taddr58, align 8
  store { i8*, i64 } %t56, { i8*, i64 }* %taddr59, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr57, { i8*, i64 }* %taddr58, { i8*, i64 }* %taddr59)
  %t60 = load { i8*, i64 }, { i8*, i64 }* %taddr57, align 8
  store { i8*, i64 } %t60, { i8*, i64 }* %taddr62, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.29, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr63, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr61, { i8*, i64 }* %taddr62, { i8*, i64 }* %taddr63)
  %t64 = load { i8*, i64 }, { i8*, i64 }* %taddr61, align 8
  %t65 = call { i8*, i64 } @p.Name()
  store { i8*, i64 } %t64, { i8*, i64 }* %taddr67, align 8
  store { i8*, i64 } %t65, { i8*, i64 }* %taddr68, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr66, { i8*, i64 }* %taddr67, { i8*, i64 }* %taddr68)
  %t69 = load { i8*, i64 }, { i8*, i64 }* %taddr66, align 8
  %t70 = load i64, i64* %i.addr14, align 8
  %t71 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %params.addr12, i32 0, i32 1
  %t72 = load i64, i64* %t71, align 8
  %t73 = icmp sge i64 %t70, 0
  %t74 = icmp slt i64 %t70, %t72
  %t75 = and i1 %t73, %t74
  br i1 %t75, label %idx.ok12, label %idx.fail13
endif10:
  br label %for.post2
else11:
  %t80 = call i1 @p.ByVal()
  br i1 %t80, label %then14, label %else16
idx.ok12:
  %t76 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %params.addr12, i32 0, i32 0
  %t77 = load { i8*, i64 }*, { i8*, i64 }** %t76, align 8
  %t78 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %t77, i64 %t70
  store { i8*, i64 } %t69, { i8*, i64 }* %taddr79
  call void @gominic_memcpy({ i8*, i64 }* %t78, { i8*, i64 }* %taddr79, i64 16)
  br label %endif10
idx.fail13:
  call void @gominic_abort()
  br label %idx.ok12
then14:
  %t81 = call i8* @p.Type()
  %t82 = call { i8*, i64 } @backend.llvmType(i8* %t81)
  store { i8*, i64 } %t82, { i8*, i64 }* %taddr84, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([8 x i8], [8 x i8]* @.str.30, i32 0, i32 0), i64 7 }, { i8*, i64 }* %taddr85, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr83, { i8*, i64 }* %taddr84, { i8*, i64 }* %taddr85)
  %t86 = load { i8*, i64 }, { i8*, i64 }* %taddr83, align 8
  %t87 = call i8* @p.ByValType()
  %t88 = call { i8*, i64 } @backend.llvmType(i8* %t87)
  store { i8*, i64 } %t86, { i8*, i64 }* %taddr90, align 8
  store { i8*, i64 } %t88, { i8*, i64 }* %taddr91, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr89, { i8*, i64 }* %taddr90, { i8*, i64 }* %taddr91)
  %t92 = load { i8*, i64 }, { i8*, i64 }* %taddr89, align 8
  store { i8*, i64 } %t92, { i8*, i64 }* %taddr94, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.31, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr95, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr93, { i8*, i64 }* %taddr94, { i8*, i64 }* %taddr95)
  %t96 = load { i8*, i64 }, { i8*, i64 }* %taddr93, align 8
  %t97 = call { i8*, i64 } @p.Name()
  store { i8*, i64 } %t96, { i8*, i64 }* %taddr99, align 8
  store { i8*, i64 } %t97, { i8*, i64 }* %taddr100, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr98, { i8*, i64 }* %taddr99, { i8*, i64 }* %taddr100)
  %t101 = load { i8*, i64 }, { i8*, i64 }* %taddr98, align 8
  %t102 = load i64, i64* %i.addr14, align 8
  %t103 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %params.addr12, i32 0, i32 1
  %t104 = load i64, i64* %t103, align 8
  %t105 = icmp sge i64 %t102, 0
  %t106 = icmp slt i64 %t102, %t104
  %t107 = and i1 %t105, %t106
  br i1 %t107, label %idx.ok17, label %idx.fail18
endif15:
  br label %endif10
else16:
  %t112 = call i8* @p.Type()
  %t113 = call { i8*, i64 } @backend.llvmType(i8* %t112)
  store { i8*, i64 } %t113, { i8*, i64 }* %taddr115, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.32, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr116, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr114, { i8*, i64 }* %taddr115, { i8*, i64 }* %taddr116)
  %t117 = load { i8*, i64 }, { i8*, i64 }* %taddr114, align 8
  %t118 = call { i8*, i64 } @p.Name()
  store { i8*, i64 } %t117, { i8*, i64 }* %taddr120, align 8
  store { i8*, i64 } %t118, { i8*, i64 }* %taddr121, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr119, { i8*, i64 }* %taddr120, { i8*, i64 }* %taddr121)
  %t122 = load { i8*, i64 }, { i8*, i64 }* %taddr119, align 8
  %t123 = load i64, i64* %i.addr14, align 8
  %t124 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %params.addr12, i32 0, i32 1
  %t125 = load i64, i64* %t124, align 8
  %t126 = icmp sge i64 %t123, 0
  %t127 = icmp slt i64 %t123, %t125
  %t128 = and i1 %t126, %t127
  br i1 %t128, label %idx.ok19, label %idx.fail20
idx.ok17:
  %t108 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %params.addr12, i32 0, i32 0
  %t109 = load { i8*, i64 }*, { i8*, i64 }** %t108, align 8
  %t110 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %t109, i64 %t102
  store { i8*, i64 } %t101, { i8*, i64 }* %taddr111
  call void @gominic_memcpy({ i8*, i64 }* %t110, { i8*, i64 }* %taddr111, i64 16)
  br label %endif15
idx.fail18:
  call void @gominic_abort()
  br label %idx.ok17
idx.ok19:
  %t129 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %params.addr12, i32 0, i32 0
  %t130 = load { i8*, i64 }*, { i8*, i64 }** %t129, align 8
  %t131 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %t130, i64 %t123
  store { i8*, i64 } %t122, { i8*, i64 }* %taddr132
  call void @gominic_memcpy({ i8*, i64 }* %t131, { i8*, i64 }* %taddr132, i64 16)
  br label %endif15
idx.fail20:
  call void @gominic_abort()
  br label %idx.ok19
then21:
  %t139 = load { i8*, i64 }, { i8*, i64 }* @ir.Void, align 8
  %t140 = call { i8*, i64 } @backend.llvmType({ i8*, i64 } %t139)
  store { i8*, i64 } %t140, { i8*, i64 }* %taddr141
  call void @gominic_memcpy({ i8*, i64 }* %result.addr135, { i8*, i64 }* %taddr141, i64 16)
  br label %endif22
endif22:
  %t170 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t171 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t170, i32 0, i32 0
  %t172 = load { i8*, i64 }, { i8*, i64 }* %result.addr135, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([8 x i8], [8 x i8]* @.str.33, i32 0, i32 0), i64 7 }, { i8*, i64 }* %taddr174, align 8
  store { i8*, i64 } %t172, { i8*, i64 }* %taddr175, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr173, { i8*, i64 }* %taddr174, { i8*, i64 }* %taddr175)
  %t176 = load { i8*, i64 }, { i8*, i64 }* %taddr173, align 8
  store { i8*, i64 } %t176, { i8*, i64 }* %taddr178, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.34, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr179, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr177, { i8*, i64 }* %taddr178, { i8*, i64 }* %taddr179)
  %t180 = load { i8*, i64 }, { i8*, i64 }* %taddr177, align 8
  store { i8*, i64 } zeroinitializer, { i8*, i64 }* %taddr181, align 8
  %t182 = load { i8*, i64 }, { i8*, i64 }* %taddr181, align 8
  store { i8*, i64 } %t180, { i8*, i64 }* %taddr184, align 8
  store { i8*, i64 } %t182, { i8*, i64 }* %taddr185, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr183, { i8*, i64 }* %taddr184, { i8*, i64 }* %taddr185)
  %t186 = load { i8*, i64 }, { i8*, i64 }* %taddr183, align 8
  store { i8*, i64 } %t186, { i8*, i64 }* %taddr188, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.35, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr189, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr187, { i8*, i64 }* %taddr188, { i8*, i64 }* %taddr189)
  %t190 = load { i8*, i64 }, { i8*, i64 }* %taddr187, align 8
  %t191 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %params.addr12, align 8
  store { { i8*, i64 }*, i64, i64 } %t191, { { i8*, i64 }*, i64, i64 }* %taddr192
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.36, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr193
  %t194 = call { i8*, i64 } @strings.Join({ { i8*, i64 }*, i64, i64 }* %taddr192, { i8*, i64 }* %taddr193)
  store { i8*, i64 } %t190, { i8*, i64 }* %taddr196, align 8
  store { i8*, i64 } %t194, { i8*, i64 }* %taddr197, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr195, { i8*, i64 }* %taddr196, { i8*, i64 }* %taddr197)
  %t198 = load { i8*, i64 }, { i8*, i64 }* %taddr195, align 8
  store { i8*, i64 } %t198, { i8*, i64 }* %taddr200, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.37, i32 0, i32 0), i64 4 }, { i8*, i64 }* %taddr201, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr199, { i8*, i64 }* %taddr200, { i8*, i64 }* %taddr201)
  %t202 = load { i8*, i64 }, { i8*, i64 }* %taddr199, align 8
  store { i8*, i64 } %t202, { i8*, i64 }* %taddr203
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t171, { i8*, i64 }* %taddr203)
  %i.addr204 = alloca i64
  store i64 0, i64* %i.addr204, align 8
  br label %for.cond32
else23:
  store { i8**, i64, i64 } zeroinitializer, { i8**, i64, i64 }* %taddr142, align 8
  %t143 = load { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr142, align 8
  store { i8**, i64, i64 } %t143, { i8**, i64, i64 }* %taddr144, align 8
  %t145 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr144, i32 0, i32 1
  %t146 = load i64, i64* %t145, align 8
  %resCount.addr147 = alloca i64
  store i64 %t146, i64* %resCount.addr147, align 8
  %t148 = load i64, i64* %resCount.addr147, align 8
  %t149 = icmp eq i64 %t148, 0
  br i1 %t149, label %then24, label %else26
then24:
  %t150 = load { i8*, i64 }, { i8*, i64 }* @ir.Void, align 8
  %t151 = call { i8*, i64 } @backend.llvmType({ i8*, i64 } %t150)
  store { i8*, i64 } %t151, { i8*, i64 }* %taddr152
  call void @gominic_memcpy({ i8*, i64 }* %result.addr135, { i8*, i64 }* %taddr152, i64 16)
  br label %endif25
endif25:
  br label %endif22
else26:
  %t153 = load i64, i64* %resCount.addr147, align 8
  %t154 = icmp eq i64 %t153, 1
  br i1 %t154, label %then27, label %else29
then27:
  store { i8**, i64, i64 } zeroinitializer, { i8**, i64, i64 }* %taddr155, align 8
  %t156 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr155, i32 0, i32 1
  %t157 = load i64, i64* %t156, align 8
  %t158 = icmp sge i64 0, 0
  %t159 = icmp slt i64 0, %t157
  %t160 = and i1 %t158, %t159
  br i1 %t160, label %idx.ok30, label %idx.fail31
endif28:
  br label %endif25
else29:
  %t167 = load { i8*, i64 }, { i8*, i64 }* @ir.Void, align 8
  %t168 = call { i8*, i64 } @backend.llvmType({ i8*, i64 } %t167)
  store { i8*, i64 } %t168, { i8*, i64 }* %taddr169
  call void @gominic_memcpy({ i8*, i64 }* %result.addr135, { i8*, i64 }* %taddr169, i64 16)
  br label %endif28
idx.ok30:
  %t161 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr155, i32 0, i32 0
  %t162 = load i8**, i8*** %t161, align 8
  %t163 = getelementptr inbounds i8*, i8** %t162, i64 0
  %t164 = load i8*, i8** %t163, align 8
  %t165 = call { i8*, i64 } @backend.llvmType(i8* %t164)
  store { i8*, i64 } %t165, { i8*, i64 }* %taddr166
  call void @gominic_memcpy({ i8*, i64 }* %result.addr135, { i8*, i64 }* %taddr166, i64 16)
  br label %endif28
idx.fail31:
  call void @gominic_abort()
  br label %idx.ok30
for.cond32:
  %t205 = load i64, i64* %i.addr204, align 8
  store { i8***, i64, i64 } zeroinitializer, { i8***, i64, i64 }* %taddr206, align 8
  %t207 = load { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr206, align 8
  store { i8***, i64, i64 } %t207, { i8***, i64, i64 }* %taddr208, align 8
  %t209 = getelementptr inbounds { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr208, i32 0, i32 1
  %t210 = load i64, i64* %t209, align 8
  %t211 = icmp slt i64 %t205, %t210
  br i1 %t211, label %for.body33, label %for.end35
for.body33:
  %t212 = load { { { { i8*, i64 }*, i64, i64 } }, i1 }*, { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr, align 8
  store { i8***, i64, i64 } zeroinitializer, { i8***, i64, i64 }* %taddr213, align 8
  %t214 = load i64, i64* %i.addr204, align 8
  %t215 = getelementptr inbounds { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr213, i32 0, i32 1
  %t216 = load i64, i64* %t215, align 8
  %t217 = icmp sge i64 %t214, 0
  %t218 = icmp slt i64 %t214, %t216
  %t219 = and i1 %t217, %t218
  br i1 %t219, label %idx.ok36, label %idx.fail37
for.post34:
  %t224 = load i64, i64* %i.addr204, align 8
  %t225 = add i64 %t224, 1
  store i64 %t225, i64* %i.addr204, align 8
  br label %for.cond32
for.end35:
  %t226 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t227 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t226, i32 0, i32 0
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.38, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr228
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t227, { i8*, i64 }* %taddr228)
  ret void
idx.ok36:
  %t220 = getelementptr inbounds { i8***, i64, i64 }, { i8***, i64, i64 }* %taddr213, i32 0, i32 0
  %t221 = load i8***, i8**** %t220, align 8
  %t222 = getelementptr inbounds i8**, i8*** %t221, i64 %t214
  %t223 = load i8**, i8*** %t222, align 8
  call void @backend.emitBasicBlock({ { { { i8*, i64 }*, i64, i64 } }, i1 }* %t212, i8** %t223)
  br label %for.post34
idx.fail37:
  call void @gominic_abort()
  br label %idx.ok36
}

define void @backend.emitBasicBlock({ { { { i8*, i64 }*, i64, i64 } }, i1 }* %e, i8** %bb) {
entry:
  %e.addr = alloca { { { { i8*, i64 }*, i64, i64 } }, i1 }* , align 8
  store { { { { i8*, i64 }*, i64, i64 } }, i1 }* %e, { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr
  %bb.addr = alloca i8** , align 8
  store i8** %bb, i8*** %bb.addr
  %t0 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t1 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t0, i32 0, i32 0
  %taddr2 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } zeroinitializer, { i8*, i64 }* %taddr2, align 8
  %t3 = load { i8*, i64 }, { i8*, i64 }* %taddr2, align 8
  %taddr4 = alloca { i8*, i64 } , align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t3, { i8*, i64 }* %taddr5, align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.39, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr6, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr4, { i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6)
  %t7 = load { i8*, i64 }, { i8*, i64 }* %taddr4, align 8
  %taddr8 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t7, { i8*, i64 }* %taddr8
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t1, { i8*, i64 }* %taddr8)
  %i.addr9 = alloca i64
  store i64 0, i64* %i.addr9, align 8
  %taddr11 = alloca { i8**, i64, i64 } , align 8
  %taddr13 = alloca { i8**, i64, i64 } , align 8
  %taddr19 = alloca { i8**, i64, i64 } , align 8
  %taddr31 = alloca { i8*, i64 } , align 8
  %taddr32 = alloca { i8*, i64 } , align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  %taddr35 = alloca { i8*, i64 } , align 8
  %taddr36 = alloca { i8*, i64 } , align 8
  %taddr37 = alloca { i8*, i64 } , align 8
  %taddr39 = alloca { i8*, i64 } , align 8
  %taddr42 = alloca i8* , align 8
  %taddr47 = alloca i8* , align 8
  %taddr50 = alloca { i8*, i64 } , align 8
  %taddr51 = alloca { i8*, i64 } , align 8
  %taddr52 = alloca { i8*, i64 } , align 8
  %taddr54 = alloca { i8*, i64 } , align 8
  %taddr55 = alloca { i8*, i64 } , align 8
  %taddr56 = alloca { i8*, i64 } , align 8
  %taddr58 = alloca { i8*, i64 } , align 8
  %taddr61 = alloca { i8*, i64 } , align 8
  br label %for.cond0
for.cond0:
  %t10 = load i64, i64* %i.addr9, align 8
  store { i8**, i64, i64 } zeroinitializer, { i8**, i64, i64 }* %taddr11, align 8
  %t12 = load { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr11, align 8
  store { i8**, i64, i64 } %t12, { i8**, i64, i64 }* %taddr13, align 8
  %t14 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr13, i32 0, i32 1
  %t15 = load i64, i64* %t14, align 8
  %t16 = icmp slt i64 %t10, %t15
  br i1 %t16, label %for.body1, label %for.end3
for.body1:
  %t17 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t18 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t17, i32 0, i32 0
  store { i8**, i64, i64 } zeroinitializer, { i8**, i64, i64 }* %taddr19, align 8
  %t20 = load i64, i64* %i.addr9, align 8
  %t21 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr19, i32 0, i32 1
  %t22 = load i64, i64* %t21, align 8
  %t23 = icmp sge i64 %t20, 0
  %t24 = icmp slt i64 %t20, %t22
  %t25 = and i1 %t23, %t24
  br i1 %t25, label %idx.ok4, label %idx.fail5
for.post2:
  %t40 = load i64, i64* %i.addr9, align 8
  %t41 = add i64 %t40, 1
  store i64 %t41, i64* %i.addr9, align 8
  br label %for.cond0
for.end3:
  store i8* null, i8** %taddr42, align 8
  %t43 = load i8*, i8** %taddr42, align 8
  %t44 = icmp ne i8* %t43, null
  br i1 %t44, label %then6, label %else8
idx.ok4:
  %t26 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr19, i32 0, i32 0
  %t27 = load i8**, i8*** %t26, align 8
  %t28 = getelementptr inbounds i8*, i8** %t27, i64 %t20
  %t29 = load i8*, i8** %t28, align 8
  %t30 = call { i8*, i64 } @ir.Instruction.String(i8* %t29)
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.40, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr32, align 8
  store { i8*, i64 } %t30, { i8*, i64 }* %taddr33, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr31, { i8*, i64 }* %taddr32, { i8*, i64 }* %taddr33)
  %t34 = load { i8*, i64 }, { i8*, i64 }* %taddr31, align 8
  store { i8*, i64 } %t34, { i8*, i64 }* %taddr36, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.41, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr37, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr35, { i8*, i64 }* %taddr36, { i8*, i64 }* %taddr37)
  %t38 = load { i8*, i64 }, { i8*, i64 }* %taddr35, align 8
  store { i8*, i64 } %t38, { i8*, i64 }* %taddr39
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t18, { i8*, i64 }* %taddr39)
  br label %for.post2
idx.fail5:
  call void @gominic_abort()
  br label %idx.ok4
then6:
  %t45 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t46 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t45, i32 0, i32 0
  store i8* null, i8** %taddr47, align 8
  %t48 = load i8*, i8** %taddr47, align 8
  %t49 = call { i8*, i64 } @ir.Instruction.String(i8* %t48)
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.42, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr51, align 8
  store { i8*, i64 } %t49, { i8*, i64 }* %taddr52, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr50, { i8*, i64 }* %taddr51, { i8*, i64 }* %taddr52)
  %t53 = load { i8*, i64 }, { i8*, i64 }* %taddr50, align 8
  store { i8*, i64 } %t53, { i8*, i64 }* %taddr55, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.43, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr56, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr54, { i8*, i64 }* %taddr55, { i8*, i64 }* %taddr56)
  %t57 = load { i8*, i64 }, { i8*, i64 }* %taddr54, align 8
  store { i8*, i64 } %t57, { i8*, i64 }* %taddr58
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t46, { i8*, i64 }* %taddr58)
  br label %endif7
endif7:
  ret void
else8:
  %t59 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t60 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t59, i32 0, i32 0
  store { i8*, i64 } { i8* getelementptr inbounds ([15 x i8], [15 x i8]* @.str.44, i32 0, i32 0), i64 14 }, { i8*, i64 }* %taddr61
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t60, { i8*, i64 }* %taddr61)
  br label %endif7
}

define { i8*, i64 } @backend.llvmType(i8* %t) {
entry:
  %t.addr = alloca i8* , align 8
  store i8* %t, i8** %t.addr
  %t0 = load i8*, i8** %t.addr, align 8
  %t1 = icmp eq i8* %t0, null
  br i1 %t1, label %then0, label %endif1
then0:
  ret { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.45, i32 0, i32 0), i64 4 }
endif1:
  %t2 = call { i8*, i64 } @t.String()
  ret { i8*, i64 } %t2
}

define void @backend.emitTargetHeader({ { { { i8*, i64 }*, i64, i64 } }, i1 }* %e, i8** %mod) {
entry:
  %e.addr = alloca { { { { i8*, i64 }*, i64, i64 } }, i1 }* , align 8
  store { { { { i8*, i64 }*, i64, i64 } }, i1 }* %e, { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr
  %mod.addr = alloca i8** , align 8
  store i8** %mod, i8*** %mod.addr
  %taddr0 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } zeroinitializer, { i8*, i64 }* %taddr0, align 8
  %t1 = load { i8*, i64 }, { i8*, i64 }* %taddr0, align 8
  %triple.addr2 = alloca { i8*, i64 }
  %taddr3 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t1, { i8*, i64 }* %taddr3
  call void @gominic_memcpy({ i8*, i64 }* %triple.addr2, { i8*, i64 }* %taddr3, i64 16)
  %t4 = load { i8*, i64 }, { i8*, i64 }* %triple.addr2, align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t4, { i8*, i64 }* %taddr5, align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.46, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr6, align 8
  %t7 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr5, i32 0, i32 0
  %t8 = load i8*, i8** %t7, align 8
  %t9 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr5, i32 0, i32 1
  %t10 = load i64, i64* %t9, align 8
  %t11 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr6, i32 0, i32 0
  %t12 = load i8*, i8** %t11, align 8
  %t13 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr6, i32 0, i32 1
  %t14 = load i64, i64* %t13, align 8
  %t15 = call i1 @gominic_str_eq(i8* %t8, i64 %t10, i8* %t12, i64 %t14)
  %taddr16 = alloca { i8*, i64 } , align 8
  %taddr18 = alloca { i8*, i64 } , align 8
  %taddr19 = alloca { i8*, i64 } , align 8
  %taddr23 = alloca { i8*, i64 } , align 8
  %taddr26 = alloca { i8*, i64 } , align 8
  %taddr28 = alloca { i8*, i64 } , align 8
  %taddr29 = alloca { i8*, i64 } , align 8
  %taddr39 = alloca { i8*, i64 } , align 8
  %taddr43 = alloca { i8*, i64 } , align 8
  %taddr44 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca { i8*, i64 } , align 8
  %taddr47 = alloca { i8*, i64 } , align 8
  %taddr48 = alloca { i8*, i64 } , align 8
  %taddr49 = alloca { i8*, i64 } , align 8
  %taddr51 = alloca { i8*, i64 } , align 8
  %taddr55 = alloca { i8*, i64 } , align 8
  %taddr56 = alloca { i8*, i64 } , align 8
  %taddr57 = alloca { i8*, i64 } , align 8
  %taddr59 = alloca { i8*, i64 } , align 8
  %taddr60 = alloca { i8*, i64 } , align 8
  %taddr61 = alloca { i8*, i64 } , align 8
  %taddr63 = alloca { i8*, i64 } , align 8
  br i1 %t15, label %then0, label %endif1
then0:
  store { i8*, i64 } { i8* getelementptr inbounds ([20 x i8], [20 x i8]* @.str.47, i32 0, i32 0), i64 19 }, { i8*, i64 }* %taddr16
  call void @gominic_memcpy({ i8*, i64 }* %triple.addr2, { i8*, i64 }* %taddr16, i64 16)
  br label %endif1
endif1:
  %t17 = load { i8*, i64 }, { i8*, i64 }* %triple.addr2, align 8
  store { i8*, i64 } %t17, { i8*, i64 }* %taddr18
  store { i8*, i64 } { i8* getelementptr inbounds ([13 x i8], [13 x i8]* @.str.48, i32 0, i32 0), i64 12 }, { i8*, i64 }* %taddr19
  %t20 = call i1 @strings.Contains({ i8*, i64 }* %taddr18, { i8*, i64 }* %taddr19)
  %t21 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t22 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t21, i32 0, i32 1
  store i1 %t20, i1* %t22, align 1
  store { i8*, i64 } zeroinitializer, { i8*, i64 }* %taddr23, align 8
  %t24 = load { i8*, i64 }, { i8*, i64 }* %taddr23, align 8
  %dl.addr25 = alloca { i8*, i64 }
  store { i8*, i64 } %t24, { i8*, i64 }* %taddr26
  call void @gominic_memcpy({ i8*, i64 }* %dl.addr25, { i8*, i64 }* %taddr26, i64 16)
  %t27 = load { i8*, i64 }, { i8*, i64 }* %dl.addr25, align 8
  store { i8*, i64 } %t27, { i8*, i64 }* %taddr28, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.49, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr29, align 8
  %t30 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr28, i32 0, i32 0
  %t31 = load i8*, i8** %t30, align 8
  %t32 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr28, i32 0, i32 1
  %t33 = load i64, i64* %t32, align 8
  %t34 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr29, i32 0, i32 0
  %t35 = load i8*, i8** %t34, align 8
  %t36 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr29, i32 0, i32 1
  %t37 = load i64, i64* %t36, align 8
  %t38 = call i1 @gominic_str_eq(i8* %t31, i64 %t33, i8* %t35, i64 %t37)
  br i1 %t38, label %then2, label %endif3
then2:
  store { i8*, i64 } { i8* getelementptr inbounds ([71 x i8], [71 x i8]* @.str.50, i32 0, i32 0), i64 70 }, { i8*, i64 }* %taddr39
  call void @gominic_memcpy({ i8*, i64 }* %dl.addr25, { i8*, i64 }* %taddr39, i64 16)
  br label %endif3
endif3:
  %t40 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t41 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t40, i32 0, i32 0
  %t42 = load { i8*, i64 }, { i8*, i64 }* %triple.addr2, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([18 x i8], [18 x i8]* @.str.51, i32 0, i32 0), i64 17 }, { i8*, i64 }* %taddr44, align 8
  store { i8*, i64 } %t42, { i8*, i64 }* %taddr45, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr43, { i8*, i64 }* %taddr44, { i8*, i64 }* %taddr45)
  %t46 = load { i8*, i64 }, { i8*, i64 }* %taddr43, align 8
  store { i8*, i64 } %t46, { i8*, i64 }* %taddr48, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.52, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr49, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr47, { i8*, i64 }* %taddr48, { i8*, i64 }* %taddr49)
  %t50 = load { i8*, i64 }, { i8*, i64 }* %taddr47, align 8
  store { i8*, i64 } %t50, { i8*, i64 }* %taddr51
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t41, { i8*, i64 }* %taddr51)
  %t52 = bitcast { { { { i8*, i64 }*, i64, i64 } }, i1 }** %e.addr to { { { { i8*, i64 }*, i64, i64 } }, i1 }*
  %t53 = getelementptr inbounds { { { { i8*, i64 }*, i64, i64 } }, i1 }, { { { { i8*, i64 }*, i64, i64 } }, i1 }* %t52, i32 0, i32 0
  %t54 = load { i8*, i64 }, { i8*, i64 }* %dl.addr25, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([22 x i8], [22 x i8]* @.str.53, i32 0, i32 0), i64 21 }, { i8*, i64 }* %taddr56, align 8
  store { i8*, i64 } %t54, { i8*, i64 }* %taddr57, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr55, { i8*, i64 }* %taddr56, { i8*, i64 }* %taddr57)
  %t58 = load { i8*, i64 }, { i8*, i64 }* %taddr55, align 8
  store { i8*, i64 } %t58, { i8*, i64 }* %taddr60, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.54, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr61, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr59, { i8*, i64 }* %taddr60, { i8*, i64 }* %taddr61)
  %t62 = load { i8*, i64 }, { i8*, i64 }* %taddr59, align 8
  store { i8*, i64 } %t62, { i8*, i64 }* %taddr63
  call void @backend.strBufWriteString({ { { i8*, i64 }*, i64, i64 } }* %t53, { i8*, i64 }* %taddr63)
  ret void
}

define { i8*, i64 } @ir.BasicType.String({ i8*, i64 } %recv) {
entry:
  %recv.addr = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %recv, { i8*, i64 }* %recv.addr
  %t0 = load { i8*, i64 }, { i8*, i64 }* @ir.t, align 8
  ret { i8*, i64 } %t0
}

define { i8*, i64 } @ir.PointerType.String({ i8* } %recv) {
entry:
  %recv.addr = alloca { i8* } , align 8
  store { i8* } %recv, { i8* }* %recv.addr
  %t0 = bitcast { i8* }* @ir.p to { i8* }*
  %t1 = getelementptr inbounds { i8* }, { i8* }* %t0, i32 0, i32 0
  %t2 = load i8*, i8** %t1, align 8
  %t3 = call { i8*, i64 } @ir.Type.String(i8* %t2)
  %taddr4 = alloca { i8*, i64 } , align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t3, { i8*, i64 }* %taddr5, align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.98, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr6, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr4, { i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6)
  %t7 = load { i8*, i64 }, { i8*, i64 }* %taddr4, align 8
  ret { i8*, i64 } %t7
}

define { i8* } @ir.PtrTo(i8* %elem) {
entry:
  %elem.addr = alloca i8* , align 8
  store i8* %elem, i8** %elem.addr
  %t0 = load i8*, i8** %elem.addr, align 8
  %taddr1 = alloca { i8* } , align 8
  %t2 = getelementptr inbounds { i8* }, { i8* }* %taddr1, i32 0, i32 0
  store i8* %t0, i8** %t2, align 8
  %t3 = load { i8* }, { i8* }* %taddr1, align 8
  ret { i8* } %t3
}

define { i8*, i64 } @ir.StructType.String({ { i8**, i64, i64 } } %recv) {
entry:
  %recv.addr = alloca { { i8**, i64, i64 } } , align 8
  store { { i8**, i64, i64 } } %recv, { { i8**, i64, i64 } }* %recv.addr
  %t0 = bitcast { { i8**, i64, i64 } }* @ir.s to { { i8**, i64, i64 } }*
  %t1 = getelementptr inbounds { { i8**, i64, i64 } }, { { i8**, i64, i64 } }* %t0, i32 0, i32 0
  %t2 = load { i8**, i64, i64 }, { i8**, i64, i64 }* %t1, align 8
  %taddr3 = alloca { i8**, i64, i64 } , align 8
  store { i8**, i64, i64 } %t2, { i8**, i64, i64 }* %taddr3, align 8
  %t4 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr3, i32 0, i32 1
  %t5 = load i64, i64* %t4, align 8
  %t6 = call i8* @gominic_makeSlice(i64 %t5, i64 %t5, i64 16)
  %t7 = bitcast i8* %t6 to { i8*, i64 }*
  %taddr8 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %t9 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 0
  store { i8*, i64 }* %t7, { i8*, i64 }** %t9, align 8
  %t10 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 1
  store i64 %t5, i64* %t10, align 8
  %t11 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 2
  store i64 %t5, i64* %t11, align 8
  %t12 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, align 8
  %parts.addr13 = alloca { { i8*, i64 }*, i64, i64 }
  %taddr14 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  store { { i8*, i64 }*, i64, i64 } %t12, { { i8*, i64 }*, i64, i64 }* %taddr14
  call void @gominic_memcpy({ { i8*, i64 }*, i64, i64 }* %parts.addr13, { { i8*, i64 }*, i64, i64 }* %taddr14, i64 24)
  %i.addr15 = alloca i64
  store i64 0, i64* %i.addr15, align 8
  %taddr20 = alloca { i8**, i64, i64 } , align 8
  %taddr46 = alloca { i8*, i64 } , align 8
  %taddr50 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %taddr51 = alloca { i8*, i64 } , align 8
  %taddr53 = alloca { i8*, i64 } , align 8
  %taddr54 = alloca { i8*, i64 } , align 8
  %taddr55 = alloca { i8*, i64 } , align 8
  %taddr57 = alloca { i8*, i64 } , align 8
  %taddr58 = alloca { i8*, i64 } , align 8
  %taddr59 = alloca { i8*, i64 } , align 8
  br label %for.cond0
for.cond0:
  %t16 = load i64, i64* %i.addr15, align 8
  %t17 = bitcast { { i8**, i64, i64 } }* @ir.s to { { i8**, i64, i64 } }*
  %t18 = getelementptr inbounds { { i8**, i64, i64 } }, { { i8**, i64, i64 } }* %t17, i32 0, i32 0
  %t19 = load { i8**, i64, i64 }, { i8**, i64, i64 }* %t18, align 8
  store { i8**, i64, i64 } %t19, { i8**, i64, i64 }* %taddr20, align 8
  %t21 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr20, i32 0, i32 1
  %t22 = load i64, i64* %t21, align 8
  %t23 = icmp slt i64 %t16, %t22
  br i1 %t23, label %for.body1, label %for.end3
for.body1:
  %t24 = bitcast { { i8**, i64, i64 } }* @ir.s to { { i8**, i64, i64 } }*
  %t25 = getelementptr inbounds { { i8**, i64, i64 } }, { { i8**, i64, i64 } }* %t24, i32 0, i32 0
  %t26 = load i64, i64* %i.addr15, align 8
  %t27 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %t25, i32 0, i32 1
  %t28 = load i64, i64* %t27, align 8
  %t29 = icmp sge i64 %t26, 0
  %t30 = icmp slt i64 %t26, %t28
  %t31 = and i1 %t29, %t30
  br i1 %t31, label %idx.ok4, label %idx.fail5
for.post2:
  %t47 = load i64, i64* %i.addr15, align 8
  %t48 = add i64 %t47, 1
  store i64 %t48, i64* %i.addr15, align 8
  br label %for.cond0
for.end3:
  %t49 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %parts.addr13, align 8
  store { { i8*, i64 }*, i64, i64 } %t49, { { i8*, i64 }*, i64, i64 }* %taddr50
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.100, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr51
  %t52 = call { i8*, i64 } @strings.Join({ { i8*, i64 }*, i64, i64 }* %taddr50, { i8*, i64 }* %taddr51)
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.99, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr54, align 8
  store { i8*, i64 } %t52, { i8*, i64 }* %taddr55, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr53, { i8*, i64 }* %taddr54, { i8*, i64 }* %taddr55)
  %t56 = load { i8*, i64 }, { i8*, i64 }* %taddr53, align 8
  store { i8*, i64 } %t56, { i8*, i64 }* %taddr58, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.101, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr59, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr57, { i8*, i64 }* %taddr58, { i8*, i64 }* %taddr59)
  %t60 = load { i8*, i64 }, { i8*, i64 }* %taddr57, align 8
  ret { i8*, i64 } %t60
idx.ok4:
  %t32 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %t25, i32 0, i32 0
  %t33 = load i8**, i8*** %t32, align 8
  %t34 = getelementptr inbounds i8*, i8** %t33, i64 %t26
  %t35 = load i8*, i8** %t34, align 8
  %t36 = call { i8*, i64 } @ir.Type.String(i8* %t35)
  %t37 = load i64, i64* %i.addr15, align 8
  %t38 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %parts.addr13, i32 0, i32 1
  %t39 = load i64, i64* %t38, align 8
  %t40 = icmp sge i64 %t37, 0
  %t41 = icmp slt i64 %t37, %t39
  %t42 = and i1 %t40, %t41
  br i1 %t42, label %idx.ok6, label %idx.fail7
idx.fail5:
  call void @gominic_abort()
  br label %idx.ok4
idx.ok6:
  %t43 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %parts.addr13, i32 0, i32 0
  %t44 = load { i8*, i64 }*, { i8*, i64 }** %t43, align 8
  %t45 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %t44, i64 %t37
  store { i8*, i64 } %t36, { i8*, i64 }* %taddr46
  call void @gominic_memcpy({ i8*, i64 }* %t45, { i8*, i64 }* %taddr46, i64 16)
  br label %for.post2
idx.fail7:
  call void @gominic_abort()
  br label %idx.ok6
}

define { i8*, i64 } @ir.ArrayType.String({ i64, i8* } %recv) {
entry:
  %recv.addr = alloca { i64, i8* } , align 8
  store { i64, i8* } %recv, { i64, i8* }* %recv.addr
  %t0 = bitcast { i64, i8* }* @ir.a to { i64, i8* }*
  %t1 = getelementptr inbounds { i64, i8* }, { i64, i8* }* %t0, i32 0, i32 0
  %t2 = load i64, i64* %t1, align 8
  %t3 = call { i8*, i64 } @strconv.Itoa(i64 %t2)
  %taddr4 = alloca { i8*, i64 } , align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.102, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr5, align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t3, { i8*, i64 }* %taddr6, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr4, { i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6)
  %t7 = load { i8*, i64 }, { i8*, i64 }* %taddr4, align 8
  %taddr8 = alloca { i8*, i64 } , align 8
  %taddr9 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t7, { i8*, i64 }* %taddr9, align 8
  %taddr10 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.103, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr10, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr8, { i8*, i64 }* %taddr9, { i8*, i64 }* %taddr10)
  %t11 = load { i8*, i64 }, { i8*, i64 }* %taddr8, align 8
  %t12 = bitcast { i64, i8* }* @ir.a to { i64, i8* }*
  %t13 = getelementptr inbounds { i64, i8* }, { i64, i8* }* %t12, i32 0, i32 1
  %t14 = load i8*, i8** %t13, align 8
  %t15 = call { i8*, i64 } @ir.Type.String(i8* %t14)
  %taddr16 = alloca { i8*, i64 } , align 8
  %taddr17 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t11, { i8*, i64 }* %taddr17, align 8
  %taddr18 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t15, { i8*, i64 }* %taddr18, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr16, { i8*, i64 }* %taddr17, { i8*, i64 }* %taddr18)
  %t19 = load { i8*, i64 }, { i8*, i64 }* %taddr16, align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t19, { i8*, i64 }* %taddr21, align 8
  %taddr22 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.104, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr22, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21, { i8*, i64 }* %taddr22)
  %t23 = load { i8*, i64 }, { i8*, i64 }* %taddr20, align 8
  ret { i8*, i64 } %t23
}

define { i8*, i64 } @ir.SliceType.String({ i8* } %recv) {
entry:
  %recv.addr = alloca { i8* } , align 8
  store { i8* } %recv, { i8* }* %recv.addr
  %t0 = bitcast { i8* }* @ir.s to { i8* }*
  %t1 = getelementptr inbounds { i8* }, { i8* }* %t0, i32 0, i32 0
  %t2 = load i8*, i8** %t1, align 8
  %t3 = call { i8* } @ir.PtrTo(i8* %t2)
  %t4 = call { i8*, i64 } @ir.PointerType.String({ i8* } %t3)
  %taddr5 = alloca { i8*, i64 } , align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.105, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr6, align 8
  %taddr7 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t4, { i8*, i64 }* %taddr7, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6, { i8*, i64 }* %taddr7)
  %t8 = load { i8*, i64 }, { i8*, i64 }* %taddr5, align 8
  %taddr9 = alloca { i8*, i64 } , align 8
  %taddr10 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t8, { i8*, i64 }* %taddr10, align 8
  %taddr11 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([13 x i8], [13 x i8]* @.str.106, i32 0, i32 0), i64 12 }, { i8*, i64 }* %taddr11, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr9, { i8*, i64 }* %taddr10, { i8*, i64 }* %taddr11)
  %t12 = load { i8*, i64 }, { i8*, i64 }* %taddr9, align 8
  ret { i8*, i64 } %t12
}

define { i8*, i64 } @ir.StringType.String({  } %recv) {
entry:
  %recv.addr = alloca {  } , align 1
  store {  } %recv, {  }* %recv.addr
  ret { i8*, i64 } { i8* getelementptr inbounds ([13 x i8], [13 x i8]* @.str.107, i32 0, i32 0), i64 12 }
}

define i8* @ir.valueData.Type({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %recv) {
entry:
  %recv.addr = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %recv, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %recv.addr
  %t0 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v, align 8
  %t1 = icmp eq { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t0, null
  br i1 %t1, label %then0, label %endif1
then0:
  ret i8* null
endif1:
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 1
  %t4 = load i8*, i8** %t3, align 8
  ret i8* %t4
}

define { i8*, i64 } @ir.valueData.Name({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %recv) {
entry:
  %recv.addr = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %recv, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %recv.addr
  %t0 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v, align 8
  %t1 = icmp eq { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t0, null
  %taddr9 = alloca { i8*, i64 } , align 8
  %taddr10 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca i1 , align 1
  %taddr25 = alloca { i8*, i64 } , align 8
  %taddr26 = alloca { i8*, i64 } , align 8
  br i1 %t1, label %then0, label %endif1
then0:
  ret { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.108, i32 0, i32 0), i64 0 }
endif1:
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 2
  %t4 = load i64, i64* %t3, align 8
  %t5 = icmp eq i64 %t4, 3
  %t6 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t7 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t6, i32 0, i32 3
  %t8 = load { i8*, i64 }, { i8*, i64 }* %t7, align 8
  store { i8*, i64 } %t8, { i8*, i64 }* %taddr9, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.109, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr10, align 8
  %t11 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr9, i32 0, i32 0
  %t12 = load i8*, i8** %t11, align 8
  %t13 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr9, i32 0, i32 1
  %t14 = load i64, i64* %t13, align 8
  %t15 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr10, i32 0, i32 0
  %t16 = load i8*, i8** %t15, align 8
  %t17 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr10, i32 0, i32 1
  %t18 = load i64, i64* %t17, align 8
  %t19 = call i1 @gominic_str_eq(i8* %t12, i64 %t14, i8* %t16, i64 %t18)
  %t20 = icmp eq i1 %t19, 0
  br i1 %t5, label %logic.rhs3, label %logic.false4
logic.end2:
  %t37 = load i1, i1* %taddr21, align 1
  br i1 %t37, label %then5, label %endif6
logic.rhs3:
  %t22 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t23 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t22, i32 0, i32 3
  %t24 = load { i8*, i64 }, { i8*, i64 }* %t23, align 8
  store { i8*, i64 } %t24, { i8*, i64 }* %taddr25, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.110, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr26, align 8
  %t27 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr25, i32 0, i32 0
  %t28 = load i8*, i8** %t27, align 8
  %t29 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr25, i32 0, i32 1
  %t30 = load i64, i64* %t29, align 8
  %t31 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr26, i32 0, i32 0
  %t32 = load i8*, i8** %t31, align 8
  %t33 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr26, i32 0, i32 1
  %t34 = load i64, i64* %t33, align 8
  %t35 = call i1 @gominic_str_eq(i8* %t28, i64 %t30, i8* %t32, i64 %t34)
  %t36 = icmp eq i1 %t35, 0
  store i1 %t36, i1* %taddr21
  br label %logic.end2
logic.false4:
  store i1 0, i1* %taddr21
  br label %logic.end2
then5:
  %t38 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t39 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t38, i32 0, i32 3
  %t40 = load { i8*, i64 }, { i8*, i64 }* %t39, align 8
  ret { i8*, i64 } %t40
endif6:
  %t41 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t42 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t41, i32 0, i32 0
  %t43 = load { i8*, i64 }, { i8*, i64 }* %t42, align 8
  ret { i8*, i64 } %t43
}

define i1 @ir.valueData.ByVal({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %recv) {
entry:
  %recv.addr = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %recv, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %recv.addr
  %t0 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v, align 8
  %t1 = icmp eq { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t0, null
  br i1 %t1, label %then0, label %endif1
then0:
  ret i1 0
endif1:
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 4
  %t4 = load i1, i1* %t3, align 1
  ret i1 %t4
}

define i8* @ir.valueData.ByValType({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %recv) {
entry:
  %recv.addr = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %recv, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %recv.addr
  %t0 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v, align 8
  %t1 = icmp eq { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t0, null
  br i1 %t1, label %then0, label %endif1
then0:
  ret i8* null
endif1:
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 5
  %t4 = load i8*, i8** %t3, align 8
  ret i8* %t4
}

define i64 @ir.valueData.Kind({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %recv) {
entry:
  %recv.addr = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %recv, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %recv.addr
  %t0 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v, align 8
  %t1 = icmp eq { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t0, null
  br i1 %t1, label %then0, label %endif1
then0:
  ret i64 0
endif1:
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** @ir.v to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 2
  %t4 = load i64, i64* %t3, align 8
  ret i64 %t4
}

define { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* @ir.NewParam({ i8*, i64 }* byval({ i8*, i64 }) %name, i8* %ty) {
entry:
  %name.addr = alloca { i8*, i64 } , align 8
  call void @gominic_memcpy({ i8*, i64 }* %name.addr, { i8*, i64 }* %name, i64 16)
  %ty.addr = alloca i8* , align 8
  store i8* %ty, i8** %ty.addr
  %t0 = load { i8*, i64 }, { i8*, i64 }* %name.addr, align 8
  %t1 = load i8*, i8** %ty.addr, align 8
  %taddr2 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } , align 8
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 0
  store { i8*, i64 } %t0, { i8*, i64 }* %t3, align 8
  %t4 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 1
  store i8* %t1, i8** %t4, align 8
  %t5 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 2
  store i64 1, i64* %t5, align 8
  %t6 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 3
  store { i8*, i64 } zeroinitializer, { i8*, i64 }* %t6, align 8
  %t7 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 4
  store i1 0, i1* %t7, align 1
  %t8 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 5
  store i8* null, i8** %t8, align 8
  %t9 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, align 8
  %taddr10 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } %t9, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr10, align 8
  ret { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr10
}

define { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* @ir.NewByValParam({ i8*, i64 }* byval({ i8*, i64 }) %name, i8* %ptr, i8* %byValTyp) {
entry:
  %name.addr = alloca { i8*, i64 } , align 8
  call void @gominic_memcpy({ i8*, i64 }* %name.addr, { i8*, i64 }* %name, i64 16)
  %ptr.addr = alloca i8* , align 8
  store i8* %ptr, i8** %ptr.addr
  %byValTyp.addr = alloca i8* , align 8
  store i8* %byValTyp, i8** %byValTyp.addr
  %t0 = load { i8*, i64 }, { i8*, i64 }* %name.addr, align 8
  %t1 = load i8*, i8** %ptr.addr, align 8
  %t2 = load i8*, i8** %byValTyp.addr, align 8
  %taddr3 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } , align 8
  %t4 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 0
  store { i8*, i64 } %t0, { i8*, i64 }* %t4, align 8
  %t5 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 1
  store i8* %t1, i8** %t5, align 8
  %t6 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 2
  store i64 1, i64* %t6, align 8
  %t7 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 3
  store { i8*, i64 } zeroinitializer, { i8*, i64 }* %t7, align 8
  %t8 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 4
  store i1 1, i1* %t8, align 1
  %t9 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 5
  store i8* %t2, i8** %t9, align 8
  %t10 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, align 8
  %taddr11 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } %t10, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr11, align 8
  ret { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr11
}

define { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* @ir.NewRegister({ i8*, i64 }* byval({ i8*, i64 }) %name, i8* %ty) {
entry:
  %name.addr = alloca { i8*, i64 } , align 8
  call void @gominic_memcpy({ i8*, i64 }* %name.addr, { i8*, i64 }* %name, i64 16)
  %ty.addr = alloca i8* , align 8
  store i8* %ty, i8** %ty.addr
  %t0 = load { i8*, i64 }, { i8*, i64 }* %name.addr, align 8
  %t1 = load i8*, i8** %ty.addr, align 8
  %taddr2 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } , align 8
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 0
  store { i8*, i64 } %t0, { i8*, i64 }* %t3, align 8
  %t4 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 1
  store i8* %t1, i8** %t4, align 8
  %t5 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 2
  store i64 2, i64* %t5, align 8
  %t6 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 3
  store { i8*, i64 } zeroinitializer, { i8*, i64 }* %t6, align 8
  %t7 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 4
  store i1 0, i1* %t7, align 1
  %t8 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, i32 0, i32 5
  store i8* null, i8** %t8, align 8
  %t9 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr2, align 8
  %taddr10 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } %t9, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr10, align 8
  ret { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr10
}

define { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* @ir.NewConstant({ i8*, i64 }* byval({ i8*, i64 }) %raw, i8* %ty) {
entry:
  %raw.addr = alloca { i8*, i64 } , align 8
  call void @gominic_memcpy({ i8*, i64 }* %raw.addr, { i8*, i64 }* %raw, i64 16)
  %ty.addr = alloca i8* , align 8
  store i8* %ty, i8** %ty.addr
  %t0 = load { i8*, i64 }, { i8*, i64 }* %raw.addr, align 8
  %t1 = load { i8*, i64 }, { i8*, i64 }* %raw.addr, align 8
  %t2 = load i8*, i8** %ty.addr, align 8
  %taddr3 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } , align 8
  %t4 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 0
  store { i8*, i64 } %t0, { i8*, i64 }* %t4, align 8
  %t5 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 1
  store i8* %t2, i8** %t5, align 8
  %t6 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 2
  store i64 3, i64* %t6, align 8
  %t7 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 3
  store { i8*, i64 } %t1, { i8*, i64 }* %t7, align 8
  %t8 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 4
  store i1 0, i1* %t8, align 1
  %t9 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, i32 0, i32 5
  store i8* null, i8** %t9, align 8
  %t10 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr3, align 8
  %taddr11 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* } %t10, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr11, align 8
  ret { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %taddr11
}

define void @ir.BinOp.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.BinOp.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %recv.addr
  %t0 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.i to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t1 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t0, i32 0, i32 0
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t1 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 0
  %t4 = load { i8*, i64 }, { i8*, i64 }* %t3, align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.111, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr6, align 8
  %taddr7 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t4, { i8*, i64 }* %taddr7, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6, { i8*, i64 }* %taddr7)
  %t8 = load { i8*, i64 }, { i8*, i64 }* %taddr5, align 8
  %taddr9 = alloca { i8*, i64 } , align 8
  %taddr10 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t8, { i8*, i64 }* %taddr10, align 8
  %taddr11 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.112, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr11, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr9, { i8*, i64 }* %taddr10, { i8*, i64 }* %taddr11)
  %t12 = load { i8*, i64 }, { i8*, i64 }* %taddr9, align 8
  %t13 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.i to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t14 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t13, i32 0, i32 1
  %t15 = load { i8*, i64 }, { i8*, i64 }* %t14, align 8
  %taddr16 = alloca { i8*, i64 } , align 8
  %taddr17 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t12, { i8*, i64 }* %taddr17, align 8
  %taddr18 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t15, { i8*, i64 }* %taddr18, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr16, { i8*, i64 }* %taddr17, { i8*, i64 }* %taddr18)
  %t19 = load { i8*, i64 }, { i8*, i64 }* %taddr16, align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t19, { i8*, i64 }* %taddr21, align 8
  %taddr22 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.113, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr22, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21, { i8*, i64 }* %taddr22)
  %t23 = load { i8*, i64 }, { i8*, i64 }* %taddr20, align 8
  %t24 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.i to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t25 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t24, i32 0, i32 2
  %t26 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t25, align 8
  %t27 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t26)
  %taddr28 = alloca { i8*, i64 } , align 8
  %taddr29 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t23, { i8*, i64 }* %taddr29, align 8
  %taddr30 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t27, { i8*, i64 }* %taddr30, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr28, { i8*, i64 }* %taddr29, { i8*, i64 }* %taddr30)
  %t31 = load { i8*, i64 }, { i8*, i64 }* %taddr28, align 8
  %taddr32 = alloca { i8*, i64 } , align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t31, { i8*, i64 }* %taddr33, align 8
  %taddr34 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.114, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr34, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr32, { i8*, i64 }* %taddr33, { i8*, i64 }* %taddr34)
  %t35 = load { i8*, i64 }, { i8*, i64 }* %taddr32, align 8
  %t36 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.i to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t37 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t36, i32 0, i32 2
  %t38 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t37, align 8
  %t39 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t38)
  %taddr40 = alloca { i8*, i64 } , align 8
  %taddr41 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t35, { i8*, i64 }* %taddr41, align 8
  %taddr42 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t39, { i8*, i64 }* %taddr42, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr40, { i8*, i64 }* %taddr41, { i8*, i64 }* %taddr42)
  %t43 = load { i8*, i64 }, { i8*, i64 }* %taddr40, align 8
  %taddr44 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t43, { i8*, i64 }* %taddr45, align 8
  %taddr46 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.115, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr46, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr44, { i8*, i64 }* %taddr45, { i8*, i64 }* %taddr46)
  %t47 = load { i8*, i64 }, { i8*, i64 }* %taddr44, align 8
  %t48 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.i to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t49 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t48, i32 0, i32 3
  %t50 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t49, align 8
  %t51 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t50)
  %taddr52 = alloca { i8*, i64 } , align 8
  %taddr53 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t47, { i8*, i64 }* %taddr53, align 8
  %taddr54 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t51, { i8*, i64 }* %taddr54, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr52, { i8*, i64 }* %taddr53, { i8*, i64 }* %taddr54)
  %t55 = load { i8*, i64 }, { i8*, i64 }* %taddr52, align 8
  ret { i8*, i64 } %t55
}

define void @ir.Return.isInstruction({ { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } %recv) {
entry:
  %recv.addr = alloca { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } , align 8
  store { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } %recv, { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.Return.String({ { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } %recv) {
entry:
  %recv.addr = alloca { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } , align 8
  store { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } %recv, { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %recv.addr
  %t0 = bitcast { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* @ir.r to { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }*
  %t1 = getelementptr inbounds { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }, { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %t0, i32 0, i32 0
  %t2 = load { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t1, align 8
  %taddr3 = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } %t2, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr3, align 8
  %t4 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr3, i32 0, i32 1
  %t5 = load i64, i64* %t4, align 8
  %t6 = icmp eq i64 %t5, 0
  %taddr10 = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } , align 8
  %taddr28 = alloca { i8*, i64 } , align 8
  %taddr29 = alloca { i8*, i64 } , align 8
  %taddr30 = alloca { i8*, i64 } , align 8
  %taddr32 = alloca { i8*, i64 } , align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  %taddr34 = alloca { i8*, i64 } , align 8
  %taddr38 = alloca { i8*, i64 } , align 8
  %taddr39 = alloca { i8*, i64 } , align 8
  %taddr40 = alloca { i8*, i64 } , align 8
  br i1 %t6, label %then0, label %endif1
then0:
  ret { i8*, i64 } { i8* getelementptr inbounds ([9 x i8], [9 x i8]* @.str.116, i32 0, i32 0), i64 8 }
endif1:
  %t7 = bitcast { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* @ir.r to { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }*
  %t8 = getelementptr inbounds { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }, { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %t7, i32 0, i32 0
  %t9 = load { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t8, align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } %t9, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr10, align 8
  %t11 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr10, i32 0, i32 1
  %t12 = load i64, i64* %t11, align 8
  %t13 = icmp eq i64 %t12, 1
  br i1 %t13, label %then2, label %endif3
then2:
  %t14 = bitcast { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* @ir.r to { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }*
  %t15 = getelementptr inbounds { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }, { { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %t14, i32 0, i32 0
  %t16 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t15, i32 0, i32 1
  %t17 = load i64, i64* %t16, align 8
  %t18 = icmp sge i64 0, 0
  %t19 = icmp slt i64 0, %t17
  %t20 = and i1 %t18, %t19
  br i1 %t20, label %idx.ok4, label %idx.fail5
endif3:
  ret { i8*, i64 } { i8* getelementptr inbounds ([39 x i8], [39 x i8]* @.str.119, i32 0, i32 0), i64 38 }
idx.ok4:
  %t21 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t15, i32 0, i32 0
  %t22 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*** %t21, align 8
  %t23 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t22, i64 0
  %t24 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t23, align 8
  %val.addr25 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t24, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %val.addr25, align 8
  %t26 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %val.addr25, align 8
  %t27 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t26)
  store { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.117, i32 0, i32 0), i64 4 }, { i8*, i64 }* %taddr29, align 8
  store { i8*, i64 } %t27, { i8*, i64 }* %taddr30, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr28, { i8*, i64 }* %taddr29, { i8*, i64 }* %taddr30)
  %t31 = load { i8*, i64 }, { i8*, i64 }* %taddr28, align 8
  store { i8*, i64 } %t31, { i8*, i64 }* %taddr33, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.118, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr34, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr32, { i8*, i64 }* %taddr33, { i8*, i64 }* %taddr34)
  %t35 = load { i8*, i64 }, { i8*, i64 }* %taddr32, align 8
  %t36 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %val.addr25, align 8
  %t37 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t36)
  store { i8*, i64 } %t35, { i8*, i64 }* %taddr39, align 8
  store { i8*, i64 } %t37, { i8*, i64 }* %taddr40, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr38, { i8*, i64 }* %taddr39, { i8*, i64 }* %taddr40)
  %t41 = load { i8*, i64 }, { i8*, i64 }* %taddr38, align 8
  ret { i8*, i64 } %t41
idx.fail5:
  call void @gominic_abort()
  br label %idx.ok4
}

define void @ir.Call.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.Call.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %recv.addr
  %t0 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t1 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t0, i32 0, i32 2
  %t2 = load { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t1, align 8
  %taddr3 = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } %t2, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr3, align 8
  %t4 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr3, i32 0, i32 1
  %t5 = load i64, i64* %t4, align 8
  %t6 = call i8* @gominic_makeSlice(i64 %t5, i64 %t5, i64 16)
  %t7 = bitcast i8* %t6 to { i8*, i64 }*
  %taddr8 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %t9 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 0
  store { i8*, i64 }* %t7, { i8*, i64 }** %t9, align 8
  %t10 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 1
  store i64 %t5, i64* %t10, align 8
  %t11 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 2
  store i64 %t5, i64* %t11, align 8
  %t12 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, align 8
  %args.addr13 = alloca { { i8*, i64 }*, i64, i64 }
  %taddr14 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  store { { i8*, i64 }*, i64, i64 } %t12, { { i8*, i64 }*, i64, i64 }* %taddr14
  call void @gominic_memcpy({ { i8*, i64 }*, i64, i64 }* %args.addr13, { { i8*, i64 }*, i64, i64 }* %taddr14, i64 24)
  %i.addr15 = alloca i64
  store i64 0, i64* %i.addr15, align 8
  %taddr20 = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } , align 8
  %taddr38 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca i1 , align 1
  %taddr55 = alloca { i8*, i64 } , align 8
  %taddr56 = alloca { i8*, i64 } , align 8
  %taddr57 = alloca { i8*, i64 } , align 8
  %taddr59 = alloca { i8*, i64 } , align 8
  %taddr60 = alloca { i8*, i64 } , align 8
  %taddr61 = alloca { i8*, i64 } , align 8
  %taddr63 = alloca { i8*, i64 } , align 8
  %taddr65 = alloca { i8*, i64 } , align 8
  %taddr67 = alloca { i8*, i64 } , align 8
  %taddr78 = alloca { i8*, i64 } , align 8
  %taddr87 = alloca i1 , align 1
  %taddr96 = alloca { i8*, i64 } , align 8
  %taddr97 = alloca { i8*, i64 } , align 8
  %taddr98 = alloca { i8*, i64 } , align 8
  %taddr100 = alloca { i8*, i64 } , align 8
  %taddr108 = alloca { i8*, i64 } , align 8
  %taddr109 = alloca { i8*, i64 } , align 8
  %taddr120 = alloca i1 , align 1
  %taddr124 = alloca { i8*, i64 } , align 8
  %taddr125 = alloca { i8*, i64 } , align 8
  %taddr140 = alloca { i8*, i64 } , align 8
  %taddr144 = alloca { i8*, i64 } , align 8
  %taddr147 = alloca { i8*, i64 } , align 8
  %taddr148 = alloca { i8*, i64 } , align 8
  %taddr149 = alloca { i8*, i64 } , align 8
  %taddr151 = alloca { i8*, i64 } , align 8
  %taddr152 = alloca { i8*, i64 } , align 8
  %taddr153 = alloca { i8*, i64 } , align 8
  %taddr156 = alloca { i8*, i64 } , align 8
  %taddr157 = alloca { i8*, i64 } , align 8
  %taddr158 = alloca { i8*, i64 } , align 8
  %taddr169 = alloca { i8*, i64 } , align 8
  %taddr173 = alloca { i8*, i64 } , align 8
  %taddr182 = alloca { i8*, i64 } , align 8
  %taddr184 = alloca { i8*, i64 } , align 8
  %taddr194 = alloca { i8*, i64 } , align 8
  %taddr195 = alloca { i8*, i64 } , align 8
  %taddr206 = alloca i1 , align 1
  %taddr212 = alloca { i8*, i64 } , align 8
  %taddr213 = alloca { i8*, i64 } , align 8
  %taddr226 = alloca { i8*, i64 } , align 8
  %taddr227 = alloca { i8*, i64 } , align 8
  %taddr238 = alloca i1 , align 1
  %taddr240 = alloca { i8*, i64 } , align 8
  %taddr241 = alloca { i8*, i64 } , align 8
  %taddr258 = alloca { i8*, i64 } , align 8
  %taddr259 = alloca { i8*, i64 } , align 8
  %taddr260 = alloca { i8*, i64 } , align 8
  %taddr262 = alloca { i8*, i64 } , align 8
  %taddr263 = alloca { i8*, i64 } , align 8
  %taddr264 = alloca { i8*, i64 } , align 8
  %taddr266 = alloca { i8*, i64 } , align 8
  %taddr268 = alloca { i8*, i64 } , align 8
  %taddr269 = alloca { i8*, i64 } , align 8
  %taddr270 = alloca { i8*, i64 } , align 8
  %taddr273 = alloca { i8*, i64 } , align 8
  %taddr274 = alloca { i8*, i64 } , align 8
  %taddr275 = alloca { i8*, i64 } , align 8
  %taddr277 = alloca { i8*, i64 } , align 8
  %taddr278 = alloca { i8*, i64 } , align 8
  %taddr279 = alloca { i8*, i64 } , align 8
  %taddr284 = alloca { i8*, i64 } , align 8
  %taddr285 = alloca { i8*, i64 } , align 8
  %taddr286 = alloca { i8*, i64 } , align 8
  %taddr288 = alloca { i8*, i64 } , align 8
  %taddr289 = alloca { i8*, i64 } , align 8
  %taddr290 = alloca { i8*, i64 } , align 8
  %taddr293 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %taddr294 = alloca { i8*, i64 } , align 8
  %taddr296 = alloca { i8*, i64 } , align 8
  %taddr297 = alloca { i8*, i64 } , align 8
  %taddr298 = alloca { i8*, i64 } , align 8
  %taddr300 = alloca { i8*, i64 } , align 8
  %taddr301 = alloca { i8*, i64 } , align 8
  %taddr302 = alloca { i8*, i64 } , align 8
  br label %for.cond0
for.cond0:
  %t16 = load i64, i64* %i.addr15, align 8
  %t17 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t18 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t17, i32 0, i32 2
  %t19 = load { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t18, align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } %t19, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr20, align 8
  %t21 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr20, i32 0, i32 1
  %t22 = load i64, i64* %t21, align 8
  %t23 = icmp slt i64 %t16, %t22
  br i1 %t23, label %for.body1, label %for.end3
for.body1:
  %t24 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t25 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t24, i32 0, i32 2
  %t26 = load i64, i64* %i.addr15, align 8
  %t27 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t25, i32 0, i32 1
  %t28 = load i64, i64* %t27, align 8
  %t29 = icmp sge i64 %t26, 0
  %t30 = icmp slt i64 %t26, %t28
  %t31 = and i1 %t29, %t30
  br i1 %t31, label %idx.ok4, label %idx.fail5
for.post2:
  %t170 = load i64, i64* %i.addr15, align 8
  %t171 = add i64 %t170, 1
  store i64 %t171, i64* %i.addr15, align 8
  br label %for.cond0
for.end3:
  %ret.addr172 = alloca { i8*, i64 }
  store { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.129, i32 0, i32 0), i64 4 }, { i8*, i64 }* %taddr173
  call void @gominic_memcpy({ i8*, i64 }* %ret.addr172, { i8*, i64 }* %taddr173, i64 16)
  %t174 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t175 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t174, i32 0, i32 3
  %t176 = load i8*, i8** %t175, align 8
  %t177 = icmp ne i8* %t176, null
  br i1 %t177, label %then29, label %endif30
idx.ok4:
  %t32 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t25, i32 0, i32 0
  %t33 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*** %t32, align 8
  %t34 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t33, i64 %t26
  %t35 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t34, align 8
  %a.addr36 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t35, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36, align 8
  %prefix.addr37 = alloca { i8*, i64 }
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.120, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr38
  call void @gominic_memcpy({ i8*, i64 }* %prefix.addr37, { i8*, i64 }* %taddr38, i64 16)
  %t39 = load i64, i64* %i.addr15, align 8
  %t40 = icmp eq i64 %t39, 0
  %t41 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t42 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t41, i32 0, i32 4
  %t43 = load i8*, i8** %t42, align 8
  %t44 = icmp ne i8* %t43, null
  br i1 %t40, label %logic.rhs7, label %logic.false8
idx.fail5:
  call void @gominic_abort()
  br label %idx.ok4
logic.end6:
  %t50 = load i1, i1* %taddr45, align 1
  br i1 %t50, label %then9, label %endif10
logic.rhs7:
  %t46 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t47 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t46, i32 0, i32 4
  %t48 = load i8*, i8** %t47, align 8
  %t49 = icmp ne i8* %t48, null
  store i1 %t49, i1* %taddr45
  br label %logic.end6
logic.false8:
  store i1 0, i1* %taddr45
  br label %logic.end6
then9:
  %t51 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t52 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t51, i32 0, i32 4
  %t53 = load i8*, i8** %t52, align 8
  %t54 = call { i8*, i64 } @ir.Type.String(i8* %t53)
  store { i8*, i64 } { i8* getelementptr inbounds ([6 x i8], [6 x i8]* @.str.121, i32 0, i32 0), i64 5 }, { i8*, i64 }* %taddr56, align 8
  store { i8*, i64 } %t54, { i8*, i64 }* %taddr57, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr55, { i8*, i64 }* %taddr56, { i8*, i64 }* %taddr57)
  %t58 = load { i8*, i64 }, { i8*, i64 }* %taddr55, align 8
  store { i8*, i64 } %t58, { i8*, i64 }* %taddr60, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.122, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr61, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr59, { i8*, i64 }* %taddr60, { i8*, i64 }* %taddr61)
  %t62 = load { i8*, i64 }, { i8*, i64 }* %taddr59, align 8
  store { i8*, i64 } %t62, { i8*, i64 }* %taddr63
  call void @gominic_memcpy({ i8*, i64 }* %prefix.addr37, { i8*, i64 }* %taddr63, i64 16)
  br label %endif10
endif10:
  %tyStr.addr64 = alloca { i8*, i64 }
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.123, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr65
  call void @gominic_memcpy({ i8*, i64 }* %tyStr.addr64, { i8*, i64 }* %taddr65, i64 16)
  %valStr.addr66 = alloca { i8*, i64 }
  store { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.124, i32 0, i32 0), i64 4 }, { i8*, i64 }* %taddr67
  call void @gominic_memcpy({ i8*, i64 }* %valStr.addr66, { i8*, i64 }* %taddr67, i64 16)
  %t68 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36, align 8
  %t69 = icmp ne { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t68, null
  br i1 %t69, label %then11, label %endif12
then11:
  %t70 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t71 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t70, i32 0, i32 1
  %t72 = load i8*, i8** %t71, align 8
  %t73 = icmp ne i8* %t72, null
  br i1 %t73, label %then13, label %endif14
endif12:
  %t145 = load { i8*, i64 }, { i8*, i64 }* %prefix.addr37, align 8
  %t146 = load { i8*, i64 }, { i8*, i64 }* %tyStr.addr64, align 8
  store { i8*, i64 } %t145, { i8*, i64 }* %taddr148, align 8
  store { i8*, i64 } %t146, { i8*, i64 }* %taddr149, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr147, { i8*, i64 }* %taddr148, { i8*, i64 }* %taddr149)
  %t150 = load { i8*, i64 }, { i8*, i64 }* %taddr147, align 8
  store { i8*, i64 } %t150, { i8*, i64 }* %taddr152, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.128, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr153, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr151, { i8*, i64 }* %taddr152, { i8*, i64 }* %taddr153)
  %t154 = load { i8*, i64 }, { i8*, i64 }* %taddr151, align 8
  %t155 = load { i8*, i64 }, { i8*, i64 }* %valStr.addr66, align 8
  store { i8*, i64 } %t154, { i8*, i64 }* %taddr157, align 8
  store { i8*, i64 } %t155, { i8*, i64 }* %taddr158, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr156, { i8*, i64 }* %taddr157, { i8*, i64 }* %taddr158)
  %t159 = load { i8*, i64 }, { i8*, i64 }* %taddr156, align 8
  %t160 = load i64, i64* %i.addr15, align 8
  %t161 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %args.addr13, i32 0, i32 1
  %t162 = load i64, i64* %t161, align 8
  %t163 = icmp sge i64 %t160, 0
  %t164 = icmp slt i64 %t160, %t162
  %t165 = and i1 %t163, %t164
  br i1 %t165, label %idx.ok27, label %idx.fail28
then13:
  %t74 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t75 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t74, i32 0, i32 1
  %t76 = load i8*, i8** %t75, align 8
  %t77 = call { i8*, i64 } @ir.Type.String(i8* %t76)
  store { i8*, i64 } %t77, { i8*, i64 }* %taddr78
  call void @gominic_memcpy({ i8*, i64 }* %tyStr.addr64, { i8*, i64 }* %taddr78, i64 16)
  br label %endif14
endif14:
  %t79 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t80 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t79, i32 0, i32 2
  %t81 = load i64, i64* %t80, align 8
  %t82 = icmp eq i64 %t81, 2
  %t83 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t84 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t83, i32 0, i32 2
  %t85 = load i64, i64* %t84, align 8
  %t86 = icmp eq i64 %t85, 1
  br i1 %t82, label %logic.true16, label %logic.rhs17
logic.end15:
  %t92 = load i1, i1* %taddr87, align 1
  br i1 %t92, label %then18, label %else20
logic.true16:
  store i1 1, i1* %taddr87
  br label %logic.end15
logic.rhs17:
  %t88 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t89 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t88, i32 0, i32 2
  %t90 = load i64, i64* %t89, align 8
  %t91 = icmp eq i64 %t90, 1
  store i1 %t91, i1* %taddr87
  br label %logic.end15
then18:
  %t93 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t94 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t93, i32 0, i32 0
  %t95 = load { i8*, i64 }, { i8*, i64 }* %t94, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.125, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr97, align 8
  store { i8*, i64 } %t95, { i8*, i64 }* %taddr98, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr96, { i8*, i64 }* %taddr97, { i8*, i64 }* %taddr98)
  %t99 = load { i8*, i64 }, { i8*, i64 }* %taddr96, align 8
  store { i8*, i64 } %t99, { i8*, i64 }* %taddr100
  call void @gominic_memcpy({ i8*, i64 }* %valStr.addr66, { i8*, i64 }* %taddr100, i64 16)
  br label %endif19
endif19:
  br label %endif12
else20:
  %t101 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t102 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t101, i32 0, i32 2
  %t103 = load i64, i64* %t102, align 8
  %t104 = icmp eq i64 %t103, 3
  %t105 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t106 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t105, i32 0, i32 3
  %t107 = load { i8*, i64 }, { i8*, i64 }* %t106, align 8
  store { i8*, i64 } %t107, { i8*, i64 }* %taddr108, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.126, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr109, align 8
  %t110 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr108, i32 0, i32 0
  %t111 = load i8*, i8** %t110, align 8
  %t112 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr108, i32 0, i32 1
  %t113 = load i64, i64* %t112, align 8
  %t114 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr109, i32 0, i32 0
  %t115 = load i8*, i8** %t114, align 8
  %t116 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr109, i32 0, i32 1
  %t117 = load i64, i64* %t116, align 8
  %t118 = call i1 @gominic_str_eq(i8* %t111, i64 %t113, i8* %t115, i64 %t117)
  %t119 = icmp eq i1 %t118, 0
  br i1 %t104, label %logic.rhs22, label %logic.false23
logic.end21:
  %t136 = load i1, i1* %taddr120, align 1
  br i1 %t136, label %then24, label %else26
logic.rhs22:
  %t121 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t122 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t121, i32 0, i32 3
  %t123 = load { i8*, i64 }, { i8*, i64 }* %t122, align 8
  store { i8*, i64 } %t123, { i8*, i64 }* %taddr124, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.127, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr125, align 8
  %t126 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr124, i32 0, i32 0
  %t127 = load i8*, i8** %t126, align 8
  %t128 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr124, i32 0, i32 1
  %t129 = load i64, i64* %t128, align 8
  %t130 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr125, i32 0, i32 0
  %t131 = load i8*, i8** %t130, align 8
  %t132 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr125, i32 0, i32 1
  %t133 = load i64, i64* %t132, align 8
  %t134 = call i1 @gominic_str_eq(i8* %t127, i64 %t129, i8* %t131, i64 %t133)
  %t135 = icmp eq i1 %t134, 0
  store i1 %t135, i1* %taddr120
  br label %logic.end21
logic.false23:
  store i1 0, i1* %taddr120
  br label %logic.end21
then24:
  %t137 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t138 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t137, i32 0, i32 3
  %t139 = load { i8*, i64 }, { i8*, i64 }* %t138, align 8
  store { i8*, i64 } %t139, { i8*, i64 }* %taddr140
  call void @gominic_memcpy({ i8*, i64 }* %valStr.addr66, { i8*, i64 }* %taddr140, i64 16)
  br label %endif25
endif25:
  br label %endif19
else26:
  %t141 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t142 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t141, i32 0, i32 0
  %t143 = load { i8*, i64 }, { i8*, i64 }* %t142, align 8
  store { i8*, i64 } %t143, { i8*, i64 }* %taddr144
  call void @gominic_memcpy({ i8*, i64 }* %valStr.addr66, { i8*, i64 }* %taddr144, i64 16)
  br label %endif25
idx.ok27:
  %t166 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %args.addr13, i32 0, i32 0
  %t167 = load { i8*, i64 }*, { i8*, i64 }** %t166, align 8
  %t168 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %t167, i64 %t160
  store { i8*, i64 } %t159, { i8*, i64 }* %taddr169
  call void @gominic_memcpy({ i8*, i64 }* %t168, { i8*, i64 }* %taddr169, i64 16)
  br label %for.post2
idx.fail28:
  call void @gominic_abort()
  br label %idx.ok27
then29:
  %t178 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t179 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t178, i32 0, i32 3
  %t180 = load i8*, i8** %t179, align 8
  %t181 = call { i8*, i64 } @ir.Type.String(i8* %t180)
  store { i8*, i64 } %t181, { i8*, i64 }* %taddr182
  call void @gominic_memcpy({ i8*, i64 }* %ret.addr172, { i8*, i64 }* %taddr182, i64 16)
  br label %endif30
endif30:
  %prefix.addr183 = alloca { i8*, i64 }
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.130, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr184
  call void @gominic_memcpy({ i8*, i64 }* %prefix.addr183, { i8*, i64 }* %taddr184, i64 16)
  %t185 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t186 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t185, i32 0, i32 0
  %t187 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t186, align 8
  %t188 = icmp ne { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t187, null
  %t189 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t190 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t189, i32 0, i32 0
  %t191 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t190 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t192 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t191, i32 0, i32 0
  %t193 = load { i8*, i64 }, { i8*, i64 }* %t192, align 8
  store { i8*, i64 } %t193, { i8*, i64 }* %taddr194, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.131, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr195, align 8
  %t196 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr194, i32 0, i32 0
  %t197 = load i8*, i8** %t196, align 8
  %t198 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr194, i32 0, i32 1
  %t199 = load i64, i64* %t198, align 8
  %t200 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr195, i32 0, i32 0
  %t201 = load i8*, i8** %t200, align 8
  %t202 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr195, i32 0, i32 1
  %t203 = load i64, i64* %t202, align 8
  %t204 = call i1 @gominic_str_eq(i8* %t197, i64 %t199, i8* %t201, i64 %t203)
  %t205 = icmp eq i1 %t204, 0
  br i1 %t188, label %logic.rhs32, label %logic.false33
logic.end31:
  %t224 = load i1, i1* %taddr206, align 1
  %t225 = load { i8*, i64 }, { i8*, i64 }* %ret.addr172, align 8
  store { i8*, i64 } %t225, { i8*, i64 }* %taddr226, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.133, i32 0, i32 0), i64 4 }, { i8*, i64 }* %taddr227, align 8
  %t228 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr226, i32 0, i32 0
  %t229 = load i8*, i8** %t228, align 8
  %t230 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr226, i32 0, i32 1
  %t231 = load i64, i64* %t230, align 8
  %t232 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr227, i32 0, i32 0
  %t233 = load i8*, i8** %t232, align 8
  %t234 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr227, i32 0, i32 1
  %t235 = load i64, i64* %t234, align 8
  %t236 = call i1 @gominic_str_eq(i8* %t229, i64 %t231, i8* %t233, i64 %t235)
  %t237 = icmp eq i1 %t236, 0
  br i1 %t224, label %logic.rhs35, label %logic.false36
logic.rhs32:
  %t207 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t208 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t207, i32 0, i32 0
  %t209 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t208 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t210 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t209, i32 0, i32 0
  %t211 = load { i8*, i64 }, { i8*, i64 }* %t210, align 8
  store { i8*, i64 } %t211, { i8*, i64 }* %taddr212, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.132, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr213, align 8
  %t214 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr212, i32 0, i32 0
  %t215 = load i8*, i8** %t214, align 8
  %t216 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr212, i32 0, i32 1
  %t217 = load i64, i64* %t216, align 8
  %t218 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr213, i32 0, i32 0
  %t219 = load i8*, i8** %t218, align 8
  %t220 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr213, i32 0, i32 1
  %t221 = load i64, i64* %t220, align 8
  %t222 = call i1 @gominic_str_eq(i8* %t215, i64 %t217, i8* %t219, i64 %t221)
  %t223 = icmp eq i1 %t222, 0
  store i1 %t223, i1* %taddr206
  br label %logic.end31
logic.false33:
  store i1 0, i1* %taddr206
  br label %logic.end31
logic.end34:
  %t252 = load i1, i1* %taddr238, align 1
  br i1 %t252, label %then37, label %endif38
logic.rhs35:
  %t239 = load { i8*, i64 }, { i8*, i64 }* %ret.addr172, align 8
  store { i8*, i64 } %t239, { i8*, i64 }* %taddr240, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.134, i32 0, i32 0), i64 4 }, { i8*, i64 }* %taddr241, align 8
  %t242 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr240, i32 0, i32 0
  %t243 = load i8*, i8** %t242, align 8
  %t244 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr240, i32 0, i32 1
  %t245 = load i64, i64* %t244, align 8
  %t246 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr241, i32 0, i32 0
  %t247 = load i8*, i8** %t246, align 8
  %t248 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr241, i32 0, i32 1
  %t249 = load i64, i64* %t248, align 8
  %t250 = call i1 @gominic_str_eq(i8* %t243, i64 %t245, i8* %t247, i64 %t249)
  %t251 = icmp eq i1 %t250, 0
  store i1 %t251, i1* %taddr238
  br label %logic.end34
logic.false36:
  store i1 0, i1* %taddr238
  br label %logic.end34
then37:
  %t253 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t254 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t253, i32 0, i32 0
  %t255 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t254 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t256 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t255, i32 0, i32 0
  %t257 = load { i8*, i64 }, { i8*, i64 }* %t256, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.135, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr259, align 8
  store { i8*, i64 } %t257, { i8*, i64 }* %taddr260, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr258, { i8*, i64 }* %taddr259, { i8*, i64 }* %taddr260)
  %t261 = load { i8*, i64 }, { i8*, i64 }* %taddr258, align 8
  store { i8*, i64 } %t261, { i8*, i64 }* %taddr263, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.136, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr264, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr262, { i8*, i64 }* %taddr263, { i8*, i64 }* %taddr264)
  %t265 = load { i8*, i64 }, { i8*, i64 }* %taddr262, align 8
  store { i8*, i64 } %t265, { i8*, i64 }* %taddr266
  call void @gominic_memcpy({ i8*, i64 }* %prefix.addr183, { i8*, i64 }* %taddr266, i64 16)
  br label %endif38
endif38:
  %t267 = load { i8*, i64 }, { i8*, i64 }* %prefix.addr183, align 8
  store { i8*, i64 } %t267, { i8*, i64 }* %taddr269, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([6 x i8], [6 x i8]* @.str.137, i32 0, i32 0), i64 5 }, { i8*, i64 }* %taddr270, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr268, { i8*, i64 }* %taddr269, { i8*, i64 }* %taddr270)
  %t271 = load { i8*, i64 }, { i8*, i64 }* %taddr268, align 8
  %t272 = load { i8*, i64 }, { i8*, i64 }* %ret.addr172, align 8
  store { i8*, i64 } %t271, { i8*, i64 }* %taddr274, align 8
  store { i8*, i64 } %t272, { i8*, i64 }* %taddr275, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr273, { i8*, i64 }* %taddr274, { i8*, i64 }* %taddr275)
  %t276 = load { i8*, i64 }, { i8*, i64 }* %taddr273, align 8
  store { i8*, i64 } %t276, { i8*, i64 }* %taddr278, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.138, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr279, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr277, { i8*, i64 }* %taddr278, { i8*, i64 }* %taddr279)
  %t280 = load { i8*, i64 }, { i8*, i64 }* %taddr277, align 8
  %t281 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }*
  %t282 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8*, i8* }* %t281, i32 0, i32 1
  %t283 = load { i8*, i64 }, { i8*, i64 }* %t282, align 8
  store { i8*, i64 } %t280, { i8*, i64 }* %taddr285, align 8
  store { i8*, i64 } %t283, { i8*, i64 }* %taddr286, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr284, { i8*, i64 }* %taddr285, { i8*, i64 }* %taddr286)
  %t287 = load { i8*, i64 }, { i8*, i64 }* %taddr284, align 8
  store { i8*, i64 } %t287, { i8*, i64 }* %taddr289, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.139, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr290, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr288, { i8*, i64 }* %taddr289, { i8*, i64 }* %taddr290)
  %t291 = load { i8*, i64 }, { i8*, i64 }* %taddr288, align 8
  %t292 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %args.addr13, align 8
  store { { i8*, i64 }*, i64, i64 } %t292, { { i8*, i64 }*, i64, i64 }* %taddr293
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.140, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr294
  %t295 = call { i8*, i64 } @strings.Join({ { i8*, i64 }*, i64, i64 }* %taddr293, { i8*, i64 }* %taddr294)
  store { i8*, i64 } %t291, { i8*, i64 }* %taddr297, align 8
  store { i8*, i64 } %t295, { i8*, i64 }* %taddr298, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr296, { i8*, i64 }* %taddr297, { i8*, i64 }* %taddr298)
  %t299 = load { i8*, i64 }, { i8*, i64 }* %taddr296, align 8
  store { i8*, i64 } %t299, { i8*, i64 }* %taddr301, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.141, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr302, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr300, { i8*, i64 }* %taddr301, { i8*, i64 }* %taddr302)
  %t303 = load { i8*, i64 }, { i8*, i64 }* %taddr300, align 8
  ret { i8*, i64 } %t303
}

define void @ir.Conv.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.Conv.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %recv.addr
  %t0 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }*
  %t1 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %t0, i32 0, i32 0
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t1 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 0
  %t4 = load { i8*, i64 }, { i8*, i64 }* %t3, align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.142, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr6, align 8
  %taddr7 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t4, { i8*, i64 }* %taddr7, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6, { i8*, i64 }* %taddr7)
  %t8 = load { i8*, i64 }, { i8*, i64 }* %taddr5, align 8
  %taddr9 = alloca { i8*, i64 } , align 8
  %taddr10 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t8, { i8*, i64 }* %taddr10, align 8
  %taddr11 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.143, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr11, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr9, { i8*, i64 }* %taddr10, { i8*, i64 }* %taddr11)
  %t12 = load { i8*, i64 }, { i8*, i64 }* %taddr9, align 8
  %t13 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }*
  %t14 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %t13, i32 0, i32 1
  %t15 = load { i8*, i64 }, { i8*, i64 }* %t14, align 8
  %taddr16 = alloca { i8*, i64 } , align 8
  %taddr17 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t12, { i8*, i64 }* %taddr17, align 8
  %taddr18 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t15, { i8*, i64 }* %taddr18, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr16, { i8*, i64 }* %taddr17, { i8*, i64 }* %taddr18)
  %t19 = load { i8*, i64 }, { i8*, i64 }* %taddr16, align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t19, { i8*, i64 }* %taddr21, align 8
  %taddr22 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.144, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr22, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21, { i8*, i64 }* %taddr22)
  %t23 = load { i8*, i64 }, { i8*, i64 }* %taddr20, align 8
  %t24 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }*
  %t25 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %t24, i32 0, i32 2
  %t26 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t25, align 8
  %t27 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t26)
  %taddr28 = alloca { i8*, i64 } , align 8
  %taddr29 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t23, { i8*, i64 }* %taddr29, align 8
  %taddr30 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t27, { i8*, i64 }* %taddr30, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr28, { i8*, i64 }* %taddr29, { i8*, i64 }* %taddr30)
  %t31 = load { i8*, i64 }, { i8*, i64 }* %taddr28, align 8
  %taddr32 = alloca { i8*, i64 } , align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t31, { i8*, i64 }* %taddr33, align 8
  %taddr34 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.145, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr34, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr32, { i8*, i64 }* %taddr33, { i8*, i64 }* %taddr34)
  %t35 = load { i8*, i64 }, { i8*, i64 }* %taddr32, align 8
  %t36 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }*
  %t37 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %t36, i32 0, i32 2
  %t38 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t37, align 8
  %t39 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t38)
  %taddr40 = alloca { i8*, i64 } , align 8
  %taddr41 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t35, { i8*, i64 }* %taddr41, align 8
  %taddr42 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t39, { i8*, i64 }* %taddr42, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr40, { i8*, i64 }* %taddr41, { i8*, i64 }* %taddr42)
  %t43 = load { i8*, i64 }, { i8*, i64 }* %taddr40, align 8
  %taddr44 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t43, { i8*, i64 }* %taddr45, align 8
  %taddr46 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.146, i32 0, i32 0), i64 4 }, { i8*, i64 }* %taddr46, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr44, { i8*, i64 }* %taddr45, { i8*, i64 }* %taddr46)
  %t47 = load { i8*, i64 }, { i8*, i64 }* %taddr44, align 8
  %t48 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }*
  %t49 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %t48, i32 0, i32 3
  %t50 = load i8*, i8** %t49, align 8
  %t51 = call { i8*, i64 } @ir.Type.String(i8* %t50)
  %taddr52 = alloca { i8*, i64 } , align 8
  %taddr53 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t47, { i8*, i64 }* %taddr53, align 8
  %taddr54 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t51, { i8*, i64 }* %taddr54, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr52, { i8*, i64 }* %taddr53, { i8*, i64 }* %taddr54)
  %t55 = load { i8*, i64 }, { i8*, i64 }* %taddr52, align 8
  ret { i8*, i64 } %t55
}

define void @ir.Alloca.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.Alloca.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }* %recv.addr
  %align.addr0 = alloca { i8*, i64 }
  %taddr1 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.147, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr1
  call void @gominic_memcpy({ i8*, i64 }* %align.addr0, { i8*, i64 }* %taddr1, i64 16)
  %t2 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }* @ir.a to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }*
  %t3 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }* %t2, i32 0, i32 2
  %t4 = load i64, i64* %t3, align 8
  %t5 = icmp sgt i64 %t4, 0
  %taddr10 = alloca { i8*, i64 } , align 8
  %taddr11 = alloca { i8*, i64 } , align 8
  %taddr12 = alloca { i8*, i64 } , align 8
  %taddr14 = alloca { i8*, i64 } , align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  %taddr22 = alloca { i8*, i64 } , align 8
  %taddr24 = alloca { i8*, i64 } , align 8
  %taddr25 = alloca { i8*, i64 } , align 8
  %taddr26 = alloca { i8*, i64 } , align 8
  %taddr32 = alloca { i8*, i64 } , align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  %taddr34 = alloca { i8*, i64 } , align 8
  %taddr37 = alloca { i8*, i64 } , align 8
  %taddr38 = alloca { i8*, i64 } , align 8
  %taddr39 = alloca { i8*, i64 } , align 8
  br i1 %t5, label %then0, label %endif1
then0:
  %t6 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }* @ir.a to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }*
  %t7 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }* %t6, i32 0, i32 2
  %t8 = load i64, i64* %t7, align 8
  %t9 = call { i8*, i64 } @strconv.FormatInt(i64 %t8, i64 10)
  store { i8*, i64 } { i8* getelementptr inbounds ([10 x i8], [10 x i8]* @.str.148, i32 0, i32 0), i64 9 }, { i8*, i64 }* %taddr11, align 8
  store { i8*, i64 } %t9, { i8*, i64 }* %taddr12, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr10, { i8*, i64 }* %taddr11, { i8*, i64 }* %taddr12)
  %t13 = load { i8*, i64 }, { i8*, i64 }* %taddr10, align 8
  store { i8*, i64 } %t13, { i8*, i64 }* %taddr14
  call void @gominic_memcpy({ i8*, i64 }* %align.addr0, { i8*, i64 }* %taddr14, i64 16)
  br label %endif1
endif1:
  %t15 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }* @ir.a to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }*
  %t16 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }* %t15, i32 0, i32 0
  %t17 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t16 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t18 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t17, i32 0, i32 0
  %t19 = load { i8*, i64 }, { i8*, i64 }* %t18, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.149, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr21, align 8
  store { i8*, i64 } %t19, { i8*, i64 }* %taddr22, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21, { i8*, i64 }* %taddr22)
  %t23 = load { i8*, i64 }, { i8*, i64 }* %taddr20, align 8
  store { i8*, i64 } %t23, { i8*, i64 }* %taddr25, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([11 x i8], [11 x i8]* @.str.150, i32 0, i32 0), i64 10 }, { i8*, i64 }* %taddr26, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr24, { i8*, i64 }* %taddr25, { i8*, i64 }* %taddr26)
  %t27 = load { i8*, i64 }, { i8*, i64 }* %taddr24, align 8
  %t28 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }* @ir.a to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }*
  %t29 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, i64 }* %t28, i32 0, i32 1
  %t30 = load i8*, i8** %t29, align 8
  %t31 = call { i8*, i64 } @ir.Type.String(i8* %t30)
  store { i8*, i64 } %t27, { i8*, i64 }* %taddr33, align 8
  store { i8*, i64 } %t31, { i8*, i64 }* %taddr34, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr32, { i8*, i64 }* %taddr33, { i8*, i64 }* %taddr34)
  %t35 = load { i8*, i64 }, { i8*, i64 }* %taddr32, align 8
  %t36 = load { i8*, i64 }, { i8*, i64 }* %align.addr0, align 8
  store { i8*, i64 } %t35, { i8*, i64 }* %taddr38, align 8
  store { i8*, i64 } %t36, { i8*, i64 }* %taddr39, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr37, { i8*, i64 }* %taddr38, { i8*, i64 }* %taddr39)
  %t40 = load { i8*, i64 }, { i8*, i64 }* %taddr37, align 8
  ret { i8*, i64 } %t40
}

define void @ir.Load.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.Load.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %recv.addr
  %align.addr0 = alloca { i8*, i64 }
  %taddr1 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.151, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr1
  call void @gominic_memcpy({ i8*, i64 }* %align.addr0, { i8*, i64 }* %taddr1, i64 16)
  %t2 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.l to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t3 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t2, i32 0, i32 2
  %t4 = load i64, i64* %t3, align 8
  %t5 = icmp sgt i64 %t4, 0
  %taddr10 = alloca { i8*, i64 } , align 8
  %taddr11 = alloca { i8*, i64 } , align 8
  %taddr12 = alloca { i8*, i64 } , align 8
  %taddr14 = alloca { i8*, i64 } , align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  %taddr22 = alloca { i8*, i64 } , align 8
  %taddr24 = alloca { i8*, i64 } , align 8
  %taddr25 = alloca { i8*, i64 } , align 8
  %taddr26 = alloca { i8*, i64 } , align 8
  %taddr32 = alloca { i8*, i64 } , align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  %taddr34 = alloca { i8*, i64 } , align 8
  %taddr36 = alloca { i8*, i64 } , align 8
  %taddr37 = alloca { i8*, i64 } , align 8
  %taddr38 = alloca { i8*, i64 } , align 8
  %taddr44 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca { i8*, i64 } , align 8
  %taddr46 = alloca { i8*, i64 } , align 8
  %taddr48 = alloca { i8*, i64 } , align 8
  %taddr49 = alloca { i8*, i64 } , align 8
  %taddr50 = alloca { i8*, i64 } , align 8
  %taddr56 = alloca { i8*, i64 } , align 8
  %taddr57 = alloca { i8*, i64 } , align 8
  %taddr58 = alloca { i8*, i64 } , align 8
  %taddr61 = alloca { i8*, i64 } , align 8
  %taddr62 = alloca { i8*, i64 } , align 8
  %taddr63 = alloca { i8*, i64 } , align 8
  br i1 %t5, label %then0, label %endif1
then0:
  %t6 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.l to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t7 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t6, i32 0, i32 2
  %t8 = load i64, i64* %t7, align 8
  %t9 = call { i8*, i64 } @strconv.FormatInt(i64 %t8, i64 10)
  store { i8*, i64 } { i8* getelementptr inbounds ([9 x i8], [9 x i8]* @.str.152, i32 0, i32 0), i64 8 }, { i8*, i64 }* %taddr11, align 8
  store { i8*, i64 } %t9, { i8*, i64 }* %taddr12, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr10, { i8*, i64 }* %taddr11, { i8*, i64 }* %taddr12)
  %t13 = load { i8*, i64 }, { i8*, i64 }* %taddr10, align 8
  store { i8*, i64 } %t13, { i8*, i64 }* %taddr14
  call void @gominic_memcpy({ i8*, i64 }* %align.addr0, { i8*, i64 }* %taddr14, i64 16)
  br label %endif1
endif1:
  %t15 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.l to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t16 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t15, i32 0, i32 0
  %t17 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t16 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t18 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t17, i32 0, i32 0
  %t19 = load { i8*, i64 }, { i8*, i64 }* %t18, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.153, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr21, align 8
  store { i8*, i64 } %t19, { i8*, i64 }* %taddr22, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21, { i8*, i64 }* %taddr22)
  %t23 = load { i8*, i64 }, { i8*, i64 }* %taddr20, align 8
  store { i8*, i64 } %t23, { i8*, i64 }* %taddr25, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([9 x i8], [9 x i8]* @.str.154, i32 0, i32 0), i64 8 }, { i8*, i64 }* %taddr26, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr24, { i8*, i64 }* %taddr25, { i8*, i64 }* %taddr26)
  %t27 = load { i8*, i64 }, { i8*, i64 }* %taddr24, align 8
  %t28 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.l to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t29 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t28, i32 0, i32 0
  %t30 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t29, align 8
  %t31 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t30)
  store { i8*, i64 } %t27, { i8*, i64 }* %taddr33, align 8
  store { i8*, i64 } %t31, { i8*, i64 }* %taddr34, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr32, { i8*, i64 }* %taddr33, { i8*, i64 }* %taddr34)
  %t35 = load { i8*, i64 }, { i8*, i64 }* %taddr32, align 8
  store { i8*, i64 } %t35, { i8*, i64 }* %taddr37, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.155, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr38, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr36, { i8*, i64 }* %taddr37, { i8*, i64 }* %taddr38)
  %t39 = load { i8*, i64 }, { i8*, i64 }* %taddr36, align 8
  %t40 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.l to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t41 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t40, i32 0, i32 1
  %t42 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t41, align 8
  %t43 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t42)
  store { i8*, i64 } %t39, { i8*, i64 }* %taddr45, align 8
  store { i8*, i64 } %t43, { i8*, i64 }* %taddr46, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr44, { i8*, i64 }* %taddr45, { i8*, i64 }* %taddr46)
  %t47 = load { i8*, i64 }, { i8*, i64 }* %taddr44, align 8
  store { i8*, i64 } %t47, { i8*, i64 }* %taddr49, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.156, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr50, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr48, { i8*, i64 }* %taddr49, { i8*, i64 }* %taddr50)
  %t51 = load { i8*, i64 }, { i8*, i64 }* %taddr48, align 8
  %t52 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.l to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t53 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t52, i32 0, i32 1
  %t54 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t53, align 8
  %t55 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t54)
  store { i8*, i64 } %t51, { i8*, i64 }* %taddr57, align 8
  store { i8*, i64 } %t55, { i8*, i64 }* %taddr58, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr56, { i8*, i64 }* %taddr57, { i8*, i64 }* %taddr58)
  %t59 = load { i8*, i64 }, { i8*, i64 }* %taddr56, align 8
  %t60 = load { i8*, i64 }, { i8*, i64 }* %align.addr0, align 8
  store { i8*, i64 } %t59, { i8*, i64 }* %taddr62, align 8
  store { i8*, i64 } %t60, { i8*, i64 }* %taddr63, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr61, { i8*, i64 }* %taddr62, { i8*, i64 }* %taddr63)
  %t64 = load { i8*, i64 }, { i8*, i64 }* %taddr61, align 8
  ret { i8*, i64 } %t64
}

define void @ir.Store.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.Store.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %recv.addr
  %align.addr0 = alloca { i8*, i64 }
  %taddr1 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.157, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr1
  call void @gominic_memcpy({ i8*, i64 }* %align.addr0, { i8*, i64 }* %taddr1, i64 16)
  %t2 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.s to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t3 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t2, i32 0, i32 2
  %t4 = load i64, i64* %t3, align 8
  %t5 = icmp sgt i64 %t4, 0
  %taddr10 = alloca { i8*, i64 } , align 8
  %taddr11 = alloca { i8*, i64 } , align 8
  %taddr12 = alloca { i8*, i64 } , align 8
  %taddr14 = alloca { i8*, i64 } , align 8
  %taddr19 = alloca { i8*, i64 } , align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  %taddr23 = alloca { i8*, i64 } , align 8
  %taddr24 = alloca { i8*, i64 } , align 8
  %taddr25 = alloca { i8*, i64 } , align 8
  %taddr31 = alloca { i8*, i64 } , align 8
  %taddr32 = alloca { i8*, i64 } , align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  %taddr35 = alloca { i8*, i64 } , align 8
  %taddr36 = alloca { i8*, i64 } , align 8
  %taddr37 = alloca { i8*, i64 } , align 8
  %taddr43 = alloca { i8*, i64 } , align 8
  %taddr44 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca { i8*, i64 } , align 8
  %taddr47 = alloca { i8*, i64 } , align 8
  %taddr48 = alloca { i8*, i64 } , align 8
  %taddr49 = alloca { i8*, i64 } , align 8
  %taddr55 = alloca { i8*, i64 } , align 8
  %taddr56 = alloca { i8*, i64 } , align 8
  %taddr57 = alloca { i8*, i64 } , align 8
  %taddr60 = alloca { i8*, i64 } , align 8
  %taddr61 = alloca { i8*, i64 } , align 8
  %taddr62 = alloca { i8*, i64 } , align 8
  br i1 %t5, label %then0, label %endif1
then0:
  %t6 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.s to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t7 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t6, i32 0, i32 2
  %t8 = load i64, i64* %t7, align 8
  %t9 = call { i8*, i64 } @strconv.FormatInt(i64 %t8, i64 10)
  store { i8*, i64 } { i8* getelementptr inbounds ([9 x i8], [9 x i8]* @.str.158, i32 0, i32 0), i64 8 }, { i8*, i64 }* %taddr11, align 8
  store { i8*, i64 } %t9, { i8*, i64 }* %taddr12, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr10, { i8*, i64 }* %taddr11, { i8*, i64 }* %taddr12)
  %t13 = load { i8*, i64 }, { i8*, i64 }* %taddr10, align 8
  store { i8*, i64 } %t13, { i8*, i64 }* %taddr14
  call void @gominic_memcpy({ i8*, i64 }* %align.addr0, { i8*, i64 }* %taddr14, i64 16)
  br label %endif1
endif1:
  %t15 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.s to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t16 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t15, i32 0, i32 0
  %t17 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t16, align 8
  %t18 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t17)
  store { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.159, i32 0, i32 0), i64 6 }, { i8*, i64 }* %taddr20, align 8
  store { i8*, i64 } %t18, { i8*, i64 }* %taddr21, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr19, { i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21)
  %t22 = load { i8*, i64 }, { i8*, i64 }* %taddr19, align 8
  store { i8*, i64 } %t22, { i8*, i64 }* %taddr24, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.160, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr25, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr23, { i8*, i64 }* %taddr24, { i8*, i64 }* %taddr25)
  %t26 = load { i8*, i64 }, { i8*, i64 }* %taddr23, align 8
  %t27 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.s to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t28 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t27, i32 0, i32 0
  %t29 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t28, align 8
  %t30 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t29)
  store { i8*, i64 } %t26, { i8*, i64 }* %taddr32, align 8
  store { i8*, i64 } %t30, { i8*, i64 }* %taddr33, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr31, { i8*, i64 }* %taddr32, { i8*, i64 }* %taddr33)
  %t34 = load { i8*, i64 }, { i8*, i64 }* %taddr31, align 8
  store { i8*, i64 } %t34, { i8*, i64 }* %taddr36, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.161, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr37, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr35, { i8*, i64 }* %taddr36, { i8*, i64 }* %taddr37)
  %t38 = load { i8*, i64 }, { i8*, i64 }* %taddr35, align 8
  %t39 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.s to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t40 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t39, i32 0, i32 1
  %t41 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t40, align 8
  %t42 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t41)
  store { i8*, i64 } %t38, { i8*, i64 }* %taddr44, align 8
  store { i8*, i64 } %t42, { i8*, i64 }* %taddr45, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr43, { i8*, i64 }* %taddr44, { i8*, i64 }* %taddr45)
  %t46 = load { i8*, i64 }, { i8*, i64 }* %taddr43, align 8
  store { i8*, i64 } %t46, { i8*, i64 }* %taddr48, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.162, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr49, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr47, { i8*, i64 }* %taddr48, { i8*, i64 }* %taddr49)
  %t50 = load { i8*, i64 }, { i8*, i64 }* %taddr47, align 8
  %t51 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* @ir.s to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }*
  %t52 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i64 }* %t51, i32 0, i32 1
  %t53 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t52, align 8
  %t54 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t53)
  store { i8*, i64 } %t50, { i8*, i64 }* %taddr56, align 8
  store { i8*, i64 } %t54, { i8*, i64 }* %taddr57, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr55, { i8*, i64 }* %taddr56, { i8*, i64 }* %taddr57)
  %t58 = load { i8*, i64 }, { i8*, i64 }* %taddr55, align 8
  %t59 = load { i8*, i64 }, { i8*, i64 }* %align.addr0, align 8
  store { i8*, i64 } %t58, { i8*, i64 }* %taddr61, align 8
  store { i8*, i64 } %t59, { i8*, i64 }* %taddr62, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr60, { i8*, i64 }* %taddr61, { i8*, i64 }* %taddr62)
  %t63 = load { i8*, i64 }, { i8*, i64 }* %taddr60, align 8
  ret { i8*, i64 } %t63
}

define void @ir.GetElementPtr.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.GetElementPtr.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* %recv.addr
  %t0 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* @ir.g to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }*
  %t1 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* %t0, i32 0, i32 3
  %t2 = load { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t1, align 8
  %taddr3 = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } %t2, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr3, align 8
  %t4 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr3, i32 0, i32 1
  %t5 = load i64, i64* %t4, align 8
  %t6 = call i8* @gominic_makeSlice(i64 %t5, i64 %t5, i64 16)
  %t7 = bitcast i8* %t6 to { i8*, i64 }*
  %taddr8 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %t9 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 0
  store { i8*, i64 }* %t7, { i8*, i64 }** %t9, align 8
  %t10 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 1
  store i64 %t5, i64* %t10, align 8
  %t11 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 2
  store i64 %t5, i64* %t11, align 8
  %t12 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, align 8
  %parts.addr13 = alloca { { i8*, i64 }*, i64, i64 }
  %taddr14 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  store { { i8*, i64 }*, i64, i64 } %t12, { { i8*, i64 }*, i64, i64 }* %taddr14
  call void @gominic_memcpy({ { i8*, i64 }*, i64, i64 }* %parts.addr13, { { i8*, i64 }*, i64, i64 }* %taddr14, i64 24)
  %i.addr15 = alloca i64
  store i64 0, i64* %i.addr15, align 8
  %taddr20 = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } , align 8
  %taddr39 = alloca { i8*, i64 } , align 8
  %taddr40 = alloca { i8*, i64 } , align 8
  %taddr41 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca { i8*, i64 } , align 8
  %taddr46 = alloca { i8*, i64 } , align 8
  %taddr47 = alloca { i8*, i64 } , align 8
  %taddr58 = alloca { i8*, i64 } , align 8
  %taddr66 = alloca { i8*, i64 } , align 8
  %taddr67 = alloca { i8*, i64 } , align 8
  %taddr68 = alloca { i8*, i64 } , align 8
  %taddr70 = alloca { i8*, i64 } , align 8
  %taddr71 = alloca { i8*, i64 } , align 8
  %taddr72 = alloca { i8*, i64 } , align 8
  %taddr78 = alloca { i8*, i64 } , align 8
  %taddr79 = alloca { i8*, i64 } , align 8
  %taddr80 = alloca { i8*, i64 } , align 8
  %taddr82 = alloca { i8*, i64 } , align 8
  %taddr83 = alloca { i8*, i64 } , align 8
  %taddr84 = alloca { i8*, i64 } , align 8
  %taddr90 = alloca { i8*, i64 } , align 8
  %taddr91 = alloca { i8*, i64 } , align 8
  %taddr92 = alloca { i8*, i64 } , align 8
  %taddr94 = alloca { i8*, i64 } , align 8
  %taddr95 = alloca { i8*, i64 } , align 8
  %taddr96 = alloca { i8*, i64 } , align 8
  %taddr102 = alloca { i8*, i64 } , align 8
  %taddr103 = alloca { i8*, i64 } , align 8
  %taddr104 = alloca { i8*, i64 } , align 8
  %taddr106 = alloca { i8*, i64 } , align 8
  %taddr107 = alloca { i8*, i64 } , align 8
  %taddr108 = alloca { i8*, i64 } , align 8
  %taddr111 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %taddr112 = alloca { i8*, i64 } , align 8
  %taddr114 = alloca { i8*, i64 } , align 8
  %taddr115 = alloca { i8*, i64 } , align 8
  %taddr116 = alloca { i8*, i64 } , align 8
  br label %for.cond0
for.cond0:
  %t16 = load i64, i64* %i.addr15, align 8
  %t17 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* @ir.g to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }*
  %t18 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* %t17, i32 0, i32 3
  %t19 = load { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t18, align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } %t19, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr20, align 8
  %t21 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr20, i32 0, i32 1
  %t22 = load i64, i64* %t21, align 8
  %t23 = icmp slt i64 %t16, %t22
  br i1 %t23, label %for.body1, label %for.end3
for.body1:
  %t24 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* @ir.g to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }*
  %t25 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* %t24, i32 0, i32 3
  %t26 = load i64, i64* %i.addr15, align 8
  %t27 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t25, i32 0, i32 1
  %t28 = load i64, i64* %t27, align 8
  %t29 = icmp sge i64 %t26, 0
  %t30 = icmp slt i64 %t26, %t28
  %t31 = and i1 %t29, %t30
  br i1 %t31, label %idx.ok4, label %idx.fail5
for.post2:
  %t59 = load i64, i64* %i.addr15, align 8
  %t60 = add i64 %t59, 1
  store i64 %t60, i64* %i.addr15, align 8
  br label %for.cond0
for.end3:
  %t61 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* @ir.g to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }*
  %t62 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* %t61, i32 0, i32 0
  %t63 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t62 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t64 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t63, i32 0, i32 0
  %t65 = load { i8*, i64 }, { i8*, i64 }* %t64, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.164, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr67, align 8
  store { i8*, i64 } %t65, { i8*, i64 }* %taddr68, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr66, { i8*, i64 }* %taddr67, { i8*, i64 }* %taddr68)
  %t69 = load { i8*, i64 }, { i8*, i64 }* %taddr66, align 8
  store { i8*, i64 } %t69, { i8*, i64 }* %taddr71, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([27 x i8], [27 x i8]* @.str.165, i32 0, i32 0), i64 26 }, { i8*, i64 }* %taddr72, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr70, { i8*, i64 }* %taddr71, { i8*, i64 }* %taddr72)
  %t73 = load { i8*, i64 }, { i8*, i64 }* %taddr70, align 8
  %t74 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* @ir.g to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }*
  %t75 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* %t74, i32 0, i32 2
  %t76 = load i8*, i8** %t75, align 8
  %t77 = call { i8*, i64 } @ir.Type.String(i8* %t76)
  store { i8*, i64 } %t73, { i8*, i64 }* %taddr79, align 8
  store { i8*, i64 } %t77, { i8*, i64 }* %taddr80, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr78, { i8*, i64 }* %taddr79, { i8*, i64 }* %taddr80)
  %t81 = load { i8*, i64 }, { i8*, i64 }* %taddr78, align 8
  store { i8*, i64 } %t81, { i8*, i64 }* %taddr83, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.166, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr84, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr82, { i8*, i64 }* %taddr83, { i8*, i64 }* %taddr84)
  %t85 = load { i8*, i64 }, { i8*, i64 }* %taddr82, align 8
  %t86 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* @ir.g to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }*
  %t87 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* %t86, i32 0, i32 1
  %t88 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t87, align 8
  %t89 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t88)
  store { i8*, i64 } %t85, { i8*, i64 }* %taddr91, align 8
  store { i8*, i64 } %t89, { i8*, i64 }* %taddr92, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr90, { i8*, i64 }* %taddr91, { i8*, i64 }* %taddr92)
  %t93 = load { i8*, i64 }, { i8*, i64 }* %taddr90, align 8
  store { i8*, i64 } %t93, { i8*, i64 }* %taddr95, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.167, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr96, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr94, { i8*, i64 }* %taddr95, { i8*, i64 }* %taddr96)
  %t97 = load { i8*, i64 }, { i8*, i64 }* %taddr94, align 8
  %t98 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* @ir.g to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }*
  %t99 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8*, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, i8* }* %t98, i32 0, i32 1
  %t100 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t99, align 8
  %t101 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t100)
  store { i8*, i64 } %t97, { i8*, i64 }* %taddr103, align 8
  store { i8*, i64 } %t101, { i8*, i64 }* %taddr104, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr102, { i8*, i64 }* %taddr103, { i8*, i64 }* %taddr104)
  %t105 = load { i8*, i64 }, { i8*, i64 }* %taddr102, align 8
  store { i8*, i64 } %t105, { i8*, i64 }* %taddr107, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.168, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr108, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr106, { i8*, i64 }* %taddr107, { i8*, i64 }* %taddr108)
  %t109 = load { i8*, i64 }, { i8*, i64 }* %taddr106, align 8
  %t110 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %parts.addr13, align 8
  store { { i8*, i64 }*, i64, i64 } %t110, { { i8*, i64 }*, i64, i64 }* %taddr111
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.169, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr112
  %t113 = call { i8*, i64 } @strings.Join({ { i8*, i64 }*, i64, i64 }* %taddr111, { i8*, i64 }* %taddr112)
  store { i8*, i64 } %t109, { i8*, i64 }* %taddr115, align 8
  store { i8*, i64 } %t113, { i8*, i64 }* %taddr116, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr114, { i8*, i64 }* %taddr115, { i8*, i64 }* %taddr116)
  %t117 = load { i8*, i64 }, { i8*, i64 }* %taddr114, align 8
  ret { i8*, i64 } %t117
idx.ok4:
  %t32 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t25, i32 0, i32 0
  %t33 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*** %t32, align 8
  %t34 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t33, i64 %t26
  %t35 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t34, align 8
  %idx.addr36 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t35, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %idx.addr36, align 8
  %t37 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %idx.addr36, align 8
  %t38 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t37)
  store { i8*, i64 } %t38, { i8*, i64 }* %taddr40, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.163, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr41, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr39, { i8*, i64 }* %taddr40, { i8*, i64 }* %taddr41)
  %t42 = load { i8*, i64 }, { i8*, i64 }* %taddr39, align 8
  %t43 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %idx.addr36, align 8
  %t44 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t43)
  store { i8*, i64 } %t42, { i8*, i64 }* %taddr46, align 8
  store { i8*, i64 } %t44, { i8*, i64 }* %taddr47, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr45, { i8*, i64 }* %taddr46, { i8*, i64 }* %taddr47)
  %t48 = load { i8*, i64 }, { i8*, i64 }* %taddr45, align 8
  %t49 = load i64, i64* %i.addr15, align 8
  %t50 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %parts.addr13, i32 0, i32 1
  %t51 = load i64, i64* %t50, align 8
  %t52 = icmp sge i64 %t49, 0
  %t53 = icmp slt i64 %t49, %t51
  %t54 = and i1 %t52, %t53
  br i1 %t54, label %idx.ok6, label %idx.fail7
idx.fail5:
  call void @gominic_abort()
  br label %idx.ok4
idx.ok6:
  %t55 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %parts.addr13, i32 0, i32 0
  %t56 = load { i8*, i64 }*, { i8*, i64 }** %t55, align 8
  %t57 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %t56, i64 %t49
  store { i8*, i64 } %t48, { i8*, i64 }* %taddr58
  call void @gominic_memcpy({ i8*, i64 }* %t57, { i8*, i64 }* %taddr58, i64 16)
  br label %for.post2
idx.fail7:
  call void @gominic_abort()
  br label %idx.ok6
}

define void @ir.ICmp.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.ICmp.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %recv.addr
  %t0 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t1 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t0, i32 0, i32 0
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t1 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 0
  %t4 = load { i8*, i64 }, { i8*, i64 }* %t3, align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.170, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr6, align 8
  %taddr7 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t4, { i8*, i64 }* %taddr7, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6, { i8*, i64 }* %taddr7)
  %t8 = load { i8*, i64 }, { i8*, i64 }* %taddr5, align 8
  %taddr9 = alloca { i8*, i64 } , align 8
  %taddr10 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t8, { i8*, i64 }* %taddr10, align 8
  %taddr11 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([9 x i8], [9 x i8]* @.str.171, i32 0, i32 0), i64 8 }, { i8*, i64 }* %taddr11, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr9, { i8*, i64 }* %taddr10, { i8*, i64 }* %taddr11)
  %t12 = load { i8*, i64 }, { i8*, i64 }* %taddr9, align 8
  %t13 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t14 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t13, i32 0, i32 1
  %t15 = load { i8*, i64 }, { i8*, i64 }* %t14, align 8
  %taddr16 = alloca { i8*, i64 } , align 8
  %taddr17 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t12, { i8*, i64 }* %taddr17, align 8
  %taddr18 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t15, { i8*, i64 }* %taddr18, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr16, { i8*, i64 }* %taddr17, { i8*, i64 }* %taddr18)
  %t19 = load { i8*, i64 }, { i8*, i64 }* %taddr16, align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t19, { i8*, i64 }* %taddr21, align 8
  %taddr22 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.172, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr22, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21, { i8*, i64 }* %taddr22)
  %t23 = load { i8*, i64 }, { i8*, i64 }* %taddr20, align 8
  %t24 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t25 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t24, i32 0, i32 2
  %t26 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t25, align 8
  %t27 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t26)
  %taddr28 = alloca { i8*, i64 } , align 8
  %taddr29 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t23, { i8*, i64 }* %taddr29, align 8
  %taddr30 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t27, { i8*, i64 }* %taddr30, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr28, { i8*, i64 }* %taddr29, { i8*, i64 }* %taddr30)
  %t31 = load { i8*, i64 }, { i8*, i64 }* %taddr28, align 8
  %taddr32 = alloca { i8*, i64 } , align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t31, { i8*, i64 }* %taddr33, align 8
  %taddr34 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.173, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr34, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr32, { i8*, i64 }* %taddr33, { i8*, i64 }* %taddr34)
  %t35 = load { i8*, i64 }, { i8*, i64 }* %taddr32, align 8
  %t36 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t37 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t36, i32 0, i32 2
  %t38 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t37, align 8
  %t39 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t38)
  %taddr40 = alloca { i8*, i64 } , align 8
  %taddr41 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t35, { i8*, i64 }* %taddr41, align 8
  %taddr42 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t39, { i8*, i64 }* %taddr42, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr40, { i8*, i64 }* %taddr41, { i8*, i64 }* %taddr42)
  %t43 = load { i8*, i64 }, { i8*, i64 }* %taddr40, align 8
  %taddr44 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t43, { i8*, i64 }* %taddr45, align 8
  %taddr46 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.174, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr46, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr44, { i8*, i64 }* %taddr45, { i8*, i64 }* %taddr46)
  %t47 = load { i8*, i64 }, { i8*, i64 }* %taddr44, align 8
  %t48 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t49 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t48, i32 0, i32 3
  %t50 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t49, align 8
  %t51 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t50)
  %taddr52 = alloca { i8*, i64 } , align 8
  %taddr53 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t47, { i8*, i64 }* %taddr53, align 8
  %taddr54 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t51, { i8*, i64 }* %taddr54, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr52, { i8*, i64 }* %taddr53, { i8*, i64 }* %taddr54)
  %t55 = load { i8*, i64 }, { i8*, i64 }* %taddr52, align 8
  ret { i8*, i64 } %t55
}

define void @ir.FCmp.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.FCmp.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %recv.addr
  %t0 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t1 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t0, i32 0, i32 0
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t1 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 0
  %t4 = load { i8*, i64 }, { i8*, i64 }* %t3, align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.175, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr6, align 8
  %taddr7 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t4, { i8*, i64 }* %taddr7, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6, { i8*, i64 }* %taddr7)
  %t8 = load { i8*, i64 }, { i8*, i64 }* %taddr5, align 8
  %taddr9 = alloca { i8*, i64 } , align 8
  %taddr10 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t8, { i8*, i64 }* %taddr10, align 8
  %taddr11 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([9 x i8], [9 x i8]* @.str.176, i32 0, i32 0), i64 8 }, { i8*, i64 }* %taddr11, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr9, { i8*, i64 }* %taddr10, { i8*, i64 }* %taddr11)
  %t12 = load { i8*, i64 }, { i8*, i64 }* %taddr9, align 8
  %t13 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t14 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t13, i32 0, i32 1
  %t15 = load { i8*, i64 }, { i8*, i64 }* %t14, align 8
  %taddr16 = alloca { i8*, i64 } , align 8
  %taddr17 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t12, { i8*, i64 }* %taddr17, align 8
  %taddr18 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t15, { i8*, i64 }* %taddr18, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr16, { i8*, i64 }* %taddr17, { i8*, i64 }* %taddr18)
  %t19 = load { i8*, i64 }, { i8*, i64 }* %taddr16, align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t19, { i8*, i64 }* %taddr21, align 8
  %taddr22 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.177, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr22, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21, { i8*, i64 }* %taddr22)
  %t23 = load { i8*, i64 }, { i8*, i64 }* %taddr20, align 8
  %t24 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t25 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t24, i32 0, i32 2
  %t26 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t25, align 8
  %t27 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t26)
  %taddr28 = alloca { i8*, i64 } , align 8
  %taddr29 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t23, { i8*, i64 }* %taddr29, align 8
  %taddr30 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t27, { i8*, i64 }* %taddr30, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr28, { i8*, i64 }* %taddr29, { i8*, i64 }* %taddr30)
  %t31 = load { i8*, i64 }, { i8*, i64 }* %taddr28, align 8
  %taddr32 = alloca { i8*, i64 } , align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t31, { i8*, i64 }* %taddr33, align 8
  %taddr34 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.178, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr34, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr32, { i8*, i64 }* %taddr33, { i8*, i64 }* %taddr34)
  %t35 = load { i8*, i64 }, { i8*, i64 }* %taddr32, align 8
  %t36 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t37 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t36, i32 0, i32 2
  %t38 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t37, align 8
  %t39 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t38)
  %taddr40 = alloca { i8*, i64 } , align 8
  %taddr41 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t35, { i8*, i64 }* %taddr41, align 8
  %taddr42 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t39, { i8*, i64 }* %taddr42, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr40, { i8*, i64 }* %taddr41, { i8*, i64 }* %taddr42)
  %t43 = load { i8*, i64 }, { i8*, i64 }* %taddr40, align 8
  %taddr44 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t43, { i8*, i64 }* %taddr45, align 8
  %taddr46 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.179, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr46, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr44, { i8*, i64 }* %taddr45, { i8*, i64 }* %taddr46)
  %t47 = load { i8*, i64 }, { i8*, i64 }* %taddr44, align 8
  %t48 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.c to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t49 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t48, i32 0, i32 3
  %t50 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t49, align 8
  %t51 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t50)
  %taddr52 = alloca { i8*, i64 } , align 8
  %taddr53 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t47, { i8*, i64 }* %taddr53, align 8
  %taddr54 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t51, { i8*, i64 }* %taddr54, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr52, { i8*, i64 }* %taddr53, { i8*, i64 }* %taddr54)
  %t55 = load { i8*, i64 }, { i8*, i64 }* %taddr52, align 8
  ret { i8*, i64 } %t55
}

define void @ir.Br.isInstruction({ { i8*, i64 } } %recv) {
entry:
  %recv.addr = alloca { { i8*, i64 } } , align 8
  store { { i8*, i64 } } %recv, { { i8*, i64 } }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.Br.String({ { i8*, i64 } } %recv) {
entry:
  %recv.addr = alloca { { i8*, i64 } } , align 8
  store { { i8*, i64 } } %recv, { { i8*, i64 } }* %recv.addr
  %t0 = bitcast { { i8*, i64 } }* @ir.b to { { i8*, i64 } }*
  %t1 = getelementptr inbounds { { i8*, i64 } }, { { i8*, i64 } }* %t0, i32 0, i32 0
  %t2 = load { i8*, i64 }, { i8*, i64 }* %t1, align 8
  %taddr3 = alloca { i8*, i64 } , align 8
  %taddr4 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([11 x i8], [11 x i8]* @.str.180, i32 0, i32 0), i64 10 }, { i8*, i64 }* %taddr4, align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t2, { i8*, i64 }* %taddr5, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr3, { i8*, i64 }* %taddr4, { i8*, i64 }* %taddr5)
  %t6 = load { i8*, i64 }, { i8*, i64 }* %taddr3, align 8
  ret { i8*, i64 } %t6
}

define void @ir.CondBr.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.CondBr.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }* %recv.addr
  %t0 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }* @ir.b to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }*
  %t1 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }* %t0, i32 0, i32 0
  %t2 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t1, align 8
  %t3 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2)
  %taddr4 = alloca { i8*, i64 } , align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.181, i32 0, i32 0), i64 6 }, { i8*, i64 }* %taddr5, align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t3, { i8*, i64 }* %taddr6, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr4, { i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6)
  %t7 = load { i8*, i64 }, { i8*, i64 }* %taddr4, align 8
  %taddr8 = alloca { i8*, i64 } , align 8
  %taddr9 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t7, { i8*, i64 }* %taddr9, align 8
  %taddr10 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([10 x i8], [10 x i8]* @.str.182, i32 0, i32 0), i64 9 }, { i8*, i64 }* %taddr10, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr8, { i8*, i64 }* %taddr9, { i8*, i64 }* %taddr10)
  %t11 = load { i8*, i64 }, { i8*, i64 }* %taddr8, align 8
  %t12 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }* @ir.b to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }*
  %t13 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }* %t12, i32 0, i32 1
  %t14 = load { i8*, i64 }, { i8*, i64 }* %t13, align 8
  %taddr15 = alloca { i8*, i64 } , align 8
  %taddr16 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t11, { i8*, i64 }* %taddr16, align 8
  %taddr17 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t14, { i8*, i64 }* %taddr17, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr15, { i8*, i64 }* %taddr16, { i8*, i64 }* %taddr17)
  %t18 = load { i8*, i64 }, { i8*, i64 }* %taddr15, align 8
  %taddr19 = alloca { i8*, i64 } , align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t18, { i8*, i64 }* %taddr20, align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([10 x i8], [10 x i8]* @.str.183, i32 0, i32 0), i64 9 }, { i8*, i64 }* %taddr21, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr19, { i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21)
  %t22 = load { i8*, i64 }, { i8*, i64 }* %taddr19, align 8
  %t23 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }* @ir.b to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }*
  %t24 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { i8*, i64 }, { i8*, i64 } }* %t23, i32 0, i32 2
  %t25 = load { i8*, i64 }, { i8*, i64 }* %t24, align 8
  %taddr26 = alloca { i8*, i64 } , align 8
  %taddr27 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t22, { i8*, i64 }* %taddr27, align 8
  %taddr28 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t25, { i8*, i64 }* %taddr28, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr26, { i8*, i64 }* %taddr27, { i8*, i64 }* %taddr28)
  %t29 = load { i8*, i64 }, { i8*, i64 }* %taddr26, align 8
  ret { i8*, i64 } %t29
}

define void @ir.CallVoid.isInstruction({ { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } %recv) {
entry:
  %recv.addr = alloca { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } , align 8
  store { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } %recv, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.CallVoid.String({ { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } %recv) {
entry:
  %recv.addr = alloca { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } , align 8
  store { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } } %recv, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %recv.addr
  %t0 = bitcast { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* @ir.c to { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }*
  %t1 = getelementptr inbounds { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %t0, i32 0, i32 1
  %t2 = load { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t1, align 8
  %taddr3 = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } %t2, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr3, align 8
  %t4 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr3, i32 0, i32 1
  %t5 = load i64, i64* %t4, align 8
  %t6 = call i8* @gominic_makeSlice(i64 %t5, i64 %t5, i64 16)
  %t7 = bitcast i8* %t6 to { i8*, i64 }*
  %taddr8 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %t9 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 0
  store { i8*, i64 }* %t7, { i8*, i64 }** %t9, align 8
  %t10 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 1
  store i64 %t5, i64* %t10, align 8
  %t11 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, i32 0, i32 2
  store i64 %t5, i64* %t11, align 8
  %t12 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %taddr8, align 8
  %args.addr13 = alloca { { i8*, i64 }*, i64, i64 }
  %taddr14 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  store { { i8*, i64 }*, i64, i64 } %t12, { { i8*, i64 }*, i64, i64 }* %taddr14
  call void @gominic_memcpy({ { i8*, i64 }*, i64, i64 }* %args.addr13, { { i8*, i64 }*, i64, i64 }* %taddr14, i64 24)
  %i.addr15 = alloca i64
  store i64 0, i64* %i.addr15, align 8
  %taddr20 = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } , align 8
  %taddr39 = alloca { i8*, i64 } , align 8
  %taddr40 = alloca { i8*, i64 } , align 8
  %taddr41 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca { i8*, i64 } , align 8
  %taddr46 = alloca { i8*, i64 } , align 8
  %taddr47 = alloca { i8*, i64 } , align 8
  %taddr58 = alloca { i8*, i64 } , align 8
  %taddr64 = alloca { i8*, i64 } , align 8
  %taddr65 = alloca { i8*, i64 } , align 8
  %taddr66 = alloca { i8*, i64 } , align 8
  %taddr68 = alloca { i8*, i64 } , align 8
  %taddr69 = alloca { i8*, i64 } , align 8
  %taddr70 = alloca { i8*, i64 } , align 8
  %taddr73 = alloca { { i8*, i64 }*, i64, i64 } , align 8
  %taddr74 = alloca { i8*, i64 } , align 8
  %taddr76 = alloca { i8*, i64 } , align 8
  %taddr77 = alloca { i8*, i64 } , align 8
  %taddr78 = alloca { i8*, i64 } , align 8
  %taddr80 = alloca { i8*, i64 } , align 8
  %taddr81 = alloca { i8*, i64 } , align 8
  %taddr82 = alloca { i8*, i64 } , align 8
  br label %for.cond0
for.cond0:
  %t16 = load i64, i64* %i.addr15, align 8
  %t17 = bitcast { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* @ir.c to { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }*
  %t18 = getelementptr inbounds { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %t17, i32 0, i32 1
  %t19 = load { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t18, align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } %t19, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr20, align 8
  %t21 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %taddr20, i32 0, i32 1
  %t22 = load i64, i64* %t21, align 8
  %t23 = icmp slt i64 %t16, %t22
  br i1 %t23, label %for.body1, label %for.end3
for.body1:
  %t24 = bitcast { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* @ir.c to { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }*
  %t25 = getelementptr inbounds { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %t24, i32 0, i32 1
  %t26 = load i64, i64* %i.addr15, align 8
  %t27 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t25, i32 0, i32 1
  %t28 = load i64, i64* %t27, align 8
  %t29 = icmp sge i64 %t26, 0
  %t30 = icmp slt i64 %t26, %t28
  %t31 = and i1 %t29, %t30
  br i1 %t31, label %idx.ok4, label %idx.fail5
for.post2:
  %t59 = load i64, i64* %i.addr15, align 8
  %t60 = add i64 %t59, 1
  store i64 %t60, i64* %i.addr15, align 8
  br label %for.cond0
for.end3:
  %t61 = bitcast { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* @ir.c to { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }*
  %t62 = getelementptr inbounds { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 } }* %t61, i32 0, i32 0
  %t63 = load { i8*, i64 }, { i8*, i64 }* %t62, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([12 x i8], [12 x i8]* @.str.185, i32 0, i32 0), i64 11 }, { i8*, i64 }* %taddr65, align 8
  store { i8*, i64 } %t63, { i8*, i64 }* %taddr66, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr64, { i8*, i64 }* %taddr65, { i8*, i64 }* %taddr66)
  %t67 = load { i8*, i64 }, { i8*, i64 }* %taddr64, align 8
  store { i8*, i64 } %t67, { i8*, i64 }* %taddr69, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.186, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr70, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr68, { i8*, i64 }* %taddr69, { i8*, i64 }* %taddr70)
  %t71 = load { i8*, i64 }, { i8*, i64 }* %taddr68, align 8
  %t72 = load { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %args.addr13, align 8
  store { { i8*, i64 }*, i64, i64 } %t72, { { i8*, i64 }*, i64, i64 }* %taddr73
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.187, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr74
  %t75 = call { i8*, i64 } @strings.Join({ { i8*, i64 }*, i64, i64 }* %taddr73, { i8*, i64 }* %taddr74)
  store { i8*, i64 } %t71, { i8*, i64 }* %taddr77, align 8
  store { i8*, i64 } %t75, { i8*, i64 }* %taddr78, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr76, { i8*, i64 }* %taddr77, { i8*, i64 }* %taddr78)
  %t79 = load { i8*, i64 }, { i8*, i64 }* %taddr76, align 8
  store { i8*, i64 } %t79, { i8*, i64 }* %taddr81, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.188, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr82, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr80, { i8*, i64 }* %taddr81, { i8*, i64 }* %taddr82)
  %t83 = load { i8*, i64 }, { i8*, i64 }* %taddr80, align 8
  ret { i8*, i64 } %t83
idx.ok4:
  %t32 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }* %t25, i32 0, i32 0
  %t33 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*** %t32, align 8
  %t34 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t33, i64 %t26
  %t35 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t34, align 8
  %a.addr36 = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t35, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36, align 8
  %t37 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36, align 8
  %t38 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t37)
  store { i8*, i64 } %t38, { i8*, i64 }* %taddr40, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.184, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr41, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr39, { i8*, i64 }* %taddr40, { i8*, i64 }* %taddr41)
  %t42 = load { i8*, i64 }, { i8*, i64 }* %taddr39, align 8
  %t43 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %a.addr36, align 8
  %t44 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t43)
  store { i8*, i64 } %t42, { i8*, i64 }* %taddr46, align 8
  store { i8*, i64 } %t44, { i8*, i64 }* %taddr47, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr45, { i8*, i64 }* %taddr46, { i8*, i64 }* %taddr47)
  %t48 = load { i8*, i64 }, { i8*, i64 }* %taddr45, align 8
  %t49 = load i64, i64* %i.addr15, align 8
  %t50 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %args.addr13, i32 0, i32 1
  %t51 = load i64, i64* %t50, align 8
  %t52 = icmp sge i64 %t49, 0
  %t53 = icmp slt i64 %t49, %t51
  %t54 = and i1 %t52, %t53
  br i1 %t54, label %idx.ok6, label %idx.fail7
idx.fail5:
  call void @gominic_abort()
  br label %idx.ok4
idx.ok6:
  %t55 = getelementptr inbounds { { i8*, i64 }*, i64, i64 }, { { i8*, i64 }*, i64, i64 }* %args.addr13, i32 0, i32 0
  %t56 = load { i8*, i64 }*, { i8*, i64 }** %t55, align 8
  %t57 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %t56, i64 %t49
  store { i8*, i64 } %t48, { i8*, i64 }* %taddr58
  call void @gominic_memcpy({ i8*, i64 }* %t57, { i8*, i64 }* %taddr58, i64 16)
  br label %for.post2
idx.fail7:
  call void @gominic_abort()
  br label %idx.ok6
}

define { { i8*, i64 }, { i8**, i64, i64 }, i8* }* @ir.NewBasicBlock({ i8*, i64 }* byval({ i8*, i64 }) %name) {
entry:
  %name.addr = alloca { i8*, i64 } , align 8
  call void @gominic_memcpy({ i8*, i64 }* %name.addr, { i8*, i64 }* %name, i64 16)
  %t0 = load { i8*, i64 }, { i8*, i64 }* %name.addr, align 8
  %taddr1 = alloca { { i8*, i64 }, { i8**, i64, i64 }, i8* } , align 8
  %t2 = getelementptr inbounds { { i8*, i64 }, { i8**, i64, i64 }, i8* }, { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %taddr1, i32 0, i32 0
  store { i8*, i64 } %t0, { i8*, i64 }* %t2, align 8
  %t3 = getelementptr inbounds { { i8*, i64 }, { i8**, i64, i64 }, i8* }, { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %taddr1, i32 0, i32 1
  store { i8**, i64, i64 } zeroinitializer, { i8**, i64, i64 }* %t3, align 8
  %t4 = getelementptr inbounds { { i8*, i64 }, { i8**, i64, i64 }, i8* }, { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %taddr1, i32 0, i32 2
  store i8* null, i8** %t4, align 8
  %t5 = load { { i8*, i64 }, { i8**, i64, i64 }, i8* }, { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %taddr1, align 8
  %taddr6 = alloca { { i8*, i64 }, { i8**, i64, i64 }, i8* } , align 8
  store { { i8*, i64 }, { i8**, i64, i64 }, i8* } %t5, { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %taddr6, align 8
  ret { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %taddr6
}

define void @ir.BasicBlock.Append({ { i8*, i64 }, { i8**, i64, i64 }, i8* }* %recv, i8* %inst) {
entry:
  %recv.addr = alloca { { i8*, i64 }, { i8**, i64, i64 }, i8* }* , align 8
  store { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %recv, { { i8*, i64 }, { i8**, i64, i64 }, i8* }** %recv.addr
  %inst.addr = alloca i8* , align 8
  store i8* %inst, i8** %inst.addr
  %t0 = bitcast { { i8*, i64 }, { i8**, i64, i64 }, i8* }** @ir.bb to { { i8*, i64 }, { i8**, i64, i64 }, i8* }*
  %t1 = getelementptr inbounds { { i8*, i64 }, { i8**, i64, i64 }, i8* }, { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %t0, i32 0, i32 1
  %t2 = load { i8**, i64, i64 }, { i8**, i64, i64 }* %t1, align 8
  %t3 = load i8*, i8** %inst.addr, align 8
  %taddr4 = alloca { i8**, i64, i64 } , align 8
  store { i8**, i64, i64 } %t2, { i8**, i64, i64 }* %taddr4, align 8
  %t5 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr4, i32 0, i32 0
  %t6 = load i8**, i8*** %t5, align 8
  %t7 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr4, i32 0, i32 1
  %t8 = load i64, i64* %t7, align 8
  %t9 = add i64 %t8, 1
  %t10 = mul i64 %t8, 8
  %t11 = call i8* @gominic_makeSlice(i64 %t9, i64 %t9, i64 8)
  %t12 = bitcast i8* %t11 to i8**
  call void @gominic_memcpy(i8** %t12, i8** %t6, i64 %t10)
  %t13 = getelementptr inbounds i8*, i8** %t12, i64 %t8
  store i8* %t3, i8** %t13, align 8
  %taddr14 = alloca { i8**, i64, i64 } , align 8
  %t15 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr14, i32 0, i32 0
  store i8** %t12, i8*** %t15, align 8
  %t16 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr14, i32 0, i32 1
  store i64 %t9, i64* %t16, align 8
  %t17 = getelementptr inbounds { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr14, i32 0, i32 2
  store i64 %t9, i64* %t17, align 8
  %t18 = load { i8**, i64, i64 }, { i8**, i64, i64 }* %taddr14, align 8
  %t19 = bitcast { { i8*, i64 }, { i8**, i64, i64 }, i8* }** @ir.bb to { { i8*, i64 }, { i8**, i64, i64 }, i8* }*
  %t20 = getelementptr inbounds { { i8*, i64 }, { i8**, i64, i64 }, i8* }, { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %t19, i32 0, i32 1
  %taddr21 = alloca { i8**, i64, i64 } , align 8
  store { i8**, i64, i64 } %t18, { i8**, i64, i64 }* %taddr21
  call void @gominic_memcpy({ i8**, i64, i64 }* %t20, { i8**, i64, i64 }* %taddr21, i64 24)
  ret void
}

define { { i8*, i64 }, { i8**, i64, i64 }, i8* }* @ir.Function.Entry({ { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* %recv) {
entry:
  %recv.addr = alloca { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* , align 8
  store { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* %recv, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** %recv.addr
  %t0 = bitcast { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** @ir.fn to { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }*
  %t1 = getelementptr inbounds { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* %t0, i32 0, i32 3
  %t2 = load { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %t1, align 8
  %taddr3 = alloca { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 } , align 8
  store { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 } %t2, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr3, align 8
  %t4 = getelementptr inbounds { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr3, i32 0, i32 1
  %t5 = load i64, i64* %t4, align 8
  %t6 = icmp eq i64 %t5, 0
  %taddr7 = alloca { i8*, i64 } , align 8
  %taddr14 = alloca { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 } , align 8
  %taddr24 = alloca { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 } , align 8
  %taddr31 = alloca { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 } , align 8
  br i1 %t6, label %then0, label %endif1
then0:
  store { i8*, i64 } { i8* getelementptr inbounds ([6 x i8], [6 x i8]* @.str.189, i32 0, i32 0), i64 5 }, { i8*, i64 }* %taddr7
  %t8 = call { { i8*, i64 }, { i8**, i64, i64 }, i8* }* @ir.NewBasicBlock({ i8*, i64 }* %taddr7)
  %entry.addr9 = alloca { { i8*, i64 }, { i8**, i64, i64 }, i8* }*
  store { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %t8, { { i8*, i64 }, { i8**, i64, i64 }, i8* }** %entry.addr9, align 8
  %t10 = bitcast { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** @ir.fn to { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }*
  %t11 = getelementptr inbounds { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* %t10, i32 0, i32 3
  %t12 = load { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %t11, align 8
  %t13 = load { { i8*, i64 }, { i8**, i64, i64 }, i8* }*, { { i8*, i64 }, { i8**, i64, i64 }, i8* }** %entry.addr9, align 8
  store { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 } %t12, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr14, align 8
  %t15 = getelementptr inbounds { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr14, i32 0, i32 0
  %t16 = load { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, { { i8*, i64 }, { i8**, i64, i64 }, i8* }*** %t15, align 8
  %t17 = getelementptr inbounds { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr14, i32 0, i32 1
  %t18 = load i64, i64* %t17, align 8
  %t19 = add i64 %t18, 1
  %t20 = mul i64 %t18, 8
  %t21 = call i8* @gominic_makeSlice(i64 %t19, i64 %t19, i64 8)
  %t22 = bitcast i8* %t21 to { { i8*, i64 }, { i8**, i64, i64 }, i8* }**
  call void @gominic_memcpy({ { i8*, i64 }, { i8**, i64, i64 }, i8* }** %t22, { { i8*, i64 }, { i8**, i64, i64 }, i8* }** %t16, i64 %t20)
  %t23 = getelementptr inbounds { { i8*, i64 }, { i8**, i64, i64 }, i8* }*, { { i8*, i64 }, { i8**, i64, i64 }, i8* }** %t22, i64 %t18
  store { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %t13, { { i8*, i64 }, { i8**, i64, i64 }, i8* }** %t23, align 8
  %t25 = getelementptr inbounds { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr24, i32 0, i32 0
  store { { i8*, i64 }, { i8**, i64, i64 }, i8* }** %t22, { { i8*, i64 }, { i8**, i64, i64 }, i8* }*** %t25, align 8
  %t26 = getelementptr inbounds { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr24, i32 0, i32 1
  store i64 %t19, i64* %t26, align 8
  %t27 = getelementptr inbounds { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr24, i32 0, i32 2
  store i64 %t19, i64* %t27, align 8
  %t28 = load { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr24, align 8
  %t29 = bitcast { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** @ir.fn to { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }*
  %t30 = getelementptr inbounds { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* %t29, i32 0, i32 3
  store { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 } %t28, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr31
  call void @gominic_memcpy({ { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %t30, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %taddr31, i64 24)
  br label %endif1
endif1:
  %t32 = bitcast { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** @ir.fn to { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }*
  %t33 = getelementptr inbounds { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* %t32, i32 0, i32 3
  %t34 = getelementptr inbounds { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %t33, i32 0, i32 1
  %t35 = load i64, i64* %t34, align 8
  %t36 = icmp sge i64 0, 0
  %t37 = icmp slt i64 0, %t35
  %t38 = and i1 %t36, %t37
  br i1 %t38, label %idx.ok2, label %idx.fail3
idx.ok2:
  %t39 = getelementptr inbounds { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }* %t33, i32 0, i32 0
  %t40 = load { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, { { i8*, i64 }, { i8**, i64, i64 }, i8* }*** %t39, align 8
  %t41 = getelementptr inbounds { { i8*, i64 }, { i8**, i64, i64 }, i8* }*, { { i8*, i64 }, { i8**, i64, i64 }, i8* }** %t40, i64 0
  %t42 = load { { i8*, i64 }, { i8**, i64, i64 }, i8* }*, { { i8*, i64 }, { i8**, i64, i64 }, i8* }** %t41, align 8
  ret { { i8*, i64 }, { i8**, i64, i64 }, i8* }* %t42
idx.fail3:
  call void @gominic_abort()
  br label %idx.ok2
}

define void @ir.Module.AddFunction({ { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }* %recv, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* %fn) {
entry:
  %recv.addr = alloca { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }* , align 8
  store { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }* %recv, { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }** %recv.addr
  %fn.addr = alloca { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* , align 8
  store { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* %fn, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** %fn.addr
  %t0 = bitcast { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }** @ir.m to { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }*
  %t1 = getelementptr inbounds { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }, { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }* %t0, i32 0, i32 0
  %t2 = load { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %t1, align 8
  %t3 = load { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }*, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** %fn.addr, align 8
  %taddr4 = alloca { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 } , align 8
  store { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 } %t2, { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %taddr4, align 8
  %t5 = getelementptr inbounds { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %taddr4, i32 0, i32 0
  %t6 = load { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }*** %t5, align 8
  %t7 = getelementptr inbounds { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %taddr4, i32 0, i32 1
  %t8 = load i64, i64* %t7, align 8
  %t9 = add i64 %t8, 1
  %t10 = mul i64 %t8, 8
  %t11 = call i8* @gominic_makeSlice(i64 %t9, i64 %t9, i64 8)
  %t12 = bitcast i8* %t11 to { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**
  call void @gominic_memcpy({ { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** %t12, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** %t6, i64 %t10)
  %t13 = getelementptr inbounds { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }*, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** %t12, i64 %t8
  store { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }* %t3, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** %t13, align 8
  %taddr14 = alloca { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 } , align 8
  %t15 = getelementptr inbounds { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %taddr14, i32 0, i32 0
  store { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }** %t12, { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }*** %t15, align 8
  %t16 = getelementptr inbounds { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %taddr14, i32 0, i32 1
  store i64 %t9, i64* %t16, align 8
  %t17 = getelementptr inbounds { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %taddr14, i32 0, i32 2
  store i64 %t9, i64* %t17, align 8
  %t18 = load { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %taddr14, align 8
  %t19 = bitcast { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }** @ir.m to { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }*
  %t20 = getelementptr inbounds { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }, { { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }, { { { i8*, i64 }, i8*, { i8*, i64 }, i64, i1 }*, i64, i64 }, { i8*, i64 }, { i8*, i64 } }* %t19, i32 0, i32 0
  %taddr21 = alloca { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 } , align 8
  store { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 } %t18, { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %taddr21
  call void @gominic_memcpy({ { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %t20, { { { i8*, i64 }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }**, i64, i64 }, { i8**, i64, i64 }, { { { i8*, i64 }, { i8**, i64, i64 }, i8* }**, i64, i64 }, i8* }**, i64, i64 }* %taddr21, i64 24)
  ret void
}

define void @ir.Bitcast.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.Bitcast.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %recv.addr
  %t0 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* @ir.b to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }*
  %t1 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %t0, i32 0, i32 0
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t1 to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 0
  %t4 = load { i8*, i64 }, { i8*, i64 }* %t3, align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.190, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr6, align 8
  %taddr7 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t4, { i8*, i64 }* %taddr7, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6, { i8*, i64 }* %taddr7)
  %t8 = load { i8*, i64 }, { i8*, i64 }* %taddr5, align 8
  %taddr9 = alloca { i8*, i64 } , align 8
  %taddr10 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t8, { i8*, i64 }* %taddr10, align 8
  %taddr11 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([12 x i8], [12 x i8]* @.str.191, i32 0, i32 0), i64 11 }, { i8*, i64 }* %taddr11, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr9, { i8*, i64 }* %taddr10, { i8*, i64 }* %taddr11)
  %t12 = load { i8*, i64 }, { i8*, i64 }* %taddr9, align 8
  %t13 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* @ir.b to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }*
  %t14 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %t13, i32 0, i32 1
  %t15 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t14, align 8
  %t16 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t15)
  %taddr17 = alloca { i8*, i64 } , align 8
  %taddr18 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t12, { i8*, i64 }* %taddr18, align 8
  %taddr19 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t16, { i8*, i64 }* %taddr19, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr17, { i8*, i64 }* %taddr18, { i8*, i64 }* %taddr19)
  %t20 = load { i8*, i64 }, { i8*, i64 }* %taddr17, align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  %taddr22 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t20, { i8*, i64 }* %taddr22, align 8
  %taddr23 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.192, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr23, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr21, { i8*, i64 }* %taddr22, { i8*, i64 }* %taddr23)
  %t24 = load { i8*, i64 }, { i8*, i64 }* %taddr21, align 8
  %t25 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* @ir.b to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }*
  %t26 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %t25, i32 0, i32 1
  %t27 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t26, align 8
  %t28 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t27)
  %taddr29 = alloca { i8*, i64 } , align 8
  %taddr30 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t24, { i8*, i64 }* %taddr30, align 8
  %taddr31 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t28, { i8*, i64 }* %taddr31, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr29, { i8*, i64 }* %taddr30, { i8*, i64 }* %taddr31)
  %t32 = load { i8*, i64 }, { i8*, i64 }* %taddr29, align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  %taddr34 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t32, { i8*, i64 }* %taddr34, align 8
  %taddr35 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.193, i32 0, i32 0), i64 4 }, { i8*, i64 }* %taddr35, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr33, { i8*, i64 }* %taddr34, { i8*, i64 }* %taddr35)
  %t36 = load { i8*, i64 }, { i8*, i64 }* %taddr33, align 8
  %t37 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* @ir.b to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }*
  %t38 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, i8* }* %t37, i32 0, i32 2
  %t39 = load i8*, i8** %t38, align 8
  %t40 = call { i8*, i64 } @ir.Type.String(i8* %t39)
  %taddr41 = alloca { i8*, i64 } , align 8
  %taddr42 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t36, { i8*, i64 }* %taddr42, align 8
  %taddr43 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t40, { i8*, i64 }* %taddr43, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr41, { i8*, i64 }* %taddr42, { i8*, i64 }* %taddr43)
  %t44 = load { i8*, i64 }, { i8*, i64 }* %taddr41, align 8
  ret { i8*, i64 } %t44
}

define void @ir.Memcpy.isInstruction({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %recv.addr
  ret void
}

define { i8*, i64 } @ir.Memcpy.String({ { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv) {
entry:
  %recv.addr = alloca { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } , align 8
  store { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* } %recv, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %recv.addr
  %t0 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.m to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t1 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t0, i32 0, i32 0
  %t2 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t1, align 8
  %t3 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2)
  %taddr4 = alloca { i8*, i64 } , align 8
  %taddr5 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([27 x i8], [27 x i8]* @.str.194, i32 0, i32 0), i64 26 }, { i8*, i64 }* %taddr5, align 8
  %taddr6 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t3, { i8*, i64 }* %taddr6, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr4, { i8*, i64 }* %taddr5, { i8*, i64 }* %taddr6)
  %t7 = load { i8*, i64 }, { i8*, i64 }* %taddr4, align 8
  %taddr8 = alloca { i8*, i64 } , align 8
  %taddr9 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t7, { i8*, i64 }* %taddr9, align 8
  %taddr10 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.195, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr10, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr8, { i8*, i64 }* %taddr9, { i8*, i64 }* %taddr10)
  %t11 = load { i8*, i64 }, { i8*, i64 }* %taddr8, align 8
  %t12 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.m to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t13 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t12, i32 0, i32 0
  %t14 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t13, align 8
  %t15 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t14)
  %taddr16 = alloca { i8*, i64 } , align 8
  %taddr17 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t11, { i8*, i64 }* %taddr17, align 8
  %taddr18 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t15, { i8*, i64 }* %taddr18, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr16, { i8*, i64 }* %taddr17, { i8*, i64 }* %taddr18)
  %t19 = load { i8*, i64 }, { i8*, i64 }* %taddr16, align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t19, { i8*, i64 }* %taddr21, align 8
  %taddr22 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.196, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr22, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21, { i8*, i64 }* %taddr22)
  %t23 = load { i8*, i64 }, { i8*, i64 }* %taddr20, align 8
  %t24 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.m to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t25 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t24, i32 0, i32 1
  %t26 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t25, align 8
  %t27 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t26)
  %taddr28 = alloca { i8*, i64 } , align 8
  %taddr29 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t23, { i8*, i64 }* %taddr29, align 8
  %taddr30 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t27, { i8*, i64 }* %taddr30, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr28, { i8*, i64 }* %taddr29, { i8*, i64 }* %taddr30)
  %t31 = load { i8*, i64 }, { i8*, i64 }* %taddr28, align 8
  %taddr32 = alloca { i8*, i64 } , align 8
  %taddr33 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t31, { i8*, i64 }* %taddr33, align 8
  %taddr34 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.197, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr34, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr32, { i8*, i64 }* %taddr33, { i8*, i64 }* %taddr34)
  %t35 = load { i8*, i64 }, { i8*, i64 }* %taddr32, align 8
  %t36 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.m to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t37 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t36, i32 0, i32 1
  %t38 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t37, align 8
  %t39 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t38)
  %taddr40 = alloca { i8*, i64 } , align 8
  %taddr41 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t35, { i8*, i64 }* %taddr41, align 8
  %taddr42 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t39, { i8*, i64 }* %taddr42, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr40, { i8*, i64 }* %taddr41, { i8*, i64 }* %taddr42)
  %t43 = load { i8*, i64 }, { i8*, i64 }* %taddr40, align 8
  %taddr44 = alloca { i8*, i64 } , align 8
  %taddr45 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t43, { i8*, i64 }* %taddr45, align 8
  %taddr46 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.198, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr46, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr44, { i8*, i64 }* %taddr45, { i8*, i64 }* %taddr46)
  %t47 = load { i8*, i64 }, { i8*, i64 }* %taddr44, align 8
  %t48 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.m to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t49 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t48, i32 0, i32 2
  %t50 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t49, align 8
  %t51 = call { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t50)
  %taddr52 = alloca { i8*, i64 } , align 8
  %taddr53 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t47, { i8*, i64 }* %taddr53, align 8
  %taddr54 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t51, { i8*, i64 }* %taddr54, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr52, { i8*, i64 }* %taddr53, { i8*, i64 }* %taddr54)
  %t55 = load { i8*, i64 }, { i8*, i64 }* %taddr52, align 8
  %taddr56 = alloca { i8*, i64 } , align 8
  %taddr57 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t55, { i8*, i64 }* %taddr57, align 8
  %taddr58 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.199, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr58, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr56, { i8*, i64 }* %taddr57, { i8*, i64 }* %taddr58)
  %t59 = load { i8*, i64 }, { i8*, i64 }* %taddr56, align 8
  %t60 = bitcast { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* @ir.m to { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }*
  %t61 = getelementptr inbounds { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }, { { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* }* %t60, i32 0, i32 2
  %t62 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %t61, align 8
  %t63 = call { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t62)
  %taddr64 = alloca { i8*, i64 } , align 8
  %taddr65 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t59, { i8*, i64 }* %taddr65, align 8
  %taddr66 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t63, { i8*, i64 }* %taddr66, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr64, { i8*, i64 }* %taddr65, { i8*, i64 }* %taddr66)
  %t67 = load { i8*, i64 }, { i8*, i64 }* %taddr64, align 8
  %taddr68 = alloca { i8*, i64 } , align 8
  %taddr69 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } %t67, { i8*, i64 }* %taddr69, align 8
  %taddr70 = alloca { i8*, i64 } , align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.200, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr70, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr68, { i8*, i64 }* %taddr69, { i8*, i64 }* %taddr70)
  %t71 = load { i8*, i64 }, { i8*, i64 }* %taddr68, align 8
  ret { i8*, i64 } %t71
}

define { i8*, i64 } @ir.formatValue({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %v) {
entry:
  %v.addr = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %v, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr
  %t0 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr, align 8
  %t1 = icmp eq { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t0, null
  %taddr10 = alloca i1 , align 1
  %taddr19 = alloca { i8*, i64 } , align 8
  %taddr20 = alloca { i8*, i64 } , align 8
  %taddr21 = alloca { i8*, i64 } , align 8
  %taddr30 = alloca { i8*, i64 } , align 8
  %taddr31 = alloca { i8*, i64 } , align 8
  %taddr42 = alloca i1 , align 1
  %taddr46 = alloca { i8*, i64 } , align 8
  %taddr47 = alloca { i8*, i64 } , align 8
  %taddr67 = alloca { i8*, i64 } , align 8
  %taddr68 = alloca { i8*, i64 } , align 8
  %taddr78 = alloca i1 , align 1
  %taddr83 = alloca { i8*, i64 } , align 8
  %taddr84 = alloca { i8*, i64 } , align 8
  %taddr98 = alloca { i8*, i64 } , align 8
  %taddr99 = alloca { i8*, i64 } , align 8
  %taddr105 = alloca { i8*, i64 } , align 8
  %taddr106 = alloca { i8*, i64 } , align 8
  %taddr107 = alloca { i8*, i64 } , align 8
  br i1 %t1, label %then0, label %endif1
then0:
  ret { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.201, i32 0, i32 0), i64 0 }
endif1:
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 2
  %t4 = load i64, i64* %t3, align 8
  %t5 = icmp eq i64 %t4, 2
  %t6 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t7 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t6, i32 0, i32 2
  %t8 = load i64, i64* %t7, align 8
  %t9 = icmp eq i64 %t8, 1
  br i1 %t5, label %logic.true3, label %logic.rhs4
logic.end2:
  %t15 = load i1, i1* %taddr10, align 1
  br i1 %t15, label %then5, label %endif6
logic.true3:
  store i1 1, i1* %taddr10
  br label %logic.end2
logic.rhs4:
  %t11 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t12 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t11, i32 0, i32 2
  %t13 = load i64, i64* %t12, align 8
  %t14 = icmp eq i64 %t13, 1
  store i1 %t14, i1* %taddr10
  br label %logic.end2
then5:
  %t16 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t17 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t16, i32 0, i32 0
  %t18 = load { i8*, i64 }, { i8*, i64 }* %t17, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([2 x i8], [2 x i8]* @.str.202, i32 0, i32 0), i64 1 }, { i8*, i64 }* %taddr20, align 8
  store { i8*, i64 } %t18, { i8*, i64 }* %taddr21, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr19, { i8*, i64 }* %taddr20, { i8*, i64 }* %taddr21)
  %t22 = load { i8*, i64 }, { i8*, i64 }* %taddr19, align 8
  ret { i8*, i64 } %t22
endif6:
  %t23 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t24 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t23, i32 0, i32 2
  %t25 = load i64, i64* %t24, align 8
  %t26 = icmp eq i64 %t25, 3
  %t27 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t28 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t27, i32 0, i32 3
  %t29 = load { i8*, i64 }, { i8*, i64 }* %t28, align 8
  store { i8*, i64 } %t29, { i8*, i64 }* %taddr30, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.203, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr31, align 8
  %t32 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr30, i32 0, i32 0
  %t33 = load i8*, i8** %t32, align 8
  %t34 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr30, i32 0, i32 1
  %t35 = load i64, i64* %t34, align 8
  %t36 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr31, i32 0, i32 0
  %t37 = load i8*, i8** %t36, align 8
  %t38 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr31, i32 0, i32 1
  %t39 = load i64, i64* %t38, align 8
  %t40 = call i1 @gominic_str_eq(i8* %t33, i64 %t35, i8* %t37, i64 %t39)
  %t41 = icmp eq i1 %t40, 0
  br i1 %t26, label %logic.rhs8, label %logic.false9
logic.end7:
  %t58 = load i1, i1* %taddr42, align 1
  br i1 %t58, label %then10, label %endif11
logic.rhs8:
  %t43 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t44 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t43, i32 0, i32 3
  %t45 = load { i8*, i64 }, { i8*, i64 }* %t44, align 8
  store { i8*, i64 } %t45, { i8*, i64 }* %taddr46, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([1 x i8], [1 x i8]* @.str.204, i32 0, i32 0), i64 0 }, { i8*, i64 }* %taddr47, align 8
  %t48 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr46, i32 0, i32 0
  %t49 = load i8*, i8** %t48, align 8
  %t50 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr46, i32 0, i32 1
  %t51 = load i64, i64* %t50, align 8
  %t52 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr47, i32 0, i32 0
  %t53 = load i8*, i8** %t52, align 8
  %t54 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr47, i32 0, i32 1
  %t55 = load i64, i64* %t54, align 8
  %t56 = call i1 @gominic_str_eq(i8* %t49, i64 %t51, i8* %t53, i64 %t55)
  %t57 = icmp eq i1 %t56, 0
  store i1 %t57, i1* %taddr42
  br label %logic.end7
logic.false9:
  store i1 0, i1* %taddr42
  br label %logic.end7
then10:
  %t59 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t60 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t59, i32 0, i32 1
  %t61 = load i8*, i8** %t60, align 8
  %t62 = icmp ne i8* %t61, null
  %t63 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t64 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t63, i32 0, i32 1
  %t65 = load i8*, i8** %t64, align 8
  %t66 = call { i8*, i64 } @ir.Type.String(i8* %t65)
  store { i8*, i64 } %t66, { i8*, i64 }* %taddr67, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.205, i32 0, i32 0), i64 6 }, { i8*, i64 }* %taddr68, align 8
  %t69 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr67, i32 0, i32 0
  %t70 = load i8*, i8** %t69, align 8
  %t71 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr67, i32 0, i32 1
  %t72 = load i64, i64* %t71, align 8
  %t73 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr68, i32 0, i32 0
  %t74 = load i8*, i8** %t73, align 8
  %t75 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr68, i32 0, i32 1
  %t76 = load i64, i64* %t75, align 8
  %t77 = call i1 @gominic_str_eq(i8* %t70, i64 %t72, i8* %t74, i64 %t76)
  br i1 %t62, label %logic.rhs13, label %logic.false14
endif11:
  %t112 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t113 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t112, i32 0, i32 0
  %t114 = load { i8*, i64 }, { i8*, i64 }* %t113, align 8
  ret { i8*, i64 } %t114
logic.end12:
  %t94 = load i1, i1* %taddr78, align 1
  br i1 %t94, label %then15, label %endif16
logic.rhs13:
  %t79 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t80 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t79, i32 0, i32 1
  %t81 = load i8*, i8** %t80, align 8
  %t82 = call { i8*, i64 } @ir.Type.String(i8* %t81)
  store { i8*, i64 } %t82, { i8*, i64 }* %taddr83, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([7 x i8], [7 x i8]* @.str.206, i32 0, i32 0), i64 6 }, { i8*, i64 }* %taddr84, align 8
  %t85 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr83, i32 0, i32 0
  %t86 = load i8*, i8** %t85, align 8
  %t87 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr83, i32 0, i32 1
  %t88 = load i64, i64* %t87, align 8
  %t89 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr84, i32 0, i32 0
  %t90 = load i8*, i8** %t89, align 8
  %t91 = getelementptr inbounds { i8*, i64 }, { i8*, i64 }* %taddr84, i32 0, i32 1
  %t92 = load i64, i64* %t91, align 8
  %t93 = call i1 @gominic_str_eq(i8* %t86, i64 %t88, i8* %t90, i64 %t92)
  store i1 %t93, i1* %taddr78
  br label %logic.end12
logic.false14:
  store i1 0, i1* %taddr78
  br label %logic.end12
then15:
  %t95 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t96 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t95, i32 0, i32 3
  %t97 = load { i8*, i64 }, { i8*, i64 }* %t96, align 8
  store { i8*, i64 } %t97, { i8*, i64 }* %taddr98
  store { i8*, i64 } { i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.207, i32 0, i32 0), i64 3 }, { i8*, i64 }* %taddr99
  %t100 = call i1 @strings.ContainsAny({ i8*, i64 }* %taddr98, { i8*, i64 }* %taddr99)
  %t101 = icmp eq i1 %t100, 0
  br i1 %t101, label %then17, label %endif18
endif16:
  %t109 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t110 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t109, i32 0, i32 3
  %t111 = load { i8*, i64 }, { i8*, i64 }* %t110, align 8
  ret { i8*, i64 } %t111
then17:
  %t102 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t103 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t102, i32 0, i32 3
  %t104 = load { i8*, i64 }, { i8*, i64 }* %t103, align 8
  store { i8*, i64 } %t104, { i8*, i64 }* %taddr106, align 8
  store { i8*, i64 } { i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.208, i32 0, i32 0), i64 2 }, { i8*, i64 }* %taddr107, align 8
  call void @gominic_str_concat({ i8*, i64 }* %taddr105, { i8*, i64 }* %taddr106, { i8*, i64 }* %taddr107)
  %t108 = load { i8*, i64 }, { i8*, i64 }* %taddr105, align 8
  ret { i8*, i64 } %t108
endif18:
  br label %endif16
}

define { i8*, i64 } @ir.typeString({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %v) {
entry:
  %v.addr = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %v, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr
  %t0 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr, align 8
  %t1 = icmp ne { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t0, null
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 1
  %t4 = load i8*, i8** %t3, align 8
  %t5 = icmp ne i8* %t4, null
  %taddr6 = alloca i1 , align 1
  br i1 %t1, label %logic.rhs1, label %logic.false2
logic.end0:
  %t11 = load i1, i1* %taddr6, align 1
  br i1 %t11, label %then3, label %endif4
logic.rhs1:
  %t7 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t8 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t7, i32 0, i32 1
  %t9 = load i8*, i8** %t8, align 8
  %t10 = icmp ne i8* %t9, null
  store i1 %t10, i1* %taddr6
  br label %logic.end0
logic.false2:
  store i1 0, i1* %taddr6
  br label %logic.end0
then3:
  %t12 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t13 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t12, i32 0, i32 1
  %t14 = load i8*, i8** %t13, align 8
  %t15 = call { i8*, i64 } @ir.Type.String(i8* %t14)
  ret { i8*, i64 } %t15
endif4:
  ret { i8*, i64 } { i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.209, i32 0, i32 0), i64 4 }
}

define i8* @ir.ValueType({ { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %v) {
entry:
  %v.addr = alloca { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* , align 8
  store { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %v, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr
  %t0 = load { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr, align 8
  %t1 = icmp eq { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t0, null
  br i1 %t1, label %then0, label %endif1
then0:
  ret i8* null
endif1:
  %t2 = bitcast { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }** %v.addr to { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }*
  %t3 = getelementptr inbounds { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }, { { i8*, i64 }, i8*, i64, { i8*, i64 }, i1, i8* }* %t2, i32 0, i32 1
  %t4 = load i8*, i8** %t3, align 8
  ret i8* %t4
}

