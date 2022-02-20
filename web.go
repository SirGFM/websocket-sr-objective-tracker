package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
)

// server wraps the httpServer and all of its components, so it may be
// gracefully stopped.
type server struct {
	// The server's HTTP server.
	httpServer *http.Server

	// Synchronizes access to the map.
	mutex sync.RWMutex

	// Data maps a given URI (the game) to all the resources for that game.
	data map[string] map[string]interface{}
}

// request wraps every object used to handle a received request.
type request struct {
	// The response object.
	w http.ResponseWriter

	// The request object.
	req *http.Request

	// The clean up request URI.
	uri []string

	// The server that received this request.
	s *server
}

// Close the running web server and clean up resourcers.
func (s *server) Close() error {
	if s.httpServer != nil {
		s.mutex.Lock()
		s.httpServer.Close()
		s.mutex.Unlock()
		s.httpServer = nil
	}

	return nil
}

// get_tracker encode the requested game as a JSON and returns it. If a
// specific field is requested, its ID is returned alongside its value.
// However, for entire games each ID is a field in the JSON. E.g.:
//
// GET: /tracker/some-game
//
//     {
//         "some-val": "text",
//         "some-num": 0
//     }
//
// GET: /tracker/some-game/some-val
//
//     {
//         "id": "some-val",
//         "value": "text"
//     }
func (r *request) get_tracker() {
	var field interface{}

	key := r.uri[1]

	r.s.mutex.RLock()
	defer r.s.mutex.RUnlock()

	game, ok := r.s.data[key]
	if ok && len(r.uri) > 2 {
		// Encode only the requested field
		id := strings.Join(r.uri[2:], "/")
		field, ok = game[id]

		if ok {
			payload := struct {
				ID string `json:"id"`
				Value interface{} `json:"value"`
			} {
				ID: id,
				Value: field,
			}

			field = &payload
		}
	} else if ok {
		// Encode the entire data
		field = game
	}

	if !ok {
		r.httpTextReply(http.StatusNotFound, "Invalid resource/game")
		log.Printf("[%s] %s - %s: Invalid resource/game", r.req.Method, key, r.req.RemoteAddr)
		return
	}

	data, err := json.Marshal(field)
	if err != nil {
		serr := "Failed to encode the response"
		r.httpTextReply(http.StatusInternalServerError, serr)
		log.Printf("[%s] %s - %s: %s (%+v)", r.req.Method, key, r.req.RemoteAddr, serr, err)
		return
	}

	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(http.StatusOK)
	writeData(data, r.w)
}

// set_value assign a game's, defined by the key, ID to the given value.
func (s *server) set_value(key, id string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	game, ok := s.data[key]
	if !ok {
		game = make(map[string]interface{})
		s.data[key] = game
	}

	log.Printf("[%s]: %s=%v\n", key, id, value)
	game[id] = value

	// TODO: Send this value over a websocket
}

// post_tracker create or update a value for a given game. The value may
// either be specified in the requests body or, if it should be simply set
// to true, it may be specified in the URI. E.g.:
//
// POST: /tracker/some-game
//
//     {
//         "id": "some-val",
//         "value": "text"
//     }
//
// POST: /tracker/some-game/some-bool
//
// This second POST simply sets some-game's some-bool to true.
func (r *request) post_tracker() {
	key := r.uri[1]

	if len(r.uri) == 2 {
		// The game name's is in the URI and the ID/value is in the body.
		var payload struct{
			ID string `json:"id"`
			Value interface{} `json:"value"`
		}
		dec := json.NewDecoder(r.req.Body)
		err := dec.Decode(&payload)
		if err != nil {
			log.Printf("[%s] %s - %s: Failed to parse request: %+v", r.req.Method, key, r.req.RemoteAddr, err)
			r.httpTextReply(http.StatusBadRequest, "Invalid data")
			return
		}

		r.s.set_value(key, payload.ID, payload.Value)
	} else {
		// This is simply assigning the ID in the URI to true.
		id := strings.Join(r.uri[2:], "/")

		r.s.set_value(key, id, true)
	}

	r.reply_no_content()
}

