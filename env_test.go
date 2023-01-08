package env

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

type unmarshaler struct {
	time.Duration
}

// TextUnmarshaler implements encoding.TextUnmarshaler.
func (d *unmarshaler) UnmarshalText(data []byte) (err error) {
	if len(data) != 0 {
		d.Duration, err = time.ParseDuration(string(data))
	} else {
		d.Duration = 0
	}
	return err
}

// nolint: maligned
type Config struct {
	String     string    `env:"STRING"`
	StringPtr  *string   `env:"STRING"`
	Strings    []string  `env:"STRINGS"`
	StringPtrs []*string `env:"STRINGS"`

	Bool     bool    `env:"BOOL"`
	BoolPtr  *bool   `env:"BOOL"`
	Bools    []bool  `env:"BOOLS"`
	BoolPtrs []*bool `env:"BOOLS"`

	Int     int    `env:"INT"`
	IntPtr  *int   `env:"INT"`
	Ints    []int  `env:"INTS"`
	IntPtrs []*int `env:"INTS"`

	Int8     int8    `env:"INT8"`
	Int8Ptr  *int8   `env:"INT8"`
	Int8s    []int8  `env:"INT8S"`
	Int8Ptrs []*int8 `env:"INT8S"`

	Int16     int16    `env:"INT16"`
	Int16Ptr  *int16   `env:"INT16"`
	Int16s    []int16  `env:"INT16S"`
	Int16Ptrs []*int16 `env:"INT16S"`

	Int32     int32    `env:"INT32"`
	Int32Ptr  *int32   `env:"INT32"`
	Int32s    []int32  `env:"INT32S"`
	Int32Ptrs []*int32 `env:"INT32S"`

	Int64     int64    `env:"INT64"`
	Int64Ptr  *int64   `env:"INT64"`
	Int64s    []int64  `env:"INT64S"`
	Int64Ptrs []*int64 `env:"INT64S"`

	Uint     uint    `env:"UINT"`
	UintPtr  *uint   `env:"UINT"`
	Uints    []uint  `env:"UINTS"`
	UintPtrs []*uint `env:"UINTS"`

	Uint8     uint8    `env:"UINT8"`
	Uint8Ptr  *uint8   `env:"UINT8"`
	Uint8s    []uint8  `env:"UINT8S"`
	Uint8Ptrs []*uint8 `env:"UINT8S"`

	Uint16     uint16    `env:"UINT16"`
	Uint16Ptr  *uint16   `env:"UINT16"`
	Uint16s    []uint16  `env:"UINT16S"`
	Uint16Ptrs []*uint16 `env:"UINT16S"`

	Uint32     uint32    `env:"UINT32"`
	Uint32Ptr  *uint32   `env:"UINT32"`
	Uint32s    []uint32  `env:"UINT32S"`
	Uint32Ptrs []*uint32 `env:"UINT32S"`

	Uint64     uint64    `env:"UINT64"`
	Uint64Ptr  *uint64   `env:"UINT64"`
	Uint64s    []uint64  `env:"UINT64S"`
	Uint64Ptrs []*uint64 `env:"UINT64S"`

	Float32     float32    `env:"FLOAT32"`
	Float32Ptr  *float32   `env:"FLOAT32"`
	Float32s    []float32  `env:"FLOAT32S"`
	Float32Ptrs []*float32 `env:"FLOAT32S"`

	Float64     float64    `env:"FLOAT64"`
	Float64Ptr  *float64   `env:"FLOAT64"`
	Float64s    []float64  `env:"FLOAT64S"`
	Float64Ptrs []*float64 `env:"FLOAT64S"`

	Duration     time.Duration    `env:"DURATION"`
	Durations    []time.Duration  `env:"DURATIONS"`
	DurationPtr  *time.Duration   `env:"DURATION"`
	DurationPtrs []*time.Duration `env:"DURATIONS"`

	Unmarshaler     unmarshaler    `env:"UNMARSHALER"`
	UnmarshalerPtr  *unmarshaler   `env:"UNMARSHALER"`
	Unmarshalers    []unmarshaler  `env:"UNMARSHALERS"`
	UnmarshalerPtrs []*unmarshaler `env:"UNMARSHALERS"`

	URL     url.URL    `env:"URL"`
	URLPtr  *url.URL   `env:"URL"`
	URLs    []url.URL  `env:"URLS"`
	URLPtrs []*url.URL `env:"URLS"`

	StringWithdefault string `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/db"`

	CustomSeparator []string `env:"SEPSTRINGS" envSeparator:":"`

	NonDefined struct {
		String string `env:"NONDEFINED_STR"`
	}

	NestedNonDefined struct {
		NonDefined struct {
			String string `env:"STR"`
		} `envPrefix:"NONDEFINED_"`
	} `envPrefix:"PRF_"`

	NotAnEnv   string
	unexported string `env:"FOO"`
}

type ParentStruct struct {
	InnerStruct *InnerStruct
	unexported  *InnerStruct
	Ignored     *http.Client
}

type InnerStruct struct {
	Inner  string `env:"innervar"`
	Number uint   `env:"innernum"`
}

type ForNestedStruct struct {
	NestedStruct
}

type NestedStruct struct {
	NestedVar string `env:"nestedvar"`
}

