package shazam

import (
    "encoding/json"
    "fmt"
    "github.com/sirupsen/logrus"
    "net/http"
    "net/url"
    "shazam_ranker/internal/apple"
    "shazam_ranker/internal/scraper"
    "sync"
)

func GetRanks(input chan *scraper.Track) chan *scraper.Track {
    output := make(chan *scraper.Track)
    go getTracksRating(input, output)

    return output
}

func getTracksRating(input, output chan *scraper.Track) {
    defer close(output)
    for track := range input {
        rating, err := getTrackRating(track.ShazamId)
        if err != nil {
            logrus.WithError(err).Errorf("unable to get rating for amid %s", track.Amid)
        } else {
            track.Rank = rating.TotalShazams
        }

        output <- track
    }
}

func getTrackRating(id string) (*Rating, error) {
    requestUrl := fmt.Sprintf("https://www.shazam.com/services/count/v2/web/track/%s", id)
    request, err := http.NewRequest("GET", requestUrl, nil)
    if err != nil {
        return nil, err
    }
    request.Header.Add("Content-Type", "application/json")

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()
    var rating Rating
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&rating)
    if err != nil {
        return nil, err
    }

    return &rating, nil
}

func MergeWithAppleTracks(input chan *apple.TrackGroup) chan *scraper.Track {
    output := make(chan *scraper.Track)
    go merge(input, output)

    return output
}

func merge(input chan *apple.TrackGroup, output chan *scraper.Track) {
    defer close(output)
    for trackGroup := range input {
        ids, err := getDetailIdsByName(trackGroup.Name)
        if err != nil {
            logrus.WithError(err).Error("unable to search tracks in shazam")
            continue
        }

        if len(ids) == 0 {
            logrus.Errorf("unable to find any tracks in shazam by isrc %s", trackGroup.Isrc)
            continue
        }

        tracks := getDetailedTracks(ids)
        groupHasMatch := false
        for _, appleTrackInfo := range trackGroup.Tracks {
            loaded, found := tracks.Load(appleTrackInfo.Id)
            if !found {
                continue
            }
            shazamTrack := loaded.(*Track)

            if shazamTrack.Isrc != appleTrackInfo.Attributes.Isrc {
                logrus.Infof("isrc mismatch, expected %s, got %s", appleTrackInfo.Attributes.Isrc, shazamTrack.Isrc)
                continue
            }
            groupHasMatch = true
            output <- &scraper.Track{
                Amid:     shazamTrack.Amid,
                ShazamId: shazamTrack.Id,
                Isrc:     shazamTrack.Isrc,
                Title:    shazamTrack.Title,
                Rank:     0,
            }
        }

        if !groupHasMatch {
            logrus.Infof("unable to find any track is shazam search by isrc %s", trackGroup.Tracks[0].Attributes.Isrc)
        }
    }
}

func getDetailedTracks(ids MatchingIds) *sync.Map {
    tracks := &sync.Map{}
    wg := &sync.WaitGroup{}
    for _, id := range ids {
        wg.Add(1)
        go func(id string, wg *sync.WaitGroup) {
            defer wg.Done()
            track, err := getTrackById(id)
            if err != nil {
                logrus.WithError(err).Errorf("unable to get track in shazam by id %s", track.Id)
                return
            }

            tracks.Store(track.Amid, track)
        }(id, wg)
    }

    wg.Wait()
    return tracks
}

func getTrackById(id string) (*Track, error) {
    requestUrl := fmt.Sprintf("https://www.shazam.com/discovery/v5/ru/RU/web/-/track/%s", id)
    request, err := http.NewRequest("GET", requestUrl, nil)
    if err != nil {
        return nil, err
    }
    request.Header.Add("Content-Type", "application/json")

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()
    var parsedResult Track
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&parsedResult)
    if err != nil {
        return nil, err
    }

    return &parsedResult, nil
}

func getDetailIdsByName(name string) (MatchingIds, error) {
    requestUrl, err := url.Parse("https://www.shazam.com/services/search/v4/ru/RU/web/search")
    if err != nil {
        return nil, err
    }
    query := requestUrl.Query()
    query.Add("term", name)
    query.Add("numResults", "5")
    query.Add("offset", "0")
    query.Add("limit", "5")
    query.Add("types", "songs")

    requestUrl.RawQuery = query.Encode()
    request, err := http.NewRequest("GET", requestUrl.String(), nil)
    if err != nil {
        return nil, err
    }
    request.Header.Add("Content-Type", "application/json")

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()
    var parsedResult SearchResult
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&parsedResult)
    if err != nil {
        return nil, err
    }

    return parsedResult.GetIds(), nil
}
