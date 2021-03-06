package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ad/gocc/ccredis"
	"github.com/ad/gocc/structs"
	uuid "github.com/nu7hatch/gouuid"
)

func ZondCreateHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "POST" {
	name := r.FormValue("name")

	u, _ := uuid.NewV4()
	var UUID = u.String()
	var msec = time.Now().Unix()

	if len(name) == 0 {
		name = UUID
	}

	ccredis.Client.SAdd("zonds", UUID)

	userUUID, _ := ccredis.Client.Get("user/UUID/" + r.Header.Get("X-Forwarded-User")).Result()
	if userUUID == "" {
		u, _ := uuid.NewV4()
		userUUID = u.String()
		ccredis.Client.Set(fmt.Sprintf("user/UUID/%s", r.Header.Get("X-Forwarded-User")), userUUID, 0)
	}

	zond := structs.Zond{UUID: UUID, Name: name, Created: msec, Creator: userUUID}
	js, _ := json.Marshal(zond)

	ccredis.Client.Set("zonds/"+UUID, string(js), 0)
	ccredis.Client.SAdd("user/zonds/"+userUUID, UUID)

	log.Println("Zond created", UUID)

	// if r.Header.Get("X-Requested-With") == "xmlhttprequest" {
	// w.Header().Set("X-CSRF-Token", csrf.Token(r))
	fmt.Fprintf(w, `{"status": "ok", "UUID": "%s"}`, UUID)
	// } else {
	// 	ShowCreateForm(w, r, UUID)
	// }
	// } else {
	// 	ShowCreateForm(w, r, "")
	// }
}
