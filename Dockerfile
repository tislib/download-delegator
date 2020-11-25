FROM openjdk:11-jdk
ARG JAR_FILE=build/libs/*.jar

COPY ${JAR_FILE} app.jar

ENTRYPOINT ["java", "-XX:MaxDirectMemorySize=300M", "-XX:MaxMetaspaceSize=91133K", "-XX:ReservedCodeCacheSize=240M", "-Xss1M", "-Xmx300M", "-jar", "/app.jar"]
