package net.tislib.downloaddelegator.data;

import lombok.Data;

import java.util.List;
import java.util.Map;
import java.util.UUID;

@Data
public class PageResponse {
    private UUID id;

    private byte[] content;

    private int httpStatus;
    private List<Map.Entry<String, String>> headers;

}
