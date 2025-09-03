package com.example.tour.model;

import jakarta.persistence.*;

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
    private String Difficulty;
    @Column(name="tags")
    private String tags;
    @Column(name="status",nullable = false)
    private Status status;
    @Column(name="price",nullable = false)
    private Double price;
    @Column(name="author_id",nullable = false)
    private Long authorId;

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

    public Tour(Integer id, String name, String description, String difficulty, String tags, Status status, Double price, Long authorId) {
        this.id = id;
        this.name = name;
        this.description = description;
        Difficulty = difficulty;
        this.tags = tags;
        this.status = status;
        this.price = price;
        this.authorId = authorId;
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
        return Difficulty;
    }

    public void setDifficulty(String difficulty) {
        Difficulty = difficulty;
    }

    public String getTags() {
        return tags;
    }

    public void setTags(String tags) {
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

    public Long getAuthorId() {
        return authorId;
    }

    public void setAuthorId(Long authorId) {
        this.authorId = authorId;
    }
}
