package net.tislib.downloaddelegator.test.base;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.buffer.ByteBuf;
import io.netty.channel.*;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.codec.http.*;
import io.netty.util.ResourceLeakDetector;
import lombok.SneakyThrows;
import org.junit.rules.TestRule;
import org.junit.runner.Description;
import org.junit.runners.model.Statement;

import java.net.InetSocketAddress;
import java.net.URL;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

import static io.netty.handler.codec.http.HttpResponseStatus.OK;
import static io.netty.handler.codec.http.HttpVersion.HTTP_1_1;

public class HttpServer implements TestRule {

    static {
        ResourceLeakDetector.setLevel(ResourceLeakDetector.Level.PARANOID);
    }

    private static final EventLoopGroup serverWorkgroup = new NioEventLoopGroup();

    private ChannelFuture serverChannel;
    private Scenario scenario;
    private AtomicInteger scenarioIndex = new AtomicInteger();

    @SneakyThrows
    public void start() {
        serverChannel = this.server(serverWorkgroup).sync();
    }

    public void stop() {
        if (serverChannel != null) {
            serverChannel.channel().closeFuture();
        }
    }

    public InetSocketAddress getAddr() {
        return (InetSocketAddress) serverChannel.channel().localAddress();
    }

    public ChannelFuture server(EventLoopGroup workerGroup) {
        ServerBootstrap b = new ServerBootstrap();
        b.group(workerGroup).channel(NioServerSocketChannel.class)
                //Setting InetSocketAddress to port 0 will assign one at random
                .localAddress(new InetSocketAddress(0))
                .childHandler(new ChannelInitializer<SocketChannel>() {
                    @Override
                    protected void initChannel(SocketChannel ch) throws Exception {
                        //HttpServerCodec is a helper ChildHandler that encompasses
                        //both HTTP request decoding and HTTP response encoding
                        ch.pipeline().addLast(new HttpServerCodec());
                        //HttpObjectAggregator helps collect chunked HttpRequest pieces into
                        //a single FullHttpRequest. If you don't make use of streaming, this is
                        //much simpler to work with.
                        ch.pipeline().addLast(new HttpObjectAggregator(1048576));
                        //Finally add your FullHttpRequest handler. Real examples might replace this
                        //with a request router
                        ch.pipeline().addLast(new SimpleChannelInboundHandler<FullHttpRequest>() {
                            @Override
                            protected void channelRead0(ChannelHandlerContext ctx, FullHttpRequest msg) throws Exception {
                                process(ctx, msg);
                            }

                            @Override
                            public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {

                            }
                        });
                    }
                });

        // Start the server & bind to a random port.
        return b.bind();
    }

    private void process(ChannelHandlerContext ctx, FullHttpRequest msg) {
        scenarioIndex.incrementAndGet();
        Scenario.Request scenarioItem = locateScenarioItem(scenarioIndex.get());

        if (scenarioItem.getResponseTime() > 0) {
            ctx.executor().schedule(() ->
                    runRequest(scenarioItem, ctx, msg), scenarioItem.getResponseTime(), TimeUnit.MILLISECONDS);
        }

        runRequest(scenarioItem, ctx, msg);
    }

    private void runRequest(Scenario.Request scenarioItem, ChannelHandlerContext ctx, FullHttpRequest msg) {
        if (scenarioItem.isCloseConnectionWithoutResponse()) {
            ctx.close();
        }

        ByteBuf content = ctx.alloc().buffer();
        content.writeBytes(scenarioItem.getResponseData());

        DefaultFullHttpResponse response = new DefaultFullHttpResponse(HTTP_1_1,
                HttpResponseStatus.valueOf(scenarioItem.getStatusCode()),
                content);

        ctx.writeAndFlush(response);

        ctx.channel().close();
    }


    private Scenario.Request locateScenarioItem(final int index) {
        int curIndex = index;

        Scenario.Request lastRequest = null;

        for (Scenario.Request request : scenario.getRequests()) {
            lastRequest = request;
            curIndex = curIndex - request.getCount();

            if (curIndex < 0) {
                break;
            }
        }

        return lastRequest;
    }


    @SneakyThrows
    @Override
    public Statement apply(Statement base, Description description) {
        return new Statement() {
            @Override
            public void evaluate() throws Throwable {
                try {
                    scenario = initScenario();
                    start();
                    base.evaluate();
                } finally {
                    stop();
                }
            }
        };
    }

    private Scenario initScenario() {
        return Scenario.builder()
                .build();
    }

    public void scenario(Scenario scenario) {
        this.scenario = scenario;
    }

    @SneakyThrows
    public URL getUrl() {
        return new URL("http://127.0.0.1:" + getAddr().getPort());
    }
}
