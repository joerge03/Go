package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"
)

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

	l, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	defer l.Sync()
}