// delete_tracker assign a vlue in the given game to false. Everything
// after the game in the URI is assumed to be the value's name.
func (r *request) delete_tracker() {
	key := r.uri[1]

	if len(r.uri) < 2 {
		log.Printf("[%s] %s - %s: No resources was specified", r.req.Method, key, r.req.RemoteAddr)
		r.httpTextReply(http.StatusBadRequest, "Missing resource")
		return
	}

	id := strings.Join(r.uri[2:], "/")

	r.s.set_value(key, id, false)
	r.reply_no_content()
}

// reply_no_content send a StatusNoContent to the requester.
func (r *request) reply_no_content() {
	// Ensure that the body is emptied before replying.
	io.Copy(io.Discard, r.req.Body)
	r.w.WriteHeader(http.StatusNoContent)
}

// ServeHTTP is called by Go's http package whenever a new HTTP request arrives.
func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	uri := cleanURL(req.URL)
	log.Printf("%s - %s - %s", req.RemoteAddr, req.Method, uri)

	res := strings.Split(uri, "/") 
	r := request {
		w: w,
		req: req,
		uri: res,
		s: s,
	}

	is_tracker := (req.Method == http.MethodPost || req.Method == http.MethodDelete)
	if is_tracker && len(res) < 2 {
		r.httpTextReply(http.StatusNotFound, "No resource was specified (expected at least '/tracker/<some-name>')")
		log.Printf("[%s] %s - %s: 404", req.Method, uri, req.RemoteAddr)
		return
	} else if is_tracker && res[0] != "tracker" {
		r.httpTextReply(http.StatusNotFound, "Invalid resource was specified (expected a '/tracker/<some-name>')")
		log.Printf("[%s] %s - %s: 404", req.Method, uri, req.RemoteAddr)
		return
	} else if len(res) == 0 && req.Method == http.MethodGet {
		// For GET requests, default to page index.html
		r.uri = []string{"index.html"}
	}

	switch req.Method {
	case http.MethodGet:
		if len(r.uri) >= 2 && r.uri[0] == "tracker" {
			r.get_tracker()
		} else {
			// TODO: handle retrieving files
			r.httpTextReply(http.StatusNotImplemented, "Can't get files yet...")
			log.Printf("[%s] %s - %s: 501", req.Method, uri, req.RemoteAddr)
		}
	case http.MethodPost:
		r.post_tracker()
	case http.MethodDelete:
		r.delete_tracker()
	default:
		r.httpTextReply(http.StatusMethodNotAllowed, "Invalid HTTP Method")
		log.Printf("[%s] %s - %s: 405", req.Method, uri, req.RemoteAddr)
	}
}

// cleanURL so everything is properly escaped/encoded and so it may be split into each of its components.
//
// Use `url.Unescape` to retrieve the unescaped path, if so desired.
func cleanURL(uri *url.URL) string {
	// Normalize and strip the URL from its leading prefix (and slash)
	resUrl := path.Clean(uri.EscapedPath())
	if len(resUrl) > 0 && resUrl[0] == '/' {
		resUrl = resUrl[1:]
	} else if len(resUrl) == 1 && resUrl[0] == '.' {
		// Clean converts an empty path into a single "."
		resUrl = ""
	}

	return resUrl
}

// httpTextReply send a simple HTTP response as a plain text.
func (r *request) httpTextReply(status int, msg string) {
	// Ensure that the body is emptied before replying.
	io.Copy(io.Discard, r.req.Body)

	r.w.Header().Set("Content-Type", "text/plain")
	r.w.WriteHeader(status)

	writeData([]byte(msg), r.w)
}

// writeData, account for incomplete writes.
func writeData(data []byte, w io.Writer) {
	for len(data) > 0 {
		n, err := w.Write(data)
		if err != nil {
			log.Printf("Failed to send: %+v", err)
			return
		}
		data = data[n:]
	}
}

// RunWeb starts the web server and return an io.Closer, so the server may
// be stopped.
func RunWeb(args Args) io.Closer {
	var srv server

	srv.httpServer = &http.Server {
		Addr: fmt.Sprintf("%s:%d", args.IP, args.Port),
		Handler: &srv,
	}
	srv.data = make(map[string] map[string]interface{})

	go func() {
		log.Printf("Waiting...")
		srv.httpServer.ListenAndServe()
	} ()

	return &srv
}
