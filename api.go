package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"crypto/sha256"
)

var execDir = func() string {
	exec, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("Failed to get executable path: %v", err))
	}
	return filepath.Dir(exec)
}()

var feedbackFilePath = filepath.Join(execDir, "data", "feedbacks.json")

type Reply struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

type Feedback struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Title   string  `json:"title"`
	Message string  `json:"message"`
	Date    string  `json:"date"`
	Replies []Reply `json:"replies"`
}

type FeedbackRequest struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ReplyRequest struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func decodeFeedbackRequest(r *http.Request) (FeedbackRequest, error) {
	var request FeedbackRequest

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return FeedbackRequest{}, err
		}
		return request, nil
	}

	if err := r.ParseForm(); err != nil {
		return FeedbackRequest{}, err
	}

	request.Name = r.FormValue("name")
	request.Title = r.FormValue("title")
	request.Message = r.FormValue("message")
	return request, nil
}

func decodeReplyRequest(r *http.Request) (ReplyRequest, error) {
	var request ReplyRequest

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return ReplyRequest{}, err
		}
		return request, nil
	}

	if err := r.ParseForm(); err != nil {
		return ReplyRequest{}, err
	}

	// id, err := strconv.Atoi(r.FormValue("id"))
	// if err != nil {
	// 	return ReplyRequest{}, err
	// }
	var id int
	idStr := r.FormValue("id")
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return ReplyRequest{}, err
	}
	request.ID = id
	request.Name = r.FormValue("name")
	request.Message = r.FormValue("message")
	return request, nil
}

func ApiNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "API endpoint not found", http.StatusNotFound)
	fmt.Printf("\033[31m[ERROR] API endpoint not found: %s\033[0m\n", r.URL.Path)
}

func NewsApiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// newsData, err := Static.ReadFile("static/data/news.json")
	newsData, err := os.ReadFile(filepath.Join(execDir, "data", "news.json"))
	if err != nil {
		http.Error(w, "Could not load news data", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read news data: %v\033[0m\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(newsData)
	fmt.Println("[INFO] News API accessed")
}

func EventApiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// eventData, err := Static.ReadFile("static/data/events.json")
	eventData, err := os.ReadFile(filepath.Join(execDir, "data", "events.json"))
	if err != nil {
		http.Error(w, "Could not load event data", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read event data: %v\033[0m\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(eventData)
	fmt.Println("[INFO] Event API accessed")
}

func PostFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := decodeFeedbackRequest(r)
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		fmt.Printf("\033[31m[ERROR] Failed to parse feedback form data: %v\033[0m\n", err)
		return
	}

	name := strings.TrimSpace(request.Name)
	title := strings.TrimSpace(request.Title)
	message := strings.TrimSpace(request.Message)
	if name == "" || title == "" || message == "" {
		http.Error(w, "Name, title, and message are required", http.StatusBadRequest)
		fmt.Println("\033[31m[ERROR] Feedback request missing required fields\033[0m")
		return
	}

	feedbacksJSON, err := os.ReadFile(feedbackFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			feedbacksJSON = []byte("[]")
		} else {
			http.Error(w, "Failed to read feedback data", http.StatusInternalServerError)
			fmt.Printf("\033[31m[ERROR] Failed to read feedback data: %v\033[0m\n", err)
			return
		}
	}

	var feedbacks []Feedback
	err = json.Unmarshal(feedbacksJSON, &feedbacks)
	if err != nil {
		http.Error(w, "Failed to parse feedback data", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to parse feedback data: %v\033[0m\n", err)
		return
	}

	feedback := Feedback{
		ID:      len(feedbacks) + 1,
		Name:    name,
		Title:   title,
		Message: message,
		Date:    time.Now().UTC().Format(time.RFC3339),
		Replies: []Reply{},
	}

	feedbacks = append(feedbacks, feedback)

	newFeedbacksJSON, err := json.MarshalIndent(feedbacks, "", "  ")
	if err != nil {
		http.Error(w, "Failed to encode feedback data", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to encode feedback data: %v\033[0m\n", err)
		return
	}

	err = os.WriteFile(feedbackFilePath, newFeedbacksJSON, 0644)
	if err != nil {
		http.Error(w, "Failed to save feedback data", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to save feedback data: %v\033[0m\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("[INFO] Feedback received and saved")
}

func AddReplyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	request, err := decodeReplyRequest(r)
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		fmt.Printf("\033[31m[ERROR] Failed to parse reply form data: %v\033[0m\n", err)
		return
	}

	feedbackID := request.ID
	name := strings.TrimSpace(request.Name)
	message := strings.TrimSpace(request.Message)
	if feedbackID == 0 || name == "" || message == "" {
		http.Error(w, "ID, name, and message are required", http.StatusBadRequest)
		fmt.Println("\033[31m[ERROR] Reply request missing required fields\033[0m")
		return
	}

	feedbacksJSON, err := os.ReadFile(feedbackFilePath)
	if err != nil {
		http.Error(w, "Failed to read feedback data", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read feedback data: %v\033[0m\n", err)
		return
	}

	var feedbacks []Feedback
	err = json.Unmarshal(feedbacksJSON, &feedbacks)
	if err != nil {
		http.Error(w, "Failed to parse feedback data", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to parse feedback data: %v\033[0m\n", err)
		return
	}

	for i, fb := range feedbacks {
		if fb.ID == feedbackID {
			reply := Reply{
				ID:      len(fb.Replies) + 1,
				Name:    name,
				Message: message,
				Date:    time.Now().UTC().Format(time.RFC3339),
			}
			feedbacks[i].Replies = append(feedbacks[i].Replies, reply)

			newFeedbacksJSON, err := json.MarshalIndent(feedbacks, "", "  ")
			if err != nil {
				http.Error(w, "Failed to encode feedback data", http.StatusInternalServerError)
				fmt.Printf("\033[31m[ERROR] Failed to encode feedback data: %v\033[0m\n", err)
				return
			}

			err = os.WriteFile(feedbackFilePath, newFeedbacksJSON, 0644)
			if err != nil {
				http.Error(w, "Failed to save feedback data", http.StatusInternalServerError)
				fmt.Printf("\033[31m[ERROR] Failed to save feedback data: %v\033[0m\n", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Println("[INFO] Reply added to feedback")
			return
		}
	}

	http.Error(w, "Feedback not found", http.StatusNotFound)
	fmt.Printf("\033[31m[ERROR] Feedback not found for ID: %d\033[0m\n", feedbackID)
}

func GetFeedbacksHandler(w http.ResponseWriter, r *http.Request) {
	feedbacksJSON, err := os.ReadFile(feedbackFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			feedbacksJSON = []byte("[]")
		} else {
			http.Error(w, "Failed to read feedback data", http.StatusInternalServerError)
			fmt.Printf("\033[31m[ERROR] Failed to read feedback data: %v\033[0m\n", err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(feedbacksJSON)
	fmt.Println("[INFO] Feedbacks API accessed")
}

func VerifySecretPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Password string `json:"password"`
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		fmt.Printf("\033[31m[ERROR] Failed to parse JSON data: %v\033[0m\n", err)
		return
	}

	hashedInput := fmt.Sprintf("%x", sha256.Sum256([]byte(request.Password)))
	expectedHash := "3afa2aac62dd41a878a70969eecab3b24f0b01d38e4a09739088596b6757e990"

	if hashedInput == expectedHash {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
		fmt.Println("[INFO] Secret password verified successfully")
	} else {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		fmt.Println("\033[31m[ERROR] Incorrect secret password attempt\033[0m")
	}
}

func SecretCommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request map[string]string

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		fmt.Printf("\033[31m[ERROR] Failed to parse JSON data: %v\033[0m\n", err)
		return
	}

	command, exists := request["command"]
	if !exists {
		http.Error(w, "Command not provided", http.StatusBadRequest)
		fmt.Println("\033[31m[ERROR] Secret command not provided in request\033[0m")
		return
	}

	switch command {
	case "delete_feedback":
		idStr, exists := request["feedback_id"]
		if !exists {
			http.Error(w, "ID not provided for delete_feedback command", http.StatusBadRequest)
			fmt.Println("\033[31m[ERROR] ID not provided for delete_feedback command\033[0m")
			return
		}

		feedbacksJSON, err := os.ReadFile(feedbackFilePath)
		if err != nil {
			http.Error(w, "Failed to read feedback data", http.StatusInternalServerError)
			fmt.Printf("\033[31m[ERROR] Failed to read feedback data: %v\033[0m\n", err)
			return
		}

		var feedbacks []Feedback
		err = json.Unmarshal(feedbacksJSON, &feedbacks)
		if err != nil {
			http.Error(w, "Failed to parse feedback data", http.StatusInternalServerError)
			fmt.Printf("\033[31m[ERROR] Failed to parse feedback data: %v\033[0m\n", err)
			return
		}

		found := false
		for i, fb := range feedbacks {
			if fmt.Sprintf("%d", fb.ID) == idStr {
				feedbacks = append(feedbacks[:i], feedbacks[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			http.Error(w, "Feedback not found", http.StatusNotFound)
			fmt.Printf("\033[31m[ERROR] Feedback not found for ID: %s\033[0m\n", idStr)
			return
		}

		newFeedbacksJSON, err := json.Marshal(feedbacks)
		if err != nil {
			http.Error(w, "Failed to encode feedback data", http.StatusInternalServerError)
			fmt.Printf("\033[31m[ERROR] Failed to encode feedback data: %v\033[0m\n", err)
			return
		}

		err = os.WriteFile(feedbackFilePath, newFeedbacksJSON, 0644)
		if err != nil {
			http.Error(w, "Failed to save feedback data", http.StatusInternalServerError)
			fmt.Printf("\033[31m[ERROR] Failed to save feedback data: %v\033[0m\n", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Println("[INFO] Feedback deleted successfully")
	case "delete_reply":
		feedbackIDStr, exists := request["feedback_id"]
		if !exists {
			http.Error(w, "Feedback ID not provided for delete_reply command", http.StatusBadRequest)
			fmt.Println("\033[31m[ERROR] Feedback ID not provided for delete_reply command\033[0m")
			return
		}
		replyIDStr, exists := request["reply_id"]
		if !exists {
			http.Error(w, "Reply ID not provided for delete_reply command", http.StatusBadRequest)
			fmt.Println("\033[31m[ERROR] Reply ID not provided for delete_reply command\033[0m")
			return
		}

		feedbacksJSON, err := os.ReadFile(feedbackFilePath)
		if err != nil {
			http.Error(w, "Failed to read feedback data", http.StatusInternalServerError)
			fmt.Printf("\033[31m[ERROR] Failed to read feedback data: %v\033[0m\n", err)
			return
		}

		var feedbacks []Feedback
		err = json.Unmarshal(feedbacksJSON, &feedbacks)
		if err != nil {
			http.Error(w, "Failed to parse feedback data", http.StatusInternalServerError)
			fmt.Printf("\033[31m[ERROR] Failed to parse feedback data: %v\033[0m\n", err)
			return
		}

		found := false
		for i, fb := range feedbacks {
			if fmt.Sprintf("%d", fb.ID) == feedbackIDStr {
				for j, rp := range fb.Replies {
					if fmt.Sprintf("%d", rp.ID) == replyIDStr {
						feedbacks[i].Replies = append(feedbacks[i].Replies[:j], feedbacks[i].Replies[j+1:]...)
						found = true
						break
					}
				}
				break
			}
		}
		if !found {
			http.Error(w, "Feedback or reply not found", http.StatusNotFound)
			fmt.Printf("\033[31m[ERROR] Feedback or reply not found for Feedback ID: %s, Reply ID: %s\033[0m\n", feedbackIDStr, replyIDStr)
			return
		}

		newFeedbacksJSON, err := json.Marshal(feedbacks)
		if err != nil {
			http.Error(w, "Failed to encode feedback data", http.StatusInternalServerError)
			fmt.Printf("\033[31m[ERROR] Failed to encode feedback data: %v\033[0m\n", err)
			return
		}

		err = os.WriteFile(feedbackFilePath, newFeedbacksJSON, 0644)
		if err != nil {
			http.Error(w, "Failed to save feedback data", http.StatusInternalServerError)
			fmt.Printf("\033[31m[ERROR] Failed to save feedback data: %v\033[0m\n", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Println("[INFO] Reply deleted successfully")
	default:
		http.Error(w, "Unknown command", http.StatusBadRequest)
		fmt.Printf("\033[31m[ERROR] Unknown secret command: %s\033[0m\n", command)
	}
}
