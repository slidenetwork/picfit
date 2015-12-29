package application

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"

	"github.com/mholt/binding"
	"github.com/slidenetwork/picfit/image"
	"github.com/thoas/gostorages"
)

type MultipartForm struct {
	Data *multipart.FileHeader `json:"data"`
}

func (f *MultipartForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&f.Data: "data",
	}
}

func (f *MultipartForm) Upload(storage gostorages.Storage) (*image.ImageFile, error) {
	var fh io.ReadCloser

	fh, err := f.Data.Open()

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	dataBytes := bytes.Buffer{}

	_, err = dataBytes.ReadFrom(fh)

	if err != nil {
		return nil, err
	}

	dataHash := fmt.Sprintf("%x", md5.Sum(dataBytes.Bytes()))
	ext := filepath.Ext(f.Data.Filename)
	filename := path.Join(dataHash[:2], dataHash[2:4], dataHash[4:]+ext)

	err = storage.Save(filename, gostorages.NewContentFile(dataBytes.Bytes()))

	if err != nil {
		return nil, err
	}

	return &image.ImageFile{
		Filepath: filename,
		Storage:  storage,
	}, nil
}
