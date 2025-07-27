package chess

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChessStats struct {
	ChessRapid struct {
		Last struct {
			Rating int `json:"rating"`
		} `json:"last"`
	} `json:"chess_rapid"`
}

func FetchChessRating(username string) (int, error) {
	url := fmt.Sprintf("https://api.chess.com/pub/player/%s/stats", username)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var stats ChessStats
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return 0, err
	}

	return stats.ChessRapid.Last.Rating, nil
}
