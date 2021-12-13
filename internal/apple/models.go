package apple

type TracksSearchResult struct {
    Data []Track `json:"data"`
}

type TrackGroup struct {
    Name   string
    Isrc   string
    Tracks []Track
}

type Track struct {
    Id         string          `json:"id"`
    Attributes TrackAttributes `json:"attributes"`
}

type TrackAttributes struct {
    ArtistName string `json:"artistName"`
    Name       string `json:"name"`
    Isrc       string `json:"isrc"`
}
