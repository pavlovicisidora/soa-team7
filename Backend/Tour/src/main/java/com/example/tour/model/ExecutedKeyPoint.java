package com.example.tour.model;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;

import java.time.LocalDateTime;

@Embeddable
public class ExecutedKeyPoint {
    private Integer keypointId;

    @Column(nullable = false)
    private LocalDateTime completedAt;

    public ExecutedKeyPoint() {
    }

    public ExecutedKeyPoint(Integer keypointId, LocalDateTime completedAt) {
        this.keypointId = keypointId;
        this.completedAt = completedAt;
    }

    public Integer getKeypointId() {
        return keypointId;
    }

    public void setKeypointId(Integer keypointId) {
        this.keypointId = keypointId;
    }

    public LocalDateTime getCompletedAt() {
        return completedAt;
    }

    public void setCompletedAt(LocalDateTime completedAt) {
        this.completedAt = completedAt;
    }
}
