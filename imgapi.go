package imgapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"io"
	// Allow image type detection
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// Image represents an image API image.
type Image struct {
	ID   []byte
	Data []byte
	Type string
}

// NewImage returns a new image from the given blob.
func NewImage(b []byte) (*Image, error) {
	f, err := ImageFormat(b)

	if err != nil {
		return nil, err
	}

	h := sha256.New()
	_, err = h.Write(b)

	if err != nil {
		return nil, err
	}

	return &Image{
		ID:   h.Sum(nil),
		Data: b,
		Type: f,
	}, nil
}

// ImageFormat returns the image format for the given image blob.
func ImageFormat(b []byte) (string, error) {
	_, f, err := image.DecodeConfig(bytes.NewBuffer(b))

	return f, err
}

// Filename generates a filename for the given ID and directory depth.
func Filename(id []byte, dirDepth int) string {
	if dirDepth < 1 {
		return hex.EncodeToString(id)
	}

	l := hex.EncodedLen(sha256.Size)
	v := make([]byte, l)

	hex.Encode(v, id)

	if l > dirDepth {
		l = dirDepth
	}

	for i, n := 1, l*2; i < n; i += 2 {
		v = append(v, 0)

		copy(v[i+1:], v[i:])

		v[i] = '/'
	}

	return string(v)
}

// Save saves the image to the given writer.
func (i *Image) Save(w io.Writer) error {
	_, err := w.Write(i.Data)

	return err
}
