# envy
Use envy to manage environment variables with your OS keychain

To use the tool, clone the repo and run `make`. To use the library, run `go get` (or just build an app which imports the library using Go modules).

## History
There have been several tools for injecting environment variables from files into a command. These can broken down into two categories, broadly speaking:

1. dotenv and copycats, which take key-value pairs from a dot-file and inject the contents into environment variables
2. envy and similar tools which store key-value pairs in a secure database of some type, ditto

The basic idea is, you execute the env-variable managing tool, which gets key-value pairs from somewhere, and then runs another command for you:

```
$ envy exec [realm] command [args ...]
```

for example,

```
$ envy exec dev curl -s -H "Authorization: Bearer $(token)" \
  https://my.server.com/add -X POST -d '{"item": "spork"}'
```

where we've previously added the token

```
$ envy add dev token=8inlknmdgoi8uap8ow3hw3.pws9jpo9jskgs....sldkfs
```

Many times the variables needed are secrets, such as a credential needed to renew an OAuth token, or perhaps the token itself. As a result, even though dot-files can be set with 0600 or 0400 permissions (only the owner has privileges), there's some risk to having these credentials in plaintext form.

envy certainly isn't unique, but I needed one or two capabilities not found elsewhere, and 

* I wanted to keep the implementation simple
* it was also a good opportunity to build an example app in Go

