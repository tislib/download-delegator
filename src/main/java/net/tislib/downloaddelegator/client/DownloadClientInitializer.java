package net.tislib.downloaddelegator.client;

import io.netty.channel.ChannelDuplexHandler;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOutboundHandlerAdapter;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.ChannelPromise;
import io.netty.channel.socket.SocketChannel;
import io.netty.handler.codec.http.HttpClientCodec;
import io.netty.handler.codec.http.HttpContentDecompressor;
import io.netty.handler.codec.http.HttpObjectAggregator;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import io.netty.handler.proxy.HttpProxyHandler;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslHandler;
import lombok.RequiredArgsConstructor;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.config.ApplicationConfig;
import net.tislib.downloaddelegator.config.Config;
import net.tislib.downloaddelegator.data.PageUrl;

import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.net.URL;
import java.util.concurrent.TimeUnit;

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

        p.addLast(new ChannelDuplexHandler(){
            @Override
            public void connect(ChannelHandlerContext ctx, SocketAddress remoteAddress, SocketAddress localAddress, ChannelPromise promise) throws Exception {
                super.connect(ctx, remoteAddress, localAddress, promise);

                ctx.executor().schedule(() -> {
                    if (ctx.channel().isOpen() || ctx.channel().isActive()) {
                        ctx.close();
                    }
                }, pageUrl.getTimeout(), TimeUnit.MILLISECONDS);
            }
        });

        if (ApplicationConfig.getBoolean(Config.TRACE_CLIENT)) {
            p.addLast(new LoggingHandler(LogLevel.ERROR));
        }

        if (pageUrl.getProxy() != null) {
            p.addLast(new HttpProxyHandler(new InetSocketAddress(pageUrl.getProxy().getHost(), pageUrl.getProxy().getPort()), pageUrl.getProxy().getUsername(), pageUrl.getProxy().getPassword()));
        }

        // Enable HTTPS if necessary.
        if (sslCtx != null) {
            SslHandler handler = sslCtx.newHandler(p.channel().alloc());

            handler.setHandshakeTimeoutMillis(pageUrl.getTimeout());

            p.addLast(handler);
        }

        p.addLast(new HttpClientCodec());

        p.addLast(new HttpContentDecompressor());

        p.addLast(new HttpObjectAggregator(1024 * 1024 * 128, true));

        p.addLast(new FullDownloadClientHandler(downloadClient, pageUrl));
    }
}
