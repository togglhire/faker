package faker

// Faker is a simple fake data generator for your own struct.
// Save your time, and Fake your data for your testing now.
import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/togglhire/faker/v3/support/slice"
)

var (
	mu = &sync.Mutex{}
	// Sets nil if the value type is struct or map and the size of it equals to zero.
	shouldSetNil = false
	//Sets random integer generation to zero for slice and maps
	testRandZero = false
	//Sets the default number of string when it is created randomly.
	randomStringLen = 25
	//Sets the boundary for random value generation. Boundaries can not exceed integer(4 byte...)
	nBoundary = numberBoundary{start: 0, end: 100}
	//Sets the random size for slices and maps.
	randomSize = 100
	// Sets the single fake data generator to generate unique values
	generateUniqueValues = false
	// Unique values are kept in memory so the generator retries if the value already exists
	uniqueValues = map[string][]interface{}{}
	// Uses randomSize as constant
	isFixedSize = false
)

type numberBoundary struct {
	start int
	end   int
}

// Supported tags
const (
	letterIdxBits         = 6                    // 6 bits to represent a letter index
	letterIdxMask         = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax          = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	maxRetry              = 10000                // max number of retry for unique values
	letterBytes           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tagName               = "faker"
	keep                  = "keep"
	unique                = "unique"
	ID                    = "uuid_digit"
	HyphenatedID          = "uuid_hyphenated"
	EmailTag              = "email"
	MacAddressTag         = "mac_address"
	DomainNameTag         = "domain_name"
	UserNameTag           = "username"
	URLTag                = "url"
	IPV4Tag               = "ipv4"
	IPV6Tag               = "ipv6"
	PASSWORD              = "password"
	LATITUDE              = "lat"
	LONGITUDE             = "long"
	CreditCardNumber      = "cc_number"
	CreditCardType        = "cc_type"
	PhoneNumber           = "phone_number"
	TollFreeNumber        = "toll_free_number"
	E164PhoneNumberTag    = "e_164_phone_number"
	TitleMaleTag          = "title_male"
	TitleFemaleTag        = "title_female"
	FirstNameTag          = "first_name"
	FirstNameMaleTag      = "first_name_male"
	FirstNameFemaleTag    = "first_name_female"
	LastNameTag           = "last_name"
	NAME                  = "name"
	UnixTimeTag           = "unix_time"
	DATE                  = "date"
	TIME                  = "time"
	MonthNameTag          = "month_name"
	YEAR                  = "year"
	DayOfWeekTag          = "day_of_week"
	DayOfMonthTag         = "day_of_month"
	TIMESTAMP             = "timestamp"
	CENTURY               = "century"
	TIMEZONE              = "timezone"
	TimePeriodTag         = "time_period"
	WORD                  = "word"
	SENTENCE              = "sentence"
	PARAGRAPH             = "paragraph"
	CurrencyTag           = "currency"
	AmountTag             = "amount"
	AmountWithCurrencyTag = "amount_with_currency"
	SKIP                  = "-"
	Length                = "len"
	BoundaryStart         = "boundary_start"
	BoundaryEnd           = "boundary_end"
	Equals                = "="
	Use                   = "use"
	comma                 = ","
)

