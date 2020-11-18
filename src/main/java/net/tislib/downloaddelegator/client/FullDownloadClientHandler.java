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

    @Override
    public void channelRead0(ChannelHandlerContext ctx, FullHttpResponse fullHttpResponse) {
        log.trace("received from: {} {} size: {}",
                ctx.channel().localAddress(),
                ctx.channel().remoteAddress(),
                fullHttpResponse.headers().get("Content-Length"));

        PageResponse response = new PageResponse();
        response.setContent(fullHttpResponse.content().copy());

        downloadClient.onFullResponse(response);
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) {
        cause.printStackTrace();
        ctx.close();
    }
}