package com.example.tour.grpc;

import com.example.tour.service.KeyPointService;
import com.example.tour.model.KeyPoint;
import com.example.tour.model.Tour;
import com.example.tour.service.TourService;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.beans.factory.annotation.Autowired;

import java.util.List;
import java.util.Optional;
@GrpcService
public class KeyPointGrpcServiceImpl extends KeyPointGrpcServiceGrpc.KeyPointGrpcServiceImplBase {
    @Autowired
    private KeyPointService keyPointService;
    @Autowired
    private TourService tourService;
    @Override
    public void createKeyPoint(CreateKeyPointRequest request, StreamObserver<CreateKeyPointResponse> responseObserver) {
        Tour tour = tourService.findById(request.getTourId());


        KeyPoint keyPointToCreate = new KeyPoint();
        keyPointToCreate.setName(request.getName());
        keyPointToCreate.setDescription(request.getDescription());
        keyPointToCreate.setLatitude(request.getLatitude());
        keyPointToCreate.setLongitude(request.getLongitude());
        keyPointToCreate.setImageUrl(request.getImageUrl());
        keyPointToCreate.setTour(tour);
        KeyPoint createdKeyPoint = keyPointService.createKeyPoint(keyPointToCreate);
        com.example.tour.grpc.KeyPoint grpcKeyPoint = com.example.tour.grpc.KeyPoint.newBuilder()
                .setId(createdKeyPoint.getId())
                .setTourId(createdKeyPoint.getTour().getId())
                .setName(createdKeyPoint.getName())
                .setDescription(createdKeyPoint.getDescription())
                .setLatitude(createdKeyPoint.getLatitude())
                .setLongitude(createdKeyPoint.getLongitude())
                .setImageUrl(createdKeyPoint.getImageUrl())
                .build();
        CreateKeyPointResponse response = CreateKeyPointResponse.newBuilder()
                .setKeypoint(grpcKeyPoint)
                .build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
    @Override
    public void getKeyPointsForTour(GetKeyPointsForTourRequest request, StreamObserver<GetKeyPointsForTourResponse> responseObserver) {

        List<KeyPoint> keyPoints = keyPointService.findByTourId(request.getTourId());
        GetKeyPointsForTourResponse.Builder responseBuilder = GetKeyPointsForTourResponse.newBuilder();
        for (KeyPoint keyPoint : keyPoints) {
            com.example.tour.grpc.KeyPoint grpcKeyPoint = com.example.tour.grpc.KeyPoint.newBuilder()
                    .setId(keyPoint.getId())
                    .setName(keyPoint.getName())
                    .setDescription(keyPoint.getDescription())
                    .setLatitude(keyPoint.getLatitude())
                    .setLongitude(keyPoint.getLongitude())
                    .setImageUrl(keyPoint.getImageUrl())
                    .setTourId(keyPoint.getTour().getId())
                    .build();
            responseBuilder.addKeyPoints(grpcKeyPoint);
        }
        responseObserver.onNext(responseBuilder.build());
        responseObserver.onCompleted();
    }
}
