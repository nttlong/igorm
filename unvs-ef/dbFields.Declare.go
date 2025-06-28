package unvsef

import "time"

type FieldFloat Field[float64]
type FieldBool Field[bool]
type FieldTime Field[time.Time]
type FieldBytes Field[[]byte]
type FieldJSON Field[interface{}]
type FieldUUID Field[string]
type FieldBigInt Field[int64]
type FieldDecimal Field[float64]

type FieldUint32 Field[uint32]
type FieldUint16 Field[uint16]
type FieldUint8 Field[uint8]
type FieldUint Field[uint]
type FieldInt64 Field[int64]
type FieldInt32 Field[int32]
type FieldInt16 Field[int16]
type FieldInt8 Field[int8]
type FieldInt Field[int]
type FieldAny Field[any]
