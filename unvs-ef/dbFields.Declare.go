package unvsef

import "time"

type FieldDateTime Field[time.Time]
type FieldString Field[string]
type FieldInt Field[int]
type FieldFloat Field[float64]
type FieldBool Field[bool]
type FieldTime Field[time.Time]
type FieldBytes Field[[]byte]
type FieldJSON Field[interface{}]
type FieldUUID Field[string]
type FieldBigInt Field[int64]
type FieldDecimal Field[float64]
