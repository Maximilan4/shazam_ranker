package shazam

type TracksSearchResult struct {
    Hits []SearchHit `json:"hits"`
}

type ShortTrack struct {
    Id string `json:"key"`
}

type SearchHit struct {
    Track ShortTrack `json:"track"`
}

type SearchResult struct {
    Tracks TracksSearchResult `json:"tracks"`
}

func (sr *SearchResult) Len() int {
    return len(sr.Tracks.Hits)
}

func (sr *SearchResult) GetIds() MatchingIds {
    ids := make(MatchingIds, sr.Len())
    for pos, hit := range sr.Tracks.Hits {
        ids[pos] = hit.Track.Id
    }

    return ids
}

type MatchingIds []string

type Track struct {
    Id    string `json:"key"`
    Title string `json:"title"`
    Isrc  string `json:"isrc"`
    Amid  string `json:"trackadamid"`
}

type Rating struct {
    TotalShazams int64 `json:"total"`
}
