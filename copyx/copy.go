package copyx

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

var (
	// ErrInvalidCopyDestination is returned when the destination is not addressable.
	ErrInvalidCopyDestination = errors.New("copy: destination must be non-nil and addressable")
	// ErrInvalidCopyFrom is returned when the source is not addressable.
	ErrInvalidCopyFrom = errors.New("copy: from must be non-nil and addressable")
	// ErrMapKeyNotMatch is returned when the map key types do not match.
	ErrMapKeyNotMatch = errors.New("copy: map's key type doesn't match")
	// ErrNotSupported is returned when the type is not supported.
	ErrNotSupported = errors.New("copy: not supported")
	// ErrFieldNameTagStartNotUpperCase is returned when the field name tag does not start with an uppercase letter.
	ErrFieldNameTagStartNotUpperCase = errors.New("copy: field name tag must be start upper case")
)

const (
	// Denotes that a destination field must be copied to. If copying fails then a panic will ensue.
	tagMust uint8 = 1 << iota
	// Indicates that a panic should not occur if a field cannot be copied.
	tagNoPanic
	// Ignore the destionaion field.
	tagIgnore
	// Override the field value if not default value.
	tagOverride
	// Has been copied.
	hasCopied

	// Some default converter types for a nicer syntax.
	String  string  = ""
	Bool    bool    = false
	Int     int     = 0
	Float32 float32 = 0
	Float64 float64 = 0
)

// Opts are the options for the copier.
type Opts struct {
	// IgnoreEmpty will ignore empty fields.
	IgnoreEmpty bool
	// CaseSensitive will ignore case sensitivity.
	CaseSensitive bool
	// DeepCopy will copy the fields deeply.
	DeepCopy bool
	// Converters are the type converters.
	Converters []TypeConverter
	// FieldNameMapping are the field name mappings.
	FieldNameMapping []FieldNameMapping
}

// WithIgnoreEmpty will ignore empty fields.
func WithIgnoreEmpty() Opt {
	return func(o *Opts) {
		o.IgnoreEmpty = true
	}
}

// WithCaseSensitive will ignore case sensitivity.
func WithCaseSensitive() Opt {
	return func(o *Opts) {
		o.CaseSensitive = true
	}
}

// WithDeepCopy will copy the fields deeply.
func WithDeepCopy() Opt {
	return func(o *Opts) {
		o.DeepCopy = true
	}
}

// WithConverters are the type converters.
func WithConverters(converters ...TypeConverter) Opt {
	return func(o *Opts) {
		o.Converters = converters
	}
}

// WithFieldNameMapping are the field name mappings.
func WithFieldNameMapping(fieldNameMapping ...FieldNameMapping) Opt {
	return func(o *Opts) {
		o.FieldNameMapping = fieldNameMapping
	}
}

// Opt is a function that configures the copy.
type Opt func(*Opts)

// DefaultOpts are the default options for the copy.
func DefaultOpts() *Opts {
	return &Opts{}
}

// Configure the copy.
func (o *Opts) Configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *Opts) converters() map[converterPair]TypeConverter {
	converters := map[converterPair]TypeConverter{}

	// save converters into map for faster lookup
	for i := range o.Converters {
		pair := converterPair{
			SrcType: reflect.TypeOf(o.Converters[i].SrcType),
			DstType: reflect.TypeOf(o.Converters[i].DstType),
		}

		converters[pair] = o.Converters[i]
	}

	return converters
}

type TypeConverter struct {
	SrcType interface{}
	DstType interface{}
	Fn      func(src interface{}) (dst interface{}, err error)
}

type converterPair struct {
	SrcType reflect.Type
	DstType reflect.Type
}

func (o *Opts) fieldNameMapping() map[converterPair]FieldNameMapping {
	mapping := map[converterPair]FieldNameMapping{}

	for i := range o.FieldNameMapping {
		pair := converterPair{
			SrcType: reflect.TypeOf(o.FieldNameMapping[i].SrcType),
			DstType: reflect.TypeOf(o.FieldNameMapping[i].DstType),
		}

		mapping[pair] = o.FieldNameMapping[i]
	}

	return mapping
}

