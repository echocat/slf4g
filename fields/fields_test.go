package fields

import "fmt"

func ExampleFields_forEach() {
	someFields := With("bar", 2).
		With("foo", 1)

	err := someFields.ForEach(func(k string, v interface{}) error {
		fmt.Printf("%s=%+v\n", k, v)
		return nil
	})

	if err != nil {
		panic(fmt.Errorf("doh!: %w", err))
	}

	// Output:
	// foo=1
	// bar=2
}

func ExampleFields_get() {
	someFields := With("bar", 2).
		With("foo", 1)

	v, _ := someFields.Get("foo")

	fmt.Printf("foo=%+v\n", v)

	// Output:
	// foo=1
}

func Example() {
	f := With("bar", 2).
		With("foo", "1").
		Withf("message", "something happened in module %s", "abc")

	err := f.ForEach(func(k string, v interface{}) error {
		fmt.Printf("%s=%+v\n", k, v)
		return nil
	})

	if err != nil {
		panic(err)
	}

	// Output:
	// message=something happened in module abc
	// foo=1
	// bar=2
}
