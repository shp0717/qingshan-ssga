package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed pages/*
var Pages embed.FS

//go:embed static/*
var Static embed.FS

func StaticFiles() {
	staticFS, err := fs.Sub(Static, "static")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(staticFS))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
}

func HandleRequests(port string) {
	// handle pages
	http.HandleFunc("/", HomePageHandler)
	http.HandleFunc("/news", NewsListPageHandler)
	http.HandleFunc("/news/", NewsPageHandler)
	http.HandleFunc("/events", EventsListPageHandler)
	http.HandleFunc("/events/", EventPageHandler)
	http.HandleFunc("/feedback", FeedbackPageHandler)
	http.HandleFunc("/feedback/new", NewFeedbackPageHandler)
	http.HandleFunc("/feedback/view", ViewFeedbackPageHandler)
	http.HandleFunc("/about", AboutPageHandler)
	http.HandleFunc("/contact", ContactPageHandler)
	http.HandleFunc("/secret", SecretPageHandler)

	// handle API
	http.HandleFunc("/api/", ApiNotFoundHandler) // catch-all for undefined API routes
	http.HandleFunc("/api/news", NewsApiHandler)
	http.HandleFunc("/api/events", EventApiHandler)
	http.HandleFunc("/api/feedbacks", GetFeedbacksHandler)
	http.HandleFunc("/api/post-feedback", PostFeedbackHandler)
	http.HandleFunc("/api/post-reply", AddReplyHandler)
	http.HandleFunc("/api/vErIfYsEcReTpAsSwOrD", VerifySecretPasswordHandler)
	http.HandleFunc("/api/secret-cmd", SecretCommandHandler)
	// handle 404 for any other routes
	fmt.Printf("[INFO] Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func renderNotFoundPage(w http.ResponseWriter, requestPath string) {
	content, err := Pages.ReadFile("pages/404.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read 404 page for %s: %v\033[0m\n", requestPath, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	w.Write(content)
	fmt.Printf("\033[33m[WARN] 404 Not Found: %s\033[0m\n", requestPath)
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderNotFoundPage(w, r.URL.Path)
		return
	}
	content, err := Pages.ReadFile("pages/home.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read home page: %v\033[0m\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Println("[INFO] HomePage accessed")
}

func NewsListPageHandler(w http.ResponseWriter, r *http.Request) {
	content, err := Pages.ReadFile("pages/news.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read news list page: %v\033[0m\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Println("[INFO] News list page accessed")
}

func EventsListPageHandler(w http.ResponseWriter, r *http.Request) {
	content, err := Pages.ReadFile("pages/events.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read events list page: %v\033[0m\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Println("[INFO] Events list page accessed")
}

func NewsPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[INFO] News page requested: %s\n", r.URL.Path)
	if r.URL.Path == "/news/" {
		http.Redirect(w, r, "/news", http.StatusMovedPermanently)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/news/")
	// requestedFile := "pages/news/news-" + id + ".html"
	// content, err := Pages.ReadFile(requestedFile)
	content, err := os.ReadFile(filepath.Join(execDir, "data", "pages", fmt.Sprintf("news-%s.html", id)))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			renderNotFoundPage(w, r.URL.Path)
			return
		}

		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read news page %s: %v\033[0m\n", id, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Printf("[INFO] News page accessed: %s\n", id)
}

func EventPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[INFO] Event page requested: %s\n", r.URL.Path)
	if r.URL.Path == "/events/" {
		http.Redirect(w, r, "/events", http.StatusMovedPermanently)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/events/")
	// requestedFile := "pages/events/event-" + id + ".html"
	// content, err := Pages.ReadFile(requestedFile)
	content, err := os.ReadFile(filepath.Join(execDir, "data", "pages", fmt.Sprintf("event-%s.html", id)))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			renderNotFoundPage(w, r.URL.Path)
			return
		}

		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read event page %s: %v\033[0m\n", id, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Printf("[INFO] Event page accessed: %s\n", id)
}

func FeedbackPageHandler(w http.ResponseWriter, r *http.Request) {
	content, err := Pages.ReadFile("pages/feedback.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read feedback page: %v\033[0m\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Println("[INFO] Feedback page accessed")
}

func NewFeedbackPageHandler(w http.ResponseWriter, r *http.Request) {
	content, err := Pages.ReadFile("pages/new-feedback.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read new feedback page: %v\033[0m\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Println("[INFO] New Feedback page accessed")
}

func ViewFeedbackPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[INFO] View Feedback page requested: %s\n", r.URL.Path)
	requestedFile := "pages/view-feedback.html"
	id := strings.TrimPrefix(r.URL.Path, "/feedback/view?id=")
	content, err := Pages.ReadFile(requestedFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			renderNotFoundPage(w, r.URL.Path)
			return
		}

		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read view feedback page %s: %v\033[0m\n", requestedFile, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Printf("[INFO] View Feedback page accessed: %s\n", id)
}

func AboutPageHandler(w http.ResponseWriter, r *http.Request) {
	content, err := Pages.ReadFile("pages/about.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read about page: %v\033[0m\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Println("[INFO] About page accessed")
}

func ContactPageHandler(w http.ResponseWriter, r *http.Request) {
	content, err := Pages.ReadFile("pages/contact.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read contact page: %v\033[0m\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Println("[INFO] Contact page accessed")
}

func SecretPageHandler(w http.ResponseWriter, r *http.Request) {
	content, err := Pages.ReadFile("pages/secret.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		fmt.Printf("\033[31m[ERROR] Failed to read secret page: %v\033[0m\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
	fmt.Println("[INFO] Secret page accessed")
}

func main() {
	StaticFiles()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	HandleRequests(port)
}