func TestParsesEnv(t *testing.T) {
	tos := func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}

	toss := func(v ...interface{}) string {
		ss := []string{}
		for _, s := range v {
			ss = append(ss, tos(s))
		}
		return strings.Join(ss, ",")
	}

	str1 := "str1"
	str2 := "str2"
	t.Setenv("STRING", str1)
	t.Setenv("STRINGS", toss(str1, str2))

	bool1 := true
	bool2 := false
	t.Setenv("BOOL", tos(bool1))
	t.Setenv("BOOLS", toss(bool1, bool2))

	int1 := -1
	int2 := 2
	t.Setenv("INT", tos(int1))
	t.Setenv("INTS", toss(int1, int2))

	var int81 int8 = -2
	var int82 int8 = 5
	t.Setenv("INT8", tos(int81))
	t.Setenv("INT8S", toss(int81, int82))

	var int161 int16 = -24
	var int162 int16 = 15
	t.Setenv("INT16", tos(int161))
	t.Setenv("INT16S", toss(int161, int162))

	var int321 int32 = -14
	var int322 int32 = 154
	t.Setenv("INT32", tos(int321))
	t.Setenv("INT32S", toss(int321, int322))

	var int641 int64 = -12
	var int642 int64 = 150
	t.Setenv("INT64", tos(int641))
	t.Setenv("INT64S", toss(int641, int642))

	var uint1 uint = 1
	var uint2 uint = 2
	t.Setenv("UINT", tos(uint1))
	t.Setenv("UINTS", toss(uint1, uint2))

	var uint81 uint8 = 15
	var uint82 uint8 = 51
	t.Setenv("UINT8", tos(uint81))
	t.Setenv("UINT8S", toss(uint81, uint82))

	var uint161 uint16 = 532
	var uint162 uint16 = 123
	t.Setenv("UINT16", tos(uint161))
	t.Setenv("UINT16S", toss(uint161, uint162))

	var uint321 uint32 = 93
	var uint322 uint32 = 14
	t.Setenv("UINT32", tos(uint321))
	t.Setenv("UINT32S", toss(uint321, uint322))

	var uint641 uint64 = 5
	var uint642 uint64 = 43
	t.Setenv("UINT64", tos(uint641))
	t.Setenv("UINT64S", toss(uint641, uint642))

	var float321 float32 = 9.3
	var float322 float32 = 1.1
	t.Setenv("FLOAT32", tos(float321))
	t.Setenv("FLOAT32S", toss(float321, float322))

	float641 := 1.53
	float642 := 0.5
	t.Setenv("FLOAT64", tos(float641))
	t.Setenv("FLOAT64S", toss(float641, float642))

	duration1 := time.Second
	duration2 := time.Second * 4
	t.Setenv("DURATION", tos(duration1))
	t.Setenv("DURATIONS", toss(duration1, duration2))

	unmarshaler1 := unmarshaler{time.Minute}
	unmarshaler2 := unmarshaler{time.Millisecond * 1232}
	t.Setenv("UNMARSHALER", tos(unmarshaler1.Duration))
	t.Setenv("UNMARSHALERS", toss(unmarshaler1.Duration, unmarshaler2.Duration))

	url1 := "https://goreleaser.com"
	url2 := "https://caarlos0.dev"
	t.Setenv("URL", tos(url1))
	t.Setenv("URLS", toss(url1, url2))

	t.Setenv("SEPSTRINGS", strings.Join([]string{str1, str2}, ":"))

	nonDefinedStr := "nonDefinedStr"
	t.Setenv("NONDEFINED_STR", nonDefinedStr)
	t.Setenv("PRF_NONDEFINED_STR", nonDefinedStr)

	cfg := Config{}
	isNoErr(t, Parse(&cfg))

	isEqual(t, str1, cfg.String)
	isEqual(t, &str1, cfg.StringPtr)
	isEqual(t, str1, cfg.Strings[0])
	isEqual(t, str2, cfg.Strings[1])
	isEqual(t, &str1, cfg.StringPtrs[0])
	isEqual(t, &str2, cfg.StringPtrs[1])

	isEqual(t, bool1, cfg.Bool)
	isEqual(t, &bool1, cfg.BoolPtr)
	isEqual(t, bool1, cfg.Bools[0])
	isEqual(t, bool2, cfg.Bools[1])
	isEqual(t, &bool1, cfg.BoolPtrs[0])
	isEqual(t, &bool2, cfg.BoolPtrs[1])

	isEqual(t, int1, cfg.Int)
	isEqual(t, &int1, cfg.IntPtr)
	isEqual(t, int1, cfg.Ints[0])
	isEqual(t, int2, cfg.Ints[1])
	isEqual(t, &int1, cfg.IntPtrs[0])
	isEqual(t, &int2, cfg.IntPtrs[1])

	isEqual(t, int81, cfg.Int8)
	isEqual(t, &int81, cfg.Int8Ptr)
	isEqual(t, int81, cfg.Int8s[0])
	isEqual(t, int82, cfg.Int8s[1])
	isEqual(t, &int81, cfg.Int8Ptrs[0])
	isEqual(t, &int82, cfg.Int8Ptrs[1])

	isEqual(t, int161, cfg.Int16)
	isEqual(t, &int161, cfg.Int16Ptr)
	isEqual(t, int161, cfg.Int16s[0])
	isEqual(t, int162, cfg.Int16s[1])
	isEqual(t, &int161, cfg.Int16Ptrs[0])
	isEqual(t, &int162, cfg.Int16Ptrs[1])

	isEqual(t, int321, cfg.Int32)
	isEqual(t, &int321, cfg.Int32Ptr)
	isEqual(t, int321, cfg.Int32s[0])
	isEqual(t, int322, cfg.Int32s[1])
	isEqual(t, &int321, cfg.Int32Ptrs[0])
	isEqual(t, &int322, cfg.Int32Ptrs[1])

	isEqual(t, int641, cfg.Int64)
	isEqual(t, &int641, cfg.Int64Ptr)
	isEqual(t, int641, cfg.Int64s[0])
	isEqual(t, int642, cfg.Int64s[1])
	isEqual(t, &int641, cfg.Int64Ptrs[0])
	isEqual(t, &int642, cfg.Int64Ptrs[1])

	isEqual(t, uint1, cfg.Uint)
	isEqual(t, &uint1, cfg.UintPtr)
	isEqual(t, uint1, cfg.Uints[0])
	isEqual(t, uint2, cfg.Uints[1])
	isEqual(t, &uint1, cfg.UintPtrs[0])
	isEqual(t, &uint2, cfg.UintPtrs[1])

	isEqual(t, uint81, cfg.Uint8)
	isEqual(t, &uint81, cfg.Uint8Ptr)
	isEqual(t, uint81, cfg.Uint8s[0])
	isEqual(t, uint82, cfg.Uint8s[1])
	isEqual(t, &uint81, cfg.Uint8Ptrs[0])
	isEqual(t, &uint82, cfg.Uint8Ptrs[1])

	isEqual(t, uint161, cfg.Uint16)
	isEqual(t, &uint161, cfg.Uint16Ptr)
	isEqual(t, uint161, cfg.Uint16s[0])
	isEqual(t, uint162, cfg.Uint16s[1])
	isEqual(t, &uint161, cfg.Uint16Ptrs[0])
	isEqual(t, &uint162, cfg.Uint16Ptrs[1])

	isEqual(t, uint321, cfg.Uint32)
	isEqual(t, &uint321, cfg.Uint32Ptr)
	isEqual(t, uint321, cfg.Uint32s[0])
	isEqual(t, uint322, cfg.Uint32s[1])
	isEqual(t, &uint321, cfg.Uint32Ptrs[0])
	isEqual(t, &uint322, cfg.Uint32Ptrs[1])

	isEqual(t, uint641, cfg.Uint64)
	isEqual(t, &uint641, cfg.Uint64Ptr)
	isEqual(t, uint641, cfg.Uint64s[0])
	isEqual(t, uint642, cfg.Uint64s[1])
	isEqual(t, &uint641, cfg.Uint64Ptrs[0])
	isEqual(t, &uint642, cfg.Uint64Ptrs[1])

	isEqual(t, float321, cfg.Float32)
	isEqual(t, &float321, cfg.Float32Ptr)
	isEqual(t, float321, cfg.Float32s[0])
	isEqual(t, float322, cfg.Float32s[1])
	isEqual(t, &float321, cfg.Float32Ptrs[0])
	isEqual(t, &float322, cfg.Float32Ptrs[1])

	isEqual(t, float641, cfg.Float64)
	isEqual(t, &float641, cfg.Float64Ptr)
	isEqual(t, float641, cfg.Float64s[0])
	isEqual(t, float642, cfg.Float64s[1])
	isEqual(t, &float641, cfg.Float64Ptrs[0])
	isEqual(t, &float642, cfg.Float64Ptrs[1])

	isEqual(t, duration1, cfg.Duration)
	isEqual(t, &duration1, cfg.DurationPtr)
	isEqual(t, duration1, cfg.Durations[0])
	isEqual(t, duration2, cfg.Durations[1])
	isEqual(t, &duration1, cfg.DurationPtrs[0])
	isEqual(t, &duration2, cfg.DurationPtrs[1])

	isEqual(t, unmarshaler1, cfg.Unmarshaler)
	isEqual(t, &unmarshaler1, cfg.UnmarshalerPtr)
	isEqual(t, unmarshaler1, cfg.Unmarshalers[0])
	isEqual(t, unmarshaler2, cfg.Unmarshalers[1])
	isEqual(t, &unmarshaler1, cfg.UnmarshalerPtrs[0])
	isEqual(t, &unmarshaler2, cfg.UnmarshalerPtrs[1])

	isEqual(t, url1, cfg.URL.String())
	isEqual(t, url1, cfg.URLPtr.String())
	isEqual(t, url1, cfg.URLs[0].String())
	isEqual(t, url2, cfg.URLs[1].String())
	isEqual(t, url1, cfg.URLPtrs[0].String())
	isEqual(t, url2, cfg.URLPtrs[1].String())

	isEqual(t, "postgres://localhost:5432/db", cfg.StringWithdefault)
	isEqual(t, nonDefinedStr, cfg.NonDefined.String)
	isEqual(t, nonDefinedStr, cfg.NestedNonDefined.NonDefined.String)

	isEqual(t, str1, cfg.CustomSeparator[0])
	isEqual(t, str2, cfg.CustomSeparator[1])

	isEqual(t, cfg.NotAnEnv, "")

	isEqual(t, cfg.unexported, "")
}

