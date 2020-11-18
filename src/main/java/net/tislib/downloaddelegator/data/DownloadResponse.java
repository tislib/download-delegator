package net.tislib.downloaddelegator.data;

import lombok.Data;

import java.util.List;

@Data
public class DownloadResponse {

    private List<PageResponse> data;

}
