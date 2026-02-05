# 1. Build Stage: Use JDK 21
FROM maven:3.9.6-eclipse-temurin-21-alpine as builder

LABEL maintainer="Opstree Solutions"
WORKDIR /java/
COPY pom.xml /java/
COPY src /java/src/
# Using -DskipTests to speed up the container build since Jenkins runs tests earlier
RUN mvn clean package -DskipTests

# 2. Runtime Stage: Use JRE 21
FROM alpine:latest
LABEL maintainer="Opstree Solutions"
USER root

# Install OpenJDK 21 instead of 17
RUN apk update && \
    apk add openjdk21-jre

COPY --from=builder /java/target/salary-0.3.0-RELEASE.jar /app/salary.jar
COPY ./config.yaml /root/config/config.yaml
ENV CONFIG_FILE "/root/config/config.yaml"

EXPOSE 8080
# Use the direct path to the java 21 binary
ENTRYPOINT ["java", "-jar", "/app/salary.jar"]
