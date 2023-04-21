package luatable

import (
	"errors"
	"math"
	"math/bits"
	"strconv"
	"strings"
	"unicode/utf8"
)

// kind represents an encoding type.
type kind uint8

const KeyAssign = "="
const NullValue = "nil"
const BoolTrue = "true"
const BoolFalse = "false"
const BeginString = "\""
const EndString = "\""
const ArrayOpen = "{"
const ArrayClose = "}"
const TableOpen = "{"
const TableClose = "}"

const (
	_ kind = (1 << iota) / 2
	key
	scalar
	objectOpen
	objectClose
	arrayOpen
	arrayClose
)

// Encoder provides methods to write out lua table constructs and values. The user is
// responsible for producing valid sequences of JSON constructs and values.
type Encoder struct {
	indent   string
	lastKind kind
	indents  []byte
	out      []byte
}

// NewEncoder returns an Encoder.
//
// If indent is a non-empty string, it causes every value for a table
// to be preceded by the indent and trailed by a newline.
func NewEncoder(indent string) (*Encoder, error) {
	err := CheckForInvalidIndentChars(indent)
	if err != nil {
		return nil, err
	}

	e := CreateEncoder(indent)
	return e, nil
}

func CheckForInvalidIndentChars(indent string) error {
	if strings.Trim(indent, " \t") != "" {
		return errors.New("indent must be space or tab characters")
	}
	return nil
}

func CreateEncoder(indent string) *Encoder {
	e := &Encoder{}
	if len(indent) > 0 {
		e.indent = indent
	}
	return e
}

// Bytes returns the content of the written bytes.
func (e *Encoder) Bytes() []byte {
	return e.out
}

// WriteNull writes out the null value.
func (e *Encoder) WriteNull() {
	e.prepareNext(scalar)
	e.out = append(e.out, NullValue...)
}

// WriteBool writes out the given boolean value.
func (e *Encoder) WriteBool(b bool) {
	e.prepareNext(scalar)
	if b {
		e.out = append(e.out, BoolTrue...)
	} else {
		e.out = append(e.out, BoolFalse...)
	}
}

// WriteString escaped the string according to rules needed for luatex
func (e *Encoder) WriteString(s string) error {
	e.prepareNext(scalar)
	e.StartString()
	defer e.EndString()
	var err error
	if e.out, err = appendString(e.out, s); err != nil {
		return err
	}
	return nil
}

func appendString(out []byte, in string) ([]byte, error) {
	i := indexNeedEscapeInString(in)
	in, out = in[i:], append(out, in[:i]...)
	for len(in) > 0 {
		switch r, n := utf8.DecodeRuneInString(in); {
		case r == utf8.RuneError && n == 1:
			return out, errors.New("the string contains invalid UTF-8 characters")
		case r < ' ' || r == '"' || r == '\\' || r == '%':
			switch r {
			case '%':
				out = append(out, "\\\\"...)
				out = append(out, byte(r))
			case '\\':
				nextChar := in[1]
				switch nextChar {
				case 'n':
					out = append(out, "\\\\\\\\"...)
					n++
				case 'r':
					n++
				default:
					out = append(out, "\\\\"...)
				}
			case '"':
				out = append(out, '\\')
				out = append(out, byte(r))
			case '\b':
				errors.New("not implemented yet")
			case '\f':
				errors.New("not implemented yet")
			case '\n':
				out = append(out, "\\\\\\\\"...)
			case '\r':
				// do nothing as \r\n and \n are reduced to line break
			case '\t':
				errors.New("not implemented yet")
			default:
				out = append(out, 'u')
				out = append(out, "0000"[1+(bits.Len32(uint32(r))-1)/4:]...)
				out = strconv.AppendUint(out, uint64(r), 16)
			}
			in = in[n:]
		default:
			i := indexNeedEscapeInString(in[n:])
			in, out = in[n+i:], append(out, in[:n+i]...)
		}
	}
	return out, nil
}

// indexNeedEscapeInString returns the index of the character that needs
// escaping. If no characters need escaping, this returns the input length.
func indexNeedEscapeInString(s string) int {
	for i, r := range s {
		if r < ' ' || r == '\\' || r == '"' || r == '%' || r == utf8.RuneError {
			return i
		}
	}
	return len(s)
}

// WriteFloat writes out the given float and bitSize in JSON number value.
func (e *Encoder) WriteFloat(n float64, bitSize int) {
	e.prepareNext(scalar)
	e.out = appendFloat(e.out, n, bitSize)
}