// FieldNameMapping is the field name mapping.
type FieldNameMapping struct {
	SrcType interface{}
	DstType interface{}
	Mapping map[string]string
}

// Tag Flags.
type flags struct {
	BitFlags  map[string]uint8
	SrcNames  tagNameMapping
	DestNames tagNameMapping
}

// Field Tag name mapping.
type tagNameMapping struct {
	FieldNameToTag map[string]string
	TagToFieldName map[string]string
}

// Copy is copying things.
func Copy(toValue interface{}, fromValue interface{}) (err error) {
	return copier(toValue, fromValue, DefaultOpts())
}

// CopyWithOption copy with option.
func CopyWithOption(toValue interface{}, fromValue interface{}, opts ...Opt) (err error) {
	options := DefaultOpts()
	options.Configure(opts...)

	return copier(toValue, fromValue, options)
}

// nolint:gocyclo
func copier(toValue interface{}, fromValue interface{}, opt *Opts) (err error) {
	var (
		isSlice    bool
		amount     = 1
		from       = indirect(reflect.ValueOf(fromValue))
		to         = indirect(reflect.ValueOf(toValue))
		converters = opt.converters()
		mappings   = opt.fieldNameMapping()
	)

	if !to.CanAddr() {
		return ErrInvalidCopyDestination
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return ErrInvalidCopyFrom
	}

	fromType, isPtrFrom := indirectType(from.Type())
	toType, _ := indirectType(to.Type())

	if fromType.Kind() == reflect.Interface {
		fromType = reflect.TypeOf(from.Interface())
	}

	if toType.Kind() == reflect.Interface {
		toType, _ = indirectType(reflect.TypeOf(to.Interface()))
		oldTo := to
		to = reflect.New(reflect.TypeOf(to.Interface())).Elem()
		defer func() {
			oldTo.Set(to)
		}()
	}

	// Just set it if possible to assign for normal types
	if from.Kind() != reflect.Slice && from.Kind() != reflect.Struct && from.Kind() != reflect.Map && (from.Type().AssignableTo(to.Type()) || from.Type().ConvertibleTo(to.Type())) {
		if !isPtrFrom || !opt.DeepCopy {
			to.Set(from.Convert(to.Type()))
		} else {
			fromCopy := reflect.New(from.Type())
			fromCopy.Set(from.Elem())
			to.Set(fromCopy.Convert(to.Type()))
		}

		return err
	}

	if from.Kind() != reflect.Slice && fromType.Kind() == reflect.Map && toType.Kind() == reflect.Map {
		if !fromType.Key().ConvertibleTo(toType.Key()) {
			return ErrMapKeyNotMatch
		}

		if to.IsNil() {
			to.Set(reflect.MakeMapWithSize(toType, from.Len()))
		}

		for _, k := range from.MapKeys() {
			toKey := indirect(reflect.New(toType.Key()))
			isSet, err := set(toKey, k, opt.DeepCopy, converters)
			if err != nil {
				return err
			}
			if !isSet {
				return fmt.Errorf("%w map, old key: %v, new key: %v", ErrNotSupported, k.Type(), toType.Key())
			}

			elemType := toType.Elem()
			if elemType.Kind() != reflect.Slice {
				elemType, _ = indirectType(elemType)
			}
			toValue := indirect(reflect.New(elemType))
			isSet, err = set(toValue, from.MapIndex(k), opt.DeepCopy, converters)
			if err != nil {
				return err
			}
			if !isSet {
				if err = copier(toValue.Addr().Interface(), from.MapIndex(k).Interface(), opt); err != nil {
					return err
				}
			}

			for {
				if elemType == toType.Elem() {
					to.SetMapIndex(toKey, toValue)
					break
				}
				elemType = reflect.PointerTo(elemType)
				toValue = toValue.Addr()
			}
		}
		return err
	}

	if from.Kind() == reflect.Slice && to.Kind() == reflect.Slice {
		// Return directly if both slices are nil
		if from.IsNil() && to.IsNil() {
			return err
		}
		if to.IsNil() {
			slice := reflect.MakeSlice(reflect.SliceOf(to.Type().Elem()), from.Len(), from.Cap())
			to.Set(slice)
		}
		if fromType.ConvertibleTo(toType) {
			for i := 0; i < from.Len(); i++ {
				if to.Len() < i+1 {
					to.Set(reflect.Append(to, reflect.New(to.Type().Elem()).Elem()))
				}
				isSet, err := set(to.Index(i), from.Index(i), opt.DeepCopy, converters)
				if err != nil {
					return err
				}
				if !isSet {
					// ignore error while copy slice element
					err = copier(to.Index(i).Addr().Interface(), from.Index(i).Interface(), opt)
					if err != nil {
						continue
					}
				}
			}

			if to.Len() > from.Len() {
				to.SetLen(from.Len())
			}

			return err
		}
	}

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		// skip not supported type
		return err
	}

	if len(converters) > 0 {
		if ok, e := set(to, from, opt.DeepCopy, converters); e == nil && ok {
			// converter supported
			return err
		}
	}

	if from.Kind() == reflect.Slice || to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}
			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		if len(converters) > 0 {
			if ok, e := set(dest, source, opt.DeepCopy, converters); e == nil && ok {
				if isSlice {
					if to.Type().Elem().Kind() == reflect.Ptr {
						to.Index(i).Set(dest.Addr())
					} else {
						if to.Len() < i+1 {
							reflect.Append(to, dest)
						} else {
							to.Index(i).Set(dest)
						}
					}
				} else {
					to.Set(dest)
				}

				continue
			}
		}

		destKind := dest.Kind()
		initDest := false
		if destKind == reflect.Interface {
			initDest = true
			dest = indirect(reflect.New(toType))
		}

		// Get tag options
		flgs, err := getFlags(dest, source, toType, fromType)
		if err != nil {
			return err
		}

		// check source
		if source.IsValid() {
			copyUnexportedStructFields(dest, source)

			// Copy from source field to dest field or method
			fromTypeFields := deepFields(fromType)
			for _, field := range fromTypeFields {
				name := field.Name

				// Get bit flags for field
				fieldFlags := flgs.BitFlags[name]

				// Check if we should ignore copying
				if (fieldFlags & tagIgnore) != 0 {
					continue
				}

				fieldNamesMapping := getFieldNamesMapping(mappings, fromType, toType)

				srcFieldName, destFieldName := getFieldName(name, flgs, fieldNamesMapping)

				if fromField := fieldByNameOrZeroValue(source, srcFieldName); fromField.IsValid() && !shouldIgnore(fromField, fieldFlags, opt.IgnoreEmpty) {
					// process for nested anonymous field
					destFieldNotSet := false
					if f, ok := dest.Type().FieldByName(destFieldName); ok {
						// only initialize parent embedded struct pointer in the path
						for idx := range f.Index[:len(f.Index)-1] {
							destField := dest.FieldByIndex(f.Index[:idx+1])

							if destField.Kind() != reflect.Ptr {
								continue
							}

							if !destField.IsNil() {
								continue
							}
							if !destField.CanSet() {
								destFieldNotSet = true
								break
							}

							// destField is a nil pointer that can be set
							newValue := reflect.New(destField.Type().Elem())
							destField.Set(newValue)
						}
					}

					if destFieldNotSet {
						break
					}

					toField := fieldByName(dest, destFieldName, opt.CaseSensitive)
					if toField.IsValid() {
						if toField.CanSet() {
							isSet, err := set(toField, fromField, opt.DeepCopy, converters)
							if err != nil {
								return err
							}
							if !isSet {
								if err := copier(toField.Addr().Interface(), fromField.Interface(), opt); err != nil {
									return err
								}
							}
							if fieldFlags != 0 {
								// Note that a copy was made
								flgs.BitFlags[name] = fieldFlags | hasCopied
							}
						}
					} else {
						// try to set to method
						var toMethod reflect.Value
						if dest.CanAddr() {
							toMethod = dest.Addr().MethodByName(destFieldName)
						} else {
							toMethod = dest.MethodByName(destFieldName)
						}

						if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
							toMethod.Call([]reflect.Value{fromField})
						}
					}
				}
			}

			// Copy from from method to dest field
			for _, field := range deepFields(toType) {
				name := field.Name
				srcFieldName, destFieldName := getFieldName(name, flgs, getFieldNamesMapping(mappings, fromType, toType))

				var fromMethod reflect.Value
				if source.CanAddr() {
					fromMethod = source.Addr().MethodByName(srcFieldName)
				} else {
					fromMethod = source.MethodByName(srcFieldName)
				}

				if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 && !shouldIgnore(fromMethod, flgs.BitFlags[name], opt.IgnoreEmpty) {
					if toField := fieldByName(dest, destFieldName, opt.CaseSensitive); toField.IsValid() && toField.CanSet() {
						values := fromMethod.Call([]reflect.Value{})
						if len(values) >= 1 {
							// nolint:errcheck
							set(toField, values[0], opt.DeepCopy, converters)
						}
					}
				}
			}
		}

		if isSlice && to.Kind() == reflect.Slice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				if to.Len() < i+1 {
					to.Set(reflect.Append(to, dest.Addr()))
				} else {
					isSet, err := set(to.Index(i), dest.Addr(), opt.DeepCopy, converters)
					if err != nil {
						return err
					}
					if !isSet {
						// ignore error while copy slice element
						err = copier(to.Index(i).Addr().Interface(), dest.Addr().Interface(), opt)
						if err != nil {
							continue
						}
					}
				}
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				if to.Len() < i+1 {
					to.Set(reflect.Append(to, dest))
				} else {
					isSet, err := set(to.Index(i), dest, opt.DeepCopy, converters)
					if err != nil {
						return err
					}
					if !isSet {
						// ignore error while copy slice element
						err = copier(to.Index(i).Addr().Interface(), dest.Interface(), opt)
						if err != nil {
							continue
						}
					}
				}
			}
		} else if initDest {
			to.Set(dest)
		}
	}

	return nil
}

