package net.tislib.downloaddelegator.netty;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.codec.http.HttpRequestDecoder;
import io.netty.handler.codec.http.HttpResponseEncoder;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import lombok.SneakyThrows;
import lombok.extern.log4j.Log4j2;

import java.net.InetSocketAddress;

//@Log4j2
public class HttpServer {

    private int port;

    // constructor

    // main method, same as simple protocol server

    public void run() throws Exception {
        ServerBootstrap b = new ServerBootstrap();


        EventLoopGroup group = new NioEventLoopGroup();

        b.channel(NioServerSocketChannel.class)
                .group(group)
                .handler(new LoggingHandler(LogLevel.INFO))
                .childHandler(new ChannelInitializer<SocketChannel>() {
                    @Override
                    protected void initChannel(SocketChannel ch) {
                        ChannelPipeline p = ch.pipeline();
                        p.addLast(new HttpRequestDecoder());
                        p.addLast(new HttpResponseEncoder());
                        p.addLast(new CustomHttpServerHandler());
                    }
                });

        ChannelFuture channel = b.bind(new InetSocketAddress(8123));

        channel.sync();


    }

    @SneakyThrows
    public static void main(String[] args) {
        new HttpServer().run();
    }
}
