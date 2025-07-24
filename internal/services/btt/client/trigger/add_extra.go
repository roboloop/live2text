package trigger

func (t Trigger) addExtra(src map[string]any) Trigger {
	merge(t, src)

	return t
}

func merge(dst, src map[string]any) map[string]any {
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			dstMap, dstIsMap := dstVal.(map[string]any)
			srcMap, srcIsMap := srcVal.(map[string]any)
			if dstIsMap && srcIsMap {
				dst[key] = merge(dstMap, srcMap)
				continue
			}
		}
		dst[key] = srcVal
	}
	return dst
}
