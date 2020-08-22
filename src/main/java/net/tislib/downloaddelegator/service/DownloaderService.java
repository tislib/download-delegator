package net.tislib.downloaddelegator.service;

import io.netty.buffer.ByteBuf;
import lombok.SneakyThrows;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.base.TimeCalc;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageResponse;
import net.tislib.downloaddelegator.data.PageUrl;
import net.tislib.downloaddelegator.util.OutputStreamPublisher;
import org.apache.commons.compress.archivers.tar.TarArchiveEntry;
import org.apache.commons.compress.archivers.tar.TarArchiveOutputStream;
import org.apache.commons.compress.compressors.gzip.GzipCompressorOutputStream;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.netty.http.client.HttpClient;
import reactor.netty.http.client.HttpClientResponse;
import reactor.netty.tcp.ProxyProvider;
import reactor.netty.tcp.TcpClient;

import java.io.IOException;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;

@Service
@Log4j2
public class DownloaderService {

    TimeCalc timeCalc = new TimeCalc();

    @Value("${downloaddelegator.page.requestTimeout}")
    private int requestTimeout;

    @Value("${downloaddelegator.page.delay}")
    private int delay;


    public Flux<DataBuffer> download(DownloadRequest downloadRequest) {
        OutputStreamPublisher<TarArchiveOutputStream> outputStreamPublisher = new OutputStreamPublisher<>(os -> {
            try {
                return new TarArchiveOutputStream(new GzipCompressorOutputStream(os));
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        });

        if (downloadRequest.getDelay() < 1) {
            downloadRequest.setDelay(delay);
        }
        if (downloadRequest.getTimeout() < 1) {
            downloadRequest.setTimeout(1);
        }

        List<PageUrl> urls = new ArrayList<>(downloadRequest.getUrls());

        Flux<PageResponse> res = Flux.fromIterable(urls)
                .delayElements(Duration.ofMillis(downloadRequest.getDelay()))
                .flatMap(this::download);

        return outputStreamPublisher.mapper(res, item -> mapResponseToTar(outputStreamPublisher, item));
    }

    @SneakyThrows
    private void mapResponseToTar(OutputStreamPublisher<TarArchiveOutputStream> outputStreamPublisher, PageResponse item) {
        if (item.getId() != null) {
            TarArchiveEntry entry = new TarArchiveEntry(String.valueOf(item.getId()));
            String head = item.getHttpStatus() + "\n";
            entry.setSize(head.getBytes().length + item.getContent().length);

            outputStreamPublisher.getStream().putArchiveEntry(entry);
            outputStreamPublisher.getStream().write(head.getBytes());
            outputStreamPublisher.getStream().write(item.getContent());
            outputStreamPublisher.getStream().closeArchiveEntry();
        }
    }

    private Mono<PageResponse> download(PageUrl pageUrl) {
        return HttpClient.create().baseUrl(pageUrl.getUrl().toString())
                .tcpConfiguration(tcpClient -> proxyConfig(pageUrl, tcpClient))
                .responseTimeout(Duration.ofMillis(requestTimeout))
                .get()
                .responseSingle((httpClientResponse, byteBufMono) -> byteBufMono.log()
                        .map(buf -> mapResponse(pageUrl, httpClientResponse, buf))
                ).onErrorResume(error -> {
                    log.error(String.format("page download failed ID: %s URL: %s", pageUrl.getId(), pageUrl.getUrl()), error);
                    return Mono.just(new PageResponse());
                });
    }

    private PageResponse mapResponse(PageUrl pageUrl, HttpClientResponse httpClientResponse, ByteBuf buf) {
        PageResponse pageResponse = new PageResponse();

        pageResponse.setHttpStatus(httpClientResponse.status().code());

        byte[] bytes = new byte[buf.readableBytes()];
        int readerIndex = buf.readerIndex();
        buf.getBytes(readerIndex, bytes);

        pageResponse.setContent(bytes);
        pageResponse.setId(pageUrl.getId());

        timeCalc.printSpeedStep();

        return pageResponse;
    }

    private TcpClient proxyConfig(final PageUrl pageUrl, final TcpClient tcpClient) {
        if (pageUrl.getProxy() == null) {
            return tcpClient;
        }

        return tcpClient.proxy(ops -> buildProxy(pageUrl, ops));
    }

    private ProxyProvider.Builder buildProxy(PageUrl pageUrl, final ProxyProvider.TypeSpec ops) {
        ProxyProvider.Builder config1 = ops.type(ProxyProvider.Proxy.HTTP).host(pageUrl.getProxy().getHost()).port(pageUrl.getProxy().getPort());
        if (pageUrl.getProxy().getUsername() != null) {
            return config1.username(pageUrl.getProxy().getUsername()).password(u -> pageUrl.getProxy().getPassword());
        } else {
            return config1;
        }
    }
}
