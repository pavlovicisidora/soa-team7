package com.example.tour.repository;

import com.example.tour.model.TourExecution;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface TourExecutionRepository extends JpaRepository<TourExecution, Integer> {
}

