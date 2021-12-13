package apple

import (
    "encoding/json"
    "fmt"
    "github.com/sirupsen/logrus"
    "net/http"
)

func SearchAll(isrcs chan string, token string) chan *TrackGroup {
    searchedTracks := make(chan *TrackGroup)
    go func(searchedTracks chan *TrackGroup) {
        defer close(searchedTracks)

        for isrc := range isrcs {
            res, err := search(isrc, token)
            if err != nil {
                logrus.WithError(err).Errorf("Error while search isrc %s", isrc)
                continue
            }
            logrus.Infof("founded %d by isrc %s", len(res.Tracks), isrc)
            searchedTracks <- res
        }

    }(searchedTracks)

    return searchedTracks
}

func search(isrc, token string) (*TrackGroup, error) {
    request, err := http.NewRequest("get", fmt.Sprintf(
        "https://amp-api.music.apple.com/v1/catalog/ru/songs?filter[isrc]=%s", isrc,
    ), nil,
    )
    if err != nil {
        return nil, err
    }

    request.Header.Add("Content-Type", "application/json")
    request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()
    var parsedResult TracksSearchResult
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&parsedResult)
    if err != nil {
        return nil, err
    }

    foundedCount := len(parsedResult.Data)
    if foundedCount == 0 {
        return nil, fmt.Errorf("unable to find any tracks in apple music by isrc %s", isrc)
    }

    tg := TrackGroup{
        Isrc:   isrc,
        Tracks: make([]Track, foundedCount),
        Name:   fmt.Sprintf("%s - %s", parsedResult.Data[0].Attributes.ArtistName, parsedResult.Data[0].Attributes.Name),
    }
    for pos, track := range parsedResult.Data {
        tg.Tracks[pos] = track
    }

    return &tg, nil
}
