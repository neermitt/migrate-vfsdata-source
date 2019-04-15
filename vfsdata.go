package migrate_vfsdata_source

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4/source"
)

func init() {
	source.Register("vfsdata", &VfsData{})
}

func Resource(path string, fs http.FileSystem) *AssetSource {
	return &AssetSource{
		Path: path,
		FS:   fs,
	}
}

type AssetSource struct {
	Path string
	FS   http.FileSystem
}

type VfsData struct {
	path       string
	as         *AssetSource
	migrations *source.Migrations
}

var (
	ErrNoAssetSource = fmt.Errorf("expects *AssetSource")
)

func WithInstance(instance interface{}) (source.Driver, error) {
	if _, ok := instance.(*AssetSource); !ok {
		return nil, ErrNoAssetSource
	}
	as := instance.(*AssetSource)

	vfsData := &VfsData{
		path:       "/",
		as:         as,
		migrations: source.NewMigrations(),
	}

	f, err := as.FS.Open(as.Path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	_ = f.Close()
	for _, fi := range list {
		m, err := source.DefaultParse(fi.Name())
		if err != nil {
			continue // ignore files that we can't parse
		}

		if !vfsData.migrations.Append(m) {
			return nil, fmt.Errorf("unable to parse file %v", fi)
		}
	}

	return vfsData, nil
}

func (f *VfsData) Open(url string) (source.Driver, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f *VfsData) Close() error {
	// nothing do to here
	return nil
}

func (f *VfsData) First() (version uint, err error) {
	if v, ok := f.migrations.First(); !ok {
		return 0, &os.PathError{Op: "first", Path: f.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (f *VfsData) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := f.migrations.Prev(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("prev for version %v", version), Path: f.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (f *VfsData) Next(version uint) (nextVersion uint, err error) {
	if v, ok := f.migrations.Next(version); !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("next for version %v", version), Path: f.path, Err: os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (f *VfsData) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := f.migrations.Up(version); ok {
		r, err := f.as.FS.Open(path.Join(f.path, m.Raw))
		if err != nil {
			return nil, "", err
		}
		return r, m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: f.path, Err: os.ErrNotExist}
}

func (f *VfsData) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := f.migrations.Down(version); ok {
		r, err := f.as.FS.Open(path.Join(f.path, m.Raw))
		if err != nil {
			return nil, "", err
		}
		return r, m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: f.path, Err: os.ErrNotExist}
}
