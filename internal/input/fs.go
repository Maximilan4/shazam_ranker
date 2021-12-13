package input

import (
    "bufio"
    "github.com/sirupsen/logrus"
    "os"
)

func ScanIsrc(path string) (chan string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    scanner := bufio.NewScanner(file)
    isrcChan := make(chan string)
    go func() {
        err = scanIsrc(scanner, isrcChan)
        if err != nil {
            logrus.Error(err)
        }
        err = file.Close()
        if err != nil {
            logrus.Error(err)
        }
    }()

    return isrcChan, nil
}

func scanIsrc(scanner *bufio.Scanner, input chan string) error {
    defer close(input)
    for scanner.Scan() {
        isrc := scanner.Text()
        if isrc == "" {
            continue
        }

        input <- isrc
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    return nil
}
