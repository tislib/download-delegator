package net.tislib.downloaddelegator.server.stats;

import io.netty.buffer.ByteBufAllocatorMetricProvider;
import io.netty.buffer.PooledByteBufAllocator;
import io.netty.buffer.UnpooledByteBufAllocator;
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

        String content = getMemoryStats((ByteBufAllocatorMetricProvider) ctx.alloc(), msg);
        content += getMemoryStats(PooledByteBufAllocator.DEFAULT, msg);

        fullHttpResponse.content().writeBytes(content.getBytes());
        ctx.writeAndFlush(fullHttpResponse);
        ctx.close();
    }

    private String getMemoryStats(ByteBufAllocatorMetricProvider metricProvider, MemoryStatsRequest msg) {
        StringBuilder stringBuilder = new StringBuilder();

        String content = String.format("Memory stats dump: %s;  %s", metricProvider.toString(), metricProvider.metric().toString());

        stringBuilder.append("<pre>");
        stringBuilder.append(content);
        stringBuilder.append("</pre>");

        if (metricProvider instanceof PooledByteBufAllocator) {
            stringBuilder.append("<pre>");
            stringBuilder.append(((PooledByteBufAllocator) metricProvider).dumpStats());
            stringBuilder.append("</pre>");
        }

        return stringBuilder.toString();
    }
}
