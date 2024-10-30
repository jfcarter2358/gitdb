package main

import (
	"fmt"

	"github.com/jfcarter2358/gitdb"
)

const testDocSingle = `//gitdb:doc:begin
//gitdb:field:foo
ABC
//gitdb:meta:bar
//DEF
//gitdb:doc:end

`

const testDocMultiple = `//gitdb:doc:begin
//gitdb:field:foo
ABC
//gitdb:meta:bar
//DEF
//gitdb:doc:end

//gitdb:doc:begin
//gitdb:meta:bar
//JKL
//gitdb:field:foo
GHI
//gitdb:doc:end

`

type Foo struct {
	Hello string `json:"foo"`
	World string `json:"bar" gitdb_meta:"true"`
}

func main() {
	f := Foo{
		Hello: "hello",
		World: "world",
	}
	r := gitdb.Repo{
		URL:    "git@github.com:jfcarter2358/gitdb",
		Path:   "test/foo.tf",
		Ref:    "main",
		Branch: "test",
	}

	if err := r.Init(); err != nil {
		panic(err)
	}

	if err := r.Pull(); err != nil {
		panic(err)
	}

	bytes, err := gitdb.Marshal(f)
	if err != nil {
		panic(err)
	}
	if err := r.Post(bytes); err != nil {
		panic(err)
	}

	if err := r.Push("Test push"); err != nil {
		panic(err)
	}

	if err := r.PR("Test PR", "This is a test PR"); err != nil {
		panic(err)
	}

	dat, err := r.Get()
	if err != nil {
		panic(err)
	}
	var ff Foo
	if err := gitdb.Unmarshal(dat, &ff); err != nil {
		panic(err)
	}
	fmt.Println(ff)

}
