package net.tislib.downloaddelegator.server;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.MessageToMessageCodec;
import io.netty.handler.codec.MessageToMessageDecoder;
import io.netty.handler.codec.http.FullHttpRequest;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.data.DownloadRequest;

import java.nio.charset.Charset;
import java.util.List;

@Log4j2
public class OperationSwitchHandler extends MessageToMessageDecoder<FullHttpRequest> {

    private static final ObjectMapper objectMapper = new ObjectMapper();

    @Override
    protected void decode(ChannelHandlerContext ctx, FullHttpRequest request, List<Object> out) throws Exception {
        String action = request.method().name() + " " + request.uri();
        String body = request.content().toString(Charset.defaultCharset());

        log.debug("Handling request: {} {}", action, ctx.channel().remoteAddress());

        if (action.startsWith("POST /download.tar.gz")) {
            DownloadRequest downloadRequest = objectMapper.readValue(body, DownloadRequest.class);

            out.add(downloadRequest);
        }
    }

}
