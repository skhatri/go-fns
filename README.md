# go-fns

[![Build](https://github.com/skhatri/go-fns/actions/workflows/build.yml/badge.svg)](https://github.com/skhatri/go-fns/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/skhatri/go-fns)](https://goreportcard.com/report/github.com/skhatri/go-fns)
[![GoDoc](https://godoc.org/github.com/skhatri/go-fns?status.svg)](https://godoc.org/github.com/skhatri/go-fns)

A collection of utility functions to boost productivity when working with Go projects. This package provides convenient functions for working with collections, file operations, data conversion, and more.

## Installation

```bash
go get github.com/skhatri/go-fns
```

## Packages

### Collections

The `collections` package provides utility functions for working with Go collections like maps and sets.

#### Maps

```go
import "github.com/skhatri/go-fns/lib/collections"

// Copy attributes between maps
source := map[string]string{"name": "John", "age": "30"}
target := make(map[string]string)
value := collections.CopyAttribute("name", source, target)

// Filter map by key
m := map[string]int{"a": 1, "b": 2, "c": 3}
filtered := collections.FilteredByKey(m, func(k string) bool {
    return k == "a" || k == "b"
})

// Map values with string keys
source := map[interface{}]interface{}{
    "name": "John",
    "age":  30,
    "nested": map[interface{}]interface{}{
        "city": "New York",
    },
}
result := collections.MapByStringKey(source)

// Deep copy map
src := map[string]interface{}{
    "name": "John",
    "age":  30,
}
dest := make(map[string]interface{})
collections.CopyMap(dest, src)
```

#### Sets

```go
// Create a new set
items := []string{"a", "b", "a", "c"}
set := collections.NewSet(items)

// Create a set with custom key function
type User struct {
    ID   string
    Name string
}
users := []User{
    {ID: "1", Name: "John"},
    {ID: "1", Name: "John"}, // Duplicate
    {ID: "2", Name: "Jane"},
}
userSet := collections.NewSetWithComparator(users, func(u User) string {
    return u.ID
})

// Check if set contains value
exists := set.Contains("a")

// Add values to set
set.AddWithKeyFunc(User{ID: "3", Name: "Bob"}, func(u User) string {
    return u.ID
})

// Convert set to list
list := set.ToList()
```

### Converters

The `converters` package provides functions for converting between different data formats and handling file-based data operations.

```go
import "github.com/skhatri/go-fns/lib/converters"

// Read and unmarshal files
var config struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// Read JSON file
err := converters.UnmarshalJsonFile("config.json", &config)

// Read YAML file
err := converters.UnmarshalFile("config.yaml", &config)

// Marshal to files
data := struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}{
    Name: "John",
    Age:  30,
}

// Write JSON file
err := converters.MarshalToJsonFile(data, "output.json")

// Write pretty JSON file
err := converters.MarshalToJsonPrettyFile(data, "output.json")

// Convert to JSON string
jsonStr := converters.MarshalToJson(data)

// Read from reader
err := converters.ReadTo(reader, &data)
```

### File System (fs)

The `fs` package provides utilities for file system operations, including directory management and file handling.

```go
import "github.com/skhatri/go-fns/lib/fs"

// Directory operations
err := fs.CreateDir("path/to/dir")
err := fs.CreateDirIfNotExists("path/to/dir")
err := fs.DeleteDir("path/to/dir")
err := fs.DeleteDirIfExists("path/to/dir")
err := fs.EnsureEmptyDir("path/to/dir")

// Read files
content, err := fs.ReadBytes("file.txt")

// Read zip entries
zipFile, err := zip.OpenReader("archive.zip")
if err != nil {
    return err
}
defer zipFile.Close()

for _, file := range zipFile.File {
    content, err := fs.ReadZipEntry(file)
    if err != nil {
        return err
    }
    // Process content...
}

// List files
files := fs.ListFiles("path/to/dir", ".txt")

// Parse password entries
password, err := fs.ParsePasswordEntry("file:passwords.txt")
```

### Expressions (expr)

The `expr` package provides utilities for evaluating expressions, particularly environment variable expressions.

```go
import "github.com/skhatri/go-fns/lib/expr"

// Evaluate environment variable expressions
result := expr.SolveEnvExpression("${ENV_VAR}")

// Compare environment variables
result := expr.SolveEnvExpression("env.USER == john")
result := expr.SolveEnvExpression("env.PORT != 8080")
```

### Types

The `types` package provides custom types and their implementations.

```go
import "github.com/skhatri/go-fns/lib/types"

// Regular Expression type
re, err := types.Compile("^test.*")
if err != nil {
    return err
}

// Must compile (panics on error)
re := types.MustCompile("^test.*")

// Marshal/Unmarshal regex
text, err := re.MarshalText()
err = re.UnmarshalText(text)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

