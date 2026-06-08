package config

import (
	"log"
	"os"
	"strconv"

	"github.com/lucas-woo/godotenv"
)

var (
	UserDatabaseName string = "user_database"
	UserCollectionName string = "user_credentials";

	CookieSessionIDString string
	CookieUserId string
	CookieSessionMaxAge int
	CookieSessionPath string
	CookieSessionDomain string
	CookieSessionSecure bool
	CookieSessionHttpOnly bool 	
)

func InitCookiesEnv() {
	userId, found := os.LookupEnv("COOKIE_USER_ID")
	if !found {
		log.Fatal("error getting cookie env")
	}
	CookieUserId = userId;	

	sessID, found := os.LookupEnv("COOKIE_SESSION_ID_STRING")
	if !found {
		log.Fatal("error getting cookie env")
	}
	CookieSessionIDString = sessID;

	sessMaxAge, found := os.LookupEnv("COOKIE_SESSION_MAX_AGE")
	if !found {
		log.Fatal("error getting cookie env")
	}
	sessMaxAgeInt, err := strconv.Atoi(sessMaxAge)
	if err != nil {
		log.Fatal("error getting cookie env")
	}
	CookieSessionMaxAge = sessMaxAgeInt;	

	sessPath, found := os.LookupEnv("COOKIE_SESSION_PATH")
	if !found {
		log.Fatal("error getting cookie env")
	}
	CookieSessionPath = sessPath;	

	sessDomain, found := os.LookupEnv("COOKIE_SESSION_DOMAIN")
	if !found {
		log.Fatal("error getting cookie env")
	}
	CookieSessionDomain = sessDomain;	

	sessSecure, found := os.LookupEnv("COOKIE_SESSION_SECURE")
	if !found {
		log.Fatal("error getting cookie env")
	}
	sessSecureBool, err := strconv.ParseBool(sessSecure)
	if err != nil {
		log.Fatal("error getting cookie env")
	}	
	CookieSessionSecure = sessSecureBool;		

	sessHttpOnly, found := os.LookupEnv("COOKIE_SESSION_HTTP_ONLY")
	if !found {
		log.Fatal("error getting cookie env")
	}
	sessHttpOnlyBool, err := strconv.ParseBool(sessHttpOnly)
	if err != nil {
		log.Fatal("error getting cookie env")
	}	
	CookieSessionHttpOnly = sessHttpOnlyBool;		
}

func InitializeEnv() error {
	return godotenv.LoadEnv()
}