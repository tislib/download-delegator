package net.tislib.downloaddelegator.data;

import lombok.Data;

import java.util.Map;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

@Data
public class AtomicPageCounter {
    private Map<UUID, Boolean> pageDoneMap = new ConcurrentHashMap<>();

    public void markDone(UUID uuid) {
        pageDoneMap.put(uuid, true);
    }

    public void markUnDone(UUID uuid) {
        pageDoneMap.put(uuid, false);
    }

    public boolean isAllDone() {
        return pageDoneMap.values()
                .stream()
                .allMatch(item -> item);
    }

    public boolean isNoneDone() {
        return pageDoneMap.values()
                .stream()
                .noneMatch(item -> item);
    }
}
