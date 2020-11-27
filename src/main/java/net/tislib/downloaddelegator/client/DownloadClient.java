package net.tislib.downloaddelegator.client;

import io.netty.bootstrap.Bootstrap;
import io.netty.buffer.PooledByteBufAllocator;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelOption;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.handler.codec.http.*;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslContextBuilder;
import io.netty.handler.ssl.util.InsecureTrustManagerFactory;
import lombok.SneakyThrows;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.base.EventLoopGroups;
import net.tislib.downloaddelegator.config.ApplicationConfig;
import net.tislib.downloaddelegator.config.Config;
import net.tislib.downloaddelegator.data.PageResponse;
import net.tislib.downloaddelegator.data.PageUrl;

import java.net.InetSocketAddress;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;

@Log4j2
public abstract class DownloadClient {

    @SneakyThrows
    public ChannelFuture connect(PageUrl pageUrl) {
        URI uri = new URI(pageUrl.getUrl().toString());
        String scheme = uri.getScheme() == null ? "http" : uri.getScheme();
        String host = uri.getHost() == null ? "127.0.0.1" : uri.getHost();
        int port = uri.getPort();
        if (port == -1) {
            if ("http".equalsIgnoreCase(scheme)) {
                port = 80;
            } else if ("https".equalsIgnoreCase(scheme)) {
                port = 443;
            }
        }

        Bootstrap bootstrap = new Bootstrap();

        bootstrap = bootstrap.channel(NioSocketChannel.class);

        bootstrap = bootstrap.group(EventLoopGroups.clientGroup);

        if (ApplicationConfig.getBoolean(Config.TRACE_CLIENT)) {
            bootstrap = bootstrap.handler(new LoggingHandler(LogLevel.ERROR));
        }

        // Configure SSL context if necessary.
        final boolean ssl = "https".equalsIgnoreCase(scheme);
        final SslContext sslCtx;
        if (ssl) {
            sslCtx = SslContextBuilder.forClient()
                    .sessionTimeout(pageUrl.getTimeout())
                    .trustManager(InsecureTrustManagerFactory.INSTANCE)
                    .build();
        } else {
            sslCtx = null;
        }

        if (pageUrl.getTimeout() <= 0) {
            pageUrl.setTimeout(Integer.parseInt(ApplicationConfig.getConfig(Config.TIMEOUT)));
        }

        bootstrap.option(ChannelOption.CONNECT_TIMEOUT_MILLIS, pageUrl.getTimeout());

        bootstrap = bootstrap.handler(new DownloadClientInitializer(sslCtx, this, pageUrl));

        if (!"http".equalsIgnoreCase(scheme) && !"https".equalsIgnoreCase(scheme)) {
            throw new UnsupportedOperationException("Only HTTP(S) is supported.");
        }

        return bootstrap.connect(new InetSocketAddress(host, port));
    }

    @SneakyThrows
    public void download(ChannelFuture channelFuture, PageUrl pageUrl) {
        channelFuture.addListener(future -> sendRequest(channelFuture, pageUrl));
    }

    private void sendRequest(ChannelFuture channelFuture,
                             PageUrl pageUrl) throws URISyntaxException {
        URI uri = new URI(pageUrl.getUrl().toString());
        // Prepare the HTTP request.
        HttpRequest request = new DefaultFullHttpRequest(
                HttpVersion.HTTP_1_1, HttpMethod.GET, uri.getRawPath(), PooledByteBufAllocator.defaultPreferDirect());

        request.headers().set(HttpHeaderNames.HOST, uri.getHost());
        request.headers().set(HttpHeaderNames.CONNECTION, HttpHeaderValues.CLOSE);
        request.headers().set(HttpHeaderNames.ACCEPT_ENCODING, HttpHeaderValues.GZIP);

        if (pageUrl.getHeaders() != null) {
            pageUrl.getHeaders().forEach(header -> {
                request.headers().set(header.getName(), header.getValue());
            });
        }
        channelFuture.channel().writeAndFlush(request);
    }

    public void onError(Throwable th) {
        log.error("onError", th);
    }

    public abstract void onFinish(PageResponse pageResponse);
}
