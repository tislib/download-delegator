package net.tislib.downloaddelegator.test.download;

import net.tislib.downloaddelegator.test.base.BaseIntegrationTest;
import net.tislib.downloaddelegator.test.base.PageData;
import net.tislib.downloaddelegator.test.base.Scenario;
import org.junit.Test;

import java.util.List;

import static org.junit.Assert.assertEquals;

public class DelayTest extends BaseIntegrationTest {

    @Test
    public void globalDelayTest() throws Exception {
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .responseData("hello-world".getBytes())
                                .build())
                        .build()
        );

        List<PageData> response = backend.call(prepareDownloadRequest(5));

        assertEquals(response.size(), 5);

        response.forEach(item -> assertEquals(new String(item.getContent()), "hello-world"));
    }
}
