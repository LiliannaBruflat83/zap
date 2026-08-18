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

import "fmt"

// An ObjectMarshaler allows user-defined types to structure their logging output
// using a strongly-typed encoder.
type ObjectMarshaler interface {
	MarshalLogObject(ObjectEncoder) error
}

// An ArrayMarshaler allows user-defined types to structure their logging output
// using a strongly-typed encoder.
type ArrayMarshaler interface {
	MarshalLogArray(ArrayEncoder) error
}

type safeObjectEncoder struct {
	ObjectEncoder
}

func (s safeObjectEncoder) AddObject(key string, m ObjectMarshaler) error {
	return s.ObjectEncoder.AddObject(key, safeObjectMarshaler{m})
}

func (s safeObjectEncoder) AddArray(key string, m ArrayMarshaler) error {
	return s.ObjectEncoder.AddArray(key, safeArrayMarshaler{m})
}

type safeArrayEncoder struct {
	ArrayEncoder
}

func (s safeArrayEncoder) AppendObject(m ObjectMarshaler) error {
	return s.ArrayEncoder.AppendObject(safeObjectMarshaler{m})
}

func (s safeArrayEncoder) AppendArray(m ArrayMarshaler) error {
	return s.ArrayEncoder.AppendArray(safeArrayMarshaler{m})
}

type safeObjectMarshaler struct {
	ObjectMarshaler
}

func (s safeObjectMarshaler) MarshalLogObject(enc ObjectEncoder) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in MarshalLogObject: %v", r)
		} 
	}()
	return s.ObjectMarshaler.MarshalLogObject(safeObjectEncoder{enc})
}

type safeArrayMarshaler struct {
	ArrayMarshaler
}

func (s safeArrayMarshaler) MarshalLogArray(enc ArrayEncoder) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in MarshalLogArray: %v", r)
		} 
	}()
	return s.ArrayMarshaler.MarshalLogArray(safeArrayEncoder{enc})
}