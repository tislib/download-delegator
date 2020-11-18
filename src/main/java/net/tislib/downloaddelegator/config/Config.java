package net.tislib.downloaddelegator.config;

import lombok.Getter;
import lombok.RequiredArgsConstructor;

@Getter
@RequiredArgsConstructor
public enum Config {
    ADDR("server.addr", "0.0.0.0"),
    PORT("server.port", "8123"),
    TRACE_SERVER("logging.server.trace", "false"),
    TRACE_CLIENT("logging.client.trace", "false");

    private final String name;
    private final String defaultValue;

}
