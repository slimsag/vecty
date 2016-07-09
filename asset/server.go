package asset

import (
	"net/http"

	"github.com/gopherjs/vecty/internal/gopherjs_http"
	"github.com/nytimes/gziphandler"
)

type Options struct {
	Dir         string
	StripPrefix string
	Gzip        bool
}

func Defaults() *Options {
	return &Options{
		Dir:         "assets",
		StripPrefix: "/assets/",
		Gzip:        true,
	}
}

func NewServer(o *Options) http.Handler {
	if o == nil {
		o = Defaults()
	}

	assets := gopherjs_http.NewFS(http.Dir(o.Dir))

	h := http.StripPrefix(o.StripPrefix, http.FileServer(assets))
	if o.Gzip {
		h = gziphandler.GzipHandler(h)
	}
	return h
}
