package lufx

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func (e *Executor) register(values map[string]any) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if len(values) < 1 {
		return errors.New("function [register] needs at least one value")
	}

	for name, value := range values {
		if err := e.go2lufx(name, reflect.ValueOf(value)); err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) go2lufx(name string, value reflect.Value) error {
	_, exists := e.commands[name]
	if exists {
		return errors.New("value [" + name + "] is already registered")
	}

	switch value.Kind() {
	case reflect.Pointer:
		return e.go2lufx(name, value.Elem())
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			fieldName := field.Name
			fieldValue := value.Field(i).Interface()

			tag := field.Tag.Get("lufx")
			if len(tag) > 0 {
				fieldName = tag
			}

			fieldName = fmt.Sprintf("%s->$%s", name, fieldName)
			e.commands[fieldName] = func(session io.Writer, args []string) (int, error) {
				_, err := io.WriteString(session, fmt.Sprintf("%v", fieldValue))
				if err != nil {
					return 0, err
				}
				return len(args), nil
			}
		}
		return nil
	case reflect.Map:
		keys := value.MapKeys()
		for _, key := range keys {
			e.commands[fmt.Sprintf("$%s->$%s", name, fmt.Sprintf("%v", key.Interface()))] = func(session io.Writer, args []string) (int, error) {
				_, err := io.WriteString(session, fmt.Sprintf("%v", value.MapIndex(key).Interface()))
				if err != nil {
					return 0, err
				}

				return len(args), nil
			}
		}
		return nil
	case reflect.Func:
		name = fmt.Sprintf("%s()", name)
		if value.Type().ConvertibleTo(reflect.TypeOf(FxFunction(nil))) {
			e.commands[name] = value.Interface().(FxFunction)
		} else {
			return errors.New(fmt.Sprintf("unsupported type: [%s]", strings.ToUpper(value.Kind().String())))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float64, reflect.Float32, reflect.String, reflect.Bool:
		e.commands[fmt.Sprintf("$%s", name)] = func(session io.Writer, args []string) (int, error) {
			_, err := io.WriteString(session, fmt.Sprintf("%v", value.Interface()))
			if err != nil {
				return 0, err
			}
			return len(args), nil
		}
	default:
		return errors.New(fmt.Sprintf("unsupported type: [%s]", strings.ToUpper(value.Kind().String())))
	}
	return nil
}
