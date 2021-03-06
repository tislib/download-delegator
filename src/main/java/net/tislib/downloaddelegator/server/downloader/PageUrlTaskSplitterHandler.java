package net.tislib.downloaddelegator.server.downloader;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelDuplexHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelPromise;
import io.netty.handler.codec.http.DefaultHttpContent;
import io.netty.handler.codec.http.DefaultHttpResponse;
import io.netty.handler.codec.http.DefaultLastHttpContent;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http.HttpVersion;
import io.netty.util.ReferenceCountUtil;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageResponse;
import net.tislib.downloaddelegator.data.PageUrl;
import net.tislib.downloaddelegator.util.UrlHelper;

import java.net.URL;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

@Log4j2
public class PageUrlTaskSplitterHandler extends ChannelDuplexHandler {

    private final Set<UUID> pageUrlSet = new HashSet<>();

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        if (msg instanceof DownloadRequest) {
            startResponse(ctx);
            processRequest(ctx, (DownloadRequest) msg);
        } else {
            super.channelRead(ctx, msg);
        }

        ctx.executor().schedule(() -> {
            if (ctx.channel().isOpen()) {
                log.warn("Closing connection for no response");
                ctx.channel().close();
            }
        }, 2, TimeUnit.MINUTES);
    }

    @Override
    public void write(ChannelHandlerContext ctx, Object msg, ChannelPromise promise) throws Exception {
        if (msg instanceof PageResponse) {
            processResponse(ctx, (PageResponse) msg, promise);
        } else {
            super.write(ctx, msg, promise);
        }
    }

    private void processResponse(ChannelHandlerContext ctx,
                                 PageResponse pageResponse,
                                 ChannelPromise promise) throws Exception {
        try {
            if (String.valueOf(pageResponse.getHttpStatus()).startsWith("3")) {
                for (Map.Entry<String, String> header : pageResponse.getHeaders()) {
                    if (header.getKey().equalsIgnoreCase("location")) {
                        String fullUrl = UrlHelper.makeFullUrl(header.getValue(), pageResponse);

                        pageResponse.getPageUrl().setUrl(new URL(fullUrl));

                        super.channelRead(ctx, pageResponse.getPageUrl());
                        return;
                    }
                }
            }

            pageUrlSet.remove(pageResponse.getPageUrl().getId());
            sendPageMetaHead(pageResponse.getPageUrl(), ctx);

            ReferenceCountUtil.touch(pageResponse.getContent(), pageResponse.getId());

            if (pageResponse.getContent() != null) {
                DefaultHttpContent defaultHttpContent = new DefaultHttpContent(pageResponse.getContent());
                ctx.writeAndFlush(defaultHttpContent);
            }

            sendPageMetaTail(pageResponse.getPageUrl(), ctx);

            log.trace("page downloaded: {} {}", pageResponse.getPageUrl().getUrl(), pageResponse.getPageUrl().getId());
        } catch (Exception e) {
            e.printStackTrace();
            pageUrlSet.remove(pageResponse.getPageUrl().getId());
        }

        if (pageUrlSet.size() == 0) {
            log.debug("last response finish page for: {}", pageResponse.getPageUrl().getUrl());
            finishResponse(ctx);
        }
    }

    private void processRequest(ChannelHandlerContext ctx, DownloadRequest downloadRequest) throws Exception {
        log.trace("starting download2: {}", downloadRequest);
        log.debug("starting download page count: {}", downloadRequest.getUrls().size());

        int globalDelay = 0;

        for (PageUrl pageUrl : downloadRequest.getUrls()) {
            pageUrlSet.add(pageUrl.getId());

            globalDelay += downloadRequest.getDelay();

            int localDelay = globalDelay + pageUrl.getDelay();

            if (localDelay == 0) {
                super.channelRead(ctx, pageUrl);
            } else {
                ctx.executor().schedule(() -> forward(ctx, pageUrl), localDelay, TimeUnit.MILLISECONDS);
            }
        }

    }

    private void forward(ChannelHandlerContext ctx, PageUrl pageUrl) {
        try {
            super.channelRead(ctx, pageUrl);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private void startResponse(ChannelHandlerContext ctx) {
        DefaultHttpResponse defaultHttpResponse = new DefaultHttpResponse(HttpVersion.HTTP_1_1, HttpResponseStatus.OK);

        ctx.writeAndFlush(defaultHttpResponse);
    }

    private void sendPageMetaHead(PageUrl pageUrl, ChannelHandlerContext ctx) {
        ByteBuf head = ctx.alloc().buffer();

        ReferenceCountUtil.touch(head);

        head.writeBytes(pageUrl.getId().toString().getBytes());
        head.writeBytes("\n".getBytes());

        // write page beginning splitter
        ctx.write(new DefaultHttpContent(head));
    }

    private void sendPageMetaTail(PageUrl pageUrl, ChannelHandlerContext ctx) {
        ByteBuf tail = ctx.alloc().buffer();

        ReferenceCountUtil.touch(tail);

        tail.writeBytes("\n".getBytes());

        tail.writeBytes(pageUrl.getId().toString().getBytes());

        tail.writeBytes("\n".getBytes()); // if is not last item, add new line after tail

        // write page ending splitter
        ctx.writeAndFlush(new DefaultHttpContent(tail));
    }

    private void finishResponse(ChannelHandlerContext ctx) {
        DefaultLastHttpContent defaultLastHttpContent = new DefaultLastHttpContent();
        ctx.writeAndFlush(defaultLastHttpContent);
        ctx.close();
    }
}
