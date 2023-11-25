package main

import (
	"io/ioutil"
	"fmt"
	"net/http"
	"regexp"
	"log"
	"time"
	"strings"
	"encoding/json"
	"bytes"
	"os"
)

func main() {
	semester := getSemester()


	resp, err := http.Get("https://raw.githubusercontent.com/gnulug/meetings/master/" + semester + "/schedule.md")
	if err != nil {
		log.Fatalf("\n%v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("\n%v", err)
	}

	configLocation := "config.json"

	if len(os.Args) >= 3 {
		if os.Args[1] == "--config" {
			configLocation = os.Args[2]
		}
	}

	sb := string(body)

	config := readConfig(configLocation)

	event := CurrentMeeting(sb)

	notifyEvent(event, config)
}

func CurrentMeeting(schedule string) string {
	
	re := regexp.MustCompile(`\|.*\r?\n`)
	headingsRaw := re.FindStringSubmatch(schedule)

	// re = regexp.MustCompile(`\| 2023-11-22.*\r?\n`) // test regex
	date := string(time.Now().Format("2006-01-02"))
	re = regexp.MustCompile(`\| ` + date + `.*\r?\n`)
	matches := re.FindStringSubmatch(schedule)
	
	var output string

	if len(matches) != 0 {
		fields := strings.Split(matches[0], `|`)
		headings := strings.Split(strings.Replace(headingsRaw[0], "\n", "", -1), `|`)

		for i := 2; i < len(headings)-1; i++ {
			output += fmt.Sprintf("%s: %s", headings[i], fields[i]);
			if i < len(headings) - 1 { output += "\n" }
		}
		return output
	}

	return "No Meeting Today"
}

func notifyEvent(event string, config map[string]interface{}) {

	postBody, err := json.Marshal(map[string]string{
		"msgtype": "m.text",
		"body": event,
	})
	if err !=nil {
		log.Fatalf("\nCould not make Json object: %v", err)
	}

	client := &http.Client{}


	token := genToken(config);
	roomURL := "https://matrix.org/_matrix/client/r0/rooms/" + config["internalRoomID"].(string) + "/send/m.room.message/?access_token=" + token

	req, err := http.NewRequest(http.MethodPut, roomURL, bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatalf("\nAn Error Occured %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("\n%v", err)
	}

	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("\n%v", err)
	}
	
	sb := string(body)
	fmt.Printf("%s", sb)
}

func genToken(config map[string]interface{}) string{
	postBody, err := json.Marshal(map[string]string{
		"type":"m.login.password",
		"user":config["username"].(string),
		"password":config["password"].(string),
	})


	if err != nil {
		log.Fatalf("\nCould not make Json object: %v", err)
	}
	
	resp, err := http.Post("https://matrix.org/_matrix/client/r0/login", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		log.Fatalf("\n%v", err)
	}
	
	defer resp.Body.Close()
	
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	
	token, ok := res["access_token"].(string)

	if !ok  {
		log.Fatalf("\nCould not find required field in response %v", token)
	}

	return token
}

func getSemester() string {
	t := time.Now()
	year := t.Year()
	month := int(t.Month())
	var season string

	if month <= 5 {
		season = "s"
	} else if month >= 8 {	// use an ENUM
		season = "f"
	} else {
		log.Fatal("\nSchools out")
	}
	return fmt.Sprintf("%d%s", year, season)
}

func readConfig(configLocation string) map[string]interface{} {
	content, err := ioutil.ReadFile(configLocation)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}
