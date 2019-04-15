# migrate-vfsdata-source
Migrate source support for vfsdata


# Usage
```go
import (
     vfsdata "github.com/neermitt/migrate-vfsdata-source"
)

    s := Resource("", data.Migrations)

    d, err := vfsdata.WithInstance(s)
    if err != nil {
        return err
    }
    
    migration.NewWithInstance("sql", d, ...)
```
