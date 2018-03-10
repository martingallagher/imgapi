package imgapi

import (
	"crypto/rand"
	"crypto/sha256"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testdata = "./testdata/"
	cacheDir = testdata + "cache/"
)

func TestMain(m *testing.M) {
	code := 0

	defer func() { os.Exit(code) }()
	defer os.RemoveAll(cacheDir)

	code = m.Run()
}

func TestNewImage(t *testing.T) {
	t.Run("Bad path; erroneous image data", func(t *testing.T) {
		b := make([]byte, 512)

		_, err := rand.Read(b)

		assert.NoError(t, err)

		i, err := NewImage(b)

		assert.NotNil(t, err)
		assert.Nil(t, i)
	})

	t.Run("Happy path", func(t *testing.T) {
		files, err := ioutil.ReadDir(testdata)

		assert.NoError(t, err)

		var f *os.File

		defer func() {
			if f != nil {
				f.Close()
			}
		}()

		for _, v := range files {
			if v.IsDir() {
				continue
			}

			t.Logf("Testing image: %s", v.Name())

			b, err := ioutil.ReadFile(testdata + v.Name())

			assert.NoError(t, err)

			i, err := NewImage(b)

			assert.NoError(t, err)
			assert.NotNil(t, i)

			name := Filename(i.ID, 3)
			name = cacheDir + name
			dir := name[:strings.LastIndexByte(name, '/')]
			err = os.MkdirAll(dir, 0755)

			assert.NoError(t, err)

			f, err = os.Create(name)

			assert.NoError(t, err)
			assert.NoError(t, i.Save(f))
		}
	})
}

func TestFilename(t *testing.T) {
	id, err := randomID()

	assert.NoError(t, err)

	for i := 0; i <= 5; i++ {
		name := Filename(id, i)

		t.Logf("Filename (depth %d): %s", i, name)
		assert.Equal(t, i, strings.Count(name, "/"))
	}
}

func randomID() ([]byte, error) {
	id := make([]byte, sha256.Size)
	_, err := rand.Read(id)

	return id, err
}
