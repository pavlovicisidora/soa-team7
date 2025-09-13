package com.example.tour.model;

import jakarta.persistence.*;
import java.time.LocalDateTime;

@Entity
public class TourExecution {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Integer id;

    @Column(nullable = false)
    private Integer tourId;

    @Column(nullable = false)
    private String touristId;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private TourExecutionStatus status;

    private LocalDateTime completionTime;
    @Column(nullable = false)
    private LocalDateTime startTime;
    
    @Column(nullable = false)
    private LocalDateTime lastActivity;

    public TourExecution() {
    }

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public Integer getTourId() {
        return tourId;
    }

    public void setTourId(Integer tourId) {
        this.tourId = tourId;
    }

    public String getTouristId() {
        return touristId;
    }

    public void setTouristId(String touristId) {
        this.touristId = touristId;
    }

    public TourExecutionStatus getStatus() {
        return status;
    }

    public void setStatus(TourExecutionStatus status) {
        this.status = status;
    }

    public LocalDateTime getCompletionTime() {
        return completionTime;
    }

    public void setCompletionTime(LocalDateTime completionTime) {
        this.completionTime = completionTime;
    }
    
    public LocalDateTime getStartTime() { return startTime; }
    public void setStartTime(LocalDateTime startTime) { this.startTime = startTime; }
    public LocalDateTime getLastActivity() { return lastActivity; }
    public void setLastActivity(LocalDateTime lastActivity) { this.lastActivity = lastActivity; }
}
