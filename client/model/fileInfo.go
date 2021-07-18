package model

type FileInfo struct {
	Name    string //包含扩展名
	Path    string //相对路径
	Size    float64
	Mode    string
	ModTime string
	IsDir   bool
}
