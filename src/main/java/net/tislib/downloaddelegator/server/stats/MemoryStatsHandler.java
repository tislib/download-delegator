package net.tislib.downloaddelegator.server.stats;

import io.netty.buffer.PooledByteBufAllocator;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.netty.handler.codec.http.FullHttpResponse;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http.HttpVersion;
import net.tislib.downloaddelegator.data.MemoryStatsRequest;

public class MemoryStatsHandler extends SimpleChannelInboundHandler<MemoryStatsRequest> {
    @Override
    protected void channelRead0(ChannelHandlerContext ctx, MemoryStatsRequest msg) throws Exception {
        FullHttpResponse fullHttpResponse = new DefaultFullHttpResponse(
                HttpVersion.HTTP_1_1,
                HttpResponseStatus.OK
        );

        String content = getMemoryStats(ctx, msg);

        fullHttpResponse.content().writeBytes(content.getBytes());
        ctx.writeAndFlush(fullHttpResponse);
        ctx.close();
    }

    private String getMemoryStats(ChannelHandlerContext ctx, MemoryStatsRequest msg) {
        PooledByteBufAllocator pooledByteBufAllocator = (PooledByteBufAllocator) ctx.alloc();
        return String.format("Memory stats dump: %s;  %s", pooledByteBufAllocator.toString(), pooledByteBufAllocator.metric().toString());
    }
}
