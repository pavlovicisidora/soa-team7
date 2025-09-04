package com.example.tour.grpc;
import com.example.tour.grpc.TourGrpcServiceGrpc;
import com.example.tour.grpc.CreateTourRequest;
import com.example.tour.grpc.CreateTourResponse;
import io.grpc.stub.StreamObserver;
import com.example.tour.model.Tour;
import com.example.tour.service.TourService;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.beans.factory.annotation.Autowired;

import java.util.List;

@GrpcService
public class TourGrpcServiceImpl extends TourGrpcServiceGrpc.TourGrpcServiceImplBase {
    @Autowired
    private TourService tourService; // Vaša postojeća poslovna logika

    @Override
    public void createTour(CreateTourRequest request, StreamObserver<CreateTourResponse> responseObserver) {

        Tour tourToCreate = new Tour();
        tourToCreate.setName(request.getName());
        tourToCreate.setDescription(request.getDescription());
        tourToCreate.setDifficulty(request.getDifficulty());
        tourToCreate.setTags(request.getTags());
        tourToCreate.setAuthorId(request.getAuthorId());

        Tour createdTour = tourService.createTour(tourToCreate);

        com.example.tour.grpc.Tour grpcTour = com.example.tour.grpc.Tour.newBuilder()
                .setId(createdTour.getId())
                .setName(createdTour.getName())
                .setDescription(createdTour.getDescription())
                .setDifficulty(createdTour.getDifficulty())
                .setTags(createdTour.getTags())
                .setStatus(createdTour.getStatus().name())
                .setPrice(createdTour.getPrice())
                .setAuthorId(createdTour.getAuthorId())
                .build();

        CreateTourResponse response = CreateTourResponse.newBuilder().setTour(grpcTour).build();

        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
    @Override
    public void getAllToursById(GetAllToursByIdRequest request, StreamObserver<GetAllToursByIdResponse> responseObserver) {

        List<Tour> tours = tourService.findAllToursById(request.getAuthorId());
        GetAllToursByIdResponse.Builder responseBuilder = GetAllToursByIdResponse.newBuilder();
        for (Tour tour : tours) {
            com.example.tour.grpc.Tour grpcTour = com.example.tour.grpc.Tour.newBuilder()
                    .setId(tour.getId())
                    .setName(tour.getName())
                    .setDescription(tour.getDescription())
                    .setDifficulty(tour.getDifficulty())
                    .setTags(tour.getTags())
                    .setStatus(tour.getStatus().name())
                    .setPrice(tour.getPrice())
                    .setAuthorId(tour.getAuthorId())
                    .build();
            responseBuilder.addTours(grpcTour);
        }
        responseObserver.onNext(responseBuilder.build());
        responseObserver.onCompleted();
    }

}
