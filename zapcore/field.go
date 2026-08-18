// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package zapcore

import (
	"fmt"
	"math"
	"time"
)

// A Field is a marshaler for a single key-value pair. The type-specific
// methods (e.g., Int64) are generally more efficient than the loose interface-
// based logging (e.g., Reflect).
type Field struct {
	Key       string
	Type      FieldType
	Integer   int64
	String    string
	Interface interface{}
}

// AddTo exports the field's data to the Logger's basic, raw context.
func (f Field) AddTo(enc ObjectEncoder) {
	var err error
	switch f.Type {
	case ArrayMarshalerType:
		err = enc.AddArray(f.Key, safeArrayMarshaler{f.Interface.(ArrayMarshaler)})
	case ObjectMarshalerType:
		err = enc.AddObject(f.Key, safeObjectMarshaler{f.Interface.(ObjectMarshaler)})
	case BinaryType:
		enc.AddBinary(f.Key, f.Interface.([]byte))
	case BoolType:
		enc.AddBool(f.Key, f.Integer == 1)
	case ByteStringType:
		enc.AddByteString(f.Key, f.Interface.([]byte))
	case Complex128Type:
		enc.AddComplex128(f.Key, f.Interface.(complex128))
	case Complex64Type:
		enc.AddComplex64(f.Key, f.Interface.(complex64))
	case DurationType:
		enc.AddDuration(f.Key, time.Duration(f.Integer))
	case Float64Type:
		enc.AddFloat64(f.Key, math.Float64frombits(uint64(f.Integer)))
	case Float32Type:
		enc.AddFloat32(f.Key, math.Float32frombits(uint32(f.Integer)))
	case Int64Type:
		enc.AddInt64(f.Key, f.Integer)
	case Int32Type:
		enc.AddInt32(f.Key, int32(f.Integer))
	case Int16Type:
		enc.AddInt16(f.Key, int16(f.Integer))
	case Int8Type:
		enc.AddInt8(f.Key, int8(f.Integer))
	case StringType:
		enc.AddString(f.Key, f.String)
	case TimeType:
		if f.Interface != nil {
			enc.AddTime(f.Key, time.Unix(0, f.Integer).In(f.Interface.(*time.Location)))
		} else {
			enc.AddTime(f.Key, time.Unix(0, f.Integer))
		}
	case TimeFullType:
		enc.AddTime(f.Key, f.Interface.(time.Time))
	case Uint64Type:
		enc.AddUint64(f.Key, uint64(f.Integer))
	case Uint32Type:
		enc.AddUint32(f.Key, uint32(f.Integer))
	case Uint16Type:
		enc.AddUint16(f.Key, uint16(f.Integer))
	case Uint8Type:
		enc.AddUint8(f.Key, uint8(f.Integer))
	case UintptrType:
		enc.AddUintptr(f.Key, uintptr(f.Integer))
	case ReflectType:
		err = enc.AddReflected(f.Key, f.Interface)
	case NamespaceType:
		enc.OpenNamespace(f.Key)
	case StringerType:
		err = encodeStringer(f.Key, f.Interface, enc)
	case ErrorType:
		err = enc.AddObject(f.Key, f.Interface.(errorObject))
	case SkipType:
		break
	default:
		panic(fmt.Sprintf("unknown field type: %v", f.Type))
	}

	if err != nil {
		enc.AddString(fmt.Sprintf("%sError", f.Key), err.Error())
	}
}

// Equals returns true if the fields are equal.
func (f Field) Equals(other Field) bool {
	if f.Type != other.Type {
		return false
	}
	if f.Key != other.Key {
		return false
	}
	if f.Integer != other.Integer {
		return false
	}
	if f.String != other.String {
		return false
	}
	return f.Interface == other.Interface
}

func encodeStringer(key string, val interface{}, enc ObjectEncoder) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("panic in Stringer method: %v", v)
		}
	}()
	enc.AddString(key, val.(fmt.Stringer).String())
	return err
}