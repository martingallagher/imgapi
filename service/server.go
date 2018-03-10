package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/martingallagher/imgapi"
	"github.com/pkg/errors"
)

type server struct {
	config *Config
}

// NewServer returns a new ImgAPI server.
func NewServer(c *Config) (ImgAPIServer, error) {
	_, err := os.Stat(c.DataDir)

	if os.IsNotExist(err) {
		err = os.Mkdir(c.DataDir, 0700)
	}

	if err != nil {
		return nil, err
	}

	return &server{c}, nil
}

func (s *server) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	i, err := imgapi.NewImage(req.Data)

	if err != nil {
		return nil, err
	}

	name := imgapi.Filename(i.ID, s.config.DirDepth)
	dir := s.config.DataDir + "/" + name[:strings.LastIndexByte(name, '/')]
	_, err = os.Stat(dir)

	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
	}

	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(s.config.DataDir+"/"+name, i.Data, 0600)

	if err != nil {
		return nil, err
	}

	return &UploadResponse{
		Id: hex.EncodeToString(i.ID[:]),
	}, nil
}

func (s *server) Download(ctx context.Context, req *DownloadRequest) (*Image, error) {
	v, err := hex.DecodeString(req.Id)

	if err != nil {
		return nil, err
	}

	name := imgapi.Filename(v, s.config.DirDepth)
	b, err := ioutil.ReadFile(s.config.DataDir + "/" + name)

	if err != nil {
		return nil, err
	}

	format, err := imgapi.ImageFormat(b)

	if err != nil {
		return nil, err
	}

	// Handle different output formats
	req.Format = strings.ToLower(req.Format)

	if req.Format != "" && req.Format != format {
		format = req.Format
		b, err = handleImageFormat(format, b)

		if err != nil {
			return nil, err
		}
	}

	return &Image{
		Id:     req.Id,
		Data:   b,
		Format: format,
	}, nil
}

func handleImageFormat(format string, b []byte) ([]byte, error) {
	if format != "png" && format != "gif" && format != "jpeg" {
		return nil, errors.Wrap(image.ErrFormat, format)
	}

	i, _, err := image.Decode(bytes.NewBuffer(b))

	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	switch format {
	case "png":
		err = png.Encode(buf, i)
	case "gif":
		err = gif.Encode(buf, i, nil)
	case "jpeg":
		err = jpeg.Encode(buf, i, nil)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
