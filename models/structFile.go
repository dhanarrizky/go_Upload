package models

type MyFile struct {
	Id   int    `uri:"id"`
	Name string `form:"name"`
	Img  string `form:"img"`
}
