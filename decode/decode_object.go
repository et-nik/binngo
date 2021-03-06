package decode

import "io"

func decodeObject(reader io.Reader, v interface{}) error {
	sz, rSize, err := readSize(reader)
	if err != nil {
		return err
	}

	cnt, rCount, err := readSize(reader)
	if err != nil {
		return err
	}

	return decodeObjectItems(reader, v, sz, rSize+rCount, cnt)
}

func decodeObjectItems(reader io.Reader, v interface{}, size int, wasRead readLen, items int) error {
	rItems := 0
	rPosition := wasRead

	for rItems < items && rPosition < readLen(size) {
		key, read, err := readObjectKey(reader)
		if err != nil {
			return err
		}
		rPosition += read

		btype, read, err := readType(reader)
		if err != nil {
			return err
		}
		rPosition += read

		bval, err := readValue(btype, reader)
		if err != nil {
			return err
		}
		rPosition += readLen(len(bval))

		err = addObjectItem(key, btype, bval, v)
		if err != nil {
			return err
		}

		rItems++
	}

	return nil
}

func readObjectKey(reader io.Reader) (string, readLen, error) {
	var bsz = make([]byte, 1)
	n, err := reader.Read(bsz)
	if err != nil {
		return "", 0, err
	}

	if n != 1 {
		return "", readLen(n), ErrIncompleteRead
	}

	sz := int(Uint8(bsz))

	var bkey = make([]byte, sz)

	n, err = reader.Read(bkey)
	if err != nil {
		return "", 0, err
	}

	if n != sz {
		return "", readLen(n), ErrIncompleteRead
	}

	return String(bkey), readLen(sz + 1), nil
}
