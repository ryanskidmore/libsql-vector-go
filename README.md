## libsql-vector-go

This library is intended to provide an easy way to work with [libSQL](https://github.com/tursodatabase/libsql) vectors which are used by providers like [Turso](https://turso.tech/). 
This library works in a similar way to [pgvector-go](https://github.com/pgvector/pgvector-go), which this code has been adapted from.

Typically vectors in databases (like those provided by pgvector) take the format of an array of floating points, e.g. `UPDATE table SET embedding = '[1,2,3]'`.
However the vectors in libSQL must be wrapped in a `vector()` function, e.g. `UPDATE table SET embedding = vector('[1,2,3]')`
which means that most other database vector libraries are incompatible.

### Getting Started

Get the library:
```shell
go get github.com/ryanskidmore/libsql-vector-go
```

Then import it:
```go
import (
	github.com/ryanskidmore/libsql-vector-go
)
```

### Gorm

If you're using gorm, you want to instead import the gorm variant which is held as a separate package:
```go
import (
    libsqlvectorgorm github.com/ryanskidmore/libsql-vector-go/gorm
)
```

You can then use the type in your models:

```go
type Product struct {
	Embedding *libsqlvectorgorm.Vector `gorm:"type:F32_BLOB(1536)"`
}
```