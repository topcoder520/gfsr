package request

import (
	"path"

	"github.com/topcoder520/gfsr/client/model"
)

func GetDirFiles(dir string) (*model.Message, error) {
	rs, err := Get(path.Join("/api/files/", dir))
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(rs))
	return ParseResult(rs)
}
