// Code generated by go-bindata.
// sources:
// pkg/generated/manifests/default-config.yaml
// DO NOT EDIT!

package generated

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _defaultConfigYaml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x94\x54\x4b\x4f\xe3\x48\x10\xbe\xfb\x57\x94\xe0\xb2\x48\xc4\x01\xed\x65\xe5\x2b\x70\x40\x8b\xd0\x2a\xec\x1f\xa8\xb8\x2b\x49\x8b\x76\xb7\xe9\x47\x58\xff\xfb\xfd\xca\xed\x40\x66\x34\x8c\x66\xa4\x3c\x6c\xd7\xe3\x7b\x54\xb5\x2f\xe9\x5e\x76\x5c\x5c\x26\x3b\xf0\x5e\x12\xe5\x40\x25\x09\xed\x42\xa4\x23\x47\x1b\x4a\x22\xf9\x2f\x4b\xf4\xec\xe8\xef\xb2\xc5\x85\x64\xa4\xf5\x61\x18\x83\x17\x9f\x53\xdb\x98\xda\xe1\x71\x6e\xd0\x35\x44\x9c\x33\xf7\x07\x89\xf3\x93\x8e\x2e\xde\x0a\x4f\xad\x0d\xeb\xd7\xbf\x52\x9f\xec\x1a\xdf\xd5\x29\xa5\x3b\xde\xb4\x7f\xb6\x37\x17\xa8\x1a\x63\x38\xda\x64\xd1\xf5\x67\x85\x67\x59\xb5\xf6\x56\x6b\x4d\xb4\x47\x89\x1b\xd9\xdb\x94\x23\x7f\x55\x5f\xb3\x56\xf1\x94\x76\x06\xee\x10\xf1\x92\xd2\x3f\x31\x6c\xe5\x8b\xf2\x53\xce\xa8\x39\x9d\x63\xf8\x90\x2f\x9a\xe6\x92\x9e\x78\x2b\x8e\x92\x38\xe9\x33\x7c\xcb\x07\xce\x34\x70\x86\xbe\x44\x3e\x18\xfc\xbe\x43\xaa\xd0\x18\x0c\x2e\x6d\x3e\xc0\x3e\x9f\x63\x70\x4e\xe2\x99\x93\xe8\xf4\xc7\x99\xbc\xeb\x0f\x1f\xaf\x88\x51\xcd\xce\x85\x77\x31\x3a\xa1\x58\x7c\x4b\xff\xa2\xfd\xd2\x33\x4f\xa3\xed\x11\x9f\xe8\xc0\x47\xa1\x14\x06\x41\xb3\x11\x72\xad\x93\x3d\x6a\xb8\xef\x41\x5c\x4b\xf3\x01\x71\xb0\x84\x42\xda\x72\xff\x2a\x7e\xee\x58\xa1\xd6\x46\xf4\x8f\x18\x0f\xfb\x28\x50\x88\x27\x0e\x03\x47\xb7\x63\x70\x65\x80\x14\x8d\x29\x00\x19\x19\x5d\x98\x06\x65\x4e\x51\xde\x8a\x05\xc7\x54\x50\x5d\x29\xcd\x2c\x29\x78\xe4\x19\x90\xcb\x60\x31\x7b\x71\x8d\x5e\x09\x80\xf3\xdd\xac\x3e\x86\xc1\x26\x84\xb7\x13\x9a\xeb\xf6\x45\x32\x61\xb6\x2e\x57\x3d\xac\x2d\xab\xad\x07\x4e\xca\xcc\x00\xd5\xb2\x53\xcb\x7e\xac\xa9\x6d\xac\xdf\x45\xc6\x98\x4b\x9f\x4b\x94\x67\x80\xbd\x2c\x03\xd2\x15\xbd\xa4\x87\x61\xcc\x93\x62\x2e\xeb\x7b\x3d\xf3\x65\x3f\xcd\xb3\xd2\xb1\x3e\x97\x01\xfb\x4e\x61\x07\x79\xa3\x83\x86\xa4\xd7\x67\xa3\x9b\x85\x7e\x31\xb2\x56\xa7\xf9\x24\x6c\x90\x37\xe3\x22\x81\x6c\x15\x95\xca\x38\x86\x58\x1d\x79\xbf\xd2\x03\x74\x72\x72\xb3\x00\x75\x74\x3b\x33\x60\xf8\x0c\xcc\x3b\x57\x12\xce\xe0\x26\x38\xa9\x36\xec\xb1\x8a\x89\xee\x5e\x1e\x97\xcd\xaf\x54\x74\x34\x67\x67\x34\x59\x23\x3d\x47\xf4\x51\xce\x6c\xbd\x52\xf1\xa1\xec\x31\x22\x89\xf0\x5c\x49\x9f\x26\xd5\x36\xfd\x27\x88\xe2\x76\x94\x26\xdc\x0f\x5d\x18\xc5\xa7\x83\xdd\xe5\x4e\xcf\x5f\xc5\xfb\x45\x72\x1f\x6f\x8e\x93\x2d\x33\xc5\x33\xc3\x74\x4b\xbf\xa5\x62\x02\xb9\xef\x5d\xc3\x07\x3b\xe6\x81\x97\x46\xee\xa5\x6d\x6a\xc6\xc3\x92\x70\xf7\x1b\xc4\x57\x9f\xe3\x5b\xd5\x2e\xab\x13\x8e\x6a\xda\x04\x8c\x07\x9a\x5e\xe1\x22\xf6\x9e\x0c\x67\x26\x83\xcd\xd6\xbd\x99\x74\x9b\x71\xca\xea\x1e\xb7\xcd\x92\xa4\x35\xf7\x36\x76\xb4\xc6\x0b\x13\x2f\x89\xed\x7a\x09\x34\xff\x07\x00\x00\xff\xff\x7a\xaa\xc7\x24\x5d\x05\x00\x00")

func defaultConfigYamlBytes() ([]byte, error) {
	return bindataRead(
		_defaultConfigYaml,
		"default-config.yaml",
	)
}

func defaultConfigYaml() (*asset, error) {
	bytes, err := defaultConfigYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "default-config.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"default-config.yaml": defaultConfigYaml,
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
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"default-config.yaml": {defaultConfigYaml, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
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

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
