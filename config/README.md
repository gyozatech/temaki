# Configuration reader

This package allows getting/setting config variables from the environment variables and from the most classical `.env` file situated at the top level of the project folder.

## Usage example

So imagine to have the following project folder:
```bash
my-project
   main.go
   .env
   .env.sample
   .gitignore
   repository
       db_connection.go
```
If we store the connection config in the `.env` file we can have a file like this:
```bash
$ cat .env
# username:password@tcp(host:port)/dbname
MYSQL_URI=localhost:3306
MYSQL_DB_NAME=civo_api_development
MYSQL_USER=root
MYSQL_PASSWORD=
```
We can use those the following way in `repository/db_connection.go`:
```go
package repository

import (
	"database/sql"
	"fmt"
	"github.com/temaki/config"
	_ "github.com/go-sql-driver/mysql"
)

// DB is the global object holding the DB connection
var DB *sql.DB = initDB()

func initDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Get("MYSQL_USER"), config.Get("MYSQL_PASSWORD"), config.Get("MYSQL_URI"), config.Get("MYSQL_DB_NAME")))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Db connection inited")
	return db
}
```

## Pre-processing rules

The lines containing the "#" char will be considered as comments and all the lines where no "=" char is present, are ignored too.
The spaces are automatically removed from the keys but not from the values where a label value, for example, can have some spaces.
The quotes (single, double, back-tick) are automatically removed from the value.
So from the following file:
```
MY_LABEL   =This label contains spaces
MY_LABEL2="This label is quoted"
EMPTY_LABEL=
``` 
We get the following Go equivalent variables:
```go
import (
   "github.com/temaki/config"
   "fmt"
)

func expectations() error{
    if config.Get("MY_LABEL") != "This label contains spaces" || 
          config.Get("MY_LABEL_2") != "This label is quoted"  ||
              config.Get("EMPTY_LABEL") != "" {
        return fmt.Error("I didn't expect that!")
    }
}
```

## Methods

The `temaki/config` package allowing the following utilities:
```go
import "github.com/temaki/config"

config.Get("MY_VARIABLE") // from the env vars or the .env file
config.Set("MY_VARIABLE", "Value1234") // sets the specified variable
config.Override(map[string]string{ "MY_VARIABLE", "mocked"}) // overrider the whole bunch of variables for testing purposes
```
