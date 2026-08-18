package zapcore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type panickingObject struct{}

func (panickingObject) MarshalLogObject(enc ObjectEncoder) error {
	panic("simulated out-of-bounds or nil pointer")
}

type panickingArray struct{}

func (panickingArray) MarshalLogArray(enc ArrayEncoder) error {
	panic("simulated out-of-bounds or nil pointer")
}

type nestedPanickingObject struct{}

func (nestedPanickingObject) MarshalLogObject(enc ObjectEncoder) error {
	panic("nested panic")
}

type outerObject struct{}

func (outerObject) MarshalLogObject(enc ObjectEncoder) error {
	return enc.AddObject("nested", nestedPanickingObject{})
}

func TestFieldPanicRecovery(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		enc := &MapObjectEncoder{Fields: make(map[string]interface{})}
		f := Field{
			Key:       "obj",
			Type:      ObjectMarshalerType,
			Interface: panickingObject{},
		}
		f.AddTo(enc)
		assert.Equal(t, "panic in MarshalLogObject: simulated out-of-bounds or nil pointer", enc.Fields["objError"])
	})

	t.Run("array", func(t *testing.T) {
		enc := &MapObjectEncoder{Fields: make(map[string]interface{})}
		f := Field{
			Key:       "arr",
			Type:      ArrayMarshalerType,
			Interface: panickingArray{},
		}
		f.AddTo(enc)
		assert.Equal(t, "panic in MarshalLogArray: simulated out-of-bounds or nil pointer", enc.Fields["arrError"])
	})

	t.Run("nested object", func(t *testing.T) {
		enc := &MapObjectEncoder{Fields: make(map[string]interface{})}
		f := Field{
			Key:       "obj",
			Type:      ObjectMarshalerType,
			Interface: outerObject{},
		}
		f.AddTo(enc)
		assert.Equal(t, "panic in MarshalLogObject: nested panic", enc.Fields["objError"])
	})
}