func TestSetEnvAndTagOptsChain(t *testing.T) {
	type config struct {
		Key1 string `mytag:"KEY1,required"`
		Key2 int    `mytag:"KEY2,required"`
	}
	envs := map[string]string{
		"KEY1": "VALUE1",
		"KEY2": "3",
	}

	cfg := config{}
	isNoErr(t, Parse(&cfg, Options{TagName: "mytag"}, Options{Environment: envs}))
	isEqual(t, "VALUE1", cfg.Key1)
	isEqual(t, 3, cfg.Key2)
}

func TestJSONTag(t *testing.T) {
	type config struct {
		Key1 string `json:"KEY1"`
		Key2 int    `json:"KEY2"`
	}

	t.Setenv("KEY1", "VALUE7")
	t.Setenv("KEY2", "5")

	cfg := config{}
	isNoErr(t, Parse(&cfg, Options{TagName: "json"}))
	isEqual(t, "VALUE7", cfg.Key1)
	isEqual(t, 5, cfg.Key2)
}

func TestParsesEnvInner(t *testing.T) {
	t.Setenv("innervar", "someinnervalue")
	t.Setenv("innernum", "8")
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
		unexported:  &InnerStruct{},
	}
	isNoErr(t, Parse(&cfg))
	isEqual(t, "someinnervalue", cfg.InnerStruct.Inner)
	isEqual(t, uint(8), cfg.InnerStruct.Number)
}

func TestParsesEnvInnerFails(t *testing.T) {
	type config struct {
		Foo struct {
			Number int `env:"NUMBER"`
		}
	}
	t.Setenv("NUMBER", "not-a-number")
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "Number" of type "int": strconv.ParseInt: parsing "not-a-number": invalid syntax`)
}

func TestParsesEnvInnerFailsMultipleErrors(t *testing.T) {
	type config struct {
		Foo struct {
			Name   string `env:"NAME,required"`
			Number int    `env:"NUMBER"`
			Bar    struct {
				Age int `env:"AGE,required"`
			}
		}
	}
	t.Setenv("NUMBER", "not-a-number")
	isErrorWithMessage(t, Parse(&config{}), `env: required environment variable "NAME" is not set; parse error on field "Number" of type "int": strconv.ParseInt: parsing "not-a-number": invalid syntax; required environment variable "AGE" is not set`)
}

func TestParsesEnvInnerNil(t *testing.T) {
	t.Setenv("innervar", "someinnervalue")
	cfg := ParentStruct{}
	isNoErr(t, Parse(&cfg))
}

func TestParsesEnvInnerInvalid(t *testing.T) {
	t.Setenv("innernum", "-547")
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
	}
	isErrorWithMessage(t, Parse(&cfg), `env: parse error on field "Number" of type "uint": strconv.ParseUint: parsing "-547": invalid syntax`)
}

func TestParsesEnvNested(t *testing.T) {
	t.Setenv("nestedvar", "somenestedvalue")
	var cfg ForNestedStruct
	isNoErr(t, Parse(&cfg))
	isEqual(t, "somenestedvalue", cfg.NestedVar)
}

func TestEmptyVars(t *testing.T) {
	cfg := Config{}
	isNoErr(t, Parse(&cfg))
	isEqual(t, "", cfg.String)
	isEqual(t, false, cfg.Bool)
	isEqual(t, 0, cfg.Int)
	isEqual(t, uint(0), cfg.Uint)
	isEqual(t, uint64(0), cfg.Uint64)
	isEqual(t, int64(0), cfg.Int64)
	isEqual(t, 0, len(cfg.Strings))
	isEqual(t, 0, len(cfg.CustomSeparator))
	isEqual(t, 0, len(cfg.Ints))
	isEqual(t, 0, len(cfg.Bools))
}

func TestPassAnInvalidPtr(t *testing.T) {
	var thisShouldBreak int
	isErrorWithMessage(t, Parse(&thisShouldBreak), "env: expected a pointer to a Struct")
}

func TestPassReference(t *testing.T) {
	cfg := Config{}
	isErrorWithMessage(t, Parse(cfg), "env: expected a pointer to a Struct")
}

func TestInvalidBool(t *testing.T) {
	t.Setenv("BOOL", "should-be-a-bool")
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Bool" of type "bool": strconv.ParseBool: parsing "should-be-a-bool": invalid syntax; parse error on field "BoolPtr" of type "*bool": strconv.ParseBool: parsing "should-be-a-bool": invalid syntax`)
}

