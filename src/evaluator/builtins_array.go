package evaluator

import (
\t\"fmt\"
\t\"strings\"
\t\"github.com/uncode/logger\"
\t\"github.com/uncode/object\"
)

// log variables for debugging
var builtinLogLevel = logger.LevelTrace

// logIfEnabled logs a message if logging is enabled for the given level
func logIfEnabled(level logger.LogLevel, format string, args ...interface{}) {
\tif logger.IsSpecialLevelEnabled(level) || logger.GetComponentLevel(logger.ComponentBuiltin) >= level {
\t\tlogger.ComponentDebug(logger.ComponentBuiltin, format, args...)
\t}
}

// registerArrayBuiltins oM¢#nDº¢pí{2Yã
func registerArrayBuiltins() {
\t// Mí#PWfáWkYã¢p
\tBuiltins[\"join\"] = &object.Builtin{
\t\tName: \"join\",
\t\tFn: func(args ...object.Object) object.Object {
\t\t\tlogIfEnabled(builtinLogLevel, \"join¢pL|s˙Uå~W_: pp=%d\", len(args))
\t\t\t
\t\t\tif len(args) != 2 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"join¢po2dnpL≈ÅgY: %dHâå~W_\", len(args))
\t\t\t\treturn createError(\"join¢po2dnpL≈ÅgY: %dHâå~W_\", len(args))
\t\t\t}
\t\t\t
\t\t\t// ,1poM
\t\t\tif args[0].Type() != object.ARRAY_OBJ {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"join¢pn,1poMgBã≈ÅLBä~Y: %s\", args[0].Type())
\t\t\t\treturn createError(\"join¢pn,1poMgBã≈ÅLBä~Y: %s\", args[0].Type())
\t\t\t}
\t\t\tarray, _ := args[0].(*object.Array)
\t\t\t
\t\t\t// ,2po:äáW
\t\t\tif args[1].Type() != object.STRING_OBJ {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"join¢pn,2poáWgBã≈ÅLBä~Y: %s\", args[1].Type())
\t\t\t\treturn createError(\"join¢pn,2poáWgBã≈ÅLBä~Y: %s\", args[1].Type())
\t\t\t}
\t\t\tdelimiter, _ := args[1].(*object.String)
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"join¢p: MÅ p=%d, :äáW='%s'\", len(array.Elements), delimiter.Value)
\t\t\t
\t\t\t// MnÅ íáWk	€
\t\t\telements := make([]string, len(array.Elements))
\t\t\tfor i, elem := range array.Elements {
\t\t\t\tswitch e := elem.(type) {
\t\t\t\tcase *object.String:
\t\t\t\t\telements[i] = e.Value
\t\t\t\tcase *object.Integer:
\t\t\t\t\telements[i] = fmt.Sprintf(\"%d\", e.Value)
\t\t\t\tcase *object.Boolean:
\t\t\t\t\telements[i] = fmt.Sprintf(\"%t\", e.Value)
\t\t\t\tdefault:
\t\t\t\t\telements[i] = e.Inspect()
\t\t\t\t}
\t\t\t}
\t\t\t
\t\t\tresult := strings.Join(elements, delimiter.Value)
\t\t\tlogIfEnabled(builtinLogLevel, \"join¢p: Pú='%s'\", result)
\t\t\treturn &object.String{Value: result}
\t\t},
\t\tReturnType: object.STRING_OBJ,
\t\tParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.STRING_OBJ},
\t}

\t// p$∑¸±Ûπí\Yã¢p
\tBuiltins[\"range\"] = &object.Builtin{
\t\tName: \"range\",
\t\tFn: func(args ...object.Object) object.Object {
\t\t\tlogIfEnabled(builtinLogLevel, \"range¢pL|s˙Uå~W_: pp=%d\", len(args))
\t\t\t
\t\t\t// pnpí¡ß√Ø: 1~_o2dnpí◊QÿQã
\t\t\tif len(args) < 1 || len(args) > 2 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"range¢po1-2npL≈ÅgY: %dHâå~W_\", len(args))
\t\t\t\treturn createError(\"range¢po1-2npL≈ÅgY: %dHâå~W_\", len(args))
\t\t\t}
\t\t\t
\t\t\tvar start, end int64
\t\t\t
\t\t\t// 1dnpn4: 0Kâ]n$~g
\t\t\tif len(args) == 1 {
\t\t\t\tif args[0].Type() != object.INTEGER_OBJ {
\t\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"range¢pnpotpgBã≈ÅLBä~Y: %s\", args[0].Type())
\t\t\t\t\treturn createError(\"range¢pnpotpgBã≈ÅLBä~Y: %s\", args[0].Type())
\t\t\t\t}
\t\t\t\tendVal, _ := args[0].(*object.Integer)
\t\t\t\t
\t\t\t\tstart = 0
\t\t\t\tend = endVal.Value
\t\t\t\tlogIfEnabled(builtinLogLevel, \"range¢p: 0Kâ%d~gnƒÚí\", end)
\t\t\t} else {
\t\t\t\t// 2dnpn4: startKâend~g
\t\t\t\tif args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
\t\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"range¢pnpotpgBã≈ÅLBä~Y\")
\t\t\t\t\treturn createError(\"range¢pnpotpgBã≈ÅLBä~Y\")
\t\t\t\t}
\t\t\t\t
\t\t\t\tstartVal, _ := args[0].(*object.Integer)
\t\t\t\tendVal, _ := args[1].(*object.Integer)
\t\t\t\t
\t\t\t\tstart = startVal.Value
\t\t\t\tend = endVal.Value
\t\t\t\tlogIfEnabled(builtinLogLevel, \"range¢p: %dKâ%d~gnƒÚí\", start, end)
\t\t\t}
\t\t\t
\t\t\t// ãÀMnLBÜMnàä'MD4oznMí‘Y
\t\t\tif start > end {
\t\t\t\tlogger.ComponentWarn(logger.ComponentBuiltin, \"range¢p: ãÀ$ %d LBÜ$ %d àä'MD_ÅzMí‘W~Y\", start, end)
\t\t\t\treturn &object.Array{Elements: []object.Object{}}
\t\t\t}
\t\t\t
\t\t\t// Mí\
\t\t\telements := make([]object.Object, end-start)
\t\t\tfor i := start; i < end; i++ {
\t\t\t\telements[i-start] = &object.Integer{Value: i}
\t\t\t}
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"range¢p: PúnMÅ p=%d\", len(elements))
\t\t\treturn &object.Array{Elements: elements}
\t\t},
\t\tReturnType: object.ARRAY_OBJ,
\t\tParamTypes: []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ},
\t}
\t
\t// map¢p - MnÅ k¢píi(Yã
\tBuiltins[\"map\"] = &object.Builtin{
\t\tName: \"map\",
\t\tFn: func(args ...object.Object) object.Object {
\t\t\tlogIfEnabled(builtinLogLevel, \"map¢pL|s˙Uå~W_: pp=%d\", len(args))
\t\t\t
\t\t\t// pnp¡ß√Ø
\t\t\tif len(args) != 2 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"map¢po2dnpL≈ÅgY: M, ¢p\")
\t\t\t\treturn createError(\"map¢po2dnpL≈ÅgY: M, ¢p\")
\t\t\t}
\t\t\t
\t\t\t// ,1pLMK¡ß√Ø
\t\t\tarr, ok := args[0].(*object.Array)
\t\t\tif !ok {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"map¢pn,1poMgBã≈ÅLBä~Y: %s\", args[0].Type())
\t\t\t\treturn createError(\"map¢pn,1poMgBã≈ÅLBä~Y: %s\", args[0].Type())
\t\t\t}
\t\t\t
\t\t\t// ,2pL¢pK¡ß√Ø
\t\t\tfn, ok := args[1].(*object.Function)
\t\t\tif !ok {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"map¢pn,2po¢pgBã≈ÅLBä~Y: %s\", args[1].Type())
\t\t\t\treturn createError(\"map¢pn,2po¢pgBã≈ÅLBä~Y: %s\", args[1].Type())
\t\t\t}
\t\t\t
\t\t\t// map¢pnpn—È·¸øozgBã≈ÅLBã
\t\t\tif len(fn.Parameters) > 0 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"map¢pk!Uå_¢po—È·¸ø¸í÷ãyMgoBä~[ì\")
\t\t\t\treturn createError(\"map¢pk!Uå_¢po—È·¸ø¸í÷ãyMgoBä~[ì\")
\t\t\t}
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"map¢p: MÅ p=%d\", len(arr.Elements))
\t\t\t
\t\t\t// PúnM
\t\t\tresultElements := make([]object.Object, 0, len(arr.Elements))
\t\t\t
\t\t\t// MnÅ k¢píi(
\t\t\tfor i, elem := range arr.Elements {
\t\t\t\tlogIfEnabled(builtinLogLevel, \"map¢p: Å  %d íÊ-: %s\", i, elem.Inspect())
\t\t\t\t
\t\t\t\t// ¢pn∞Éí·5Wf<Uk˛(nÅ í-ö
\t\t\t\textendedEnv := object.NewEnclosedEnvironment(fn.Env)
\t\t\t\textendedEnv.Set(\"<U\", elem)
\t\t\t\t
\t\t\t\t// ¢píU°
\t\t\t\tresult := Eval(fn.ASTBody, extendedEnv)
\t\t\t\t
\t\t\t\t// ®È¸¡ß√Ø
\t\t\t\tif errObj, ok := result.(*object.Error); ok {
\t\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"map¢pnÊ-k®È¸Lz: %s\", errObj.Message)
\t\t\t\t\treturn errObj
\t\t\t\t}
\t\t\t\t
\t\t\t\t// ReturnValueí¢ÛÈ√◊
\t\t\t\tif retVal, ok := result.(*object.ReturnValue); ok {
\t\t\t\t\tresult = retVal.Value
\t\t\t\t}
\t\t\t\t
\t\t\t\tlogIfEnabled(builtinLogLevel, \"map¢p: Å  %d nÊPú: %s\", i, result.Inspect())
\t\t\t\t
\t\t\t\t// PúíMk˝†
\t\t\t\tresultElements = append(resultElements, result)
\t\t\t}
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"map¢p: ÊåÜ, PúnMÅ p=%d\", len(resultElements))
\t\t\treturn &object.Array{Elements: resultElements}
\t\t},
\t\tReturnType: object.ARRAY_OBJ,
\t\tParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
\t}
\t
\t// filter¢p - aˆkÙYãÅ níΩ˙Yã
\tBuiltins[\"filter\"] = &object.Builtin{
\t\tName: \"filter\",
\t\tFn: func(args ...object.Object) object.Object {
\t\t\tlogIfEnabled(builtinLogLevel, \"filter¢pL|s˙Uå~W_: pp=%d\", len(args))
\t\t\t
\t\t\t// pnp¡ß√Ø
\t\t\tif len(args) != 2 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"filter¢po2dnpL≈ÅgY: M, ¢p\")
\t\t\t\treturn createError(\"filter¢po2dnpL≈ÅgY: M, ¢p\")
\t\t\t}
\t\t\t
\t\t\t// ,1pLMK¡ß√Ø
\t\t\tarr, ok := args[0].(*object.Array)
\t\t\tif !ok {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"filter¢pn,1poMgBã≈ÅLBä~Y: %s\", args[0].Type())
\t\t\t\treturn createError(\"filter¢pn,1poMgBã≈ÅLBä~Y: %s\", args[0].Type())
\t\t\t}
\t\t\t
\t\t\t// ,2pL¢pK¡ß√Ø
\t\t\tfn, ok := args[1].(*object.Function)
\t\t\tif !ok {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"filter¢pn,2po¢pgBã≈ÅLBä~Y: %s\", args[1].Type())
\t\t\t\treturn createError(\"filter¢pn,2po¢pgBã≈ÅLBä~Y: %s\", args[1].Type())
\t\t\t}
\t\t\t
\t\t\t// filter¢pnpn—È·¸øozgBã≈ÅLBã
\t\t\tif len(fn.Parameters) > 0 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"filter¢pk!Uå_¢po—È·¸ø¸í÷ãyMgoBä~[ì\")
\t\t\t\treturn createError(\"filter¢pk!Uå_¢po—È·¸ø¸í÷ãyMgoBä~[ì\")
\t\t\t}
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"filter¢p: MÅ p=%d\", len(arr.Elements))
\t\t\t
\t\t\t// PúnM
\t\t\tresultElements := make([]object.Object, 0)
\t\t\t
\t\t\t// MnÅ kaˆ¢píi(
\t\t\tfor i, elem := range arr.Elements {
\t\t\t\tlogIfEnabled(builtinLogLevel, \"filter¢p: Å  %d íÊ-: %s\", i, elem.Inspect())
\t\t\t\t
\t\t\t\t// ¢pn∞Éí·5Wf<Uk˛(nÅ í-ö
\t\t\t\textendedEnv := object.NewEnclosedEnvironment(fn.Env)
\t\t\t\textendedEnv.Set(\"<U\", elem)
\t\t\t\t
\t\t\t\t// aˆ¢píU°
\t\t\t\tresult := Eval(fn.ASTBody, extendedEnv)
\t\t\t\t
\t\t\t\t// ®È¸¡ß√Ø
\t\t\t\tif errObj, ok := result.(*object.Error); ok {
\t\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"filter¢pnÊ-k®È¸Lz: %s\", errObj.Message)
\t\t\t\t\treturn errObj
\t\t\t\t}
\t\t\t\t
\t\t\t\t// ReturnValueí¢ÛÈ√◊
\t\t\t\tif retVal, ok := result.(*object.ReturnValue); ok {
\t\t\t\t\tresult = retVal.Value
\t\t\t\t}
\t\t\t\t
\t\t\t\t// PúLn4Å íPúMk˝†
\t\t\t\tif isTruthy(result) {
\t\t\t\t\tlogIfEnabled(builtinLogLevel, \"filter¢p: Å  %d oaˆíÄ_W~Y\", i)
\t\t\t\t\tresultElements = append(resultElements, elem)
\t\t\t\t} else {
\t\t\t\t\tlogIfEnabled(builtinLogLevel, \"filter¢p: Å  %d oaˆíÄ_W~[ì\", i)
\t\t\t\t}
\t\t\t}
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"filter¢p: ÊåÜ, PúnMÅ p=%d\", len(resultElements))
\t\t\treturn &object.Array{Elements: resultElements}
\t\t},
\t\tReturnType: object.ARRAY_OBJ,
\t\tParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
\t}
}