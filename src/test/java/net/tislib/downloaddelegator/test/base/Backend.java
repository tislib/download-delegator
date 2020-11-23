package net.tislib.downloaddelegator.test.base;

import com.fasterxml.jackson.databind.ObjectMapper;
import kong.unirest.HttpResponse;
import kong.unirest.Unirest;
import lombok.SneakyThrows;
import net.tislib.downloaddelegator.data.DownloadRequest;
import net.tislib.downloaddelegator.server.Server;
import org.junit.platform.commons.util.StringUtils;
import org.junit.rules.ExternalResource;

import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLConnection;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;
import java.util.function.Consumer;
import java.util.function.Supplier;

public class Backend {

    private final Server server = new Server();
    private final ObjectMapper objectMapper = new ObjectMapper();

    public Backend() {
        server.run();
    }

    @SneakyThrows
    public void callAsync(DownloadRequest downloadRequest, Consumer<List<PageData>> consumer) {
        new Thread(() ->
                consumer.accept(call(downloadRequest)))
                .start();
    }

    @SneakyThrows
    public List<PageData> call(DownloadRequest downloadRequest) {
        System.out.println("CALLING BACKEND");
        URL url = new URL("http://127.0.0.1:8123/download.tar.gz");

        HttpURLConnection conn = (HttpURLConnection) url.openConnection();

        conn.setRequestMethod("POST");
        conn.setDoOutput(true);
        conn.setDoInput(true);

        objectMapper.writeValue(conn.getOutputStream(), downloadRequest);

        String body = new String(conn.getInputStream().readAllBytes());

        List<PageData> response = new ArrayList<>();

        StringBuilder content = new StringBuilder();
        PageData currentResponse = null;

        for (String line : body.split("\\n")) {
            System.out.println(line);
            if (StringUtils.isBlank(line)) {
                continue;
            }
            if (currentResponse == null) {
                UUID id = UUID.fromString(line);
                currentResponse = new PageData();
                currentResponse.setId(id);
                content.setLength(0);
                continue;
            }

            if (currentResponse.getId().toString().equals(line)) {
                currentResponse.setContent(content.toString().getBytes());
                response.add(currentResponse);
                currentResponse = null;
                continue;
            }

            content.append(line);
        }

        if (currentResponse != null) {
            throw new RuntimeException("response is not complete");
        }


        return response;
    }
}
