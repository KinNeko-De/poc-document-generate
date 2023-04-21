package luatable

type Escaper interface {
	EscapeString(out []byte, in string) ([]byte, error)
}
