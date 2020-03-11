package censo

// CBas will create new basic censorship schema -> will set to null value of field
func CBas(f string) (c C) {
	return C{Field: f}
}

// CSim will create new simple censorship schema
func CSim(f string, p interface{}) (c C) {
	return C{Field: f, Placeholder: p}
}

// CFunc will create new censorship based on placeholder func schema
func CFunc(f string, pf func(interface{}) interface{}) (c C) {
	return C{Field: f, Placeholder: pf}
}
