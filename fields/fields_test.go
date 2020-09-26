package fields

import "fmt"

var someFields Fields

func ExampleFields_forEach() {
	err := someFields.ForEach(func(k string, v interface{}) error {
		fmt.Printf("%s=%+v\n", k, v)
		return nil
	})

	if err != nil {
		panic(fmt.Errorf("doh!: %w", err))
	}
}

func ExampleFields_get() {
	v := someFields.Get("foo")

	fmt.Printf("foo=%+v\n", v)
}
