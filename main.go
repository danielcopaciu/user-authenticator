package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("User Manager", "Application that manages users")

	port := app.String(cli.StringOpt{
		Name:   "port",
		Desc:   "Port on which the application is running",
		EnvVar: "PORT",
	})

	dbUsername := app.String(cli.StringOpt{
		Name:   "dbUsername",
		Desc:   "Database username",
		EnvVar: "DB_USERNAME",
	})

	dbPassword := app.String(cli.StringOpt{
		Name:   "dbPassword",
		Desc:   "Database password",
		EnvVar: "DB_USERNAME",
	})

	dbName := app.String(cli.StringOpt{
		Name:   "dbUsername",
		Desc:   "Database username",
		EnvVar: "DB_USERNAME",
	})

	rw, err := newUserReaderWriter(*dbUsername, *dbPassword, *dbName)
	if err != nil {
		panic(err)
	}

	userManager := newUserManager(rw)
	handler := newHTTPHandler(userManager)

	app.Action = func() {
		go listen(*port, handler)
		waitForSignal()
		rw.close()
	}

	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func listen(port string, handler httpHandler) {
	r := mux.NewRouter()
	r.HandleFunc("/login", handler.doAuthentication).Methods("POST")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}

}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
