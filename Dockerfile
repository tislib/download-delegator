FROM openjdk:11-jdk
ARG JAR_FILE=build/libs/*.jar

COPY ${JAR_FILE} app.jar

ENTRYPOINT ["java", "-XX:MaxDirectMemorySize=300M", "-XX:MaxMetaspaceSize=91133K", "-XX:ReservedCodeCacheSize=240M", "-Dreactor.netty.tcp.sshHandshakeTimeout=120000", "-Xss1M", "-Xmx300M", "-Dreactor.netty.tcp.sshHandshakeTimeout=500000", "-jar", "/app.jar"]
