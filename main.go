package main

import (
	"github.com/stretchr/gomniauth/providers/google"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
	trace "github.com/youthtrouble/gotrace"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	//var addr = flag.String("addr", ":8080", "The addr of the application.")
	//flag.Parse()
	var addr = os.Getenv("PORT")
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		google.New(os.Getenv("G_KEY"), os.Getenv("G_CLIENT_ID"), "https://apoti.herokuapp.com/auth/callback/google"),
		)
	r := newRoom(avatars)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(307)
	})
	http.Handle("/room", r)
	go r.run()
	// start the web server
	log.Println("Starting web server on", addr)
	if err := http.ListenAndServe(":" + addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
