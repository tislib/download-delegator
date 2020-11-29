package net.tislib.downloaddelegator.data;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import lombok.Data;

import java.util.List;

@Data
@JsonIgnoreProperties(ignoreUnknown = true)
public class DownloadRequest {

    private List<PageUrl> urls;

    private int delay;

    private int timeout;


}
