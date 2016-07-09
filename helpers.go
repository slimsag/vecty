// +build !js

package vecty

import (
	"bytes"
	"encoding/gob"
	"html/template"
	"io"
)

type Page struct {
	Title   string
	Payload interface{}
	Script  string
}

type internalPayload struct {
	*Page
	GobPayload []byte
}

var t = template.Must(template.New("").Parse(`
	<html>
	    <head>
	        <title>{{.Title}}</title>
	        <script>
	            window.Payload = {{.GobPayload}};
	        </script>
	        {{with .Script}}
		        <script src="{{.Script}}"></script>
	        {{else}}
		        <script src="/assets/frontend/frontend.js"></script>
	        {{end}}
	    </head>
	</html>
`))

func RenderPage(w io.Writer, p *Page) {
	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).Encode(&p.Payload); err != nil {
		panic(err)
	}
	if err := t.Execute(w, &internalPayload{p, buf.Bytes()}); err != nil {
		panic(err)
	}
}
