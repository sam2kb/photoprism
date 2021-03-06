package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetDownloadFileID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := SetDownloadFileID("exampleFileName.jpg", 1000000)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("filename empty", func(t *testing.T) {
		err := SetDownloadFileID("", 1000000)
		if err == nil {
			t.Fatal()
		}
		assert.Equal(t, "sync: can't update, filename empty", err.Error())
	})
}
