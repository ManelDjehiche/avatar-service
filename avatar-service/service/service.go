package service

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/emplorium/auth-service/database"
	"github.com/emplorium/auth-service/env"
	"github.com/gorilla/mux"
)

type ID struct {
	id string `json:"id"`
}

type Avatar struct {
	Path       string `json:"path"`
	CreatedAt  int64  `json:"createdAt"`
	ModifiedAt int64  `json:"modifiedAt"`
}

var db *sql.DB
var err error

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func uploadFile(r *http.Request) (string, error) {
	var pathSave = ".\\avatars"

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("avatar")

	if err != nil {
		fmt.Println("error retrieving file from form-data")
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	extension := filepath.Ext(handler.Filename)
	tempFile, err := ioutil.TempFile(pathSave, "*-avatar"+extension)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	destination := tempFile.Name()
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	return destination, nil
}

func Start() {
	db, err = database.StartDatabase()
	if err != nil {
		log.Fatal("Connection to databse faild ", err)
	}
	defer db.Close()

	router := mux.NewRouter()
	fmt.Println("Application Starting ...", env.Settings.AgentService.Host+":"+env.Settings.AgentService.Port)
	if err != nil {
		log.Fatal("Connection to agent domain service faild ", err)
	}

	//account
	router.HandleFunc("/addAccountAvatar", AddAccountAvatar)
	router.HandleFunc("/updateAccountAvatar", UpdateAccountAvatar)
	router.HandleFunc("/deleteAccountAvatar", DeleteAccountAvatar)
	router.HandleFunc("/getAccountAvatar", GetAccountAvatar)

	//agent
	router.HandleFunc("/addAgentAvatar", AddAgentAvatar)
	router.HandleFunc("/updateAgentAvatar", UpdateAgentAvatar)
	router.HandleFunc("/deleteAgentAvatar", DeleteAgentAvatar)
	router.HandleFunc("/getAgentAvatar", GetAgentAvatar)

	//lead
	router.HandleFunc("/addLeadAvatar", AddLeadAvatar)
	router.HandleFunc("/updateLeadAvatar", UpdateLeadAvatar)
	router.HandleFunc("/deleteLeadAvatar", DeleteLeadAvatar)
	router.HandleFunc("/getLeadAvatar", GetAccountAvatar)

	log.Fatal(http.ListenAndServe(":"+env.Settings.HTTPPort, router))
}
