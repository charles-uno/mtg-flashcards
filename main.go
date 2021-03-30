package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "time"

    "github.com/charles-uno/mtgserver/lib"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)


type openingHand struct {
    Hand        []string    `json:"hand"`
    Library     []string    `json:"library"`
    OnThePlay   bool        `json:"on_the_play"`
}


func handleOpeningHand(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    log.Println("endpoint hit: /api/hand")
    deck := lib.LoadDeck()
    oh := openingHand{
        Hand: deck[:7],
        Library: deck[7:],
        OnThePlay: flip(),
    }
    json.NewEncoder(w).Encode(oh)
}


func handleSequencing(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    oh := openingHand{}
    err := json.NewDecoder(r.Body).Decode(&oh)
    if err != nil {
        reply := map[string]string{"error": err.Error()}
        b, _ := json.Marshal(reply)
        http.Error(w, string(b), http.StatusBadRequest)
        log.Println("bad payload at /api/play")
        return
    }
    game := lib.NewGame(oh.Hand, lib.Shuffled(oh.Library), oh.OnThePlay)
    for game.IsNotDone() {
        game = game.NextTurn()
    }
    fmt.Fprintf(w, game.ToJSON())
    log.Println("done with calculation at /api/play")
    fmt.Println(game.Pretty())
}


func main() {
    log.Println("launching service")
    router := mux.NewRouter()
    router.HandleFunc("/api/hand", handleOpeningHand).Methods("GET")
    router.HandleFunc("/api/play", handleSequencing).Methods("POST")
    cors := handlers.AllowedOrigins([]string{"*"})
    log.Fatal(http.ListenAndServe(":5001", handlers.CORS(cors)(router)))
}


func flip() bool {
    // Random generator should be seeded from shuffling, but let's be sure. We
    // only call this once per game anyway.
    rand.Seed(time.Now().UTC().UnixNano())
    return rand.Intn(2) == 0
}
