// Package encoder implements BINN encoding.
package encode

func Marshal(v interface{}) ([]byte, error) {
	return marshal(v)
}
