package com.example.tour.service;

import com.example.tour.grpc.Tour;
import com.example.tour.model.Status;
import com.example.tour.model.TourExecution;
import com.example.tour.model.TourExecutionStatus;
import com.example.tour.repository.TourExecutionRepository;
import com.example.tour.repository.TourRepository;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import java.time.LocalDateTime;

@Service
public class TourExecutionService {

    @Autowired
    private TourExecutionRepository tourExecutionRepository;
    @Autowired
    private TourRepository tourRepository;

    public TourExecution startTour(Integer tourId, String touristId) {
        com.example.tour.model.Tour tourToStart = tourRepository.findById(tourId)
                .orElseThrow(() -> new RuntimeException("Tour with id " + tourId + " not found"));
        
        if (tourToStart.getStatus() == Status.DRAFT) {
            throw new IllegalStateException("Tour cannot be started because its status is not PUBLISHED or ARCHIVED. Current status: " + tourToStart.getStatus());
        }
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

    public TourExecution getTourExecution(Integer executionId) {
        return tourExecutionRepository.findById(executionId)
                .orElseThrow(() -> new RuntimeException("TourExecution not found"));
    }
}
