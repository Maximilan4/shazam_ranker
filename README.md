## shazam ranker
loads shazam rating by track`s isrc

## Algo
- search in apple by isrc
- search in shazam by name and track title
- match isrc and amid in search result
- if match founded -> loads shazam rank
- write all results to file

## RUN
```bash
go build -o ranker cmd/ranker/main.go

./ranker -i isrc.txt -o output -w 2
```