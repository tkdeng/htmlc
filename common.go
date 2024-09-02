package htmlc

import (
	"fmt"
)

func isCharAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

//todo: consider adding safeFor to goutil (maybe under a different name)

type safeForMethods struct {
	buf *[]byte
	i   *int
	b   *bool
}

// safeFor runs a for loop over bytes in a safer way
//
// @cb:
//   - return true to continue loop
//   - return false to break the loop
func safeFor(buf *[]byte, cb func(i *int, b func(int) byte, m safeForMethods) bool) {
	end := false
	m := safeForMethods{
		buf: buf,
		b:   &end,
	}

	for i := 0; i < len(*buf); i++ {
		m.i = &i

		if !cb(&i, func(s int) byte {
			if i+s < len(*buf) {
				return (*buf)[i+s]
			}
			return 0
		}, m) {
			break
		}

		if end {
			break
		}
	}
}

// inc increments i and will break the loop if i >= len(buf)
func (m *safeForMethods) inc(size int) {
	*m.i += size
	if *m.i >= len(*m.buf) {
		*m.b = true
	}
}

// end breaks the loop
func (m *safeForMethods) end() {
	*m.b = true
}

// loop creates an inner loop that continues to verify the array length
//
// if loop is not incramented after the callback, the loop will automatically break and log a warning
//
// @cb:
//   - return true to continue loop
//   - return false to break the loop
func (m *safeForMethods) loop(logic func() bool, cb func() bool) {
	lastInd := *m.i
	for *m.i < len(*m.buf) && logic() {
		if !cb() {
			break
		}

		// prevent accidental infinite loop
		if *m.i == lastInd {
			fmt.Println("Warning: Loop Not Incramemted!")
			break
		}
		lastInd = *m.i
	}
}

func (m *safeForMethods) replace(ind *[2]int, rep *[]byte) {
	*m.buf = append((*m.buf)[:(*ind)[0]], append(*rep, (*m.buf)[(*ind)[1]:]...)...)
	*m.i += ((*ind)[0] - (*ind)[1]) + len(*rep)
}

func (m *safeForMethods) getBuf(size int) []byte {
	if *m.i + size < len(*m.buf) {
		return (*m.buf)[*m.i:*m.i + size]
	}
	return []byte{}
}