var defaultTag = map[string]string{
	EmailTag:              EmailTag,
	MacAddressTag:         MacAddressTag,
	DomainNameTag:         DomainNameTag,
	URLTag:                URLTag,
	UserNameTag:           UserNameTag,
	IPV4Tag:               IPV4Tag,
	IPV6Tag:               IPV6Tag,
	PASSWORD:              PASSWORD,
	CreditCardType:        CreditCardType,
	CreditCardNumber:      CreditCardNumber,
	LATITUDE:              LATITUDE,
	LONGITUDE:             LONGITUDE,
	PhoneNumber:           PhoneNumber,
	TollFreeNumber:        TollFreeNumber,
	E164PhoneNumberTag:    E164PhoneNumberTag,
	TitleMaleTag:          TitleMaleTag,
	TitleFemaleTag:        TitleFemaleTag,
	FirstNameTag:          FirstNameTag,
	FirstNameMaleTag:      FirstNameMaleTag,
	FirstNameFemaleTag:    FirstNameFemaleTag,
	LastNameTag:           LastNameTag,
	NAME:                  NAME,
	UnixTimeTag:           UnixTimeTag,
	DATE:                  DATE,
	TIME:                  TimeFormat,
	MonthNameTag:          MonthNameTag,
	YEAR:                  YearFormat,
	DayOfWeekTag:          DayOfWeekTag,
	DayOfMonthTag:         DayOfMonthFormat,
	TIMESTAMP:             TIMESTAMP,
	CENTURY:               CENTURY,
	TIMEZONE:              TIMEZONE,
	TimePeriodTag:         TimePeriodFormat,
	WORD:                  WORD,
	SENTENCE:              SENTENCE,
	PARAGRAPH:             PARAGRAPH,
	CurrencyTag:           CurrencyTag,
	AmountTag:             AmountTag,
	AmountWithCurrencyTag: AmountWithCurrencyTag,
	ID:                    ID,
	HyphenatedID:          HyphenatedID,
}

// TaggedFunction used as the standard layout function for tag providers in struct.
// This type also can be used for custom provider.
type TaggedFunction func(v reflect.Value) (interface{}, error)

var mapperTag = map[string]TaggedFunction{
	EmailTag:              GetNetworker().Email,
	MacAddressTag:         GetNetworker().MacAddress,
	DomainNameTag:         GetNetworker().DomainName,
	URLTag:                GetNetworker().URL,
	UserNameTag:           GetNetworker().UserName,
	IPV4Tag:               GetNetworker().IPv4,
	IPV6Tag:               GetNetworker().IPv6,
	PASSWORD:              GetNetworker().Password,
	CreditCardType:        GetPayment().CreditCardType,
	CreditCardNumber:      GetPayment().CreditCardNumber,
	LATITUDE:              GetAddress().Latitude,
	LONGITUDE:             GetAddress().Longitude,
	PhoneNumber:           GetPhoner().PhoneNumber,
	TollFreeNumber:        GetPhoner().TollFreePhoneNumber,
	E164PhoneNumberTag:    GetPhoner().E164PhoneNumber,
	TitleMaleTag:          GetPerson().TitleMale,
	TitleFemaleTag:        GetPerson().TitleFeMale,
	FirstNameTag:          GetPerson().FirstName,
	FirstNameMaleTag:      GetPerson().FirstNameMale,
	FirstNameFemaleTag:    GetPerson().FirstNameFemale,
	LastNameTag:           GetPerson().LastName,
	NAME:                  GetPerson().Name,
	UnixTimeTag:           GetDateTimer().UnixTime,
	DATE:                  GetDateTimer().Date,
	TIME:                  GetDateTimer().Time,
	MonthNameTag:          GetDateTimer().MonthName,
	YEAR:                  GetDateTimer().Year,
	DayOfWeekTag:          GetDateTimer().DayOfWeek,
	DayOfMonthTag:         GetDateTimer().DayOfMonth,
	TIMESTAMP:             GetDateTimer().Timestamp,
	CENTURY:               GetDateTimer().Century,
	TIMEZONE:              GetDateTimer().TimeZone,
	TimePeriodTag:         GetDateTimer().TimePeriod,
	WORD:                  GetLorem().Word,
	SENTENCE:              GetLorem().Sentence,
	PARAGRAPH:             GetLorem().Paragraph,
	CurrencyTag:           GetPrice().Currency,
	AmountTag:             GetPrice().Amount,
	AmountWithCurrencyTag: GetPrice().AmountWithCurrency,
	ID:                    GetIdentifier().Digit,
	HyphenatedID:          GetIdentifier().Hyphenated,
}

