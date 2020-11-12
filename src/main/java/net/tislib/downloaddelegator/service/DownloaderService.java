package net.tislib.downloaddelegator.service;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelOption;
import lombok.SneakyThrows;
import lombok.extern.log4j.Log4j2;
import lombok.val;
import net.tislib.downloaddelegator.base.TimeCalc;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageResponse;
import net.tislib.downloaddelegator.data.PageUrl;
import net.tislib.downloaddelegator.util.OutputStreamPublisher;
import org.apache.commons.compress.archivers.tar.TarArchiveEntry;
import org.apache.commons.compress.archivers.tar.TarArchiveOutputStream;
import org.apache.commons.compress.compressors.gzip.GzipCompressorInputStream;
import org.apache.commons.compress.compressors.gzip.GzipCompressorOutputStream;
import org.apache.commons.compress.utils.IOUtils;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.netty.http.client.HttpClient;
import reactor.netty.http.client.HttpClientResponse;
import reactor.netty.resources.ConnectionProvider;
import reactor.netty.tcp.ProxyProvider;
import reactor.netty.tcp.TcpClient;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.net.InetSocketAddress;
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
            downloadRequest.setTimeout(requestTimeout);
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
            entry.setSize(item.getContent().length);

            outputStreamPublisher.getStream().putArchiveEntry(entry);
            outputStreamPublisher.getStream().write(item.getContent());
            outputStreamPublisher.getStream().closeArchiveEntry();

            TarArchiveEntry infoEntry = new TarArchiveEntry(item.getId() + ".info");
            StringBuilder header = new StringBuilder();

            header.append(item.getHttpStatus()).append("\n");

            item.getHeaders().forEach((val) -> header.append(val.getKey()).append(": ").append(val.getValue()).append("\n"));

            String headerData = header.toString();

            infoEntry.setSize(headerData.getBytes().length);

            outputStreamPublisher.getStream().putArchiveEntry(infoEntry);
            outputStreamPublisher.getStream().write(headerData.getBytes());
            outputStreamPublisher.getStream().closeArchiveEntry();
        }
    }

    private Mono<PageResponse> download(PageUrl pageUrl) {
        int requestTimeoutLocal = requestTimeout;

        val provider = ConnectionProvider.newConnection();
        TcpClient tcpClient = TcpClient.create(provider);

        tcpClient = tcpConfig(pageUrl, tcpClient);
//        tcpClient = tcpClient.option(ChannelOption.SO_TIMEOUT, 50000);
        tcpClient = tcpClient.option(ChannelOption.CONNECT_TIMEOUT_MILLIS, 500000);

        return HttpClient.from(tcpClient)
                .baseUrl(pageUrl.getUrl().toString())
                .followRedirect(true)

                .responseTimeout(Duration.ofMillis(requestTimeoutLocal))
                .headers(item -> {
                    if (pageUrl.getHeaders() != null) {
                        pageUrl.getHeaders().forEach(h -> {
                            item.add(h.getName(), h.getValue());
                        });
                    }
                })
                .get()
                .responseSingle((httpClientResponse, byteBufMono) -> byteBufMono.map(buf -> mapResponse(pageUrl, httpClientResponse, buf))
                ).onErrorResume(error -> {
                    log.error(String.format("page download failed ID: %s URL: %s message: %s", pageUrl.getId(), pageUrl.getUrl(), error.getMessage()));
                    return Mono.just(new PageResponse());
                });
    }

    private PageResponse mapResponse(PageUrl pageUrl, HttpClientResponse httpClientResponse, ByteBuf buf) {
        PageResponse pageResponse = new PageResponse();

        pageResponse.setHttpStatus(httpClientResponse.status().code());

        byte[] bytes = new byte[buf.readableBytes()];
        int readerIndex = buf.readerIndex();
        buf.getBytes(readerIndex, bytes);
        buf.release();

        String contentEncoding = httpClientResponse.responseHeaders().get("content-encoding");

        if (contentEncoding != null && contentEncoding.contains("gzip")) {
            pageResponse.setContent(gzipDecompress(bytes));
        } else {
            pageResponse.setContent(bytes);
        }

        pageResponse.setHeaders(httpClientResponse.responseHeaders().entries());

        pageResponse.setId(pageUrl.getId());

        timeCalc.printSpeedStep();

        return pageResponse;
    }

    @SneakyThrows
    private byte[] gzipDecompress(byte[] compressed) {
        try (GzipCompressorInputStream gzipCompressorInputStream = new GzipCompressorInputStream(new ByteArrayInputStream(compressed))) {
            ByteArrayOutputStream boas = new ByteArrayOutputStream();
            IOUtils.copy(gzipCompressorInputStream, boas);

            return boas.toByteArray();
        }
    }

    private TcpClient tcpConfig(final PageUrl pageUrl, TcpClient tcpClient) {
        if (pageUrl.getBind() != null) {
            tcpClient = tcpClient.bindAddress(() -> new InetSocketAddress(pageUrl.getBind(), (int) (30000 + Math.random() * 30000)));
        }

        if (pageUrl.getProxy() != null) {
            return tcpClient.proxy(ops -> buildProxy(pageUrl, ops));
        }

        return tcpClient;
    }

    private ProxyProvider.Builder buildProxy(PageUrl pageUrl, final ProxyProvider.TypeSpec ops) {
        ProxyProvider.Builder config1 = ops.type(ProxyProvider.Proxy.HTTP)
                .host(pageUrl.getProxy().getHost())
                .port(pageUrl.getProxy().getPort());

        config1 = config1.connectTimeoutMillis(500000);

        if (pageUrl.getProxy().getUsername() != null) {
            return config1.username(pageUrl.getProxy().getUsername())
                    .password(u -> pageUrl.getProxy().getPassword());
        } else {
            return config1;
        }
    }
}
