package net.tislib.downloaddelegator.data;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.Singular;

import java.net.URL;
import java.util.List;
import java.util.UUID;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
@JsonIgnoreProperties(ignoreUnknown = true)
public class PageUrl {

    private UUID id;

    private URL url;

    private String method;

    private Proxy proxy;

    private String bind;

    private int timeout;

    private int delay;

    private byte[] body;

    @Singular
    private List<Header> headers;

}
