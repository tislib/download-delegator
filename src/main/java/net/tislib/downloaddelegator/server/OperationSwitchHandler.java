package net.tislib.downloaddelegator.server;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.MessageToMessageCodec;
import io.netty.handler.codec.MessageToMessageDecoder;
import io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.netty.handler.codec.http.FullHttpRequest;
import io.netty.handler.codec.http.FullHttpResponse;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http.HttpVersion;
import io.netty.util.ReferenceCountUtil;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.MemoryStatsRequest;

import java.nio.charset.Charset;
import java.util.List;

@Log4j2
public class OperationSwitchHandler extends MessageToMessageDecoder<FullHttpRequest> {

    private static final ObjectMapper objectMapper = new ObjectMapper();

    @Override
    protected void decode(ChannelHandlerContext ctx, FullHttpRequest request, List<Object> out) throws Exception {
        String action = request.method().name() + " " + request.uri();
        String body = request.content().toString(Charset.defaultCharset());

        ReferenceCountUtil.touch(request);

        log.debug("Handling request: {} {}", action, ctx.channel().remoteAddress());

        if (action.startsWith("POST /download")) {
            DownloadRequest downloadRequest = objectMapper.readValue(body, DownloadRequest.class);

            out.add(downloadRequest);
        } else if (action.startsWith("GET /memory/stats")) {
            MemoryStatsRequest memoryStatsRequest = new MemoryStatsRequest();
            out.add(memoryStatsRequest);
        } else {
            FullHttpResponse fullHttpResponse = new DefaultFullHttpResponse(
                    HttpVersion.HTTP_1_1,
                    HttpResponseStatus.NOT_FOUND
            );
            ctx.writeAndFlush(fullHttpResponse);
            ctx.close();
        }
    }

}
