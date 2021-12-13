package output

import (
    "encoding/csv"
    "os"
    "path"
    "path/filepath"
    "shazam_ranker/internal/scraper"
    "strconv"
    "time"
)

var csvHeader = []string{"Shazam ID", "Apple Music ID", "ISRC", "Title", "Rank"}

func Write2Csv(inputChan chan *scraper.Track, outputDir string) error {
    dir, err := createOutputIfNotExists(outputDir)
    if err != nil {
        return err
    }

    file, err := createCsv(dir)
    if err != nil {
        return err
    }
    defer file.Close()
    csvWriter := csv.NewWriter(file)
    err = csvWriter.Write(csvHeader)
    if err != nil {
        return err
    }

    defer csvWriter.Flush()

    for track := range inputChan {
        err = csvWriter.Write(track.CsvRow())
        if err != nil {
            return err
        }
        csvWriter.Flush()
    }

    return nil
}

func createCsv(dir *string) (*os.File, error) {
    return os.Create(path.Join(*dir, strconv.FormatInt(time.Now().Unix(), 10)+".csv"))
}

func createOutputIfNotExists(outputDir string) (*string, error) {
    var dir string
    var err error
    if !path.IsAbs(outputDir) {
        dir, err = filepath.Abs(outputDir)
        if err != nil {
            return nil, err
        }
    } else {
        dir = outputDir
    }

    if _, err = os.Stat(dir); err != nil {
        if os.IsNotExist(err) {
            err = os.Mkdir(outputDir, os.ModePerm)
            if err != nil {
                return nil, err
            }
            return &dir, nil
        }

        return nil, err
    }

    return &dir, err
}
