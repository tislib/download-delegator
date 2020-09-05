package net.tislib.downloaddelegator.data;

import lombok.*;

import java.net.URL;
import java.util.List;
import java.util.UUID;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class PageUrl {

    private UUID id;

    private URL url;

    private String method;

    private Proxy proxy;

    private int timeout;

    private byte[] body;

    @Singular
    private List<Header> headers;
}
