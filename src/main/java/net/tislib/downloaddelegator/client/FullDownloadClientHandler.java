package net.tislib.downloaddelegator.client;

import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.FullHttpResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.data.PageResponse;

@Log4j2
@RequiredArgsConstructor
public class FullDownloadClientHandler extends SimpleChannelInboundHandler<FullHttpResponse> {
    private final DownloadClient downloadClient;
    private boolean isResponded;

    @Override
    public void channelRead0(ChannelHandlerContext ctx, FullHttpResponse fullHttpResponse) {
        log.trace("received from: {} {} size: {}",
                ctx.channel().localAddress(),
                ctx.channel().remoteAddress(),
                fullHttpResponse.headers().get("Content-Length"));

        PageResponse response = new PageResponse();
        response.setContent(fullHttpResponse.content().copy());

        downloadClient.onFullResponse(response);
        isResponded = true;
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) {
        ctx.close();
        downloadClient.onError(cause);
    }

    @Override
    public void channelUnregistered(ChannelHandlerContext ctx) throws Exception {
        super.channelUnregistered(ctx);
        downloadClient.onClose(isResponded);
    }
}
