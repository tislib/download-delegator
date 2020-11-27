FROM openjdk:11-jdk
ARG JAR_FILE=build/libs/*.jar

COPY ${JAR_FILE} app.jar

ENTRYPOINT ["java", "-XX:MaxDirectMemorySize=2000M", "-XX:MaxMetaspaceSize=91133K", "-XX:ReservedCodeCacheSize=240M", "-XX:-PrintGC", "-XX:-PrintGCDetails", "-Xss1M", "-Xmx1000M", "-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5125", "-Dio.netty.leakDetection.level=paranoid", "-Dio.netty.allocator.type=unpooled", "-jar", "/app.jar"]
