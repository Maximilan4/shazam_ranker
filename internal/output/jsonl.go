package output

import (
    "bufio"
    "encoding/json"
    "shazam_ranker/internal/scraper"
)

func Write2Jsonl(inputChan chan *scraper.Track, outputDir string) error {
    dir, err := createOutputIfNotExists(outputDir)
    if err != nil {
        return err
    }

    file, err := createFile(dir, ".jsonl")
    if err != nil {
        return err
    }
    defer file.Close()
    jsonlWriter := bufio.NewWriter(file)
    defer jsonlWriter.Flush()

    var row []byte
    for track := range inputChan {
        row, err = json.Marshal(track)
        if err != nil {
            return err
        }

        row = append(row, []byte("\n")...)
        _, err = jsonlWriter.Write(row)
        if err != nil {
            return err
        }
        jsonlWriter.Flush()
    }

    return nil
}