func getFieldNamesMapping(mappings map[converterPair]FieldNameMapping, fromType reflect.Type, toType reflect.Type) map[string]string {
	var fieldNamesMapping map[string]string

	if len(mappings) > 0 {
		pair := converterPair{
			SrcType: fromType,
			DstType: toType,
		}
		if v, ok := mappings[pair]; ok {
			fieldNamesMapping = v.Mapping
		}
	}
	return fieldNamesMapping
}

func fieldByNameOrZeroValue(source reflect.Value, fieldName string) (value reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
			value = reflect.Value{}
		}
	}()

	return source.FieldByName(fieldName)
}

func copyUnexportedStructFields(to, from reflect.Value) {
	if from.Kind() != reflect.Struct || to.Kind() != reflect.Struct || !from.Type().AssignableTo(to.Type()) {
		return
	}

	// create a shallow copy of 'to' to get all fields
	tmp := indirect(reflect.New(to.Type()))
	tmp.Set(from)

	// revert exported fields
	for i := 0; i < to.NumField(); i++ {
		if tmp.Field(i).CanSet() {
			tmp.Field(i).Set(to.Field(i))
		}
	}
	to.Set(tmp)
}

func shouldIgnore(v reflect.Value, bitFlags uint8, ignoreEmpty bool) bool {
	return ignoreEmpty && bitFlags&tagOverride == 0 && v.IsZero()
}

