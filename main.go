package main

import (
        "encoding/json"
        "fmt"
        "html/template"
        "net/http"
        "os"
        "strconv"
        "time"
)

type Paste struct {
        ID      int
        Content string
        Created time.Time
        Expires time.Time
}

var pastes []Paste
var port = "9000"
var indexFile = "index.html"
var pasteFile = "paste.html"
var pastesFile = "../pastes.json"
var templates = template.New("")

func main() {
        loadPastes()
        templates.Funcs(template.FuncMap{
                "truncate":   truncate,
                "formatTime": formatTime,
        })
        _, err := templates.ParseFiles(indexFile, pasteFile)
        if err != nil {
                panic(err)
        }

        http.HandleFunc("/", homeHandler)
        http.HandleFunc("/new", newPasteHandler)
        http.HandleFunc("/paste/", viewPasteHandler)
        err = http.ListenAndServe(":"+port, nil)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

}

func truncate(s string, length int) string {
        if len(s) > length {
                return s[:length] + "..."
        }
        return s
}

func formatTime(t time.Time) string {
        return t.Format("2006-01-02 15:04:05")
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
        err := templates.ExecuteTemplate(w, "index.html", pastes)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

func newPasteHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        r.ParseForm()
        content := r.FormValue("content")

        if content == "" {
                http.Error(w, "Content cannot be empty", http.StatusBadRequest)
                return
        }

        created := time.Now()
        expires := time.Now().AddDate(0, 0, 7)

        // Check for expired pastes or empty spots
        var id int
        for i := range pastes {
                if pastes[i].Content == "" || pastes[i].Expires.Before(time.Now()) {
                        id = i + 1
                        pastes[i] = Paste{
                                ID:      id,
                                Content: content,
                                Created: created,
                                Expires: expires,
                        }
                        break
                }
        }

        // If no expired pastes or empty spots, append a new paste
        if id == 0 {
                id = len(pastes) + 1
                paste := Paste{
                        ID:      id,
                        Content: content,
                        Created: created,
                        Expires: expires,
                }
                pastes = append(pastes, paste)
        }

        savePastes()

        http.Redirect(w, r, fmt.Sprintf("/paste/%d", id), http.StatusFound)
}

func viewPasteHandler(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.URL.Path[len("/paste/"):])
        if err != nil || id < 1 {
                http.Redirect(w, r, "/", http.StatusPermanentRedirect)
                return
        }

        if id > len(pastes) {
                http.NotFound(w, r)
                return
        }

        paste := pastes[id-1]

        if paste.Created.Before(time.Now().Add(-7 * 24 * time.Hour)) {
                http.NotFound(w, r)
                return
        }

        err = templates.ExecuteTemplate(w, "paste.html", paste)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

func savePastes() {
        f, err := os.Create(pastesFile)

        if err != nil {
                fmt.Println(err)
                return
        }

        defer f.Close()

        encoder := json.NewEncoder(f)
        encoder.Encode(pastes)
}

func loadPastes() {
        f, err := os.Open(pastesFile)
        if err != nil {
                return
        }
        defer f.Close()

        decoder := json.NewDecoder(f)
        var validPastes []Paste
        err = decoder.Decode(&validPastes)
        if err != nil {
                return
        }

        // Initialize pastes slice with empty pastes
        pastes = make([]Paste, len(validPastes))
        for i := range pastes {
                pastes[i].ID = i + 1
                pastes[i].Content = ""
                pastes[i].Expires = time.Time{}
        }

        // Fill pastes slice with valid pastes from JSON file
        for _, p := range validPastes {
                if p.ID <= len(pastes) {
                        pastes[p.ID-1] = p
                }
        }

        // Remove expired pastes
        now := time.Now()
        for i := range pastes {
                if pastes[i].Expires.Before(now) {
                        pastes[i].Content = ""
                        pastes[i].Expires = time.Time{}
                }
        }

        savePastes()

        // Check for expired pastes every minute
        ticker := time.NewTicker(time.Minute)
        go func() {
                for {
                        <-ticker.C
                        now := time.Now()
                        for i := range pastes {
                                if pastes[i].Expires.Before(now) {
                                        pastes[i].Content = ""
                                        pastes[i].Expires = time.Time{}
                                }
                        }
                        savePastes()
                }
        }()
}
