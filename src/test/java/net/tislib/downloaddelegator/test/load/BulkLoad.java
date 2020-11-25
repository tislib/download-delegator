package net.tislib.downloaddelegator.test.load;

import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import kong.unirest.HttpResponse;
import kong.unirest.Unirest;
import lombok.SneakyThrows;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.Handler;
import net.tislib.downloaddelegator.data.Header;
import net.tislib.downloaddelegator.data.PageUrl;
import net.tislib.downloaddelegator.test.base.Backend;
import net.tislib.downloaddelegator.test.base.HttpServer;
import net.tislib.downloaddelegator.test.base.Scenario;
import org.apache.commons.compress.archivers.ArchiveEntry;
import org.apache.commons.compress.archivers.tar.TarArchiveInputStream;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.core.config.Configurator;
import org.junit.Rule;
import org.junit.Test;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.IOException;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;
import java.util.concurrent.CountDownLatch;
import java.util.zip.GZIPInputStream;

public class BulkLoad {

    private static ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
    private static Backend backend = new Backend();

    @Rule
    public HttpServer httpServer = new HttpServer();

    static {
        Configurator.setRootLevel(Level.INFO);
    }

    @SneakyThrows
    @Test
    public void load1() {
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .responseData("hello-world".getBytes())
                                .build())
                        .build()
        );

        for (int i = 0; i < 1000; i++) {
            new Thread(() -> {
                try {
                    send1();
                } catch (MalformedURLException | JsonProcessingException e) {
                    e.printStackTrace();
                }
            }).start();
            Thread.sleep(1000);
        }
    }

    @SneakyThrows
    @Test
    public void load2() {
        byte[] data = new byte[1024 * 1024 * 100];

        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .responseData(data)
                                .build())
                        .build()
        );

        for (int i = 0; i < 1000; i++) {
            new Thread(() -> {
                try {
                    send2();
                } catch (MalformedURLException | JsonProcessingException e) {
                    e.printStackTrace();
                }
            }).start();
            Thread.sleep(1000);
        }
    }

    private void send1() throws MalformedURLException, JsonProcessingException {
        DownloadRequest downloadRequest = new DownloadRequest();
        List<PageUrl> urls = new ArrayList<>();

        for (int i = 0; i < 500; i++) {
            UUID pageId = UUID.randomUUID();

            urls.add(PageUrl.builder()
                    .id(pageId)
                    .url(new URL(httpServer.getUrl() + "/" + pageId))
                    .method("GET")
                    .bind("172.20.11.45")
                    .build());
        }

        downloadRequest.setUrls(urls);
        downloadRequest.setDelay(1);

        backend.call(downloadRequest);
    }

    @SneakyThrows
    private void send2() throws MalformedURLException, JsonProcessingException {
        DownloadRequest downloadRequest = new DownloadRequest();
        List<PageUrl> urls = new ArrayList<>();

        for (int i = 0; i < 5; i++) {
            UUID pageId = UUID.randomUUID();

            urls.add(PageUrl.builder()
                    .id(pageId)
                    .url(new URL(httpServer.getUrl() + "/" + pageId))
                    .method("GET")
                    .bind("172.20.11.45")
                    .build());
        }

        downloadRequest.setUrls(urls);
        downloadRequest.setDelay(1);

        CountDownLatch latch = new CountDownLatch(1);

        backend.callAsync(downloadRequest, list -> {
            list.size();
            latch.countDown();
        });

        latch.await();
    }

}