func TestInvalidInt(t *testing.T) {
	t.Setenv("INT", "should-be-an-int")
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Int" of type "int": strconv.ParseInt: parsing "should-be-an-int": invalid syntax; parse error on field "IntPtr" of type "*int": strconv.ParseInt: parsing "should-be-an-int": invalid syntax`)
}

func TestInvalidUint(t *testing.T) {
	t.Setenv("UINT", "-44")
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Uint" of type "uint": strconv.ParseUint: parsing "-44": invalid syntax; parse error on field "UintPtr" of type "*uint": strconv.ParseUint: parsing "-44": invalid syntax`)
}

func TestInvalidFloat32(t *testing.T) {
	t.Setenv("FLOAT32", "AAA")
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Float32" of type "float32": strconv.ParseFloat: parsing "AAA": invalid syntax; parse error on field "Float32Ptr" of type "*float32": strconv.ParseFloat: parsing "AAA": invalid syntax`)
}

func TestInvalidFloat64(t *testing.T) {
	t.Setenv("FLOAT64", "AAA")
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Float64" of type "float64": strconv.ParseFloat: parsing "AAA": invalid syntax; parse error on field "Float64Ptr" of type "*float64": strconv.ParseFloat: parsing "AAA": invalid syntax`)
}

func TestInvalidUint64(t *testing.T) {
	t.Setenv("UINT64", "AAA")
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Uint64" of type "uint64": strconv.ParseUint: parsing "AAA": invalid syntax; parse error on field "Uint64Ptr" of type "*uint64": strconv.ParseUint: parsing "AAA": invalid syntax`)
}

func TestInvalidInt64(t *testing.T) {
	t.Setenv("INT64", "AAA")
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Int64" of type "int64": strconv.ParseInt: parsing "AAA": invalid syntax; parse error on field "Int64Ptr" of type "*int64": strconv.ParseInt: parsing "AAA": invalid syntax`)
}

func TestInvalidInt64Slice(t *testing.T) {
	t.Setenv("BADINTS", "A,2,3")
	type config struct {
		BadFloats []int64 `env:"BADINTS"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "BadFloats" of type "[]int64": strconv.ParseInt: parsing "A": invalid syntax`)
}

func TestInvalidUInt64Slice(t *testing.T) {
	t.Setenv("BADINTS", "A,2,3")
	type config struct {
		BadFloats []uint64 `env:"BADINTS"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "BadFloats" of type "[]uint64": strconv.ParseUint: parsing "A": invalid syntax`)
}

func TestInvalidFloat32Slice(t *testing.T) {
	t.Setenv("BADFLOATS", "A,2.0,3.0")
	type config struct {
		BadFloats []float32 `env:"BADFLOATS"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "BadFloats" of type "[]float32": strconv.ParseFloat: parsing "A": invalid syntax`)
}

func TestInvalidFloat64Slice(t *testing.T) {
	t.Setenv("BADFLOATS", "A,2.0,3.0")
	type config struct {
		BadFloats []float64 `env:"BADFLOATS"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "BadFloats" of type "[]float64": strconv.ParseFloat: parsing "A": invalid syntax`)
}

func TestInvalidBoolsSlice(t *testing.T) {
	t.Setenv("BADBOOLS", "t,f,TRUE,faaaalse")
	type config struct {
		BadBools []bool `env:"BADBOOLS"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "BadBools" of type "[]bool": strconv.ParseBool: parsing "faaaalse": invalid syntax`)
}

func TestInvalidDuration(t *testing.T) {
	t.Setenv("DURATION", "should-be-a-valid-duration")
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Duration" of type "time.Duration": unable to parse duration: time: invalid duration "should-be-a-valid-duration"; parse error on field "DurationPtr" of type "*time.Duration": unable to parse duration: time: invalid duration "should-be-a-valid-duration"`)
}

func TestInvalidDurations(t *testing.T) {
	t.Setenv("DURATIONS", "1s,contains-an-invalid-duration,3s")
	isErrorWithMessage(t, Parse(&Config{}), `env: parse error on field "Durations" of type "[]time.Duration": unable to parse duration: time: invalid duration "contains-an-invalid-duration"; parse error on field "DurationPtrs" of type "[]*time.Duration": unable to parse duration: time: invalid duration "contains-an-invalid-duration"`)
}

func TestParseStructWithoutEnvTag(t *testing.T) {
	cfg := Config{}
	isNoErr(t, Parse(&cfg))
	isEqual(t, cfg.NotAnEnv, "")
}

func TestParseStructWithInvalidFieldKind(t *testing.T) {
	type config struct {
		WontWorkByte byte `env:"BLAH"`
	}
	t.Setenv("BLAH", "a")
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "WontWorkByte" of type "uint8": strconv.ParseUint: parsing "a": invalid syntax`)
}

func TestUnsupportedSliceType(t *testing.T) {
	type config struct {
		WontWork []map[int]int `env:"WONTWORK"`
	}

	t.Setenv("WONTWORK", "1,2,3")
	isErrorWithMessage(t, Parse(&config{}), `env: no parser found for field "WontWork" of type "[]map[int]int"`)
}

func TestBadSeparator(t *testing.T) {
	type config struct {
		WontWork []int `env:"WONTWORK" envSeparator:":"`
	}

	t.Setenv("WONTWORK", "1,2,3,4")
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "WontWork" of type "[]int": strconv.ParseInt: parsing "1,2,3,4": invalid syntax`)
}

func TestNoErrorRequiredSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required"`
	}

	cfg := &config{}

	t.Setenv("IS_REQUIRED", "")
	isNoErr(t, Parse(cfg))
	isEqual(t, "", cfg.IsRequired)
}

func TestHook(t *testing.T) {
	type config struct {
		Something string `env:"SOMETHING" envDefault:"important"`
		Another   string `env:"ANOTHER"`
	}

	cfg := &config{}
	t.Setenv("ANOTHER", "1")

	type onSetArgs struct {
		tag       string
		key       interface{}
		isDefault bool
	}

	var onSetCalled []onSetArgs

	isNoErr(t, Parse(cfg, Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			onSetCalled = append(onSetCalled, onSetArgs{tag, value, isDefault})
		},
	}))
	isEqual(t, "important", cfg.Something)
	isEqual(t, "1", cfg.Another)
	isEqual(t, 2, len(onSetCalled))
	isEqual(t, onSetArgs{"SOMETHING", "important", true}, onSetCalled[0])
	isEqual(t, onSetArgs{"ANOTHER", "1", false}, onSetCalled[1])
}

