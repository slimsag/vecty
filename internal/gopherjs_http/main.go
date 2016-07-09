// Package gopherjs_http provides helpers for compiling Go using GopherJS and serving it over HTTP.
package gopherjs_http

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"sync"

	"github.com/shurcooL/gopherjslib"
)

func handleJsError(jsCode string, err error) string {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return `console.error("` + template.JSEscapeString(err.Error()) + `");`
	}
	return jsCode
}

// Needed to prevent race condition until https://github.com/go-on/gopherjslib/issues/2 is resolved.
var gopherjslibLock sync.Mutex

func goReadersToJs(names []string, goReaders []io.Reader) (jsCode string, err error) {
	gopherjslibLock.Lock()
	defer gopherjslibLock.Unlock()

	var out bytes.Buffer
	builder := gopherjslib.NewBuilder(&out, &gopherjslib.Options{Minify: true})

	for i, goReader := range goReaders {
		builder.Add(names[i], goReader)
	}

	err = builder.Build()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
