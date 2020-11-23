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

public class FailTest extends BaseIntegrationTest {

    @Test
    public void failDownloadTest() throws Exception {
//        Thread.sleep(100000000);
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .closeConnectionWithoutResponse(true)
                                .count(3)
                                .build())
                        .request(Scenario.Request.builder()
                                .responseData("hello-world".getBytes())
                                .build())
                        .build()
        );

        List<PageData> response = backend.call(prepareDownloadRequest(5));

        assertEquals(response.size(), 90);

        response.forEach(item -> assertEquals(new String(item.getContent()), "hello-world"));
    }
}
