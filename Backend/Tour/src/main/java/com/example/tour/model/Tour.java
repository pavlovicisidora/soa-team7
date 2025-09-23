package com.example.tour.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.persistence.*;

import java.time.LocalDateTime;
import java.util.List;

@Entity
public class Tour {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;
    @Column(name="name",nullable = false)
    private String name;
    @Column(name="description")
    private String description;
    @Column(name="Difficulty",nullable = false)
    private String difficulty;
    @Column(name="tags")
    private List<String> tags;
    @Column(name="status",nullable = false)
    private Status status;
    @Column(name="price",nullable = false)
    private Double price;
    @Column(name="author_id",nullable = false)
    private String authorId;
    
    @JsonProperty("distance_km")
    @Column(name = "distance_km")
    private Double distanceInKm;
    @Column(name = "published_date_time")
    private LocalDateTime publishedDateTime;
    @Column(name = "archived_date_time")
    private LocalDateTime archivedDateTime;

    @PrePersist
    public void setInitialValues(){
        if(status==null){
            status=Status.DRAFT;
        }
        if(price==null){
            price=0.0;
        }
    }

    public Tour() {
    }

    public Tour(Integer id, String name, String description, String difficulty, List<String> tags, Status status, Double price, String authorId, double distanceInKm, LocalDateTime publishedDateTime, LocalDateTime archivedDateTime) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.difficulty = difficulty;
        this.tags = tags;
        this.status = status;
        this.price = price;
        this.authorId = authorId;
        this.distanceInKm = distanceInKm;
        this.publishedDateTime = publishedDateTime;
        this.archivedDateTime = archivedDateTime;
    }

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getDifficulty() {
        return difficulty;
    }

    public void setDifficulty(String difficulty) {
        this.difficulty = difficulty;
    }

    public List<String> getTags() {
        return tags;
    }

    public void setTags(List<String> tags) {
        this.tags = tags;
    }

    public Status getStatus() {
        return status;
    }

    public void setStatus(Status status) {
        this.status = status;
    }

    public Double getPrice() {
        return price;
    }

    public void setPrice(Double price) {
        this.price = price;
    }

    public String getAuthorId() {
        return authorId;
    }

    public void setAuthorId(String authorId) {
        this.authorId = authorId;
    }

    public LocalDateTime getPublishedDateTime() { return publishedDateTime; }

    public void setPublishedDateTime(LocalDateTime publishedDateTime) { this.publishedDateTime = publishedDateTime; }
    
    public LocalDateTime getArchivedDateTime() { return archivedDateTime; }
    
    public void setArchivedDateTime(LocalDateTime archivedDateTime) { this.archivedDateTime = archivedDateTime; }

    public Double getDistanceInKm() {
        return distanceInKm;
    }

    public void setDistanceInKm(Double distanceInKm) {
        this.distanceInKm = distanceInKm;
    }
}
