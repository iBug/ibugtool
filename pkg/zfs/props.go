package zfs

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Property struct {
	Name     string
	Type     string
	Creation time.Time
}

func propertyName(f reflect.StructField) string {
	name := strings.ToLower(f.Name)
	if tag, ok := f.Tag.Lookup("zfs"); ok {
		parts := strings.Split(tag, ",")
		if parts[0] != "" {
			name = parts[0]
		}
	}
	if name == "_" {
		name = "-"
	}
	return name
}

func buildPropertiesArgument(p any) string {
	t := reflect.TypeOf(p)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		qualName := t.PkgPath() + "." + t.Name()
		panic(fmt.Errorf("internal error: %q is not a struct", qualName))
	}

	props := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		name := propertyName(t.Field(i))
		if name != "-" {
			props = append(props, name)
		}
	}

	return strings.Join(props, ",")
}

func parsePropertyValue(v reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		if v.OverflowInt(i) {
			return fmt.Errorf("integer overflow for %s: %d", v.Type().String(), i)
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		if v.OverflowUint(i) {
			return fmt.Errorf("unsigned integer overflow for %s: %d", v.Type().String(), i)
		}
		v.SetUint(i)
	default:
		switch v.Interface().(type) {
		case time.Time:
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(time.Unix(i, 0)))
		default:
			return errors.New("unknown type")
		}
	}
	return nil
}

func parseProperties(p any, s string) error {
	fields := strings.Split(strings.TrimSpace(s), "\t")

	v := reflect.ValueOf(p)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		qualName := v.Type().PkgPath() + "." + v.Type().Name()
		return fmt.Errorf("internal error: %q is not a struct", qualName)
	}

	for i, nf := 0, 0; i < v.NumField(); i++ {
		if propertyName(v.Type().Field(i)) == "-" {
			continue
		}
		if err := parsePropertyValue(v.Field(i), fields[nf]); err != nil {
			return err
		}
		nf++
	}
	return nil
}

func GetProperties(p any, dataset string) (err error) {
	props := buildPropertiesArgument(p)

	cmd := ZfsCommand("list", "-Hp", "-o", props, dataset)
	r, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	if err = cmd.Start(); err != nil {
		return
	}
	defer cmd.Wait()

	b, err := io.ReadAll(r)
	if err != nil {
		return
	}
	s := strings.TrimSpace(string(b))
	return parseProperties(p, s)
}
