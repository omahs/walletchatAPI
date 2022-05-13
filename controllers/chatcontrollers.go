package controllers

import (
	"encoding/json"
	"io/ioutil"

	"fmt"

	"net/http"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"time"

	"github.com/gorilla/mux"
)

//GetAllInbox get all inboxes data
// func GetAllInbox(w http.ResponseWriter, r *http.Request) {
// 	var inbox []entity.Inbox
// 	database.Connector.Find(&inbox)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(inbox)
// }

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//GetInboxByID returns the latest message for each unique conversation
//TODO: properly design the relational DB structs to optimize this search/retrieve
func GetInboxByOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["address"]

	//fmt.Printf("GetInboxByOwner: %#v\n", key)

	//get all items that relate to passed in owner/address
	var chat []entity.Chatitem
	database.Connector.Where("fromaddr = ?", key).Or("toaddr = ?", key).Find(&chat)

	//get unique conversation addresses
	var uniqueChatMembers []string
	for _, chatitem := range chat {
		//fmt.Printf("search for unique addrs")
		if chatitem.Fromaddr != key {
			if !stringInSlice(chatitem.Fromaddr, uniqueChatMembers) {
				//fmt.Printf("Unique Addr Found: %#v\n", chatitem.Fromaddr)
				uniqueChatMembers = append(uniqueChatMembers, chatitem.Fromaddr)
			}
		}
		if chatitem.Toaddr != key {
			if !stringInSlice(chatitem.Toaddr, uniqueChatMembers) {
				//fmt.Printf("Unique Addr Found: %#v\n", chatitem.Toaddr)
				uniqueChatMembers = append(uniqueChatMembers, chatitem.Toaddr)
			}
		}
	}

	//fmt.Printf("find first message now")
	//for each unique chat member that is not the owner addr, get the latest message
	var firstItem entity.Chatitem
	var secondItem entity.Chatitem
	var userInbox []entity.Chatitem
	for _, chatmember := range uniqueChatMembers {
		database.Connector.Where("fromaddr = ?", chatmember).Where("toaddr = ?", key).First(&firstItem)
		database.Connector.Where("fromaddr = ?", key).Where("toaddr = ?", chatmember).First(&secondItem)

		//pick the most recent message
		if firstItem.Fromaddr != "" {
			if secondItem.Fromaddr == "" {
				userInbox = append(userInbox, firstItem)
			} else {
				layout := "2006-01-02T15:04:05.000Z"
				firstTime, error := time.Parse(layout, firstItem.Timestamp)
				if error != nil {
					fmt.Println(error)
					return
				}
				secondTime, error := time.Parse(layout, secondItem.Timestamp)
				if error != nil {
					fmt.Println(error)
					return
				}

				if firstTime.Before(secondTime) {
					userInbox = append(userInbox, firstItem)
				} else {
					userInbox = append(userInbox, secondItem)
				}
			}
		} else if secondItem.Fromaddr != "" {
			userInbox = append(userInbox, secondItem)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInbox)
}

//CreateInbox creates inbox
// func CreateInbox(w http.ResponseWriter, r *http.Request) {
// 	requestBody, _ := ioutil.ReadAll(r.Body)
// 	var inbox entity.Chatitem
// 	json.Unmarshal(requestBody, &inbox)

// 	database.Connector.Create(inbox)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(inbox)
// }

//UpdateInboxByOwner updates inbox with respective owner address
// func UpdateInboxByOwner(w http.ResponseWriter, r *http.Request) {
// 	requestBody, _ := ioutil.ReadAll(r.Body)
// 	var inbox entity.Chatitem
// 	json.Unmarshal(requestBody, &inbox)
// 	database.Connector.Save(&inbox)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(inbox)
// }

//DeletePersonByID delete's person with specific ID
// func DeleteInboxByOwner(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	key := vars["address"]

// 	var inbox entity.Inbox
// 	//id, _ := strconv.ParseString(key, 10, 64)
// 	database.Connector.Where("address = ?", key).Delete(&inbox)
// 	w.WriteHeader(http.StatusNoContent)
// }

//*********chat info*********************
//GetAllChatitems get all chat data
func GetAllChatitems(w http.ResponseWriter, r *http.Request) {
	//log.Println("get all chats")
	var chat []entity.Chatitem
	database.Connector.Find(&chat)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chat)
}

//GetChatFromAddressToOwner returns all chat items from user to owner
func GetChatFromAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["address"]

	var chat []entity.Chatitem
	database.Connector.Where("fromaddr = ?", key).Or("toaddr = ?", key).Find(&chat)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

//CreateChatitem creates Chatitem
func CreateChatitem(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var chat entity.Chatitem
	json.Unmarshal(requestBody, &chat)

	database.Connector.Create(chat)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}

//UpdateInboxByOwner updates person with respective owner address
func UpdateChatitemByOwner(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var chat entity.Chatitem
	//database.Connector.Where("fromaddr = ?", owner).Where("toaddr = ?", to).Find(&chat)
	json.Unmarshal(requestBody, &chat)
	database.Connector.Save(&chat)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chat)
}

func DeleteAllChatitemsToAddressByOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	to := vars["toAddr"]
	owner := vars["fromAddr"]

	var chat entity.Chatitem
	//id, _ := strconv.ParseString(key, 10, 64)
	database.Connector.Where("toAddr = ?", to).Where("fromAddr = ?", owner).Delete(&chat)
	w.WriteHeader(http.StatusNoContent)
}
