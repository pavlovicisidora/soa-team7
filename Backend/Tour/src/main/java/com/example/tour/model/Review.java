package com.example.tour.model;

import jakarta.persistence.*;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.List;

@Entity
public class Review {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    Integer id;
    @Column(nullable = false)
    Integer rating;
    @Column(nullable = false)
    private String touristId;
    private LocalDate visitDate;
    private LocalDateTime createdDate;
    private List<String> images;
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name="tour_id",nullable = false)
    private Tour tour;
}
