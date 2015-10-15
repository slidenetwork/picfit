package application

import (
	"encoding/json"
	"github.com/mholt/binding"
	"github.com/thoas/muxer"
	"net/http"
)

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 not found", http.StatusNotFound)
	})
}

type Handler func(muxer.Response, *Request, *Application)

var ImageHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, true, true)

	if err != nil {
		panic(err)
	}

	res.SetHeaders(file.Headers, true)
	res.ResponseWriter.Write(file.Content())
}

var UploadHandler = func(res muxer.Response, req *http.Request, app *Application) {
	if !app.EnableUpload {
		res.Forbidden()
		return
	}

	if app.SourceStorage == nil {
		res.Abort(500, "Your application doesn't have a source storage")
		return
	}

	var err error

	multipartForm := new(MultipartForm)
	errs := binding.Bind(req, multipartForm)
	if errs.Handle(res) {
		return
	}

	file, err := multipartForm.Upload(app.SourceStorage)

	if err != nil {
		panic(err)
	}

	content, err := json.Marshal(map[string]string{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})

	if err != nil {
		panic(err)
	}

	res.Header().Add("Location", file.URL())
	res.ContentType("application/json")
	res.ResponseWriter.WriteHeader(http.StatusCreated)
	res.ResponseWriter.Write(content)
}

var DeleteHandler = func(res muxer.Response, req *muxer.Request, app *Application) {
	if app.SourceStorage == nil {
		res.Abort(500, "Your application doesn't have a source storage")
		return
	}

	filename := req.Params["path"]

	// Delete the image from the source storage right away.
	// Upcoming requests will be able to read stuff from cache, but never
	// create new things after we've deleted the source image.
	//
	// Don't care about an error here. If it's not there, then whatever..
	// Keep going...
	app.Logger.Infof("Deleting source storage file: %s", filename)
	_ = app.SourceStorage.Delete(filename)

	// We can clean up the rest now. Again, don't care about errors.
	// ImageCleanup always succeeds.
	app.ImageCleanup(filename)

	content, err := json.Marshal(map[string]string{
		"filename": filename,
	})

	if err != nil {
		panic(err)
	}

	res.ContentType("application/json")
	res.ResponseWriter.Write(content)
}

var GetHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, false, false)

	if err != nil {
		panic(err)
	}

	content, err := json.Marshal(map[string]string{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})

	if err != nil {
		panic(err)
	}

	res.ContentType("application/json")
	res.ResponseWriter.Write(content)
}

var RedirectHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, false, false)

	if err != nil {
		panic(err)
	}

	res.PermanentRedirect(file.URL())
}
