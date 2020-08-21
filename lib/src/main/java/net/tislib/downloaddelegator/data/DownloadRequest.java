package net.tislib.downloaddelegator.data;

import lombok.Data;

import java.util.List;

@Data
public class DownloadRequest {

    private List<PageUrl> urls;
}
