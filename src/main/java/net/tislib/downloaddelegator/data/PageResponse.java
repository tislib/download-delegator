package net.tislib.downloaddelegator.data;

import lombok.Data;

import java.util.UUID;

@Data
public class PageResponse {
    private UUID id;

    private byte[] content;

    private int httpStatus;
}