// Generic Error Messages for tags
// 		ErrUnsupportedKindPtr: Error when get fake from ptr
// 		ErrUnsupportedKind: Error on passing unsupported kind
// 		ErrValueNotPtr: Error when value is not pointer
// 		ErrTagNotSupported: Error when tag is not supported
// 		ErrTagAlreadyExists: Error when tag exists and call AddProvider
// 		ErrMoreArguments: Error on passing more arguments
// 		ErrNotSupportedPointer: Error when passing unsupported pointer
var (
	ErrUnsupportedKindPtr  = "Unsupported kind: %s Change Without using * (pointer) in Field of %s"
	ErrUnsupportedKind     = "Unsupported kind: %s"
	ErrValueNotPtr         = "Not a pointer value"
	ErrTagNotSupported     = "Tag unsupported: %s"
	ErrTagAlreadyExists    = "Tag exists"
	ErrMoreArguments       = "Passed more arguments than is possible : (%d)"
	ErrNotSupportedPointer = "Use sample:=new(%s)\n faker.FakeData(sample) instead"
	ErrSmallerThanZero     = "Size:%d is smaller than zero."
	ErrUniqueFailure       = "Failed to generate a unique value for field \"%s\""

	ErrStartValueBiggerThanEnd = "Start value can not be bigger than end value."
	ErrWrongFormattedTag       = "Tag \"%s\" is not written properly"
	ErrUnknownType             = "Unknown Type"
	ErrNotSupportedTypeForTag  = "Type is not supported by tag."
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// SetSeed allows custom seeds for rand
func SetSeed(seed int64) {
	rand.Seed(seed)
}

// ResetUnique is used to forget generated unique values.
// Call this when you're done generating a dataset.
func ResetUnique() {
	uniqueValues = map[string][]interface{}{}
}

// SetGenerateUniqueValues allows to set the single fake data generator functions to generate unique data.
func SetGenerateUniqueValues(unique bool) {
	generateUniqueValues = unique
}

// SetNilIfLenIsZero allows to set nil for the slice and maps, if size is 0.
func SetNilIfLenIsZero(setNil bool) {
	shouldSetNil = setNil
}

// SetRandomStringLength sets a length for random string generation
func SetRandomStringLength(size int) error {
	if size < 0 {
		return fmt.Errorf(ErrSmallerThanZero, size)
	}
	randomStringLen = size
	return nil
}

// SetFixedMapAndSliceSize sets the fixed size for maps and slices for generation.
func SetFixedMapAndSliceSize(size int) error {
	if size < 0 {
		return fmt.Errorf(ErrSmallerThanZero, size)
	}
	randomSize = size
	isFixedSize = true
	return nil
}

// SetRandomMapAndSliceSize sets the size for maps and slices for random generation.
func SetRandomMapAndSliceSize(size int) error {
	if size < 0 {
		return fmt.Errorf(ErrSmallerThanZero, size)
	}
	randomSize = size
	isFixedSize = false
	return nil
}

// SetRandomNumberBoundaries sets boundary for random number generation
func SetRandomNumberBoundaries(start, end int) error {
	if start > end {
		return errors.New(ErrStartValueBiggerThanEnd)
	}
	nBoundary = numberBoundary{start: start, end: end}
	return nil
}

// FakeData is the main function. Will generate a fake data based on your struct.  You can use this for automation testing, or anything that need automated data.
// You don't need to Create your own data for your testing.
func FakeData(a interface{}) error {

	reflectType := reflect.TypeOf(a)

	if reflectType.Kind() != reflect.Ptr {
		return errors.New(ErrValueNotPtr)
	}

	if reflect.ValueOf(a).IsNil() {
		return fmt.Errorf(ErrNotSupportedPointer, reflectType.Elem().String())
	}

	rval := reflect.ValueOf(a)

	finalValue, err := getValue(a)
	if err != nil {
		return err
	}

	rval.Elem().Set(finalValue.Elem().Convert(reflectType.Elem()))
	return nil
}

// AddProvider extend faker with tag to generate fake data with specified custom algorithm
// Example:
// 		type Gondoruwo struct {
// 			Name       string
// 			Locatadata int
// 		}
//
// 		type Sample struct {
// 			ID                 int64     `faker:"customIdFaker"`
// 			Gondoruwo          Gondoruwo `faker:"gondoruwo"`
// 			Danger             string    `faker:"danger"`
// 		}
//
// 		func CustomGenerator() {
// 			// explicit
// 			faker.AddProvider("customIdFaker", func(v reflect.Value) (interface{}, error) {
// 			 	return int64(43), nil
// 			})
// 			// functional
// 			faker.AddProvider("danger", func() faker.TaggedFunction {
// 				return func(v reflect.Value) (interface{}, error) {
// 					return "danger-ranger", nil
// 				}
// 			}())
// 			faker.AddProvider("gondoruwo", func(v reflect.Value) (interface{}, error) {
// 				obj := Gondoruwo{
// 					Name:       "Power",
// 					Locatadata: 324,
// 				}
// 				return obj, nil
// 			})
// 		}
//
// 		func main() {
// 			CustomGenerator()
// 			var sample Sample
// 			faker.FakeData(&sample)
// 			fmt.Printf("%+v", sample)
// 		}
//
// Will print
// 		{ID:43 Gondoruwo:{Name:Power Locatadata:324} Danger:danger-ranger}
// Notes: when using a custom provider make sure to return the same type as the field
func AddProvider(tag string, provider TaggedFunction) error {
	if _, ok := mapperTag[tag]; ok {
		return errors.New(ErrTagAlreadyExists)
	}

	mapperTag[tag] = provider

	return nil
}

func getValue(a interface{}) (reflect.Value, error) {
	t := reflect.TypeOf(a)
	if t == nil {
		return reflect.Value{}, fmt.Errorf("interface{} not allowed")
	}
	k := t.Kind()

	switch k {
	case reflect.Ptr:
		v := reflect.New(t.Elem())
		var val reflect.Value
		var err error
		if a != reflect.Zero(reflect.TypeOf(a)).Interface() {
			val, err = getValue(reflect.ValueOf(a).Elem().Interface())
			if err != nil {
				return reflect.Value{}, err
			}
		} else {
			val, err = getValue(v.Elem().Interface())
			if err != nil {
				return reflect.Value{}, err
			}
		}
		v.Elem().Set(val.Convert(t.Elem()))
		return v, nil
	case reflect.Struct:
		switch t.String() {
		case "time.Time":
			ft := time.Now().Add(time.Duration(rand.Int63()))
			return reflect.ValueOf(ft), nil
		default:
			originalDataVal := reflect.ValueOf(a)
			v := reflect.New(t).Elem()
			retry := 0 // error if cannot generate unique value after maxRetry tries
			for i := 0; i < v.NumField(); i++ {
				if !v.Field(i).CanSet() {
					continue // to avoid panic to set on unexported field in struct
				}
				tags := decodeTags(t, i)
				switch {
				case tags.keepOriginal:
					zero, err := isZero(reflect.ValueOf(a).Field(i))
					if err != nil {
						return reflect.Value{}, err
					}
					if zero {
						err := setDataWithTag(v.Field(i).Addr(), tags.fieldType)
						if err != nil {
							return reflect.Value{}, err
						}
						continue
					}
					v.Field(i).Set(reflect.ValueOf(a).Field(i))
				case tags.fieldType == "":
					val, err := getValue(v.Field(i).Interface())
					if err != nil {
						return reflect.Value{}, err
					}
					val = val.Convert(v.Field(i).Type())
					v.Field(i).Set(val)
				case tags.fieldType == SKIP:
					item := originalDataVal.Field(i).Interface()
					if v.CanSet() && item != nil {
						v.Field(i).Set(reflect.ValueOf(item))
					}
				default:
					err := setDataWithTag(v.Field(i).Addr(), tags.fieldType)
					if err != nil {
						return reflect.Value{}, err
					}
				}

				if tags.unique {

					if retry >= maxRetry {
						return reflect.Value{}, fmt.Errorf(ErrUniqueFailure, reflect.TypeOf(a).Field(i).Name)
					}

					value := v.Field(i).Interface()
					if slice.ContainsValue(uniqueValues[tags.fieldType], value) { // Retry if unique value already found
						i--
						retry++
						continue
					}
					retry = 0
					uniqueValues[tags.fieldType] = append(uniqueValues[tags.fieldType], value)
				} else {
					retry = 0
				}

			}
			return v, nil
		}

	case reflect.String:
		res := randomString(randomStringLen)
		return reflect.ValueOf(res), nil
	case reflect.Array, reflect.Slice:
		len := randomSliceAndMapSize()
		if shouldSetNil && len == 0 {
			return reflect.Zero(t), nil
		}
		v := reflect.MakeSlice(t, len, len)
		for i := 0; i < v.Len(); i++ {
			val, err := getValue(v.Index(i).Interface())
			if err != nil {
				return reflect.Value{}, err
			}
			v.Index(i).Set(val)
		}
		return v, nil
	case reflect.Int:
		return reflect.ValueOf(randomInteger()), nil
	case reflect.Int8:
		return reflect.ValueOf(int8(randomInteger())), nil
	case reflect.Int16:
		return reflect.ValueOf(int16(randomInteger())), nil
	case reflect.Int32:
		return reflect.ValueOf(int32(randomInteger())), nil
	case reflect.Int64:
		return reflect.ValueOf(int64(randomInteger())), nil
	case reflect.Float32:
		return reflect.ValueOf(rand.Float32()), nil
	case reflect.Float64:
		return reflect.ValueOf(rand.Float64()), nil
	case reflect.Bool:
		val := rand.Intn(2) > 0
		return reflect.ValueOf(val), nil

	case reflect.Uint:
		return reflect.ValueOf(uint(randomInteger())), nil

	case reflect.Uint8:
		return reflect.ValueOf(uint8(randomInteger())), nil

	case reflect.Uint16:
		return reflect.ValueOf(uint16(randomInteger())), nil

	case reflect.Uint32:
		return reflect.ValueOf(uint32(randomInteger())), nil

	case reflect.Uint64:
		return reflect.ValueOf(uint64(randomInteger())), nil

	case reflect.Map:
		len := randomSliceAndMapSize()
		if shouldSetNil && len == 0 {
			return reflect.Zero(t), nil
		}
		v := reflect.MakeMap(t)
		for i := 0; i < len; i++ {
			keyInstance := reflect.New(t.Key()).Elem().Interface()
			key, err := getValue(keyInstance)
			if err != nil {
				return reflect.Value{}, err
			}

			valueInstance := reflect.New(t.Elem()).Elem().Interface()
			val, err := getValue(valueInstance)
			if err != nil {
				return reflect.Value{}, err
			}
			v.SetMapIndex(key, val)
		}
		return v, nil
	default:
		err := fmt.Errorf("no support for kind %+v", t)
		return reflect.Value{}, err
	}

}

func isZero(field reflect.Value) (bool, error) {
	if field.Kind() == reflect.Map {
		return field.Len() == 0, nil
	}

	for _, kind := range []reflect.Kind{reflect.Struct, reflect.Slice, reflect.Array} {
		if kind == field.Kind() {
			return false, fmt.Errorf("keep not allowed on struct")
		}
	}
	return reflect.Zero(field.Type()).Interface() == field.Interface(), nil
}

func decodeTags(typ reflect.Type, i int) structTag {
	tags := strings.Split(typ.Field(i).Tag.Get(tagName), ",")

	keepOriginal := false
	uni := false
	res := make([]string, 0)
	for _, tag := range tags {
		if tag == keep {
			keepOriginal = true
			continue
		} else if tag == unique {
			uni = true
			continue
		}
		res = append(res, tag)
	}

	return structTag{
		fieldType:    strings.Join(res, ","),
		unique:       uni,
		keepOriginal: keepOriginal,
	}
}

type structTag struct {
	fieldType    string
	unique       bool
	keepOriginal bool
}

func setDataWithTag(v reflect.Value, tag string) error {
	if v.Kind() != reflect.Ptr {
		return errors.New(ErrValueNotPtr)
	}
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Ptr:
		if strings.Contains(tag, Use) {
			t := v.Type()
			newv := reflect.New(t.Elem()).Elem()
			err := setDataWithTagSwitch(newv, tag)
			if err != nil {
				return err
			}
			v.Set(newv.Addr())
			return nil
		}

		if _, exist := mapperTag[tag]; !exist {
			return fmt.Errorf(ErrTagNotSupported, tag)
		}
		if _, def := defaultTag[tag]; !def {
			res, err := mapperTag[tag](v)
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(res))
			return nil
		}

		t := v.Type()
		newv := reflect.New(t.Elem())
		res, err := mapperTag[tag](newv.Elem())
		if err != nil {
			return err
		}
		rval := reflect.ValueOf(res)
		newv.Elem().Set(rval)
		v.Set(newv)
		return nil
	default:
		return setDataWithTagSwitch(v, tag)
	}
}

