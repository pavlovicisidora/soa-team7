package com.example.tour.service;

import com.example.tour.model.Tour;
import com.example.tour.repository.TourRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class TourService {
    @Autowired
    private TourRepository tourRepository;

    public Tour createTour(Tour Tour){
        return tourRepository.save(Tour);
    }
}
