package net.tislib.downloaddelegator.client;

import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOutboundHandlerAdapter;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.socket.SocketChannel;
import io.netty.handler.codec.http.HttpClientCodec;
import io.netty.handler.codec.http.HttpContentDecompressor;
import io.netty.handler.codec.http.HttpObjectAggregator;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslHandler;
import lombok.RequiredArgsConstructor;
import lombok.extern.log4j.Log4j2;

@Log4j2
@RequiredArgsConstructor
public class DownloadClientInitializer extends ChannelInitializer<SocketChannel> {
    private final SslContext sslCtx;
    private final DownloadClient downloadClient;

    @Override
    protected void initChannel(SocketChannel ch) {
        ChannelPipeline p = ch.pipeline();
        log.debug("connected to: {} {}", ch.localAddress(), ch.remoteAddress());

        p.addLast(new ChannelOutboundHandlerAdapter());

        // Enable HTTPS if necessary.
        if (sslCtx != null) {
            SslHandler handler = sslCtx.newHandler(p.channel().alloc());

            handler.setHandshakeTimeoutMillis(500000);

            p.addLast(handler);
        }

        p.addLast(new HttpClientCodec());

//        p.addLast(new LoggingHandler(LogLevel.ERROR));

        p.addLast(new HttpContentDecompressor());

        p.addLast(new HttpObjectAggregator(1024 * 1024, true));

        p.addLast(new FullDownloadClientHandler(downloadClient));
    }
}