func TestErrorRequiredWithDefault(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required" envDefault:"important"`
	}

	cfg := &config{}

	t.Setenv("IS_REQUIRED", "")
	isNoErr(t, Parse(cfg))
	isEqual(t, "", cfg.IsRequired)
}

func TestErrorRequiredNotSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: required environment variable "IS_REQUIRED" is not set`)
}

func TestNoErrorNotEmptySet(t *testing.T) {
	t.Setenv("IS_REQUIRED", "1")
	type config struct {
		IsRequired string `env:"IS_REQUIRED,notEmpty"`
	}
	isNoErr(t, Parse(&config{}))
}

func TestNoErrorRequiredAndNotEmptySet(t *testing.T) {
	t.Setenv("IS_REQUIRED", "1")
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required,notEmpty"`
	}
	isNoErr(t, Parse(&config{}))
}

func TestErrorNotEmptySet(t *testing.T) {
	t.Setenv("IS_REQUIRED", "")
	type config struct {
		IsRequired string `env:"IS_REQUIRED,notEmpty"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: environment variable "IS_REQUIRED" should not be empty`)
}

func TestErrorRequiredAndNotEmptySet(t *testing.T) {
	t.Setenv("IS_REQUIRED", "")
	type config struct {
		IsRequired string `env:"IS_REQUIRED,notEmpty,required"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: environment variable "IS_REQUIRED" should not be empty`)
}

func TestErrorRequiredNotSetWithDefault(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required" envDefault:"important"`
	}

	cfg := &config{}
	isNoErr(t, Parse(cfg))
	isEqual(t, "important", cfg.IsRequired)
}

func TestParseExpandOption(t *testing.T) {
	type config struct {
		Host        string `env:"HOST" envDefault:"localhost"`
		Port        int    `env:"PORT" envDefault:"3000" envExpand:"True"`
		SecretKey   string `env:"SECRET_KEY" envExpand:"True"`
		ExpandKey   string `env:"EXPAND_KEY"`
		CompoundKey string `env:"HOST_PORT" envDefault:"${HOST}:${PORT}" envExpand:"True"`
		Default     string `env:"DEFAULT" envDefault:"def1"  envExpand:"True"`
	}

	t.Setenv("HOST", "localhost")
	t.Setenv("PORT", "3000")
	t.Setenv("EXPAND_KEY", "qwerty12345")
	t.Setenv("SECRET_KEY", "${EXPAND_KEY}")

	cfg := config{}
	err := Parse(&cfg)

	isNoErr(t, err)
	isEqual(t, "localhost", cfg.Host)
	isEqual(t, 3000, cfg.Port)
	isEqual(t, "qwerty12345", cfg.SecretKey)
	isEqual(t, "qwerty12345", cfg.ExpandKey)
	isEqual(t, "localhost:3000", cfg.CompoundKey)
	isEqual(t, "def1", cfg.Default)
}

func TestParseUnsetRequireOptions(t *testing.T) {
	type config struct {
		Password string `env:"PASSWORD,unset,required"`
	}
	cfg := config{}

	isErrorWithMessage(t, Parse(&cfg), `env: required environment variable "PASSWORD" is not set`)
	t.Setenv("PASSWORD", "superSecret")
	isNoErr(t, Parse(&cfg))

	isEqual(t, "superSecret", cfg.Password)
	unset, exists := os.LookupEnv("PASSWORD")
	isEqual(t, "", unset)
	isEqual(t, false, exists)
}

func TestCustomParser(t *testing.T) {
	type foo struct {
		name string
	}

	type bar struct {
		Name string `env:"OTHER_CUSTOM"`
		Foo  *foo   `env:"BLAH_CUSTOM"`
	}

	type config struct {
		Var   foo  `env:"VAR_CUSTOM"`
		Foo   *foo `env:"BLAH_CUSTOM"`
		Other *bar
	}

	t.Setenv("VAR_CUSTOM", "test")
	t.Setenv("OTHER_CUSTOM", "test2")
	t.Setenv("BLAH_CUSTOM", "test3")

	runtest := func(t *testing.T) {
		t.Helper()
		cfg := &config{
			Other: &bar{},
		}
		err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
			reflect.TypeOf(foo{}): func(v string) (interface{}, error) {
				return foo{name: v}, nil
			},
		})

		isNoErr(t, err)
		isEqual(t, cfg.Var.name, "test")
		isEqual(t, cfg.Foo.name, "test3")
		isEqual(t, cfg.Other.Name, "test2")
		isEqual(t, cfg.Other.Foo.name, "test3")
	}

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("%d", i), runtest)
	}
}

func TestIssue226(t *testing.T) {
	type config struct {
		Inner struct {
			Abc []byte `env:"ABC" envDefault:"asdasd"`
			Def []byte `env:"DEF" envDefault:"a"`
		}
		Hij []byte `env:"HIJ"`
		Lmn []byte `env:"LMN"`
	}

	t.Setenv("HIJ", "a")
	t.Setenv("LMN", "b")

	cfg := &config{}
	isNoErr(t, ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf([]byte{0}): func(v string) (interface{}, error) {
			if v == "a" {
				return []byte("nope"), nil
			}
			return []byte(v), nil
		},
	}))
	isEqual(t, cfg.Inner.Abc, []byte("asdasd"))
	isEqual(t, cfg.Inner.Def, []byte("nope"))
	isEqual(t, cfg.Hij, []byte("nope"))
	isEqual(t, cfg.Lmn, []byte("b"))
}

func TestParseWithFuncsNoPtr(t *testing.T) {
	type foo struct{}
	isErrorWithMessage(t, ParseWithFuncs(foo{}, nil), "env: expected a pointer to a Struct")
}

func TestParseWithFuncsInvalidType(t *testing.T) {
	var c int
	isErrorWithMessage(t, ParseWithFuncs(&c, nil), "env: expected a pointer to a Struct")
}

