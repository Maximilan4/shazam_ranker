package main

import "testing"

func BenchmarkSingleWorker(b *testing.B) {
    start(
        "../../test/isrc.txt",
        "output",
        "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IldlYlBsYXlLaWQifQ.eyJpc3MiOiJBTVBXZWJQbGF5IiwiaWF0IjoxNjI1NzgxODY3LCJleHAiOjE2NDEzMzM4Njd9.yWOQkHcO59ydmtgIzP9TDB_Oasd_u-VNSzP-WJ1Fo_GUlICKq_LU9or5ABFx3EAF9geYHvBkIXvuCbVApN12sg",
        1,
    )
    b.ReportAllocs()
}

func BenchmarkTwoWorkers(b *testing.B) {
    start(
        "../../test/isrc.txt",
        "output",
        "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IldlYlBsYXlLaWQifQ.eyJpc3MiOiJBTVBXZWJQbGF5IiwiaWF0IjoxNjI1NzgxODY3LCJleHAiOjE2NDEzMzM4Njd9.yWOQkHcO59ydmtgIzP9TDB_Oasd_u-VNSzP-WJ1Fo_GUlICKq_LU9or5ABFx3EAF9geYHvBkIXvuCbVApN12sg",
        2,
    )
    b.ReportAllocs()
}
