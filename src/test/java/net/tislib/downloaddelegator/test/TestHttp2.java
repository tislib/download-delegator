package net.tislib.downloaddelegator.test;

import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.databind.ObjectMapper;
import kong.unirest.HttpResponse;
import kong.unirest.Unirest;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageUrl;

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

        for (int i = 0; i < 10; i++) {
            urls.add(PageUrl.builder()
                    .id(UUID.randomUUID())
                    .url(new URL("https://medium.com/@proustibat/how-to-fix-error-node-sass-does-not-yet-support-your-current-environment-os-x-64-bit-with-c1b3298e4af0"))
                    .method("GET")
                    .build());
        }

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
