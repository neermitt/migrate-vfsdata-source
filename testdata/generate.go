// +build ignore

package main

import (
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/shurcooL/vfsgen"
)

func main() {
    var cwd, _ = os.Getwd()
    migrations := http.Dir(filepath.Join(cwd, "testdata/sql"))
    if err := vfsgen.Generate(migrations, vfsgen.Options{
        Filename:     "testdata/migrations_vfsdata.go",
        PackageName:  "testdata",
        BuildTags:    "!deploy_build",
        VariableName: "Migrations",
    }); err != nil {
        log.Fatalln(err)
    }
}
