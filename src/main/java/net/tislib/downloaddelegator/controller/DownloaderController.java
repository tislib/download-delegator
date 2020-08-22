package net.tislib.downloaddelegator.controller;

import lombok.RequiredArgsConstructor;
import net.tislib.downloaddelegator.service.DownloaderService;
import net.tislib.downloaddelegator.data.DownloadRequest;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;

@RestController
@RequestMapping("/")
@RequiredArgsConstructor
public class DownloaderController {

    private final DownloaderService downloaderService;

    @PostMapping
    public Flux<DataBuffer> download(@RequestBody DownloadRequest downloadRequest) {
        return downloaderService.download(downloadRequest);
    }

}
