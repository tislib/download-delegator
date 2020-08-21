package net.tislib.downloaddelegator;

import net.tislib.downloaddelegator.server.Backend;

public class Main {

    public static void main(String[] args) {
        Backend backend = new Backend();

        backend.start();
    }

}
