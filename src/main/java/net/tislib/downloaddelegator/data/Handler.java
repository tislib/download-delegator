package net.tislib.downloaddelegator.data;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class Handler {

    private String name;

    private int connectTimeout;

    private int socketTimeout;

    private int delay;

}
