package com.example.tour.grpc;

import com.example.tour.grpc.TourGrpcServiceGrpc;
import com.example.tour.model.Status;
import com.example.tour.model.Tour;
import com.example.tour.model.TourExecution;
import com.example.tour.service.TourService;
import com.google.protobuf.Timestamp;
import com.example.tour.service.TourExecutionService;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.beans.factory.annotation.Autowired;

import java.time.ZoneOffset;
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
            com.example.tour.grpc.Tour.Builder grpcTourBuilder = com.example.tour.grpc.Tour.newBuilder()
                    .setId(tour.getId())
                    .setName(tour.getName())
                    .setDescription(tour.getDescription())
                    .setDifficulty(tour.getDifficulty())
                    .addAllTags(tour.getTags())
                    .setStatus(tour.getStatus().name())
                    .setPrice(tour.getPrice())
                    .setAuthorId(tour.getAuthorId())
                    .setDistanceInKm(tour.getDistanceInKm());

            if (tour.getArchivedDateTime() != null) {
                grpcTourBuilder.setArchivedDateTime(
                        Timestamp.newBuilder()
                                .setSeconds(tour.getArchivedDateTime().toEpochSecond(ZoneOffset.UTC))
                                .setNanos(tour.getArchivedDateTime().getNano())
                                .build()
                );
            }

            if (tour.getPublishedDateTime() != null) {
                grpcTourBuilder.setPublishedDateTime(
                        Timestamp.newBuilder()
                                .setSeconds(tour.getPublishedDateTime().toEpochSecond(ZoneOffset.UTC))
                                .setNanos(tour.getPublishedDateTime().getNano())
                                .build()
                );
            }

            responseBuilder.addTours(grpcTourBuilder.build());
        }

        responseObserver.onNext(responseBuilder.build());
        responseObserver.onCompleted();
    }
    
    @Override
    public void getAllTours(GetAllToursRequest request, StreamObserver<GetAllToursResponse> responseObserver) {
        List<Tour> tours = tourService.findAllTours();
        GetAllToursResponse.Builder responseBuilder = GetAllToursResponse.newBuilder();

        for (Tour tour : tours) {
            if(tour.getStatus() == Status.PUBLISHED) {
                com.example.tour.grpc.Tour.Builder grpcTourBuilder = com.example.tour.grpc.Tour.newBuilder()
                        .setId(tour.getId())
                        .setName(tour.getName())
                        .setDescription(tour.getDescription())
                        .setDifficulty(tour.getDifficulty())
                        .addAllTags(tour.getTags())
                        .setStatus(tour.getStatus().name())
                        .setPrice(tour.getPrice())
                        .setAuthorId(tour.getAuthorId())
                        .setDistanceInKm(tour.getDistanceInKm());

                if (tour.getArchivedDateTime() != null) {
                    grpcTourBuilder.setArchivedDateTime(
                            Timestamp.newBuilder()
                                    .setSeconds(tour.getArchivedDateTime().toEpochSecond(ZoneOffset.UTC))
                                    .setNanos(tour.getArchivedDateTime().getNano())
                                    .build()
                    );
                }

                if (tour.getPublishedDateTime() != null) {
                    grpcTourBuilder.setPublishedDateTime(
                            Timestamp.newBuilder()
                                    .setSeconds(tour.getPublishedDateTime().toEpochSecond(ZoneOffset.UTC))
                                    .setNanos(tour.getPublishedDateTime().getNano())
                                    .build()
                    );
                }

                responseBuilder.addTours(grpcTourBuilder.build());
            }
        }
        responseObserver.onNext(responseBuilder.build());
        responseObserver.onCompleted();
    }
    
    @Override
    public void getTourById(GetTourByIdRequest request, StreamObserver<GetTourByIdResponse> responseObserver) {
        Tour tour = tourService.findById(request.getTourId());

        com.example.tour.grpc.Tour.Builder grpcTourBuilder = com.example.tour.grpc.Tour.newBuilder()
                .setId(tour.getId())
                .setName(tour.getName())
                .setDescription(tour.getDescription())
                .setDifficulty(tour.getDifficulty())
                .addAllTags(tour.getTags())
                .setStatus(tour.getStatus().name())
                .setPrice(tour.getPrice())
                .setAuthorId(tour.getAuthorId())
                .setDistanceInKm(tour.getDistanceInKm());

        if (tour.getArchivedDateTime() != null) {
            grpcTourBuilder.setArchivedDateTime(
                    Timestamp.newBuilder()
                            .setSeconds(tour.getArchivedDateTime().toEpochSecond(ZoneOffset.UTC))
                            .setNanos(tour.getArchivedDateTime().getNano())
                            .build()
            );
        }

        if (tour.getPublishedDateTime() != null) {
            grpcTourBuilder.setPublishedDateTime(
                    Timestamp.newBuilder()
                            .setSeconds(tour.getPublishedDateTime().toEpochSecond(ZoneOffset.UTC))
                            .setNanos(tour.getPublishedDateTime().getNano())
                            .build()
            );
        }

        GetTourByIdResponse response = GetTourByIdResponse.newBuilder()
                .setTour(grpcTourBuilder.build())
                .build();

        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }


    private com.example.tour.grpc.TourExecution toGrpcExecution(TourExecution execution) {
        Timestamp startTime = Timestamp.newBuilder()
                .setSeconds(execution.getStartTime().toEpochSecond(ZoneOffset.UTC))
                .build();
        
        Timestamp lastActivity = Timestamp.newBuilder()
                .setSeconds(execution.getLastActivity().toEpochSecond(ZoneOffset.UTC))
                .build();

        return com.example.tour.grpc.TourExecution.newBuilder()
                .setId(execution.getId())
                .setTourId(execution.getTourId())
                .setTouristId(execution.getTouristId())
                .setStatus(execution.getStatus().name())
                .setStartTime(startTime)
                .setLastActivity(lastActivity)
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

    @Override
    public void getTourExecution(GetTourExecutionRequest request, StreamObserver<GetTourExecutionResponse> responseObserver) {
        TourExecution execution = tourExecutionService.getTourExecution(request.getTourExecutionId());
        
        Tour tour = tourService.findById(execution.getTourId());
        
        com.example.tour.grpc.TourExecution grpcExecution = toGrpcExecution(execution);
        com.example.tour.grpc.Tour grpcTour = toGrpcTour(tour); 

        GetTourExecutionResponse response = GetTourExecutionResponse.newBuilder()
                .setTourExecution(grpcExecution)
                .setTour(grpcTour)
                .build();
                
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    private com.example.tour.grpc.Tour toGrpcTour(Tour tour) {
        com.example.tour.grpc.Tour.Builder builder = com.example.tour.grpc.Tour.newBuilder()
                .setId(tour.getId())
                .setName(tour.getName())
                .setDescription(tour.getDescription())
                .setDifficulty(tour.getDifficulty())
                .addAllTags(tour.getTags())
                .setStatus(tour.getStatus().name())
                .setPrice(tour.getPrice())
                .setAuthorId(tour.getAuthorId())
                .setDistanceInKm(tour.getDistanceInKm());
        

        if (tour.getPublishedDateTime() != null) {
            Timestamp publishedTs = Timestamp.newBuilder()
                    .setSeconds(tour.getPublishedDateTime().toEpochSecond(ZoneOffset.UTC))
                    .build();
            builder.setPublishedDateTime(publishedTs);
        }

        if (tour.getArchivedDateTime() != null) {
            Timestamp archivedTs = Timestamp.newBuilder()
                    .setSeconds(tour.getArchivedDateTime().toEpochSecond(ZoneOffset.UTC))
                    .build();
            builder.setArchivedDateTime(archivedTs);
        }

        return builder.build();
    }

    @Override
    public void updateTour(UpdateTourRequest request, StreamObserver<UpdateTourResponse> responseObserver) {

        com.example.tour.grpc.Tour grpcTour = request.getTour();

        Tour tourToUpdate = new Tour();
        tourToUpdate.setId(grpcTour.getId());
        tourToUpdate.setName(grpcTour.getName());
        tourToUpdate.setDescription(grpcTour.getDescription());
        tourToUpdate.setDifficulty(grpcTour.getDifficulty());
        tourToUpdate.setTags(grpcTour.getTagsList());
        tourToUpdate.setStatus(com.example.tour.model.Status.valueOf(grpcTour.getStatus()));
        tourToUpdate.setPrice(grpcTour.getPrice());
        tourToUpdate.setAuthorId(grpcTour.getAuthorId());
        tourToUpdate.setDistanceInKm(grpcTour.getDistanceInKm());

        if (grpcTour.hasPublishedDateTime()) {
            Timestamp ts = grpcTour.getPublishedDateTime();
            tourToUpdate.setPublishedDateTime(
                    java.time.LocalDateTime.ofInstant(java.time.Instant.ofEpochSecond(ts.getSeconds(), ts.getNanos()), ZoneOffset.UTC)
            );
        }
        if (grpcTour.hasArchivedDateTime()) {
            Timestamp ts = grpcTour.getArchivedDateTime();
            tourToUpdate.setArchivedDateTime(
                    java.time.LocalDateTime.ofInstant(java.time.Instant.ofEpochSecond(ts.getSeconds(), ts.getNanos()), ZoneOffset.UTC)
            );
        }


        Tour updatedTourFromDb = tourService.updateTour(tourToUpdate);


        com.example.tour.grpc.Tour responseGrpcTour = toGrpcTour(updatedTourFromDb);


        UpdateTourResponse response = UpdateTourResponse.newBuilder()
                .setTour(responseGrpcTour)
                .build();

        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
}
