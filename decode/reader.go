package decode

import (
	"fmt"
	"io"

	"github.com/et-nik/binngo/binn"
	"github.com/et-nik/binngo/encode"
)

func readValue(btype binn.Type, reader io.Reader) ([]byte, error) {
	tp := btype &^ binn.StorageTypeMask

	var readingSize int
	var containerSize int

	var bytes []byte

	switch tp {
	case binn.StorageNoBytes:
		readingSize = 0
	case binn.StorageByte:
		readingSize = 1
	case binn.StorageWord:
		readingSize = 2
	case binn.StorageDWord:
		readingSize = 4
	case binn.StorageQWord:
		readingSize = 8
	case binn.StorageString:
		dataSize, _, err := readSize(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read string storage size: %w", err)
		}
		readingSize = dataSize + 1 // data size and null terminator
	case binn.StorageBlob, binn.StorageContainer:
		s, l, err := readSize(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read storage size: %w", err)
		}

		containerSize = s

		bytes = append(bytes, encode.Int(s)...)

		readingSize = containerSize - 1 - int(l) // minus container type byte and size byte
	default:
		return nil, ErrUnknownType
	}

	if readingSize == 0 {
		return []byte{byte(btype)}, nil
	}

	b := make([]byte, readingSize)

	_, err := reader.Read(b)
	bytes = append(bytes, b...)

	if err != nil {
		return nil, fmt.Errorf("failed to read storage: %w", err)
	}

	return bytes, nil
}

func readType(reader io.Reader) (binn.Type, readLen, error) {
	var bt = make([]byte, 1)

	_, err := reader.Read(bt)
	if err != nil {
		return binn.Null, 0, ErrFailedToReadType
	}

	return Type(bt), 1, nil
}

func readSize(reader io.Reader) (int, readLen, error) {
	var bsz = make([]byte, 1)
	_, err := reader.Read(bsz)
	if err != nil {
		return 0, 0, ErrFailedToReadSize
	}

	read := 1

	sz := int(Uint8(bsz))

	if sz > maxOneByteSize {
		var bszOtherBytes = make([]byte, 3)
		_, err := reader.Read(bszOtherBytes)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to read long size: %w", err)
		}
		read += 3
		sz ^= 0x80000000

		sz = int(Uint32([]byte{
			byte(sz),
			bszOtherBytes[0],
			bszOtherBytes[1],
			bszOtherBytes[2],
		}))
	}

	return sz, readLen(read), nil
}
