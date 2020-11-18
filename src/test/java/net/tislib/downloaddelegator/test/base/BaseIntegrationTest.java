package net.tislib.downloaddelegator.test.base;

import lombok.SneakyThrows;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageUrl;
import org.junit.Rule;

import java.net.URL;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

public class BaseIntegrationTest {

    @Rule
    public HttpServer httpServer = new HttpServer();

    @Rule
    public Backend backend = new Backend();

    @SneakyThrows
    protected DownloadRequest prepareDownloadRequest(int count) {
        DownloadRequest downloadRequest = new DownloadRequest();
        List<PageUrl> urls = new ArrayList<>();

        for (int i = 0; i < count; i++) {
            urls.add(PageUrl.builder()
                    .id(UUID.randomUUID())
                    .url(httpServer.getUrl())
                    .method("GET")
                    .build());
        }

        downloadRequest.setUrls(urls);
//        downloadRequest.setDelay(1000);

        return downloadRequest;
    }

}
