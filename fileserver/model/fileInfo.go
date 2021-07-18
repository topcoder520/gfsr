package model

type FileInfo struct {
	Name    string //包含扩展名
	Path    string //相对路径
	Size    int
	Mode    string
	ModTime string
	IsDir   bool
}
