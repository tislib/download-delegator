package net.tislib.downloaddelegator.test;

import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.databind.ObjectMapper;
import kong.unirest.Unirest;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.Handler;
import net.tislib.downloaddelegator.data.Header;
import net.tislib.downloaddelegator.data.PageUrl;

import java.io.File;
import java.io.IOException;
import java.net.URL;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

public class Test {

    private static ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());

    public static void main(String[] args) throws IOException {
        DownloadRequest downloadRequest = new DownloadRequest();
        List<PageUrl> urls = new ArrayList<>();
        List<Handler> handlers = new ArrayList<>();

        for (int i = 0; i < 10000; i++) {
            if (i % 2 == 0) {
                urls.add(PageUrl.builder()
                        .id(UUID.randomUUID())
                        .url(new URL("http://localhost"))
                        .method("GET")
//                        .delay(1)
                        .header(Header.builder()
                                .name("user-agent")
                                .value("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3945.88 Safari/537.36")
                                .build())
                        .build());
            } else {
//                urls.add(PageUrl.builder()
//                        .id(UUID.randomUUID())
//                        .url(new URL("https://www.allmovie.com/artist/nicolas-cage-p10155"))
//                        .method("GET")
//                        .delay(100)
//                        .header(Header.builder()
//                                .name("user-agent")
//                                .value("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3945.88 Safari/537.36")
//                                .build())
//                        .build());
            }
        }

        Handler handler1 = new Handler();
        Handler handler2 = new Handler();

        handler1.setName("imdb.com");
        handler1.setConnectTimeout(10000);
        handler1.setSocketTimeout(10000);
        handler1.setDelay(15);

        handler2.setName("allmovie.com");
        handler2.setConnectTimeout(10000);
        handler2.setSocketTimeout(10000);
        handler2.setDelay(10);

        handlers.add(handler1);
        handlers.add(handler2);

        downloadRequest.setUrls(urls);

        String body = objectMapper.writeValueAsString(downloadRequest);
        objectMapper.writeValue(new File("/home/taleh/temp/ddreq4.json"), downloadRequest);

        System.out.println(body);

        byte[] data = Unirest.post("http://127.0.0.1:8080")
                .body(body)
                .header("Content-type", "application/json")
                .asBytes()
                .getBody();

//        DownloadResponse downloadResponse = objectMapper.readValue(new ByteArrayInputStream(data), DownloadResponse.class);
//
//        System.out.println(downloadResponse.getData().size());

        System.out.println(new String(data));
    }

}
