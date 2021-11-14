package service

import (
	"encoding/json"
	"log"
	"net/http"
)

func AddLeadAvatar(w http.ResponseWriter, r *http.Request) {
	var id int32
	leadId := r.PostFormValue("id")

	if leadId == "" {
		log.Println("leadId is empty")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This leadId is empty"))
		return
	}

	err := db.QueryRow(`
	SELECT id
	FROM
	leads
	WHERE id = $1
	`, leadId).Scan(&id)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This leadId is wrong"))
		return
	}

	err = db.QueryRow(`
	SELECT id
	FROM
	leads_avatar
	WHERE lead_id = $1
	`, leadId).Scan(&id)

	if err == nil {
		log.Println("already exist!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This leadId has already a picture!"))
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
		leads_avatar
		(created_at, modified_at, path, lead_id)
		VALUES
		($1, $2, $3, $4)
		RETURNING id
	`,
		now,
		now,
		path,
		leadId).Scan(&id)

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

func GetLeadAvatar(w http.ResponseWriter, r *http.Request) {

	var id int32
	leadId := r.PostFormValue("id")
	var leadAvatar Avatar
	var idLead string

	err = db.QueryRow(`
		SELECT 
		id, lead_id, created_at, modified_at, path
		FROM 
		leads_avatar 
		WHERE 
		lead_id = $1
	`, leadId).Scan(&id, &idLead, &leadAvatar.CreatedAt, &leadAvatar.ModifiedAt, &leadAvatar.Path)

	if err != nil {
		log.Println("this lead doesnt have picture before!")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("this lead doesnt have picture before!"))
		return
	}

	jsonBytes, err := json.Marshal(leadAvatar)

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

func UpdateLeadAvatar(w http.ResponseWriter, r *http.Request) {

	var id int32
	leadId := r.PostFormValue("id")

	err = db.QueryRow(`
	SELECT id
	FROM
	leads_avatar
	WHERE lead_id = $1
	`, leadId).Scan(&id)

	if err != nil {
		log.Println("this lead doesnt have picture before!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This lead doesnt have a picture!"))
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
		leads_avatar
		SET
		modified_at = $1, path = $2
		WHERE
		lead_id = $3
	`,
		now,
		path,
		leadId)

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

func DeleteLeadAvatar(w http.ResponseWriter, r *http.Request) {
	var id int32
	leadId := r.PostFormValue("id")

	err = db.QueryRow(`
	SELECT id
	FROM
	leads_avatar
	WHERE lead_id = $1
	`, leadId).Scan(&id)

	if err != nil {
		log.Println("alredy deleted!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("alredy deleted!"))
		return
	}

	err := db.QueryRow(`
	DELETE 
	FROM 
	leads_avatar 
	WHERE lead_id = $1 RETURNING id;`,
		leadId).Scan(&id)

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
