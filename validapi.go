package validapi

import (
	"net/http"
)

type (
	//ValidAPI ...
	ValidAPI struct {
		tree            Router
		NotFoundHandler func(http.ResponseWriter, *http.Request)
		CORS            bool
	}
)
