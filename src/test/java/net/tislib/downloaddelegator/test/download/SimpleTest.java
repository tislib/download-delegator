package net.tislib.downloaddelegator.test.download;

import net.tislib.downloaddelegator.test.base.BaseIntegrationTest;
import net.tislib.downloaddelegator.test.base.PageData;
import net.tislib.downloaddelegator.test.base.Scenario;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.core.config.Configurator;
import org.junit.Ignore;
import org.junit.Test;

import java.util.List;

import static org.junit.Assert.assertEquals;

public class SimpleTest extends BaseIntegrationTest {

    @Test
    public void singleDownloadTest() throws Exception {
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .responseData("hello-world".getBytes())
                                .build())
                        .build()
        );

        List<PageData> response = backend.call(prepareDownloadRequest(1));

        assertEquals(response.size(), 1);

        response.forEach(item -> assertEquals(new String(item.getContent()), "hello-world"));
    }

    @Test
    public void manyDownloadTest() throws Exception {
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .responseData("hello-world".getBytes())
                                .build())
                        .build()
        );

        List<PageData> response = backend.call(prepareDownloadRequest(50));

        assertEquals(response.size(), 50);

        response.forEach(item -> assertEquals(new String(item.getContent()), "hello-world"));
    }

    @Test
    public void multiDownloadTest() throws Exception {
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .responseData("hello-world".getBytes())
                                .build())
                        .build()
        );

        for (int i = 0; i < 100; i++) {
            List<PageData> response = backend.call(prepareDownloadRequest(50));

            assertEquals(response.size(), 50);
        }
    }

    @Test
    @Ignore
    public void simpleMassiveRequests() {
        try {
            Configurator.setRootLevel(Level.ERROR);

            httpServer.scenario(
                    Scenario.builder()
                            .request(Scenario.Request.builder()
                                    .responseData("hello-world".getBytes())
                                    .build())
                            .build()
            );

            List<PageData> response = backend.call(prepareDownloadRequest(5000));

            assertEquals(response.size(), 5000);

            response.forEach(item -> assertEquals(new String(item.getContent()), "hello-world"));
        } finally {
            Configurator.setRootLevel(Level.TRACE);
        }
    }

}
