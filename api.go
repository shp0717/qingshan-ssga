package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"strconv"
)

var feedbackFilePath = func() string {
	exec, err := os.Executable()
	if err != nil {
		fmt.Printf("\033[31m[ERROR] Failed to get executable path: %v\033[0m\n", err)
		return "feedback.json" // fallback to current directory
	}
	path := filepath.Join(filepath.Dir(exec), "data/feedback.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create an empty feedback file if it doesn't exist
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			fmt.Printf("\033[31m[ERROR] Failed to create directories for feedback file: %v\033[0m\n", err)
			return "feedback.json" // fallback to current directory
		}
		err = os.WriteFile(path, []byte("[]"), 0644)
		if err != nil {
			fmt.Printf("\033[31m[ERROR] Failed to create feedback file: %v\033[0m\n", err)
			return "feedback.json" // fallback to current directory
		}
	}
	return path
}()

type Reply struct {
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

type feedbackRequest struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type replyRequest struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func decodeFeedbackRequest(r *http.Request) (feedbackRequest, error) {
	var request feedbackRequest

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return feedbackRequest{}, err
		}
		return request, nil
	}

	if err := r.ParseForm(); err != nil {
		return feedbackRequest{}, err
	}

	request.Name = r.FormValue("name")
	request.Title = r.FormValue("title")
	request.Message = r.FormValue("message")
	return request, nil
}

func decodeReplyRequest(r *http.Request) (replyRequest, error) {
	var request replyRequest

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return replyRequest{}, err
		}
		return request, nil
	}

	if err := r.ParseForm(); err != nil {
		return replyRequest{}, err
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return replyRequest{}, err
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

	newsData, err := Static.ReadFile("static/data/news.json")
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

	eventData, err := Static.ReadFile("static/data/events.json")
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
