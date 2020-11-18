package net.tislib.downloaddelegator.test.base;

import io.netty.buffer.ByteBuf;
import lombok.Data;

import java.util.List;
import java.util.Map;
import java.util.UUID;

@Data
public class PageData {
    private UUID id;

    private byte[] content;

    private int httpStatus;
    private List<Map.Entry<String, String>> headers;

}
