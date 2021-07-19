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
	return ParseResult(rs)
}

func HandleCdCmd(cmd, cdPath string) (*model.Message, error) {
	rs, err := Get(path.Join("/api/handlecmd/", cmd) + "?cdpath=" + cdPath)
	if err != nil {
		return nil, err
	}
	return ParseResult(rs)
}
