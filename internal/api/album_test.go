package api

import (
	"net/http"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/stretchr/testify/assert"
)

func TestGetAlbums(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		GetAlbums(router, conf)
		r := PerformRequest(app, "GET", "/api/v1/albums?count=10")
		count := gjson.Get(r.Body.String(), "#")
		assert.LessOrEqual(t, int64(3), count.Int())
		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("invalid request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		GetAlbums(router, conf)
		r := PerformRequest(app, "GET", "/api/v1/albums?xxx=10")
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestGetAlbum(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		GetAlbum(router, conf)
		r := PerformRequest(app, "GET", "/api/v1/albums/at9lxuqxpogaaba8")
		val := gjson.Get(r.Body.String(), "Slug")
		assert.Equal(t, "holiday-2030", val.String())
		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("invalid request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		GetAlbum(router, conf)
		r := PerformRequest(app, "GET", "/api/v1/albums/999000")
		val := gjson.Get(r.Body.String(), "error")
		assert.Equal(t, "Album not found", val.String())
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
}

func TestCreateAlbum(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		CreateAlbum(router, conf)
		r := PerformRequestWithBody(app, "POST", "/api/v1/albums", `{"Title": "New created album", "Notes": "", "Favorite": true}`)
		val := gjson.Get(r.Body.String(), "Slug")
		assert.Equal(t, "new-created-album", val.String())
		val2 := gjson.Get(r.Body.String(), "Favorite")
		assert.Equal(t, "true", val2.String())
		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("invalid request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		CreateAlbum(router, conf)
		r := PerformRequestWithBody(app, "POST", "/api/v1/albums", `{"Title": 333, "Description": "Created via unit test", "Notes": "", "Favorite": true}`)
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}
func TestUpdateAlbum(t *testing.T) {
	app, router, conf := NewApiTest()
	CreateAlbum(router, conf)
	r := PerformRequestWithBody(app, "POST", "/api/v1/albums", `{"Title": "Update", "Description": "To be updated", "Notes": "", "Favorite": true}`)
	assert.Equal(t, http.StatusOK, r.Code)
	uid := gjson.Get(r.Body.String(), "UID").String()

	t.Run("successful request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		UpdateAlbum(router, conf)
		r := PerformRequestWithBody(app, "PUT", "/api/v1/albums/"+uid, `{"Title": "Updated01", "Notes": "", "Favorite": false}`)
		val := gjson.Get(r.Body.String(), "Slug")
		assert.Equal(t, "updated01", val.String())
		val2 := gjson.Get(r.Body.String(), "Favorite")
		assert.Equal(t, "false", val2.String())
		assert.Equal(t, http.StatusOK, r.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		UpdateAlbum(router, conf)
		r := PerformRequestWithBody(app, "PUT", "/api/v1/albums"+uid, `{"Title": 333, "Description": "Created via unit test", "Notes": "", "Favorite": true}`)
		assert.Equal(t, http.StatusNotFound, r.Code)
	})

	t.Run("not found", func(t *testing.T) {
		app, router, conf := NewApiTest()
		UpdateAlbum(router, conf)
		r := PerformRequestWithBody(app, "PUT", "/api/v1/albums/xxx", `{"Title": "Update03", "Description": "Created via unit test", "Notes": "", "Favorite": true}`)
		val := gjson.Get(r.Body.String(), "error")
		assert.Equal(t, "Album not found", val.String())
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
}
func TestDeleteAlbum(t *testing.T) {
	app, router, conf := NewApiTest()
	CreateAlbum(router, conf)
	r := PerformRequestWithBody(app, "POST", "/api/v1/albums", `{"Title": "Delete", "Description": "To be deleted", "Notes": "", "Favorite": true}`)
	assert.Equal(t, http.StatusOK, r.Code)
	uid := gjson.Get(r.Body.String(), "UID").String()

	t.Run("delete existing album", func(t *testing.T) {
		app, router, conf := NewApiTest()
		DeleteAlbum(router, conf)
		r := PerformRequest(app, "DELETE", "/api/v1/albums/"+uid)
		assert.Equal(t, http.StatusOK, r.Code)
		val := gjson.Get(r.Body.String(), "Slug")
		assert.Equal(t, "delete", val.String())
		GetAlbums(router, conf)
		r2 := PerformRequest(app, "GET", "/api/v1/albums/"+uid)
		assert.Equal(t, http.StatusNotFound, r2.Code)
	})
	t.Run("delete not existing album", func(t *testing.T) {
		app, router, conf := NewApiTest()
		DeleteAlbum(router, conf)
		r := PerformRequest(app, "DELETE", "/api/v1/albums/999000")
		val := gjson.Get(r.Body.String(), "error")
		assert.Equal(t, "Album not found", val.String())
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
}

func TestLikeAlbum(t *testing.T) {
	t.Run("like not existing album", func(t *testing.T) {
		app, router, ctx := NewApiTest()

		LikeAlbum(router, ctx)

		r := PerformRequest(app, "POST", "/api/v1/albums/xxx/like")
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
	t.Run("like existing album", func(t *testing.T) {
		app, router, ctx := NewApiTest()

		LikeAlbum(router, ctx)
		r := PerformRequest(app, "POST", "/api/v1/albums/at9lxuqxpogaaba7/like")
		assert.Equal(t, http.StatusOK, r.Code)
		GetAlbum(router, ctx)
		r2 := PerformRequest(app, "GET", "/api/v1/albums/at9lxuqxpogaaba7")
		val := gjson.Get(r2.Body.String(), "Favorite")
		assert.Equal(t, "true", val.String())
	})
}

func TestDislikeAlbum(t *testing.T) {
	t.Run("dislike not existing album", func(t *testing.T) {
		app, router, conf := NewApiTest()

		DislikeAlbum(router, conf)

		r := PerformRequest(app, "DELETE", "/api/v1/albums/5678/like")
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
	t.Run("dislike existing album", func(t *testing.T) {
		app, router, conf := NewApiTest()

		DislikeAlbum(router, conf)

		r := PerformRequest(app, "DELETE", "/api/v1/albums/at9lxuqxpogaaba8/like")
		assert.Equal(t, http.StatusOK, r.Code)
		GetAlbum(router, conf)
		r2 := PerformRequest(app, "GET", "/api/v1/albums/at9lxuqxpogaaba8")
		val := gjson.Get(r2.Body.String(), "Favorite")
		assert.Equal(t, "false", val.String())
	})
}

func TestAddPhotosToAlbum(t *testing.T) {
	app, router, conf := NewApiTest()
	CreateAlbum(router, conf)
	r := PerformRequestWithBody(app, "POST", "/api/v1/albums", `{"Title": "Add photos", "Description": "", "Notes": "", "Favorite": true}`)
	assert.Equal(t, http.StatusOK, r.Code)
	uid := gjson.Get(r.Body.String(), "UID").String()

	t.Run("successful request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		AddPhotosToAlbum(router, conf)
		r := PerformRequestWithBody(app, "POST", "/api/v1/albums/"+uid+"/photos", `{"photos": ["pt9jtdre2lvl0y12", "pt9jtdre2lvl0y11"]}`)
		val := gjson.Get(r.Body.String(), "message")
		assert.Equal(t, "photos added to album", val.String())
		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("add one photo to album", func(t *testing.T) {
		app, router, conf := NewApiTest()
		AddPhotosToAlbum(router, conf)
		r := PerformRequestWithBody(app, "POST", "/api/v1/albums/"+uid+"/photos", `{"photos": ["pt9jtdre2lvl0y12"]}`)
		val := gjson.Get(r.Body.String(), "message")
		assert.Equal(t, "photos added to album", val.String())
		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("invalid request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		AddPhotosToAlbum(router, conf)
		r := PerformRequestWithBody(app, "POST", "/api/v1/albums/"+uid+"/photos", `{"photos": [123, "pt9jtdre2lvl0yxx"]}`)
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
	t.Run("not found", func(t *testing.T) {
		app, router, conf := NewApiTest()
		AddPhotosToAlbum(router, conf)
		r := PerformRequestWithBody(app, "POST", "/api/v1/albums/xxx/photos", `{"photos": ["pt9jtdre2lvl0yxx"]}`)
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
}

func TestRemovePhotosFromAlbum(t *testing.T) {
	app, router, conf := NewApiTest()
	CreateAlbum(router, conf)
	r := PerformRequestWithBody(app, "POST", "/api/v1/albums", `{"Title": "Remove photos", "Description": "", "Notes": "", "Favorite": true}`)
	assert.Equal(t, http.StatusOK, r.Code)
	uid := gjson.Get(r.Body.String(), "UID").String()
	AddPhotosToAlbum(router, conf)
	r2 := PerformRequestWithBody(app, "POST", "/api/v1/albums/"+uid+"/photos", `{"photos": ["pt9jtdre2lvl0y12", "pt9jtdre2lvl0y11"]}`)
	assert.Equal(t, http.StatusOK, r2.Code)

	t.Run("successful request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		RemovePhotosFromAlbum(router, conf)
		r := PerformRequestWithBody(app, "DELETE", "/api/v1/albums/"+uid+"/photos", `{"photos": ["pt9jtdre2lvl0y12", "pt9jtdre2lvl0y11"]}`)
		val := gjson.Get(r.Body.String(), "message")
		assert.Equal(t, "entries removed from album", val.String())
		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("no photos selected", func(t *testing.T) {
		app, router, conf := NewApiTest()
		RemovePhotosFromAlbum(router, conf)
		r := PerformRequestWithBody(app, "DELETE", "/api/v1/albums/at9lxuqxpogaaba7/photos", `{"photos": []}`)
		val := gjson.Get(r.Body.String(), "error")
		assert.Equal(t, "No photos selected", val.String())
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
	t.Run("invalid request", func(t *testing.T) {
		app, router, conf := NewApiTest()
		RemovePhotosFromAlbum(router, conf)
		r := PerformRequestWithBody(app, "DELETE", "/api/v1/albums/"+uid+"/photos", `{"photos": [123, "pt9jtdre2lvl0yxx"]}`)
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
	t.Run("album not found", func(t *testing.T) {
		app, router, conf := NewApiTest()
		RemovePhotosFromAlbum(router, conf)
		r := PerformRequestWithBody(app, "DELETE", "/api/v1/albums/xxx/photos", `{"photos": ["pt9jtdre2lvl0yxx"]}`)
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
}

func TestDownloadAlbum(t *testing.T) {
	t.Run("download not existing album", func(t *testing.T) {
		app, router, conf := NewApiTest()

		DownloadAlbum(router, conf)

		r := PerformRequest(app, "GET", "/api/v1/albums/5678/dl?t="+conf.DownloadToken())
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
	t.Run("download existing album", func(t *testing.T) {
		app, router, conf := NewApiTest()

		DownloadAlbum(router, conf)

		r := PerformRequest(app, "GET", "/api/v1/albums/at9lxuqxpogaaba8/dl?t="+conf.DownloadToken())
		assert.Equal(t, http.StatusOK, r.Code)
	})
}

func TestAlbumThumbnail(t *testing.T) {
	t.Run("invalid type", func(t *testing.T) {
		app, router, conf := NewApiTest()
		AlbumThumbnail(router, conf)
		r := PerformRequest(app, "GET", "/api/v1/albums/at9lxuqxpogaaba7/t/"+conf.PreviewToken()+"/xxx")

		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("album has no photo (because is not existing)", func(t *testing.T) {
		app, router, conf := NewApiTest()
		AlbumThumbnail(router, conf)
		r := PerformRequest(app, "GET", "/api/v1/albums/987-986435/t/"+conf.PreviewToken()+"/tile_500")
		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("album: could not find original", func(t *testing.T) {
		app, router, conf := NewApiTest()
		AlbumThumbnail(router, conf)
		r := PerformRequest(app, "GET", "/api/v1/albums/at9lxuqxpogaaba8/t/"+conf.PreviewToken()+"/tile_500")
		assert.Equal(t, http.StatusOK, r.Code)
	})
}
