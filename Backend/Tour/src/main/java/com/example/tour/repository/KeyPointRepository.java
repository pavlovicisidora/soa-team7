package com.example.tour.repository;

import com.example.tour.model.KeyPoint;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;

public interface KeyPointRepository extends JpaRepository<KeyPoint, Integer> {
    List<KeyPoint> findByTourId(@Param("id") Integer tourId);
}
