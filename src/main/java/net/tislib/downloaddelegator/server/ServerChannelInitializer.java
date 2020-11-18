package net.tislib.downloaddelegator.server;

import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.socket.SocketChannel;
import io.netty.handler.codec.http.HttpContentCompressor;
import io.netty.handler.codec.http.HttpObjectAggregator;
import io.netty.handler.codec.http.HttpRequestDecoder;
import io.netty.handler.codec.http.HttpResponseEncoder;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import lombok.RequiredArgsConstructor;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.config.ApplicationConfig;
import net.tislib.downloaddelegator.config.Config;
import net.tislib.downloaddelegator.server.downloader.PageDownloadHandler;
import net.tislib.downloaddelegator.server.downloader.PageUrlTaskSplitterHandler;

@Log4j2
@RequiredArgsConstructor
public class ServerChannelInitializer extends ChannelInitializer<SocketChannel> {

    @Override
    protected void initChannel(SocketChannel ch) {
        log.debug("init channel for connection: {} {}", ch.localAddress(), ch.remoteAddress());

        ChannelPipeline p = ch.pipeline();

        p.addLast(new HttpRequestDecoder());
        p.addLast(new HttpResponseEncoder());

        if (ApplicationConfig.getBoolean(Config.TRACE_SERVER)) {
            p.addLast(new LoggingHandler(LogLevel.ERROR));
        }

        p.addLast(new HttpContentCompressor());
        p.addLast(new HttpObjectAggregator(1024 * 1024, true));

        // switch operation and convert event to operation specific event
        p.addLast("operationSwitcher", new OperationSwitchHandler());

        // handle download request
        p.addLast(new PageUrlTaskSplitterHandler());
        p.addLast(new PageDownloadHandler());
    }
}
