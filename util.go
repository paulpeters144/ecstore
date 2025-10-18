package ecstore

import "reflect"

func getTypeKey(e Entity) (string, error) {
	if e == nil {
		return "", ErrInvalidEntityPointer
	}

	val := reflect.ValueOf(e)

	if val.Kind() == reflect.Pointer && val.IsNil() {
		return "", ErrInvalidEntityPointer
	}
	if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Struct {
		return "", ErrInvalidEntityPointer
	}
	return reflect.TypeOf(e).Elem().Name(), nil
}
