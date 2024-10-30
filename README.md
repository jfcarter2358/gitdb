# GitDB

## About

GitDB is a "database" using Git as a backend. THIS SHOULD NOT BE USED FOR THE VAST MAJORITY OF APPS. This package was built for a use-case of being able to manage Terraform files using an API. As such, this "database" is primarily static with occasional edits, which this package helps facilitate.

## Usage

In order to use `gitdb`, you'll need to setup a struct for whatever object you want to manage. When adding tags to struct fields, only those with `json` tags will be written to the file in Git. Additionally, if you want a value to be commented out in the file (e.g. a metadata value that you don't want Terraform to try and fail at applying) then add the `gitdb_meta:true` tag to that struct field. For example, the following struct

```go
type Foo struct {
	Hello string `json:"foo"`
	World string `json:"bar" gitdb_meta:"true"`
}
```

corresponds to the following repo file contents:

```hcl
//gitdb:field:foo
<Hello value here>
//gitdb:meta:bar
//<World value here>
//gitdb:doc:end
```

Below is an example of using `gitdb` in an application:

```go
package main

import (
	"fmt"

	"github.com/jfcarter2358/gitdb"
)

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
		URL:  "git@github.com:jfcarter2358/gitdb",
		Path: "test/foo.tf",
		Ref:  "main",
	}

    // setup the repo
	if err := r.Init(); err != nil {
		panic(err)
	}

    // pull from the latest ref, this normally would be run in a cron somewhere, not right after init
	if err := r.Update(); err != nil {
		panic(err)
	}

    // write our object to the repo file
	bytes, err := gitdb.Marshal(f)
	if err != nil {
		panic(err)
	}
	if err := r.Post(bytes); err != nil {
		panic(err)
	}

    // get the contents of the repo file
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

```

## TODO

- [ ] Add ability to create GitHub PR
