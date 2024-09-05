package htmlc

func isCharAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

//todo: consider adding safeFor to goutil (maybe under a different name)

type safeForMethods struct {
	buf *[]byte
	i   *int
	b   *bool
}
