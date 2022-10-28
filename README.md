download-delegator ![CircleCI](https://img.shields.io/circleci/build/github/tislib/download-delegator)
====

# About
Download Delegator is high concurrent webpage downloader, where you can send list of pages(thousands of pages) and it will download them, compress them and will return to you.

With Download Delegator, you can specify how it will behave with network, you can configure compression, max concurrency, etc.

Download Delegator also have built-in html parsing and script interpreting abilities, you can send your JS code and Download Delegator will execute this JS codes after download is done, it is good for minimising outbound traffic.

Download Delegator can also be deployed as aws lambda(s) and thousands of download delegators can be managed in parallel

# Getting started

*Request(curl):*
```
curl -X POST --location "http://localhost:8234/bulk" \
    -d "{
    \"url\": [\"https://imdb.com\", \"https://www.rottentomatoes.com/browse/opening\"],
    \"compression\": {
        \"algo\": \"bzip2\"
    },
    \"maxConcurrency\": 100
}" | bunzip2
```

*Request(goland http):*
```
POST http://localhost:8234/bulk

{
    "url": ["https://imdb.com", "https://www.rottentomatoes.com/browse/opening"],
    "compression": {
        "algo": "bzip2"
    },
    "maxConcurrency": 100
}
```

*Response:*
```
[
    {
        "Url":"https://imdb.com",
        "Content":"<downloaded content>",
        "StatusCode":404,
        "Duration":260201227,
        "DurationMS":260,
        "Error":"",
        "Index":0,
        "Retried":0
    },
    {
        "Url":"https://www.rottentomatoes.com/browse/opening",
        "Content":"<downloaded content>",
        "StatusCode":404,
        "Duration":260201227,
        "DurationMS":260,
        "Error":"",
        "Index":0,
        "Retried":0
    }
]
```

# Running

## Inside Docker

```
docker run -p 8234:8234 tislib/download-delegator:v2
```