func TestCustomParserError(t *testing.T) {
	type foo struct {
		name string
	}

	customParserFunc := func(v string) (interface{}, error) {
		return nil, errors.New("something broke")
	}

	t.Run("single", func(t *testing.T) {
		type config struct {
			Var foo `env:"VAR"`
		}

		t.Setenv("VAR", "single")
		cfg := &config{}
		err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
			reflect.TypeOf(foo{}): customParserFunc,
		})

		isEqual(t, cfg.Var.name, "")
		isErrorWithMessage(t, err, `env: parse error on field "Var" of type "env.foo": something broke`)
	})

	t.Run("slice", func(t *testing.T) {
		type config struct {
			Var []foo `env:"VAR2"`
		}
		t.Setenv("VAR2", "slice,slace")

		cfg := &config{}
		err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
			reflect.TypeOf(foo{}): customParserFunc,
		})

		isEqual(t, cfg.Var, nil)
		isErrorWithMessage(t, err, `env: parse error on field "Var" of type "[]env.foo": something broke`)
	})
}

func TestCustomParserBasicType(t *testing.T) {
	type ConstT int32

	type config struct {
		Const ConstT `env:"CONST_"`
	}

	exp := ConstT(123)
	t.Setenv("CONST_", fmt.Sprintf("%d", exp))

	customParserFunc := func(v string) (interface{}, error) {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		r := ConstT(i)
		return r, nil
	}

	cfg := &config{}
	err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(ConstT(0)): customParserFunc,
	})

	isNoErr(t, err)
	isEqual(t, exp, cfg.Const)
}

func TestCustomParserUint64Alias(t *testing.T) {
	type T uint64

	var one T = 1

	type config struct {
		Val T `env:"" envDefault:"1x"`
	}

	parserCalled := false

	tParser := func(value string) (interface{}, error) {
		parserCalled = true
		trimmed := strings.TrimSuffix(value, "x")
		i, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, err
		}
		return T(i), nil
	}

	cfg := config{}

	err := ParseWithFuncs(&cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(one): tParser,
	})

	isTrue(t, parserCalled)
	isNoErr(t, err)
	isEqual(t, T(1), cfg.Val)
}

func TestTypeCustomParserBasicInvalid(t *testing.T) {
	type ConstT int32

	type config struct {
		Const ConstT `env:"CONST_"`
	}

	t.Setenv("CONST_", "foobar")

	customParserFunc := func(_ string) (interface{}, error) {
		return nil, errors.New("random error")
	}

	cfg := &config{}
	err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(ConstT(0)): customParserFunc,
	})

	isEqual(t, cfg.Const, ConstT(0))
	isErrorWithMessage(t, err, `env: parse error on field "Const" of type "env.ConstT": random error`)
}

func TestCustomParserNotCalledForNonAlias(t *testing.T) {
	type T uint64
	type U uint64

	type config struct {
		Val   uint64 `env:"" envDefault:"33"`
		Other U      `env:"OTHER_NAME" envDefault:"44"`
	}

	tParserCalled := false

	tParser := func(value string) (interface{}, error) {
		tParserCalled = true
		return T(99), nil
	}

	cfg := config{}

	err := ParseWithFuncs(&cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(T(0)): tParser,
	})

	isFalse(t, tParserCalled)
	isNoErr(t, err)
	isEqual(t, uint64(33), cfg.Val)
	isEqual(t, U(44), cfg.Other)
}

func TestCustomParserBasicUnsupported(t *testing.T) {
	type ConstT struct {
		A int
	}

	type config struct {
		Const ConstT `env:"CONST_"`
	}

	t.Setenv("CONST_", "42")

	cfg := &config{}
	err := Parse(cfg)

	isEqual(t, cfg.Const, ConstT{0})
	isErrorWithMessage(t, err, `env: no parser found for field "Const" of type "env.ConstT"`)
}

func TestUnsupportedStructType(t *testing.T) {
	type config struct {
		Foo http.Client `env:"FOO"`
	}
	t.Setenv("FOO", "foo")
	isErrorWithMessage(t, Parse(&config{}), `env: no parser found for field "Foo" of type "http.Client"`)
}

func TestEmptyOption(t *testing.T) {
	type config struct {
		Var string `env:"VAR,"`
	}

	cfg := &config{}

	t.Setenv("VAR", "")
	isNoErr(t, Parse(cfg))
	isEqual(t, "", cfg.Var)
}

func TestErrorOptionNotRecognized(t *testing.T) {
	type config struct {
		Var string `env:"VAR,not_supported!"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: tag option "not_supported!" not supported`)
}

func TestTextUnmarshalerError(t *testing.T) {
	type config struct {
		Unmarshaler unmarshaler `env:"UNMARSHALER"`
	}
	t.Setenv("UNMARSHALER", "invalid")
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "Unmarshaler" of type "env.unmarshaler": time: invalid duration "invalid"`)
}

func TestTextUnmarshalersError(t *testing.T) {
	type config struct {
		Unmarshalers []unmarshaler `env:"UNMARSHALERS"`
	}
	t.Setenv("UNMARSHALERS", "1s,invalid")
	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "Unmarshalers" of type "[]env.unmarshaler": time: invalid duration "invalid"`)
}

func TestParseURL(t *testing.T) {
	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL" envDefault:"https://google.com"`
	}
	var cfg config
	isNoErr(t, Parse(&cfg))
	isEqual(t, "https://google.com", cfg.ExampleURL.String())
}

func TestParseInvalidURL(t *testing.T) {
	type config struct {
		ExampleURL url.URL `env:"EXAMPLE_URL_2"`
	}
	t.Setenv("EXAMPLE_URL_2", "nope://s s/")

	isErrorWithMessage(t, Parse(&config{}), `env: parse error on field "ExampleURL" of type "url.URL": unable to parse URL: parse "nope://s s/": invalid character " " in host name`)
}

func ExampleParse() {
	type inner struct {
		Foo string `env:"FOO" envDefault:"foobar"`
	}
	type config struct {
		Home         string `env:"HOME,required"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
		Inner        inner
	}
	os.Setenv("HOME", "/tmp/fakehome")
	var cfg config
	if err := Parse(&cfg); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Home:/tmp/fakehome Port:3000 IsProduction:false Inner:{Foo:foobar}}
}

func ExampleParse_onSet() {
	type config struct {
		Home         string `env:"HOME,required"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	var cfg config
	if err := Parse(&cfg, Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			fmt.Printf("Set %s to %v (default? %v)\n", tag, value, isDefault)
		},
	}); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output: Set HOME to /tmp/fakehome (default? false)
	// Set PORT to 3000 (default? true)
	// Set PRODUCTION to  (default? false)
	// {Home:/tmp/fakehome Port:3000 IsProduction:false}
}

func ExampleParse_defaults() {
	type config struct {
		A string `env:"FOO" envDefault:"foo"`
		B string `env:"FOO"`
	}

	// env FOO is not set

	cfg := config{
		A: "A",
		B: "B",
	}
	if err := Parse(&cfg); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {A:foo B:B}
}

