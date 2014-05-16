// Package autoflags provides a convenient way of exposing fields of struct as
// command line flags. Exposed fields should have special tag attached:
//
//	var config = struct {
//		Name    string `flag:"name,name of user"`
//		Age     uint   `flag:"age"`
//		Married bool   // this won't be exposed
//	}{
//		Name: "John Doe", // default values
//		Age:  34,
//	}
//
// After declaring your flags and their default values as above, just register
// flags with flag package and call flag.Parse() as usually:
//
// 	if err := autoflags.Define(&config) ; err != nil {
// 		log.Fatal(err)
// 	}
// 	flag.Parse()
//
// Now config struct has its fields populated from command line flags.
package autoflags

import (
	"errors"
	"flag"
	"reflect"
	"strings"
	"time"
)

var (
	// ErrPointerWanted is returned when passed argument is not a pointer
	ErrPointerWanted = errors.New("pointer expected")
	// ErrInvalidArgument is returned when passed argument is nil pointer or
	// pointer to a non-struct value
	ErrInvalidArgument = errors.New("non-nil pointer to struct expected")
)

// Define takes pointer to struct and declares flags for its flag-tagged fields.
// Valid tags have the following form: `flag:"flagname"` or
// `flag:"flagname,usage string"`.
func Define(config interface{}) error {
	st := reflect.ValueOf(config)
	if st.Kind() != reflect.Ptr {
		return ErrPointerWanted
	}
	st = reflect.Indirect(st)
	if !st.IsValid() || st.Type().Kind() != reflect.Struct {
		return ErrInvalidArgument
	}
	for i := 0; i < st.NumField(); i++ {
		val := st.Field(i)
		if !val.CanAddr() {
			continue
		}
		typ := st.Type().Field(i)
		var name, usage string
		tag := typ.Tag.Get("flag")
		if len(tag) == 0 {
			continue
		}
		flagData := strings.SplitN(tag, ",", 2)
		switch len(flagData) {
		case 1:
			name = flagData[0]
		case 2:
			name, usage = flagData[0], flagData[1]
		}
		addr := val.Addr()
		switch d := val.Interface().(type) {
		case int:
			flag.IntVar(addr.Interface().(*int), name, d, usage)
		case int64:
			flag.Int64Var(addr.Interface().(*int64), name, d, usage)
		case uint:
			flag.UintVar(addr.Interface().(*uint), name, d, usage)
		case uint64:
			flag.Uint64Var(addr.Interface().(*uint64), name, d, usage)
		case float64:
			flag.Float64Var(addr.Interface().(*float64), name, d, usage)
		case bool:
			flag.BoolVar(addr.Interface().(*bool), name, d, usage)
		case string:
			flag.StringVar(addr.Interface().(*string), name, d, usage)
		case time.Duration:
			flag.DurationVar(addr.Interface().(*time.Duration), name, d, usage)
		}
	}
	return nil
}
