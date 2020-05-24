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

// ListAll handles GET requests and returns all current users
func (u *Users) ListAll(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("[DEBUG] get all records")
	rw.Header().Add("Content-Type", "application/json")

	users := data.GetAllUsers()

	err := data.ToJSON(users, rw)
	if err != nil {
		// we should never be here but log the error just incase
		u.l.Println("[ERROR] serializing product", err)
	}
}

// GetUser handles GET requests and returns the user with ID
func (u *Users) GetUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		io.WriteString(rw, `{{"error": "id not found"}}`)
		return
	}

	user, err := data.FindUserByID(id)

	if err != nil {
		io.WriteString(rw, `{{"error": "user not found"}}`)
		return
	}

	err = data.ToJSON(user, rw)
}
