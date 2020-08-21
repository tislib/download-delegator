package net.tislib.downloaddelegator.service;

import net.tislib.downloaddelegator.base.TimeCalc;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageResponse;
import net.tislib.downloaddelegator.data.PageUrl;
import net.tislib.downloaddelegator.util.OutputStreamPublisher;
import org.apache.commons.compress.archivers.tar.TarArchiveEntry;
import org.apache.commons.compress.archivers.tar.TarArchiveOutputStream;
import org.apache.commons.compress.compressors.gzip.GzipCompressorOutputStream;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.netty.http.client.HttpClient;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

@Service
public class DownloaderService {

    TimeCalc timeCalc = new TimeCalc();

    public Flux<DataBuffer> download(DownloadRequest downloadRequest) {
//        DataBufferFactory dataBufferFactory = new DefaultDataBufferFactory();

        OutputStreamPublisher<TarArchiveOutputStream> outputStreamPublisher = new OutputStreamPublisher<>(os -> {
            try {
                return new TarArchiveOutputStream(new GzipCompressorOutputStream(os));
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        });

        List<PageUrl> urls = new ArrayList<>(downloadRequest.getUrls());

//        for (int i = 0; i < 100; i++) {
//            urls.addAll(downloadRequest.getUrls());
//        }

        Flux<PageResponse> res = Flux.fromIterable(urls)
                .flatMap(this::download);

        return outputStreamPublisher.mapper(res, item -> {
            try {
                TarArchiveEntry entry = new TarArchiveEntry(String.valueOf(item.getId()));
                String head = item.getHttpStatus() + "\n";
                entry.setSize(head.getBytes().length + item.getContent().length);

                outputStreamPublisher.getStream().putArchiveEntry(entry);
                outputStreamPublisher.getStream().write(head.getBytes());
                outputStreamPublisher.getStream().write(item.getContent());
                outputStreamPublisher.getStream().closeArchiveEntry();
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        });

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