var (
	deepFieldsLock sync.RWMutex
	deepFieldsMap  = make(map[reflect.Type][]reflect.StructField)
)

func deepFields(reflectType reflect.Type) []reflect.StructField {
	deepFieldsLock.RLock()
	cache, ok := deepFieldsMap[reflectType]
	deepFieldsLock.RUnlock()
	if ok {
		return cache
	}
	var res []reflect.StructField
	if reflectType, _ = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		fields := make([]reflect.StructField, 0, reflectType.NumField())

		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			// PkgPath is the package path that qualifies a lower case (unexported)
			// field name. It is empty for upper case (exported) field names.
			// See https://golang.org/ref/spec#Uniqueness_of_identifiers
			if v.PkgPath == "" {
				fields = append(fields, v)
				if v.Anonymous {
					// also consider fields of anonymous fields as fields of the root
					fields = append(fields, deepFields(v.Type)...)
				}
			}
		}
		res = fields
	}

	deepFieldsLock.Lock()
	deepFieldsMap[reflectType] = res
	deepFieldsLock.Unlock()

	return res
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) (_ reflect.Type, isPtr bool) {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
		isPtr = true
	}
	return reflectType, isPtr
}

// nolint:gocyclo
func set(to, from reflect.Value, deepCopy bool, converters map[converterPair]TypeConverter) (bool, error) {
	if !from.IsValid() {
		return true, nil
	}
	if ok, err := lookupAndCopyWithConverter(to, from, converters); err != nil {
		return false, err
	} else if ok {
		return true, nil
	}

	if to.Kind() == reflect.Ptr {
		// nolint:gocritic
		if from.Kind() == reflect.Ptr && from.IsNil() {
			to.Set(reflect.Zero(to.Type()))
			return true, nil
		} else if to.IsNil() {
			// `from`         -> `to`
			// sql.NullString -> *string
			if fromValuer, ok := driverValuer(from); ok {
				v, err := fromValuer.Value()
				if err != nil {
					return true, err
				}

				// if `from` is not valid do nothing with `to`
				if v == nil {
					return true, nil
				}
			}
			// allocate new `to` variable with default value (eg. *string -> new(string))
			to.Set(reflect.New(to.Type().Elem()))
		} else if from.Kind() != reflect.Ptr && from.IsZero() {
			to.Set(reflect.Zero(to.Type()))
			return true, nil
		}
		// depointer `to`
		to = to.Elem()
	}

	if deepCopy {
		toKind := to.Kind()
		if toKind == reflect.Interface && to.IsNil() {
			if reflect.TypeOf(from.Interface()) != nil {
				to.Set(reflect.New(reflect.TypeOf(from.Interface())).Elem())
				toKind = reflect.TypeOf(to.Interface()).Kind()
			}
		}
		if from.Kind() == reflect.Ptr && from.IsNil() {
			to.Set(reflect.Zero(to.Type()))
			return true, nil
		}
		if _, ok := to.Addr().Interface().(sql.Scanner); !ok && (toKind == reflect.Struct || toKind == reflect.Map || toKind == reflect.Slice) {
			return false, nil
		}
	}

	// try convert directly
	if from.Type().ConvertibleTo(to.Type()) {
		to.Set(from.Convert(to.Type()))
		return true, nil
	}

	// try Scanner
	if toScanner, ok := to.Addr().Interface().(sql.Scanner); ok {
		// `from`  -> `to`
		// *string -> sql.NullString
		if from.Kind() == reflect.Ptr {
			// if `from` is nil do nothing with `to`
			if from.IsNil() {
				return true, nil
			}
			// depointer `from`
			from = indirect(from)
		}
		// `from` -> `to`
		// string -> sql.NullString
		// set `to` by invoking method Scan(`from`)
		err := toScanner.Scan(from.Interface())
		if err == nil {
			return true, nil
		}
	}

	// try Valuer
	if fromValuer, ok := driverValuer(from); ok {
		// `from`         -> `to`
		// sql.NullString -> string
		v, err := fromValuer.Value()
		if err != nil {
			return false, err
		}

		// if `from` is not valid do nothing with `to`
		if v == nil {
			return true, nil
		}

		rv := reflect.ValueOf(v)
		if rv.Type().AssignableTo(to.Type()) {
			to.Set(rv)
			return true, nil
		}
		if to.CanSet() && rv.Type().ConvertibleTo(to.Type()) {
			to.Set(rv.Convert(to.Type()))
			return true, nil
		}
		return false, nil
	}

	// from is ptr
	if from.Kind() == reflect.Ptr {
		return set(to, from.Elem(), deepCopy, converters)
	}

	return false, nil
}

