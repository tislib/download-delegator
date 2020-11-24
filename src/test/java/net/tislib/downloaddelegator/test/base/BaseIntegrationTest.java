package net.tislib.downloaddelegator.test.base;

import lombok.SneakyThrows;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageUrl;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.core.config.Configurator;
import org.junit.Rule;
import org.junit.rules.ErrorCollector;
import org.junit.rules.ExpectedException;

import java.net.URL;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

public class BaseIntegrationTest {

    @Rule
    public HttpServer httpServer = new HttpServer();

    @Rule
    public ExpectedException exception = ExpectedException.none();

    @Rule
    public ErrorCollector collector = new ErrorCollector();

    public static final Backend backend;

    static {
        backend = new Backend();
        Configurator.setRootLevel(Level.TRACE);
    }

    @SneakyThrows
    protected DownloadRequest prepareDownloadRequest(int count) {
        DownloadRequest downloadRequest = new DownloadRequest();
        List<PageUrl> urls = new ArrayList<>();

        for (int i = 0; i < count; i++) {
            UUID pageId = UUID.randomUUID();
            urls.add(PageUrl.builder()
                    .id(pageId)
                    .url(new URL(httpServer.getUrl() + "/" + pageId))
                    .method("GET")
                    .build());
        }

        downloadRequest.setUrls(urls);
//        downloadRequest.setDelay(1000);

        return downloadRequest;
    }

}
