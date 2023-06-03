package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"encoding/json"
	"os"

 	//"context"
    _ "github.com/godror/godror"
	_ "github.com/mattn/go-sqlite3"

	"github.com/iambenzo/dirtyhttp"
	)

var api dirtyhttp.Api = dirtyhttp.Api{}

// Handler/Controller struct
type httpHandler struct{}

// Implement http.Handler
//
// Logic goes here
func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT,HEAD,DELETE, authorization")
    w.Header().Set("Access-Control-Allow-Headers","Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    w.Header().Set("Access-Control-Expose-Headers", "Authorization")

	switch r.Method {
	case http.MethodGet:
		// get the URL query parameters
		var queryParameters = r.URL.Query()

        if (queryParameters.Get("PostSearch") != "")&&(queryParameters.Get("PostSearchUser") != "") {
            var out []PostSearch
            var err error
             out, err = getSearchPosts(queryParameters.Get("PostSearch"),queryParameters.Get("PostSearchUser"), r.Context())

             if err != nil {
                api.HttpErrorWriter.InternalServerError(w, "Unable to retrieve data from database")
                return
             }

             dirtyhttp.EncodeResponseAsJSON(out, w)
             return
        }else {
			api.HttpErrorWriter.BadRequest(w, "User does not exist")
			return
		}
	case http.MethodPut:

		// get the URL query parameters
		var queryParameters = r.URL.Query()

        if queryParameters.Get("post") != "" {
            // Get user object from request body
            d := json.NewDecoder(r.Body)
            var page PostSearch
            err := d.Decode(&page)
            if err != nil {
                api.Logger.Error(fmt.Sprintf("%v", err))
                api.HttpErrorWriter.InternalServerError(w, "Unable to parse request body")
                return
            }
            // Update our DB
            u, err := editPost(page, r.Context())
            if err != nil {
                api.HttpErrorWriter.InternalServerError(w, "User doesn't exist")
            }

            dirtyhttp.EncodeResponseAsJSON(u, w)
            return

        }else{
            api.HttpErrorWriter.BadParameters(w, "id")
            return
        }


	case http.MethodDelete:
        var queryParameters = r.URL.Query()

        if queryParameters.Get("user") != "" {
            if deleteUser(queryParameters.Get("user"), r.Context())== nil {
                w.WriteHeader(http.StatusNoContent)
                return
            } else {
                api.HttpErrorWriter.BadRequest(w, "User does not exist")
                return
            }
        }else {
			api.HttpErrorWriter.WriteError(w, http.StatusBadRequest, "Please include an 'user_id' or 'content_id' or 'company_page' parameter")
			return
		}

	case http.MethodOptions:
	    w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT,HEAD,DELETE, authorization")
        w.Header().Set("Access-Control-Allow-Headers", "content_type, content-type")
        w.Header().Set("Access-Control-Max-Age", "3600")
        w.WriteHeader(http.StatusNoContent)
    	return


	case http.MethodPost:
		// Get user object from request body
		var queryParameters = r.URL.Query()
		if queryParameters.Get("page") != "" {
            d := json.NewDecoder(r.Body)
            var page PageSearch
            err := d.Decode(&page)
            if err != nil {
                api.Logger.Error(fmt.Sprintf("%v", err))
                api.HttpErrorWriter.InternalServerError(w, "Unable to parse request body")
                return
            }
            out, err:= createPage(page, r.Context())
            // Makes database request
            dirtyhttp.EncodeResponseAsJSON(out, w)

        }
	default:
		// Write a timestamped log entry
		api.Logger.Error("A non-implemented method was attempted")



		// Return a pre-defined error with a custom message
		api.HttpErrorWriter.MethodNotAllowed(w, "Naughty, naughty.")
		return
	}
}
func getDbConnection(cnf *dirtyhttp.EnvConfig) *sql.DB {
	var db *sql.DB
	var err error

	if cnf.DbUrl == "" {
		db, err = sql.Open("sqlite3", "./Nexus.Sqlite")
		api.Logger.Info("Using local SQLite DB")
	}else{
	    os.Setenv("TNS_ADMIN", "./Wallet_LoadoutDatabaseOne/ewallet.p12")
	    api.Logger.Info("Trying to remote the Oracle DB")
        connStr := "admin:Patience1973!@adb.uk-london-1.oraclecloud.com:1522/loadoutdatabaseone_high?wallet_location=/Wallet_LoadoutDatabaseOne/ewallet.p12&wallet_password=SavageM0nk3y!"
        db, err = sql.Open("godror", connStr)
        if err != nil {
            fmt.Println("Error connecting to the database:", err)
            return nil
        }

	}

	if err != nil {
		api.Logger.Fatal("Couldn't connect to database")
	}

	return db
}
func main() {
	// Initialisation
    // Use custom config to remove auth
	config := dirtyhttp.EnvConfig{}
    config.ApiPort = "8080" // change port here
    api.InitWithConfig(&config)

    // set up DB connection
    api.Upstream.SetDatabase(getDbConnection(api.Config))


	// Register a handler
	handler := &httpHandler{}
	api.RegisterHandler("/", *handler)

	// Go, baby, go!
    api.StartServiceNoAuth()
}
