package service

import (
	"encoding/json"
	"log"
	"net/http"
)

func AddAccountAvatar(w http.ResponseWriter, r *http.Request) {
	var id int32
	accountId := r.PostFormValue("id")

	if accountId == "" {
		log.Println("accountId is empty")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This accountId is empty"))
		return
	}

	err := db.QueryRow(`
	SELECT id
	FROM
	accounts
	WHERE account_id = $1
	`, accountId).Scan(&id)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This accountId is wrong"))
		return
	}

	err = db.QueryRow(`
	SELECT id
	FROM
	accounts_avatar
	WHERE account_id = $1
	`, accountId).Scan(&id)

	if err == nil {
		log.Println("already exist!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This accountId has already a picture!"))
		return
	}

	path, err := uploadFile(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("This file cant be loaded"))
		return
	}

	now := timestamp()
	err = db.QueryRow(`
		INSERT INTO
		accounts_avatar
		(created_at, modified_at, path, account_id)
		VALUES
		($1, $2, $3, $4)
		RETURNING id
	`,
		now,
		now,
		path,
		accountId).Scan(&id)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This file path cant be saved in the database"))
		return
	}

	avatar := Avatar{
		Path:       path,
		CreatedAt:  now,
		ModifiedAt: now,
	}

	jsonBytes, err := json.Marshal(avatar)

	if err != nil {
		log.Println("err json: ", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonBytes)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonBytes)
	return
}

func GetAccountAvatar(w http.ResponseWriter, r *http.Request) {

	var id int32
	accountId := r.PostFormValue("id")
	var accountAvatar Avatar
	var idAccount string

	err = db.QueryRow(`
		SELECT 
		id, account_id, created_at, modified_at, path
		FROM 
		accounts_avatar 
		WHERE 
		account_id = $1
	`, accountId).Scan(&id, &idAccount, &accountAvatar.CreatedAt, &accountAvatar.ModifiedAt, &accountAvatar.Path)

	if err != nil {
		log.Println("this account doesnt have picture before!")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("this account doesnt have picture before!"))
		return
	}

	jsonBytes, err := json.Marshal(accountAvatar)

	if err != nil {
		log.Println("err json: ", err)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Can't produce json object"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonBytes)
	return
}

func UpdateAccountAvatar(w http.ResponseWriter, r *http.Request) {

	var id int32
	accountId := r.PostFormValue("id")

	err = db.QueryRow(`
	SELECT id
	FROM
	accounts_avatar
	WHERE account_id = $1
	`, accountId).Scan(&id)

	if err != nil {
		log.Println("this account doesnt have picture before!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This account doesnt have a picture!"))
		return
	}

	path, err := uploadFile(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("This file cant be loaded"))
		return
	}

	now := timestamp()

	_, err = db.Exec(
		`
		UPDATE
		accounts_avatar
		SET
		modified_at = $1, path = $2
		WHERE
		account_id = $3
	`,
		now,
		path,
		accountId)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This file path cant be saved in the database"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("updated successfully"))
	return
}

func DeleteAccountAvatar(w http.ResponseWriter, r *http.Request) {
	var id int32
	accountId := r.PostFormValue("id")

	err = db.QueryRow(`
	SELECT id
	FROM
	accounts_avatar
	WHERE account_id = $1
	`, accountId).Scan(&id)

	if err != nil {
		log.Println("alredy deleted!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("alredy deleted!"))
		return
	}

	err := db.QueryRow(`
	DELETE 
	FROM 
	accounts_avatar 
	WHERE account_id = $1 RETURNING id;`,
		accountId).Scan(&id)

	if err != nil {
		log.Println("Can't be deleted!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Can't be deleted!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("deleted successfully"))
	return
}
