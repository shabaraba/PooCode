// getTypeName はオブジェクトの型名を返す（nilの場合は「未指定」）
func getTypeName(obj object.Object) string {
	if obj == nil {
		return \未指定\
	}
	return string(obj.Type())
}
