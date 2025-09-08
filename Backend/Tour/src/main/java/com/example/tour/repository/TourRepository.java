package com.example.tour.repository;

import com.example.tour.model.Tour;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;

public interface TourRepository extends JpaRepository<Tour, Integer> {
    @Query("SELECT t FROM Tour t WHERE t.authorId=:authorId")
    List<Tour> findToursForAuthor(@Param("authorId") String authorId);
}
