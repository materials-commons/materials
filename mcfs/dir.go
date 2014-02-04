package mcfs

import (
	"fmt"
	"github.com/materials-commons/base/mc"
	"github.com/materials-commons/gohandy/file"
	"github.com/materials-commons/mcfs/protocol"
	"strings"
)

func (c *Client) CreateDir(projectID, projectName, path string) (dataDirID string, err error) {
	i := strings.Index(path, projectName)
	if i == -1 {
		return "", mc.ErrInvalid
	}

	properPath := path[i:] // only send up portion starting from project
	req := protocol.CreateDirReq{
		ProjectID: projectID,
		Path:      file.NormalizePath(properPath),
	}

	resp, err := c.doRequest(req)
	if resp == nil {
		return "", err
	}

	switch t := resp.(type) {
	case protocol.CreateResp:
		return t.ID, err
	default:
		fmt.Printf("2 %s %T\n", ErrBadResponseType, t)
		return "", ErrBadResponseType
	}
}

func (c *Client) UploadDirectory(projectID, dataDirID string, path string) ([]DataFileUpload, error) {
	switch {
	case dataDirID == "":
		return c.uploadNewDirectory(projectID, path)
	default:
		return c.uploadExistingDirectory(projectID, dataDirID, path)
	}
}

func (c *Client) uploadNewDirectory(projectID, path string) ([]DataFileUpload, error) {
	return nil, nil
}

func (c *Client) uploadExistingDirectory(projectID, dataDirID, path string) ([]DataFileUpload, error) {

	return nil, nil
}
