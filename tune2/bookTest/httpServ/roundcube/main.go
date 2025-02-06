package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type loginFunc = func(w http.ResponseWriter, r *http.Request, log *zap.Logger)

func funcLogHandler(log *zap.Logger, l loginFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l(w, r, log)
	}
}

func login(w http.ResponseWriter, r *http.Request, log *zap.Logger) {
	log.Info("Loggin attempt!",
		zap.String("username", r.FormValue("_user")),
		zap.String("password", r.FormValue("_pass")),
		zap.String("user-agent", r.UserAgent()),
		zap.String("ip_address", r.RemoteAddr),
		zap.String("time", r.FormValue("_timezone")),
	)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {

	f, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"test.txt"}
	l, err := config.Build()
	// l, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	defer l.Sync()

	r := mux.NewRouter()

	r.HandleFunc("/login", funcLogHandler(l, login))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	// err = http.ListenAndServe(":8080", r)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err1 := http.Error`()
	log.Fatal(http.ListenAndServe(":8080", r))

	// http.ListenAndServe("")
}
