package autoflags

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func ResetForTesting(usage func()) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.Usage = usage
}

func TestDefineErrPointerWanted(t *testing.T) {
	ResetForTesting(nil)
	if err := Define(1); err != ErrPointerWanted {
		t.Fatalf("should fail with error ErrPointerWanted, got %q", err)
	}
}

func TestDefineErrInvalidArgument(t *testing.T) {
	ResetForTesting(nil)
	var testConfig *struct{}
	if err := Define(testConfig); err != ErrInvalidArgument {
		t.Fatalf("should fail with error ErrInvalidArgument, got %q", err)
	}
}

func TestDefineParseEmpty(t *testing.T) {
	ResetForTesting(nil)
	reference := config{
		String: "foo",
		Int:    42,
	}
	conf := reference
	if err := Define(&conf); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := flag.CommandLine.Parse([]string{}); err != nil {
		t.Fatal("parsing failed:", err)
	}
	if !reflect.DeepEqual(reference, conf) {
		t.Fatalf("result differs after parsing empty arguments; "+
			"want: %+v, got %+v", reference, conf)
	}
}

func TestDefineParse(t *testing.T) {
	ResetForTesting(nil)
	reference := configBig{
		String:   "whales",
		Int:      42,
		Int64:    100 << 30,
		Uint:     7,
		Uint64:   24,
		Float64:  1.55,
		Bool:     true,
		Duration: 15 * time.Minute,
	}
	conf := configBig{}
	if err := Define(&conf); err != nil {
		t.Fatal("unexpected error:", err)
	}
	args := []string{
		"-string", "whales", "-int", "42",
		"-int64", "107374182400", "-uint", "7",
		"-uint64", "24", "-float64", "1.55", "-bool",
		"-duration", "15m",
	}
	if err := flag.CommandLine.Parse(args); err != nil {
		t.Fatal("parsing failed:", err)
	}
	if !reflect.DeepEqual(reference, conf) {
		t.Fatalf("result differs after parsing arguments; "+
			"want: %+v, got %+v", reference, conf)
	}
}

func TestDefineFlagSetErrInvalidFlagSet(t *testing.T) {
	if err := DefineFlagSet(nil, &struct{}{}); err != ErrInvalidFlagSet {
		t.Fatalf("should fail with error ErrInvalidFlagSet, got %q", err)
	}
}

func ExampleDefineFlagSet() {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var config = struct {
		Name    string `flag:"name,name of user"`
		Age     uint   `flag:"age"`
		Married bool   // this won't be exposed
	}{
		Name: "John Doe", // default values
		Age:  34,
	}
	args := []string{"-name", "Jane Roe", "-age", "29"}

	fmt.Printf("before:\n%+v\n", config)
	fmt.Println("define err:", DefineFlagSet(fs, &config))
	fmt.Println("parse err:", fs.Parse(args))
	fmt.Printf("after:\n%+v\n", config)
	// Output:
	//
	// before:
	// {Name:John Doe Age:34 Married:false}
	// define err: <nil>
	// parse err: <nil>
	// after:
	// {Name:Jane Roe Age:29 Married:false}
}

type config struct {
	String string `flag:"name"`
	Int    int    `flag:"num,integer number"`
}

type configBig struct {
	String   string        `flag:"string,string flag example"`
	Int      int           `flag:"int,int flag example"`
	Int64    int64         `flag:"int64,int64 flag example"`
	Uint     uint          `flag:"uint,uint flag example"`
	Uint64   uint64        `flag:"uint64"`
	Float64  float64       `flag:"float64"`
	Bool     bool          `flag:"bool"`
	Duration time.Duration `flag:"duration"`

	NonAddressable unsafe.Pointer `flag:"nil"` // non-addressable
	Invalid        bool           `flag:""`    // invalid flag definition
	NonExposed     int            // does not have flag attached
}
