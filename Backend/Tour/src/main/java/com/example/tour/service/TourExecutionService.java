package com.example.tour.service;

import com.example.tour.model.TourExecution;
import com.example.tour.model.TourExecutionStatus;
import com.example.tour.repository.TourExecutionRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import java.time.LocalDateTime;

@Service
public class TourExecutionService {

    @Autowired
    private TourExecutionRepository tourExecutionRepository;

    public TourExecution startTour(Integer tourId, String touristId) {
        TourExecution newExecution = new TourExecution();
        newExecution.setTourId(tourId);
        newExecution.setTouristId(touristId);
        newExecution.setStatus(TourExecutionStatus.IN_PROGRESS);
        newExecution.setStartTime(LocalDateTime.now());
        newExecution.setLastActivity(LocalDateTime.now());
        
        return tourExecutionRepository.save(newExecution);
    }
    
    public TourExecution abandonTour(Integer executionId) {
        TourExecution execution = tourExecutionRepository.findById(executionId)
                .orElseThrow(() -> new RuntimeException("TourExecution with id " + executionId + " not found"));
        
        execution.setStatus(TourExecutionStatus.ABANDONED);
        execution.setCompletionTime(LocalDateTime.now());
        execution.setLastActivity(LocalDateTime.now());
        
        return tourExecutionRepository.save(execution);
    }

    public TourExecution completeTour(Integer executionId) {
        TourExecution execution = tourExecutionRepository.findById(executionId)
                .orElseThrow(() -> new RuntimeException("TourExecution with id " + executionId + " not found"));
        
        execution.setStatus(TourExecutionStatus.COMPLETED);
        execution.setCompletionTime(LocalDateTime.now());
        execution.setLastActivity(LocalDateTime.now());
        
        return tourExecutionRepository.save(execution);
    }
}
