package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func WriteTo(w io.Writer, v interface{}) error {
	b, err := MarshalJSON(v)
	fmt.Fprint(w, template.HTML(b))
	return err
}

func MarshalJSON(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Print("unable to marshal json")
		return nil, err
	}
	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	return b, nil
}

func GetJSON(request string, i interface{}) error {
	resp, err := http.Get(request)

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Issue reading json")
		return err
	}
	defer resp.Body.Close()
	return json.Unmarshal(contents, &i)
}