func TestIgnoresUnexported(t *testing.T) {
	type unexportedConfig struct {
		home  string `env:"HOME"`
		Home2 string `env:"HOME"`
	}
	cfg := unexportedConfig{}

	t.Setenv("HOME", "/tmp/fakehome")
	isNoErr(t, Parse(&cfg))
	isEqual(t, cfg.home, "")
	isEqual(t, "/tmp/fakehome", cfg.Home2)
}

type LogLevel int8

func (l *LogLevel) UnmarshalText(text []byte) error {
	txt := string(text)
	switch txt {
	case "debug":
		*l = DebugLevel
	case "info":
		*l = InfoLevel
	default:
		return fmt.Errorf("unknown level: %q", txt)
	}

	return nil
}

const (
	DebugLevel LogLevel = iota - 1
	InfoLevel
)

func TestPrecedenceUnmarshalText(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_LEVELS", "debug,info")

	type config struct {
		LogLevel  LogLevel   `env:"LOG_LEVEL"`
		LogLevels []LogLevel `env:"LOG_LEVELS"`
	}
	var cfg config

	isNoErr(t, Parse(&cfg))
	isEqual(t, DebugLevel, cfg.LogLevel)
	isEqual(t, []LogLevel{DebugLevel, InfoLevel}, cfg.LogLevels)
}

