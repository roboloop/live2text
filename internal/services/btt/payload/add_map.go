package payload

func (p Payload) AddMap(src map[string]any) Payload {
	merge(p, src)

	return p
}

func merge(dst map[string]any, src map[string]any) map[string]any {
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			dstMap, dstIsMap := dstVal.(map[string]interface{})
			srcMap, srcIsMap := srcVal.(map[string]interface{})
			if dstIsMap && srcIsMap {
				dst[key] = merge(dstMap, srcMap)
				continue
			}
		}
		dst[key] = srcVal
	}
	return dst
}
