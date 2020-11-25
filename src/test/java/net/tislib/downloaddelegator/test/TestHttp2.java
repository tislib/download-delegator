package net.tislib.downloaddelegator.test;

import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.databind.ObjectMapper;
import kong.unirest.HttpResponse;
import kong.unirest.Unirest;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.Header;
import net.tislib.downloaddelegator.data.PageUrl;
import net.tislib.downloaddelegator.data.Proxy;

import java.io.IOException;
import java.net.URL;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

public class TestHttp2 {

    private static ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());

    public static void main(String[] args) throws IOException {
        DownloadRequest downloadRequest = new DownloadRequest();
        List<PageUrl> urls = new ArrayList<>();

        for (int i = 0; i < 1; i++) {
            urls.add(PageUrl.builder()
                    .id(UUID.randomUUID())
                    .url(new URL("https://www.gsmarena.com/samsung-phones-9.php"))
                    .proxy(Proxy.builder()
                            .host("23.229.40.183")
                            .port(43758)
                            .username("talehsmail")
                            .password("W4MTMS712TOQU0BS55RWL806")
                            .build())
                    .header(Header.builder()
                            .name("Cookie")
                            .value("")
                            .build())
                    .header(Header.builder()
                            .name("User-Agent")
                            .value("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.16; rv:82.0) Gecko/20100101 Firefox/82.0")
                            .build())
                    .header(Header.builder()
                            .name("Accept")
                            .value("text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
                            .build())
                    .method("GET")
                    .build());
        }

//        GET /samsung-phones-9.php HTTP/1.1
//        Host: www.gsmarena.com
//        User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.16; rv:82.0) Gecko/20100101 Firefox/82.0
//        Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
//Accept-Language: en-US,en;q=0.5
//Accept-Encoding: gzip, deflate, br
//Proxy-Authorization: Basic dGFsZWhzbWFpbDpXNE1UTVM3MTJUT1FVMEJTNTVSV0w4MDY=
//Connection: keep-alive
//Cookie: __cfduid=dac9679f4add409057128e94a3f49d7371606265979; __unid=34884959-501b-035d-8040-4d2af09f7f5f; _ga=GA1.2.1744527464.1606265987; _gid=GA1.2.1734441614.1606265987; __gads=ID=712074b12f9c3654-22ea99fd10c5006f:T=1606266000:S=ALNI_MZm-5VPY2Pzxm4F2mXe3IL10XonVg
//Upgrade-Insecure-Requests: 1

        downloadRequest.setUrls(urls);
        downloadRequest.setDelay(1);

        String body = objectMapper.writeValueAsString(downloadRequest);
//        objectMapper.writeValue(new File("/home/taleh/temp/ddreq4.json"), downloadRequest);

        System.out.println(body);

        HttpResponse<byte[]> resp = Unirest.post("http://127.0.0.1:8123/download")
                .body(body)
                .header("Content-type", "application/json")
                .header("Accept-Encoding", "gzip")
                .asBytes();

        System.out.println(resp.getHeaders().get("Content-type"));

        System.out.println(new String(resp.getBody()));
    }

}
