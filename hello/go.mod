module example.com/hello

go 1.21.3

replace example.com/greetings => ../greetings

require (
	example.com/greetings v0.0.0-00010101000000-000000000000
	golang.org/x/example/hello v0.0.0-20231013143937-1d6d2400d402
)
