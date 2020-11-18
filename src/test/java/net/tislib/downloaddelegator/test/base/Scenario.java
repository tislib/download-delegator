package net.tislib.downloaddelegator.test.base;

import lombok.Builder;
import lombok.Data;
import lombok.Singular;

import java.util.ArrayList;
import java.util.List;

@Data
@Builder
public class Scenario {

    @Singular
    private final List<Request> requests;

    @Data
    @Builder
    public static class Request {
        private int count;
        private byte[] responseData;
        private int statusCode;
        private int responseTime;
        private boolean closeConnectionWithoutResponse; //@todo improve fail scenarios
        private int protocol; // 0-> http, 1 -> ssl, 2 -> http2 (convert to enum)
    }
}
