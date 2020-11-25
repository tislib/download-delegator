package net.tislib.downloaddelegator.server.downloader;

import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.base.TimeCalc;
import net.tislib.downloaddelegator.client.DownloadClient;
import net.tislib.downloaddelegator.data.PageResponse;
import net.tislib.downloaddelegator.data.PageUrl;

@Log4j2
public class PageDownloadHandler extends SimpleChannelInboundHandler<PageUrl> {

    private static final TimeCalc timeCalc = new TimeCalc();

    @Override
    protected void channelRead0(ChannelHandlerContext ctx, PageUrl pageUrl) {
        log.trace("downloading page: {} {} {}", pageUrl.getUrl(), pageUrl.getId(), ctx.channel().localAddress());

        startDownload(ctx, pageUrl);
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        log.error(cause);
        ctx.close();
    }

    private void startDownload(ChannelHandlerContext ctx, PageUrl pageUrl) {
        DownloadClient downloadClient = new DownloadClient() {
            @Override
            public void onFinish(PageResponse pageResponse) {
                onDownload(pageResponse, pageUrl, ctx);
            }
        };

        ChannelFuture channelFuture = downloadClient.connect(pageUrl);

        downloadClient.download(channelFuture, pageUrl.getUrl());
    }

    private void onDownload(PageResponse pageResponse, PageUrl pageUrl, ChannelHandlerContext ctx) {
        pageResponse.setId(pageUrl.getId());
        pageResponse.setPageUrl(pageUrl);

        log.trace("response page for: {} {}", pageUrl.getUrl(), pageUrl.getId());

        ctx.writeAndFlush(pageResponse);

        timeCalc.printSpeedStep();
    }
}
