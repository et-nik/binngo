package binngo_test

import (
	"encoding/json"
	"testing"

	"github.com/et-nik/binngo"
)

func BenchmarkDecodeList(b *testing.B) {
	binnBinary := []byte{
		0xE0,                          // [type] list (container)
		23,                            // [size] container total size
		0x02,                          // [count] items
		0xA0,                          // [type] = string
		0x05,                          // [size]
		'h', 'e', 'l', 'l', 'o', 0x00, // [data] (null terminated)
		0xA0,                          // [type] = string
		0x05,                          // [size]
		'w', 'o', 'r', 'l', 'd', 0x00, // [data] (null terminated)
	}

	for i := 0; i < b.N; i++ {
		items := []string{}
		err := binngo.Unmarshal(binnBinary, &items)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeListJSON(b *testing.B) {
	jsonData := []byte("[\"hello\", \"world\"]")

	for i := 0; i < b.N; i++ {
		items := []string{}
		err := json.Unmarshal(jsonData, &items)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeMap(b *testing.B) {
	binnBinary := []byte{
		0xE2,									// type = object (container)
		0x11,									// container total size
		0x01,									// key/value pairs count
		0x05, 'h', 'e', 'l', 'l', 'o',			// key
		0xA0,									// type = string
		0x05, 'w', 'o', 'r', 'l', 'd', 0x00,	// value (null terminated)
	}

	for i := 0; i < b.N; i++ {
		m := map[string]string{}
		err := binngo.Unmarshal(binnBinary, &m)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeMapJSON(b *testing.B) {
	jsonData := []byte("{\"hello\": \"world\"}")

	for i := 0; i < b.N; i++ {
		m := map[string]string{}
		err := json.Unmarshal(jsonData, &m)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeObject(b *testing.B) {
	binnBinary := []byte{
		0xE2,									// type = object (container)
		0x11,									// container total size
		0x01,									// key/value pairs count
		0x05, 'h', 'e', 'l', 'l', 'o',			// key
		0xA0,									// type = string
		0x05, 'w', 'o', 'r', 'l', 'd', 0x00,	// value (null terminated)
	}
	type structure struct {
		Hello string `binn:"hello"`
	}

	for i := 0; i < b.N; i++ {
		var v structure
		err := binngo.Unmarshal(binnBinary, &v)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeObjectJSON(b *testing.B) {
	jsonData := []byte("{\"hello\": \"world\"}")
	type structure struct {
		Hello string `binn:"hello"`
	}

	for i := 0; i < b.N; i++ {
		var v structure
		err := json.Unmarshal(jsonData, &v)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeList(b *testing.B) {
	v := []string{"hello", "world"}

	for i := 0; i < b.N; i++ {
		_, err := binngo.Marshal(v)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeListJSON(b *testing.B) {
	v := []string{"hello", "world"}

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(v)
		if err != nil {
			b.Fatal(err)
		}
	}
}
