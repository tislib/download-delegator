**Download Delegator is a multipurpose batch downloader**
** **

**Purpose**

This application helps to download multiple files and stream the result in tar.gz format, this application may used for delegating download, achieving multi proxy, delegating network traffic, etc.

**Setup**
    
    docker run -d -p 8123:8123 docker.io/tisserv/download-delegator:latest

**Usage**

    curl -X POST http://localhost:8124 --data '{"urls": [{"id": "bb91c1ee-d03c-4281-b3af-adee589b939d", "url": "https://www.imdb.com/name/nm1869101/?ref_=fn_al_nm_1"}]}' --header 'Content-type: application/json' > result.tar.gz
    tar -zxf result.tar.gz
    cat bb91c1ee-d03c-4281-b3af-adee589b939d
    cat bb91c1ee-d03c-4281-b3af-adee589b939d.info


**Reference**

Post body: 

    urls    : PageUrl[] // urls which application will download
    delay   : number // delays between starting each download (downloads occours in parallel, this parameter is mainly to control download speed limit)
    timeout : number // timeout for response
    
PageUrl:

        id      : UUID // id is used to distinguish pages and find pages inside tar archive
        url     : URL
    
        proxy: {
                host: string
                port: number
                username: number // optional
                password: number // optional
        } // optional
