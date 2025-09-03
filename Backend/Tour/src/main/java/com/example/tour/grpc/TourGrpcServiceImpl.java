package com.example.tour.grpc;
import com.example.tour.grpc.TourGrpcServiceGrpc;
import com.example.tour.grpc.CreateTourRequest;
import com.example.tour.grpc.CreateTourResponse;
import io.grpc.stub.StreamObserver;
import com.example.tour.model.Tour;
import com.example.tour.service.TourService;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.beans.factory.annotation.Autowired;
@GrpcService
public class TourGrpcServiceImpl extends TourGrpcServiceGrpc.TourGrpcServiceImplBase {
    @Autowired
    private TourService tourService; // Vaša postojeća poslovna logika

    @Override
    public void createTour(CreateTourRequest request, StreamObserver<CreateTourResponse> responseObserver) {

        // Kreiramo naš JPA entitet iz gRPC zahteva
        Tour tourToCreate = new Tour();
        tourToCreate.setName(request.getName());
        tourToCreate.setDescription(request.getDescription());
        tourToCreate.setDifficulty(request.getDifficulty());
        tourToCreate.setTags(request.getTags());
        // Postavljamo authorId koji smo dobili od Gateway-a
        tourToCreate.setAuthorId(request.getAuthorId());

        // Pozivamo vaš postojeći servisni sloj koji čuva turu u bazi.
        // On će automatski postaviti status i cenu preko @PrePersist.
        Tour createdTour = tourService.createTour(tourToCreate);

        // Mapiramo sačuvani entitet nazad u gRPC odgovor
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
}
