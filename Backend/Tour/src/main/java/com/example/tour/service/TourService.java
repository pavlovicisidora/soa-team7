package com.example.tour.service;

import com.example.tour.model.KeyPoint;
import com.example.tour.model.Status;
import com.example.tour.model.Tour;
import com.example.tour.repository.TourRepository;
import com.example.tour.repository.KeyPointRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.NoSuchElementException;

@Service
public class TourService {
    @Autowired
    private TourRepository tourRepository;

    @Autowired
    private KeyPointRepository keypointRepository;

    public Tour createTour(Tour tour){
        return tourRepository.save(tour);
    }
    public List<Tour> findAllToursById(String id){return tourRepository.findToursForAuthor(id);}
    public Tour findById(Integer id) {
        return tourRepository.findById(id)
            .orElseThrow(() -> new NoSuchElementException("Tour with id " + id + " not found"));
    }


    @Transactional
    public Tour updateTour(Tour tourWithUpdates) {
        if (tourWithUpdates.getId() == null) {
            throw new IllegalArgumentException("ID must be provided for an update.");
        }

        Tour existingTour = findById(tourWithUpdates.getId());

        existingTour.setName(tourWithUpdates.getName());
        existingTour.setDescription(tourWithUpdates.getDescription());
        existingTour.setDifficulty(tourWithUpdates.getDifficulty());
        existingTour.setTags(tourWithUpdates.getTags());
        existingTour.setPrice(tourWithUpdates.getPrice());
        existingTour.setDistanceInKm(tourWithUpdates.getDistanceInKm());

        Status newStatus = tourWithUpdates.getStatus();


        if (newStatus == Status.PUBLISHED) {
            validateForPublishing(existingTour);
            existingTour.setStatus(newStatus);
            existingTour.setPublishedDateTime(LocalDateTime.now());
            existingTour.setArchivedDateTime(null);
        } else {
            existingTour.setStatus(newStatus);
            if (newStatus == Status.ARCHIVED) {

                existingTour.setPublishedDateTime(null);
                existingTour.setArchivedDateTime(LocalDateTime.now());
            }
            else{
                existingTour.setPublishedDateTime(null);
                existingTour.setArchivedDateTime(null);
            }
        }


        return tourRepository.save(existingTour);
    }


    private void validateForPublishing(Tour tour) {

        if (tour.getName() == null || tour.getName().isBlank() ||
                tour.getDescription() == null || tour.getDescription().isBlank() ||
                tour.getDifficulty() == null || tour.getDifficulty().isBlank() ||
                tour.getTags() == null || tour.getTags().isEmpty()) {
            throw new IllegalStateException("Name, description, difficulty, and tags must be provided to publish.");
        }


        List<KeyPoint> keypoints = keypointRepository.findByTourId(tour.getId());
        if (keypoints.size() < 2) {
            throw new IllegalStateException("At least two keypoints must be defined to publish. Found: " + keypoints.size());
        }



        if (tour.getPrice() == null || tour.getPrice() <= 0) {
            throw new IllegalStateException("A valid price greater than zero must be set to publish.");
        }
    }


    public boolean deleteTour(Integer id){
        if (tourRepository.existsById(id)) {
            tourRepository.deleteById(id);
            return true;
        } else {
            return false;
        }
    }
    public List<Tour> findAllTours(){return tourRepository.findAll();}
}
