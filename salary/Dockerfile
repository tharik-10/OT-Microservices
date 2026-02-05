# 1. Build Stage: Java 17
FROM maven:3.9.6-eclipse-temurin-17-alpine AS builder

LABEL maintainer="Opstree Solutions"
WORKDIR /java

COPY pom.xml .
RUN mvn dependency:go-offline

COPY src src
RUN mvn clean package -DskipTests

# 2. Runtime Stage: Java 17
FROM eclipse-temurin:17-jre-alpine

LABEL maintainer="Opstree Solutions"
WORKDIR /app

COPY --from=builder /java/target/salary-0.3.0-RELEASE.jar app.jar
COPY config.yaml /root/config/config.yaml

ENV CONFIG_FILE=/root/config/config.yaml

EXPOSE 8080
ENTRYPOINT ["java","-jar","/app/app.jar"]

