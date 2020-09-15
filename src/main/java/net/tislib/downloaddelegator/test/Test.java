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

        urls.add(PageUrl.builder()
                .id(UUID.randomUUID())
                .url(new URL("http://tisserv.net"))
                .method("GET")
                .bind("172.20.11.45")
                .build());

        downloadRequest.setUrls(urls);
        downloadRequest.setDelay(1000);

        String body = objectMapper.writeValueAsString(downloadRequest);
        objectMapper.writeValue(new File("/home/taleh/temp/ddreq4.json"), downloadRequest);

        System.out.println(body);

        byte[] data = Unirest.post("http://127.0.0.1:8123")
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
