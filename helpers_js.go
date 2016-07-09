package vecty

import (
	"encoding/base64"
	"encoding/gob"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

func Payload() interface{} {
	gobPayload := js.Global.Get("window").Get("Payload").String()

	var x interface{}
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(gobPayload))
	if err := gob.NewDecoder(dec).Decode(&x); err != nil {
		panic(err)
	}

	return x
}
