// +build integration

package service

import (
	"bytes"
	"context"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	client     ImgAPIClient
	grpcServer *grpc.Server
)

func TestMain(m *testing.M) {
	config, err := LoadConfig("../config.yml")

	if err != nil {
		panic(err)
	}

	code := 0

	defer func() { os.Exit(code) }()

	go startServer(config)

	// Allow server to spin-up
	time.Sleep(time.Second / 2)

	conn, err := grpc.Dial(
		config.Address,
		grpc.WithInsecure(),
		grpc.WithDialer(func(string, time.Duration) (net.Conn, error) {
			return net.Dial(config.Network, config.Address)
		}),
	)

	if err != nil {
		panic(err)
	}

	defer conn.Close()
	defer grpcServer.Stop()

	client = NewImgAPIClient(conn)
	code = m.Run()
}

func startServer(config *Config) {
	server, err := NewServer(config)

	if err != nil {
		panic(err)
	}

	l, err := net.Listen(config.Network, config.Address)

	if err != nil {
		panic(err)
	}

	defer l.Close()

	grpcServer = grpc.NewServer()

	RegisterImgAPIServer(grpcServer, server)

	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}
}

func TestService(t *testing.T) {
	var imageID string

	t.Run("Upload image", func(t *testing.T) {
		b, err := ioutil.ReadFile("../testdata/gopher.png")

		assert.NoError(t, err)

		resp, err := client.Upload(context.Background(), &UploadRequest{
			Data: b,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Id)

		imageID = resp.Id
	})

	t.Run("Download image in multiple formats", func(t *testing.T) {
		formats := []string{"", "png", "gif", "jpeg"}

		for _, format := range formats {
			resp, err := client.Download(context.Background(), &DownloadRequest{
				Id:     imageID,
				Format: format,
			})

			assert.NoError(t, err)
			assert.NotEmpty(t, resp.Data)

			// Test actual returned blob is valid image data
			switch format {
			case "png":
				_, err = png.Decode(bytes.NewBuffer(resp.Data))
			case "gif":
				_, err = gif.Decode(bytes.NewBuffer(resp.Data))
			case "jpeg":
				_, err = jpeg.Decode(bytes.NewBuffer(resp.Data))
			}

			assert.NoError(t, err)
		}
	})

	t.Run("Attempt to download a non-existant image", func(t *testing.T) {
		_, err := client.Download(context.Background(), &DownloadRequest{
			Id: "123",
		})

		assert.NotNil(t, err)
	})

	t.Run("Download image in unsupported format", func(t *testing.T) {
		_, err := client.Download(context.Background(), &DownloadRequest{
			Id:     imageID,
			Format: "xpeg",
		})

		assert.NotNil(t, err)
	})
}
