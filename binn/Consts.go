package binn

const (
	StorageNoBytes   = 0x00
	StorageByte      = 0x20
	StorageWord      = 0x40
	StorageDWord     = 0x60
	StorageQWord     = 0x80
	StorageString    = 0xA0
	StorageBlob      = 0xC0
	StorageContainer = 0xE0
	StorageVirtual   = 0x80000
)

const (
	StorageMin = StorageNoBytes
	StorageMax = StorageContainer
)

const (
	StorageMask       = 0xE0
	StorageMask16     = 0xE000
	StorageHasMore    = 0x10
	StorageTypeMask   = 0x0F
	StorageTypeMask16 = 0x0FFF
)

const (
	ListType   = 0xE0
	MapType    = 0xE1
	ObjectType = 0xE2

	Null  = 0x00
	True  = 0x01
	False = 0x02

	Uint8Type  = 0x20
	Int8Type   = 0x21
	Uint16Type = 0x40
	Int16Type  = 0x41
	Uint32Type = 0x60
	Int32Type  = 0x61
	Uint64Type = 0x80
	Int64Type  = 0x81

	Schar = Int8Type
	Uchar = Uint8Type

	StringType      = 0xA0
	DateTimeType    = 0xA1
	DateType        = 0xA2
	TimeType        = 0xA3
	DecimalType     = 0xA4
	CurrencyStrType = 0xA5
	SingleStrType   = 0xA6
	DoubleStrType   = 0xA7

	Float32Type = 0x62
	Float64Type = 0x82
	FloatType   = Float32Type
	SingleType  = Float32Type
	DoubleType  = Float64Type

	CurrencyType = 0x83
	BlobType     = 0xC0
)

const (
	HTML       = 0xB001
	XML        = 0xB002
	JSON       = 0xB003
	JavaScript = 0xB004
	CSS        = 0xB005
)

const (
	JPEG = 0xD001
	GIF  = 0xD002
	PNG  = 0xD003
	BMP  = 0xD004
)

type Type int