I have deliberately minimized the dependencies, which are basically the [Bolt database](https://github.com/boltdb/bolt) and [go-keyring](https://github.com/zalando/go-keyring). I have also avoided the many layers of abstraction typical of ["enterprise Fizzbuzz"](https://github.com/EnterpriseQualityCoding/FizzBuzzEnterpriseEdition) style development.

## How it works
Variables (key-value pairs) are grouped into "realms" which is just a shorter way to type "namespaces". Because these variables are primarily used as environment variables, they're stored in a map of string keys to string values.

envy maintains a Bolt database in your "user config" directory, for example, `$HOME/Library/Application Support` on macOS. That database has a bucket for each realm, and an entry in the bucket for each key-value pair.

With the variable is some metadata: we keep the last-modified timestamp, size, and a secure hash of the value part of the key-value pair. The hash is also used with AES-GCM when that value is encrypted. The encrypted data and the metadata in JSON form are converted to Base64 encoding and then stored together a single object identified by the key. Only the (possibly secret) value is encrypted; the metadata isn't, but if the hash is changed, decryption fails.

The secret key needed to run AES-GCM is stored in your system's secure keychain, which on macOS means the login keychain that's visible in Keychain Access. (Note that you can see and even edit the secret key in Keychain Access or using the `security` command -- but if you change or delete that key, you'll never get your data back out of the Bolt database.)

The secret key is added once to the keychain when you first run envy. If you want to wipe everything and start over, then

1. remove the key, which is named `matt4biz-envy-secret-key`
2. remove the database, which is `envy/envy.db` in your config directory

## Commands
There are five commands, but one of them just lists the version of the program; you can also type `envy -h` to see usage:

```
envy: a tool to securely store and retrieve environment variables.

Variables are key-value pairs stored in a "realm" (or "namespace") of which
there may be one or more. All data is stored in a DB within the user's "config"
directory, encrypted with a per-user secret key stored in the system keychain.

All operations take place in one of five subcommands. Add will create a realm
if it doesn't exist, or overwrite keys in a realm that already exists. Drop
may be used to delete one key or an entire realm. Exec will execute a command
with arguments, with value(s) from the realm injected as environment variables.

Usage:
  -h show this help message

  add       realm key=value [key=value ...]
  drop      realm[/key]
  exec      realm[/key] command [args...]
  list [-d] realm
    -d decrypt & show secrets also
  version

Listing a realm displays a timestamp, size, and hash for each key-value pair.
```

### Add
The add subcommand adds one or more keys to a realm. The realm will be created if it doesn't exist. If it exists already, the key(s) you set will overwrite any matching key in the realm.

For example, assuming a new database:

```
$ envy add test a=1 b=2
```

will set up two key-value pairs. Note that keys are case-sensitive.

If you then run

```
$ envy add test a=3
```

the value for key "a" will change, but other keys will not be disturbed.

### List
The list command with no realm will list the available realms in the database. For example, after the commands above,

```
$ envy list
test
```

and

```
$ envy list test
a   2020-10-11T23:28:05-06:00  1  3cf3aef
b   2020-10-11T23:28:05-06:00  1  39c6844
```

where just the first seven characters of the hash are shown.

The `-d` option will also show the decrypted data:

```
$ envy list -d test
a   2020-10-11T23:28:05-06:00  1  3cf3aef   3
b   2020-10-11T23:28:05-06:00  1  39c6844   2
```

### Drop
The drop command can delete one key from a realm, or the entire realm.

For example,

```
$ envy drop test/b
$ envy list test
a   2020-10-11T23:28:05-06:00  1  3cf3aef
```

while

```
$ envy drop test
$ envy list test
fetching test: realm test: not found
$ envy list
```

shows that we've returned the database to its empty state.

### Exec
Of course, the exec command is the main reason for this tool. Given a realm (or a specific key from a realm), envy will execute another command with its environment variables augmented by data that envy stores. (See the example above.)

envy can pass (some) signals through to its child process, particularly control-C, so it's possible to kill off the child if you need to. The childs standard input, output, and error output mirror envy's environment.

## As a library
envy is not just a command-line tool, it's also a library that can be used in building another tool.

To get started, you just need to create the `Envy` object, passing in the directory where the database lives:

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/matt4biz/envy"
)

func main() {
	dir, err := os.UserConfigDir()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(dir)

	// the standard location is config-dir/envy

	e, err := envy.New(path.Join(dir, "envy"))

	if err != nil {
		log.Fatal(err)
	}

	m, err := e.FetchAsJSON("test")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(m))
}
```

with the output on macOS (given the examples above):

```
$ go run .
/Users/<your-login>/Library/Application Support
{"a":"3","b":"2"}
```

## Details
The repo is organized simply:

```
$ tree
.
├── LICENSE
├── Makefile
├── README.md
├── cmd
│   ├── add.go
│   ├── add_test.go
│   ├── app.go
│   ├── drop.go
│   ├── drop_test.go
│   ├── exec.go
│   ├── exec_test.go
│   ├── list.go
│   ├── list_test.go
│   ├── main.go
│   ├── test_common.go
│   └── version.go
├── envy.go
├── envy_test.go
├── go.mod
├── go.sum
├── hack
│   └── main.go
├── internal
│   ├── db.go
│   ├── db_test.go
│   ├── extract.go
│   ├── extract_test.go
│   ├── ring.go
│   ├── ring_test.go
│   ├── sealer.go
│   └── sealer_test.go
└── test
    └── test.sh
```

The top-level library API is in `envy.go`; everything it needs is in the `internal` sub-package. The CLI and subcommands are in `cmd`.

The makefile has only a few targets:

- envy (the default)
- lint, assuming you have [golangci-lint](https://github.com/golangci/golangci-lint) installed (at the moment there's no special config file, so linting uses the defaults)
- test, which runs all UTs with code coverage
- demo, which depends on another target, `child` (from `hack/main.go`)

The demo will run `child` as a subcommand; the child will print some environment variables and exit in 10 seconds unless it gets a control-C sooner.

The test script in `test/test.sh` is used by unit tests, and shouldn't be changed.

The unit tests use a mock keyring in memory and auto-delete their temporary Bolt DB, so they have no effect on your "real" envy secret key and secure DB.

Code coverage is above 70% except for the subcommands; that's only about 50% (more error path coverage needed).

The design of the CLI was influenced by Carl Johnson's [_Writing Go CLIs With Just Enough Architecture_](https://blog.carlmjohnson.net/post/2020/go-cli-how-to-and-advice/).
