package net.tislib.downloaddelegator.downloader;

import lombok.SneakyThrows;
import net.tislib.downloaddelegator.base.TimeCalc;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageResponse;
import net.tislib.downloaddelegator.data.PageUrl;
import reactor.core.publisher.Mono;
import reactor.netty.http.client.HttpClient;

import java.io.OutputStream;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;

public class Downloader {

    TimeCalc timeCalc = new TimeCalc();

    @SneakyThrows
    public void download(DownloadRequest downloadRequest, OutputStream responseBody) {
        List<CompletableFuture<Void>> responses = new ArrayList<>();

        for (PageUrl pageUrl : downloadRequest.getUrls()) {
            responses.add(download(pageUrl).toFuture().thenAccept(resp -> {
                onResponse(resp, responseBody);
            }));
            if (pageUrl.getDelay() > 0) {
                Thread.sleep(pageUrl.getDelay());
            }
        }

        CompletableFuture.allOf(responses.toArray(new CompletableFuture[0]))
                .get();
    }

    @SneakyThrows
    private void onResponse(PageResponse resp, OutputStream responseBody) {
        responseBody.write(resp.getContent());
    }

    private Mono<PageResponse> download(PageUrl pageUrl) {
        return HttpClient.create().baseUrl(pageUrl.getUrl().toString()).get()
                .responseSingle((httpClientResponse, byteBufMono) -> byteBufMono.log().map(buf -> {
                    PageResponse pageResponse = new PageResponse();
                    pageResponse.setHttpStatus(httpClientResponse.status().code());

                    byte[] bytes = new byte[buf.readableBytes()];
                    int readerIndex = buf.readerIndex();
                    buf.getBytes(readerIndex, bytes);

                    pageResponse.setContent(bytes);
                    pageResponse.setId(pageUrl.getId());

                    timeCalc.printSpeedStep();

                    return pageResponse;
                })).log("asd");

    }
}
