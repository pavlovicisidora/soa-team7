package com.example.tour.grpc;

import com.example.tour.model.Review;
import com.example.tour.model.Tour;
import com.example.tour.service.ReviewService;
import com.example.tour.service.TourService;
import com.google.protobuf.Timestamp;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;
import com.example.tour.grpc.ReviewGrpcServiceGrpc;
import org.springframework.beans.factory.annotation.Autowired;

import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.List;

@GrpcService
public class ReviewGrpcServiceImpl extends ReviewGrpcServiceGrpc.ReviewGrpcServiceImplBase {
    @Autowired
    private ReviewService reviewService;
    @Autowired
    private TourService tourService;
    @Override
    public void createReview(CreateReviewRequest request, StreamObserver<ReviewResponse> responseObserver){
        Tour tour = tourService.findById(request.getTourId());
        Review review = new Review();
        review.setRating(request.getRating());
        review.setComment(request.getComment());
        review.setTouristId(request.getTouristId());
        review.setVisitDate(LocalDate.parse(request.getVisitDate()));
        review.setCreatedDate(LocalDateTime.now());
        review.setImages(request.getImagesList());
        review.setTour(tour);

        Review createdReview =  reviewService.createReview(review);
        com.example.tour.grpc.Review.Builder grpcReview = com.example.tour.grpc.Review.newBuilder();
        grpcReview.setId(createdReview.getId());
        grpcReview.setRating(createdReview.getRating());
        grpcReview.setComment(createdReview.getComment());
        grpcReview.setTouristId(createdReview.getTouristId());
        grpcReview.setVisitDate(createdReview.getVisitDate().toString());
        if (createdReview.getCreatedDate() != null) {
            Instant instant = createdReview.getCreatedDate().atZone(ZoneId.systemDefault()).toInstant();
            Timestamp timestamp = Timestamp.newBuilder()
                    .setSeconds(instant.getEpochSecond())
                    .setNanos(instant.getNano())
                    .build();
            grpcReview.setCreatedDate(timestamp);
        }
        grpcReview.addAllImages(createdReview.getImages());
        grpcReview.setTourId(createdReview.getTour().getId());
        grpcReview.build();
        ReviewResponse response = ReviewResponse.newBuilder()
                .setReview(grpcReview)
                .build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
    @Override
    public void getReviewsForTour(GetReviewForTourRequest request, StreamObserver<ReviewsResponse> responseObserver){
        List<Review> reviews = reviewService.getAllReviewForTour(request.getTourId());
        ReviewsResponse.Builder responseBuilder = ReviewsResponse.newBuilder();
        for(Review review : reviews){
            com.example.tour.grpc.Review.Builder grpcReview = com.example.tour.grpc.Review.newBuilder()
                    .setId(review.getId())
                    .setRating(review.getRating())
                    .setComment(review.getComment())
                    .setTouristId(review.getTouristId())
                    .setVisitDate(review.getVisitDate().toString());
            if (review.getCreatedDate() != null) {
                Instant instant = review.getCreatedDate().atZone(ZoneId.systemDefault()).toInstant();
                Timestamp timestamp = Timestamp.newBuilder()
                        .setSeconds(instant.getEpochSecond())
                        .setNanos(instant.getNano())
                        .build();
                grpcReview.setCreatedDate(timestamp);
            }
            grpcReview.addAllImages(review.getImages());
            grpcReview.setTourId(review.getTour().getId());
            grpcReview.build();
            responseBuilder.addReviews(grpcReview);
        }
        responseObserver.onNext(responseBuilder.build());
        responseObserver.onCompleted();
    }
}
