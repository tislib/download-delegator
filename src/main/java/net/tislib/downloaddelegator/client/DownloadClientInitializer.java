package net.tislib.downloaddelegator.client;

import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOutboundHandlerAdapter;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.socket.SocketChannel;
import io.netty.handler.codec.http.HttpClientCodec;
import io.netty.handler.codec.http.HttpContentDecompressor;
import io.netty.handler.codec.http.HttpObjectAggregator;
import io.netty.handler.proxy.HttpProxyHandler;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslHandler;
import lombok.RequiredArgsConstructor;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.data.PageUrl;

import java.net.InetSocketAddress;
import java.net.URL;

@Log4j2
@RequiredArgsConstructor
public class DownloadClientInitializer extends ChannelInitializer<SocketChannel> {
    private final SslContext sslCtx;
    private final DownloadClient downloadClient;
    private final PageUrl pageUrl;

    @Override
    protected void initChannel(SocketChannel ch) {
        ChannelPipeline p = ch.pipeline();
        log.debug("connected to: {} {} {}", ch.localAddress(), ch.remoteAddress(), pageUrl.getId());

        p.addLast(new ChannelOutboundHandlerAdapter());

        // Enable HTTPS if necessary.
        if (sslCtx != null) {
            SslHandler handler = sslCtx.newHandler(p.channel().alloc());

            handler.setHandshakeTimeoutMillis(pageUrl.getTimeout());

            p.addLast(handler);
        }

        if (pageUrl.getProxy() != null) {
            p.addLast(new HttpProxyHandler(new InetSocketAddress(pageUrl.getProxy().getHost(), pageUrl.getProxy().getPort()), pageUrl.getProxy().getUsername(), pageUrl.getProxy().getPassword()));
        }

        p.addLast(new HttpClientCodec());

        p.addLast(new HttpContentDecompressor());

        p.addLast(new HttpObjectAggregator(1024 * 1024, true));

        p.addLast(new FullDownloadClientHandler(downloadClient, pageUrl));
    }
}
