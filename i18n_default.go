package i18n

var (
	DefaultString string  = ""
	DefaultInt    int64   = -1
	DefaultFloat  float64 = -1.0
)

func requireDefault[T any](def T, defs ...T) T {
	if len(defs) > 0 {
		return defs[0]
	}
	return def
}
