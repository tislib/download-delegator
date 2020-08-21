package net.tislib.downloaddelegator.base;

import lombok.experimental.UtilityClass;

import java.io.IOException;
import java.util.Properties;

@UtilityClass
public class Config {

    private static final Properties properties;

    public static final String CONFIG_PROPERTIES_FILE = "config.properties";

    static {
        properties = new Properties();
        try {
            properties.load(Thread.currentThread().getContextClassLoader().getResourceAsStream(CONFIG_PROPERTIES_FILE));
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    public int maxDownloadQueueSize() {
        return Integer.parseInt(properties.getProperty("downloader.maxQueueSize", "1000"));
    }

    public static int downloaderThreadCount() {
        return Integer.parseInt(properties.getProperty("downloader.threadCount", "50"));
    }
}
