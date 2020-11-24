package net.tislib.downloaddelegator.test.download;

import net.tislib.downloaddelegator.test.base.BaseIntegrationTest;
import net.tislib.downloaddelegator.test.base.PageData;
import net.tislib.downloaddelegator.test.base.Repeat;
import net.tislib.downloaddelegator.test.base.RepeatRule;
import net.tislib.downloaddelegator.test.base.Scenario;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.core.config.Configurator;
import org.hamcrest.core.Is;
import org.junit.Ignore;
import org.junit.Rule;
import org.junit.Test;

import java.util.Arrays;
import java.util.List;
import java.util.Objects;

import static org.junit.Assert.assertEquals;

public class FailTest extends BaseIntegrationTest {

    @Test
    public void failDownloadTest() throws Exception {
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .scenarioKind(Scenario.ScenarioKind.CLOSE_CONNECTION)
                                .count(10)
                                .build())
                        .request(Scenario.Request.builder()
                                .responseData("hello-world".getBytes())
                                .build())
                        .build()
        );

        List<PageData> response = backend.call(prepareDownloadRequest(100));

        assertEquals(100, response.size());

        assertEquals(90, response.stream().filter(item -> Arrays.equals(item.getContent(), "hello-world".getBytes())).count());
        assertEquals(10, response.stream().filter(item -> Arrays.equals(item.getContent(), "".getBytes())).count());
    }

    @Test
    public void failDownloadTest2() throws Exception {
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .count(10)
                                .responseData("hello-world".getBytes())
                                .build())
                        .request(Scenario.Request.builder()
                                .scenarioKind(Scenario.ScenarioKind.CLOSE_CONNECTION)
                                .count(10)
                                .build())
                        .request(Scenario.Request.builder()
                                .count(10)
                                .responseData("hello-world".getBytes())
                                .build())
                        .request(Scenario.Request.builder()
                                .scenarioKind(Scenario.ScenarioKind.CLOSE_CONNECTION)
                                .count(10)
                                .build())
                        .request(Scenario.Request.builder()
                                .responseData("hello-world".getBytes())
                                .build())
                        .build()
        );

        List<PageData> response = backend.call(prepareDownloadRequest(100));

        assertEquals(response.size(), 100);

        assertEquals(80, response.stream().filter(item -> Arrays.equals(item.getContent(), "hello-world".getBytes())).count());
        assertEquals(20, response.stream().filter(item -> Arrays.equals(item.getContent(), "".getBytes())).count());
    }

    @Test
    public void failDownloadTestHttpError() throws Exception {
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .count(10)
                                .responseData("hello-world".getBytes())
                                .build())
                        .request(Scenario.Request.builder()
                                .scenarioKind(Scenario.ScenarioKind.CLOSE_HTTP)
                                .count(10)
                                .build())
                        .request(Scenario.Request.builder()
                                .count(10)
                                .responseData("hello-world".getBytes())
                                .build())
                        .request(Scenario.Request.builder()
                                .scenarioKind(Scenario.ScenarioKind.CORRUPT_HTTP)
                                .count(10)
                                .build())
                        .request(Scenario.Request.builder()
                                .responseData("hello-world".getBytes())
                                .build())
                        .build()
        );

        List<PageData> response = backend.call(prepareDownloadRequest(100));

        assertEquals(response.size(), 100);

        assertEquals(80, response.stream().filter(item -> Arrays.equals(item.getContent(), "hello-world".getBytes())).count());
        assertEquals(20, response.stream().filter(item -> Arrays.equals(item.getContent(), "".getBytes())).count());
    }

    @Test
    public void failDownloadTestHttpError2() throws Exception {
        httpServer.scenario(
                Scenario.builder()
                        .request(Scenario.Request.builder()
                                .scenarioKind(Scenario.ScenarioKind.CLOSE_CONNECTION)
                                .count(10)
                                .build())
                        .request(Scenario.Request.builder()
                                .scenarioKind(Scenario.ScenarioKind.CORRUPT_HTTP)
                                .count(10)
                                .build())
                        .build()
        );

        List<PageData> response = backend.call(prepareDownloadRequest(100));

        assertEquals(100, response.size());

        response.forEach(item -> assertEquals(new String(item.getContent()), ""));
    }
}
