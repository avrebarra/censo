package censo

func CSetToCPMap(c []C) (cfpmap CFPMap) {
	cfpmap = CFPMap{}

	for _, entry := range c {
		cfpmap[entry.Field] = entry.Placeholder
	}

	return
}
