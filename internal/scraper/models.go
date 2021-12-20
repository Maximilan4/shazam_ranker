package scraper

import "strconv"

type Track struct {
    Amid     string `json:"amid"`
    ShazamId string `json:"shazamId"`
    Isrc     string `json:"isrc"`
    Title    string `json:"title"`
    Rank     int64  `json:"rank"`
}

func (t *Track) CsvRow() []string {
    return []string{t.ShazamId, t.Amid, t.Isrc, t.Title, strconv.FormatInt(t.Rank, 10)}
}
