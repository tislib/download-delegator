FROM openjdk:11-jdk
ARG JAR_FILE=build/libs/*.jar

COPY ${JAR_FILE} app.jar

ENTRYPOINT ["java", "-XX:MaxDirectMemorySize=2000M", "-XX:MaxMetaspaceSize=91133K", "-Dio.netty.leakDetection.level=paranoid", "-XX:ReservedCodeCacheSize=240M", "-XX:-PrintGC", "-XX:-PrintGCDetails", "-Xss1M", "-Xmx1000M", "-jar", "/app.jar"]