// appendFloat formats given float in bitSize, and appends to the given []byte.
func appendFloat(out []byte, n float64, bitSize int) []byte {
	switch {
	case math.IsNaN(n):
		return append(out, `"NaN"`...)
	case math.IsInf(n, +1):
		return append(out, `"Infinity"`...)
	case math.IsInf(n, -1):
		return append(out, `"-Infinity"`...)
	}

	// JSON number formatting logic based on encoding/json.
	// See floatEncoder.encode for reference.
	fmt := byte('f')
	if abs := math.Abs(n); abs != 0 {
		if bitSize == 64 && (abs < 1e-6 || abs >= 1e21) ||
			bitSize == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}
	out = strconv.AppendFloat(out, n, fmt, -1, bitSize)
	if fmt == 'e' {
		n := len(out)
		if n >= 4 && out[n-4] == 'e' && out[n-3] == '-' && out[n-2] == '0' {
			out[n-2] = out[n-1]
			out = out[:n-1]
		}
	}
	return out
}

func (e *Encoder) WriteInt(n int64) {
	e.prepareNext(scalar)
	e.out = append(e.out, strconv.FormatInt(n, 10)...)
}

func (e *Encoder) WriteUint(n uint64) {
	e.prepareNext(scalar)
	e.out = append(e.out, strconv.FormatUint(n, 10)...)
}

func (e *Encoder) WriteNumber(number string) {
	e.prepareNext(scalar)
	e.out = append(e.out, number...)
}

func (e *Encoder) StartObject() {
	e.prepareNext(objectOpen)
	e.out = append(e.out, TableOpen...)
}

func (e *Encoder) EndObject() {
	e.prepareNext(objectClose)
	e.out = append(e.out, TableClose...)
}

func (e *Encoder) WriteKey(s string) error {
	e.prepareNext(key)
	e.out = append(e.out, s...)
	e.WriteKeyAssign()
	return nil
}

func (e *Encoder) StartArray() {
	e.prepareNext(arrayOpen)
	e.out = append(e.out, ArrayOpen...)
}

func (e *Encoder) EndArray() {
	e.prepareNext(arrayClose)
	e.out = append(e.out, ArrayClose...)
}

func (e *Encoder) StartString() {
	e.out = append(e.out, BeginString...)
}

func (e *Encoder) EndString() {
	e.out = append(e.out, EndString...)
}

func (e *Encoder) WriteIndexedList(i int) {
	e.prepareNext(key)
	e.out = append(e.out, "["...)
	e.out = append(e.out, strconv.FormatInt(int64(i), 10)...)
	e.out = append(e.out, "]"...)
	e.WriteKeyAssign()
}

func (e *Encoder) WriteKeyAssign() {
	if len(e.indent) != 0 {
		e.out = append(e.out, " "...)
	}
	e.out = append(e.out, KeyAssign...)
	if len(e.indent) != 0 {
		e.out = append(e.out, " "...)
	}
}

// prepareNext adds possible comma and indentation for the next value based
// on last type and indent option. It also updates lastKind to next.
func (e *Encoder) prepareNext(next kind) {
	defer func() {
		// Set lastKind to next.
		e.lastKind = next
	}()

	if len(e.indent) == 0 {
		// Need to add comma on the following condition.
		if e.lastKind&(scalar|objectClose|arrayClose) != 0 &&
			next&(key|scalar|objectOpen|arrayOpen) != 0 {
			e.out = append(e.out, ',')
		}
		return
	}

	switch {
	case e.lastKind&(objectOpen|arrayOpen) != 0:
		// If next type is NOT closing, add indent and newline.
		if next&(objectClose|arrayClose) == 0 {
			e.indents = append(e.indents, e.indent...)
			e.out = append(e.out, '\n')
			e.out = append(e.out, e.indents...)
		}

	case e.lastKind&(scalar|objectClose|arrayClose) != 0:
		switch {
		// If next type is either a value or name, add comma and newline.
		case next&(key|scalar|objectOpen|arrayOpen) != 0:
			e.out = append(e.out, ',', '\n')

		// If next type is a closing object or array, adjust indentation.
		case next&(objectClose|arrayClose) != 0:
			e.indents = e.indents[:len(e.indents)-len(e.indent)]
			e.out = append(e.out, '\n')
		}
		e.out = append(e.out, e.indents...)
	}

}