// lookupAndCopyWithConverter looks up the type pair, on success the TypeConverter Fn func is called to copy src to dst field.
func lookupAndCopyWithConverter(to, from reflect.Value, converters map[converterPair]TypeConverter) (copied bool, err error) {
	pair := converterPair{
		SrcType: from.Type(),
		DstType: to.Type(),
	}

	if cnv, ok := converters[pair]; ok {
		result, err := cnv.Fn(from.Interface())
		if err != nil {
			return false, err
		}

		if result != nil {
			to.Set(reflect.ValueOf(result))
		} else {
			// in case we've got a nil value to copy
			to.Set(reflect.Zero(to.Type()))
		}

		return true, nil
	}

	return false, nil
}

// parseTags Parses struct tags and returns uint8 bit flags.
func parseTags(tag string) (flg uint8, name string, err error) {
	for _, t := range strings.Split(tag, ",") {
		switch t {
		case "-":
			flg = tagIgnore
			return
		case "must":
			flg |= tagMust
		case "nopanic":
			flg |= tagNoPanic
		case "override":
			flg |= tagOverride
		default:
			r, _ := utf8.DecodeRuneInString(t)
			if unicode.IsUpper(r) {
				name = strings.TrimSpace(t)
			} else {
				err = ErrFieldNameTagStartNotUpperCase
			}
		}
	}
	return
}

