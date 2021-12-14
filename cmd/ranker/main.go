package main

import (
    "flag"
    "github.com/sirupsen/logrus"
    "os"
    "shazam_ranker/internal/apple"
    "shazam_ranker/internal/input"
    "shazam_ranker/internal/output"
    "shazam_ranker/internal/scraper"
    "shazam_ranker/internal/shazam"
)

var (
    isrcFile, outputDir, token *string
    workersCount               *int
)

func init() {
    isrcFile = flag.String("i", "", "-i <path_to_input_file>")
    outputDir = flag.String("o", "output", "-p <path_to_output_dir>")
    token = flag.String("t", "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IldlYlBsYXlLaWQifQ.eyJpc3MiOiJBTVBXZWJQbGF5IiwiaWF0IjoxNjI1NzgxODY3LCJleHAiOjE2NDEzMzM4Njd9.yWOQkHcO59ydmtgIzP9TDB_Oasd_u-VNSzP-WJ1Fo_GUlICKq_LU9or5ABFx3EAF9geYHvBkIXvuCbVApN12sg", "-t <apple_auth_token>")
    workersCount = flag.Int("w", 2, "-w <workers_count>")
}

func main() {
    flag.Parse()
    if *isrcFile == "" || *outputDir == "" {
        flag.Usage()
        return
    }

    tokenInEnv := os.Getenv("APPLE_TOKEN")
    if *token == "" && tokenInEnv == "" {
        flag.Usage()
        return
    } else if tokenInEnv != "" {
        token = &tokenInEnv
    }

    start(*isrcFile, *outputDir, *token, *workersCount)
}

func start(inputFile, outputDir, token string, workersCount int) {
    isrcChan, err := input.ScanIsrc(inputFile)
    if err != nil {
        logrus.Fatal(err)
    }

    rankedTracks := make(chan *scraper.Track)

    for i := 1; i <= workersCount; i++ {
        go worker(isrcChan, rankedTracks, token)
    }

    if err = output.Write2Csv(rankedTracks, outputDir); err != nil {
        logrus.Fatal(err)
    }
}

func worker(input chan string, output chan *scraper.Track, token string) {
    defer close(output)
    appleTrackChan := apple.SearchAll(input, token)
    mergedTracks := shazam.MergeWithAppleTracks(appleTrackChan)
    rankedTracks := shazam.GetRanks(mergedTracks)

    for track := range rankedTracks {
        output <- track
    }
}
