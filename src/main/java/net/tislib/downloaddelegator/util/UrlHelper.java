package net.tislib.downloaddelegator.util;

import lombok.experimental.UtilityClass;
import net.tislib.downloaddelegator.data.PageResponse;

@UtilityClass
public class UrlHelper {

    public static String makeFullUrl(String url, PageResponse pageResponse) {
        if (!url.startsWith("http")) {
            String baseUrl = pageResponse.getPageUrl().getUrl().getProtocol() +
                    "://" + pageResponse.getPageUrl().getUrl().getHost();

            if (!url.startsWith("/")) {
                url = "/" + url;
            }

            url = baseUrl + url;
        }

        return url;
    }
}
