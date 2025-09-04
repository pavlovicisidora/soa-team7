package com.example.tour.service;

import com.example.tour.model.Tour;
import com.example.tour.repository.TourRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class TourService {
    @Autowired
    private TourRepository tourRepository;

    public Tour createTour(Tour tour){
        return tourRepository.save(tour);
    }
    public List<Tour> findAllToursById(String id){return tourRepository.findAllToursById(id);}
}
