package net.tislib.downloaddelegator.client;

import io.netty.bootstrap.Bootstrap;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelFuture;
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

import java.net.InetSocketAddress;
import java.net.URI;
import java.net.URL;

@Log4j2
public abstract class DownloadClient {

    @SneakyThrows
    public ChannelFuture connect(URL url) {
        URI uri = new URI(url.toString());
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
                    .sessionTimeout(500000)
                    .trustManager(InsecureTrustManagerFactory.INSTANCE).build();
        } else {
            sslCtx = null;
        }

        bootstrap = bootstrap.handler(new DownloadClientInitializer(sslCtx, this));

        if (!"http".equalsIgnoreCase(scheme) && !"https".equalsIgnoreCase(scheme)) {
            throw new UnsupportedOperationException("Only HTTP(S) is supported.");
        }

        return bootstrap.connect(new InetSocketAddress(host, port));
    }

    @SneakyThrows
    public void download(ChannelFuture channelFuture, URL url) {
        URI uri = new URI(url.toString());
        // Prepare the HTTP request.
        HttpRequest request = new DefaultFullHttpRequest(
                HttpVersion.HTTP_1_1, HttpMethod.GET, uri.getRawPath(), Unpooled.EMPTY_BUFFER);

        request.headers().set(HttpHeaderNames.HOST, uri.getHost());
        request.headers().set(HttpHeaderNames.CONNECTION, HttpHeaderValues.CLOSE);
        request.headers().set(HttpHeaderNames.ACCEPT_ENCODING, HttpHeaderValues.GZIP);

        channelFuture.addListener(future -> {
            channelFuture.channel().writeAndFlush(request);
        });
    }

    public abstract void onFullResponse(PageResponse pageResponse);

    public void onError(Throwable th) {
        log.error("onError", th);
    }

    public void onClose(boolean isResponded) {

    }
}
