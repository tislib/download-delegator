package net.tislib.downloaddelegator.test.download;

import net.tislib.downloaddelegator.test.base.BaseIntegrationTest;
import net.tislib.downloaddelegator.test.base.Scenario;
import org.junit.Test;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

import static org.hamcrest.core.Is.is;

public class ConcurrentTest extends BaseIntegrationTest {

    @Test
    public void concurrentDownloadTest() throws InterruptedException {
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .responseData("hello-world".getBytes())
                                .build())
                        .build()
        );

        int concurrentCount = 10;

        CountDownLatch latch = new CountDownLatch(concurrentCount);

        for (int i = 0; i < concurrentCount; i++) {
            backend.callAsync(prepareDownloadRequest(5), response -> {
                try {
                    response.forEach(item -> collector.checkThat(new String(item.getContent()), is("hello-world")));
                } finally {
                    latch.countDown();
                }
            });
        }

        if (!latch.await(10, TimeUnit.SECONDS)) {
            throw new RuntimeException("operation is not completed");
        }
    }

}
