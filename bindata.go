package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _db_migrations_0001_create_agents_table_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x54\x8e\xcf\x4a\x86\x40\x14\xc5\xf7\xf3\x14\x67\xf9\x49\x09\xb6\x76\x75\xd3\x1b\x0e\x8d\xa3\x5c\xef\x14\xb6\x11\xc9\x41\x5c\x68\x61\x42\xaf\x9f\x15\x44\xdf\xee\xfc\xe1\x1c\x7e\x69\x8a\x9b\x75\x99\xf7\xf1\x88\x08\xef\xa6\x10\x26\x65\x28\xdd\x3b\xc6\x38\xc7\xed\xf8\xc0\xc5\xe0\x57\x0e\xcb\x84\x8e\xc5\x92\x43\x2b\xb6\x26\xe9\xf1\xc8\xfd\xed\x59\x6f\xe3\x1a\xf1\x44\x52\x54\x24\x97\xbb\x2c\x4b\xbe\xc3\xd7\x3d\x9e\xaf\x13\xd4\xd6\xdc\x29\xd5\x2d\x9e\xad\x56\x3f\x16\x2f\x8d\x67\xf8\x46\xe1\x83\x73\x28\xf9\x81\x82\x53\x14\x41\x84\xbd\x0e\x7f\x0b\x93\xe4\xc6\xfc\x47\x2c\xdf\x3e\x37\x53\x4a\xd3\x5e\x21\xe6\xe6\x2b\x00\x00\xff\xff\xc3\xbe\x48\xbc\xc7\x00\x00\x00")

func db_migrations_0001_create_agents_table_sql_bytes() ([]byte, error) {
	return bindata_read(
		_db_migrations_0001_create_agents_table_sql,
		"db/migrations/0001_create_agents_table.sql",
	)
}

func db_migrations_0001_create_agents_table_sql() (*asset, error) {
	bytes, err := db_migrations_0001_create_agents_table_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "db/migrations/0001_create_agents_table.sql", size: 199, mode: os.FileMode(420), modTime: time.Unix(1427176216, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _db_migrations_0002_create_variables_table_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x5c\x8f\x41\x4b\xc3\x40\x10\x46\xef\xfb\x2b\xbe\x63\x8b\x16\xaa\xd7\x9e\xc6\x64\xa4\x8b\x9b\x4d\x98\xce\x2a\xf5\x52\x56\xbb\xc8\x82\x8d\x92\x46\xfd\xfb\x46\x0f\x4b\xc9\xf1\xe3\xbd\x81\x79\xab\x15\xae\x4e\xf9\x6d\x88\x63\x42\xf8\x34\x95\x30\x29\x43\xe9\xce\x31\xbe\xe3\x90\xe3\xcb\x7b\x3a\x63\x61\x50\xd6\x21\x1f\xb1\x63\xb1\xe4\xd0\x89\x6d\x48\xf6\x78\xe0\xfd\xf5\x64\xf4\xf1\x94\xf0\x48\x52\x6d\x49\x16\x37\xeb\xf5\x12\xbe\x55\xf8\xe0\xdc\x1f\xfd\xea\xf3\x78\x2e\xf8\x76\x46\x5f\x87\x34\xbd\x70\x84\xda\x86\x77\x4a\x4d\x87\x27\xab\xdb\xff\x89\xe7\xd6\x73\x91\x51\xf3\x3d\x05\xa7\xa8\x82\x08\x7b\x3d\x94\x0b\xb3\xdc\x18\x73\xd9\x53\x7f\xfc\xf4\xa6\x96\xb6\x9b\xf7\x6c\xcc\x6f\x00\x00\x00\xff\xff\xd8\x32\x56\x33\xf7\x00\x00\x00")

func db_migrations_0002_create_variables_table_sql_bytes() ([]byte, error) {
	return bindata_read(
		_db_migrations_0002_create_variables_table_sql,
		"db/migrations/0002_create_variables_table.sql",
	)
}

func db_migrations_0002_create_variables_table_sql() (*asset, error) {
	bytes, err := db_migrations_0002_create_variables_table_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "db/migrations/0002_create_variables_table.sql", size: 247, mode: os.FileMode(420), modTime: time.Unix(1427176221, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _db_migrations_0003_create_data_table_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x90\x4f\x4b\x03\x31\x10\x47\xef\xf9\x14\xbf\xe3\x2e\x6e\x41\xc1\x5b\x4f\x71\x77\x8a\xc1\xfd\xc7\x74\x82\xd4\x4b\x89\x34\x94\x40\x5b\xa5\x46\xfd\xfa\xc6\xad\xdd\xdd\x93\x87\x1c\x02\xef\x0d\x6f\x66\xb1\xc0\xcd\x31\xec\xcf\x2e\x7a\xd8\x77\x55\x32\x69\x21\x88\x7e\xa8\x09\x3b\x17\x1d\x32\x05\xb8\xbd\x3f\xc5\x6d\xd8\xc1\xb4\x82\xb6\x4b\xcf\xd6\x35\x98\x56\xc4\xd4\x96\xb4\xbe\x00\x1f\xc8\xae\x60\x5e\x24\xeb\xcb\x9d\x83\x7b\x3d\xf8\xff\xc4\x2b\x93\xdc\x19\x3e\xe8\x31\x1c\x3d\xc4\x34\xb4\x16\xdd\xf4\x78\x36\xf2\x38\x7c\xf1\xd2\xb5\x34\x0d\xab\x68\xa5\x6d\x2d\x28\x2d\xa7\x99\xb2\x1d\x8d\x4b\xc2\xe1\xd3\x27\xae\x21\x36\x65\x76\x77\x5b\xe0\x3e\x1f\xd5\x5f\xa0\x67\xd3\x68\xde\xe0\x89\x36\x53\x7d\x31\x4f\x2f\x86\x90\x5c\xe5\x4b\xa5\xe6\xc7\xaa\xde\xbe\x4f\xaa\xe2\xae\xff\x3b\xd6\xb8\xc9\x52\xfd\x04\x00\x00\xff\xff\x22\xaf\x99\x0d\x54\x01\x00\x00")

func db_migrations_0003_create_data_table_sql_bytes() ([]byte, error) {
	return bindata_read(
		_db_migrations_0003_create_data_table_sql,
		"db/migrations/0003_create_data_table.sql",
	)
}

func db_migrations_0003_create_data_table_sql() (*asset, error) {
	bytes, err := db_migrations_0003_create_data_table_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "db/migrations/0003_create_data_table.sql", size: 340, mode: os.FileMode(420), modTime: time.Unix(1427176539, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"db/migrations/0001_create_agents_table.sql":    db_migrations_0001_create_agents_table_sql,
	"db/migrations/0002_create_variables_table.sql": db_migrations_0002_create_variables_table_sql,
	"db/migrations/0003_create_data_table.sql":      db_migrations_0003_create_data_table_sql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() (*asset, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"db": &_bintree_t{nil, map[string]*_bintree_t{
		"migrations": &_bintree_t{nil, map[string]*_bintree_t{
			"0001_create_agents_table.sql":    &_bintree_t{db_migrations_0001_create_agents_table_sql, map[string]*_bintree_t{}},
			"0002_create_variables_table.sql": &_bintree_t{db_migrations_0002_create_variables_table_sql, map[string]*_bintree_t{}},
			"0003_create_data_table.sql":      &_bintree_t{db_migrations_0003_create_data_table_sql, map[string]*_bintree_t{}},
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	if err != nil { // File
		return RestoreAsset(dir, name)
	} else { // Dir
		for _, child := range children {
			err = RestoreAssets(dir, path.Join(name, child))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
