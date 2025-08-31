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
		Best struct {
			Rating int `json:"rating"`
		}
	} `json:"chess_rapid"`
	ChessBullet struct {
		Last struct {
			Rating int `json:"rating"`
		} `json:"last"`
	} `json:"chess_bullet"`
	ChessBlitz struct {
		Last struct {
			Rating int `json:"rating"`
		} `json:"last"`
	} `json:"chess_blitz"`
}

type ChessUser struct {
	Avatar string `json:"avatar"`
}

func FetchChessUser(username string) (*ChessUser, error) {
	url := fmt.Sprintf("https://api.chess.com/pub/player/%s", username)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Chess.com a renvoyé le code %d", resp.StatusCode)
	}

	var user ChessUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing de la réponse JSON : %v", err)
	}

	return &user, nil
}

func FetchChessRating(username string) (*ChessStats, error) {
	url := fmt.Sprintf("https://api.chess.com/pub/player/%s/stats", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var stats ChessStats
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
