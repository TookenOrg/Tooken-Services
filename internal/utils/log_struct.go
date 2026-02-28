package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func Dump(v any) string {
	return dumpValue(reflect.ValueOf(v), 0)
}

func dumpValue(v reflect.Value, indent int) string {
	if !v.IsValid() {
		return "nil"
	}

	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return "nil"
		}
		v = v.Elem()
	}

	prefix := strings.Repeat("  ", indent)

	switch v.Kind() {

	case reflect.Struct:
		var sb strings.Builder
		sb.WriteString("{\n")

		for i := 0; i < v.NumField(); i++ {
			fieldType := v.Type().Field(i)
			fieldValue := v.Field(i)

			if fieldType.PkgPath != "" {
				continue
			}

			sb.WriteString(fmt.Sprintf(
				"%s  %s: %s\n",
				prefix,
				fieldType.Name,
				dumpValue(fieldValue, indent+1),
			))
		}

		sb.WriteString(prefix + "}")
		return sb.String()

	case reflect.Slice, reflect.Array:
		var sb strings.Builder
		sb.WriteString("[\n")
		for i := 0; i < v.Len(); i++ {
			sb.WriteString(fmt.Sprintf(
				"%s  - %s\n",
				prefix,
				dumpValue(v.Index(i), indent+1),
			))
		}
		sb.WriteString(prefix + "]")
		return sb.String()

	case reflect.Map:
		var sb strings.Builder
		sb.WriteString("{\n")
		for _, key := range v.MapKeys() {
			sb.WriteString(fmt.Sprintf(
				"%s  %v: %s\n",
				prefix,
				key.Interface(),
				dumpValue(v.MapIndex(key), indent+1),
			))
		}
		sb.WriteString(prefix + "}")
		return sb.String()

	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
