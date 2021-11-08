download-delegator ![CircleCI](https://img.shields.io/circleci/build/github/tislib/download-delegator)
====

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
