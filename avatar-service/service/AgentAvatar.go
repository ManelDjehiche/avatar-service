package service

import (
	"encoding/json"
	"log"
	"net/http"
)

func AddAgentAvatar(w http.ResponseWriter, r *http.Request) {

	var id string
	agentId := r.PostFormValue("id")

	if agentId == "" {
		log.Println("agentId is empty")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This agentId is empty"))
		return
	}

	err := db.QueryRow(`
	SELECT id
	FROM
	agents
	WHERE id = $1
	`, agentId).Scan(&id)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This agentId is wrong"))
		return
	}

	err = db.QueryRow(`
	SELECT id
	FROM
	agents_avatar
	WHERE agent_id = $1
	`, agentId).Scan(&id)

	if err == nil {
		log.Println("already exist!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This agentId has already a picture!"))
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
		agents_avatar
		(created_at, modified_at, path, agent_id)
		VALUES
		($1, $2, $3, $4)
		RETURNING id
	`,
		now,
		now,
		path,
		agentId).Scan(&id)

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

func GetAgentAvatar(w http.ResponseWriter, r *http.Request) {

	var id string
	agentId := r.PostFormValue("id")
	var agentAvatar Avatar
	var idAgent string

	err = db.QueryRow(`
		SELECT 
		id, agent_id, created_at, modified_at, path
		FROM 
		agents_avatar 
		WHERE 
		agent_id = $1
	`, agentId).Scan(&id, &idAgent, &agentAvatar.CreatedAt, &agentAvatar.ModifiedAt, &agentAvatar.Path)

	if err != nil {
		log.Println("this agent doesnt have picture before!")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("this agent doesnt have picture before!"))
		return
	}

	jsonBytes, err := json.Marshal(agentAvatar)

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

func UpdateAgentAvatar(w http.ResponseWriter, r *http.Request) {

	var id string
	agentId := r.PostFormValue("id")

	err = db.QueryRow(`
	SELECT id
	FROM
	agents_avatar
	WHERE agent_id = $1
	`, agentId).Scan(&id)

	if err != nil {
		log.Println("this agent doesnt have picture before!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This agent doesnt have a picture!"))
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
		agents_avatar
		SET
		modified_at = $1, path = $2
		WHERE
		agent_id = $3
	`,
		now,
		path,
		agentId)

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

func DeleteAgentAvatar(w http.ResponseWriter, r *http.Request) {
	var id string
	agentId := r.PostFormValue("id")

	err = db.QueryRow(`
	SELECT id
	FROM
	agents_avatar
	WHERE agent_id = $1
	`, agentId).Scan(&id)

	if err != nil {
		log.Println("alredy deleted!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("alredy deleted!"))
		return
	}

	err := db.QueryRow(`
	DELETE 
	FROM 
	agents_avatar 
	WHERE agent_id = $1 RETURNING id;`,
		agentId).Scan(&id)

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
