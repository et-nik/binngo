package decode

import "io"

func decodeList(reader io.Reader, v interface{}) error {
	sz, rSize, err := readSize(reader)
	if err != nil {
		return err
	}

	cnt, rCount, err := readSize(reader)
	if err != nil {
		return err
	}

	return decodeListItems(reader, v, sz, rSize+rCount, cnt)
}

func decodeListItems(reader io.Reader, v interface{}, size int, wasRead readLen, items int) error {
	rItems := 0
	readPosition := wasRead

	for rItems < items && readPosition < readLen(size) {
		btype, rlen, err := readType(reader)
		if err != nil {
			return err
		}

		bval, err := readValue(btype, reader)
		if err != nil {
			return err
		}

		err = addSliceItem(btype, bval, v)
		if err != nil {
			return err
		}

		rItems++
		readPosition += readLen(len(bval)) + rlen
	}

	return nil
}
