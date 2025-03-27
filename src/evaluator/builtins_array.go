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

// registerArrayBuiltins oM�#nD��p�{2Y�
func registerArrayBuiltins() {
\t// M�#PWf�WkY��p
\tBuiltins[\"join\"] = &object.Builtin{
\t\tName: \"join\",
\t\tFn: func(args ...object.Object) object.Object {
\t\t\tlogIfEnabled(builtinLogLevel, \"join�pL|s�U�~W_: pp=%d\", len(args))
\t\t\t
\t\t\tif len(args) != 2 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"join�po2dnpLŁgY: %dH��~W_\", len(args))
\t\t\t\treturn createError(\"join�po2dnpLŁgY: %dH��~W_\", len(args))
\t\t\t}
\t\t\t
\t\t\t// ,1poM
\t\t\tif args[0].Type() != object.ARRAY_OBJ {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"join�pn,1poMgB�ŁLB�~Y: %s\", args[0].Type())
\t\t\t\treturn createError(\"join�pn,1poMgB�ŁLB�~Y: %s\", args[0].Type())
\t\t\t}
\t\t\tarray, _ := args[0].(*object.Array)
\t\t\t
\t\t\t// ,2po:��W
\t\t\tif args[1].Type() != object.STRING_OBJ {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"join�pn,2po�WgB�ŁLB�~Y: %s\", args[1].Type())
\t\t\t\treturn createError(\"join�pn,2po�WgB�ŁLB�~Y: %s\", args[1].Type())
\t\t\t}
\t\t\tdelimiter, _ := args[1].(*object.String)
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"join�p: M� p=%d, :��W='%s'\", len(array.Elements), delimiter.Value)
\t\t\t
\t\t\t// Mn� ��Wk	�
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
\t\t\tlogIfEnabled(builtinLogLevel, \"join�p: P�='%s'\", result)
\t\t\treturn &object.String{Value: result}
\t\t},
\t\tReturnType: object.STRING_OBJ,
\t\tParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.STRING_OBJ},
\t}

\t// p$����\Y��p
\tBuiltins[\"range\"] = &object.Builtin{
\t\tName: \"range\",
\t\tFn: func(args ...object.Object) object.Object {
\t\t\tlogIfEnabled(builtinLogLevel, \"range�pL|s�U�~W_: pp=%d\", len(args))
\t\t\t
\t\t\t// pnp���ï: 1~_o2dnp��Q�Q�
\t\t\tif len(args) < 1 || len(args) > 2 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"range�po1-2npLŁgY: %dH��~W_\", len(args))
\t\t\t\treturn createError(\"range�po1-2npLŁgY: %dH��~W_\", len(args))
\t\t\t}
\t\t\t
\t\t\tvar start, end int64
\t\t\t
\t\t\t// 1dnpn4: 0K�]n$~g
\t\t\tif len(args) == 1 {
\t\t\t\tif args[0].Type() != object.INTEGER_OBJ {
\t\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"range�pnpotpgB�ŁLB�~Y: %s\", args[0].Type())
\t\t\t\t\treturn createError(\"range�pnpotpgB�ŁLB�~Y: %s\", args[0].Type())
\t\t\t\t}
\t\t\t\tendVal, _ := args[0].(*object.Integer)
\t\t\t\t
\t\t\t\tstart = 0
\t\t\t\tend = endVal.Value
\t\t\t\tlogIfEnabled(builtinLogLevel, \"range�p: 0K�%d~gn��\", end)
\t\t\t} else {
\t\t\t\t// 2dnpn4: startK�end~g
\t\t\t\tif args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
\t\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"range�pnpotpgB�ŁLB�~Y\")
\t\t\t\t\treturn createError(\"range�pnpotpgB�ŁLB�~Y\")
\t\t\t\t}
\t\t\t\t
\t\t\t\tstartVal, _ := args[0].(*object.Integer)
\t\t\t\tendVal, _ := args[1].(*object.Integer)
\t\t\t\t
\t\t\t\tstart = startVal.Value
\t\t\t\tend = endVal.Value
\t\t\t\tlogIfEnabled(builtinLogLevel, \"range�p: %dK�%d~gn��\", start, end)
\t\t\t}
\t\t\t
\t\t\t// ��MnLB�Mn��'MD4oznM��Y
\t\t\tif start > end {
\t\t\t\tlogger.ComponentWarn(logger.ComponentBuiltin, \"range�p: ��$ %d LB�$ %d ��'MD_�zM��W~Y\", start, end)
\t\t\t\treturn &object.Array{Elements: []object.Object{}}
\t\t\t}
\t\t\t
\t\t\t// M�\
\t\t\telements := make([]object.Object, end-start)
\t\t\tfor i := start; i < end; i++ {
\t\t\t\telements[i-start] = &object.Integer{Value: i}
\t\t\t}
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"range�p: P�nM� p=%d\", len(elements))
\t\t\treturn &object.Array{Elements: elements}
\t\t},
\t\tReturnType: object.ARRAY_OBJ,
\t\tParamTypes: []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ},
\t}
\t
\t// map�p - Mn� k�p�i(Y�
\tBuiltins[\"map\"] = &object.Builtin{
\t\tName: \"map\",
\t\tFn: func(args ...object.Object) object.Object {
\t\t\tlogIfEnabled(builtinLogLevel, \"map�pL|s�U�~W_: pp=%d\", len(args))
\t\t\t
\t\t\t// pnp��ï
\t\t\tif len(args) != 2 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"map�po2dnpLŁgY: M, �p\")
\t\t\t\treturn createError(\"map�po2dnpLŁgY: M, �p\")
\t\t\t}
\t\t\t
\t\t\t// ,1pLMK��ï
\t\t\tarr, ok := args[0].(*object.Array)
\t\t\tif !ok {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"map�pn,1poMgB�ŁLB�~Y: %s\", args[0].Type())
\t\t\t\treturn createError(\"map�pn,1poMgB�ŁLB�~Y: %s\", args[0].Type())
\t\t\t}
\t\t\t
\t\t\t// ,2pL�pK��ï
\t\t\tfn, ok := args[1].(*object.Function)
\t\t\tif !ok {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"map�pn,2po�pgB�ŁLB�~Y: %s\", args[1].Type())
\t\t\t\treturn createError(\"map�pn,2po�pgB�ŁLB�~Y: %s\", args[1].Type())
\t\t\t}
\t\t\t
\t\t\t// map�pnpn�����ozgB�ŁLB�
\t\t\tif len(fn.Parameters) > 0 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"map�pk!U�_�po�������֋yMgoB�~[�\")
\t\t\t\treturn createError(\"map�pk!U�_�po�������֋yMgoB�~[�\")
\t\t\t}
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"map�p: M� p=%d\", len(arr.Elements))
\t\t\t
\t\t\t// P�nM
\t\t\tresultElements := make([]object.Object, 0, len(arr.Elements))
\t\t\t
\t\t\t// Mn� k�p�i(
\t\t\tfor i, elem := range arr.Elements {
\t\t\t\tlogIfEnabled(builtinLogLevel, \"map�p: �  %d ��-: %s\", i, elem.Inspect())
\t\t\t\t
\t\t\t\t// �pn����5Wf<Uk�(n� �-�
\t\t\t\textendedEnv := object.NewEnclosedEnvironment(fn.Env)
\t\t\t\textendedEnv.Set(\"<U\", elem)
\t\t\t\t
\t\t\t\t// �p�U�
\t\t\t\tresult := Eval(fn.ASTBody, extendedEnv)
\t\t\t\t
\t\t\t\t// �����ï
\t\t\t\tif errObj, ok := result.(*object.Error); ok {
\t\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"map�pn�-k���Lz: %s\", errObj.Message)
\t\t\t\t\treturn errObj
\t\t\t\t}
\t\t\t\t
\t\t\t\t// ReturnValue������
\t\t\t\tif retVal, ok := result.(*object.ReturnValue); ok {
\t\t\t\t\tresult = retVal.Value
\t\t\t\t}
\t\t\t\t
\t\t\t\tlogIfEnabled(builtinLogLevel, \"map�p: �  %d n�P�: %s\", i, result.Inspect())
\t\t\t\t
\t\t\t\t// P��Mk��
\t\t\t\tresultElements = append(resultElements, result)
\t\t\t}
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"map�p: ���, P�nM� p=%d\", len(resultElements))
\t\t\treturn &object.Array{Elements: resultElements}
\t\t},
\t\tReturnType: object.ARRAY_OBJ,
\t\tParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
\t}
\t
\t// filter�p - a�k�Y�� n���Y�
\tBuiltins[\"filter\"] = &object.Builtin{
\t\tName: \"filter\",
\t\tFn: func(args ...object.Object) object.Object {
\t\t\tlogIfEnabled(builtinLogLevel, \"filter�pL|s�U�~W_: pp=%d\", len(args))
\t\t\t
\t\t\t// pnp��ï
\t\t\tif len(args) != 2 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"filter�po2dnpLŁgY: M, �p\")
\t\t\t\treturn createError(\"filter�po2dnpLŁgY: M, �p\")
\t\t\t}
\t\t\t
\t\t\t// ,1pLMK��ï
\t\t\tarr, ok := args[0].(*object.Array)
\t\t\tif !ok {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"filter�pn,1poMgB�ŁLB�~Y: %s\", args[0].Type())
\t\t\t\treturn createError(\"filter�pn,1poMgB�ŁLB�~Y: %s\", args[0].Type())
\t\t\t}
\t\t\t
\t\t\t// ,2pL�pK��ï
\t\t\tfn, ok := args[1].(*object.Function)
\t\t\tif !ok {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"filter�pn,2po�pgB�ŁLB�~Y: %s\", args[1].Type())
\t\t\t\treturn createError(\"filter�pn,2po�pgB�ŁLB�~Y: %s\", args[1].Type())
\t\t\t}
\t\t\t
\t\t\t// filter�pnpn�����ozgB�ŁLB�
\t\t\tif len(fn.Parameters) > 0 {
\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"filter�pk!U�_�po�������֋yMgoB�~[�\")
\t\t\t\treturn createError(\"filter�pk!U�_�po�������֋yMgoB�~[�\")
\t\t\t}
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"filter�p: M� p=%d\", len(arr.Elements))
\t\t\t
\t\t\t// P�nM
\t\t\tresultElements := make([]object.Object, 0)
\t\t\t
\t\t\t// Mn� ka��p�i(
\t\t\tfor i, elem := range arr.Elements {
\t\t\t\tlogIfEnabled(builtinLogLevel, \"filter�p: �  %d ��-: %s\", i, elem.Inspect())
\t\t\t\t
\t\t\t\t// �pn����5Wf<Uk�(n� �-�
\t\t\t\textendedEnv := object.NewEnclosedEnvironment(fn.Env)
\t\t\t\textendedEnv.Set(\"<U\", elem)
\t\t\t\t
\t\t\t\t// a��p�U�
\t\t\t\tresult := Eval(fn.ASTBody, extendedEnv)
\t\t\t\t
\t\t\t\t// �����ï
\t\t\t\tif errObj, ok := result.(*object.Error); ok {
\t\t\t\t\tlogger.ComponentError(logger.ComponentBuiltin, \"filter�pn�-k���Lz: %s\", errObj.Message)
\t\t\t\t\treturn errObj
\t\t\t\t}
\t\t\t\t
\t\t\t\t// ReturnValue������
\t\t\t\tif retVal, ok := result.(*object.ReturnValue); ok {
\t\t\t\t\tresult = retVal.Value
\t\t\t\t}
\t\t\t\t
\t\t\t\t// P�Ln4� �P�Mk��
\t\t\t\tif isTruthy(result) {
\t\t\t\t\tlogIfEnabled(builtinLogLevel, \"filter�p: �  %d oa���_W~Y\", i)
\t\t\t\t\tresultElements = append(resultElements, elem)
\t\t\t\t} else {
\t\t\t\t\tlogIfEnabled(builtinLogLevel, \"filter�p: �  %d oa���_W~[�\", i)
\t\t\t\t}
\t\t\t}
\t\t\t
\t\t\tlogIfEnabled(builtinLogLevel, \"filter�p: ���, P�nM� p=%d\", len(resultElements))
\t\t\treturn &object.Array{Elements: resultElements}
\t\t},
\t\tReturnType: object.ARRAY_OBJ,
\t\tParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
\t}
}