// getTagFlags Parses struct tags for bit flags, field name.
// nolint:gocyclo
func getFlags(dest, src reflect.Value, toType, fromType reflect.Type) (flags, error) {
	flgs := flags{
		BitFlags: map[string]uint8{},
		SrcNames: tagNameMapping{
			FieldNameToTag: map[string]string{},
			TagToFieldName: map[string]string{},
		},
		DestNames: tagNameMapping{
			FieldNameToTag: map[string]string{},
			TagToFieldName: map[string]string{},
		},
	}

	var toTypeFields, fromTypeFields []reflect.StructField
	if dest.IsValid() {
		toTypeFields = deepFields(toType)
	}
	if src.IsValid() {
		fromTypeFields = deepFields(fromType)
	}

	// Get a list dest of tags
	for _, field := range toTypeFields {
		tags := field.Tag.Get("copy")

		if tags != "" {
			var name string
			var err error

			if flgs.BitFlags[field.Name], name, err = parseTags(tags); err != nil {
				return flags{}, err
			} else if name != "" {
				flgs.DestNames.FieldNameToTag[field.Name] = name
				flgs.DestNames.TagToFieldName[name] = field.Name
			}
		}
	}

	// Get a list source of tags
	for _, field := range fromTypeFields {
		tags := field.Tag.Get("copy")
		if tags != "" {
			var name string
			var err error

			if _, name, err = parseTags(tags); err != nil {
				return flags{}, err
			} else if name != "" {
				flgs.SrcNames.FieldNameToTag[field.Name] = name
				flgs.SrcNames.TagToFieldName[name] = field.Name
			}
		}
	}

	return flgs, nil
}

func getFieldName(fieldName string, flgs flags, fieldNameMapping map[string]string) (srcFieldName string, destFieldName string) {
	if name, ok := fieldNameMapping[fieldName]; ok {
		return srcFieldName, name
	}

	if srcTagName, ok := flgs.SrcNames.FieldNameToTag[fieldName]; ok {
		destFieldName = srcTagName
		if destTagName, ok := flgs.DestNames.TagToFieldName[srcTagName]; ok {
			destFieldName = destTagName
		}
	} else {
		if destTagName, ok := flgs.DestNames.TagToFieldName[fieldName]; ok {
			destFieldName = destTagName
		}
	}
	if destFieldName == "" {
		destFieldName = fieldName
	}

	// get source field name
	if destTagName, ok := flgs.DestNames.FieldNameToTag[fieldName]; ok {
		srcFieldName = destTagName
		if srcField, ok := flgs.SrcNames.TagToFieldName[destTagName]; ok {
			srcFieldName = srcField
		}
	} else {
		if srcField, ok := flgs.SrcNames.TagToFieldName[fieldName]; ok {
			srcFieldName = srcField
		}
	}

	if srcFieldName == "" {
		srcFieldName = fieldName
	}

	return srcFieldName, destFieldName
}

func driverValuer(v reflect.Value) (i driver.Valuer, ok bool) {
	if !v.CanAddr() {
		i, ok = v.Interface().(driver.Valuer)
		return
	}

	i, ok = v.Addr().Interface().(driver.Valuer)
	return
}

func fieldByName(v reflect.Value, name string, caseSensitive bool) reflect.Value {
	if caseSensitive {
		return v.FieldByName(name)
	}

	return v.FieldByNameFunc(func(n string) bool { return strings.EqualFold(n, name) })
}
