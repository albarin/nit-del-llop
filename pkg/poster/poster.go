package poster

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/albarin/poster/pkg/webhooks"

	"github.com/goodsign/monday"
)

type Poster struct {
	Title  string
	Guest  string
	Date   time.Time
	Time   string
	PicURL string
	Type   string
}

func Parse(w webhooks.Webhook) Poster {
	poster := Poster{}

	for _, answer := range w.FormResponse.Answers {
		switch answer.Field.Ref {
		case "title":
			poster.Title = answer.Text
		case "guest":
			poster.Guest = answer.Text
		case "date":
			date, _ := time.Parse("2006-01-02", answer.Date)
			poster.Date = date
		case "time":
			poster.Time = answer.Text
		case "type":
			poster.Type = answer.Choice.Label
		case "pic":
			poster.PicURL = answer.PicURL
		}
	}

	return poster
}

func (p Poster) When() string {
	dayName := strings.Title(monday.Format(p.Date, "Monday", monday.LocaleCaES))
	dayNumber := monday.Format(p.Date, "02", monday.LocaleCaES)
	monthName := formatMonth(monday.Format(p.Date, "January", monday.LocaleCaES))

	return fmt.Sprintf("%s %s %s a les %s", dayName, dayNumber, monthName, p.Time)
}

func formatMonth(month string) string {
	if isVowel(month[:1]) {
		return fmt.Sprintf("d'%s", month)
	}

	return fmt.Sprintf("de %s", month)
}

func isVowel(char string) bool {
	if char == "a" || char == "e" || char == "i" || char == "o" || char == "u" {
		return true
	}

	return false
}

func (p Poster) Where() string {
	types := map[string]string{
		"Cena":    "sopar tertúlia amb l'autor",
		"Cuentos": "copa de vi i montaditos",
	}

	return fmt.Sprintf("a l'Orfeó Catalònia, %s", types[p.Type])
}

func (p Poster) Picture() (string, error) {
	filepath := "pic.png"

	err := downloadFile(filepath, p.PicURL)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}
