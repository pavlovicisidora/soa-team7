package com.example.tour.grpc;

import com.example.tour.grpc.TourGrpcServiceGrpc;
import com.example.tour.model.Tour;
import com.example.tour.model.TourExecution;
import com.example.tour.service.TourService;
import com.example.tour.service.TourExecutionService;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.beans.factory.annotation.Autowired;

import java.util.List;

@GrpcService
public class TourGrpcServiceImpl extends TourGrpcServiceGrpc.TourGrpcServiceImplBase {
    @Autowired
    private TourService tourService;
    
    @Autowired
    private TourExecutionService tourExecutionService;
    
    @Override
    public void createTour(CreateTourRequest request, StreamObserver<CreateTourResponse> responseObserver) {
        Tour tourToCreate = new Tour();
        tourToCreate.setName(request.getName());
        tourToCreate.setDescription(request.getDescription());
        tourToCreate.setDifficulty(request.getDifficulty());
        tourToCreate.setTags(request.getTagsList());
        tourToCreate.setAuthorId(request.getAuthorId());

        Tour createdTour = tourService.createTour(tourToCreate);

        com.example.tour.grpc.Tour grpcTour = com.example.tour.grpc.Tour.newBuilder()
                .setId(createdTour.getId())
                .setName(createdTour.getName())
                .setDescription(createdTour.getDescription())
                .setDifficulty(createdTour.getDifficulty())
                .addAllTags(createdTour.getTags())
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
                    .addAllTags(tour.getTags())
                    .setStatus(tour.getStatus().name()) 
                    .setPrice(tour.getPrice())
                    .setAuthorId(tour.getAuthorId())
                    .build();
            responseBuilder.addTours(grpcTour);
        }
        responseObserver.onNext(responseBuilder.build());
        responseObserver.onCompleted();
    }
    
    @Override
    public void getAllTours(GetAllToursRequest request, StreamObserver<GetAllToursResponse> responseObserver) {
        List<Tour> tours = tourService.findAllTours();
        GetAllToursResponse.Builder responseBuilder = GetAllToursResponse.newBuilder();
        for (Tour tour : tours) {
            com.example.tour.grpc.Tour grpcTour = com.example.tour.grpc.Tour.newBuilder()
                    .setId(tour.getId())
                    .setName(tour.getName())
                    .setDescription(tour.getDescription())
                    .setDifficulty(tour.getDifficulty())
                    .addAllTags(tour.getTags())
                    .setStatus(tour.getStatus().name())
                    .setPrice(tour.getPrice())
                    .setAuthorId(tour.getAuthorId())
                    .build();
            responseBuilder.addTours(grpcTour);
        }
        responseObserver.onNext(responseBuilder.build());
        responseObserver.onCompleted();
    }
    
    @Override
    public void getTourById(GetTourByIdRequest request, StreamObserver<GetTourByIdResponse> responseObserver) {
        Tour tour = tourService.findById(request.getTourId());
        com.example.tour.grpc.Tour grpcTour = com.example.tour.grpc.Tour.newBuilder()
                .setId(tour.getId())
                .setName(tour.getName())
                .setDescription(tour.getDescription())
                .setDifficulty(tour.getDifficulty())
                .addAllTags(tour.getTags())
                .setStatus(tour.getStatus().name())
                .setPrice(tour.getPrice())
                .setAuthorId(tour.getAuthorId())
                .build();
        GetTourByIdResponse response= GetTourByIdResponse.newBuilder().setTour(grpcTour).build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }


    private com.example.tour.grpc.TourExecution toGrpcExecution(TourExecution execution) {
        return com.example.tour.grpc.TourExecution.newBuilder()
                .setId(execution.getId())
                .setTourId(execution.getTourId())
                .setTouristId(execution.getTouristId())
                .setStatus(execution.getStatus().name())
                .build();
    }

    @Override
    public void startTour(StartTourRequest request, StreamObserver<StartTourResponse> responseObserver) {
        TourExecution createdExecution = tourExecutionService.startTour(request.getTourId(), request.getTouristId());
        com.example.tour.grpc.TourExecution grpcExecution = toGrpcExecution(createdExecution);
        StartTourResponse response = StartTourResponse.newBuilder().setTourExecution(grpcExecution).build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void abandonTour(AbandonTourRequest request, StreamObserver<AbandonTourResponse> responseObserver) {
        TourExecution updatedExecution = tourExecutionService.abandonTour(request.getTourExecutionId());
        com.example.tour.grpc.TourExecution grpcExecution = toGrpcExecution(updatedExecution);
        AbandonTourResponse response = AbandonTourResponse.newBuilder().setTourExecution(grpcExecution).build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void completeTour(CompleteTourRequest request, StreamObserver<CompleteTourResponse> responseObserver) {
        TourExecution updatedExecution = tourExecutionService.completeTour(request.getTourExecutionId());
        com.example.tour.grpc.TourExecution grpcExecution = toGrpcExecution(updatedExecution);
        CompleteTourResponse response = CompleteTourResponse.newBuilder().setTourExecution(grpcExecution).build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
}