func setDataWithTagSwitch(v reflect.Value, tag string) error {
	switch v.Kind() {
	case reflect.String:
		return userDefinedString(v, tag)
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return userDefinedNumber(v, tag)
	case reflect.Slice, reflect.Array:
		return userDefinedArray(v, tag)
	case reflect.Map:
		return userDefinedMap(v, tag)
	case reflect.Bool:
		return userDefinedBool(v, tag)
	default:
		if _, exist := mapperTag[tag]; !exist {
			return fmt.Errorf(ErrTagNotSupported, tag)
		}
		res, err := mapperTag[tag](v)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(res))
	}
	return nil
}

func userDefinedBool(v reflect.Value, tag string) error {
	var res interface{}
	var err error

	if strings.Contains(tag, Use) {
		res, err = extractBoolFromUseTag(tag)
		if err != nil {
			return err
		}
	}
	if res == nil {
		return fmt.Errorf(ErrTagNotSupported, tag)
	}
	val, _ := res.(bool)
	v.SetBool(val)
	return nil
}

func userDefinedMap(v reflect.Value, tag string) error {
	if tagFunc, ok := mapperTag[tag]; ok {
		res, err := tagFunc(v)
		if err != nil {
			return err
		}

		v.Set(reflect.ValueOf(res))
		return nil
	}

	len := randomSliceAndMapSize()
	if shouldSetNil && len == 0 {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	definedMap := reflect.MakeMap(v.Type())
	for i := 0; i < len; i++ {
		key, err := getValueWithTag(v.Type().Key(), tag)
		if err != nil {
			return err
		}
		val, err := getValueWithTag(v.Type().Elem(), tag)
		if err != nil {
			return err
		}
		definedMap.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	}
	v.Set(definedMap)
	return nil
}

func getValueWithTag(t reflect.Type, tag string) (interface{}, error) {
	switch t.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64:
		res, err := extractNumberFromTag(tag, t)
		if err != nil {
			return nil, err
		}
		return res, nil
	case reflect.String:
		res, err := extractStringFromTag(tag)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		return 0, errors.New(ErrUnknownType)
	}
}

