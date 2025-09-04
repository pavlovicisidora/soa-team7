package com.example.tour.controller;

import com.example.tour.model.Tour;
import com.example.tour.service.TourService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/tour",produces = MediaType.APPLICATION_JSON_VALUE)
public class TourController {
    @Autowired
    private TourService tourService;
    @PostMapping("/create")
    public ResponseEntity<Tour> createTour(@RequestBody Tour tour){
        Tour tourData = this.tourService.createTour(tour);
        return ResponseEntity.ok(tourData);
    }
}
