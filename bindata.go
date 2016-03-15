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

var _db_migrations_0006_create_users_table_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x6c\x51\x3d\x6f\xc2\x30\x10\xdd\xf3\x2b\xde\x08\x6a\x59\xaa\x76\x62\x3a\x92\x43\x44\x75\x1c\xea\xd8\xad\xe8\x12\x59\x4d\x04\x96\x20\xa0\x38\x15\x7f\xbf\x26\x02\x9a\xaa\xd9\x6c\xbd\xaf\xbb\x77\xb3\x19\x1e\x0e\x6e\xdb\xda\xae\x86\x39\x45\xb1\x62\xd2\x0c\x4d\x0b\xc1\xf8\xf6\x75\xeb\x31\x89\xd0\xbf\x4a\x57\xa1\x60\x95\x92\xc0\x5a\xa5\x19\xa9\x0d\x5e\x79\xf3\x18\xd0\xfa\x60\xdd\x1e\xef\xa4\xe2\x15\xa9\xc9\xd3\xcb\xf3\x14\x32\xd7\x90\x46\x08\x18\x99\xbe\x19\xbe\xb0\x4e\xd6\xfb\xf3\xb1\xad\x4a\xd7\xd5\x21\xce\x1d\x1b\x8f\x54\xea\x3b\xf5\x0f\xc7\xdb\x7d\x87\xc5\x46\x33\x8d\xe3\x3b\xeb\x77\x23\xb8\xf3\xa5\xad\x0e\xae\xc1\x22\xcf\x05\x93\xfc\x9d\x23\xe1\x25\x19\xa1\xb1\x24\x51\xf4\xe3\x7c\xb5\x75\xd8\xb9\x82\x4e\x33\x2e\x34\x65\x6b\x7c\xa4\x7a\xd5\x7f\xf1\x99\x4b\xfe\x2f\x8d\x8d\x52\x2c\x75\x79\x57\x44\xd3\x79\x14\x91\xd0\xac\xae\x85\xd9\x6d\xdd\x74\x1e\x94\x24\x88\x73\x61\x32\x89\xe3\xb9\x09\xcd\xdd\xea\x1b\xae\x0b\xc5\x4b\x0e\x7e\x31\x17\xb7\xa2\xaf\xb4\x8b\xeb\xf0\x2c\x49\x30\x19\x8b\x49\x54\xbe\x1e\xcd\x09\xfa\x1e\x1b\x5c\x71\x1e\xfd\x04\x00\x00\xff\xff\xf9\xfb\xd1\x76\xe9\x01\x00\x00")

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

	info := bindata_file_info{name: "db/migrations/0006_create_users_table.sql", size: 489, mode: os.FileMode(420), modTime: time.Unix(1458015538, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _db_migrations_0007_agents_table_add_token_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x48\x4c\x4f\xcd\x2b\x29\x56\x70\x74\x71\x51\x70\xf6\xf7\x09\xf5\xf5\x53\x28\xc9\xcf\x4e\xcd\x53\x08\x73\x0c\x72\xf6\x70\x0c\xd2\x30\x34\x30\xd0\x54\xf0\xf3\x0f\x51\xf0\x0b\xf5\xf1\xb1\xe6\xe2\x42\x36\xca\x25\xbf\x3c\x0f\x9b\x61\x2e\x41\xfe\x01\x28\xa6\x59\x73\x01\x02\x00\x00\xff\xff\xcb\xb6\x9f\x66\x82\x00\x00\x00")

func db_migrations_0007_agents_table_add_token_sql_bytes() ([]byte, error) {
	return bindata_read(
		_db_migrations_0007_agents_table_add_token_sql,
		"db/migrations/0007_agents_table_add_token.sql",
	)
}

func db_migrations_0007_agents_table_add_token_sql() (*asset, error) {
	bytes, err := db_migrations_0007_agents_table_add_token_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "db/migrations/0007_agents_table_add_token.sql", size: 130, mode: os.FileMode(420), modTime: time.Unix(1458015607, 0)}
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
	"db/migrations/0007_agents_table_add_token.sql":                     db_migrations_0007_agents_table_add_token_sql,
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
			"0007_agents_table_add_token.sql":                     &_bintree_t{db_migrations_0007_agents_table_add_token_sql, map[string]*_bintree_t{}},
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
