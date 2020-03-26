package faker_test

import (
	"fmt"

	"github.com/togglhire/faker/v3"
)

type CustomString string

// You can set length for your random strings also set boundary for your integers.
func Example_withTagsUse() {

	// SomeStruct ...
	type SomeStruct struct {
		Inta  int   `faker:"use=10"`
		Int8  int8  `faker:"use=11"`
		Int16 int16 `faker:"use=12"`
		Int32 int32 `faker:"use=13"`
		Int64 int64 `faker:"use=14"`

		UInta  uint   `faker:"use=15"`
		UInt8  uint8  `faker:"use=16"`
		UInt16 uint16 `faker:"use=17"`
		UInt32 uint32 `faker:"use=18"`
		UInt64 uint64 `faker:"use=19"`

		Float32 float32 `faker:"use=20.1"`
		Float64 float64 `faker:"use=20.2"`

		String       string `faker:"use=string"`
		CustomString string `faker:"use=custom string"`
	}

	a := SomeStruct{}
	faker.FakeData(&a)
	fmt.Printf("%+v", a)

	// Result:
	/*
		{
			Inta:10
			Int8:11
			Int16:12
			Int32:13
			Int64:14

			UInta:15
			UInt8:16
			UInt16:17
			UInt32:18
			UInt64:19

			Float32:20.1
			Float64:20.2

			String:string
			CustomString:custom string
		}
	*/
}
