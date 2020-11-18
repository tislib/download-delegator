package net.tislib.downloaddelegator.server.downloader;

import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageDecoder;
import lombok.extern.log4j.Log4j2;
import net.tislib.downloaddelegator.data.AtomicPageCounter;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.data.PageUrl;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

@Log4j2
public class PageUrlTaskSplitterHandler extends MessageToMessageDecoder<DownloadRequest> {

    @Override
    protected void decode(ChannelHandlerContext ctx, DownloadRequest downloadRequest, List<Object> out) {
        log.trace("starting download2: {}", downloadRequest);

        AtomicPageCounter atomicPageCounter = new AtomicPageCounter();

        List<PageUrl> urls = new ArrayList<>(downloadRequest.getUrls());

        AtomicInteger globalDelay = new AtomicInteger();

        urls.forEach(item -> {
            atomicPageCounter.markUnDone(item.getId());
            item.setPageCounter(atomicPageCounter);

            globalDelay.addAndGet(downloadRequest.getDelay());

            item.setDelay(item.getDelay() + globalDelay.get());
        });

        out.addAll(urls);
    }
}
