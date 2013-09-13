package straightshot

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"appengine"
	"appengine/blobstore"
	"appengine/user"
)

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Internal Server Error")
	c.Errorf("%v", err)
}

var rootTmpl = template.Must(template.ParseFiles("templates/base.html",
	"straightshot/templates/index.html"))

func handleRoot(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uploadURL, err := blobstore.UploadURL(c, "/upload", nil)
	if err != nil {
		serveError(c, w, err)
		return
	}
	//w.Header().Set("Content-Type", "text/html")

	model := make(map[string]interface{})
	model["title"] = "Straightshot"
	model["uploadURL"] = uploadURL

	if err := rootTmpl.Execute(w, model); err != nil {
		serveError(c, w, err)
		return
	}
}

func handleServe(w http.ResponseWriter, r *http.Request) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		serveError(c, w, err)
		return
	}
	file := blobs["file"]
	if len(file) == 0 {
		c.Errorf("no file uploaded")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/serve/?blobKey="+string(file[0].BlobKey), http.StatusFound)
}

func handleGetUploadURL(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uploadURL, err := blobstore.UploadURL(c, "/upload", nil)
	if err != nil {
		serveError(c, w, err)
		return
	}
	io.WriteString(w, uploadURL.String())
}

func handleSecure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	c := appengine.NewContext(r)
	//u := user.Current(c)
	u, _ := user.CurrentOAuth(c, "https://www.googleapis.com/auth/userinfo.email")
	if u == nil {
		url, _ := user.LoginURL(c, "/secure")
		fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
		return
	}
	url, _ := user.LogoutURL(c, "/secure")
	fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}

func init() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/serve/", handleServe)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/api/getuploadurl", handleGetUploadURL)
	http.HandleFunc("/secure", handleSecure)
}
