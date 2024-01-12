package models

import "fmt"

type Image struct {
	BaseUrl string
	Project string
	Name    string
	Tag     string
}

func (i *Image) Print() {
	fmt.Println(i.ToString())
}

func (i *Image) ToString() string {
	return i.BaseUrl + "/" + i.Project + "/" + i.Name + "/" + i.Tag
}
