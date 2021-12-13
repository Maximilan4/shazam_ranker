package scraper

import "strconv"

type Track struct {
    Amid, ShazamId, Isrc, Title string
    Rank                        int64
}

func (t *Track) CsvRow() []string {
    return []string{t.ShazamId, t.Amid, t.Isrc, t.Title, strconv.FormatInt(t.Rank, 10)}
}
