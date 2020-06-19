package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	data "traceability/data"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/gorilla/mux"
)

// UpdateArchView handles PATCH requests and updates archview
func (aw *ArchViews) UpdateArchView(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	jsonBody, err := ioutil.ReadAll(r.Body)

	if !ok {
		io.WriteString(rw, `{{"error": "id not found"}}`)
		return
	}

	archView, err := data.FindArchViewByID(id)

	if err != nil {
		http.Error(rw, `{{"error": "architecture view not found"}}`, http.StatusNotFound)
		return
	}
	jsonArch, err := json.Marshal(archView)
	if err != nil {
		http.Error(rw, `{{"error": "architecture view not found"}}`, http.StatusInternalServerError)
		return
	}
	modifiedJSON, err := jsonpatch.MergePatch(jsonArch, jsonBody)
	modifiedArchView := &data.ArchView{}
	err = json.Unmarshal(modifiedJSON, modifiedArchView)
	data.UpdateArchView(*modifiedArchView)
	err = data.ToJSON(modifiedArchView, rw)
}
