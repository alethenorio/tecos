package v1

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ByteFlinger/tecos/backend"
	"github.com/go-chi/chi"
)

// Routes returns a v1 Router
func Routes(backend backend.Storage) chi.Router {
	r := chi.NewRouter()
	r.Get("/", ListModules(backend))
	r.Get("/{namespace}", ListModules(backend))

	return r
}

// ListModules lists all modules in the registry
func ListModules(storage backend.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.URL.Query()

		_ = chi.URLParam(r, "namespace")

		list := ModuleList{}

		list.Modules = []ModuleInfo{}

		for _, m := range storage.ListModules() {
			list.Modules = append(list.Modules, FromData(m))
		}

		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(list)
		if err != nil {
			http.Error(w, "500 - Internal Server Error", 500)
			return
		}

		w.Write(b.Bytes())

	}
}
