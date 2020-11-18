package net.tislib.downloaddelegator.base;

import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import lombok.experimental.UtilityClass;

@UtilityClass
public class EventLoopGroups {
    public static final EventLoopGroup serverGroup = new NioEventLoopGroup();
    public static final EventLoopGroup clientGroup = new NioEventLoopGroup();
}
