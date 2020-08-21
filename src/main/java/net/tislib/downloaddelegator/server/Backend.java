package net.tislib.downloaddelegator.server;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import lombok.SneakyThrows;
import net.tislib.downloaddelegator.base.Operational;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.downloader.Downloader;
import org.msgpack.jackson.dataformat.MessagePackFactory;

import java.io.IOException;
import java.net.InetSocketAddress;

public class Backend implements Operational, HttpHandler {

    private HttpServer server;
    private final Downloader downloader = new Downloader();

    private ObjectMapper objectMapper = new ObjectMapper(new MessagePackFactory());

    @SneakyThrows
    public void start() {
        server = HttpServer.create(new InetSocketAddress(8123), 0);
        server.createContext("/", this);
        server.setExecutor(null); // creates a default executor
        server.start();
    }

    public void stop() {
        server.stop(3000);
    }

    @Override
    public void handle(HttpExchange httpExchange) throws IOException {
        DownloadRequest downloadRequest = objectMapper.readValue(httpExchange.getRequestBody(), DownloadRequest.class);

        httpExchange.sendResponseHeaders(200, 0);
        downloader.download(downloadRequest, httpExchange.getResponseBody());

        httpExchange.getResponseBody().close();
    }
}