func userDefinedArray(v reflect.Value, tag string) error {
	len := randomSliceAndMapSize()
	if shouldSetNil && len == 0 {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	array := reflect.MakeSlice(v.Type(), len, len)
	for i := 0; i < len; i++ {
		res, err := getValueWithTag(v.Type().Elem(), tag)
		if err != nil {
			return err
		}
		array.Index(i).Set(reflect.ValueOf(res))
	}
	v.Set(array)
	return nil
}

func userDefinedString(v reflect.Value, tag string) error {
	var res interface{}
	var err error

	if tagFunc, ok := mapperTag[tag]; ok {
		res, err = tagFunc(v)
		if err != nil {
			return err
		}
	} else if strings.Contains(tag, Length) {
		res, err = extractStringFromTag(tag)
		if err != nil {
			return err
		}
	} else {
		res, err = extractStringFromUseTag(tag)
		if err != nil {
			return err
		}
	}
	if res == nil {
		return fmt.Errorf(ErrTagNotSupported, tag)
	}
	val, _ := res.(string)
	v.SetString(val)
	return nil
}

func userDefinedNumber(v reflect.Value, tag string) error {
	var res interface{}
	var err error

	if tagFunc, ok := mapperTag[tag]; ok {
		res, err = tagFunc(v)
		if err != nil {
			return err
		}
	} else if strings.Contains(tag, BoundaryStart) {
		res, err = extractNumberFromTag(tag, v.Type())
		if err != nil {
			return err
		}
	} else {
		res, err = extractNumberFromUseTag(tag, v.Type())
		if err != nil {
			return err
		}
	}
	if res == nil {
		return fmt.Errorf(ErrTagNotSupported, tag)
	}

	v.Set(reflect.ValueOf(res))
	return nil
}

func extractStringFromTag(tag string) (interface{}, error) {
	if !strings.Contains(tag, Length) {
		return nil, fmt.Errorf(ErrTagNotSupported, tag)
	}
	len, err := extractNumberFromText(tag)
	if err != nil {
		return nil, err
	}
	res := randomString(int(len))
	return res, nil
}

func extractBoolFromUseTag(tag string) (interface{}, error) {
	if !strings.Contains(tag, Use) {
		return nil, fmt.Errorf(ErrTagNotSupported, tag)
	}
	str, err := extractStringFromText(tag)
	if err != nil {
		return nil, err
	}
	return strconv.ParseBool(str)
}

func extractStringFromUseTag(tag string) (interface{}, error) {
	if !strings.Contains(tag, Use) {
		return nil, fmt.Errorf(ErrTagNotSupported, tag)
	}
	return extractStringFromText(tag)
}

func extractNumberFromUseTag(tag string, t reflect.Type) (interface{}, error) {
	if !strings.Contains(tag, Use) {
		return nil, fmt.Errorf(ErrTagNotSupported, tag)
	}
	number, err := extractNumberFromText(tag)
	if err != nil {
		return nil, err
	}

	switch t.Kind() {
	case reflect.Uint:
		return uint(number), nil
	case reflect.Uint8:
		return uint8(number), nil
	case reflect.Uint16:
		return uint16(number), nil
	case reflect.Uint32:
		return uint32(number), nil
	case reflect.Uint64:
		return uint64(number), nil
	case reflect.Int:
		return int(number), nil
	case reflect.Int8:
		return int8(number), nil
	case reflect.Int16:
		return int16(number), nil
	case reflect.Int32:
		return int32(number), nil
	case reflect.Int64:
		return int64(number), nil
	case reflect.Float32:
		return float32(number), nil
	case reflect.Float64:
		return float64(number), nil
	default:
		return nil, errors.New(ErrNotSupportedTypeForTag)
	}
}

func extractNumberFromTag(tag string, t reflect.Type) (interface{}, error) {
	if !strings.Contains(tag, BoundaryStart) || !strings.Contains(tag, BoundaryEnd) {
		return nil, fmt.Errorf(ErrTagNotSupported, tag)
	}
	valuesStr := strings.SplitN(tag, comma, -1)
	if len(valuesStr) != 2 {
		return nil, fmt.Errorf(ErrWrongFormattedTag, tag)
	}
	startBoundary, err := extractNumberFromText(valuesStr[0])
	if err != nil {
		return nil, err
	}
	endBoundary, err := extractNumberFromText(valuesStr[1])
	if err != nil {
		return nil, err
	}
	boundary := numberBoundary{start: int(startBoundary), end: int(endBoundary)}
	switch t.Kind() {
	case reflect.Uint:
		return uint(randomIntegerWithBoundary(boundary)), nil
	case reflect.Uint8:
		return uint8(randomIntegerWithBoundary(boundary)), nil
	case reflect.Uint16:
		return uint16(randomIntegerWithBoundary(boundary)), nil
	case reflect.Uint32:
		return uint32(randomIntegerWithBoundary(boundary)), nil
	case reflect.Uint64:
		return uint64(randomIntegerWithBoundary(boundary)), nil
	case reflect.Int:
		return randomIntegerWithBoundary(boundary), nil
	case reflect.Int8:
		return int8(randomIntegerWithBoundary(boundary)), nil
	case reflect.Int16:
		return int16(randomIntegerWithBoundary(boundary)), nil
	case reflect.Int32:
		return int32(randomIntegerWithBoundary(boundary)), nil
	case reflect.Int64:
		return int64(randomIntegerWithBoundary(boundary)), nil
	default:
		return nil, errors.New(ErrNotSupportedTypeForTag)
	}
}

func extractNumberFromText(text string) (float64, error) {
	text = strings.TrimSpace(text)
	texts := strings.SplitN(text, Equals, -1)
	if len(texts) != 2 {
		return 0, fmt.Errorf(ErrWrongFormattedTag, text)
	}
	result, err := strconv.ParseFloat(texts[1], 64)
	if err != nil {
		return 0, fmt.Errorf(ErrWrongFormattedTag, text)
	}
	return result, nil
}

func extractStringFromText(text string) (string, error) {
	text = strings.TrimSpace(text)
	texts := strings.SplitN(text, Equals, -1)
	if len(texts) != 2 {
		return "", fmt.Errorf(ErrWrongFormattedTag, text)
	}
	return texts[1], nil
}

func randomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// randomIntegerWithBoundary returns a random integer between input start and end boundary. [start, end)
func randomIntegerWithBoundary(boundary numberBoundary) int {
	return rand.Intn(boundary.end-boundary.start) + boundary.start
}

// randomInteger returns a random integer between start and end boundary. [start, end)
func randomInteger() int {
	return rand.Intn(nBoundary.end-nBoundary.start) + nBoundary.start
}

// randomSliceAndMapSize returns a random integer between [0,randomSliceAndMapSize). If the testRandZero is set, returns 0
// Written for test purposes for shouldSetNil
func randomSliceAndMapSize() int {
	if testRandZero {
		return 0
	}
	if isFixedSize {
		return randomSize
	}
	return rand.Intn(randomSize)
}

func randomElementFromSliceString(s []string) string {
	return s[rand.Int()%len(s)]
}
func randomStringNumber(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(numberBytes) {
			b[i] = numberBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandomInt Get three parameters , only first mandatory and the rest are optional
// 		If only set one parameter :  This means the minimum number of digits and the total number
// 		If only set two parameters : First this is min digit and second max digit and the total number the difference between them
// 		If only three parameters: the third argument set Max count Digit
func RandomInt(parameters ...int) (p []int, err error) {
	switch len(parameters) {
	case 1:
		minCount := parameters[0]
		p = rand.Perm(minCount)
		for i := range p {
			p[i] += minCount
		}
	case 2:
		minDigit, maxDigit := parameters[0], parameters[1]
		p = rand.Perm(maxDigit - minDigit + 1)

		for i := range p {
			p[i] += minDigit
		}
	default:
		err = fmt.Errorf(ErrMoreArguments, len(parameters))
	}
	return p, err
}

func generateUnique(dataType string, fn func() interface{}) (interface{}, error) {
	for i := 0; i < maxRetry; i++ {
		value := fn()
		if !slice.ContainsValue(uniqueValues[dataType], value) { // Retry if unique value already found
			uniqueValues[dataType] = append(uniqueValues[dataType], value)
			return value, nil
		}
	}
	return reflect.Value{}, fmt.Errorf(ErrUniqueFailure, dataType)
}

func singleFakeData(dataType string, fn func() interface{}) interface{} {
	if generateUniqueValues {
		v, err := generateUnique(dataType, fn)
		if err != nil {
			panic(err)
		}
		return v
	}
	return fn()
}
