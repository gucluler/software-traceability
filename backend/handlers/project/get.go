package handlers

import (
	"io"
	"net/http"

	data "traceability/data"

	"github.com/gorilla/mux"
)

// swagger:route GET /users listUsers
// Return a list of users from the database
// responses:
// 200: usersResponse

// ListAll handles GET requests and returns all current projects
func (p *Projects) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] get all records")

	vars := mux.Vars(r)
	userID, ok := vars["userID"]

	if !ok {
		io.WriteString(rw, `{{"error": "id not found"}}`)
		return
	}
	projects := data.GetAllUserProjects(userID)

	err := data.ToJSON(projects, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing project", err)
	}
}

// GetProject handles GET requests and returns the project by ID
func (p *Projects) GetProject(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["projectID"]
	if !ok {
		io.WriteString(rw, `{{"error": "id not found"}}`)
		return
	}

	project, err := data.FindProjectByID(id)

	if err != nil {
		io.WriteString(rw, `{{"error": "user not found"}}`)
		return
	}

	err = data.ToJSON(project, rw)
}
