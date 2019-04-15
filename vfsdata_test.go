package migrate_vfsdata_source_test

import (
	"testing"

	st "github.com/golang-migrate/migrate/v4/source/testing"
	. "github.com/neermitt/migrate-vfsdata-source"

	"github.com/neermitt/migrate-vfsdata-source/testdata"
)

func Test(t *testing.T) {
	// wrap assets into Resource first
	s := Resource("", testdata.Migrations)

	d, err := WithInstance(s)
	if err != nil {
		t.Fatal(err)
	}
	st.Test(t, d)
}

func TestWithInstance(t *testing.T) {
	// wrap assets into Resource
	s := Resource("", testdata.Migrations)

	_, err := WithInstance(s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOpen(t *testing.T) {
	b := &VfsData{}
	_, err := b.Open("")
	if err == nil {
		t.Fatal("expected err, because it's not implemented yet")
	}
}
