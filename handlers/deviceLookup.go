package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria-go/index"
)

// DeviceLookup is a HTTP handler that is used to do the lookup rest api.
//
// It forwards the request to the query endpoint.
type DeviceLookup struct {
	store *index.Store
}

// NewDeviceLookupHandler creates a HTTP handler for device lookup
func NewDeviceLookupHandler(store *index.Store) *DeviceLookup {
	return &DeviceLookup{store: store}
}

// CabHandler handles the
func (d *DeviceLookup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	device := vars["device"]
	component, foundComponent := vars["component"]
	register, foundRegister := vars["register"]
	bitfield, foundBitfield := vars["bitfield"]

	query := fmt.Sprintf("atmel;device:%v", device)
	if foundComponent {
		query += fmt.Sprintf(";comp:%v", component)
	}
	if foundRegister {
		query += fmt.Sprintf(";register:%v", register)
	}
	if foundBitfield {
		query += fmt.Sprintf(";bitfield:%v", bitfield)
	}

	keywordResults := d.store.LookupKeyword(query)

	log.Printf("[device-lookup] lookup for %s/%s/%s/%s returned %v", device, component, register, bitfield, keywordResults)
	log.Printf("[device-lookup] based on %s", query)

	if len(keywordResults) == 0 {
		http.Error(w, fmt.Sprintf("No results for query '%v'", query), http.StatusNotFound)
	} else {
		result := keywordResults[0]
		url := fmt.Sprintf("http://content.alexandria.atmel.com/webhelp/%v/index.html?%v", result.BookID, strings.TrimSuffix(result.Filename, filepath.Ext(result.Filename)))
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
