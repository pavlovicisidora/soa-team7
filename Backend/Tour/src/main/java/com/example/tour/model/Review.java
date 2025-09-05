package com.example.tour.model;

import jakarta.persistence.*;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.List;

@Entity
public class Review {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;
    @Column(nullable = false)
    private Integer rating;
    private String comment;
    @Column(nullable = false)
    private String touristId;
    private LocalDate visitDate;
    private LocalDateTime createdDate;
    private List<String> images;
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name="tour_id",nullable = false)
    private Tour tour;

    public Review() {
    }

    public Review(Integer id, Integer rating, String comment, String touristId, LocalDate visitDate, LocalDateTime createdDate, List<String> images, Tour tour) {
        this.id = id;
        this.rating = rating;
        this.comment = comment;
        this.touristId = touristId;
        this.visitDate = visitDate;
        this.createdDate = createdDate;
        this.images = images;
        this.tour = tour;
    }

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public Integer getRating() {
        return rating;
    }

    public void setRating(Integer rating) {
        this.rating = rating;
    }

    public String getComment() {
        return comment;
    }

    public void setComment(String comment) {
        this.comment = comment;
    }

    public String getTouristId() {
        return touristId;
    }

    public void setTouristId(String touristId) {
        this.touristId = touristId;
    }

    public LocalDate getVisitDate() {
        return visitDate;
    }

    public void setVisitDate(LocalDate visitDate) {
        this.visitDate = visitDate;
    }

    public LocalDateTime getCreatedDate() {
        return createdDate;
    }

    public void setCreatedDate(LocalDateTime createdDate) {
        this.createdDate = createdDate;
    }

    public List<String> getImages() {
        return images;
    }

    public void setImages(List<String> images) {
        this.images = images;
    }

    public Tour getTour() {
        return tour;
    }

    public void setTour(Tour tour) {
        this.tour = tour;
    }
}
