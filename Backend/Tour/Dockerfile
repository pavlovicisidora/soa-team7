FROM maven:3.9.0-eclipse-temurin-17 AS builder
WORKDIR /app

COPY Tour/pom.xml .
COPY Tour/src ./src
COPY Shopping/go.mod Shopping/go.sum ./Shopping/

RUN mvn clean package -DskipTests

FROM eclipse-temurin:17-jdk-alpine
WORKDIR /app
COPY --from=builder /app/target/*.jar /app/tour.jar
EXPOSE 8083
ENTRYPOINT ["java", "-jar", "/app/tour.jar"]