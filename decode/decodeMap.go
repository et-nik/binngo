package decode

import "io"

func decodeMap(reader io.Reader, v interface{}) error {
	sz, rSize, err := readSize(reader)
	if err != nil {
		return err
	}

	cnt, rCount, err := readSize(reader)
	if err != nil {
		return err
	}

	return decodeMapItems(reader, v, sz, rSize+rCount, cnt)
}

func decodeMapItems(reader io.Reader, v interface{}, size int, wasRead readLen, items int) error {
	readItems := 0
	readPosition := wasRead

	for readItems < items && readPosition < readLen(size) {
		key, read, err := readMapKey(reader)
		if err != nil {
			return err
		}
		readPosition += read

		t, read, err := readType(reader)
		if err != nil {
			return err
		}
		readPosition += read

		val, err := readValue(t, reader)
		if err != nil {
			return err
		}
		readPosition += readLen(len(val))

		err = addMapItem(key, t, val, v)
		if err != nil {
			return err
		}

		readItems++
	}

	return nil
}

func readMapKey(reader io.Reader) (int, readLen, error) {
	var bk = make([]byte, 4)
	_, err := reader.Read(bk)
	if err != nil {
		return 0, 0, err
	}

	return int(DecodeInt32(bk)), 4, nil
}
