package net.tislib.downloaddelegator.netty;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.DecoderResult;
import io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.netty.handler.codec.http.FullHttpResponse;
import io.netty.handler.codec.http.HttpContent;
import io.netty.handler.codec.http.HttpHeaders;
import io.netty.handler.codec.http.HttpObject;
import io.netty.handler.codec.http.HttpRequest;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http.LastHttpContent;
import io.netty.handler.codec.http.QueryStringDecoder;
import io.netty.util.CharsetUtil;
import lombok.extern.log4j.Log4j2;

import static io.netty.handler.codec.http.HttpResponseStatus.*;
import static io.netty.handler.codec.http.HttpVersion.HTTP_1_1;

//@Log4j2
public class CustomHttpServerHandler extends SimpleChannelInboundHandler<Object> {

    private static final String CRLF = "\r\n";
    private final StringBuilder builder = new StringBuilder();

    @Override
    public void channelReadComplete(ChannelHandlerContext ctx) throws Exception {
        ctx.flush();
    }

    @Override
    protected void channelRead0(ChannelHandlerContext ctx, Object msg) throws Exception {
        HttpRequest request = null;

        if (msg instanceof HttpRequest) {
            request = (HttpRequest) msg;
            handleRequest(ctx, request);
        }

        if (msg instanceof HttpContent) {
            handleContent(ctx, request, (HttpContent) msg);
        }
    }

    private void handleRequest(ChannelHandlerContext ctx, HttpRequest request) {
        builder.setLength(0);
        builder.append("WELCOME TO THE WILD WILD WEB SERVER").append(CRLF);
        builder.append("===================================").append(CRLF);

        builder.append("VERSION: ").append(request.getProtocolVersion()).append(CRLF);
        builder.append("HOSTNAME: ").append(HttpHeaders.getHost(request, "unknown")).append(CRLF);
        builder.append("REQUEST_URI: ").append(request.getUri()).append(CRLF).append(CRLF);

        HttpHeaders headers = request.headers();
        headers.forEach(h -> {
            builder.append("HEADER: ")
                    .append(h.getKey())
                    .append(" = ")
                    .append(h.getValue())
                    .append(CRLF);
        });

        QueryStringDecoder queryStringDecoder = new QueryStringDecoder(request.getUri());
        queryStringDecoder.parameters().forEach((key, values) -> {
            builder.append("PARAM: ")
                    .append(key)
                    .append(" = ")
                    .append(values)
                    .append(CRLF);
        });

        appendDecoderResult(builder, request);
    }

    private void handleContent(ChannelHandlerContext ctx, HttpRequest request, HttpContent httpContent) {
        ByteBuf content = httpContent.content();

        if (content != null && content.isReadable()) {
            builder.append("CONTENT: ")
                    .append(content.toString(CharsetUtil.UTF_8))
                    .append(CRLF);
            appendDecoderResult(builder, request);
        }

        if (httpContent instanceof LastHttpContent) {
            builder.append("END OF CONTENT\r\n");

            LastHttpContent trailer = (LastHttpContent) httpContent;
            if (!trailer.trailingHeaders().isEmpty()) {
                builder.append(CRLF);

                trailer.trailingHeaders().names().forEach(name -> {
                    trailer.trailingHeaders().getAll(name).forEach(value -> {
                        builder.append("TRAINING HEADER: ")
                                .append(name)
                                .append(" = ")
                                .append(value)
                                .append(CRLF);
                    });
                });
                builder.append(CRLF);
            }
            writeResponse(ctx, request, trailer);

            ctx.writeAndFlush(Unpooled.EMPTY_BUFFER).addListener(ChannelFutureListener.CLOSE);
        }

    }

    private void writeResponse(ChannelHandlerContext ctx, HttpRequest req, HttpObject currentObj) {
        // Decide whether to close the connection or not.

        // Build the response object.
        HttpResponseStatus status = currentObj.getDecoderResult().isSuccess() ? OK : BAD_REQUEST;
        ByteBuf content = Unpooled.copiedBuffer(builder.toString(), CharsetUtil.UTF_8);

        FullHttpResponse response = new DefaultFullHttpResponse(HTTP_1_1, status, content);

        // Write the response
        ctx.write(response);

    }

    private static void appendDecoderResult(StringBuilder builder, HttpObject o) {
        DecoderResult dr = o.getDecoderResult();
        if (dr.isSuccess()) {
            return;
        }
        builder.append(".. WITH DECODER FAILURE: ");
        builder.append(dr.cause());
        builder.append(CRLF);
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        if (cause != null) {
//            log.error("ERROR:", cause);
        }
        ctx.close();
    }
}
