package net.tislib.downloaddelegator.server;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.ChannelFuture;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.base.EventLoopGroups;
import net.tislib.downloaddelegator.config.ApplicationConfig;
import net.tislib.downloaddelegator.config.Config;

import java.net.InetSocketAddress;

@Log4j2
public class Server {
    private ServerBootstrap bootstrap;
    private boolean running;
    private ChannelFuture channelFuture;

    public Server() {
        bootstrap = new ServerBootstrap();
    }

    public synchronized void run() {
        if (running) {
            throw new IllegalStateException("server is already running");
        }

        log.trace("using nio server socket channel");
        bootstrap = bootstrap.channel(NioServerSocketChannel.class);

        log.trace("using default serverGroup event pool");
        bootstrap = bootstrap.group(EventLoopGroups.serverGroup);

        if (ApplicationConfig.getBoolean(Config.TRACE_SERVER)) {
            bootstrap = bootstrap.handler(new LoggingHandler(LogLevel.ERROR));
        }

        bootstrap.childHandler(new ServerChannelInitializer());

        String bindAddr = ApplicationConfig.getConfig(Config.ADDR);
        int bindPort = Integer.parseInt(ApplicationConfig.getConfig(Config.PORT));

        log.trace("listening to address: {}:{}", bindAddr, bindPort);
        channelFuture = bootstrap.bind(new InetSocketAddress(bindAddr, bindPort));

        try {
            channelFuture.sync();
        } catch (InterruptedException e) {
            log.error(e.getMessage(), e);
            channelFuture.channel().closeFuture();
        }

        running = true;
    }

    public void stop() {
        channelFuture.channel().closeFuture();
    }
}
