package net.tislib.downloaddelegator.server.downloader;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.*;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.base.TimeCalc;
import net.tislib.downloaddelegator.client.DownloadClient;
import net.tislib.downloaddelegator.data.PageResponse;
import net.tislib.downloaddelegator.data.PageUrl;

import java.util.concurrent.TimeUnit;

@Log4j2
public class PageDownloadHandler extends SimpleChannelInboundHandler<PageUrl> {

    private final TimeCalc timeCalc = new TimeCalc();

    @Override
    protected void channelRead0(ChannelHandlerContext ctx, PageUrl pageUrl) {
        log.trace("downloading page: {}", pageUrl.getUrl());
        if (pageUrl.getDelay() == 0) {
            startDownload(ctx, pageUrl);
        } else {
            ctx.executor().schedule(() -> startDownload(ctx, pageUrl), pageUrl.getDelay(), TimeUnit.MILLISECONDS);
        }
    }

    private void startDownload(ChannelHandlerContext ctx, PageUrl pageUrl) {
        DownloadClient downloadClient = new DownloadClient() {
            @Override
            public void onFullResponse(PageResponse response) {
                ctx.executor().execute(() -> onDownload(pageUrl, ctx, response));
            }
        };

        ChannelFuture channelFuture = downloadClient.connect(pageUrl.getUrl());

        downloadClient.download(channelFuture, pageUrl.getUrl());
    }

    private void onDownload(PageUrl pageUrl,
                            ChannelHandlerContext ctx,
                            PageResponse pageResponse) {
        pageResponse.setId(pageUrl.getId());
        log.trace("response page for: {}", pageUrl.getUrl());

        timeCalc.printSpeedStep();

        // send headers if is first page response
        if (pageUrl.getPageCounter().isNoneDone()) {
            startResponse(ctx);
        }

        pageUrl.getPageCounter().markDone(pageResponse.getId());

        sendPageMetaHead(pageUrl, ctx);

        DefaultHttpContent defaultHttpContent = new DefaultHttpContent(pageResponse.getContent());
        ctx.writeAndFlush(defaultHttpContent);

        sendPageMetaTail(pageUrl, ctx);

        if (pageUrl.getPageCounter().isAllDone()) {
            finishResponse(pageUrl, ctx);
        }
    }

    private void finishResponse(PageUrl pageUrl, ChannelHandlerContext ctx) {
        DefaultLastHttpContent defaultLastHttpContent = new DefaultLastHttpContent();
        ctx.writeAndFlush(defaultLastHttpContent);
        ctx.close();

        log.trace("last response finish page for: {}", pageUrl.getUrl());
    }

    private void startResponse(ChannelHandlerContext ctx) {
        DefaultHttpResponse defaultHttpResponse = new DefaultHttpResponse(HttpVersion.HTTP_1_1, HttpResponseStatus.OK);

        ctx.writeAndFlush(defaultHttpResponse);
    }

    private void sendPageMetaHead(PageUrl pageUrl, ChannelHandlerContext ctx) {
        ByteBuf head = ctx.alloc().buffer();
        head.writeBytes(pageUrl.getId().toString().getBytes());
        head.writeBytes("\n".getBytes());

        // write page beginning splitter
        ctx.write(new DefaultHttpContent(head));
    }

    private void sendPageMetaTail(PageUrl pageUrl, ChannelHandlerContext ctx) {
        ByteBuf tail = ctx.alloc().buffer();
        tail.writeBytes("\n".getBytes());

        tail.writeBytes(pageUrl.getId().toString().getBytes());

        if (!pageUrl.getPageCounter().isAllDone()) {
            tail.writeBytes("\n".getBytes()); // if is not last item, add new line after tail
        }

        // write page ending splitter
        ctx.writeAndFlush(new DefaultHttpContent(tail));
    }
}