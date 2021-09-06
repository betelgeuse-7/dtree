/*
 Get tree representation of the current directory as json data.
*/
package dtree

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	DIRECTORY = "directory"
	FILE      = "file"
)

type fsItem struct {
	Name  string    `json:"name"`
	IsDir bool      `json:"is_dir"`
	Path  string    `json:"path"`
	Items []*fsItem `json:"items"`
	Bytes int64     `json:"bytes"`
}

type fsItems []fsItem

func (fis fsItems) JSON() (string, error) {
	xs, err := json.MarshalIndent(fis, "", " ")
	if err != nil {
		return "", err
	}
	return string(xs), nil
}

func ReadDir(dir string) (fsItems, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return []fsItem{}, err
	}
	var res []fsItem
	for _, f := range fileInfos {
		var r fsItem

		r.Name = f.Name()
		r.IsDir = f.IsDir()
		r.Bytes = f.Size()
		if r.IsDir {
			r.Path = dir + r.Name
		} else {
			r.Path = fmt.Sprintf("%s%s", dir, r.Name)
		}
		if r.IsDir {
			r.setItems()
		}
		res = append(res, r)
	}
	return res, nil
}

// populate f.Items
func (f *fsItem) setItems() {
	fileInfos, err := ioutil.ReadDir(f.Path)
	if err != nil {
		log.Fatalln(err)
	}
	for _, fi := range fileInfos {
		fItem := fsItem{
			Name: fi.Name(),
			Path: fmt.Sprintf("%s/%s", f.Path, fi.Name()),
		}
		fItem.IsDir = fi.IsDir()
		if fItem.IsDir {
			fItem.setItems()
		}
		f.Items = append(f.Items, &fItem)
	}
}
