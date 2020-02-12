package censo

func convertCSetToCPMap(c []C) (cpmap map[string]interface{}) {
	cpmap = map[string]interface{}{}

	for _, entry := range c {
		cpmap[entry.Field] = entry.Placeholder
	}

	return
}
