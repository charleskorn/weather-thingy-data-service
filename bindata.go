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

var _db_migrations_0002_create_variables_table_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x5c\xcf\x41\x4f\xc3\x30\x0c\x05\xe0\x7b\x7e\xc5\x3b\x6e\x82\x49\x83\xeb\x4e\xa6\x35\x5a\x44\x9a\x16\xcf\x01\x8d\xcb\x14\x58\x84\x22\xb1\x82\xba\x02\x7f\x9f\x96\x43\xa9\x38\x3e\xf9\xb3\xe5\xb7\x5a\xe1\xe2\x94\x5f\xbb\xd8\x27\x84\x0f\x53\x08\x93\x32\x94\x6e\x1c\xe3\x2b\x76\x39\x3e\xbf\xa5\x33\x16\x06\x53\x3a\xe4\x23\x76\x2c\x96\x1c\x1a\xb1\x15\xc9\x1e\x77\xbc\xbf\x1c\x44\x1b\x4f\x09\x0f\x24\xc5\x96\x64\x71\xb5\x5e\x2f\xe1\x6b\x85\x0f\xce\x21\x78\x7b\x1f\x78\x44\x9f\x6d\xee\xcf\x93\xba\x9e\xa1\x71\xfa\xd2\xa5\xe1\x93\x23\xd4\x56\xbc\x53\xaa\x1a\x3c\x5a\xdd\xfe\x46\x3c\xd5\x9e\xff\x2e\x96\x7c\x4b\xc1\x29\x8a\x20\xc2\x5e\x0f\xd3\x86\x59\x6e\x8c\x99\xd7\x2a\xdf\xbf\x5b\x53\x4a\xdd\xfc\xaf\xb5\x31\x3f\x01\x00\x00\xff\xff\x66\xec\xe1\xcf\xfe\x00\x00\x00")

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

	info := bindata_file_info{name: "db/migrations/0002_create_variables_table.sql", size: 254, mode: os.FileMode(420), modTime: time.Unix(1428379540, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _db_migrations_0003_create_data_table_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x90\xcd\x4a\x03\x31\x14\x46\xf7\x79\x8a\x6f\x39\x83\x53\x50\x70\xd7\x55\x9c\xb9\xc5\xe0\xfc\x71\x7b\x83\xd4\x4d\x89\x34\x94\x40\x5b\xa5\x46\x7d\x7d\xe3\x94\xb6\x59\x75\x91\x45\xe0\x9c\xcb\xe1\x9b\xcd\x70\xb7\x0f\xdb\xa3\x8b\x1e\xf6\x53\xd5\x4c\x5a\x08\xa2\x9f\x5a\xc2\xc6\x45\x87\x42\x01\x6e\xeb\x0f\x71\x1d\x36\x30\xbd\xa0\x1f\xd2\xb3\x6d\x0b\xa6\x05\x31\xf5\x35\x2d\x4f\xc0\x17\x8a\x33\x58\x56\xc9\xfa\x71\xc7\xe0\xde\x77\xfe\x96\x78\x66\x92\x9b\xe1\x93\x1e\xc3\xde\x43\x4c\x47\x4b\xd1\xdd\x88\x57\x23\xcf\xd3\x17\x6f\x43\x4f\xd7\x63\x0d\x2d\xb4\x6d\x05\xb5\xe5\x74\x53\xd6\x17\xe3\x94\xb0\xfb\xf6\x89\xeb\x88\x4d\x5d\x3c\xdc\x57\x78\x2c\x2f\xea\x3f\x30\xb2\xe9\x34\xaf\xf0\x42\xab\x6b\x7d\x95\xa7\x57\x53\x48\xa9\xca\xb9\x52\xf9\x58\xcd\xc7\xef\x41\x35\x3c\x8c\xd9\x58\x73\xf5\x17\x00\x00\xff\xff\xef\x23\xd0\x7e\x4f\x01\x00\x00")

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

	info := bindata_file_info{name: "db/migrations/0003_create_data_table.sql", size: 335, mode: os.FileMode(420), modTime: time.Unix(1427330306, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _db_migrations_0004_variables_table_add_display_decimal_places_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x28\x4b\x2c\xca\x4c\x4c\xca\x49\x2d\x56\x70\x74\x71\x51\x70\xf6\xf7\x09\xf5\xf5\x53\x48\xc9\x2c\x2e\xc8\x49\xac\x8c\x4f\x49\x4d\xce\xcc\x4d\xcc\x89\x07\x72\x92\x81\x2a\x3c\xfd\x42\x14\xfc\xfc\x81\x38\xd4\xc7\x47\xc1\xc5\xd5\xcd\x31\xd4\x27\x44\x41\xc3\x40\xd3\x1a\x97\x89\x60\x51\xfc\x66\xba\x04\xf9\x07\xc0\xcc\xb2\xe6\xe2\x42\x76\xa9\x4b\x7e\x79\x1e\x0e\x93\xc1\xba\xf0\x1a\x6c\xcd\x05\x08\x00\x00\xff\xff\xbc\x26\xc0\xe2\xf5\x00\x00\x00")

func db_migrations_0004_variables_table_add_display_decimal_places_sql_bytes() ([]byte, error) {
	return bindata_read(
		_db_migrations_0004_variables_table_add_display_decimal_places_sql,
		"db/migrations/0004_variables_table_add_display_decimal_places.sql",
	)
}

func db_migrations_0004_variables_table_add_display_decimal_places_sql() (*asset, error) {
	bytes, err := db_migrations_0004_variables_table_add_display_decimal_places_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "db/migrations/0004_variables_table_add_display_decimal_places.sql", size: 245, mode: os.FileMode(420), modTime: time.Unix(1430807482, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _db_migrations_0005_agents_table_enforce_name_not_null_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x48\x4c\x4f\xcd\x2b\x29\x56\x80\x08\x39\xfb\xfb\x84\xfa\xfa\x29\xe4\x25\xe6\xa6\x2a\x04\xbb\x86\x28\xf8\xf9\x03\x71\xa8\x8f\x8f\x35\x17\x17\xb2\x11\x2e\xf9\xe5\x79\xc4\x19\xe2\x12\xe4\x1f\x80\x64\x0a\x20\x00\x00\xff\xff\x48\x44\x96\xb5\x88\x00\x00\x00")

func db_migrations_0005_agents_table_enforce_name_not_null_sql_bytes() ([]byte, error) {
	return bindata_read(
		_db_migrations_0005_agents_table_enforce_name_not_null_sql,
		"db/migrations/0005_agents_table_enforce_name_not_null.sql",
	)
}

func db_migrations_0005_agents_table_enforce_name_not_null_sql() (*asset, error) {
	bytes, err := db_migrations_0005_agents_table_enforce_name_not_null_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "db/migrations/0005_agents_table_enforce_name_not_null.sql", size: 136, mode: os.FileMode(420), modTime: time.Unix(1453676673, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _db_migrations_0006_create_users_table_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x6c\x51\xcd\x6e\xf2\x30\x10\xbc\xe7\x29\xe6\x08\xfa\x3e\x2e\x55\x7b\xe2\xb4\x24\x8b\x88\xea\x38\xd4\xb1\x5b\xd1\x4b\x64\x35\x16\x58\x82\x80\xe2\x54\xbc\x7e\x4d\x04\x14\x54\x6e\x89\xe7\x6f\x77\x76\x32\xc1\xbf\x9d\x5f\x77\xb6\x77\x30\x87\x24\x55\x4c\x9a\xa1\x69\x26\x18\xdf\xc1\x75\x01\xa3\x04\xc3\x57\xed\x1b\x54\xac\x72\x12\x58\xaa\xbc\x20\xb5\xc2\x2b\xaf\xfe\x47\xd4\xed\xac\xdf\xe2\x9d\x54\xba\x20\x35\x7a\x7a\x79\x1e\x43\x96\x1a\xd2\x08\x01\x23\xf3\x37\xc3\x27\xd6\xc1\x86\x70\xdc\x77\x4d\xed\x7b\x17\xe3\xfc\xbe\x0d\xc8\xa5\xbe\x52\xef\x38\xc1\x6e\x7b\xcc\x56\x9a\xe9\x31\xbe\xb1\x61\xf3\x00\xf7\xa1\xb6\xcd\xce\xb7\x98\x95\xa5\x60\x92\xbf\x73\x64\x3c\x27\x23\x34\xe6\x24\xaa\x61\x9c\xaf\xce\xc5\x9d\x1b\xe8\xbc\xe0\x4a\x53\xb1\xc4\x47\xae\x17\xc3\x2f\x3e\x4b\xc9\x7f\xa5\xa9\x51\x8a\xa5\xae\xaf\x8a\x64\x3c\x4d\x12\x12\x9a\xd5\xb9\x30\xbb\x76\x6d\x1f\x40\x59\x86\xb4\x14\xa6\x90\xd8\x1f\xdb\xd8\xdc\xa5\xbe\xdb\x75\xa1\x78\xce\xd1\x2f\xe5\xea\x52\xf4\x99\x76\x72\xbd\x3d\x4b\x16\x4d\x1e\xc5\x64\xaa\x5c\xde\xe7\xf8\x26\x4a\x87\xe7\x9b\x03\x4e\x93\x9f\x00\x00\x00\xff\xff\xdb\xfb\x3e\x99\xe4\x01\x00\x00")

func db_migrations_0006_create_users_table_sql_bytes() ([]byte, error) {
	return bindata_read(
		_db_migrations_0006_create_users_table_sql,
		"db/migrations/0006_create_users_table.sql",
	)
}

func db_migrations_0006_create_users_table_sql() (*asset, error) {
	bytes, err := db_migrations_0006_create_users_table_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "db/migrations/0006_create_users_table.sql", size: 484, mode: os.FileMode(420), modTime: time.Unix(1453681917, 0)}
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
	"db/migrations/0001_create_agents_table.sql":                        db_migrations_0001_create_agents_table_sql,
	"db/migrations/0002_create_variables_table.sql":                     db_migrations_0002_create_variables_table_sql,
	"db/migrations/0003_create_data_table.sql":                          db_migrations_0003_create_data_table_sql,
	"db/migrations/0004_variables_table_add_display_decimal_places.sql": db_migrations_0004_variables_table_add_display_decimal_places_sql,
	"db/migrations/0005_agents_table_enforce_name_not_null.sql":         db_migrations_0005_agents_table_enforce_name_not_null_sql,
	"db/migrations/0006_create_users_table.sql":                         db_migrations_0006_create_users_table_sql,
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
			"0001_create_agents_table.sql":                        &_bintree_t{db_migrations_0001_create_agents_table_sql, map[string]*_bintree_t{}},
			"0002_create_variables_table.sql":                     &_bintree_t{db_migrations_0002_create_variables_table_sql, map[string]*_bintree_t{}},
			"0003_create_data_table.sql":                          &_bintree_t{db_migrations_0003_create_data_table_sql, map[string]*_bintree_t{}},
			"0004_variables_table_add_display_decimal_places.sql": &_bintree_t{db_migrations_0004_variables_table_add_display_decimal_places_sql, map[string]*_bintree_t{}},
			"0005_agents_table_enforce_name_not_null.sql":         &_bintree_t{db_migrations_0005_agents_table_enforce_name_not_null_sql, map[string]*_bintree_t{}},
			"0006_create_users_table.sql":                         &_bintree_t{db_migrations_0006_create_users_table_sql, map[string]*_bintree_t{}},
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
