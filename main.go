package main


import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"Student_managment/Project/handler"
	"Student_managment/Project/storage/postgres"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"github.com/justinas/nosurf"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
)

var sessionManager *scs.SessionManager



func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	decoder := form.NewDecoder()

	postgreStorage, err := postgres.NewPostgresStorage(config)
	if err != nil {
		log.Fatalln(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
        log.Fatalln(err)
    }

    if err := goose.Up(postgreStorage.DB.DB, "migrations"); err != nil {
        log.Fatalln(err)
    }

	lt := config.GetDuration("session.lifetime")
	it := config.GetDuration("session.idletime")
	sessionManager = scs.New()
	sessionManager.Lifetime = lt * time.Hour
	sessionManager.IdleTimeout = it * time.Minute
	sessionManager.Cookie.Name = "web-session"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = true
	sessionManager.Store = NewSQLXStore(postgreStorage.DB)

	chi := handler.NewHandler(sessionManager, decoder, postgreStorage)
	p := config.GetInt("server.port")
	nosurfHandler := nosurf.New(chi)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", p), nosurfHandler); err != nil {
		log.Fatalf("%#v", err)
	}
}
