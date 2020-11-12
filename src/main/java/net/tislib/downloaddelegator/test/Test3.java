package net.tislib.downloaddelegator.test;

import com.fasterxml.jackson.databind.ObjectMapper;
import kong.unirest.Unirest;
import lombok.SneakyThrows;
import net.tislib.downloaddelegator.base.TimeCalc;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageUrl;
import net.tislib.downloaddelegator.data.Proxy;
import org.apache.commons.compress.archivers.ArchiveEntry;
import org.apache.commons.compress.archivers.tar.TarArchiveInputStream;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.net.URL;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;
import java.util.zip.GZIPInputStream;

public class Test3 {

    private final static ObjectMapper objectMapper = new ObjectMapper();
    private final static TimeCalc timeCalc = new TimeCalc();

    @SneakyThrows
    public static void main(String[] args) {
        while (true) {
            try {
                main2(args);
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
    }

    @SneakyThrows
    public static void main2(String[] args) {
        String imdbUrl = "https://www.imdb.com/";

        DownloadRequest downloadRequest = new DownloadRequest();
        List<PageUrl> urls = new ArrayList<>();

        for (int i = 0; i < 50; i++) {
            urls.add(PageUrl.builder()
                    .id(UUID.randomUUID())
                    .url(new URL(imdbUrl))
                    .method("GET")
                    .proxy(Proxy.builder()
                            .host("83.149.70.159")
                            .port(13012)
                            .build())
                    .build());
        }

        downloadRequest.setUrls(urls);
//        downloadRequest.setDelay(1000);

        String body = objectMapper.writeValueAsString(downloadRequest);

        byte[] data = Unirest.post("http://127.0.0.1:8123")
                .body(body)
                .header("Content-type", "application/json")
                .asBytes()
                .getBody();

        System.out.println(data.length);

        try (GZIPInputStream gis = new GZIPInputStream(new ByteArrayInputStream(data));
             TarArchiveInputStream tais = new TarArchiveInputStream(gis)) {
            ArchiveEntry entry;

            while ((entry = tais.getNextEntry()) != null) {
                timeCalc.printSpeedStep();
            }
        }

    }

}
