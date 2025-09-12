package com.example.tour.service;

import com.example.tour.model.Review;
import com.example.tour.repository.ReviewRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class ReviewService {
    @Autowired
    ReviewRepository reviewRepository;
    public Review createReview(Review review) {return reviewRepository.save(review);}
    public List<Review> getAllReviewForTour(Integer tourId) { return reviewRepository.findByTourId(tourId);}
}
