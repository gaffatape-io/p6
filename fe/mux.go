package fe

import (
	"net/http"
)

func NewMux(development bool) *http.ServeMux {
	mux := http.NewServeMux()			
	mux.Handle("/", &IndexPage{CacheTemplate:!development})
	return mux
}
