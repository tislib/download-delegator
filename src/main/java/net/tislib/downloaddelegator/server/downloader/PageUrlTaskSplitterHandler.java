package net.tislib.downloaddelegator.server.downloader;

import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageDecoder;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.data.AtomicPageCounter;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageUrl;

import java.util.ArrayList;
import java.util.List;

@Log4j2
public class PageUrlTaskSplitterHandler extends MessageToMessageDecoder<DownloadRequest> {

    @Override
    protected void decode(ChannelHandlerContext ctx, DownloadRequest downloadRequest, List<Object> out) {
        log.trace("starting download2: {}", downloadRequest);

        AtomicPageCounter atomicPageCounter = new AtomicPageCounter();

        List<PageUrl> urls = new ArrayList<>(downloadRequest.getUrls());
        urls.forEach(item -> {
            atomicPageCounter.markUnDone(item.getId());
            item.setPageCounter(atomicPageCounter);
        });

        urls.forEach(url -> {
            if (downloadRequest.getDelay() > 0) {
                out.add(url);
            } else {
                out.add(url);
            }
        });
    }
}
