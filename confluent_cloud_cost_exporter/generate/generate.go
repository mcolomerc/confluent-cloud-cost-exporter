package generate

// Use a go:generate directive to build the Go structs for `example.avsc`
// These files are used for all of the example projects
// Source files will be in a package called `example/avro`

//go:generate $GOPATH/bin/gogen-avro ../pkg/avro ./avro/cost.avsc
