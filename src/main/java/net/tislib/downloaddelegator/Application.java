package net.tislib.downloaddelegator;

import net.tislib.downloaddelegator.server.Server;

public class Application {
    public static void main(String[] args) {
        Server server = new Server();

        server.run();
    }
}
