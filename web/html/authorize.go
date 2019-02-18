package html

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// Authorize authorizes user.
func (ar *Router) Authorize() http.HandlerFunc {
	tmpl, err := template.ParseFiles("../../static/pm.html")

	return func(w http.ResponseWriter, r *http.Request) {
		serveTemplate := func(errorMessage, IDToken, callbackURL string) {
			if err != nil {
				ar.Logger.Printf("Error parsing template: %v", err)
				ar.Error(w, err, 500, "Error parsing template")
				return
			}

			data := map[string]interface{}{
				"Error":       errorMessage,
				"IDToken":     IDToken,
				"CallbackURL": callbackURL,
			}

			if err := tmpl.Execute(w, data); err != nil {
				ar.Logger.Printf("Error executing template: %v", err)
				ar.Error(w, err, 500, "Error executing template")
				return
			}
		}

		context := r.Context()

		if err := appIDErrorFromContext(context); err != "" {
			ar.Logger.Printf(err)
			serveTemplate(err, "", "")
			return
		}

		app := appFromContext(context)
		if app == nil {
			message := "Error getting App"
			ar.Logger.Println(message)
			serveTemplate(message, "", "")
			return
		}

		userID, err := getCookie(r, "identifo-user")
		if err != nil || userID == "" {
			serveTemplate("not authorized", "", app.RedirectURL())
			return
		}

		user, err := ar.UserStorage.UserByID(userID)
		if err != nil {
			ar.Logger.Printf("Error: getting UserByID: %v, userID: %v", err, userID)
			serveTemplate("invalid user token", "", app.RedirectURL())
			return
		}

		scopesJSON := strings.TrimSpace(r.URL.Query().Get("scopes"))
		scopes := []string{}
		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("Error: Invalid scopes %v", scopesJSON)
			serveTemplate("invalid scopes", "", app.RedirectURL())
			return
		}

		scopes, err = ar.UserStorage.RequestScopes(user.ID(), scopes)
		if err != nil {
			ar.Logger.Printf("Error: invalid scopes %v for userID: %v", scopes, user.ID())
			message := fmt.Sprintf("user not allowed to access this scopes %v", scopes)
			serveTemplate(message, "", app.RedirectURL())
			return
		}

		token, err := ar.TokenService.NewToken(user, scopes, app)
		if err != nil {
			ar.Logger.Printf("Error creating token: %v", err)
			serveTemplate("server error", "", app.RedirectURL())
			return
		}

		tokenString, err := ar.TokenService.String(token)
		if err != nil {
			ar.Logger.Printf("Error stringifying token: %v", err)
			serveTemplate("server error", "", app.RedirectURL())
			return
		}

		serveTemplate("", tokenString, app.RedirectURL())
		return
	}
}