func ExampleParseWithFuncs() {
	type thing struct {
		desc string
	}

	type conf struct {
		Thing thing `env:"THING"`
	}

	os.Setenv("THING", "my thing")

	c := conf{}

	err := ParseWithFuncs(&c, map[reflect.Type]ParserFunc{
		reflect.TypeOf(thing{}): func(v string) (interface{}, error) {
			return thing{desc: v}, nil
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c.Thing.desc)
	// Output:
	// my thing
}

func TestFile(t *testing.T) {
	type config struct {
		SecretKey string `env:"SECRET_KEY,file"`
	}

	dir := t.TempDir()
	file := filepath.Join(dir, "sec_key")
	isNoErr(t, os.WriteFile(file, []byte("secret"), 0o660))

	t.Setenv("SECRET_KEY", file)

	cfg := config{}
	isNoErr(t, Parse(&cfg))
	isEqual(t, "secret", cfg.SecretKey)
}

func TestFileNoParam(t *testing.T) {
	type config struct {
		SecretKey string `env:"SECRET_KEY,file"`
	}

	cfg := config{}
	isNoErr(t, Parse(&cfg))
}

func TestFileNoParamRequired(t *testing.T) {
	type config struct {
		SecretKey string `env:"SECRET_KEY,file,required"`
	}
	isErrorWithMessage(t, Parse(&config{}), `env: required environment variable "SECRET_KEY" is not set`)
}

func TestFileBadFile(t *testing.T) {
	type config struct {
		SecretKey string `env:"SECRET_KEY,file"`
	}

	filename := "not-a-real-file"
	t.Setenv("SECRET_KEY", filename)

	oserr := "no such file or directory"
	if runtime.GOOS == "windows" {
		oserr = "The system cannot find the file specified."
	}
	isErrorWithMessage(t, Parse(&config{}), fmt.Sprintf(`env: could not load content of file "%s" from variable SECRET_KEY: open %s: %s`, filename, filename, oserr))
}

func TestFileWithDefault(t *testing.T) {
	type config struct {
		SecretKey string `env:"SECRET_KEY,file" envDefault:"${FILE}" envExpand:"true"`
	}

	dir := t.TempDir()
	file := filepath.Join(dir, "sec_key")
	isNoErr(t, os.WriteFile(file, []byte("secret"), 0o660))

	t.Setenv("FILE", file)

	cfg := config{}
	isNoErr(t, Parse(&cfg))
	isEqual(t, "secret", cfg.SecretKey)
}

func TestCustomSliceType(t *testing.T) {
	type customslice []byte

	type config struct {
		SecretKey customslice `env:"SECRET_KEY"`
	}

	t.Setenv("SECRET_KEY", "somesecretkey")

	var cfg config
	isNoErr(t, ParseWithFuncs(&cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(customslice{}): func(value string) (interface{}, error) {
			return customslice(value), nil
		},
	}))
}

type MyTime time.Time

func (t *MyTime) UnmarshalText(text []byte) error {
	tt, err := time.Parse("2006-01-02", string(text))
	*t = MyTime(tt)
	return err
}

func TestCustomTimeParser(t *testing.T) {
	type config struct {
		SomeTime MyTime `env:"SOME_TIME"`
	}

	t.Setenv("SOME_TIME", "2021-05-06")

	var cfg config
	isNoErr(t, Parse(&cfg))
	isEqual(t, 2021, time.Time(cfg.SomeTime).Year())
	isEqual(t, time.Month(5), time.Time(cfg.SomeTime).Month())
	isEqual(t, 6, time.Time(cfg.SomeTime).Day())
}

func TestRequiredIfNoDefOption(t *testing.T) {
	type Tree struct {
		Fruit string `env:"FRUIT"`
	}
	type config struct {
		Name  string `env:"NAME"`
		Genre string `env:"GENRE" envDefault:"Unknown"`
		Tree
	}
	var cfg config

	t.Run("missing", func(t *testing.T) {
		isErrorWithMessage(t, Parse(&cfg, Options{RequiredIfNoDef: true}), `env: required environment variable "NAME" is not set; required environment variable "FRUIT" is not set`)
		t.Setenv("NAME", "John")
		isErrorWithMessage(t, Parse(&cfg, Options{RequiredIfNoDef: true}), `env: required environment variable "FRUIT" is not set`)
	})

	t.Run("all set", func(t *testing.T) {
		t.Setenv("NAME", "John")
		t.Setenv("FRUIT", "Apple")

		// should not trigger an error for the missing 'GENRE' env because it has a default value.
		isNoErr(t, Parse(&cfg, Options{RequiredIfNoDef: true}))
	})
}

func TestRequiredIfNoDefNested(t *testing.T) {
	type Server struct {
		Host string `env:"HOST"`
		Port uint16 `env:"PORT"`
	}
	type API struct {
		Server
		Token string `env:"TOKEN"`
	}
	type config struct {
		API API `envPrefix:"SERVER_"`
	}
	var cfg config

	t.Run("missing", func(t *testing.T) {
		t.Setenv("SERVER_HOST", "https://google.com")
		t.Setenv("SERVER_TOKEN", "0xdeadfood")

		isErrorWithMessage(t, Parse(&cfg, Options{RequiredIfNoDef: true}), `env: required environment variable "SERVER_PORT" is not set`)
	})

	t.Run("all set", func(t *testing.T) {
		t.Setenv("SERVER_HOST", "https://google.com")
		t.Setenv("SERVER_PORT", "443")
		t.Setenv("SERVER_TOKEN", "0xdeadfood")

		isNoErr(t, Parse(&cfg, Options{RequiredIfNoDef: true}))
	})
}

func TestPrefix(t *testing.T) {
	type Config struct {
		Home string `env:"HOME"`
	}
	type ComplexConfig struct {
		Foo   Config `envPrefix:"FOO_"`
		Bar   Config `envPrefix:"BAR_"`
		Clean Config
	}
	cfg := ComplexConfig{}
	isNoErr(t, Parse(&cfg, Options{Environment: map[string]string{"FOO_HOME": "/foo", "BAR_HOME": "/bar", "HOME": "/clean"}}))
	isEqual(t, "/foo", cfg.Foo.Home)
	isEqual(t, "/bar", cfg.Bar.Home)
	isEqual(t, "/clean", cfg.Clean.Home)
}

func TestPrefixPointers(t *testing.T) {
	type Test struct {
		Str string `env:"TEST"`
	}
	type ComplexConfig struct {
		Foo   *Test `envPrefix:"FOO_"`
		Bar   *Test `envPrefix:"BAR_"`
		Clean *Test
	}

	cfg := ComplexConfig{
		Foo:   &Test{},
		Bar:   &Test{},
		Clean: &Test{},
	}
	isNoErr(t, Parse(&cfg, Options{Environment: map[string]string{"FOO_TEST": "kek", "BAR_TEST": "lel", "TEST": "clean"}}))
	isEqual(t, "kek", cfg.Foo.Str)
	isEqual(t, "lel", cfg.Bar.Str)
	isEqual(t, "clean", cfg.Clean.Str)
}

func TestNestedPrefixPointer(t *testing.T) {
	type ComplexConfig struct {
		Foo struct {
			Str string `env:"STR"`
		} `envPrefix:"FOO_"`
	}
	cfg := ComplexConfig{}
	isNoErr(t, Parse(&cfg, Options{Environment: map[string]string{"FOO_STR": "foo_str"}}))
	isEqual(t, "foo_str", cfg.Foo.Str)

	type ComplexConfig2 struct {
		Foo struct {
			Bar struct {
				Str string `env:"STR"`
			} `envPrefix:"BAR_"`
			Bar2 string `env:"BAR2"`
		} `envPrefix:"FOO_"`
	}
	cfg2 := ComplexConfig2{}
	isNoErr(t, Parse(&cfg2, Options{Environment: map[string]string{"FOO_BAR_STR": "kek", "FOO_BAR2": "lel"}}))
	isEqual(t, "lel", cfg2.Foo.Bar2)
	isEqual(t, "kek", cfg2.Foo.Bar.Str)
}

func TestComplePrefix(t *testing.T) {
	type Config struct {
		Home string `env:"HOME"`
	}
	type ComplexConfig struct {
		Foo   Config `envPrefix:"FOO_"`
		Clean Config
		Bar   Config `envPrefix:"BAR_"`
		Blah  string `env:"BLAH"`
	}
	cfg := ComplexConfig{}
	isNoErr(t, Parse(&cfg, Options{
		Prefix: "T_",
		Environment: map[string]string{
			"T_FOO_HOME": "/foo",
			"T_BAR_HOME": "/bar",
			"T_BLAH":     "blahhh",
			"T_HOME":     "/clean",
		},
	}))
	isEqual(t, "/foo", cfg.Foo.Home)
	isEqual(t, "/bar", cfg.Bar.Home)
	isEqual(t, "/clean", cfg.Clean.Home)
	isEqual(t, "blahhh", cfg.Blah)
}

func TestNonStructPtrValues(t *testing.T) {
	type Foo struct {
		FltPtr *float64 `env:"FLT_PRT"`
	}

	type ComplexConfig struct {
		StrPtr *string `env:"STR_PTR"`
		Foo    Foo     `env:"FOO_"`
	}

	cfg1 := ComplexConfig{}

	isNoErr(t, Parse(&cfg1))
	isEqual(t, nil, cfg1.StrPtr)
	isEqual(t, nil, cfg1.Foo.FltPtr)

	strPtr := "str_ptr"
	fltPtr := 3.16
	cfg2 := ComplexConfig{
		StrPtr: &strPtr,
		Foo: Foo{
			FltPtr: &fltPtr,
		},
	}

	t.Setenv("STR_PTR", "env_str_ptr")
	t.Setenv("FLT_PRT", "5.16")

	isNoErr(t, Parse(&cfg2))
	isEqual(t, "env_str_ptr", *cfg2.StrPtr)
	isEqual(t, 5.16, *cfg2.Foo.FltPtr)

	var strPtrNill *string
	var fltPtrNill *float64
	cfg3 := ComplexConfig{
		StrPtr: strPtrNill,
		Foo: Foo{
			FltPtr: fltPtrNill,
		},
	}

	isNoErr(t, Parse(&cfg3))
	isEqual(t, "env_str_ptr", *cfg3.StrPtr)
	isEqual(t, 5.16, *cfg3.Foo.FltPtr)
}

func isTrue(tb testing.TB, b bool) {
	tb.Helper()

	if !b {
		tb.Fatalf("expected true, got false")
	}
}

func isFalse(tb testing.TB, b bool) {
	tb.Helper()

	if b {
		tb.Fatalf("expected false, got true")
	}
}

func isErrorWithMessage(tb testing.TB, err error, msg string) {
	tb.Helper()

	if err == nil {
		tb.Fatalf("expected error, got nil")
	}

	if msg != err.Error() {
		tb.Fatalf("expected error message %q, got %q", msg, err.Error())
	}
}

func isNoErr(tb testing.TB, err error) {
	tb.Helper()

	if err != nil {
		tb.Fatalf("unexpected error: %v", err)
	}
}

func isEqual(tb testing.TB, a, b interface{}) {
	tb.Helper()

	if areEqual(a, b) {
		return
	}

	tb.Fatalf("expected %#v (type %T) == %#v (type %T)", a, a, b, b)
}

// copied from https://github.com/matryer/is
func areEqual(a, b interface{}) bool {
	if isNil(a) && isNil(b) {
		return true
	}
	if isNil(a) || isNil(b) {
		return false
	}
	if reflect.DeepEqual(a, b) {
		return true
	}
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)
	return aValue == bValue
}

// copied from https://github.com/matryer/is
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}
