package telegram

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"
)

type Fileable interface{}

type FileUploader interface {
	Name() string
	Reader() (io.Reader, error)
	Size() int64
}

var fileUploaderInterface = reflect.TypeOf((*FileUploader)(nil)).Elem()

type fileUpload struct {
	params    map[string]string
	fieldname string
	file      FileUploader
	err       error
}

func isFileUpload(req interface{}) (upload fileUpload, isUpload bool) {
	val := reflect.ValueOf(req)
	if val.Kind() != reflect.Ptr {
		return fileUpload{}, false
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return fileUpload{}, false
	}

	var uploader reflect.Value
	var uploaderField reflect.StructField
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if f.Kind() == reflect.Interface && !f.IsZero() && f.Elem().Type().Implements(fileUploaderInterface) {
			uploader = f
			uploaderField = val.Type().Field(i)
			isUpload = true
			break
		}
	}

	if !isUpload {
		return fileUpload{}, false
	}

	upload.params = make(map[string]string)
	upload.fieldname, _ = parseTag(uploaderField.Tag.Get("json"))
	upload.file = uploader.Interface().(FileUploader)

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if f == uploader {
			continue
		}

		fieldType := val.Type().Field(i)
		fieldName, tagOpts := parseTag(fieldType.Tag.Get("json"))

		if tagOpts.Contains("omitempty") && isEmptyValue(f) {
			continue
		}

		kind := f.Kind()
		if kind == reflect.String {
			upload.params[fieldName] = f.String()
		} else {
			data, err := json.Marshal(f.Interface())
			if err != nil {
				upload.err = err
				return
			}

			upload.params[fieldName] = string(data)
		}
	}

	return upload, true
}

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
