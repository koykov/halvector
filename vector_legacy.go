package halvector

// ParseStr parses source string.
// DEPRECATED: use ParseString instead.
func (vec *Vector) ParseStr(s string) error {
	return vec.ParseString(s)
}

// ParseCopyStr copies source string and parse it.
// DEPRECATED: use ParseCopyString instead.
func (vec *Vector) ParseCopyStr(s string) error {
	return vec.ParseCopyString(s)
}
