package net.tislib.downloaddelegator.config;

import lombok.experimental.UtilityClass;
import lombok.extern.log4j.Log4j2;
import org.apache.commons.lang3.StringUtils;

import static java.lang.Boolean.parseBoolean;

@UtilityClass
@Log4j2
public class ApplicationConfig {

    public String getConfig(Config config) {
        String systemProperty = System.getProperty(config.getName());
        String envProperty = System.getenv(config.getName());

        if (!StringUtils.isBlank(systemProperty)) {
            return systemProperty;
        } else if (!StringUtils.isBlank(envProperty)) {
            return systemProperty;
        } else {
            return config.getDefaultValue();
        }
    }

    public static boolean getBoolean(Config config) {
        return parseBoolean(getConfig(config));
    }
}
