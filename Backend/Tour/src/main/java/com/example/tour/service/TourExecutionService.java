package com.example.tour.service;

//import com.example.tour.grpc.Tour;
import com.example.tour.model.*;
import com.example.tour.repository.KeyPointRepository;
import com.example.tour.repository.TourExecutionRepository;
import com.example.tour.repository.TourRepository;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
//import java.util.Optional;

@Service
public class TourExecutionService {

    @Autowired
    private TourExecutionRepository tourExecutionRepository;
    @Autowired
    private TourRepository tourRepository;
    @Autowired
    private KeyPointRepository keyPointRepository;

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
        int KPNumber = keyPointRepository.findByTourId(tourId).size();
        newExecution.setKPtoBeCompleted(KPNumber);

        
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

    public TourExecution getTourExecution(Integer executionId) {
        return tourExecutionRepository.findById(executionId)
                .orElseThrow(() -> new RuntimeException("TourExecution not found"));
    }

    @Transactional
    public TourExecution completeKeyPoint(Integer executionId, int keyPointId) {
        TourExecution tourExecution = getTourExecution(executionId);
        List<ExecutedKeyPoint> executedList = tourExecution.getExecutedKeyPoints();

        if(keyPointId == 0){
            tourExecution.setLastActivity(LocalDateTime.now());
            return tourExecutionRepository.save(tourExecution);
        }
        
        for(ExecutedKeyPoint e : executedList){
            if(e.getKeypointId() == keyPointId)
                throw new RuntimeException("Keypoint already completed!");
        }

        ExecutedKeyPoint executed = new ExecutedKeyPoint(keyPointId, LocalDateTime.now());
        executedList.add(executed);

        tourExecution.setExecutedKeyPoints(executedList);
        tourExecution.setLastActivity(LocalDateTime.now());
        //Ako je odradio sve KT koje je trebao onda je i zavrsio turu uspesno
        if(executedList.size() == tourExecution.getKPtoBeCompleted()){

            tourExecution.setStatus(TourExecutionStatus.COMPLETED);
            tourExecution.setCompletionTime(LocalDateTime.now());

        }
            return tourExecutionRepository.save(tourExecution);

    }
}
